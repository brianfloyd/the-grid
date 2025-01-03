package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/brianfloyd/the-grid/internal/config"
	"github.com/brianfloyd/the-grid/internal/db/pg"
	"github.com/brianfloyd/the-grid/internal/logger"
	is "github.com/brianfloyd/the-grid/internal/service"
	rh "github.com/brianfloyd/the-grid/rest/handler"
	vh "github.com/brianfloyd/the-grid/view/handler"
	vs "github.com/brianfloyd/the-grid/view/service"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	now := time.Now()
	fmt.Println("The grid is starting up!")

	fmt.Println("Intializing configuration.")
	c := config.NewConfig()
	err := c.Read()
	if err != nil {
		panic("Could not read the configuration file.")
	}

	fmt.Println("Intializing logging.")
	adapter := logger.NewZeroLogAdapater(config.GetLogLevel(c))
	logger.Init(adapter)

	logger.Info(context.Background(), "Initializing postgres.")
	pool, err := pg.NewPostgreSQL(config.GetPostgresConfig(c))
	if err != nil {
		panic(err)
	}

	logger.Info(context.Background(), "Initializing chi.")
	router := chi.NewRouter()
	router.Use(middleware.RequestID, logger.Middleware, middleware.Recoverer)

	logger.Info(context.Background(), "Initializing services.")
	// Internal Services
	userService := setupUserService(pool)
	workoutService := setupWorkoutService(pool, userService)
	exerciseService := setupExercisesService(pool)

	// View Services
	loginViewService := setupLoginViewService(userService)
	exerciseViewService := setupExerciseViewService(exerciseService)

	// Rest handlers
	setupRestUserHandler(router, userService)
	setupRestWorkoutHandler(router, workoutService)
	setupRestExercisesHandler(router, exerciseService)

	// View handlers
	setupViewExercisesHandler(router, exerciseViewService)
	setupViewLoginHandler(router, loginViewService)
	setupViewAppHandler(router)

	logger.Info(context.Background(), "Intializing static file content.")
	fs := http.FileServer(http.Dir("static/"))
	router.Handle("/static/*", http.StripPrefix("/static/", fs))

	logger.InfoArgs(context.Background(), "The grid is about to listen. Startup took %d ms.", time.Since(now).Milliseconds())
	http.ListenAndServe(config.GetServerAddress(), router)
}

func setupUserService(pool *pgxpool.Pool) is.IUserService {
	userRepo := pg.NewUser(pool)
	userSvc := is.NewUserService(userRepo)
	return userSvc
}

func setupRestUserHandler(router *chi.Mux, userService is.IUserService) {
	rh.NewUserHandler(userService).Register(router)
}

func setupWorkoutService(pool *pgxpool.Pool, userSvc is.IUserService) is.IWorkoutService {
	workoutRepo := pg.NewWorkoutRepository(pool)
	workoutSvc := is.NewWorkoutService(workoutRepo, userSvc)
	return workoutSvc
}

func setupRestWorkoutHandler(router *chi.Mux, workoutService is.IWorkoutService) {
	rh.NewWorkoutHandler(workoutService).Register(router)
}

func setupExercisesService(pool *pgxpool.Pool) is.IExercisesService {
	repo := pg.NewExercisesRepository(pool)
	exercisesService := is.NewExerciseService(repo)
	return exercisesService
}

func setupRestExercisesHandler(router *chi.Mux, exercisesService is.IExercisesService) {
	rh.NewExerciseHandler(exercisesService).Register(router)
}

func setupExerciseViewService(exercisesService is.IExercisesService) vs.IExerciseViewService {
	exerciseViewService := vs.NewExerciseViewService(exercisesService)
	return exerciseViewService
}

func setupViewExercisesHandler(router *chi.Mux, exerciseViewService vs.IExerciseViewService) {
	vh.NewExerciseViewHandler(exerciseViewService).Register(router)
}

func setupLoginViewService(userService is.IUserService) vs.ILoginViewService {
	loginViewService := vs.NewLoginViewService(userService)
	return loginViewService
}

func setupViewLoginHandler(router *chi.Mux, loginViewService vs.ILoginViewService) {
	vh.NewLoginViewHandler(loginViewService).Register(router)
}

func setupViewAppHandler(router *chi.Mux) {
	vh.NewAppViewHandler().Register(router)
}
