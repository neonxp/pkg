package main

type Config struct {
	Title    string    `json:"title"`
	Host     string    `json:"host"`
	Packages *Packages `json:"packages"`
}

type Packages map[string]Package

type Package struct {
	Pkg         string `json:"pkg"`
	VCS         string `json:"vcs"`
	Repo        string `json:"repo"`
	Description string `json:"desc"`
}

type pageRenderContext struct {
	Title   string
	Package *Package
	Host    string
	Doc     string
}

type indexRenderContext struct {
	Title    string
	Packages *Packages
	Host     string
	Doc      string
}
