package stopgo

import (
	"bytes"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"os"
)

func mainHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("web/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	resume := Load("resume.json")

	t.Execute(w, resume)
}

func copy(source, destination string) {
	cpCmd := exec.Command("cp", "-rf", source, destination)
	err := cpCmd.Run()
	if err != nil {
		panic(err)
	}
}

func rmrf(path string) {
	cmd := exec.Command("rm", "-rf", path)
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

func mkdir(path string) {
	cmd := exec.Command("mkdir", path)
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

func Server() {

	r := mux.NewRouter()

	r.HandleFunc("/", mainHandler)
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./web/assets/"))))
	http.Handle("/", r)

	log.Printf("Listening on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	fmt.Println(err.Error())
}

func Build(overlayPath, outputDir string) {

	// load template & model
	packageDir := fmt.Sprintf("%s/src/github.com/jtunison/stopgo", os.Getenv("GOPATH"))
	t, err := template.ParseFiles(fmt.Sprintf("%s/web/index.html", packageDir))
	if err != nil {
		panic(err)
	}
	model := Load("resume.json")

	// housekeeping
	rmrf("public")
	mkdir("public")
	copy(fmt.Sprintf("%s/web/assets", packageDir), fmt.Sprintf("%s/assets", outputDir))
	if overlayPath!="" {
		copy(fmt.Sprintf("%s/assets", overlayPath), outputDir)		
	}

	// generate website
	var doc bytes.Buffer
	t.Execute(&doc, model)
	err = ioutil.WriteFile(fmt.Sprintf("%s/index.html", outputDir), doc.Bytes(), 0644)
	if err != nil {
		panic(err)
	}

	// generate pdf
	Write(fmt.Sprintf("%s/%s", outputDir, model.PdfFilename), model)

}
