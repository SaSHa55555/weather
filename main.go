package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"os"
)

var weatherData struct {
	Weather []struct {
		Main        string `json:"main"`
		Description string `json:"description"`
	} `json:"weather"`
	Main struct {
		Temperature float64 `json:"temp"`
		Humidity    int     `json:"humidity"`
	} `json:"main"`
}

var api = "e7704bc895b4a8d2dfd4a29d404285b6"
var tpl = template.Must(template.ParseFiles("template/index.html"))
var tpl1 = template.Must(template.ParseFiles("template/ans.html"))

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tpl.Execute(w, nil)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	u, err := url.Parse(r.URL.String())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}

	params := u.Query()
	searchKey := params.Get("q")

	city := searchKey

	search := &weatherData

	endpoint := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%v&appid=%v", city, api)
	resp, err := http.Get(endpoint)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewDecoder(resp.Body).Decode(&search)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = tpl1.Execute(w, search)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/search", searchHandler)
	mux.HandleFunc("/", indexHandler)
	http.ListenAndServe(":"+port, mux)
}
