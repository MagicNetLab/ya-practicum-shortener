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

	route.HandleFunc(`/`, encodeLinkHeader)
	route.HandleFunc(`/sl/`, decodeLinkHeader)

	return route, nil
}

func encodeLinkHeader(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		response.Header().Set("content-type", "text/plain")
		response.WriteHeader(http.StatusForbidden)
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

	if string(link) == "" {
		response.Header().Set("content-type", "text/plain")
		response.WriteHeader(http.StatusBadRequest)
		body := "Missing link"
		_, err := response.Write([]byte(body))
		if err != nil {
			panic(err)
		}
		return
	}

	short := shortgen.GetShort(7)

	linkStore[short] = string(link)

	response.Header().Set("content-type", "text/plain")
	response.WriteHeader(http.StatusCreated)
	body := fmt.Sprintf("http://localhost:8080/sl/%s", short)
	_, err = response.Write([]byte(body))
	if err != nil {
		panic(err)
	}
}

func decodeLinkHeader(resp http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		resp.Header().Set("content-type", "text/plain")
		resp.WriteHeader(http.StatusForbidden)
		body := "Method not allowed"
		_, err := resp.Write([]byte(body))
		if err != nil {
			panic(err)
		}
		return
	}

	short := strings.ReplaceAll(req.RequestURI, "/sl/", "")

	link, ok := linkStore[short]
	if !ok {
		http.NotFound(resp, req)
	}

	http.Redirect(resp, req, link, http.StatusTemporaryRedirect)
}
