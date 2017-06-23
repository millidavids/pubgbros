package main

import (
	"fmt"
	"log"
	"net/http"
  "os"
  "io/ioutil"
)

func main() {
  http.HandleFunc("/", handle)
	http.HandleFunc("/_ah/health", healthCheckHandler)
	log.Print("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handle(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

  client := &http.Client{}
  req, _ := http.NewRequest("GET", "https://pubgtracker.com/api/profile/pc/millidavids", nil)
  req.Header.Add("TRN-API-KEY", os.Getenv("TRN_API_KEY"))
  resp, _ := client.Do(req)
  body, _ := ioutil.ReadAll(resp.Body)
  bodyString := string(body)
	fmt.Fprint(w, bodyString)
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "ok")
}
