package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

const github_project_url = "https://github.com/schreibe72/8021x-caffeinate"

type githubVersion struct {
	ID                      int    `json:"id"`
	TagName                 string `json:"tag_name"`
	UpdateURL               string `json:"update_url"`
	UpdateAuthenticityToken string `json:"update_authenticity_token"`
	DeleteURL               string `json:"delete_url"`
	DeleteAuthenticityToken string `json:"delete_authenticity_token"`
	EditURL                 string `json:"edit_url"`
}

func check4update(v string) (string, string) {
	u := fmt.Sprintf("%s/releases/latest", github_project_url)
	client := &http.Client{}
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		log.Printf("error building Request: %s", err)
		return "", ""
	}
	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Printf("error fetching Request: %s", err)
		return "", ""
	}
	var g githubVersion

	if err := json.Unmarshal(body, &g); err != nil {
		log.Printf("error unmarshaling body: %s\n%s", err, string(body))
		return "", ""
	}
	if g.TagName == version {
		return "", ""
	}
	return fmt.Sprintf("https://github.com%s", g.UpdateURL), g.TagName
}
