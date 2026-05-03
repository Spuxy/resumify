package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"text/template"

	"github.com/Spuxy/resume-generator/config"
	d "github.com/Spuxy/resume-generator/data"
	yam "gopkg.in/yaml.v2"
)

type CV struct {
	Name          d.Name `yaml:"name"`
	Email         string
	Status        string             `yaml:"status"`
	Social        d.Social           `yaml:"social"`
	CurPosition   d.Current_position `yaml:"current_position"`
	Theme         d.Theme            `yaml:"theme"`
	About         string
	Photo         string                 `yaml:"photo"`
	Assets        string                 `yaml:"assets"`
	Service       []d.Service            `yaml:"service"`
	Education     []d.Education          `yaml:"education"`
	PrevPositions []d.Previous_positions `yaml:"positions"`
	Repos         []d.Repos              `yaml:"repos"`
	Skills        []d.Skill              `yaml:"skills"`
	Galleries     []d.Gallery            `yaml:"galleries"`
}

func main() {
	build := flag.Bool("build", false, "generate index.html and exit (CI/GitHub Pages mode)")
	src := flag.String("src", "", "override YAML source (URL or local file path)")
	out := flag.String("out", "index.html", "output path for --build mode")
	flag.Parse()

	cfg, err := config.New("src")
	if err != nil {
		panic(err)
	}

	source := cfg.Get("src")
	if *src != "" {
		source = *src
	}

	cv, err := loadCV(source)
	if err != nil {
		panic(err)
	}

	if *build {
		if err := generate(cv, *out); err != nil {
			panic(err)
		}
		fmt.Println("generated", *out)
		return
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
		if err != nil {
			panic(err)
		}
		defer f.Close()
		tmpl.Execute(f, cv)
	})

	fmt.Println("server is running on port", cfg.Get("port"))
	fmt.Println(http.ListenAndServe(fmt.Sprintf(":%s", cfg.Get("port")), nil))
}

func loadCV(source string) (CV, error) {
	var cv CV
	var content []byte

	if strings.HasPrefix(source, "http://") || strings.HasPrefix(source, "https://") {
		resp, err := http.Get(source)
		if err != nil {
			return cv, err
		}
		defer resp.Body.Close()
		content, err = io.ReadAll(resp.Body)
		if err != nil {
			return cv, err
		}
	} else {
		var err error
		content, err = os.ReadFile(source)
		if err != nil {
			return cv, err
		}
	}

	if err := yam.Unmarshal(content, &cv); err != nil {
		return cv, err
	}
	return cv, nil
}

func generate(cv CV, out string) error {
	tmpl, err := template.ParseFiles("public/index.html")
	if err != nil {
		return err
	}
	f, err := os.Create(out)
	if err != nil {
		return err
	}
	defer f.Close()
	return tmpl.Execute(f, cv)
}
