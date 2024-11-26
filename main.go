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
	router.Use(middleware.Logger)

	setupUser(pool, router)

	http.ListenAndServe(conf.LoadServerAddress(), router)
}

func setupUser(pool *pgxpool.Pool, router *chi.Mux) {
	userRepo := pg.NewUser(pool)
	userSvc := service.NewUser(userRepo)
	rest.NewUserHandler(userSvc).Register(router)
}
