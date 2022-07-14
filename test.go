package main

import (
	"log"
)

func main() {

	//find lowest common multiple

	data := []int32{0, 1, 2, 3, 4, 5}

	log.Println(data)

	data = removeAtIndex(data, 5)

	log.Println(data)

}

func removeAtIndex(v []int32, index int) []int32 {
	if index >= len(v) || index < 0 {
		return v
	}
	return append(v[:index], v[index+1:]...)
}
