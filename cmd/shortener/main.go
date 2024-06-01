package main

import (
	"fmt"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/app/shortgen"
	"io"
	"net/http"
	"strings"
)

var linkStore = make(map[string]string)

func main() {
	err := runServer()
	if err != nil {
		panic(err)
	}
}

func runServer() error {
	route, err := routerInit()
	if err != nil {
		return err
	}

	err = http.ListenAndServe(`localhost:8080`, route)
	if err != nil {
		return err
	}

	return nil
}

func routerInit() (*http.ServeMux, error) {

	route := http.NewServeMux()

	route.HandleFunc(`/`, defaultHandler)

	return route, nil
}

func defaultHandler(resp http.ResponseWriter, req *http.Request) {
	if req.RequestURI == string('/') {
		encodeLink(resp, req)
		return
	}

	decodeLink(resp, req)
}

func encodeLink(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		response.WriteHeader(http.StatusForbidden)
		response.Header().Set("content-type", "text/plain")
		body := "Method not allowed"
		_, err := response.Write([]byte(body))
		if err != nil {
			panic(err)
		}
		return

	}

	link, err := io.ReadAll(request.Body)
	if err != nil {
		panic(err)
	}

	short := shortgen.GetShort(7)

	linkStore[short] = string(link)

	response.WriteHeader(http.StatusCreated)
	response.Header().Set("content-type", "text/plain")
	body := fmt.Sprintf("http://localhost:8080/%s", short)
	_, err = response.Write([]byte(body))
	if err != nil {
		panic(err)
	}
}

func decodeLink(resp http.ResponseWriter, req *http.Request) {

	short := strings.ReplaceAll(req.RequestURI, "/", "")

	link, ok := linkStore[short]
	if !ok {
		http.NotFound(resp, req)
	}

	http.Redirect(resp, req, link, http.StatusTemporaryRedirect)
}
