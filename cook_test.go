package main

import (
    "fmt"
    "testing"
)

func TestIngredientString(t *testing.T) {
    i := &Ingredient{
        Name: "Chicken",
        Amount: 1,
        Unit: "kg",
    }

    if i.String() != "1 kg Chicken" {
        t.Error("Expected 1 kg Chicken, got ", i)
    }
}

func TestRecipeString(t *testing.T) {
    r_should := "Chicken\n\n1 kg Chicken\n1 TL Salt\n4 EL Honey\n1 Piece Peperoni\n\n1. Cook chicken.\n\n2. Eat chicken.\n\nchicken.jpg\nhttp://chicken.go/"
    r_is := &Recipe{
        Title: "Chicken",
        Ingredients: []Ingredient{
            {
                Name: "Chicken",
                Amount: 1,
                Unit: "kg",
            },
            {
                Name: "Salt",
                Amount: 1,
                Unit: "TL",
            },
            {
                Name: "Honey",
                Amount: 4,
                Unit: "EL",
            },
            {
                Name: "Peperoni",
                Amount: 1,
                Unit: "Piece",
            },
        },
        Steps: []string{
            "1. Cook chicken.",
            "2. Eat chicken.",
        },
        Image: "chicken.jpg",
        Source: "http://chicken.go/",
    }

    if r_is.String() != r_should {
        t.Error("Recipe string didn't match expectation.")
    }
}

