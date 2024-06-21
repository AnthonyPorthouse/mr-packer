package cli

import (
	"github.com/anthonyporthouse/mr-packer/internal"
	"github.com/urfave/cli/v2"
	"log"
)

func MakeApp() *cli.App {
	app := &cli.App{
		Name:  "mr-packer",
		Usage: "Unpack Modrinth Modpacks",
		Commands: []*cli.Command{
			{
				Name:    "client",
				Aliases: []string{"c"},
				Usage:   "Extract a client side modpack",
				Action: func(ctx *cli.Context) error {
					log.Println("Extracting client side modpack")

					valid, err := internal.ValidateFile(ctx.Args().First())
					if err != nil {
						return err
					}

					if !valid {
						return nil
					}

					return nil
				},
			},
		},
	}

	return app
}
