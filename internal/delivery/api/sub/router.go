package sub

import (
	"github.com/go-chi/chi/v5"
	_ "github.com/maYkiss56/subscription-aggregation-service/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(subs *HandlerSub) chi.Router {
	r := chi.NewRouter()

	// Swagger
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	r.Route("/api/subs", func(r chi.Router) {

		r.Get("/", subs.GetAllSubs)
		r.Get("/{user_id}", subs.GetSubByUserID)
		r.Post("/total", subs.CalculateTotalCost)
		r.Post("/create", subs.CreateSub)
		r.Patch("/update/{id}", subs.UpdateSub)
		r.Delete("/delete/{id}", subs.DeleteSub)
	})

	return r
}
