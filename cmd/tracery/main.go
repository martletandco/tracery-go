package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/martletandco/tracery-go"
)

// Flags to add
// Flatten expression (default #origin#)
// Random seed
// Rules to add (`symbol:value` format)

func main() {
	g := tracery.NewGrammar()

	readInRuleSet(g)

	r := g.Flatten("#origin#")

	os.Stdout.WriteString(r)
	os.Stdout.WriteString("\n")
	os.Exit(0)
}

func readInRuleSet(g tracery.Grammar) {
	fi, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	if fi.Mode()&os.ModeNamedPipe == 0 {
		return
	}

	var buf []byte
	buf, err = ioutil.ReadAll(os.Stdin)
	if err != nil {
		// @cleanup: Add better error
		bail(err)
		return
	}

	var set tracery.RuleSet
	if err := json.Unmarshal(buf, &set); err != nil {
		// @cleanup: Add better error
		bail(err)
	}

	g.PushRuleSet(set)
}

func bail(err error) {
	fmt.Fprintln(os.Stderr, "error:", err)
	os.Exit(1)
}
