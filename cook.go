package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/wgnrj/search"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Ingredient struct {
	Name   string
	Amount float64
	Unit   string
}

func (i *Ingredient) String() string {
	if i.Amount == 0 {
		return i.Name
	}
	return fmt.Sprintf("%v%v %v", i.Amount, i.Unit, i.Name)
}

// Uninterpreted/raw string literal
var validIngredient = regexp.MustCompile(`^((\d+([,\.]\d+)?)([a-zA-Z]*))? ?([A-Za-z0-9äöüÄÖÜß\(\)\- ]+)$`)

func Ingredients(s string) ([]Ingredient, error) {
	var ingredients []Ingredient
	for _, v := range strings.Split(s, "\n") {
		v := strings.Trim(v, "\r\n")
		if v != "" {
			m := validIngredient.FindStringSubmatch(v)
			if m == nil {
				return nil, errors.New("No match")
			}
			var f float64
			var err error
			if m[2] != "" {
				f, err = strconv.ParseFloat(strings.Replace(m[2], ",", ".", 1), 64)
				if err != nil {
					return nil, err
				}
			}
			ingredients = append(ingredients, Ingredient{
				Name:   m[5],
				Amount: f,
				Unit:   m[4],
			})
		}
	}

	return ingredients, nil
}

type Recipe struct {
	Title       string
	Ingredients []Ingredient
	Steps       []string
	Image       string
	Source      string
	Tags        []string
}

func (r *Recipe) String() string {
	var i, s, t string
	for _, v := range r.Ingredients {
		i += "\n" + v.String()
	}
	for _, v := range r.Steps {
		s += "\n\n" + v
	}
	for _, v := range r.Tags {
		t += " " + v
	}
	t = strings.Trim(t, " ")
	return fmt.Sprintf(fmt.Sprintf("%v\n%v%v\n\n%v\n%v\n%v", r.Title, i, s, t, r.Image, r.Source))
}

func (r *Recipe) save() error {
	filename := "data/" + r.Title + ".txt"
	b, err := json.Marshal(r)
	if err != nil {
		log.Printf("Error: Could not encode Recipe %v to JSON (%v).\n", r.Title, err)
	}
	return ioutil.WriteFile(filename, b, 0600)
}

func loadRecipe(title string) (*Recipe, error) {
	filename := "data/" + title + ".txt"
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("Error: Could not read file %v (%v).\n", filename, err)
		return nil, err
	}
	var r Recipe
	err = json.Unmarshal(b, &r)
	if err != nil {
		log.Printf("Error: Could not decode JSON to Recipe %v (%v).\n", title, err)
		return nil, err
	}
	return &r, nil
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fi, err := ioutil.ReadDir("data/")
	if err != nil {
		log.Println(err.Error())
	}
	var s []string
	for _, v := range fi {
		s = append(s, strings.TrimSuffix(v.Name(), ".txt"))
	}
	err = templates.ExecuteTemplate(w, "main.html", s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	recipe, err := loadRecipe(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", recipe)
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	http.Redirect(w, r, "/edit/"+title, http.StatusFound)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	recipe, err := loadRecipe(title)
	if err != nil {
		recipe = &Recipe{Title: title}
	}
	renderTemplate(w, "edit", recipe)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	newTitle := r.FormValue("title")
	// Remove `\r` from the ingredients form value?
	ingredients, err := Ingredients(r.FormValue("ingredients"))
	if err != nil {
		log.Fatal("No ingredients found:", err.Error())
	}
	var steps []string
	for _, v := range strings.Split(r.FormValue("steps"), "\n") {
		v := strings.Trim(v, "\r\n")
		if v != "" {
			steps = append(steps, v)
		}
	}
	image := r.FormValue("image")
	source := r.FormValue("source")
	tags := strings.Fields(r.FormValue("tags"))
	for i, v := range tags {
		if !strings.HasPrefix(v, "#") {
			tags[i] = "#" + v
		}
	}
	recipe := &Recipe{
		Title:       newTitle,
		Ingredients: ingredients,
		Steps:       steps,
		Image:       image,
		Source:      source,
		Tags:        tags,
	}
	if title != newTitle {
		err = os.Remove("data/" + title + ".txt")
		if err != nil {
			log.Printf("Problem removing the file: %v (%v)\n", title, err.Error())
		}
	}
	err = recipe.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+newTitle, http.StatusFound)
}

func removeHandler(w http.ResponseWriter, r *http.Request, title string) {
	if err := os.Remove("data/" + title + ".txt"); err != nil {
		log.Printf("Problem removing the file: %v (%v)\n", title, err.Error())
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	// TODO search for more than one tag, use sets, ...
	var tags []string
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m != nil {
		tags = append(tags, m[2])
	}
	t := r.FormValue("tags")
	if t != "" {
		tags = append(tags, strings.Fields(t)...)
	}
	if len(tags) == 0 {
		http.NotFound(w, r)
		return
	}
	for i, v := range tags {
		if !strings.HasPrefix(v, "#") {
			tags[i] = "#" + v
		}
	}

	files, err := search.Search("data/", tags[0])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound) // TODO
	}
	err = templates.ExecuteTemplate(w, "results.html", files)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	//w.Write([]byte("Not implemented yet."))
}

var templates = template.Must(template.ParseFiles("tmpl/edit.html", "tmpl/view.html", "tmpl/main.html", "tmpl/results.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, r *Recipe) {
	err := templates.ExecuteTemplate(w, tmpl+".html", r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var validPath = regexp.MustCompile(`^/(edit|save|view|delete|search)/([A-Za-zÄÖÜäöüß0-9\- ]+)$`)

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	http.HandleFunc("/delete/", makeHandler(removeHandler))
	http.HandleFunc("/add/", addHandler)
	http.HandleFunc("/search/", searchHandler)

	log.Fatal(http.ListenAndServe(":8081", nil))
}
