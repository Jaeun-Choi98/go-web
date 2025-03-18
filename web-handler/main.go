package main

import (
	"net/http"

	"turker.web/handler"
)

func main() {
	http.ListenAndServe(":3000", handler.NewHandler())
}
