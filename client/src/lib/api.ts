import { get } from "svelte/store"
import { idToken, refreshToken } from "./auth/auth"

export interface RequestData {
	path: string
	body?: any
}

const baseUrl = `${import.meta.env.VITE_API_BASE_URL}/api/v1`

export const getRequest = async ({ path }: RequestData) => {
	try {
		await refreshToken()
		return await fetch(`${baseUrl}${path}`, {
			method: "GET",
			headers: {
				Authorization: `Bearer ${get(idToken)}`
			}
		})
	} catch (error) {
		console.log(error)
	}
}

export const postRequest = async ({ path, body }: RequestData) => {
	try {
		await refreshToken()
		return await fetch(`${baseUrl}${path}`, {
			method: "POST",
			headers: {
				Authorization: `Bearer ${get(idToken)}`
			},
			body
		})
	} catch (error) {
		console.log(error)
	}
}

export const putRequest = async ({ path, body }: RequestData) => {
	try {
		await refreshToken()
		return await fetch(`${baseUrl}${path}`, {
			method: "PUT",
			headers: {
				Authorization: `Bearer ${get(idToken)}`
			},
			body
		})
	} catch (error) {
		console.log(error)
	}
}

export const deleteRequest = async ({ path }: RequestData) => {
	try {
		await refreshToken()
		return await fetch(`${baseUrl}${path}`, {
			method: "DELETE",
			headers: {
				Authorization: `Bearer ${get(idToken)}`
			}
		})
	} catch (error) {
		console.log(error)
	}
}
