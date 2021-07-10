package main

import (
	"github.com/urfave/cli/v2"
	"github.com/wongpinter/webnovel-scraper/webnovel"
	"log"
	"os"
	"sort"
)

func main() {
	myFlags := []cli.Flag{
		&cli.StringFlag{
			Name: "webnovel-scraper",
			Value: "",
			Usage: "Scraping webnovel from https://readnovelfull.com",
		},
	}
	app := &cli.App{
		Name: "tes",
		Usage: "Example: webnovel-scraper --url https://readnovelfull.com/",
		Flags: myFlags,
		Commands: []*cli.Command{
			{
				Name: "all",
				Usage: "Download all chapters",
				Action: func(c *cli.Context) error {
					s := new(webnovel.Scraper)
					s.Fetch("https://readnovelfull.com/if-you-dont-fall-in-love-youll-die.html")
					c.String("host")
					return nil
				},
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
