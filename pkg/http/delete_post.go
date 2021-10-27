package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/seggga/prometheus/pkg/storage"
	"github.com/seggga/prometheus/pkg/storage/model"
	"go.uber.org/zap"
)

func Delete(stor storage.CropURLStorage, slogger *zap.SugaredLogger) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		// obtain request body
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			slogger.Errorf("Unable to parse request body", err)
			JSONError(rw, err, http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// decode entered data in a structure
		linkPair := new(model.LinkData)
		err = json.Unmarshal(body, linkPair)
		if err != nil {
			slogger.Errorw("Unable to unmarshal JSON", err)
			JSONError(rw, err, http.StatusBadRequest)
			return
		}

		err = stor.Delete(linkPair.ShortID)
		if err != nil {
			slogger.Errorw("error deleting short-to-long pair", err)
			JSONError(rw, err, http.StatusBadRequest)
			return
		}

		slogger.Infof("short ID %s was deleted from database", linkPair.ShortID)
		// output data
		rw.Header().Set("Application", "CropURL")
		rw.WriteHeader(http.StatusOK)
	}
}
