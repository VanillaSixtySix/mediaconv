package main

import (
	"encoding/json"
	"io"
	"log"
	"mediaconv/internal/mediaconv"
	"net/http"
	"os"
	"os/exec"
	"strconv"
)

func main() {
	config, err := mediaconv.LoadConfig("config.json")
	if err != nil {
		log.Fatalln("failed to load config:", err)
	}

	http.HandleFunc("POST /", func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println("failed to read body bytes:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if r.Header.Get("Content-Type") != "application/json" {
			log.Println("request content-type is not application/json")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var convRequest mediaconv.Request
		if err = json.Unmarshal(bodyBytes, &convRequest); err != nil {
			log.Println("failed to unmarshal json body:", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if convRequest.URL == "" {
			log.Println("convert request has no \"url\"")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		resp, err := http.Get(convRequest.URL)
		if err != nil {
			log.Println("failed to get source url:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		tempSrcFile, err := os.CreateTemp("", "*.mp4")
		if err != nil {
			log.Println("failed to create temp source file:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if _, err = io.Copy(tempSrcFile, resp.Body); err != nil {
			log.Println("failed to copy source file to temp file:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		tempDestFile, err := os.CreateTemp("", "*-mediaconv.gif")
		if err != nil {
			log.Println("failed to create temp destination file:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		cmd := exec.Command("ffmpeg", "-y", "-i", tempSrcFile.Name(), "-filter_complex", "split[s0][s1];[s0]palettegen=max_colors=32[p];[s1][p]paletteuse=dither=bayer", tempDestFile.Name())
		cmd.Stdout = os.Stderr
		cmd.Stderr = os.Stderr
		if err = cmd.Start(); err != nil {
			log.Println("failed to start ffmpeg:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err = cmd.Wait(); err != nil {
			log.Println("failed to finish conversion:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if _, err = tempDestFile.Seek(0, 0); err != nil {
			log.Println("failed to seek start of temp destination file:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if _, err = io.Copy(w, tempDestFile); err != nil {
			log.Println("failed to copy temp destination file to http.ResponseWriter:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err = tempSrcFile.Close(); err != nil {
			log.Println("failed to close temp source file:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err = os.Remove(tempSrcFile.Name()); err != nil {
			log.Println("failed to remove temp source file:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err = tempDestFile.Close(); err != nil {
			log.Println("failed to close temp destination file:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err = os.Remove(tempDestFile.Name()); err != nil {
			log.Println("failed to remove temp destination file:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
	port := strconv.Itoa(config.Port)
	log.Println("Listening on :" + port)
	log.Fatalln(http.ListenAndServe(":"+port, nil))
}
