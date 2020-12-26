package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// Developer is the golang struct to
// the values stored in developers.json
type Developer struct {
	Subdomain   string `json:"subdomain"`
	Name        string `json:"name"`
	Github      string `json:"github"`
	Email       string `json:"email"`
	Youtube     string `json:"youtube"`
	Description string `json:"description"`
}

var (
	tpl        *template.Template
	developers []Developer
	domain     string
	port       string
)

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.html"))

	devsJSON, err := ioutil.ReadFile("developers.json")
	if err != nil {
		log.Println(err)
	}

	err = json.Unmarshal(devsJSON, &developers)
	if err != nil {
		log.Println(err)
	}

	domainEnv, exists := os.LookupEnv("DOMAIN")
	if exists {
		domain = domainEnv
	} else {
		domain = "localhost"
	}

	portEnv, exists := os.LookupEnv("PORT")
	if exists {
		port = portEnv
	} else {
		port = "8080"
	}
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := tpl.ExecuteTemplate(w, "404.html", nil)
		if err != nil {
			log.Fatalln("template didn't execute: ", err)
		}
	})

	for _, developer := range developers {
		http.HandleFunc(fmt.Sprintf("%v.%v/", developer.Subdomain, domain), func(w http.ResponseWriter, r *http.Request) {
			err := tpl.ExecuteTemplate(w, "index.html", developer)

			if err != nil {
				log.Fatalln("template didn't execute: ", err)
			}
		})
	}

	http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
}
