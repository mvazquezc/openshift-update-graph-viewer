package main

import (
	"crypto/tls"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

var version = "v0.0.1"
var graphDataGitHubUrl = "https://api.github.com/repos/openshift/cincinnati-graph-data/contents/channels"

// ReturnIndex renders the index template
func ReturnIndex(w http.ResponseWriter, r *http.Request) {
	indexTemplate := template.Must(template.ParseFiles("templates/index-template.html"))
	err := indexTemplate.ExecuteTemplate(w, "index-template.html", "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// ReturnOpenShiftChannels gets the channel list from GitHub and returns it as json
func ReturnOpenShiftChannels(w http.ResponseWriter, r *http.Request) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	req, _ := http.NewRequest("GET", graphDataGitHubUrl, nil)

	res, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write([]byte(string(body)))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// ReturnCincinnatiOutput gets the update graph from the cincinnati server + channel specified
func ReturnCincinnatiOutput(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	cincinnatiChannel := vars["channel"]
	cincinnatiApiUrl := vars["api"]
	if strings.HasSuffix(cincinnatiApiUrl, "/") {
		cincinnatiApiUrl = strings.TrimSuffix(cincinnatiApiUrl, "/")
	}
	cincinnatiUrl := "https://" + cincinnatiApiUrl + "/api/upgrades_info/v1/graph?channel=" + cincinnatiChannel

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	req, _ := http.NewRequest("GET", cincinnatiUrl, nil)
	req.Header.Add("Accept", "application/json")
	res, err := client.Do(req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write([]byte(string(body)))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

//getEnv returns the value for a given Env Var
func getEnv(varName, defaultValue string) string {
	if varValue, ok := os.LookupEnv(varName); ok {
		return varValue
	}
	return defaultValue
}

func main() {
	port := getEnv("APP_PORT", "8080")
	ip := getEnv("IP", "127.0.0.1")
	log.Println("Starting OpenShift Update Graph", version)
	log.Println("Listening on", ip+":"+port)
	router := mux.NewRouter()
	router.HandleFunc("/", ReturnIndex).Methods("GET")
	router.HandleFunc("/channels", ReturnOpenShiftChannels).Methods("GET")
	router.HandleFunc("/cincinnati/{channel}/{api}", ReturnCincinnatiOutput).Methods("GET")
	log.Fatal(http.ListenAndServe(ip+":"+port, router))
}
