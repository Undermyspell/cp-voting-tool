package dtos

type NewQuestionDto struct {
	Text string `json:"text" binding:"required"`
}

type QuestionDto struct {
	Id       string
	Text     string
	Votes    int
	Answered bool
}
