import { writable } from "svelte/store"
import { idToken } from "./auth/auth"

export const questions = writable([])
export const sessionActive = writable(false)

export const getQuestions = async () => {
	let token = ""
	const unsub = idToken.subscribe((value) => {
		token = value
	})
	const repsonse = await fetch("http://localhost:3333/api/v1/question/session", {
		headers: {
			Authorization: `Bearer ${token}`
		}
	})
	const data = await repsonse.json()
	questions.set(data)
	unsub()
}
