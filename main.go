/* SPDX-FileCopyrightText: 2026 (c) Joachim MARIE <moonstroke+github@live.fr>
   SPDX-License-Identifier: MIT */

package main

import (
	"os"
)

func main() {
	/* If requested, set up a game; by default, look for one instead */
	if len(os.Args) == 1 {
		joinGame()
	} else if os.Args[1] == "host" {
		hostGame()
	} else if os.Args[1] == "join" {
		joinGame()
	} else {
		os.Stderr.WriteString("Usage: ggg ACTION?\n")
		os.Stderr.WriteString("\twhere ACTION is either \"host\" or \"join\"\n")
		os.Stderr.WriteString("\tif ACTION is unspecified, the default is \"join\"\n")
		os.Exit(1)
	}
}

func hostGame() {
	// TODO set up game and make available on network
}

func joinGame() {
	// TODO find game to join on network
}
