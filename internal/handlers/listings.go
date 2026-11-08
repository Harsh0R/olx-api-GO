package handlers

import (
	"database/sql"
	"net/http"
)

func List(db *sql.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("Lists okkk"))
	}
}
