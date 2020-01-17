package console

import (
	"fmt"
	"github.com/peterh/liner"
	"os"
	"path"
	"strings"
)

var (
	historyPath = path.Join(os.Getenv("HOME"), ".gotikvcli_history") // $HOME/.gorediscli_history
)

func LoadHistory(line *liner.State) {
	if f, err := os.Open(historyPath); err == nil {
		line.ReadHistory(f)
		f.Close()
	}
}

func SaveHistory(line *liner.State) {
	if f, err := os.Create(historyPath); err != nil {
		fmt.Printf("Error writing history file: %s", err.Error())
	} else {
		line.WriteHistory(f)
		f.Close()
	}
}

func AppendHistory(cmds []string, line *liner.State) {
	// make a copy of cmds
	cloneCmds := make([]string, len(cmds))
	for i, cmd := range cmds {
		cloneCmds[i] = cmd
	}

	line.AppendHistory(strings.Join(cloneCmds, " "))
}