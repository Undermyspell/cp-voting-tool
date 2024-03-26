package voting_usecases_events

type UserConnected struct {
	UserCount int
}

type UserDisconnected struct {
	UserCount int
}

type QuestionUpvoted struct {
	Id    string
	Votes int
	Voted bool
}

type QuestionAnswered struct {
	Id string
}

type QuestionDeleted struct {
	Id string
}

type QuestionUpdated struct {
	Id        string
	Text      string
	Creator   string
	Anonymous bool
}

type QuestionCreated struct {
	Id        string
	Text      string
	Votes     int
	Voted     bool
	Answered  bool
	Creator   string
	Anonymous bool
	Owned     bool
}

const PayloadEmpty = "{}"
