package main

import (
	"fmt"
	"math"
)

func cal(p float64, r float64, n float64) (m float64) {
	var res0, res1, res2 float64

	res0 = r / 12
	res1 = math.Pow(1+res0, n)
	res2 = p * res0
	m = res2 * res1 / (res1 - 1)

	return m

}

func main() {
	c := cal(1, 2, 3)

	fmt.Println(c)
}
