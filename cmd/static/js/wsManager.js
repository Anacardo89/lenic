import * as notifs from './notifs.js';
import { DMModule, DMChatModule } from './auth.js';
import * as dms from './dms.js';

export let ws = null;

export const MSG_COMMENT_RATE = ' has rated your comment.';
export const MSG_POST_RATE = ' has rated your post.';
export const MSG_COMMENT_TAG = ' has tagged you in their comment';
export const MSG_POST_TAG = ' has tagged you in their post';
export const MSG_COMMENT_ON_POST = ' has commented on your post';
export const MSG_FOLLOW_ACCEPT = ' has accepted your follow request.';
export const MSG_FOLLOW_REQUEST = ' has requested to follow you.';

export const TYPE_COMMENT_RATE = 'rate_comment';
export const TYPE_POST_RATE = 'rate_post';
export const TYPE_COMMENT_TAG = 'comment_tag';
export const TYPE_POST_TAG = 'post_tag';
export const TYPE_COMMENT_ON_POST = 'comment_on_post';
export const TYPE_FOLLOW_ACCEPT = 'follow_accept';
export const TYPE_FOLLOW_REQUEST = 'follow_request';
export const TYPE_DM = 'dm';


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
        const dmButton = $('.dm-button');
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
            case TYPE_COMMENT_TAG:
                handleCommentTag(message);
                if (!message.is_read) {
                    notifButton.css('--notif-display', 'block');
                }
                break;
            case TYPE_POST_TAG:
                handlePostTag(message);
                if (!message.is_read) {
                    notifButton.css('--notif-display', 'block');
                }
                break;
            case TYPE_COMMENT_ON_POST:
                handleCommentOnPost(message);
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
            case TYPE_DM:
                const $dmWindow = $('#dm-window');
                const $dmTitle = $('#dm-title');
                const $dmContent = $('#dm-content');
                if (!$dmWindow.hasClass('hidden')) {
                    if ($dmTitle.text() === message.fromuser.username) {
                        DMChatModule.appendMessage(message.msg, 'received');
                        $dmContent.scrollTop($dmContent[0].scrollHeight);
                        DMChatModule.readConversation(message.resource_id);
                    }    
                } else {
                    DMModule.clearAndFetchConversations()
                    dmButton.css('--dm-display', 'block');
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

function handleCommentTag(notification) {
    const notifContainer = $('.notif-body');
    const notif = notifs.makeCommentTagNotif(notification);
    notifContainer.prepend(notif);
}

function handlePostTag(notification) {
    const notifContainer = $('.notif-body');
    const notif = notifs.makePostTagNotif(notification);
    notifContainer.prepend(notif);
}

function handleCommentOnPost(notification) {
    const notifContainer = $('.notif-body');
    const notif = notifs.makeCommentOnPostNotif(notification);
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