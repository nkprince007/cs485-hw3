package main

import (
	"ethos/altEthos"
	"ethos/syscall"
	"log"
)

func init() {
	altEthos.LogToDirectory("application/ethosChatService")
	SetupChatRpcListChatRooms(listChatRooms)
}

func listChatRooms() (ChatRpcProcedure) {
	path := "/user/nobody/chatrooms"
	fd, status := altEthos.DirectoryOpen(path)
	if status != syscall.StatusOk {
		log.Printf("DirectoryOpen failed %v\n", status)
		altEthos.Exit(status)
	}
	defer altEthos.Close(fd)

	files, status := altEthos.SubFiles(path)
	if status != syscall.StatusOk {
		log.Printf("altEthos.SubFiles failed %v\n", status)
		altEthos.Exit(status)
	}

	chatRooms := make([]ChatRoom, len(files))

	for i := 0; i < len(files); i++ {
		status = altEthos.ReadStream(fd, &chatRooms[i])
		if status != syscall.StatusOk {
			break
		}
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
		log.Println(status)
	}
}
