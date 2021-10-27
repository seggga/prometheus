package handler

import (
	"encoding/json"
	"net/http"
)

func JSONError(rw http.ResponseWriter, err error, code int) {
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.Header().Set("Accept", "application/json")
	rw.Header().Set("X-Content-Type-Options", "nosniff")
	rw.WriteHeader(code)
	_ = json.NewEncoder(rw).Encode(err)
}
