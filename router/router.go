package router

import (
	"net/http"
)

func init() {
	http.HandleFunc("/sse/log", seeLog)
}
