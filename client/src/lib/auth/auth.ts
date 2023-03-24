import { derived, get, writable, type Writable } from "svelte/store"
import { loginRequest, msalConfig } from "./auth.config"
import { PublicClientApplication, InteractionRequiredAuthError, type AuthenticationResult } from "@azure/msal-browser"
import { initEventSource } from "../eventsource"

const msalInstance = new PublicClientApplication(msalConfig)

const authResult: Writable<AuthenticationResult> = writable(null)

export const refreshToken = async () => {
	try {
		const tokenExpired = !!get(authResult) && +get(authResult).idTokenClaims["exp"] * 1000 < Date.now()
		console.log("[Force refresh] ", tokenExpired)
		const refreshResult: AuthenticationResult = await msalInstance.acquireTokenSilent({
			scopes: ["User.Read"],
			forceRefresh: tokenExpired
		})
		const initEvSource = get(authResult)?.idTokenClaims["exp"] !== refreshResult.idTokenClaims["exp"]
		authResult.set(refreshResult)

		if (initEvSource) {
			initEventSource()
		}
	} catch (error) {
		console.error("[Acquire Token Error]: ", error)
		if (error instanceof InteractionRequiredAuthError) {
			msalInstance.acquireTokenRedirect({ scopes: ["User.Read"] })
		}
	}
}

export const idToken = derived(authResult, ($values) => $values?.idToken)
export const user = derived(authResult, ($values) => $values?.account)
const roles = derived(authResult, ($values) => ($values.idTokenClaims["roles"] as string[]) ?? [])
export const isAdmin = () => get(roles).filter((role) => role === "admin").length > 0
export const isSessionAdmin = () => get(roles).filter((role) => role === "session_admin").length > 0

export const authenticate = async () => {
	try {
		await msalInstance.handleRedirectPromise()
		msalInstance.getAllAccounts()[0] ?? (await msalInstance.loginRedirect(loginRequest))
		const accounts = msalInstance.getAllAccounts()
		if (accounts.length > 0) {
			msalInstance.setActiveAccount(accounts[0])
			await refreshToken()
		}
	} catch (error) {
		if (error instanceof InteractionRequiredAuthError) {
			msalInstance.acquireTokenRedirect({ scopes: ["User.Read"] })
		}
	}
}
