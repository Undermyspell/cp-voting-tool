import { get, writable } from "svelte/store"
import { idToken } from "./auth/auth"

export const eventSource = writable<EventSource | null>(null)

export const initEventSource = () => {
	eventSource.set(
		new EventSource("http://localhost:3333/api/v1/events", {
			headers: {
				Authorization: `Bearer ${get(idToken)}`
			}
		} as any)
	)
}

const subscribe = () => {}
