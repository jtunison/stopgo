package stopgo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

type Resume struct {
	Name        string
	Title       string
	Location    string
	Email       string
	Phone       string
	Summary     string
	History     []Experience
	Supplements []Supplement
	Links       Links
	Education   []Education
	PdfFilename string
	GoogleAnalyticsCode string
}

type Experience struct {
	Role       string
	Company    string
	CompanyUrl string
	StartYear  int
	EndYear    int
	Bullets    []string
}

type Education struct {
	Institution string
	Degree      string
	StartYear   int
	EndYear     int
}

type Supplement struct {
	Heading string
	Bullets []string
}

type Links struct {
	Website    string
	Twitter    string
	Github     string
	GooglePlus string
	LinkedIn   string
}

// see https://gobyexample.com/json
func Load(path string) *Resume {
	log.Printf("reading resume")
	content, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	resume := &Resume{}
	if err := json.Unmarshal([]byte(content), &resume); err != nil {
		panic(err)
	}

	//	if debug {
	//		// echo what we read
	//		enc := json.NewEncoder(os.Stdout)
	//		enc.Encode(resume)
	//	}
	resume.PdfFilename = fmt.Sprintf("%s_Resume.pdf", strings.Replace(resume.Name, " ", "_", -1))

	return resume
}
