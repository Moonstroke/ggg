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

func usage() {
	os.Stderr.WriteString("Usage: " + os.Args[0] + ` ACTION NAME [PLAYER_COUNT]
Where:
	ACTION is either "host" or "join"
	NAME is a non-empty string defining the player name
	PLAYER_COUNT is a positive integer greater than 1 specifying the number
	             of players for the game; mandatory if ACTION is "host",
	             ignored otherwise
`)
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
		usage()
	}
	name := os.Args[2]
	if name == "" {
		DEBUG.Println("Empty name")
		usage()
	}
	action := os.Args[1]
	switch action {
	case "host":
		if len(os.Args) == 3 {
			DEBUG.Println("Missing player count")
			usage()
		}
		playerCount, err := strconv.Atoi(os.Args[3])
		if err != nil || playerCount <= 1 {
			DEBUG.Println("Invalid player count", os.Args[3])
			usage()
		}

		hostGame(name, playerCount)
	case "join":
		joinGame(name)
	default:
		DEBUG.Println("Unknown action", action)
		usage()
	}
}

func recvJoinRequest(conn *net.UDPConn, buffer []byte) (string, *net.UDPAddr) {
	n, addr, err := conn.ReadFromUDP(buffer)
	if err != nil {
		ERROR.Fatalln(err)
	}
	msg := string(buffer[:n])
	DEBUG.Println("Received", msg, "from", addr)
	var playerName string
	if n, err = fmt.Sscanf(msg, JOIN_MSG_FMT, &playerName); err != nil {
		ERROR.Println(err)
		return "", nil
	}
	if n != 1 {
		return "", nil
	}
	return playerName, addr
}

func sendJoinAck(conn *net.UDPConn, addr *net.UDPAddr, name, playerName string) {
	conn.WriteToUDP(fmt.Appendf(nil, ACCEPT_MSG_FMT, name, playerName), addr)
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
		playerName, addr := recvJoinRequest(conn, buffer)
		if addr == nil {
			continue
		}
		DEBUG.Println("Acepting player", playerName)
		sendJoinAck(conn, addr, name, playerName)
		players = append(players, player{playerName, addr})
		if len(players) == playerCount {
			break
		}
	}
}

func sendJoinRequest(conn *net.UDPConn, name string) {
	payload := fmt.Appendf(nil, JOIN_MSG_FMT, name)
	if _, err := conn.Write(payload); err != nil {
		ERROR.Fatalln(err)
	}
}

func recvJoinAck(conn *net.UDPConn, buffer []byte, replyFmt string) string {
	var n int
	var err error
	if n, _, err = conn.ReadFromUDP(buffer); err != nil {
		ERROR.Println(err)
		return ""
	}
	reply := string(buffer[:n])
	DEBUG.Println("Received", reply, "from", conn.RemoteAddr())
	var hostName string
	if n, err = fmt.Sscanf(reply, replyFmt, &hostName); err != nil {
		ERROR.Println(err)
		return ""
	}
	if n == 1 {
		/* Host accepted our request */
		DEBUG.Println("Join request accepted by host", hostName, "@", conn.RemoteAddr())
		return hostName
	}
	return ""
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

	buffer := make([]byte, BUFFER_SIZE)
	/* Dirty hack, but the only way I found to format only one flag */
	replyFmt := fmt.Sprintf(ACCEPT_MSG_FMT, "%s", name)
	for {
		sendJoinRequest(conn, name)
		conn.SetReadDeadline(time.Now().Add(time.Second))
		hostName := recvJoinAck(conn, buffer, replyFmt)
		if hostName != "" {
			break
		}
	}
}
