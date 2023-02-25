import { writable } from "svelte/store"
import { loginRequest, msalConfig } from "./auth.config"
import { PublicClientApplication, InteractionRequiredAuthError } from "@azure/msal-browser"

const msalInstance = new PublicClientApplication(msalConfig)

export const accessToken = writable(null)
export const idToken = writable(null)
export const user = writable(null)

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
			accessToken.set(response.accessToken)
			idToken.set(response.idToken)
			user.set(response.account)
		}
	} catch (error) {
		if (error instanceof InteractionRequiredAuthError) {
			msalInstance.acquireTokenRedirect({ scopes: ["User.Read"] })
		}
	}
}
