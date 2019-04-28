package handlers

import (
	"net/http"
)

/* TODO: implement a CORS middleware handler, as described
in https://drstearns.github.io/tutorials/cors/ that responds
with the following headers to all requests:

  Access-Control-Allow-Origin: *
  Access-Control-Allow-Methods: GET, PUT, POST, PATCH, DELETE
  Access-Control-Allow-Headers: Content-Type, Authorization
  Access-Control-Expose-Headers: Authorization
  Access-Control-Max-Age: 600
*/

const accessControlAllowOrigin = "*"
const accessControlAllowMethods = "GET, PUT, POST, PATCH, DELETE"
const accessControlAllowHeaders = "Content-Type, Authorization"
const accessControlExposeHeaders = "Authorization"

type ResponseHeader struct {
	handler http.Handler
}

func NewResponseHeader(handler http.Handler) *ResponseHeader {
	return &ResponseHeader{handler}
}

func (rh *ResponseHeader) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", accessControlAllowOrigin)
	w.Header().Add("Access-Control-Allow-Methods", accessControlAllowMethods)
	w.Header().Add("Access-Control-Allow-Headers", accessControlAllowHeaders)
	w.Header().Add("Access-Control-Expose-Headers", accessControlExposeHeaders)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	rh.handler.ServeHTTP(w, r)
}
