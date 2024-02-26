import { Centrifuge } from 'centrifuge';
import { refreshToken } from './auth/auth';
export const initCentrifuge = () => {
    // Use WebSocket transport endpoint.
    const centrifuge = new Centrifuge('ws://localhost:3333/api/v1/connection/websocket', {
        getToken: () => refreshToken()
    });

    // const centrifuge = new Centrifuge('ws://localhost:3333/api/v1/connection/websocket');
    // Allocate Subscription to a channel.
    const sub = centrifuge.newSubscription('voting');

    // React on `news` channel real-time publications.
    sub.on('publication', function(ctx) {
        console.log("Received centrifuge msg: ", ctx.data);
    });

    // Trigger subscribe process.
    sub.subscribe();
    console.log("Centrifuge subscribed")

    // Trigger actual connection establishement.
    centrifuge.connect();
}

