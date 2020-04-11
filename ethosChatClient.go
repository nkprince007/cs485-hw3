package main

import (
	"ethos/altEthos"
	"ethos/syscall"
	// "bufio"
	"fmt"
	// "os"
	// "regexp"
	// "strings"
	"log"
)

func init() {
	altEthos.LogToDirectory("application/ethosChatClient")
	SetupChatRpcListChatRoomsReply(listChatRoomsReply)
}

func listChatRoomsReply(rooms []ChatRoom) (ChatRpcProcedure) {
	log.Println("Received listChatRoom reply: ", rooms)
	fmt.Println(rooms)
	return nil
}

func main() {
	log.Println("ethosChatClient started")
	fd, status := altEthos.IpcRepeat("ethosChat", "", nil)
	if status != syscall.StatusOk {
		log.Println("Ipc failed: ", status)
		altEthos.Exit(status)
	}
	defer altEthos.Close(fd)

	// reader := bufio.NewReader(os.Stdin)
	// fmt.Println("Ethos Chat")
	// fmt.Println("---------------------")

	// for {
	// 	fmt.Print("? ")
	// 	text, _ := reader.ReadString('\n')
	// 	text = strings.TrimSpace(text)

	// 	if ok, _ := regexp.MatchString(`> list`, text); ok {
	// 		call := &ChatRpcListChatRooms{}
	// 		status = altEthos.ClientCall(fd, call)
	// 		if status != syscall.StatusOk {
	// 			log.Println("clientCall failed: ", status)
	// 			altEthos.Exit(status)
	// 		}
	// 	}
	// }

}
