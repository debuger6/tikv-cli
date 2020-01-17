package console

import (
	"context"
	"fmt"
	"github.com/tikv/client-go/rawkv"
	"strconv"
	"strings"
)

var (
	startKey = ""
)

func verifyCmd(cmd string) bool {
	if _, ok := cmdMap[cmd]; !ok {
		fmt.Printf("\n%s is invalid, please check again\n", cmd)
		return false
	}
	return true
}

func CliSendCmd(client *rawkv.Client, cmds []string) {
	cmd := strings.ToUpper(cmds[0])
	if !verifyCmd(cmd) {
		return
	}

	args := cmds[1:]

	switch cmd {
	case "GET":
		if len(args) != 1 {
			fmt.Printf("usage: %s\n", cmdMap[cmd])
			return
		}
		r, e := client.Get(context.Background(), []byte(args[0]))
		if e != nil {
			fmt.Println(e)
			return
		}
		fmt.Println(string(r))
	case "SET":
		if len(args) != 2 {
			fmt.Printf("usage: %s\n", cmdMap[cmd])
			return
		}
		if err := client.Put(context.Background(), []byte(args[0]), []byte(args[1])); err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("ok.")
	case "DEL":
		if len(args) != 1 {
			fmt.Printf("usage: %s\n", cmdMap[cmd])
			return
		}
		if err := client.Delete(context.Background(), []byte(args[0])); err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("ok.")
	case "SCAN":
		if len(args) != 3 {
			fmt.Printf("usage: %s\n", cmdMap[cmd])
			return
		}
		limit, err := strconv.Atoi(args[2])
		if err != nil {
			fmt.Println(err)
			return
		}
		keys, values, err := client.Scan(context.Background(), []byte(args[0]), []byte(args[1]), limit)
		if err != nil {
			fmt.Println(err)
			return
		}
		lines := len(keys)
		fl := len(strconv.Itoa(lines))
		formatStr := fmt.Sprintf("%%%ds) ", fl)
		for i := 0; i < lines; i++ {
			num := strconv.Itoa(i+1)
			fmt.Printf(formatStr, num)
			fmt.Println(string(keys[i]), ":", string(values[i]))
		}
	case "KEYS":
		if len(args) != 1 {
			fmt.Printf("usage: %s\n", cmdMap[cmd])
			return
		}
		keys, _, err := client.Scan(context.Background(), []byte(args[0]), nil, 20)
		if err != nil {
			fmt.Println(err)
		}
		lines := len(keys)
		if lines == 0 {
			return
		}
		fl := len(strconv.Itoa(lines))
		formatStr := fmt.Sprintf("%%%ds) ", fl)
		for i := 0; i < lines; i++ {
			num := strconv.Itoa(i+1)
			fmt.Printf(formatStr, num)
			fmt.Println(string(keys[i]))
		}
		startKey = string(keys[lines-1])
		fmt.Println(`Type "it" for more`)
	case "IT":
		if startKey == "" {
			return
		}
		keys, _, err := client.Scan(context.Background(), []byte(startKey), nil, 20)
		if err != nil {
			fmt.Println(err)
		}
		lines := len(keys)
		if lines <= 1 {
			return
		}
		fl := len(strconv.Itoa(lines))
		formatStr := fmt.Sprintf("%%%ds) ", fl)
		for i := 1; i < lines; i++ {
			num := strconv.Itoa(1)
			fmt.Printf(formatStr, num)
			fmt.Println(string(keys[i]))
		}
		startKey = string(keys[lines-1])
		fmt.Println(`Type "it" for more`)
	}
}