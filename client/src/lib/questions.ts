import { get, writable, type Writable } from "svelte/store"
import type { Question } from "../models/question"
import { deleteRequest, getRequest, postRequest, putRequest } from "./api"
import { eventSource } from "./eventsource"
import { activeSessison } from "./session"

export const questions: Writable<Question[]> = writable([])
export const sessionActive = writable(false)

const unsub = eventSource.subscribe((eventSource) => {
	if (eventSource) {
		eventSource.addEventListener("new_question", (event) => {
			const data = JSON.parse(event.data)
			questionAdded(data)
		})
		eventSource.addEventListener("upvote_question", (event) => {
			const data = JSON.parse(event.data)
			questionVoted(data)
		})
		eventSource.addEventListener("update_question", (event) => {
			const data = JSON.parse(event.data)
			questionEdited(data)
		})
		eventSource.addEventListener("delete_question", (event) => {
			const data = JSON.parse(event.data)
			questionDeleted(data)
		})
		eventSource.addEventListener("answer_question", (event) => {
			const data = JSON.parse(event.data)
			questionAnswered(data)
		})
	}
})

const questionMap = new Map<string, Question>()

export const getQuestions = async () => {
	try {
		const repsonse = await getRequest({ path: "/question/session" })
		if (repsonse.ok) {
			activeSessison.set(true)
			const data = await repsonse.json()
			console.log(data)
			data.forEach((question) => {
				questionMap.set(question.Id, question)
			})
			sortAndUpdateQuestions()
		}
	} catch (error) {
		console.log(error)
	}
}

export function getQuestion(id: string) {
	return questionMap.get(id)
}

export function clearQuestions() {
	questionMap.clear()
	sortAndUpdateQuestions()
}

export const postQuestion = async (questionText) => {
	await postRequest({ path: "/question/new", body: JSON.stringify({ anonymous: true, text: questionText }) })
}

export const updateQuestion = async (payload: { Id: string; Anonymous: boolean; Text: string }) => {
	await putRequest({ path: "/question/update", body: JSON.stringify({ id: payload.Id, text: payload.Text, anonymous: payload.Anonymous }) })
}

export const voteQuestion = async (questionId) => {
	await putRequest({ path: `/question/upvote/${questionId}` })
}

export const answerQuestion = async (questionId) => {
	await putRequest({ path: `/question/answer/${questionId}` })
}

export const deleteQuestion = async (questionId) => {
	await deleteRequest({ path: `/question/delete/${questionId}` })
}

const questionAdded = (question: Question) => {
	questionMap.set(question.Id, question)
	sortAndUpdateQuestions()
}

const questionDeleted = (payload: { Id: string }) => {
	questionMap.delete(payload.Id)
	sortAndUpdateQuestions()
}

const questionVoted = (payload: { Id: string; Votes: number }) => {
	const votedQuestion = questionMap.get(payload.Id)
	questionMap.set(payload.Id, Object.assign({}, votedQuestion, { Votes: payload.Votes }))
	sortAndUpdateQuestions()
}

const questionEdited = (payload: { Id: string; Text: string; Creator: string; Anonymous: boolean }) => {
	const editedQuestion = questionMap.get(payload.Id)
	questionMap.set(payload.Id, Object.assign({}, editedQuestion, { ...payload }))
	sortAndUpdateQuestions()
}

const questionAnswered = (payload: { Id: string }) => {
	const editedQuestion = questionMap.get(payload.Id)
	questionMap.set(payload.Id, Object.assign({}, editedQuestion, { Answered: true }))
	sortAndUpdateQuestions()
}

const sortAndUpdateQuestions = () => {
	questions.set([...questionMap.values()].sort((a, b) => b.Votes - a.Votes))
}
