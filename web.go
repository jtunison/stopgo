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

func Build() {

	// load template & model
	t, err := template.ParseFiles("../stopgo/web/index.html")
	if err != nil {
		panic(err)
	}
	model := Load("resume.json")

	// housekeeping
	rmrf("public")
	mkdir("public")
	copy("../stopgo/web/assets", "public/assets")
	copy("overlay/assets", "public")

	// generate website
	var doc bytes.Buffer
	t.Execute(&doc, model)
	err = ioutil.WriteFile("public/index.html", doc.Bytes(), 0644)
	if err != nil {
		panic(err)
	}

	// generate pdf
	Write("public/"+model.PdfFilename, model)

}
