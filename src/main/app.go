package main

import (
	"repository"
	"resource"
	"log"
	"github.com/pressly/chi"
	"net/http"
	"strings"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"env"
	"github.com/pressly/chi/middleware"
	"github.com/unrolled/render"
	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
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
	dbx.MapperFunc(strings.ToLower)

	en := env.NewEnv()
	en.GithubConfig = initOauth2GithubConfig()
	en.JwtSigningKey = []byte(os.Getenv("JWT_SIGNING_KEY"))
	en.Render = initRenderJSON()

	categoryRepository := repository.NewCategoryRepository(dbx)
	categoryResource := resource.NewCategoryResource(categoryRepository, en)

	authResource := resource.NewAuthResource(en)

	jwtMiddleware := initJwtMiddleware(en).Handler

	a.Router = chi.NewRouter()

	a.Router.Use(middleware.RequestID)
	a.Router.Use(middleware.RealIP)
	a.Router.Use(middleware.Logger)
	a.Router.Use(middleware.Recoverer)

	categoryRouter := categoryResource.Router(jwtMiddleware)

	a.Router.Mount("/category", categoryRouter)
	a.Router.Mount("/auth", authResource.Router())
}

func (a *App) run() {
	log.Fatal(http.ListenAndServe(":3333", a.Router))
}

func initOauth2GithubConfig() *env.GithubConfig {
	conf := &oauth2.Config{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		Scopes:       []string{"user"},
		Endpoint:     github.Endpoint,
	}
	state := "githubStateStr"
	return &env.GithubConfig{
		OAuth2GithubConf:conf,
		GithubStateStr:state,
	}
}

func initRenderJSON() *render.Render {
	r := render.New()
	return r
}

func initJwtMiddleware(env *env.Env) *jwtmiddleware.JWTMiddleware {
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(*jwt.Token) (interface{}, error) {
			return env.JwtSigningKey, nil
		},
		SigningMethod:jwt.SigningMethodHS256,
	})
	return jwtMiddleware
}