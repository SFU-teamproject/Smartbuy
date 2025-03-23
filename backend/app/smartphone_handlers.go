package app

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/sfu-teamproject/smartbuy/backend/models"
)

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
				app.ErrorJSON(w, r, fmt.Errorf("%w: invalid id(%s): %w", errBadRequest, IDStr, err))
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
