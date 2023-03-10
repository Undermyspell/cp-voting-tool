import { writable } from "svelte/store"
import { postRequest } from "./api"
import { eventSource } from "./eventsource"
import { clearQuestions } from "./questions"

export const activeSessison = writable(false)
export const userOnline = writable(0)

const unsub = eventSource.subscribe((eventSource) => {
	if (eventSource) {
		eventSource.addEventListener("start_session", (event) => {
			activeSessison.set(true)
		})
		eventSource.addEventListener("user_connected", (event) => {
			const data = JSON.parse(event.data)
			userOnline.set(data.UserCount)
		})
		eventSource.addEventListener("user_disconnected", (event) => {
			const data = JSON.parse(event.data)
			userOnline.set(data.UserCount)
		})
		eventSource.addEventListener("stop_session", (event) => {
			activeSessison.set(false)
		})
	}
})

export const startSession = async () => {
	try {
		const repsonse = await postRequest({ path: "/question/session/start" })
	} catch (error) {
		console.log(error)
	}
}

export const stopSession = async () => {
	try {
		const repsonse = await postRequest({ path: "/question/session/stop" })
		clearQuestions()
	} catch (error) {
		console.log(error)
	}
}
