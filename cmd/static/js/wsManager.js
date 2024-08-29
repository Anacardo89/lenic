
export let ws = null;

export const MSG_COMMENT_RATE = ' has rated your comment.'
export const MSG_POST_RATE = ' has rated your post.'

// WebSocket connection
export function connectWS(user_name) {
    if (ws && ws.readyState === WebSocket.OPEN) {
        console.log('WebSocket connection already open');
        return;
    }

    const wsUrl = `wss://${window.location.host}/ws?user_name=${user_name}`;
    ws = new WebSocket(wsUrl);

    ws.onopen = function() {
        console.log('WebSocket connection established');
    };

    ws.onmessage = function(event) {
        const message = JSON.parse(event.data);

        switch (message.type) {
            case 'rate-comment':
                handleRateComment(message.data);
                break;
            case 'rate-post':
                handleRatePost(message.data);
                break;
            case 'command':
                executeCommand(message.data);
                break;
            default:
                console.warn('Unknown message type:', message.type);
        }
        console.log('Message from server:', event.data);
    };

    ws.onerror = function(error) {
        console.error('WebSocket error:', error);
    };

    ws.onclose = function(event) {
        if (event.wasClean) {
            console.log('WebSocket connection closed cleanly:', event);
        } else {
            console.error('WebSocket connection closed with error:', event);
        }
    };
}

export function sendWSmsg(message) {
    if (ws && ws.readyState === WebSocket.OPEN) {
        ws.send(JSON.stringify(message));
        console.log('Message sent to server:', message);
    } else {
        console.error('WebSocket is not open. Cannot send message.');
    }
}

export function closeWS() {
    if (ws && ws.readyState === WebSocket.OPEN) {
        ws.close();
    }
    ws = null;
}