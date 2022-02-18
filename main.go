package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"text/template"

	"github.com/Spuxy/resume-generator/config"
	d "github.com/Spuxy/resume-generator/data"
	yam "gopkg.in/yaml.v2"
)

type CV struct {
	Name          d.Name `yaml:"name"`
	Email         string
	Social        d.Social `yaml:"social"`
	Theme         d.Theme  `yaml:"theme"`
	About         string
	Service       []d.Service            `yaml:"service"`
	Repos         []d.Repos              `yaml:"repos"`
	Education     []d.Education          `yaml:"education"`
	PrevPositions []d.Previous_positions `yaml:"positions"`
	CurPosition   d.Current_position     `yaml:"current_position"`
	Skills        []d.Skill              `yaml:"skills"`
	Galleries     []d.Gallery            `yaml:"galleries"`
}

func main() {
	var cv CV

	cfg, err := config.New("src")
	if err != nil {
		panic(err)
	}

	resp, err := http.Get(cfg.Get("src"))
	if err != nil {
		panic(err)
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	err = yam.Unmarshal(content, &cv)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/preview", func(rw http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("public/index.html")

		if err != nil {
			panic(err)
		}

		fmt.Println("preview the cv")

		tmpl.Execute(rw, cv)
	})

	http.HandleFunc("/generate", func(rw http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("public/index.html")
		if err != nil {
			panic(err)
		}

		fmt.Println("generate the cv")

		tmpl.Execute(rw, cv)

		f, err := os.Create("index.html")
		defer f.Close()
		if err != nil {
			panic(err)
		}

		tmpl.Execute(f, cv)
	})

	fmt.Println("server is running on port", cfg.Get("port"))
	fmt.Println(http.ListenAndServe(fmt.Sprintf(":%s", cfg.Get("port")), nil))
}
