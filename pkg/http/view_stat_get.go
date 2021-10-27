package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/seggga/prometheus/pkg/storage"
	"go.uber.org/zap"
)

func ViewStatistics(stor storage.CropURLStorage, slogger *zap.SugaredLogger) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		// define shortID from users query
		shortID := chi.URLParam(r, "shortID")
		// define data about the specified short ID
		urlData, err := stor.ViewStat(shortID)
		if err != nil {
			slogger.Debugf("get statistics error %w", err)
			JSONError(rw, err, http.StatusBadRequest)
			return
		}
		slogger.Debugf("statistics query on %s", shortID)
		// output in json
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		err = json.NewEncoder(rw).Encode(urlData)
		if err != nil {
			slogger.Errorf("statistics query on %w", err)
		}
	}
}
