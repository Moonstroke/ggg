/* SPDX-FileCopyrightText: 2026 (c) Joachim MARIE <moonstroke+github@live.fr>
   SPDX-License-Identifier: MIT */

package main

import (
	"log"
	"net"
	"os"
)

var DEBUG = log.New(os.Stderr, "[DEBUG] ", log.Lshortfile)
var ERROR = log.New(os.Stderr, "[ERROR] ", log.LstdFlags | log.Lshortfile)

func main() {
	/* If requested, set up a game; by default, look for one instead */
	if len(os.Args) == 1 {
		joinGame()
	} else {
		action := os.Args[1]
		switch action {
		case "host":
			hostGame()
		case "join":
			joinGame()
		default:
			DEBUG.Println("Unknown action", action)
			os.Stderr.WriteString("Usage: ggg ACTION?\n")
			os.Stderr.WriteString("\twhere ACTION is either \"host\" or \"join\"\n")
			os.Stderr.WriteString("\tif ACTION is unspecified, the default is \"join\"\n")
			os.Exit(1)
		}
	}
}

func getUDPAddr(address string) *net.UDPAddr {
	udpAddr, err := net.ResolveUDPAddr("udp4", address)
	if err != nil {
		ERROR.Fatalln(err)
	}
	return udpAddr
}

func hostGame() {
	DEBUG.Println("Hosting game")
	// TODO set up game and make available on network
}

func joinGame() {
	DEBUG.Println("Joining game")
	// TODO find game to join on network
}
