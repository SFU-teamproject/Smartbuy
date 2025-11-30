package app

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/sfu-teamproject/smartbuy/backend/apperrors"
	"github.com/sfu-teamproject/smartbuy/backend/models"
)

// @Summary      Get a Smartphone
// @Description  Get a specific smartphone by ID
// @Tags         smartphones
// @Produce      json
// @Param        id  path int true "Smartphone ID"
// @Success      200  {object}   models.Smartphone
// @Failure      400  {object}  apperrors.ErrorResponse "Bad Request"
// @Router       /smartphones/{smartphone_id} [get]
func (app *App) GetSmartphone(w http.ResponseWriter, r *http.Request) {
	smartphoneID, err := app.ExtractPathValue(r, "smartphone_id")
	if err != nil {
		app.ErrorJSON(w, r, err)
		return
	}
	sm, err := app.DB.GetSmartphone(smartphoneID)
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("error getting smartphone %d: %w", smartphoneID, err))
		return
	}
	reviews, err := app.DB.GetReviews(sm.ID)
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("error getting reviews for smartphone %d: %w", smartphoneID, err))
		return
	}
	sm.Reviews = reviews
	app.Encode(w, r, sm)
}

// GetSmartphones lists smartphones
// @Summary      List Smartphones
// @Description  Get a list of all smartphones or filter by IDs
// @Tags         smartphones
// @Accept       json
// @Produce      json
// @Param        ids  query string false "Comma separated IDs (e.g. 1,2,3)"
// @Success      200  {array}   models.Smartphone
// @Failure      400  {object}  apperrors.ErrorResponse "Bad Request"
// @Router       /smartphones [get]
func (app *App) GetSmartphones(w http.ResponseWriter, r *http.Request) {
	var sm []models.Smartphone
	var err error
	IDsParam := r.URL.Query().Get("ids")
	if IDsParam == "" {
		sm, err = app.DB.GetSmartphones()
	} else {
		IDsStr := strings.Split(IDsParam, ",")
		IDs := make([]int, len(IDsStr))
		for i, IDStr := range IDsStr {
			ID, err := strconv.Atoi(IDStr)
			if err != nil {
				app.ErrorJSON(w, r, fmt.Errorf("%w: invalid id(%s): %w", apperrors.ErrBadRequest, IDStr, err))
			}
			IDs[i] = ID
		}
		sm, err = app.DB.GetSmartphonesByIDs(IDs)
	}
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("error getting smartphones: %w", err))
		return
	}
	app.Encode(w, r, sm)
}
