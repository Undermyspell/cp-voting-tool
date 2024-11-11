import { Centrifuge  } from 'centrifuge';
import { refreshToken } from './auth/auth';
import { writable } from 'svelte/store';
import { PUBLIC_API_BASE_URL_WS } from '$env/static/public';
export const centrifuge = writable<Centrifuge | null>(null);

export const initCentrifuge = () => {
    // Use WebSocket transport endpoint.
    const c = new Centrifuge(`${PUBLIC_API_BASE_URL_WS}/api/v1/connection/websocket`, {
        getToken: () => refreshToken()
    });
    centrifuge.set(c)
    c.connect();
}

