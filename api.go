package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/gorilla/mux"
)

type APIServer struct {
	store      Storage
	listenAddr string
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			//handle error
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

func NewAPIServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/api/meta", makeHTTPHandleFunc(s.handleGetMeta))
	http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) handleGetMeta(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "POST" {
		return s.handleGetPageMeta(w, r)
	}

	return fmt.Errorf("method %s not supported", r.Method)
}

func (s *APIServer) handleGetPageMeta(w http.ResponseWriter, r *http.Request) error {
	getPageMetaRequest := new(GetPageMetaRequest)
	if err := json.NewDecoder(r.Body).Decode(getPageMetaRequest); err != nil {
		return err
	}

	urlReq := getPageMetaRequest.Url
	_, urlError := url.ParseRequestURI(urlReq)
	if urlError != nil {
		return WriteJSON(w, http.StatusInternalServerError, "Invalid Url")
	}
	pageMeta, pageMetaErr := s.store.GetPageMetaByUrl(urlReq)

	if pageMetaErr != nil || pageMeta == nil {
		apiKey, apiKeyErr := s.store.GetRandomApiKey()
		fmt.Println(apiKey.Key)
		if apiKeyErr != nil {
			fmt.Print(apiKeyErr.Error())
		}
		response, resErr := http.Get("https://iframe.ly/api/oembed?url=" + urlReq + "&api_key=" + apiKey.Key)

		if resErr != nil {
			fmt.Print(resErr.Error())
		}
		defer response.Body.Close()
		pageMetaResponse := new(PageMeta)

		json.NewDecoder(response.Body).Decode(pageMetaResponse)
		videoID, videoErr := getYouTubeVideoID(urlReq)
		pageMetaResponse.DataIframelyUrl = strings.Contains(pageMetaResponse.Html, "data-iframely-url")
		pageMetaResponse.Url = urlReq
		if videoErr == nil {
			pageMetaResponse.YoutubeVideoId = videoID
		}
		newData, err := s.store.CreatePageMeta(pageMetaResponse)
		if err != nil {
			fmt.Print(err.Error())
		}
		return WriteJSON(w, http.StatusOK, newData)
	}

	return WriteJSON(w, http.StatusOK, pageMeta)

}

func getYouTubeVideoID(url string) (string, error) {
	// Regular expression to match YouTube video ID
	regex := regexp.MustCompile(`(?:https?:\/\/)?(?:www\.)?(?:youtube\.com\/(?:[^\/\n\s]+\/\S+\/|(?:v|e(?:mbed)?)\/|\S*?[?&]v=)|youtu\.be\/)([a-zA-Z0-9_-]{11})`)

	// Find submatch in the URL
	match := regex.FindStringSubmatch(url)

	// Check if a match is found
	if len(match) < 2 {
		return "", fmt.Errorf("YouTube video ID not found in the URL")
	}

	// Return the YouTube video ID
	return match[1], nil
}
