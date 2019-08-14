package main

import (
	"flag"
	"fmt"
	"github.com/gocolly/colly"
	"html/template"
	"os"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	fmt.Println("Compiling Bill of Material with license text and copyright info")
	orgIDPtr := flag.String("orgID", "", "Your Org ID")
	endpointAPIPtr := flag.String("api", "https://snyk.io/api", "Your API endpoint")
	apiTokenPtr := flag.String("token", "", "Your API token")
	flag.Parse()

	var orgID = *orgIDPtr
	// var projectID string = *projectIDPtr
	var endpointAPI = *endpointAPIPtr
	var apiToken = *apiTokenPtr

	//getDependenciesPageSimpler(endpointAPI+"/v1/org/"+orgID+"/dependencies", apiToken)
	results := GetAllDependencies(endpointAPI, orgID, apiToken)
	var resultsAsMap []map[string]interface{}
	licensesMap = make(map[string]map[string]string)
	for _, result := range results {
		r := result.Raw().(map[string]interface{})
		//fmt.Printf("%s ", r["name"])
		resultLicense := r["licenses"]

		if len(resultLicense.([]interface{})) > 0 {
			rg := resultLicense.([]interface{})[0]
			//fmt.Println(rg.(map[string]interface{})["title"])
			r["licenseTitle"] = rg.(map[string]interface{})["title"]
			r["licenseText"] = consolidateLicensesText(r["licenseTitle"].(string), rg.(map[string]interface{})["id"].(string))
			r["licenseCopyrights"] = ""
			// fmt.Printf("\nPackage: %s %s\n", r["name"], r["version"])
			// fmt.Printf("Used in %d projects", len(r["projects"].([]interface{})))
			// fmt.Printf("\nLicense(s): %s \n", r["licenseTitle"])
			// fmt.Printf("\nCopyrights:\n%s\n", r["licenseCopyrights"])
			// fmt.Printf("Text:\n%s\n", r["licenseText"])
		} else {
			r["licenseTitle"] = "no license"
			r["licenseText"] = "no license"
		}

		//fmt.Println("------------")
		resultsAsMap = append(resultsAsMap, r)
	}

	type TodoPageData struct {
		PageTitle string
		Packages  []map[string]interface{}
	}
	fmap := template.FuncMap{
		"countProjects": func(projectRange []interface{}) int {
			return len(projectRange)
		},
		"returnHTML": func(text string) template.HTML {
			return template.HTML(text)
		},
	}
	tmpl, err := template.New("template.html").Funcs(fmap).ParseFiles("template.html")
	if err != nil {
		panic(err)
	}

	data := TodoPageData{
		PageTitle: "Bill of Material",
		Packages:  resultsAsMap,
	}

	f, err := os.Create("output.html")
	check(err)
	defer f.Close()

	tmpl.Execute(f, data)

	fmt.Printf("Found %d dependencies in use", len(resultsAsMap))

}

var licensesMap map[string]map[string]string

func consolidateLicensesText(licenseTitle string, licenseID string) string {
	if strings.Contains(licenseID, "_OR_") {
		licenseString := strings.Replace(licenseTitle, "(", "", -1)
		licenseString = strings.Replace(licenseTitle, ")", "", -1)
		licenseString = strings.Replace(licenseTitle, " ", "", -1)
		licenseString = strings.Replace(licenseTitle, "Multiple licenses: ", "", -1)
		licenseString = strings.Replace(licenseTitle, "Dual license: ", "", -1)
		licenseStringArray := strings.Split(licenseString, " OR ")
		if len(licenseStringArray) < 2 {
			licenseStringArray = strings.Split(licenseString, ",")
		}

		licenseIDPrefix := "snyk:lic:::"
		var consolidatedLicenseTexts string
		for _, license := range licenseStringArray {
			consolidatedLicenseTexts += getLicenseText(license, licenseIDPrefix+license)
		}
		return consolidatedLicenseTexts
	}
	return getLicenseText(licenseTitle, licenseID)
}

func getLicenseText(licenseTitle string, licenseID string) string {
	if strings.Contains(licenseID, "OR") {

	}
	link := "https://snyk.io/vuln/" + licenseID

	_, ok := licensesMap[licenseTitle]
	if ok {
		return licensesMap[licenseTitle]["Text"]
	}

	var text string
	c := colly.NewCollector()

	c.OnHTML("div[class=license-text]", func(e *colly.HTMLElement) {
		if e.Attr("property") != "spdx:standardLicenseHeader" {
			text, _ = e.DOM.Html()
			licenseData := make(map[string]string)
			licenseData["Text"] = text
			licenseData["Copyrights"] = ""

			licensesMap[licenseTitle] = licenseData
		}
	})
	c.Visit(link)

	return text
}
