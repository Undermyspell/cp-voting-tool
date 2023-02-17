package dtos

type NewQuestionDto struct {
	Text      string `json:"text" binding:"required"`
	Anonymous bool   `json:"anonymous"`
}

type QuestionDto struct {
	Id        string
	Text      string
	Votes     int
	Answered  bool
	Creator   string
	Anonymous bool
	Owned     bool
}
