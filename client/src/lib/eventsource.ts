import { get, writable } from "svelte/store"
import { idToken, refreshToken } from "./auth/auth"

export const eventSource = writable<EventSource | null>(null)

export const initEventSource = () => {
	const source = new EventSource(`${import.meta.env.VITE_API_BASE_URL}/api/v1/events`, {
		headers: {
			Authorization: `Bearer ${get(idToken)}`
		}
	} as any)
	source.addEventListener("heart_beat", () => console.log("[Heart Beat]"))
	source.addEventListener("error", (event: any) => {
		console.log("[ERROR]", event)
		source.close()
		if (event.status === 401) {
			refreshToken()
		} else {
			initEventSource()
		}
	})

	eventSource.set(source)
}
