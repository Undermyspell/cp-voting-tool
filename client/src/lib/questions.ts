import { get, writable, type Writable } from "svelte/store"
import type { Question } from "../models/question"
import { getRequest, postRequest, putRequest } from "./api"
import { idToken } from "./auth/auth"
import { activeSessison } from "./session"

export const questions: Writable<Question[]> = writable([])
export const sessionActive = writable(false)
let eventSource: EventSource | null = null
const questionMap = new Map<string, Question>()

export const getQuestions = async () => {
	if (eventSource === null) {
		eventSource = new EventSource("http://localhost:3333/api/v1/events", {
			headers: {
				Authorization: `Bearer ${get(idToken)}`
			}
		} as any)

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
	try {
		const repsonse = await getRequest({ path: "/question/session" })
		if (repsonse.ok) {
			activeSessison.set(true)
			const data = await repsonse.json()
			data.forEach((question) => {
				questionMap.set(question.Id, question)
			})
			sortAndUpdateQuestions()
		}
	} catch (error) {
		console.log(error)
	}
}

export const postQuestion = async (questionText) => {
	await postRequest({ path: "/question/new", body: JSON.stringify({ anonymous: true, text: questionText }) })
}

export const voteQuestion = async (questionId) => {
	await putRequest({ path: `/question/upvote/${questionId}` })
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
