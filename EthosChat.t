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
}

ChatRpc interface {
    ListChatRooms() (rooms []ChatRoom)
    CreateChatRoom(owner User, name string) (room ChatRoom, status bool)
    BlacklistUser(roomName string, user User) (status bool)
}
