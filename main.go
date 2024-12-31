package main

import (
	"net/http"

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
	pool, err := pg.NewPostgreSQL(conf.LoadPgConf())
	if err != nil {
		panic(err)
	}

	router := chi.NewRouter()
	router.Use(middleware.RequestID, middleware.Logger, middleware.Recoverer)

	userSvc := setupUser(pool, router)
	setupWorkout(pool, router, userSvc)
	setupExerciseViewService(router)
	loginViewService := setupLoginViewService(router, userSvc)
	setupAppViewHandler(router, loginViewService)

	fs := http.FileServer(http.Dir("static/"))
	router.Handle("/static/*", http.StripPrefix("/static/", fs))

	http.ListenAndServe(conf.LoadServerAddress(), router)
}

func setupUser(pool *pgxpool.Pool, router *chi.Mux) is.IUserService {
	userRepo := pg.NewUser(pool)
	userSvc := is.NewUserService(userRepo)
	rh.NewUserHandler(userSvc).Register(router)
	return userSvc
}

func setupWorkout(pool *pgxpool.Pool, router *chi.Mux, userSvc is.IUserService) is.IWorkoutService {
	workoutRepo := pg.NewWorkout(pool)
	workoutSvc := is.NewWorkoutService(workoutRepo, userSvc)
	rh.NewWorkoutHandler(workoutSvc).Register(router)
	return workoutSvc
}

func setupExerciseViewService(router *chi.Mux) vs.IExerciseViewService {
	exerciseViewService := vs.NewExerciseViewService()
	vh.NewExerciseViewHandler(exerciseViewService).Register(router)
	return exerciseViewService
}

func setupLoginViewService(router *chi.Mux, userService is.IUserService) vs.ILoginViewService {
	loginViewService := vs.NewLoginViewService(userService)
	vh.NewLoginViewHandler(loginViewService).Register(router)
	return loginViewService
}

func setupAppViewHandler(router *chi.Mux, loginViewService vs.ILoginViewService) {
	vh.NewAppViewHandler(loginViewService).Register(router)
}
