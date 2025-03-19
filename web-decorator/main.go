package main

import (
	"net/http"
	"root/handler"
)

func main() {
	hdr := handler.NewHandler()
	http.ListenAndServe(":3000", hdr)
}
