package app

import (
	"net/http"
)

func (app *App) NewRouter() http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("GET /api/v1/smartphones", app.GetSmartphones)
	router.HandleFunc("GET /api/v1/smartphones/{smartphone_id}", app.GetSmartphone)

	router.HandleFunc("GET /api/v1/users", app.Auth(app.GetUsers))
	router.HandleFunc("GET /api/v1/users/{user_id}", app.Auth(app.GetUser))
	router.HandleFunc("POST /api/v1/login", app.Login)
	router.HandleFunc("POST /api/v1/signup", app.Signup)

	router.HandleFunc("GET /api/v1/smartphones/{smartphone_id}/reviews", app.GetReviews)
	router.HandleFunc("GET /api/v1/smartphones/{smartphone_id}/reviews/{review_id}", app.GetReview)
	router.HandleFunc("POST /api/v1/smartphones/{smartphone_id}/reviews", app.Auth(app.CreateReview))
	router.HandleFunc("PATCH /api/v1/smartphones/{smartphone_id}/reviews/{review_id}", app.Auth(app.UpdateReview))
	router.HandleFunc("DELETE /api/v1/smartphones/{smartphone_id}/reviews/{review_id}", app.Auth(app.DeleteReview))

	router.HandleFunc("GET /api/v1/carts", app.Auth(app.GetCarts))
	router.HandleFunc("GET /api/v1/carts/{cart_id}", app.Auth(app.GetCart))

	router.HandleFunc("GET /api/v1/carts/{cart_id}/items", app.Auth(app.GetCartItems))
	router.HandleFunc("POST /api/v1/carts/{cart_id}/items", app.Auth(app.AddToCart))
	router.HandleFunc("PATCH /api/v1/carts/{cart_id}/items/{item_id}", app.Auth(app.SetQuantity))
	router.HandleFunc("DELETE /api/v1/carts/{cart_id}/items/{item_id}", app.Auth(app.DeleteFromCart))

	return app.RecoverPanic(app.LogRequests(router))
}
