package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Hello does
type Hello struct {
	l *log.Logger
}

// NewHello takes in a logger struct and returns a Hello 
// handler containing said struct
func NewHello(l *log.Logger) *Hello {
	return &Hello{l}
}

func (h *Hello) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.l.Println("hello world!")

	d, err := ioutil.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "Oops, something went wrong", http.StatusBadRequest)
		return
	}

	h.l.Printf("Data %s\n", d)

	fmt.Fprintf(w, "Hello %s\n", d)
}
