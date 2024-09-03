package main

import (
	"fmt"
	"io"
	"os"
	"testing"
)

var (
	PathToEventJsonFile = "../events/event.json"
)

func TestHandler(t *testing.T) {
	jsonFile, err := os.Open(PathToEventJsonFile)
	if err != nil {
		t.Fatal("Unable to open event.json")
	}

	fmt.Println("Successfully Opened event.json")

	defer jsonFile.Close()

	jsonStream, _ := io.ReadAll(jsonFile)

	event, _ := UnmarshalEvent(jsonStream)

	t.Run("Successful Request", func(t *testing.T) {
		err := handler(nil, event)
		if err != nil {
			t.Fatal("Everything should be ok")
		}
	})
}
