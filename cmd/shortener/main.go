package main

import (
	"fmt"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/app/shortgen"
	"github.com/go-chi/chi/v5"
	"io"
	"log"
	"net/http"
)

var linkStore = make(map[string]string)

func main() {
	ParseInitFlag()
	err := runServer(AppFlags)
	if err != nil {
		panic(err)
	}
}

func runServer(flags AppFlagsStruct) error {

	baseRoute, shortRoute, err := routerInit()
	if err != nil {
		return err
	}

	go func() { log.Fatal(http.ListenAndServe(flags.GetBaseAddress(), baseRoute)) }()
	go func() { log.Fatal(http.ListenAndServe(flags.GetShortAddress(), shortRoute)) }()
	select {}
}

func routerInit() (*chi.Mux, *chi.Mux, error) {

	baseRoute := chi.NewRouter()
	baseRoute.Post(`/`, encodeLinkHeader)

	shortRoute := chi.NewRouter()
	shortRoute.Get(`/{short}`, decodeLinkHeader)

	return baseRoute, shortRoute, nil
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
	body := fmt.Sprintf("http://%s/%s", AppFlags.GetShortAddress(), short)
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

	short := chi.URLParam(req, "short")

	link, ok := linkStore[short]
	if !ok {
		http.NotFound(resp, req)
	}

	http.Redirect(resp, req, link, http.StatusTemporaryRedirect)
}
