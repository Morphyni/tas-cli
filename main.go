package main

import (
	"fmt"
	"os"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "boom"
	app.Usage = "make an explosive entrance"
	app.Version = "1.2.3"
	app.Copyright = "2018 @boom.go"
	app.Commands = []cli.Command{
		{
			Name:  "eula",
			Usage: "Display the end use license agreement",
		},
		{
			Name:  "login",
			Usage: "Log in the user to server",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "password, p",
					Usage: "the user's password",
					Value: "<password>",
				},
				cli.StringFlag{
					Name:  "org, o",
					Usage: "the organization name",
				},
			},
			Action: func(c *cli.Context) {
				nflags := c.NumFlags()
				arg1 := c.String("password")
				arg2 := c.String("o")
				fmt.Println("Args: ", arg1, arg2)
				fmt.Println("Nflags: ", nflags)
			},
		},
		{
			Name:  "list",
			Usage: "List all elements",
			Action: func(c *cli.Context) {
				fmt.Println("List all elements")
			},
			Subcommands: []cli.Command{
				{
					Name:  "users",
					Usage: "Display all users",
					Action: func(c *cli.Context) {
						fmt.Println("List all users")
					},
				},
				{
					Name:  "orgs",
					Usage: "Display all orgs",
					Action: func(c *cli.Context) {
						fmt.Println("List all orgs")
					},
				},
			},
		},
	}
	app.Action = func(c *cli.Context) {
		c.Args()
		fmt.Println("boom displayed: Hello!")
		fmt.Println()
	}
	app.Run(os.Args)
}
