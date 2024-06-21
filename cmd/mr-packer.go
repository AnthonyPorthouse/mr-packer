package main

import (
	"github.com/anthonyporthouse/mr-packer/internal/cli"
	"log"
	"os"
)

func main() {
	app := cli.MakeApp()

	if err := app.Run(os.Args); err != nil {
		log.Fatalln(err)
	}
}
