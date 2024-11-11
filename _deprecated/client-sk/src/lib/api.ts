import { get } from 'svelte/store';
import { idToken, refreshToken } from './auth/auth';
import { PUBLIC_API_BASE_URL } from '$env/static/public';

export interface RequestData {
	path: string;
	body?: string | null;
}

const baseUrl = `${PUBLIC_API_BASE_URL}/api/v1`;

export const getRequest = async ({ path }: RequestData) => {
	await refreshToken();
	return await fetch(`${baseUrl}${path}`, {
		method: 'GET',
		headers: {
			Authorization: `Bearer ${get(idToken)}`
		}
	});
};

export const postRequest = async ({ path, body }: RequestData) => {
	await refreshToken();
	return await fetch(`${baseUrl}${path}`, {
		method: 'POST',
		headers: {
			Authorization: `Bearer ${get(idToken)}`
		},
		body
	});
};

export const putRequest = async ({ path, body }: RequestData) => {
	await refreshToken();
	return await fetch(`${baseUrl}${path}`, {
		method: 'PUT',
		headers: {
			Authorization: `Bearer ${get(idToken)}`
		},
		body
	});
};

export const deleteRequest = async ({ path }: RequestData) => {
	await refreshToken();
	return await fetch(`${baseUrl}${path}`, {
		method: 'DELETE',
		headers: {
			Authorization: `Bearer ${get(idToken)}`
		}
	});
};
