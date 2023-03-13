package sse

type UserConnected struct {
	UserCount int
}

type UserDisconnected struct {
	UserCount int
}
