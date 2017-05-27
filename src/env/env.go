package env

import (
	"golang.org/x/oauth2"
	"github.com/unrolled/render"
)


type Env struct {
	*GithubConfig
	Render *render.Render
	JwtSigningKey []byte
}

func NewEnv() *Env {
	return &Env{}
}

func (e *Env) SetGithubConfig(conf *GithubConfig) {
	e.GithubConfig = conf
}

type GithubConfig struct {
	OAuth2GithubConf *oauth2.Config
	GithubStateStr   string
}