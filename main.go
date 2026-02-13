/* SPDX-FileCopyrightText: 2026 (c) Joachim MARIE <moonstroke+github@live.fr>
   SPDX-License-Identifier: MIT */

package main

import (
	"log"
	"net"
	"os"
	"time"
)

const DEFAULT_ADDRESS = ":10042"
const BUFFER_SIZE = 256

var DEBUG = log.New(os.Stderr, "[DEBUG] ", log.Lshortfile)
var ERROR = log.New(os.Stderr, "[ERROR] ", log.LstdFlags|log.Lshortfile)

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
			os.Stderr.WriteString("Usage: ggg [ACTION]\n")
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
	remoteAddr := getUDPAddr(DEFAULT_ADDRESS)
	conn, err := net.ListenUDP("udp4", remoteAddr)
	if err != nil {
		ERROR.Fatalln(err)
	}
	DEBUG.Println(conn.LocalAddr(), "is connected to", conn.RemoteAddr())
	defer conn.Close()

	buffer := make([]byte, BUFFER_SIZE)
	for {
		n, addr, err := conn.ReadFrom(buffer)
		if err != nil {
			ERROR.Fatalln(err)
		}
		DEBUG.Println("Received", string(buffer[:n]), "from", addr)
		// TODO if addr is a player, register and notify it
	}
}

func joinGame() {
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

	payload := []byte("I wanna play!")
	for {
		_, err := conn.Write(payload)
		if err != nil {
			ERROR.Fatalln(err)
		}
		time.Sleep(1 * time.Second)
		// TODO if host replies, stop sending
	}
}
