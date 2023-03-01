import { get, writable, type Writable } from "svelte/store"
import type { Question } from "../models/question"
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
	}
	try {
		const repsonse = await fetch("http://localhost:3333/api/v1/question/session", {
			headers: {
				Authorization: `Bearer ${get(idToken)}`
			}
		})
		if (repsonse.ok) {
			activeSessison.set(true)
			const data = await repsonse.json()
			data.forEach((question) => {
				questionMap.set(question.Id, question)
			})
			questions.set([...questionMap.values()].sort((a, b) => b.Votes - a.Votes))
		}
	} catch (error) {
		console.log(error)
	}
}

const questionAdded = (question: Question) => {
	questions.update((questions) => {
		questionMap.set(question.Id, question)
		return [...questionMap.values()].sort((a, b) => a.Votes - b.Votes)
	})
}

const questionVoted = (payload: { Id: string; Votes: number }) => {
	const votedQuestion = questionMap.get(payload.Id)
	questionMap.set(payload.Id, Object.assign({}, votedQuestion, { Votes: payload.Votes }))
	questions.set([...questionMap.values()].sort((a, b) => a.Votes - b.Votes))
}
