package voting_repositories

import (
	"context"
	shared "voting/shared/helper"
	voting_models "voting/voting/models"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oklog/ulid/v2"
	"github.com/sirupsen/logrus"
)

type Postgresql struct {
	pool      *pgxpool.Pool
	sessionId string
}

func (session *Postgresql) Start() {
	_, err := session.pool.Exec(context.Background(), "INSERT INTO Sessions (id, sessionSecret) VALUES ($1, $2)", session.sessionId, shared.GetRandomId(30))

	if err != nil {
		logrus.Error(err.Error())
		return
	}
}

func (session *Postgresql) Stop() {
	_, err := session.pool.Exec(context.Background(), "DELETE FROM Sessions WHERE id = $1", session.sessionId)

	if err != nil {
		logrus.Error(err.Error())
		return
	}
}

func (session *Postgresql) GetQuestion(id string) (voting_models.Question, bool) {
	var question voting_models.Question
	err := session.pool.QueryRow(context.Background(),
		"SELECT q.id,q.text,q.answered,q.anonymous,q.creatorName,q.creatorHash,COUNT(uv.questionid) as votes "+
			"FROM questions q LEFT JOIN uservotes uv on q.id = uv.questionid "+
			"WHERE id = $1 "+
			"GROUP BY q.id,q.text,q.answered,q.anonymous,q.creatorName,q.creatorHash ", id).Scan(&question.Id, &question.Text, &question.Answered, &question.Anonymous, &question.CreatorName, &question.CreatorHash, &question.Votes)

	if err != nil {
		logrus.Error(err.Error())
		return voting_models.Question{}, false
	}

	return question, true
}

func (session *Postgresql) IsRunning() bool {
	var count int
	err := session.pool.QueryRow(context.Background(), "SELECT COUNT(*) FROM Sessions WHERE id = $1", session.sessionId).Scan(&count)

	if err != nil {
		logrus.Error(err.Error())
		return false
	}

	return count > 0
}

func (session *Postgresql) GetSecret() string {
	var secret string
	err := session.pool.QueryRow(context.Background(), "SELECT sessionSecret FROM Sessions WHERE id = $1", session.sessionId).Scan(&secret)

	if err != nil {
		logrus.Error(err.Error())
		return ""
	}

	return secret
}

func (session *Postgresql) GetQuestions() map[string]voting_models.Question {
	rows, err := session.pool.Query(context.Background(),
		"SELECT q.id,q.text,q.answered,q.anonymous,q.creatorName,q.creatorHash,COUNT(uv.questionid) as votes "+
			"FROM questions q LEFT JOIN uservotes uv on q.id = uv.questionid "+
			"WHERE sessionId = $1 "+
			"GROUP BY q.id,q.text,q.answered,q.anonymous,q.creatorName,q.creatorHash ", session.sessionId)

	if err != nil {
		logrus.Error(err.Error())
		return make(map[string]voting_models.Question)
	}

	defer rows.Close()

	questions := make(map[string]voting_models.Question)
	for rows.Next() {
		var question voting_models.Question
		if err = rows.Scan(&question.Id, &question.Text, &question.Answered, &question.Anonymous, &question.CreatorName, &question.CreatorHash, &question.Votes); err != nil {
			logrus.Fatalf(err.Error())
		}
		questions[question.Id] = question
	}

	return questions
}

func (session *Postgresql) GetUserVotes() map[string]map[string]bool {
	rows, err := session.pool.Query(context.Background(),
		"SELECT questionId,userHash FROM UserVotes")

	if err != nil {
		logrus.Error(err.Error())
		return make(map[string]map[string]bool)
	}

	defer rows.Close()

	userVotes := make(map[string]map[string]bool)
	for rows.Next() {
		var userHash string
		var questionId string
		if err = rows.Scan(&questionId, &userHash); err != nil {
			logrus.Fatalf(err.Error())
		}

		if userVotes[userHash] == nil {
			userVotes[userHash] = make(map[string]bool)
		}

		userVotes[userHash][questionId] = true
	}

	return userVotes
}

func (session *Postgresql) AddQuestion(text string, anonymous bool, creatorName, creatorHash string) voting_models.Question {
	questionId := ulid.Make().String()

	question := createQuestion(questionId, text, 0, false, anonymous, creatorName, creatorHash)

	_, err := session.pool.Exec(context.Background(),
		"INSERT INTO Questions"+
			"(id, sessionId,text,answered,anonymous,creatorName,creatorHash)"+
			"VALUES ($1,$2,$3,$4,$5,$6,$7)", question.Id, session.sessionId, question.Text, question.Anonymous, question.Anonymous, question.CreatorName, question.CreatorHash)

	if err != nil {
		logrus.Error(err.Error())
		return voting_models.Question{}
	}

	session.Vote(creatorHash, questionId)
	question.Voted = true
	return question
}

func (session *Postgresql) UpdateQuestion(id, text, creatorName string, anonymous bool) voting_models.Question {
	if anonymous {
		creatorName = ""
	}

	_, err := session.pool.Exec(context.Background(),
		"UPDATE Questions SET text = $1, creatorName = $2, anonymous = $3 WHERE id = $4", text, creatorName, anonymous, id)

	if err != nil {
		logrus.Error(err.Error())
	}

	question, _ := session.GetQuestion(id)

	return question

}

func (session *Postgresql) AnswerQuestion(id string) {
	_, err := session.pool.Exec(context.Background(),
		"UPDATE Questions SET answered = true WHERE id = $1", id)

	if err != nil {
		logrus.Error(err.Error())
	}
}

func (session *Postgresql) DeleteQuestion(id string) {
	_, err := session.pool.Exec(context.Background(), "DELETE FROM Questions WHERE id = $1", id)

	if err != nil {
		logrus.Error(err.Error())
		return
	}
}

func (session *Postgresql) Vote(userHash, id string) {
	_, err := session.pool.Exec(context.Background(),
		"INSERT INTO UserVotes"+
			"(questionId, userHash)"+
			"VALUES ($1,$2)", id, userHash)

	if err != nil {
		logrus.Error(err.Error())
		return
	}
}

func (session *Postgresql) UndoVote(userHash, id string) {
	_, err := session.pool.Exec(context.Background(),
		"DELETE FROM UserVotes WHERE questionId = $1 AND userHash = $2", id, userHash)

	if err != nil {
		logrus.Error(err.Error())
		return
	}
}

func createQuestion(questionId string, text string, votes int, answered, anonymous bool, creatorName, creatorHash string) voting_models.Question {
	if anonymous {
		creatorName = ""
	}
	return voting_models.NewQuestion(
		questionId,
		text,
		votes,
		answered,
		false,
		anonymous,
		creatorName,
		creatorHash)
}
