/* SPDX-FileCopyrightText: 2026 (c) Joachim MARIE <moonstroke+github@live.fr>
   SPDX-License-Identifier: MIT */

package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

const DEFAULT_PORT = 10042
const BUFFER_SIZE = 256

var JOIN_MSG_FMT = "%s wants to join"
var ACCEPT_MSG_FMT = "%s welcomes %s"

var DEBUG = log.New(os.Stderr, "[DEBUG] ", log.Lshortfile)
var ERROR = log.New(os.Stderr, "[ERROR] ", log.LstdFlags|log.Lshortfile)

func usage(execName string) {
	os.Stderr.WriteString("Usage: " + execName + " ACTION NAME [PLAYER_COUNT]\n")
	os.Stderr.WriteString("Where:\n")
	os.Stderr.WriteString("\tACTION is either \"host\" or \"join\"\n")
	os.Stderr.WriteString("\tNAME is a non-empty string defining the player name\n")
	os.Stderr.WriteString("\tPLAYER_COUNT is a positive integer greater than 1 specifying the number\n")
	os.Stderr.WriteString("\t             of players for the game; mandatory if ACTION is \"host\",\n")
	os.Stderr.WriteString("\t             ignored otherwise\n")
	os.Exit(1)
}

type player struct {
	name string
	addr net.Addr
}

func (p player) String() string {
	return p.name + "@" + p.addr.String()
}

func main() {
	if len(os.Args) < 3 {
		usage(os.Args[0])
	}
	name := os.Args[2]
	if name == "" {
		DEBUG.Println("Empty name")
		usage(os.Args[0])
	}
	action := os.Args[1]
	switch action {
	case "host":
		if len(os.Args) == 3 {
			DEBUG.Println("Missing player count")
			usage(os.Args[0])
		}
		playerCount, err := strconv.Atoi(os.Args[3])
		if err != nil || playerCount <= 1 {
			DEBUG.Println("Invalid player count", os.Args[3])
			usage(os.Args[0])
		}

		hostGame(name, playerCount)
	case "join":
		joinGame(name)
	default:
		DEBUG.Println("Unknown action", action)
		usage(os.Args[0])
	}
}

func hostGame(name string, playerCount int) {
	players := make([]player, 0, playerCount)
	DEBUG.Println("Hosting game")
	remoteAddr := &net.UDPAddr{Port: DEFAULT_PORT}
	conn, err := net.ListenUDP("udp4", remoteAddr)
	if err != nil {
		ERROR.Fatalln(err)
	}
	DEBUG.Println(conn.LocalAddr(), "is connected to", conn.RemoteAddr())
	defer conn.Close()

	players = append(players, player{name, conn.LocalAddr()})
	buffer := make([]byte, BUFFER_SIZE)
	for {
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			ERROR.Fatalln(err)
		}
		msg := string(buffer[:n])
		DEBUG.Println("Received", msg, "from", addr)
		var playerName string
		if n, err = fmt.Sscanf(msg, JOIN_MSG_FMT, &playerName); err != nil {
			ERROR.Println(err)
			continue
		}
		if n == 1 {
			/* Message was a join request: accept player */
			DEBUG.Println("Acepting player", playerName)
			conn.WriteToUDP(fmt.Appendf(nil, ACCEPT_MSG_FMT, name, playerName), addr)
			players = append(players, player{playerName, addr})
			if len(players) == playerCount {
				break
			}
		}
	}
}

func joinGame(name string) {
	DEBUG.Println("Joining game")
	localAddr := &net.UDPAddr{Port: 0}
	remoteAddr := &net.UDPAddr{Port: DEFAULT_PORT}
	DEBUG.Println("local address =", localAddr, "; remote address =", remoteAddr)
	conn, err := net.DialUDP("udp4", localAddr, remoteAddr)
	if err != nil {
		ERROR.Fatalln(err)
	}
	DEBUG.Println(conn.LocalAddr(), "is connected to", conn.RemoteAddr())
	defer conn.Close()

	payload := fmt.Appendf(nil, JOIN_MSG_FMT, name)
	buffer := make([]byte, BUFFER_SIZE)
	/* Dirty hack, but the only way I found to format only one flag */
	replyFmt := fmt.Sprintf(ACCEPT_MSG_FMT, "%s", name)
	for {
		_, err := conn.Write(payload)
		if err != nil {
			ERROR.Fatalln(err)
		}
		conn.SetReadDeadline(time.Now().Add(time.Second))
		var n int
		if n, _, err = conn.ReadFromUDP(buffer); err != nil {
			ERROR.Println(err)
			continue
		}
		reply := string(buffer[:n])
		DEBUG.Println("Received", reply, "from", conn.RemoteAddr())
		var hostName string
		if n, err = fmt.Sscanf(reply, replyFmt, &hostName); err != nil {
			ERROR.Println(err)
			continue
		}
		if n == 1 {
			/* Host accepted our request */
			DEBUG.Println("Join request accepted by host", hostName, "@", conn.RemoteAddr())
			break
		}
	}
}
