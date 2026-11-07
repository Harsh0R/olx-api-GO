package handlers

import "net/http"

func Health(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(`{"Status" : "ok"}`))
}
