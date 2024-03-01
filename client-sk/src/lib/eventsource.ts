// import { get, writable } from 'svelte/store';
// import { idToken, refreshToken } from './auth/auth';
// import { PUBLIC_API_BASE_URL } from '$env/static/public';

// export const eventSource = writable<EventSource | null>(null);

// export const initEventSource = () => {
// 	if (get(eventSource)) {
// 		get(eventSource).close();
// 	}
// 	const source = new EventSource(`${PUBLIC_API_BASE_URL}/api/v1/events`, {
// 		headers: {
// 			Authorization: `Bearer ${get(idToken)}`
// 		}
// 	} as any);
// 	source.addEventListener('heart_beat', () => console.log('[Heart Beat]'));
// 	source.addEventListener('error', (event: any) => {
// 		console.log('[ERROR]', event);
// 		if (event.status === 401) {
// 			refreshToken();
// 		} else {
// 			initEventSource();
// 		}
// 	});

// 	eventSource.set(source);
// };
