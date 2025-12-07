package v1

import (
	"log/slog"

	customerHandler "user-service/internal/http/v1/customer"
	customerService "user-service/internal/service/customer"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func SetupRoutes(
	r chi.Router,
	customerSvc *customerService.Service,
	log *slog.Logger,
) {
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

	customerH := customerHandler.NewHandler(log, customerSvc)

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/customers", func(r chi.Router) {
			r.Post("/", customerH.CreateCustomer)
			r.Get("/", customerH.GetAllCustomers)
			r.Get("/{id}", customerH.GetCustomer)
			r.Put("/{id}", customerH.UpdateCustomer)
		})
	})
}
