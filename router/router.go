package router

import (
	"net/http"
)

func init() {
	http.HandleFunc("/see/log", seeLog)
}
