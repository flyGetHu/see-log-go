package router

import "net/http"

func init() {
	http.Handle("/sse", broker())
}
