package main

import (
	"fmt"
	"io"
	"log"
)

func main() {
	client := &KSeFClient{}

	resp, err := client.GetAuthChallenge()
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(body))
}
