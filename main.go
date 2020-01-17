package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
	"tikv-client/console"

	"github.com/peterh/liner"
	"github.com/tikv/client-go/config"
	"github.com/tikv/client-go/rawkv"
)

var (
	pdAddr      = flag.String("addr", "127.0.0.1:2379", "pd-server address split by comma")
)

var (
	line *liner.State

	startKey = ""

	client *rawkv.Client
)

func main() {
	flag.Parse()

	line = liner.NewLiner()
	defer line.Close()
	line.SetCtrlCAborts(true)

	// Start interactive mode when no command is provided
	if flag.NArg() == 0 {
		repl()
	}
}

func cliConnect() {
	if client == nil {
		conf := config.Default()

		addrs := strings.Split(*pdAddr, ",")
		var err error
		client, err = rawkv.NewClient(context.Background(), addrs, conf)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(-1)
		}
	}
}

func repl() {
	console.LoadHistory(line)
	defer console.SaveHistory(line)

	reg, _ := regexp.Compile(`'.*?'|".*?"|\S+`)
	prompt := ""

	cliConnect()

	for {
		prompt = fmt.Sprintf("%s> ", *pdAddr)

		cmd, err := line.Prompt(prompt)
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			return
		}

		cmds := reg.FindAllString(cmd, -1)
		if len(cmds) == 0 {
			continue
		} else {
			console.AppendHistory(cmds, line)

			cmd := strings.ToLower(cmds[0])
			if cmd == "help" || cmd == "?" {
				console.PrintHelp(cmds)
			} else if cmd == "quit" || cmd == "exit" {
				os.Exit(0)
			} else {
				console.CliSendCmd(client, cmds)
			}
		}
	}
}
