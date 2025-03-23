package app

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sfu-teamproject/smartbuy/backend/models"
)

func (app *App) GetReview(w http.ResponseWriter, r *http.Request) {
	reviewID, err := app.ExtractPathValue(r, "review_id")
	if err != nil {
		app.ErrorJSON(w, r, err)
		return
	}
	review, err := app.DB.GetReview(reviewID)
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("error getting review %d: %w", reviewID, err))
		return
	}
	app.Encode(w, r, review)
}

func (app *App) GetReviews(w http.ResponseWriter, r *http.Request) {
	smartphoneID, err := app.ExtractPathValue(r, "smartphone_id")
	if err != nil {
		app.ErrorJSON(w, r, err)
		return
	}
	reviews, err := app.DB.GetReviews(smartphoneID)
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("error getting reviews(%d): %w", smartphoneID, err))
		return
	}
	app.Encode(w, r, reviews)
}

func (app *App) CreateReview(w http.ResponseWriter, r *http.Request) {
	userID, _, err := app.GetClaims(r)
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("%w: error extracting claims: %w", errUnauthorized, err))
		return
	}
	smartphoneID, err := app.ExtractPathValue(r, "smartphone_id")
	if err != nil {
		app.ErrorJSON(w, r, err)
		return
	}
	var review models.Review
	err = json.NewDecoder(r.Body).Decode(&review)
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("%w: error decoding review: %w", errBadRequest, err))
		return
	}
	review.SmartphoneID = smartphoneID
	review.UserID = userID
	newReview, err := app.DB.CreateReview(review)
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("error creating review: %w", err))
		return
	}
	w.WriteHeader(http.StatusCreated)
	app.Encode(w, r, newReview)
}

func (app *App) UpdateReview(w http.ResponseWriter, r *http.Request) {
	userID, role, err := app.GetClaims(r)
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("%w: error extracting claims: %w", errUnauthorized, err))
		return
	}
	var review models.Review
	review.ID, err = app.ExtractPathValue(r, "review_id")
	if err != nil {
		app.ErrorJSON(w, r, err)
		return
	}
	err = json.NewDecoder(r.Body).Decode(&review)
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("%w: error decoding review: %w", errBadRequest, err))
		return
	}
	existingReview, err := app.DB.GetReview(review.ID)
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("error getting review: %w", err))
		return
	}
	if existingReview.UserID != userID && role != models.RoleAdmin {
		app.ErrorJSON(w, r, errForbidden)
		return
	}
	smID, err := app.ExtractPathValue(r, "smartphone_id")
	if err != nil {
		app.ErrorJSON(w, r, err)
		return
	}
	review.SmartphoneID = smID
	if existingReview.SmartphoneID != review.SmartphoneID {
		app.ErrorJSON(w, r, fmt.Errorf("%w: review %d is not for smartphone %d",
			errBadRequest, review.ID, smID))
		return
	}
	updatedReview, err := app.DB.UpdateReview(review)
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("error updating review: %w", err))
		return
	}
	app.Encode(w, r, updatedReview)
}

func (app *App) DeleteReview(w http.ResponseWriter, r *http.Request) {
	userID, role, err := app.GetClaims(r)
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("%w: error extracting claims: %w", errUnauthorized, err))
		return
	}
	reviewID, err := app.ExtractPathValue(r, "review_id")
	if err != nil {
		app.ErrorJSON(w, r, err)
		return
	}
	existingReview, err := app.DB.GetReview(reviewID)
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("error getting review: %w", err))
		return
	}
	if existingReview.UserID != userID && role != models.RoleAdmin {
		app.ErrorJSON(w, r, errForbidden)
		return
	}
	smID, err := app.ExtractPathValue(r, "smartphone_id")
	if err != nil {
		app.ErrorJSON(w, r, err)
		return
	}
	if existingReview.SmartphoneID != smID {
		app.ErrorJSON(w, r, fmt.Errorf("%w: review %d is not for smartphone %d",
			errBadRequest, reviewID, smID))
		return
	}
	deletedReview, err := app.DB.DeleteReview(reviewID)
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("error deleting review: %w", err))
		return
	}
	app.Encode(w, r, deletedReview)
}
