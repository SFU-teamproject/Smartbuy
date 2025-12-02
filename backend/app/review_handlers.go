package app

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sfu-teamproject/smartbuy/backend/apperrors"
	"github.com/sfu-teamproject/smartbuy/backend/models"
)

// GetReview gets a single review
// @Summary      Get a Review
// @Tags         reviews
// @Produce      json
// @Param        smartphone_id path int true "Smartphone ID"
// @Param        review_id path int true "Review ID"
// @Success      200  {object}  models.Review
// @Failure      404  {object}  apperrors.ErrorResponse "Not Found"
// @Router       /smartphones/{smartphone_id}/reviews/{review_id} [get]
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

// GetReviews gets reviews for a phone
// @Summary      List Reviews
// @Description  Get all reviews for a specific smartphone
// @Tags         reviews
// @Produce      json
// @Param        smartphone_id path int true "Smartphone ID"
// @Success      200  {array}   models.Review
// @Failure      404  {object}  apperrors.ErrorResponse "Not Found"
// @Router       /smartphones/{smartphone_id}/reviews [get]
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

// CreateReview adds a review
// @Summary      Post a Review
// @Tags         reviews
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        smartphone_id path int true "Smartphone ID"
// @Param        input body models.ReviewRequest true "Review Body"
// @Success      201  {object}  models.Review
// @Failure      401  {object}  apperrors.ErrorResponse "Unauthorized"
// @Router       /smartphones/{smartphone_id}/reviews [post]
func (app *App) CreateReview(w http.ResponseWriter, r *http.Request) {
	userID, _, err := app.GetClaims(r)
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("%w: error extracting claims: %w", apperrors.ErrUnauthorized, err))
		return
	}
	smartphoneID, err := app.ExtractPathValue(r, "smartphone_id")
	if err != nil {
		app.ErrorJSON(w, r, err)
		return
	}
	var reviewreq models.ReviewRequest
	err = json.NewDecoder(r.Body).Decode(&reviewreq)
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("%w: error decoding review: %w", apperrors.ErrBadRequest, err))
		return
	}
	review := models.Review{Rating: reviewreq.Rating, Comment: reviewreq.Comment}
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

// UpdateReview edits a review
// @Summary      Edit a Review
// @Tags         reviews
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        smartphone_id path int true "Smartphone ID"
// @Param        review_id path int true "Review ID"
// @Param        input body models.ReviewRequest true "Review Body"
// @Success      200  {object}  models.Review
// @Failure      403  {object}  apperrors.ErrorResponse "Forbidden"
// @Router       /smartphones/{smartphone_id}/reviews/{review_id} [patch]
func (app *App) UpdateReview(w http.ResponseWriter, r *http.Request) {
	userID, role, err := app.GetClaims(r)
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("%w: error extracting claims: %w", apperrors.ErrUnauthorized, err))
		return
	}
	reviewID, err := app.ExtractPathValue(r, "review_id")
	if err != nil {
		app.ErrorJSON(w, r, err)
		return
	}
	var reviewreq models.ReviewRequest
	err = json.NewDecoder(r.Body).Decode(&reviewreq)
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("%w: error decoding review: %w", apperrors.ErrBadRequest, err))
		return
	}
	existingReview, err := app.DB.GetReview(reviewID)
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("error getting review: %w", err))
		return
	}
	if existingReview.UserID != userID && role != models.RoleAdmin {
		app.ErrorJSON(w, r, apperrors.ErrForbidden)
		return
	}
	smID, err := app.ExtractPathValue(r, "smartphone_id")
	if err != nil {
		app.ErrorJSON(w, r, err)
		return
	}
	review := models.Review{Rating: reviewreq.Rating, Comment: reviewreq.Comment, ID: reviewID}
	review.SmartphoneID = smID
	if existingReview.SmartphoneID != review.SmartphoneID {
		app.ErrorJSON(w, r, fmt.Errorf("%w: review %d is not for smartphone %d",
			apperrors.ErrBadRequest, review.ID, smID))
		return
	}
	updatedReview, err := app.DB.UpdateReview(review)
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("error updating review: %w", err))
		return
	}
	app.Encode(w, r, updatedReview)
}

// @Summary      Deletes a review
// @Description  Deletes a review of a certain smartphone
// @Tags         reviews
// @Security     BearerAuth
// @Produce      json
// @Param        smartphone_id path int true "Smartphone ID"
// @Param        review_id path int true "Review ID"
// @Success      200  {object}  models.Review
// @Failure      401  {object}  apperrors.ErrorResponse "Unauthorized"
// @Failure      403  {object}  apperrors.ErrorResponse "Forbidden"
// @Router       /smartphones/{smartphone_id}/reviews/{review_id} [delete]
func (app *App) DeleteReview(w http.ResponseWriter, r *http.Request) {
	userID, role, err := app.GetClaims(r)
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("%w: error extracting claims: %w", apperrors.ErrUnauthorized, err))
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
		app.ErrorJSON(w, r, apperrors.ErrForbidden)
		return
	}
	smID, err := app.ExtractPathValue(r, "smartphone_id")
	if err != nil {
		app.ErrorJSON(w, r, err)
		return
	}
	if existingReview.SmartphoneID != smID {
		app.ErrorJSON(w, r, fmt.Errorf("%w: review %d is not for smartphone %d",
			apperrors.ErrBadRequest, reviewID, smID))
		return
	}
	deletedReview, err := app.DB.DeleteReview(reviewID)
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("error deleting review: %w", err))
		return
	}
	app.Encode(w, r, deletedReview)
}
