import { Centrifuge  } from 'centrifuge';
import { refreshToken } from './auth/auth';
import { writable } from 'svelte/store';

export const centrifuge = writable<Centrifuge | null>(null);

export const initCentrifuge = () => {
    // Use WebSocket transport endpoint.
    const c = new Centrifuge('ws://localhost:3333/api/v1/connection/websocket', {
        getToken: () => refreshToken()
    });
    centrifuge.set(c)

    // const centrifuge = new Centrifuge('ws://localhost:3333/api/v1/connection/websocket');
    // Allocate Subscription to a channel.
    // const sub = centrifuge.newSubscription('voting');

    // React on `news` channel real-time publications.
    // sub.on('publication', function(ctx) {
    //     console.log("Received centrifuge msg: ", ctx.data);
    // });

    // Trigger subscribe process.
    // sub.subscribe();
    // console.log("Centrifuge subscribed")

    // Trigger actual connection establishement.
    c.connect();
}

