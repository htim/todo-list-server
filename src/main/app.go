package main

import (
	"repository"
	"resource"
	"log"
	"github.com/pressly/chi"
	"net/http"
	"os"
)

type App struct {
	Router *chi.Mux
}

func (a *App) init() {
	dbx, err := repository.NewDB(os.Getenv("POSTGRES_URL"))
	if err != nil {
		log.Fatal(err)
	}

	a.Router = chi.NewRouter()

	categoryRepository := repository.NewCategoryRepository(dbx)
	categoryResource := resource.NewCategoryResource(categoryRepository)

	a.Router.Mount("/category", categoryResource.Routes())
}

func (a *App) run() {
	http.ListenAndServe(":" + os.Getenv("PORT"), a.Router)
}
