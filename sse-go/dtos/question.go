package dtos

type NewQuestionDto struct {
	Text string `json:"text" binding:"required"`
}
