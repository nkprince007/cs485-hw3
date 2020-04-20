package main

import (
	"ethos/altEthos"
	"ethos/syscall"
	"log"
	"path"
)

var chatRoomsDir = "/user/nobody/chatrooms"

func init() {
	altEthos.LogToDirectory("application/ethosChatService")
	SetupChatRpcListChatRooms(listChatRooms)
	SetupChatRpcCreateChatRoom(createChatRoom)
	SetupChatRpcBlacklistUser(blacklistUser)
	SetupChatRpcSelectChatRoom(selectChatRoom)
	SetupChatRpcGetMessages(getMessages)
	SetupChatRpcPostMessage(postMessage)
}

func checkUserPermissions(user User, name string) (ChatRoom, bool) {
	fd, status := altEthos.DirectoryOpen(chatRoomsDir)
	if status != syscall.StatusOk {
		log.Printf("DirectoryOpen failed %v\n", status)
		altEthos.Exit(status)
	}
	defer altEthos.Close(fd)

	var chatRoom ChatRoom
	status = altEthos.ReadVar(fd, name, &chatRoom)
	if status != syscall.StatusOk {
		log.Println("ReaVar failed:", chatRoom, status)
		altEthos.Exit(status)
	}

	for _, bUser := range chatRoom.BlacklistedUsers {
		if bUser == user {
			return chatRoom, false
		}
	}
	return chatRoom, true
}

func postMessage(msg Message) (ChatRpcProcedure) {
	log.Println("Received message:", msg)
	room := msg.ChatRoom
	dirPath := "/user/nobody/" + room.Name
	fd, status := altEthos.DirectoryOpen(dirPath)
	if status != syscall.StatusOk {
		log.Printf("DirectoryOpen failed %v\n", status)
		return &ChatRpcPostMessageReply{false, "No chat room found."}
	}
	defer altEthos.Close(fd)


	if _, ok := checkUserPermissions(msg.SentBy, msg.ChatRoom.Name); !ok {
		return &ChatRpcPostMessageReply{false, "User blacklisted."}
	}

	status = altEthos.WriteStream(fd, &msg)
	if status != syscall.StatusOk {
		log.Println("WriteStream failed: ", status, msg)
		return &ChatRpcPostMessageReply{false, "WriteStream failure."}
	}

	return &ChatRpcPostMessageReply{true, ""}
}

func getMessages(room ChatRoom, user User) (ChatRpcProcedure) {
	log.Println("GetMessages request received:", room.Name)
	messages := []Message{}
	dirPath := "/user/nobody/" + room.Name
	fd, status := altEthos.DirectoryOpen(dirPath)
	if status != syscall.StatusOk {
		log.Printf("DirectoryOpen failed %v\n", status)
		return &ChatRpcGetMessagesReply{messages}
	}
	defer altEthos.Close(fd)

	if _, ok := checkUserPermissions(user, room.Name); !ok {
		return &ChatRpcGetMessagesReply{nil}
	}

	for {
		msg := Message{}
		status = altEthos.ReadStream(fd, &msg)
		if status != syscall.StatusOk {
			break
		}
		messages = append(messages, msg)
	}

	return &ChatRpcGetMessagesReply{messages}
}

func listChatRooms(user User) (ChatRpcProcedure) {
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

		// hide blacklisted chatrooms
		for _, bUser := range chatRoom.BlacklistedUsers {
			if bUser == user {
				continue
			}
		}
		chatRooms = append(chatRooms, chatRoom)
	}

	return &ChatRpcListChatRoomsReply{chatRooms}
}

func selectChatRoom(name string, user User) (ChatRpcProcedure) {
	chatRoom, ok := checkUserPermissions(user, name);
	return &ChatRpcSelectChatRoomReply{chatRoom, ok}
}

func createChatRoom(owner User, name string) (ChatRpcProcedure) {
	log.Printf("createChatRoom request received for '%s' from '%s' \n", name, owner)
	fd, status := altEthos.DirectoryOpen(chatRoomsDir)
	if status != syscall.StatusOk {
		log.Printf("DirectoryOpen failed %v\n", status)
		altEthos.Exit(status)
	}
	defer altEthos.Close(fd)

	if altEthos.IsFile(path.Join(chatRoomsDir, name)) {
		log.Printf("Chatroom %s already exists\n", name)
		failedChatRoom := ChatRoom{[]User{}, owner, name}
		return &ChatRpcCreateChatRoomReply{failedChatRoom, false}
	}

	chatRoom := ChatRoom{[]User{}, owner, name}
	status = altEthos.WriteVar(fd, name, &chatRoom)
	if status != syscall.StatusOk {
		log.Println("WriteStream failed: ", chatRoom, status)
		altEthos.Exit(status)
	}

	dirPath := "/user/nobody/" + name
	var msg Message
	status = altEthos.DirectoryCreate(dirPath, &msg, "all")
	if status != syscall.StatusOk {
		log.Println("CreateDirectory failed:", chatRoom, status)
		altEthos.Exit(status)
	}

	log.Printf("Chatroom created: %s\n", chatRoom.Name)
	return &ChatRpcCreateChatRoomReply{chatRoom, true}
}

func blacklistUser(roomName string, user User) (ChatRpcProcedure) {
	filePath := path.Join(chatRoomsDir, roomName)
	if !altEthos.IsFile(filePath) {
		log.Printf("Chatroom [%s] does not exist", roomName)
		return &ChatRpcBlacklistUserReply{false}
	}

	var chatRoom ChatRoom
	status := altEthos.Read(filePath, &chatRoom)
	if status != syscall.StatusOk {
		log.Println("Read failed: ", chatRoom, status)
		altEthos.Exit(status)
	}

	chatRoom.BlacklistedUsers = append(chatRoom.BlacklistedUsers, user)
	status = altEthos.Write(filePath, &chatRoom)
	if status != syscall.StatusOk {
		log.Println("Write failed: ", chatRoom, status)
		altEthos.Exit(status)
	}

	return &ChatRpcBlacklistUserReply{true}
}

func main() {
	log.Println("ethosChatService started...")
	listeningFd, status := altEthos.Advertise("ethosChat")
	if status != syscall.StatusOk {
		log.Println("Advertising service failed: ", status)
		altEthos.Exit(status)
	}
	log.Println("ethosChatService advertised: ", status)

	event, status := altEthos.ImportAsync(listeningFd, &ChatRpc{}, altEthos.HandleImport)
	if status != syscall.StatusOk {
		log.Println("Import failed:", status)
		altEthos.Exit(status)
	}

	var tree altEthos.EventTreeSlice
	next := []syscall.EventId{event}

	for {
		tree = altEthos.WaitTreeCreateOr(next)
		tree, _ = altEthos.Block(tree)
		completed, pending := altEthos.GetTreeEvents(tree)
		for _, eventId := range completed {
			eventInfo, status := altEthos.OnComplete(eventId)
			if status != syscall.StatusOk {
				log.Println("OnComplete_failed", eventInfo, status)
				return
			}
			eventInfo.Do()
		}

		next = nil
		next = append(next, pending...)
		next = append(next, altEthos.RetrievePostedEvents()...)
	}
}
