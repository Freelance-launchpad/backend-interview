package httpauth

import (
	"context"
	"net/url"
	"os"
	"sync"

	"github.com/Freelance-launchpad/backend-interview/common/jhttp/v2"
)

type Auth struct {
	sync.RWMutex

	clientID     string
	clientSecret string
	token        authToken

	skipAuth bool
}

func NewAuth(appName, appVersion, clientSecret string, authURL url.URL) jhttp.BearerAuthentication {
	return &Auth{
		clientID:     appName,
		clientSecret: clientSecret,
		skipAuth:     os.Getenv("SKIP_M2M_AUTH") == "true",
	}
}

func (auth *Auth) GetBearer(ctx context.Context) (token string, err error) {
	return "", nil
}
