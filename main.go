/* SPDX-FileCopyrightText: 2026 (c) Joachim MARIE <moonstroke+github@live.fr>
   SPDX-License-Identifier: MIT */

package main

import (
	"os"
)

func main() {
	/* If requested, set up a game; by default, look for one instead */
	if len(os.Args) > 1 && os.Args[1] == "host" {
		hostGame()
	} else {
		joinGame()
	}
}

func hostGame() {
	// TODO set up game and make available on network
}

func joinGame() {
	// TODO find game to join on network
}
