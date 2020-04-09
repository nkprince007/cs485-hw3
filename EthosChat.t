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

Chat interface {
    ListChatRooms() []ChatRoom
    SelectChatRoom(ChatRoom) bool
    CreateChatRoom(string) ChatRoom
    BlacklistUser(ChatRoom, User)
}
