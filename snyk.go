package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/michael-go/go-jsn/jsn"
	"github.com/tomnomnom/linkheader"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type filter struct {
	// Languages    []string `json:"languages"`
	// Projects     []string `json:"projects"`
	// Dependencies []string `json:"dependencies"`
	// Licenses     []string `json:"licenses"`
	// DepStatus    string   `json:"depStatus"`
}

type dependenciesFilters struct {
	Filters filter `json:"filters"`
}
type Project []struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
type result []struct {
	ID                     string        `json:"id"`
	Name                   string        `json:"name"`
	Version                string        `json:"version"`
	Type                   string        `json:"type"`
	DependenciesWithIssues []interface{} `json:"dependenciesWithIssues"`
	Licenses               []struct {
		ID      string `json:"id"`
		Title   string `json:"title"`
		License string `json:"license"`
	} `json:"licenses"`
	Projects                   Project   `json:"projects"`
	LatestVersion              string    `json:"latestVersion"`
	LatestVersionPublishedDate time.Time `json:"latestVersionPublishedDate"`
	FirstPublishedDate         time.Time `json:"firstPublishedDate"`
	IsDeprecated               bool      `json:"isDeprecated"`
}

type dependenciesResults struct {
	Results result `json:"results"`
	Total   int    `json:"total"`
}

// GetAllDependencies - Exporting All Dependencies function
func GetAllDependencies(endpointAPI string, orgID string, token string) []jsn.Json {
	var resultsSet []jsn.Json
	fmt.Println("Retrieving data pages")
	nextPage := endpointAPI + "/v1/org/" + orgID + "/dependencies"
	lastPage := false
	for lastPage == false {

		resultsSetPage, links := getDependenciesPage(nextPage, token)
		//resultsSet = append(resultsSet, resultSetPage...)

		r := resultsSetPage.K("results")
		for _, e := range r.Array().Elements() {
			resultsSet = append(resultsSet, e)
		}
		for _, link := range links {
			// fmt.Printf("URL: %s; Rel: %s\n", link.URL, link.Rel)
			if strings.Contains(link.Rel, "next") {
				nextPage = link.URL
				lastPage = false
			}
			if link.Rel == "last" {
				lastPage = true
			}
		}
	}
	fmt.Println("Done retrieving data pages")
	return resultsSet
}

func getDependenciesPage(url string, token string) (jsn.Json, linkheader.Links) {
	dependenciesBody := dependenciesFilters{
		filter{
			// Languages:    []string{},
			// Projects:     []string{},
			// Dependencies: []string{},
			// Licenses:     []string{},
			// DepStatus:    "",
		},
	}

	b, err := json.Marshal(dependenciesBody)
	request, _ := http.NewRequest("POST", url, bytes.NewBuffer(b))
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", "token "+token)
	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	links := linkheader.Parse(response.Header["Link"][0])

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var jsonData jsn.Json

	err2 := json.Unmarshal(responseData, &jsonData)
	if err2 != nil {
		log.Fatal(err)
	}
	return jsonData, links
}
