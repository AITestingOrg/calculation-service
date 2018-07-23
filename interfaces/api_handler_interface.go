package interfaces

import "net/http"

type ApiHandlerInterface interface {
	Handle(w http.ResponseWriter, r *http.Request)
	GetPath() string
	GetRequestType() string
}
