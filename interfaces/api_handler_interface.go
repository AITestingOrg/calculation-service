package interfaces

import (
	"github.com/gorilla/mux"
)

type ApiHandlerInterface interface {
	AddHandlerToRouter(r *mux.Router)
}
