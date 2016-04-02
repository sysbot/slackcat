package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/codegangsta/cli"
	"github.com/fatih/color"
)

var version = "dev-build"

func readIn(lines chan string, tee bool) {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		lines <- scanner.Text()
		if tee {
			fmt.Println(scanner.Text())
		}
	}
	fmt.Println("just testing")
	fmt.Println("new branch")

	close(lines)
}

func writeTemp(lines chan string) string {
	tmp, err := ioutil.TempFile(os.TempDir(), "slackcat-")
	failOnError(err, "unable to create tmpfile", false)

	w := bufio.NewWriter(tmp)
	for line := range lines {
		fmt.Fprintln(w, line)
	}
	w.Flush()

	return tmp.Name()
}

func output(s string) {
	cyan := color.New(color.Bold).SprintFunc()
	fmt.Printf("%s %s\n", cyan("slackcat"), s)
}

func failOnError(err error, msg string, appendErr bool) {
	if err != nil {
		if appendErr {
			exitErr(fmt.Errorf("%s: %s", msg, err))
		} else {
			exitErr(fmt.Errorf("%s", msg))
		}
	}
}

func exitErr(err error) {
	output(color.RedString(err.Error()))
	os.Exit(1)
}

func main() {
	app := cli.NewApp()
	app.Name = "slackcat"
	app.Usage = "redirect a file to slack"
	app.Version = version
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "tee, t",
			Usage: "Print stdin to screen before posting",
		},
		cli.BoolFlag{
			Name:  "stream, s",
			Usage: "Stream messages to Slack continuously instead of uploading a single snippet",
		},
		cli.BoolFlag{
			Name:  "plain, p",
			Usage: "Write messages as plain texts instead of code blocks",
		},
		cli.BoolFlag{
			Name:  "noop",
			Usage: "Skip posting file to Slack. Useful for testing",
		},
		cli.BoolFlag{
			Name:  "configure",
			Usage: "Configure Slackcat via oauth",
		},
		cli.StringFlag{
			Name:  "channel, c",
			Usage: "Slack channel or group to post to",
		},
		cli.StringFlag{
			Name:  "filename, n",
			Usage: "Filename for upload. Defaults to current timestamp",
		},
	}

	app.Action = func(c *cli.Context) {
		if c.Bool("configure") {
			configureOA()
			os.Exit(0)
		}

		config := readConfig()
		fileName := c.String("filename")
		team, channel, err := config.parseChannelOpt(c.String("channel"))
		failOnError(err, "", true)

		token := config.teams[team]
		if token == "" {
			exitErr(fmt.Errorf("no such team: %s", team))
		}

		if !c.Bool("stream") && c.Bool("plain") {
			exitErr(fmt.Errorf("'plain' flag requires 'stream' mode!"))
		}

		slackcat, err := newSlackCat(token, channel)
		failOnError(err, "Slack API Error", true)

		if len(c.Args()) > 0 {
			if c.Bool("stream") {
				output("filepath provided, ignoring stream option")
			}
			filePath := c.Args()[0]
			if fileName == "" {
				fileName = filepath.Base(filePath)
			}
			slackcat.postFile(filePath, fileName, c.Bool("noop"))
			os.Exit(0)
		}

		lines := make(chan string)
		go readIn(lines, c.Bool("tee"))

		if c.Bool("stream") {
			output("starting stream")
			go slackcat.addToStreamQ(lines)
			go slackcat.processStreamQ(c.Bool("noop"), c.Bool("plain"))
			go slackcat.trap()
			select {}
		} else {
			filePath := writeTemp(lines)
			defer os.Remove(filePath)
			slackcat.postFile(filePath, fileName, c.Bool("noop"))
			os.Exit(0)
		}
	}

	app.Run(os.Args)

}
