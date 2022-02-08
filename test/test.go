package main

import (
	"fmt"
)

func main() {
	m0 := map[string]int{
		"Nurlan": 25,
		"Zangar": 28,
	}

	fmt.Printf("%#v\n", m0)

	m1 := map[person]int{}

	m1[person{"Nurlan", 27}] = 1
	m1[person{"Zangar", 27}] = 2

	fmt.Println(m1)
	fmt.Printf("%#v\n", m1)

	var age = 28
	m2 := map[*int]int{&age: 1}
	fmt.Println("m2=", m2)
	fmt.Println(m2[&age])
}

type person struct {
	name string
	age  int
}
