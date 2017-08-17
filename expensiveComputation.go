package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

type Answer struct {
	id     int
	result float64
}

func expensiveCompuation(data int, id int, answer chan Answer) {
	a := Answer{}
	result := float64(data)

	for result > 1.0001 {
		result = math.Sqrt(result)
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(500)))
	}
	a.id = id
	a.result = result
	answer <- a
}

func main() {
	doneCount := 0
	waitTag := 0
	answer1 := make(chan Answer)
	answer2 := make(chan Answer)
	go expensiveCompuation(9999, 1, answer1)
	go expensiveCompuation(10, 2, answer2)

	for doneCount < 2 {
		select {
		case a1 := <-answer1:
			doneCount++
			fmt.Printf("%d -> %g\n", a1.id, a1.result)
		case a2 := <-answer2:
			doneCount++
			fmt.Printf("%d -> %g\n", a2.id, a2.result)
		default:
			if waitTag%10000 == 0 {
				fmt.Printf("Waiting ... %d\n", waitTag)
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(500)))
			}
			waitTag++
		}
	}
	fmt.Println()
}
