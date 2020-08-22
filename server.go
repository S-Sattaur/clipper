package main

import (
	"encoding/json"
	"fmt"
	"github.com/kkdai/youtube"
	"log"
	"net/http"
)

type VideoBlueprint struct {
	Url string
}

func main() {
	fileServer := http.FileServer(http.Dir("./public")) // New code
	http.Handle("/", fileServer)                        // New code
	http.HandleFunc("/api/grab", urlHandler)

	fmt.Printf("Starting server at port 3000\n")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal(err)
	}
}

func urlHandler(res http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/api/grab" {
		http.Error(res, "404 not found.", http.StatusNotFound)
		return
	}
	if req.Method != "POST" {
		http.Error(res, "Method is not supported.", http.StatusNotFound)
		return
	}

	var video VideoBlueprint

	json.NewDecoder(req.Body).Decode(&video)

	log.Println(video.Url)

	client := youtube.Client{}
	videoObj, err := client.GetVideo(video.Url)
	if err != nil {
		log.Println("Could not get video ", video.Url)
		fmt.Fprint(res, `{ "success": false }`)
		return
	}

	res.Header().Set("Content-Type", "application/json")

	if len(videoObj.Formats) > 0 {
		fmt.Fprintf(res, `{
			"success": true,
			"stream": "%s"
		}`, videoObj.Formats[0].URL)
	} else {
		fmt.Fprint(res, `{ "success": false }`)
	}
}