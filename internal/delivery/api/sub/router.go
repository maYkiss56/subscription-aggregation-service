package sub

import "github.com/go-chi/chi/v5"

func NewRouter(subs *HandlerSub) chi.Router {
	r := chi.NewRouter()

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
