User string

ChatRoom struct {
    BlacklistedUsers []User
    Owner User
    Name string
}

Message struct {
    ChatRoom ChatRoom
    SentBy User
    CreatedAt int64
    Content string
}

ChatRpc interface {
    ListChatRooms() (rooms []ChatRoom)
    CreateChatRoom(owner User, name string) (room ChatRoom, status bool)
    BlacklistUser(roomName string, user User) (status bool)
    SelectChatRoom(name string, user User) (room ChatRoom, status bool)
    GetMessages(room ChatRoom) (messages []Message)
    PostMessage(msg Message) (status bool, issue string)
}
