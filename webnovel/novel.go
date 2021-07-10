package webnovel

import (
	"encoding/json"
	"github.com/gosimple/slug"
	"io/ioutil"
)

type Webnovel struct {
	Name         string
	Author       []string
	Genre        []string
	TotalChapter int
	URL          string
	Chapters
}

func (w Webnovel) Save() {
	jsonName := slug.Make(w.Name)

	file, _ := json.MarshalIndent(w.Chapters, "", " ")
	_ = ioutil.WriteFile(jsonName+".json", file, 0644) //nolint:gosec
}