package handlers

import (
	"fmt"
	"net/http"
	"time"
)

func TimeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add(headerContentType, contentTypeText)
	fmt.Fprintf(w, "The current time is %v", time.Now().Format("Mon Jan _2 15:04:05 2006"))
}
