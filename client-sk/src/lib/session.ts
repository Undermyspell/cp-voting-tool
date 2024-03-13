import { writable } from 'svelte/store';
import { getRequest, postRequest } from './api';
import { clearQuestions } from './questions';

export const activeSessison = writable(false);
export const userOnline = writable(0);

export const getSession = async () => {
	try {
		const response = await getRequest({ path: '/question/session' });
		if (response.status === 200) {
			activeSessison.set(true);
		}
	} catch (error) {
		console.log(error);
	}
};

export const startSession = async () => {
	try {   
		await postRequest({ path: '/question/session/start' });
	} catch (error) {
		console.log(error);
	}
};

export const stopSession = async () => {
	try {
		await postRequest({ path: '/question/session/stop' });
		clearQuestions();
	} catch (error) {
		console.log(error);
	}
};
