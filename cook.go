package main

import (
	"log"
	"net/http"
	"fmt"
	"io/ioutil"
)

type Ingredient struct {
	Name   string
	Amount float64
	Unit   string
}

func (i *Ingredient) String() string {
    return fmt.Sprintf("%v %v %v", i.Amount, i.Unit, i.Name)
}

type Recipe struct {
	Title       string
	Ingredients []Ingredient
	Steps       []string
	Image       string
	Source      string
}

func (r *Recipe) String() string {
    var i, s string
    for _, v := range r.Ingredients {
        i += "\n"+v.String()
    }
    for _, v := range r.Steps {
        s += "\n\n"+v
    }
    return fmt.Sprintf(fmt.Sprintf("%v\n%v%v\n\n%v\n%v", r.Title, i, s, r.Image, r.Source))
}

func (r *Recipe) save() error {
    filename := "data/" + r.Title + ".txt"
    return ioutil.WriteFile(filename, []byte(r.String()), 0600)
}

func loadRecipe(name string) *Recipe {
    return &Recipe{}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world!"))
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
}

func editHandler(w http.ResponseWriter, r *http.Request) {
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
}

func main() {
	http.HandleFunc("/", rootHandler)

	log.Fatal(http.ListenAndServe(":8081", nil))
}
