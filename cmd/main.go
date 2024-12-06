package main

import (
	"bccChat/internal/app"
)

func main() {

	app := app.New()
	err := app.Start()
	if err != nil {
		panic(err)
	}

}
