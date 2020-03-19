package handlers

import (
	"fmt"
	"net/http"
)

type RootHandlerFunc struct {}

func (*RootHandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}
