package resource

import (
	"github.com/pressly/chi"
	"net/http"
	"golang.org/x/oauth2"
	"github.com/google/go-github/github"
	"log"
	"github.com/dgrijalva/jwt-go"
	"github.com/alioygur/gores"
	"env"
)

type AuthResource struct {
	env *env.Env
}


func NewAuthResource(env *env.Env) *AuthResource {
	return &AuthResource{env:env}
}

func (a *AuthResource) Router() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/github", a.GithubAuth)
	r.Get("/github-callback", a.GithubAuthCallback)
	return r
}

func (a *AuthResource) GithubAuth(w http.ResponseWriter, r *http.Request) {
	url := a.env.OAuth2GithubConf.AuthCodeURL(a.env.GithubStateStr, oauth2.AccessTypeOnline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (a *AuthResource) GithubAuthCallback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	if state != a.env.GithubStateStr {
		log.Println("Invalid state: " + state)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	code := r.FormValue("code")
	ctx := r.Context()
	token, err := a.env.OAuth2GithubConf.Exchange(ctx, code)
	if err != nil {
		log.Println("error while getting token: ", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	oauthClient := a.env.OAuth2GithubConf.Client(ctx, token)
	client := github.NewClient(oauthClient)
	user, _, err := client.Users.Get(ctx, "")
	if err != nil {
		log.Println("Error while getting user info", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	log.Println("Logged: ", user)

	claims := jwt.MapClaims{
		"login" : user.Login,
		"authprovider":"github",
		"id":user.ID,
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := jwtToken.SignedString(a.env.JwtSigningKey)
	if err != nil {
		log.Println(err)
	}
	gores.JSON(w, http.StatusOK, tokenString)
}