import { get, writable, type Writable } from "svelte/store"
import type { Question } from "../models/question"
import { idToken } from "./auth/auth"
import { activeSessison } from "./session"

export const questions: Writable<Question[]> = writable([])
export const sessionActive = writable(false)
let source: EventSource | null = null

export const getQuestions = async (eventSource: EventSource) => {
	source = eventSource
	source.addEventListener("new_question", (event) => {
		const data = JSON.parse(event.data)
		questionAdded(data)
	})
	try {
		const repsonse = await fetch("http://localhost:3333/api/v1/question/session", {
			headers: {
				Authorization: `Bearer ${get(idToken)}`
			}
		})
		if (repsonse.ok) {
			activeSessison.set(true)
			const data = await repsonse.json()
			questions.set(data)
		}
	} catch (error) {
		console.log(error)
	}
}

const questionAdded = (question: Question) => {
	questions.update((questions) => {
		return (questions = [...questions, question])
	})
}
