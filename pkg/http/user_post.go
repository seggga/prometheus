package handler

import (
	"fmt"
	"net/http"
)

func UserLogin() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(rw, "user login")
	}
}

func UserLogout() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(rw, "user logout")
	}
}