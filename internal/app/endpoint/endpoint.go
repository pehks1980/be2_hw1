package endpoint

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"log"
	"net/http"
	"pehks1980/be2_hw1/internal/pkg/model"
)

// App - application contents & methods
type App struct {
	Logger *zap.Logger
	Repository RepoIf
	CTX context.Context
}
// RepoIf - repository interface (PG)
type RepoIf interface {
	New(ctx context.Context, filename string) RepoIf
	CloseConn()
	// user
	AuthUser(ctx context.Context, user model.User) (string, error)
	GetUser(ctx context.Context, name string) (model.User, error)
	AddUpdUser(ctx context.Context, user model.User) (string, error)
	DelUser(ctx context.Context, id uuid.UUID) (error)
	GetUserEnvs(ctx context.Context, name string) (model.Envs, error)
	// env
	AddUpdEnv(ctx context.Context, env model.Environment) (string, error)
	GetEnv(ctx context.Context, title string) (model.Environment, error)
	DelEnv(ctx context.Context, id uuid.UUID) (error)
	GetEnvUsers(ctx context.Context, title string) (model.Users, error)
}
// RegisterPublicHTTP - регистрация роутинга путей типа urls.py для обработки сервером
func (app *App) RegisterPublicHTTP() *mux.Router {
	r := mux.NewRouter()
	// authorization
	r.HandleFunc("/user/auth", app.postAuth()).Methods(http.MethodPost)
	// user crud
	r.HandleFunc("/user/", app.putUser()).Methods(http.MethodPost)
	r.HandleFunc("/user/{uid}", app.getUser()).Methods(http.MethodGet)
	r.HandleFunc("/user/{uid}", app.putUser()).Methods(http.MethodPut)
	r.HandleFunc("/user/{uid}", app.delUser()).Methods(http.MethodDelete)
	// env crud
	r.HandleFunc("/env/", app.putEnv()).Methods(http.MethodPost)
	r.HandleFunc("/env/{uid}", app.getEnv()).Methods(http.MethodGet)
	r.HandleFunc("/env/{uid}", app.putEnv()).Methods(http.MethodPut)
	r.HandleFunc("/env/{uid}", app.delEnv()).Methods(http.MethodDelete)

	// GetUserEnvs
	r.HandleFunc("/user/envs", app.postUserEnvs()).Methods(http.MethodPost)
	// GetEnvUsers
	r.HandleFunc("/env/users", app.postEnvUsers()).Methods(http.MethodPost)

	//...
	return r
}
// write response in json format
func writeResponse(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	_, _ = w.Write([]byte(message))
	_, _ = w.Write([]byte("\n"))
}
// write response in json format
func writeJsonResponse(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, fmt.Sprintf("can't marshal data: %s", err))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	writeResponse(w, status, string(response))
}
// postAuth - user auth method
func (app *App) postAuth() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, request *http.Request) {

		ctx := request.Context()

		defer func() {
			// update Prom objects AuthCounter tries
		}()

		//json header check
		contentType := request.Header.Get("Content-Type")
		if contentType != "application/json" {
			return
		}

		user := model.User{}

		err := json.NewDecoder(request.Body).Decode(&user)
		if err != nil {
			return
		}

		UID, err1 := app.Repository.AuthUser(ctx, user)

		if err1 != nil || UID == "" {
			log.Printf("USER %s Log in error.\n", user.Name)
			return
		}
		log.Printf("USER %s Logged in.\n", user.Name)

		writeJsonResponse(w, http.StatusOK, UID)

	}
}
// putUser - add or update user
func (app *App) putUser() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		defer func() {
			// update Prom objects
		}()
		user := model.User{}
		id, _ := app.Repository.AddUpdUser(ctx, user)
		writeJsonResponse(w, http.StatusOK, id)
	}
}
// delUser - delete user
func (app *App) delUser() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		defer func() {
			// update Prom objects
		}()
		var id uuid.UUID
		_ = app.Repository.DelUser(ctx, id)
		writeJsonResponse(w, http.StatusOK, "")
	}
}
// getUser - get user
func (app *App) getUser() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		defer func() {
			// update Prom objects
		}()
		var (
			name string
			user model.User
		)
		user, _ = app.Repository.GetUser(ctx, name)
		log.Printf("getUser = %v", user)
		writeJsonResponse(w, http.StatusOK, "")
	}
}
// putEnv - add or update environment
func (app *App) putEnv() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, request *http.Request) {
		//ctx := request.Context()
		defer func() {
			// update Prom objects
		}()
		writeJsonResponse(w, http.StatusOK, "")
	}
}
// getEnv - get Env
func (app *App) getEnv() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, request *http.Request) {
		//ctx := request.Context()
		defer func() {
			// update Prom objects
		}()
		writeJsonResponse(w, http.StatusOK, "")
	}
}
// delEnv - delete Env
func (app *App) delEnv() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, request *http.Request) {
		//ctx := request.Context()
		defer func() {
			// update Prom objects
		}()
		writeJsonResponse(w, http.StatusOK, "")
	}
}
// postUserEnvs - get envs of which user is member of
func (app *App) postUserEnvs() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		defer func() {
			// update Prom objects
		}()
		var (
			envs model.Envs
			name string
		)
		envs, _ = app.Repository.GetUserEnvs(ctx,name)
		log.Printf("GetUserEnvs(%s) = %v", name, envs)
		writeJsonResponse(w, http.StatusOK, "")
	}
}
// postEnvUsers - get Users which have this env membership
func (app *App) postEnvUsers() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		defer func() {
			// update Prom objects
		}()
		var (
			users model.Users
			title string
		)
		users, _ = app.Repository.GetEnvUsers(ctx, title)
		log.Printf("GetEnvUsers(%s) = %v", title, users)
		writeJsonResponse(w, http.StatusOK, "")
	}
}
