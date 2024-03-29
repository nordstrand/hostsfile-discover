package main

import (
	"embed"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
)

//go:embed templates/*
var Assets embed.FS

func main() {
	entries, err := getHostFileEntries()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%d entries matching TLD .'%s' found in '%s':\n %v \n",
		len(entries),
		CONFIG.TLD(),
		CONFIG.HOSTS_FILE_PATH(), entries)

	handler := func(w http.ResponseWriter, req *http.Request) {
		host := getHost(req)
		entries, err := getEntriesMatching(host)

		if err != nil {
			log.Println("Couldn't get hosts")
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, "Couldnt get host\n")
			return
		}

		tmpl, err := template.ParseFS(Assets, "templates/index.html", "templates/pico.classless.css")
		if err != nil {
			log.Fatal(err)
		}
		split := strings.Split(strings.TrimSuffix(host, "."), ".")
		parentHost := host
                if len(split) > 1 {
			parentHost = strings.Join(split[1:], ".")+"."
		} 

		data := struct {
			Host string
			ParentHost string
			List []hostfile_entry
		}{
			Host: host,
			ParentHost: parentHost,
			List: entries,
		}
		err = tmpl.Execute(w, data)
		if err != nil {
			log.Fatal(err)
		}
	}

	http.HandleFunc("/", handler)
	log.Printf("Listening for requests at :%d", CONFIG.PORT())
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", CONFIG.PORT()), logRequest(http.DefaultServeMux)))
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s\n", r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func getHost(req *http.Request) string {
	log.Printf("%s", req.Host)

	if !(strings.HasPrefix(req.Host, "localhost") || strings.HasPrefix(req.Host, "127.0.0.1")) {
		return req.Host
	}

	host, ok := req.URL.Query()["host"]

	if !ok || len(host[0]) < 1 {
		log.Println("Url Param 'host' is missing")
		return ""
	}

	h := host[0]

	log.Printf("Url Param 'host' is: %s\n", string(h))
	return h
}
