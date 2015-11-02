package main

import (
	"log"
)

func test(tmp *map[string]string) {
	log.Println(*tmp)
}

func main() {
	tmp := make(map[string]string)
	test(&tmp)
}
