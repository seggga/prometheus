package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/seggga/prometheus/pkg/storage"
	"github.com/seggga/prometheus/pkg/storage/model"
	"go.uber.org/zap"
)

func NewLink(stor storage.CropURLStorage, slogger *zap.SugaredLogger) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		// obtain request body
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			err = fmt.Errorf("unable to parse request body %w", err)
			slogger.Debug(err)
			JSONError(rw, err, http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// compose dataset
		linkPair := new(model.LinkData)
		err = json.Unmarshal(body, linkPair)
		if err != nil {
			err = fmt.Errorf("unable to unmarshal JSON %w", err)
			slogger.Debug(err)
			JSONError(rw, err, http.StatusBadRequest)
			return
		}

		longURL, err := url.Parse(linkPair.LongURL)
		// check user's input: incorrect URL format
		if err != nil {
			err = fmt.Errorf("entered long URL cannot be recognized: (%s) %w", linkPair.LongURL, err)
			slogger.Debug(err)
			JSONError(rw, err, http.StatusBadRequest)
			return
		}
		// check user's input: scheme is empty
		if longURL.Scheme == "" {
			err = fmt.Errorf("protocol should be set (http:// or https:// or ...): %s %w", linkPair.LongURL, err)
			slogger.Debug(err)
			JSONError(rw, err, http.StatusBadRequest)
			return
		}
		// check user's input: host not set
		if longURL.Host == "" {
			err = fmt.Errorf("host address was not set: %s %w", linkPair.LongURL, err)
			slogger.Debug(err)
			JSONError(rw, err, http.StatusBadRequest)
			return
		}

		// check if shortID is free
		if stor.IsSet(linkPair.ShortID) {
			err = fmt.Errorf("short URL %s is already in use %w", linkPair.ShortID, err)
			slogger.Debug(err)
			JSONError(rw, err, http.StatusBadRequest)
			return
		}

		// add dataset to the database
		err = stor.AddURL(linkPair)
		if err != nil {
			err = fmt.Errorf("error creating new short-to-long pair %w", err)
			slogger.Errorw("error creating new short-to-long pair", err)
			JSONError(rw, err, http.StatusBadRequest)
			return
		}

		slogger.Infof("a new short ID added to database %+v", linkPair)
		// output data
		rw.Header().Set("Application", "CropURL")
		rw.WriteHeader(http.StatusOK)
	}
}
