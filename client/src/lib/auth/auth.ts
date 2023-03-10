import { writable } from "svelte/store"
import { loginRequest, msalConfig } from "./auth.config"
import { PublicClientApplication, InteractionRequiredAuthError, type AuthenticationResult } from "@azure/msal-browser"
import { initEventSource } from "../eventsource"

const msalInstance = new PublicClientApplication(msalConfig)
const refreshInterval = 60 * 1000 * 15

export const idToken = writable(null)
export const user = writable(null)

export const refreshToken = async () => {
	const refreshResult: AuthenticationResult = await msalInstance.acquireTokenSilent({ scopes: ["User.Read"] })
	idToken.set(refreshResult.idToken)
	initEventSource()
}

export const authenticate = async () => {
	try {
		const res = await msalInstance.handleRedirectPromise()
		msalInstance.getAllAccounts()[0] ?? (await msalInstance.loginRedirect(loginRequest))
		const accounts = await msalInstance.getAllAccounts()
		if (accounts.length > 0) {
			msalInstance.setActiveAccount(accounts[0])
			const response = await msalInstance.acquireTokenSilent({
				scopes: ["User.Read"]
			})
			idToken.set(response.idToken)
			user.set(response.account)
			initEventSource()
			setInterval(refreshToken, refreshInterval)
		}
	} catch (error) {
		if (error instanceof InteractionRequiredAuthError) {
			msalInstance.acquireTokenRedirect({ scopes: ["User.Read"] })
		}
	}
}
