package main

import (
	"ethos/altEthos"
	"ethos/syscall"
	"log"
)

var chatRoomsDir = "/user/nobody/chatrooms"

func init() {
	altEthos.LogToDirectory("application/ethosChatService")
	SetupChatRpcListChatRooms(listChatRooms)
	SetupChatRpcCreateChatRoom(createChatRoom)
}

func createChatRoom(owner User, name string) (ChatRpcProcedure) {
	log.Printf("createChatRoom request received for '%s' from '%s' \n", name, owner)
	fd, status := altEthos.DirectoryOpen(chatRoomsDir)
	if status != syscall.StatusOk {
		log.Printf("DirectoryOpen failed %v\n", status)
		altEthos.Exit(status)
	}
	defer altEthos.Close(fd)

	chatRoom := ChatRoom{[]User{}, owner, name}
	status = altEthos.WriteStream(fd, &chatRoom)
	if status != syscall.StatusOk {
		log.Println("WriteStream failed: ", chatRoom, status)
		altEthos.Exit(status)
	}

	log.Printf("Chatroom created: %s\n", chatRoom.Name)
	return &ChatRpcCreateChatRoomReply{chatRoom}
}

func listChatRooms() (ChatRpcProcedure) {
	fd, status := altEthos.DirectoryOpen(chatRoomsDir)
	if status != syscall.StatusOk {
		log.Printf("DirectoryOpen failed %v\n", status)
		altEthos.Exit(status)
	}
	defer altEthos.Close(fd)

	chatRooms := []ChatRoom{}

	for {
		chatRoom := ChatRoom{}
		status = altEthos.ReadStream(fd, &chatRoom)
		if status != syscall.StatusOk {
			break
		}
		chatRooms = append(chatRooms, chatRoom)
	}

	return &ChatRpcListChatRoomsReply{chatRooms}
}

func main() {
	log.Println("ethosChatService started...")
	listeningFd, status := altEthos.Advertise("ethosChat")
	if status != syscall.StatusOk {
		log.Println("Advertising service failed: ", status)
		altEthos.Exit(status)
	}
	log.Println("ethosChatService advertised: ", status)

	for {
		_, fd, status := altEthos.Import(listeningFd)
		if status != syscall.StatusOk {
			log.Println("Error importing service: ", fd, status)
			altEthos.Exit(status)
		}

		log.Println("Accepted connection: ", status)
		t := ChatRpc{}
		status = altEthos.Handle(fd, &t)
	}
}
