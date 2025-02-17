package api

import (
	"net/http"

	"github.com/odit-bit/cloudfs/internal/blob"
	"github.com/odit-bit/cloudfs/internal/user"
	"github.com/sirupsen/logrus"
)

type handlerCategory string

const (
	_account handlerCategory = "auth-handler"
	_object  handlerCategory = "object-handler"
)

type App struct {
	// session  *scs.SessionManager
	accounts *user.Users
	objects  *blob.Blobs
	logger   *logrus.Logger
}

func New(accounts *user.Users, objects *blob.Blobs, logger *logrus.Logger) *App {
	return &App{accounts: accounts, objects: objects, logger: logger}
}

func (v *App) serviceErr(w http.ResponseWriter, _ *http.Request, handler string, err error) {
	v.logger.Errorf("handler: %v, error: %v", handler, err)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
