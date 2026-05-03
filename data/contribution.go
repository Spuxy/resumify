package data

type Contribution struct {
	Repo_url string `yaml:"repo_url"`
	Year     int    `yaml:"year"`
	Name     string `yaml:"name"`
	Desc     string `yaml:"desc"`
}
