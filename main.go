package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

func mustEnv(key string) string {
	v, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("Missing env: %v", key)
	}
	return v
}

var (
	baseURL = mustEnv("PAPERLESS_URL")
	user    = mustEnv("PAPERLESS_USER")
	pass    = mustEnv("PAPERLESS_PASS")
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/Image.pdf", func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "PUT" {
			http.NotFound(w, req)
			return
		}

		log.Println("PDF received. Forwarding to Paperless...")
		buf := &bytes.Buffer{}

		m := multipart.NewWriter(buf)
		f, err := m.CreateFormFile("document", fmt.Sprintf("scan_%s.pdf", time.Now().Format("2006-01-02_15-04-05")))
		if err != nil {
			http.Error(w, fmt.Sprintf("failed creating form file: %v", err), http.StatusInternalServerError)
			return
		}

		if _, err := io.Copy(f, req.Body); err != nil {
			http.Error(w, fmt.Sprintf("failed copying data: %v", err), http.StatusInternalServerError)
			return
		}

		if err := m.Close(); err != nil {
			http.Error(w, fmt.Sprintf("failed closing form: %v", err), http.StatusInternalServerError)
			return
		}

		pReq, err := http.NewRequest("POST", baseURL+"/api/documents/post_document/", buf)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed creating request: %v", err), http.StatusInternalServerError)
			return
		}
		pReq.SetBasicAuth(user, pass)
		pReq.Header.Add("Content-Type", m.FormDataContentType())

		resp, err := http.DefaultClient.Do(pReq)
		if err != nil {
			return
		}

		if resp.StatusCode != 200 {
			http.Error(w, fmt.Sprintf("invalid response: %v", resp.StatusCode), http.StatusInternalServerError)
			return
		}

		log.Println("Document uploaded!")
		w.WriteHeader(http.StatusOK)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
