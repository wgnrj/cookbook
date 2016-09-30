package main

import (
	"encoding/json"
	"testing"
)

func TestIngredientString(t *testing.T) {
	cases := []struct {
		in   Ingredient
		want string
	}{
		{
			in:   Ingredient{Name: "Chicken", Amount: 1, Unit: "kg"},
			want: "1kg Chicken",
		},
		{
			in:   Ingredient{Name: "Banana", Amount: 0, Unit: ""},
			want: "Banana",
		},
		{
			in:   Ingredient{Name: "Cinnamon", Amount: 1.5, Unit: "EL"},
			want: "1.5EL Cinnamon",
		},
	}

	for _, c := range cases {
		got := c.in.String()
		if got != c.want {
			t.Errorf("Got %q, want %q", got, c.want)
		}
	}
}

func TestStringIngredient(t *testing.T) {
	cases := []struct {
		in   string
		want []Ingredient
	}{
		{
			in: "1kg Apples\n4TL Honey\n0.5TL Cinnamon",
			want: []Ingredient{
				{Name: "Apples", Amount: 1, Unit: "kg"},
				{Name: "Honey", Amount: 4, Unit: "TL"},
				{Name: "Cinnamon", Amount: 0.5, Unit: "TL"},
			},
		},
		{
			in: "1EL Mayo\n\n",
			want: []Ingredient{
				{Name: "Mayo", Amount: 1, Unit: "EL"},
			},
		},
		{
			in: "Banana",
			want: []Ingredient{
				{Name: "Banana", Amount: 0, Unit: ""},
			},
		},
        {
			in: "70g Parmesan\n70g Pinienkerne\r\n1 Knoblauchzehe\r\n100ml Olivenöl\nSalz\n15 getrocknete Tomaten in Öl eingelegt\r\n3EL Fruchtessig Mango\n\n",
			want: []Ingredient{
				{Name: "Parmesan", Amount: 70, Unit: "g"},
				{Name: "Pinienkerne", Amount: 70, Unit: "g"},
				{Name: "Knoblauchzehe", Amount: 1, Unit: ""},
				{Name: "Olivenöl", Amount: 100, Unit: "ml"},
				{Name: "Salz", Amount: 0, Unit: ""},
				{Name: "getrocknete Tomaten in Öl eingelegt", Amount: 15, Unit: ""},
				{Name: "Fruchtessig Mango", Amount: 3, Unit: "EL"},
			},
        },
	}

	for _, c := range cases {
		got, err := Ingredients(c.in)
		if err != nil {
			t.Errorf("Ingredients(%q) returned with error: %q", c.in, err.Error())
		}
		if !equals(got, c.want) {
			t.Errorf("Ingredients(%q) == %q, want %q", c.in, got, c.want)
		}
	}
}

func equals(a, b []Ingredient) bool {
	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func TestRecipeString(t *testing.T) {
	r_want := "Chicken\n\n1kg Chicken\n1TL Salt\n4EL Honey\n1Piece Peperoni\n\n1. Cook chicken.\n\n2. Eat chicken.\n\n#Hühnchen #Hauptspeise\nchicken.jpg\nhttp://chicken.go/"
	r_in := &Recipe{
		Title: "Chicken",
		Ingredients: []Ingredient{
			{Name: "Chicken", Amount: 1, Unit: "kg"},
			{Name: "Salt", Amount: 1, Unit: "TL"},
			{Name: "Honey", Amount: 4, Unit: "EL"},
			{Name: "Peperoni", Amount: 1, Unit: "Piece"},
		},
		Steps:  []string{"1. Cook chicken.", "2. Eat chicken."},
		Image:  "chicken.jpg",
		Source: "http://chicken.go/",
        Tags: []string{"Hühnchen", "Hauptspeise"},
	}

    if got := r_in.String(); got != r_want {
        t.Errorf("Got:\n%q\nWant:\n%q", got, r_want)
	}
}

func TestJson(t *testing.T) {
	r_in := &Recipe{
		Title: "Chicken",
		Ingredients: []Ingredient{
			{Name: "Chicken", Amount: 1, Unit: "kg"},
			{Name: "Salt", Amount: 1, Unit: "TL"},
			{Name: "Honey", Amount: 4, Unit: "EL"},
			{Name: "Peperoni", Amount: 1, Unit: "Piece"},
		},
		Steps:  []string{"1. Cook chicken.", "2. Eat chicken."},
		Image:  "chicken.jpg",
		Source: "http://chicken.go/",
	}

	j, err := json.Marshal(r_in)
	if err != nil {
		t.Error("Error in the encoding of Recipe: ", err)
	}
	var J Recipe
	err = json.Unmarshal(j, &J)
	if err != nil {
		t.Error("Error in the decoding of Recipe: ", err)
	}
	if r_in.String() != J.String() {
		t.Error("Error: Expected value doesn't match the en- and decoded struct.")
	}
}
