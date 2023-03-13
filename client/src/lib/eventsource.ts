import { get, writable } from "svelte/store"
import { idToken, refreshToken } from "./auth/auth"

export const eventSource = writable<EventSource | null>(null)

export const initEventSource = () => {
	const source = new EventSource(`${import.meta.env.VITE_API_BASE_URL}/api/v1/events`, {
		headers: {
			Authorization: `Bearer ${get(idToken)}`
		}
	} as any)
	source.addEventListener("error", (event: any) => {
		if (event.status === 401) {
			refreshToken()
		}
	})

	eventSource.set(source)
}
