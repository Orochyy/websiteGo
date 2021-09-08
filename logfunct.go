package main

import (
	"fmt"
	"sort"
	"strings"
)

//func slices() {
//	slice1 := []int{1, 2, 3}
//	slice2 := append(slice1, 4, 5)
//	fmt.Println(slice1, slice2)
//}
//
//func maps() {
//	elements := map[string]map[string]string{
//		"H": map[string]string{
//			"name":  "Hydrogen",
//			"state": "gas",
//		},
//		"He": map[string]string{
//			"name":  "Helium",
//			"state": "gas",
//		},
//		"Li": map[string]string{
//			"name":  "Lithium",
//			"state": "solid",
//		},
//		"Be": map[string]string{
//			"name":  "Beryllium",
//			"state": "solid",
//		},
//		"B": map[string]string{
//			"name":  "Boron",
//			"state": "solid",
//		},
//		"C": map[string]string{
//			"name":  "Carbon",
//			"state": "solid",
//		},
//		"N": map[string]string{
//			"name":  "Nitrogen",
//			"state": "gas",
//		},
//		"O": map[string]string{
//			"name":  "Oxygen",
//			"state": "gas",
//		},
//		"F": map[string]string{
//			"name":  "Fluorine",
//			"state": "gas",
//		},
//		"Ne": map[string]string{
//			"name":  "Neon",
//			"state": "gas",
//		},
//	}
//
//	if el, ok := elements["Ne"]; ok {
//		fmt.Println(el["name"], el["state"])
//	}
//}
//
///////////
//
//func main() {
//	var x [5]float64
//	x[0] = 98
//	x[1] = 93
//	x[2] = 77
//	x[3] = 82
//	x[4] = 83
//
//	//x := [5]float64{ 98, 93, 77, 82, 83 }
//
//	var total float64 = 0
//	for _, value := range x {
//		total += value
//	}
//
//	fmt.Println(total / float64(len(x)))
//
//	slices()
//	maps()
//}

func slice() {
	users := []string{"Tom", "Alice", "Kate"}
	users = append(users, "Bob")

	for _, value := range users {
		fmt.Println(value)
	}
}

func MinAndMux(x []int) (min int, max int) {

	min = x[0]
	max = x[0]

	for _, value := range x {
		if value < min {
			min = value
		}
		if value > max {
			max = value
		}
	}
	return min, max

}

var listOfStrings = []string{
	"mars bar",
	"milk-duds",
	"awdvawvdawvd23",
	"Mars bar",
	"milk",
	"milky-way",
	"Milk",
	"Milky-way",
	"qilawvdvawway",
	"mars",
	"wavdwavd",
}

type Alphabetic []string

func (list Alphabetic) Len() int {
	return len(list)
}

func (list Alphabetic) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

func (list Alphabetic) Less(i, j int) bool {
	var si = list[i]
	var sj = list[j]
	var si_lower = strings.ToLower(si)
	var sj_lower = strings.ToLower(sj)
	if si_lower == sj_lower {
		return si < sj
	}
	return si_lower < sj_lower
}

func printStrings(slice []string) {
	for i := 0; i < len(slice); i++ {
		fmt.Println(slice[i])
	}
}

func mnoj() {
	var x float64 = 32132
	var y float64 = 42452
	var mn float64 = (x * y)
	fmt.Println(mn)
}
func main() {
	x := []int{
		48, 96, 86, 68,
		57, -82, 63, 70,
		37, 34, 83, 27,
		19, 97, 9, 17,
	}

	min, max := MinAndMux(x)

	fmt.Println("Min:", min)
	fmt.Println("Max:", max)

	//slice()

	sort.Sort(Alphabetic(listOfStrings))
	fmt.Println()
	fmt.Println("SORTED ALPHABETICALLY")
	printStrings(listOfStrings)

	mnoj()
}

///qq~~13reqqsaqnqnqaqvqeqrqesqhqkqo~~oqrqoqcqhqy~~iqlqoqvqeq~~yqoquqsqoqmquqr
//wadwawdzdawdawd
