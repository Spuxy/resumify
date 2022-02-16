package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"text/template"

	yam "gopkg.in/yaml.v2"
)

type social struct {
	Github   string
	Twitter  string
	Linkedin string
}

type theme struct {
	Style string
	Color string
}

type service struct {
	Details string `yaml:"details"`
	Year    int    `yaml:"year"`
	Url     string `yaml:"url"`
}

type repos struct {
	Repo_url string `yaml:"repo_url"`
	Year     int    `yaml:"year"`
	Name     string `yaml:"name"`
	Desc     string `yaml:"desc"`
}

type education struct {
	School   string `yaml:"school"`
	Location string `yaml:"location"`
	Degree   string `yaml:"degree"`
	Dates    string `yaml:"dates"`
}

type skill struct {
	Title   string `yaml:"title"`
	Details string `yaml:"details"`
}

type previous_positions struct {
	Place         string `yaml:"place"`
	Title         string `yaml:"title"`
	Inline_detail string `yaml:"inline_detail"`
	Dates         string `yaml:"dates"`
}

type current_position struct {
	Place    string `yaml:"place"`
	Location string `yaml:"location"`
	Title    string `yaml:"title"`
	Dates    string `yaml:"dates"`
	Website  string `yaml:"website"`
}

type name struct {
	First string `yaml:"first"`
	Last  string `yaml:"last"`
}

type CV struct {
	Name          name `yaml:"name"`
	Email         string
	Social        social
	Theme         theme
	About         string
	Service       []service            `yaml:"service"`
	Repos         []repos              `yaml:"repos"`
	Education     []education          `yaml:"education"`
	PrevPositions []previous_positions `yaml:"positions"`
	CurPosition   current_position     `yaml:"current_position"`
	Skills        []skill              `yaml:"skills"`
}

func main() {
	var cv CV

	content, err := ioutil.ReadFile("cv/cv.yml")
	if err != nil {
		panic(err)
	}

	err = yam.Unmarshal(content, &cv)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/generate", func(rw http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("public/index.html")
		if err != nil {
			panic(err)
		}

		f, err := os.Create("index.html")
		if err != nil {
			panic(err)
		}
		defer f.Close()

		tmpl.Execute(f, cv)
		tmpl.Execute(rw, cv)
	})

	http.ListenAndServe(":8090", nil)
}
