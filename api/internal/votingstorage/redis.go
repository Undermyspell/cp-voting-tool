package votingstorage

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sse/internal/helper"
	"sse/internal/models"

	goredis "github.com/go-redis/redis/v8"
	"github.com/nitishm/go-rejson/v4"
	"github.com/oklog/ulid/v2"
	"github.com/sirupsen/logrus"
)

type redisQuestion struct {
	Id          string
	Text        string
	Votes       int
	Answered    bool
	Voted       bool
	CreatorHash string
	CreatorName string
	Anonymous   bool
}

func newRedisQuestion(text string, anonymous bool, creatorName, creatorHash string) redisQuestion {
	if anonymous {
		creatorName = ""
	}
	return redisQuestion{
		Id:          ulid.Make().String(),
		Text:        text,
		Votes:       0,
		Answered:    false,
		Voted:       false,
		CreatorHash: creatorHash,
		CreatorName: creatorName,
		Anonymous:   anonymous,
	}
}

type VotingSession struct {
	Questions     map[string]redisQuestion
	UserVotes     map[string]map[string]bool
	SessionSecret string
}

type Redis struct {
	redisHandler *rejson.Handler
	goRedisCli   *goredis.Client
}

func (session *Redis) Start() {
	votingSession := VotingSession{
		UserVotes:     make(map[string]map[string]bool),
		Questions:     make(map[string]redisQuestion),
		SessionSecret: helper.GetRandomId(),
	}

	_, err := session.redisHandler.JSONSet("voting_session", ".", votingSession)
	if err != nil {
		log.Fatalf(err.Error())
		return
	}
}

func (session *Redis) Stop() {
	if err := session.goRedisCli.FlushAll(context.Background()).Err(); err != nil {
		logrus.Fatalf("goredis - failed to flush: %v", err)
	}
}

func (session *Redis) GetQuestion(id string) (models.Question, bool) {
	res, err := session.redisHandler.JSONGet("voting_session", fmt.Sprintf(".Questions.%s", id))
	if err != nil {
		return models.Question{}, false
	}

	redisQuestion := redisQuestion{}
	json.Unmarshal(res.([]byte), &redisQuestion)

	question := questionFromRedisQuestion(redisQuestion)

	return question, true
}

func (session *Redis) IsRunning() bool {
	_, err := session.redisHandler.JSONGet("voting_session", ".Questions")
	return err == nil
}

func (session *Redis) GetSecret() string {
	res, err := session.redisHandler.JSONGet("voting_session", ".SessionSecret")
	if err != nil {
		log.Fatalf(err.Error())
		return ""
	}

	return string(res.([]byte))
}

func (session *Redis) GetQuestions() map[string]models.Question {
	res, err := session.redisHandler.JSONGet("voting_session", ".Questions")
	if err != nil {
		log.Fatalf(err.Error())
		return make(map[string]models.Question)
	}

	redisQuestions := make(map[string]redisQuestion)
	json.Unmarshal(res.([]byte), &redisQuestions)

	questions := make(map[string]models.Question)

	for id, question := range redisQuestions {
		questions[id] = questionFromRedisQuestion(question)
	}

	return questions
}

func (session *Redis) GetUserVotes() map[string]map[string]bool {
	res, err := session.redisHandler.JSONGet("voting_session", ".UserVotes")
	if err != nil {
		log.Fatalf(err.Error())
		return map[string]map[string]bool{}
	}

	userVotes := make(map[string]map[string]bool)
	json.Unmarshal(res.([]byte), &userVotes)

	return userVotes
}

func (session *Redis) AddQuestion(text string, anonymous bool, creatorName, creatorHash string) models.Question {
	question := newRedisQuestion(text, anonymous, creatorName, creatorHash)

	_, err := session.redisHandler.JSONSet("voting_session", fmt.Sprintf(".Questions.%s", question.Id), question)
	if err != nil {
		log.Fatalf(err.Error())
		return models.Question{}
	}

	return models.Question{}
}

func (session *Redis) UpdateQuestion(id, text, creatorName string, anonymous bool) models.Question {
	if anonymous {
		creatorName = ""
	}

	session.redisHandler.JSONSet("voting_session", fmt.Sprintf(".Questions.%s.Text", id), text)
	session.redisHandler.JSONSet("voting_session", fmt.Sprintf(".Questions.%s.Anonymous", id), anonymous)
	session.redisHandler.JSONSet("voting_session", fmt.Sprintf(".Questions.%s.CreatorName", id), creatorName)

	res, err := session.redisHandler.JSONGet("voting_session", fmt.Sprintf(".Questions.%s", id))
	if err != nil {
		log.Fatalf(err.Error())
		return models.Question{}
	}

	redisQuestion := redisQuestion{}
	json.Unmarshal(res.([]byte), &redisQuestion)

	return questionFromRedisQuestion(redisQuestion)
}

func (session *Redis) AnswerQuestion(id string) {
	session.redisHandler.JSONSet("voting_session", fmt.Sprintf(".Questions.%s.Answered", id), true)
}

func (session *Redis) DeleteQuestion(id string) {
	session.redisHandler.JSONDel("voting_session", fmt.Sprintf(".Questions.%s", id))
}

func (session *Redis) Vote(userHash, id string) {
	if _, err := session.redisHandler.JSONGet("voting_session", fmt.Sprintf(".UserVotes.%s", userHash)); err != nil {
		votedQuestions := make(map[string]bool)
		session.redisHandler.JSONSet("voting_session", fmt.Sprintf(".UserVotes.%s", userHash), votedQuestions)
	}

	if _, err := session.redisHandler.JSONSet("voting_session", fmt.Sprintf(".UserVotes.%s.%s", userHash, id), true); err != nil {
		logrus.Fatal(err.Error())
	}

	session.redisHandler.JSONNumIncrBy("voting_session", fmt.Sprintf(".Questions.%s.Votes", id), 1)
}

func questionFromRedisQuestion(question redisQuestion) models.Question {
	return models.NewQuestion(
		question.Id,
		question.Text,
		question.Votes,
		question.Answered,
		question.Voted,
		question.Anonymous,
		question.CreatorName,
		question.CreatorHash)
}
