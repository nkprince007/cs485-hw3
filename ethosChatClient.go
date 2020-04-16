package main

import (
	"ethos/altEthos"
	"ethos/syscall"
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"log"
)

var owner = User(altEthos.GetUser())

func init() {
	altEthos.LogToDirectory("application/ethosChatClient")
	SetupChatRpcListChatRoomsReply(listChatRoomsReply)
	SetupChatRpcCreateChatRoomReply(createChatRoomReply)
}

func listChatRoomsReply(rooms []ChatRoom) (ChatRpcProcedure) {
	log.Println("Received listChatRoom reply: ", rooms)
	if len(rooms) > 0 {
		for _, room := range(rooms) {
			fmt.Printf("[%s]\n", room.Name)
		}
	} else {
		fmt.Println("No chatrooms found.")
	}
	return nil
}

func createChatRoomReply(room ChatRoom, status bool) (ChatRpcProcedure) {
	log.Println("Received createChatRoom reply: ", room, status)
	if status {
		fmt.Println("New chatroom created: ", room.Name)
	} else {
		fmt.Printf("Failed creating room %s, because a room with name already exists\n", room.Name)
	}
	return nil
}

func checkRpcStatus(status syscall.Status) {
	if status != syscall.StatusOk {
		log.Println("clientCall failed: ", status)
		altEthos.Exit(status)
	}
}

func printUsage() {
	fmt.Println("All commands start with a > sign. Please use it responsibly.")
	fmt.Println("> list\t\t- Get list of channels")
	fmt.Println("> help\t\t- Show help info")
	fmt.Println("> create <name>\t- Create a channel with given name")
	fmt.Println("> quit\t\t- Exit application")
}

func main() {
	log.Println("ethosChatClient started")
	fd, status := altEthos.IpcRepeat("ethosChat", "", nil)
	if status != syscall.StatusOk {
		log.Println("Ipc failed: ", status)
		altEthos.Exit(status)
	}
	defer altEthos.Close(fd)

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Ethos Chat")
	fmt.Println(strings.Repeat("-", 20))
	printUsage()

	for {
		fmt.Print("? ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)

		if ok, _ := regexp.MatchString(`> list`, text); ok {
			log.Println("Sending listChatRoom request.")
			call := &ChatRpcListChatRooms{}
			status = altEthos.ClientCall(fd, call)
			log.Println("reached here")
			checkRpcStatus(status)
		} else if ok, _ := regexp.MatchString(`> create [A-Za-z0-9_\-]+`, text); ok {
			name := strings.TrimSpace(strings.Split(text, " ")[2])
			log.Println("Sending createChatRoom request: ", name)
			call := &ChatRpcCreateChatRoom{owner, name}
			status = altEthos.ClientCall(fd, call)
			checkRpcStatus(status)
		} else if ok, _ := regexp.MatchString(`> quit`, text); ok {
			log.Println("Quitting application.")
			altEthos.Exit(syscall.StatusOk)
		} else if ok, _ := regexp.MatchString(`> help`, text); ok {
			printUsage()
		}
	}
}
