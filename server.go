package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"text/template"
	"time"
	"url-shortener/store"
)

func StartServer(base *url.URL) {
	tmpl := template.Must(template.New("").ParseGlob("./templates/*"))
	s := store.NewMemoryStore("shorten")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := tmpl.ExecuteTemplate(w, "index.go.html", nil)
		if err != nil {
			slog.Error(fmt.Sprintf("Error handling '/': %s", err))
		}
	})

	http.HandleFunc("/shorten", func(w http.ResponseWriter, r *http.Request) {
		to := r.FormValue("url")
		if to == "" {
			http.Error(w, "Invalid URL", http.StatusBadRequest)
			return
		}

		u, err := url.Parse(to)
		if err != nil {
			http.Error(w, "Invalid URL", http.StatusBadRequest)
			return
		}

		if u.Scheme == "" {
			u, err = url.Parse("https://" + to)
			if err != nil {
				http.Error(w, "Invalid URL", http.StatusBadRequest)
				return
			}
		}

		next := s.SetNext(u.String())
		pathURL, err := url.Parse(next)
		if err != nil {
			panic(err)
		}

		err = tmpl.ExecuteTemplate(w, "shorten.go.html", struct {
			URL          string
			ShortenedURL string
		}{
			URL:          u.String(),
			ShortenedURL: base.ResolveReference(pathURL).String(),
		})
		if err != nil {
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("/{param}", func(w http.ResponseWriter, r *http.Request) {
		param := r.PathValue("param")

		u, ok := s.Get(param)
		if !ok {
			http.Error(w, "Not a valid shortened path", http.StatusNotFound)
			return
		}

		http.Redirect(w, r, u, http.StatusFound)
	})

	slog.Info("Server is starting")
	server := &http.Server{
		Addr:         ":" + base.Port(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	slog.Error(server.ListenAndServe().Error())
}
