import { derived, get, writable, type Writable } from 'svelte/store';
import { loginRequest, msalConfig } from './auth.config';
import {
	PublicClientApplication,
	InteractionRequiredAuthError,
	type AuthenticationResult
} from '@azure/msal-browser';
import { initEventSource } from '../eventsource';
import { getSession } from '$lib/session';
import { initCentrifuge } from '$lib/centrifuge';

const msalInstance = new PublicClientApplication(msalConfig);
await msalInstance.initialize();

const authResult: Writable<AuthenticationResult | null> = writable(null);

export const refreshToken = async (): Promise<string> => {
	try {
		const tokenExpired =
			!!get(authResult) && (+get(authResult).idTokenClaims['exp'] - 60) * 1000 < Date.now();
		const refreshResult: AuthenticationResult = await msalInstance.acquireTokenSilent({
			scopes: ['User.Read'],
			forceRefresh: tokenExpired
		});
		const initEvSource =
			get(authResult)?.idTokenClaims['exp'] !== refreshResult.idTokenClaims['exp'];
		authResult.set(refreshResult);

		if (initEvSource) {
			// initEventSource();
		}

		return refreshResult.idToken
	} catch (error) {
		console.error('[Acquire Token Error]: ', error);
		if (error instanceof InteractionRequiredAuthError) {
			msalInstance.acquireTokenRedirect({ scopes: ['User.Read'] });
		}
		return ''
	}
};

export const idToken = derived(authResult, ($values) => $values?.idToken);
export const user = derived(authResult, ($values) => $values?.account);
const roles = derived(authResult, ($values) => {
	return ($values?.idTokenClaims['roles'] as string[]) ?? [];
});
export const isAdmin = derived(roles, ($values) => $values.includes('admin'));

export const isSessionAdmin = derived(roles, ($values) => $values.includes('session_admin'));

export const authenticate = async () => {
	try {
		await msalInstance.handleRedirectPromise();
		msalInstance.getAllAccounts()[0] ?? (await msalInstance.loginRedirect(loginRequest));
		const accounts = msalInstance.getAllAccounts();
		if (accounts.length > 0) {
			msalInstance.setActiveAccount(accounts[0]);
			await refreshToken();
			await getSession();
			initCentrifuge()
		}
	} catch (error) {
		if (error instanceof InteractionRequiredAuthError) {
			msalInstance.acquireTokenRedirect({ scopes: ['User.Read'] });
		}
	}
};
