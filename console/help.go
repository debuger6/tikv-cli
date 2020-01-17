package console

import (
	"fmt"
	"strings"
)

var (
	cmdMap = map[string]string{
		"GET":   "GET key",
		"SET":   "SET key value",
		"DEL":   "DEL key [key ...]",
		"SCAN":  "SCAN start_key end_key limit",
		"RSCAN": "RSCAN start_key end_key limit",
		"KEYS":  "KEYS pattern",
		"IT": "it",
	}
)

func printGenericHelp() {
	msg :=
		`tikv-cli
Type:	"help <command>" for help on <command>
Commands:
        GET
        SET
        DEL
        SCAN
        RSCAN
        KEYS
	`
	fmt.Println(msg)
}

func printCommandHelp(cmd string) {
	fmt.Println()
	usage, ok := cmdMap[cmd]
	if !ok {
		fmt.Printf("cmd %s is missing, please checkout again\n\n", cmd)
		return
	}

	fmt.Println(usage)
	fmt.Println()
}

func PrintHelp(cmds []string) {
	args := cmds[1:]
	if len(args) == 0 {
		printGenericHelp()
	} else if len(args) > 1 {
		fmt.Println()
	} else {
		cmd := strings.ToUpper(args[0])
		printCommandHelp(cmd)
	}
}