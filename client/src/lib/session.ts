import { writable } from "svelte/store"
import { postRequest } from "./api"
import { eventSource } from "./eventsource"

export const activeSessison = writable(false)

const unsub = eventSource.subscribe((eventSource) => {
	console.log(eventSource)
	if (eventSource) {
		eventSource.addEventListener("start_session", (event) => {
			activeSessison.set(true)
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
	} catch (error) {
		console.log(error)
	}
}
