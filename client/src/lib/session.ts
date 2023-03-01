import { get, writable } from "svelte/store"
import { postRequest } from "./api"
import { idToken } from "./auth/auth"

export const activeSessison = writable(false)

export const startSession = async () => {
	try {
		const repsonse = await postRequest({path: "/question/session/start"})

		if (repsonse.ok) {
			activeSessison.set(true)
		}
	} catch (error) {
		console.log(error)
	}
}

export const stopSession = async () => {
	try {
		const repsonse = await postRequest({path: "/question/session/stop"})
		if (repsonse.ok) {
			activeSessison.set(false)
		}
	} catch (error) {
		console.log(error)
	}
}
