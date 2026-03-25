/* SPDX-FileCopyrightText: 2026 (c) Joachim MARIE <moonstroke+github@live.fr>
   SPDX-License-Identifier: MIT */

package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

const DEFAULT_ADDRESS = ":10042"
const BUFFER_SIZE = 256

var JOIN_MSG_FMT = "%s wants to join"

var DEBUG = log.New(os.Stderr, "[DEBUG] ", log.Lshortfile)
var ERROR = log.New(os.Stderr, "[ERROR] ", log.LstdFlags|log.Lshortfile)

func usage() {
	os.Stderr.WriteString("Usage: ggg [ACTION] NAME\n")
	os.Stderr.WriteString("Where:\n")
	os.Stderr.WriteString("\tACTION is either \"host\" or \"join\"\n")
	os.Stderr.WriteString("\tNAME is a non-empty string defining the player name\n")
	os.Exit(1)
}

func main() {
	/* If requested, set up a game; by default, look for one instead */
	if len(os.Args) != 3 {
		usage()
	} else {
		name := os.Args[2]
		if name == "" {
			DEBUG.Println("Empty name")
			usage()
		}
		action := os.Args[1]
		switch action {
		case "host":
			hostGame(name)
		case "join":
			joinGame(name)
		default:
			DEBUG.Println("Unknown action", action)
			usage()
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

func hostGame(name string) {
	DEBUG.Println("Hosting game")
	remoteAddr := getUDPAddr(DEFAULT_ADDRESS)
	conn, err := net.ListenUDP("udp4", remoteAddr)
	if err != nil {
		ERROR.Fatalln(err)
	}
	DEBUG.Println(conn.LocalAddr(), "is connected to", conn.RemoteAddr())
	defer conn.Close()

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
			conn.WriteToUDP([]byte("Welcome, "+playerName+"!"), addr)
			// TODO register {playerName, addr}
		}
	}
}

func joinGame(name string) {
	DEBUG.Println("Joining game")
	localAddr := getUDPAddr(":0")
	remoteAddr := getUDPAddr(DEFAULT_ADDRESS)
	DEBUG.Println("local address =", localAddr, "; remote address =", remoteAddr)
	conn, err := net.DialUDP("udp4", localAddr, remoteAddr)
	if err != nil {
		ERROR.Fatalln(err)
	}
	DEBUG.Println(conn.LocalAddr(), "is connected to", conn.RemoteAddr())
	defer conn.Close()

	payload := []byte(fmt.Sprintf(JOIN_MSG_FMT, name))
	buffer := make([]byte, BUFFER_SIZE)
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
		if reply == "Welcome, "+name+"!" {
			DEBUG.Println("Join request accepted by host", conn.RemoteAddr())
			break
		}
	}
}
