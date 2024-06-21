package cli

import (
	"log"

	"github.com/anthonyporthouse/mr-packer/internal/modrinth"
	"github.com/spf13/afero"
	"github.com/urfave/cli/v2"
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

					appFs := afero.NewOsFs()

					archivePath := ctx.Args().First()

					manifest, err := modrinth.ValidateFile(archivePath, appFs)
					if err != nil {
						return err
					}

					modrinth.ExtractOverrides(archivePath, modrinth.Client, appFs)

					err = modrinth.DownloadFiles(manifest, modrinth.Client, appFs)
					if err != nil {
						return err
					}

					return nil
				},
			},
			{
				Name:    "server",
				Aliases: []string{"s"},
				Usage:   "Extract a server side modpack",
				Action: func(ctx *cli.Context) error {
					log.Println("Extracting server side modpack")

					appFs := afero.NewOsFs()

					manifest, err := modrinth.ValidateFile(ctx.Args().First(), appFs)
					if err != nil {
						return err
					}

					err = modrinth.DownloadFiles(manifest, modrinth.Server, appFs)
					if err != nil {
						return err
					}

					return nil
				},
			},
		},
	}

	return app
}
