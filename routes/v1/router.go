package v1

import (
	"github.com/go-chi/chi"
)

func Router(r chi.Router) {
	r.Route("/users", userRouter)
}
