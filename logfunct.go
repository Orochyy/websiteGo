package main

import "fmt"

func slices() {
	slice1 := []int{1, 2, 3}
	slice2 := append(slice1, 4, 5)
	fmt.Println(slice1, slice2)
}

func maps() {
	x := make(map[string]int)
	x["key"] = 10
	fmt.Println(x["key"])
}

func main() {
	var x [5]float64
	x[0] = 98
	x[1] = 93
	x[2] = 77
	x[3] = 82
	x[4] = 83

	//x := [5]float64{ 98, 93, 77, 82, 83 }

	var total float64 = 0
	for _, value := range x {
		total += value
	}

	fmt.Println(total / float64(len(x)))

	//slices()
	maps()
}
