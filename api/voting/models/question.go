package voting_models

const (
	MaxLength int = 500
)

type Question struct {
	Id          string
	Text        string
	Votes       int
	Answered    bool
	Voted       bool
	CreatorHash string
	CreatorName string
	Anonymous   bool
}

func NewQuestion(id, text string, votes int, answered, voted, anonymous bool, creatorName, creatorHash string) Question {
	if anonymous {
		creatorName = ""
	}
	return Question{
		Id:          id,
		Text:        text,
		Votes:       votes,
		Answered:    answered,
		Voted:       voted,
		CreatorHash: creatorHash,
		CreatorName: creatorName,
		Anonymous:   anonymous,
	}
}
