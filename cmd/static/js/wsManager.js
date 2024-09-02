import * as notifs from './notifs.js';

export let ws = null;

export const MSG_COMMENT_RATE = ' has rated your comment.';
export const MSG_POST_RATE = ' has rated your post.';
export const MSG_FOLLOW_ACCEPT = ' has accepted your follow request.';
export const MSG_FOLLOW_REQUEST = ' has requested to follow you.';

export const TYPE_COMMENT_RATE = 'rate_comment';
export const TYPE_POST_RATE = 'rate_post';
export const TYPE_FOLLOW_ACCEPT = 'follow_accept';
export const TYPE_FOLLOW_REQUEST = 'follow_request';

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
        const notifButton = $('.notif-button');
        console.log(message);

        switch (message.type) {
            case TYPE_COMMENT_RATE:
                handleRateComment(message);
                if (!message.is_read) {
                    notifButton.css('--notif-display', 'block');
                }
                break;
            case TYPE_POST_RATE:
                handleRatePost(message);
                if (!message.is_read) {
                    notifButton.css('--notif-display', 'block');
                }
                break;
            case TYPE_FOLLOW_ACCEPT:
                handleFollowAccept(message);
                if (!message.is_read) {
                    notifButton.css('--notif-display', 'block');
                }
                break;
            case TYPE_FOLLOW_REQUEST:
                handleFollowRequest(message);
                if (!message.is_read) {
                    notifButton.css('--notif-display', 'block');
                }
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

function handleRateComment(notification) {
    const notifContainer = $('.notif-body');
    const notif = notifs.makeCommentRateNotif(notification);
    notifContainer.prepend(notif);
}

function handleRatePost(notification) {
    const notifContainer = $('.notif-body');
    const notif = notifs.makePostRateNotif(notification);
    notifContainer.prepend(notif);
}

function handleFollowRequest(notification) {
    const notifContainer = $('.notif-body');
    const notif = notifs.makeFollowRequestNotif(notification);
    notifContainer.prepend(notif);
}

function handleFollowAccept(notification) {
    const notifContainer = $('.notif-body');
    const notif = notifs.makeFollowAcceptNotif(notification);
    notifContainer.prepend(notif);
}