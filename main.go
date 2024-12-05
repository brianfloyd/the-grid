package main

import (
	"net/http"

	"github.com/brianfloyd/the-grid/internal/conf"
	"github.com/brianfloyd/the-grid/internal/db/pg"
	"github.com/brianfloyd/the-grid/internal/rest"
	"github.com/brianfloyd/the-grid/internal/service"
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

	fs := http.FileServer(http.Dir("static/"))
	router.Handle("/static/*", http.StripPrefix("/static/", fs))
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})

	http.ListenAndServe(conf.LoadServerAddress(), router)
}

func setupUser(pool *pgxpool.Pool, router *chi.Mux) service.UserService {
	userRepo := pg.NewUser(pool)
	userSvc := service.NewUser(userRepo)
	rest.NewUserHandler(userSvc).Register(router)
	return userSvc
}

func setupWorkout(pool *pgxpool.Pool, router *chi.Mux, userSvc service.UserService) service.WorkoutService {
	workoutRepo := pg.NewWorkout(pool)
	workoutSvc := service.NewWorkout(workoutRepo, userSvc)
	rest.NewWorkoutHandler(workoutSvc).Register(router)
	return workoutSvc
}
