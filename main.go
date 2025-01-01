package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/brianfloyd/the-grid/internal/conf"
	"github.com/brianfloyd/the-grid/internal/db/pg"
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
	fmt.Printf("The grid is starting up!\n")

	fmt.Println("Intializing postgres.")
	pool, err := pg.NewPostgreSQL(conf.LoadPgConf())
	if err != nil {
		panic(err)
	}

	fmt.Println("Intializing chi.")
	router := chi.NewRouter()
	router.Use(middleware.RequestID, middleware.Logger, middleware.Recoverer)

	fmt.Println("Intializing services.")

	// Internal Services
	userService := setupUserService(pool)
	workoutService := setupWorkoutService(pool, userService)
	exerciseService := setupExercisesService(pool)

	// View Services
	loginViewService := setupLoginViewService(router, userService)
	exerciseViewService := setupExerciseViewService(router, exerciseService)

	// Rest handlers
	setupRestUserHandler(router, userService)
	setupRestWorkoutHandler(router, workoutService)
	setupRestExercisesHandler(router, exerciseService)

	// View handlers
	setupViewExercisesHandler(router, exerciseViewService)
	setupViewLoginHandler(router, loginViewService)
	setupViewAppHandler(router, loginViewService)

	fmt.Println("Intializing static file content.")
	fs := http.FileServer(http.Dir("static/"))
	router.Handle("/static/*", http.StripPrefix("/static/", fs))

	fmt.Printf("The grid is about to listen. Startup took %d ms.\n", time.Since(now).Milliseconds())
	http.ListenAndServe(conf.LoadServerAddress(), router)
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
	workoutRepo := pg.NewWorkout(pool)
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

func setupExerciseViewService(router *chi.Mux, exercisesService is.IExercisesService) vs.IExerciseViewService {
	exerciseViewService := vs.NewExerciseViewService(exercisesService)
	return exerciseViewService
}

func setupViewExercisesHandler(router *chi.Mux, exerciseViewService vs.IExerciseViewService) {
	vh.NewExerciseViewHandler(exerciseViewService).Register(router)
}

func setupLoginViewService(router *chi.Mux, userService is.IUserService) vs.ILoginViewService {
	loginViewService := vs.NewLoginViewService(userService)
	return loginViewService
}

func setupViewLoginHandler(router *chi.Mux, loginViewService vs.ILoginViewService) {
	vh.NewLoginViewHandler(loginViewService).Register(router)
}

func setupViewAppHandler(router *chi.Mux, loginViewService vs.ILoginViewService) {
	vh.NewAppViewHandler(loginViewService).Register(router)
}
