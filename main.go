package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Morphyni/tas-cli/consts"
	"github.com/Morphyni/tas-cli/eula"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	listenSignals()

	app := cli.NewApp()
	app.Name = consts.CLI_MODULE_NAME
	app.Usage = consts.CLI_MODULE_USAGE
	app.Version = consts.CLI_VERSION
	app.Copyright = consts.CLI_COPYRIGHT_MESSAGE

	app.CommandNotFound = func(ctx *cli.Context, command string) {
		fmt.Fprintf(ctx.App.Writer, "Command '%v' does not exist. Type tibcli -h to list valid commands.\n", command)
	}

	app.Before = setLogLevel
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug, d",
			Usage: "Enable debug logging.",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:      "eula",
			Usage:     "Display the End User License Agreement (EULA).",
			ArgsUsage: " ", // space is needed, otherwise [arguments...] will be displayed in help
			Action: func(c *cli.Context) {
				fmt.Println(eula.Tas_eula)
			},
		},
		{
			Name:  "login",
			Usage: "Log in the user to server",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "username, u",
					Usage: "The username. If username is specified, the password option has to be specified as well.",
					Value: "<username>",
				},
				cli.StringFlag{
					Name:  "password, p",
					Usage: "the user's password",
					Value: "<password>",
				},
				cli.StringFlag{
					Name:  "org, o",
					Usage: "the organization name",
				},
				cli.StringFlag{
					Name:  "region, r",
					Usage: "Passing region in command line argument.",
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
	// app.Action = func(c *cli.Context) {
	// 	c.Args()
	// 	fmt.Println("boom displayed: Hello!")
	// 	fmt.Println()
	// }
	app.Run(os.Args)
}

// listenSignals listening the os interrupt signal like ctrl+c and do os.Exit
func listenSignals() {
	log.Debug("Listening system signals ...")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		// check whether the output device is a terminal or not
		if terminal.IsTerminal(int(os.Stdout.Fd())) {
			// Remember the state before any invocation on terminal like terminal.ReadPassword()
			log.Debug("You're a terminal.")
			state, err := terminal.GetState(int(os.Stdin.Fd()))
			if err != nil {
				log.Debug("terminal.GetState() failed, exit.")
				os.Exit(1)
			}

			if c != nil {
				<-c
				fmt.Println()
				// Restore the state before any invocation on terminal like terminal.ReadPassword(),
				// without this step the terminal runs the tibcli will make all the following input after exit() get invisible
				terminal.Restore(int(os.Stdin.Fd()), state)
				log.Debug("System interrupt signal received, exit.")
				//TODO utils.ExecAllShutdownHandlers()
				os.Exit(1)
			}
		} else {
			log.Debug("You're not a terminal.")
		}

	}()
}

func setLogLevel(c *cli.Context) error {
	if c.Bool("debug") {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
	return nil
}
