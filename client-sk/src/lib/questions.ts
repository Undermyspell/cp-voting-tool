import { derived, writable, type Writable } from "svelte/store"
import type { Question } from "./models/question"
import { deleteRequest, getRequest, postRequest, putRequest } from "./api"
import { eventSource } from "./eventsource"
import { activeSessison } from "./session"

const allQuestions: Writable<Question[]> = writable([])
export const questions = derived(allQuestions, ($allQuestions: Question[]) => $allQuestions.filter((q) => q.Answered === false))
export const answeredQuestions = derived(allQuestions, ($allQuestions: Question[]) => $allQuestions.filter((q) => q.Answered === true))
export const sessionActive = writable(false)

eventSource.subscribe((eventSource) => {
	if (eventSource) {
		eventSource.addEventListener("new_question", (event) => {
			const data = JSON.parse(event.data)
			questionAdded(data)
		})
		eventSource.addEventListener("upvote_question", (event) => {
			const data = JSON.parse(event.data)
			questionVoted(data)
		})
		eventSource.addEventListener("undo_upvote_question", (event) => {
			const data = JSON.parse(event.data)
			questionVoteUndone(data)
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
			const data = await repsonse.json() as Question[]
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

export const postQuestion = async (questionText: string, anonymous: boolean) => {
	await postRequest({ path: "/question/new", body: JSON.stringify({ anonymous: anonymous, text: questionText }) })
}

export const updateQuestion = async (payload: { Id: string; Anonymous: boolean; Text: string }) => {
	await putRequest({ path: "/question/update", body: JSON.stringify({ id: payload.Id, text: payload.Text, anonymous: payload.Anonymous }) })
}

export const voteQuestion = async (questionId: string) => {
	await putRequest({ path: `/question/upvote/${questionId}` })
}

export const undoVoteQuestion = async (questionId: string) => {
	await putRequest({ path: `/question/undovote/${questionId}` })
}

export const answerQuestion = async (questionId: string) => {
	await putRequest({ path: `/question/answer/${questionId}` })
}

export const deleteQuestion = async (questionId: string) => {
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

const questionVoted = (payload: { Id: string; Votes: number; Voted: boolean }) => {
	const votedQuestion = questionMap.get(payload.Id)
	questionMap.set(payload.Id, Object.assign({}, votedQuestion, { ...payload }))
	sortAndUpdateQuestions()
}

const questionVoteUndone = (payload: { Id: string; Votes: number; Voted: boolean }) => {
	const votedQuestion = questionMap.get(payload.Id)
	questionMap.set(payload.Id, Object.assign({}, votedQuestion, { ...payload }))
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
	allQuestions.set([...questionMap.values()].sort((a, b) => b.Votes - a.Votes))
}
