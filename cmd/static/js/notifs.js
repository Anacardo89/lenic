import { session_username, session_encoded } from "./auth.js"
import * as wsoc from './wsManager.js';


export function makeCommentRateNotif(notification) {
    const postGuid = notification.parent_id;
    const notif = document.createElement('div');
    notif.classList.add('notif-item');
    if (!notification.is_read) {
        notif.classList.add('notif-item-unread');
    }
    const authorInline = document.createElement('div');
    authorInline.classList.add('author-info-inline');
    const profilePic = document.createElement('img');
    profilePic.classList.add('profile-pic-mini');
    if (notification.fromuser.profile_pic === '') {
        profilePic.src = '/static/img/no-profile-pic.jpg';
    } else {
        profilePic.src = '/action/profile-pic?user-encoded=' + notification.fromuser.encoded
    }
    const notifMsg = document.createElement('div');
    notifMsg.innerHTML = '<strong>' + notification.fromuser.username + '</strong> ' + notification.msg;
    const idHidden = document.createElement('input');
    idHidden.type = 'hidden';
    idHidden.value = notification.id;
    const readHidden = document.createElement('input');
    readHidden.type = 'hidden';
    readHidden.value = notification.is_read;
    authorInline.append(profilePic);
    authorInline.append(notifMsg);
    authorInline.append(idHidden);
    authorInline.append(readHidden);
    notif.append(authorInline);

    notif.addEventListener('click', function() {
        $.ajax({
            url: '/action/user/' + session_encoded + '/notifications/' + idHidden.value + '/read',
            method: 'PUT',
            success: function() {
                window.location.href = '/post/' +  postGuid;
            },
            error: function(xhr) {
                const errorMessage = xhr.responseText;
                window.location.href = '/error?message=' + encodeURIComponent(errorMessage);
            }
        });
    });
    return notif;
}

export function makePostRateNotif(notification) {
    const postGuid = notification.resource_id;
    const notif = document.createElement('div');
    notif.classList.add('notif-item');
    if (!notification.is_read) {
        notif.classList.add('notif-item-unread');
    }
    const authorInline = document.createElement('div');
    authorInline.classList.add('author-info-inline');
    const profilePic = document.createElement('img');
    profilePic.classList.add('profile-pic-mini');
    if (notification.fromuser.profile_pic === '') {
        profilePic.src = '/static/img/no-profile-pic.jpg';
    } else {
        profilePic.src = '/action/profile-pic?user-encoded=' + notification.fromuser.encoded
    }
    const notifMsg = document.createElement('div');
    notifMsg.innerHTML = '<strong>' + notification.fromuser.username + '</strong> ' + notification.msg;
    const idHidden = document.createElement('input');
    idHidden.type = 'hidden';
    idHidden.value = notification.id;
    const readHidden = document.createElement('input');
    readHidden.type = 'hidden';
    readHidden.value = notification.is_read;
    authorInline.append(profilePic);
    authorInline.append(notifMsg);
    authorInline.append(idHidden);
    authorInline.append(readHidden); 
    notif.append(authorInline);

    notif.addEventListener('click', function() {
        $.ajax({
            url: '/action/user/' + session_encoded + '/notifications/' + idHidden.value + '/read',
            method: 'PUT',
            success: function() {
                window.location.href = '/post/' +  postGuid;
            },
            error: function(xhr) {
                const errorMessage = xhr.responseText;
                window.location.href = '/error?message=' + encodeURIComponent(errorMessage);
            }
        });
    });
    return notif;
}

export function makeCommentOnPostNotif(notification) {
    const postGuid = notification.parent_id;
    const notif = document.createElement('div');
    notif.classList.add('notif-item');
    if (!notification.is_read) {
        notif.classList.add('notif-item-unread');
    }
    const authorInline = document.createElement('div');
    authorInline.classList.add('author-info-inline');
    const profilePic = document.createElement('img');
    profilePic.classList.add('profile-pic-mini');
    if (notification.fromuser.profile_pic === '') {
        profilePic.src = '/static/img/no-profile-pic.jpg';
    } else {
        profilePic.src = '/action/profile-pic?user-encoded=' + notification.fromuser.encoded
    }
    const notifMsg = document.createElement('div');
    notifMsg.innerHTML = '<strong>' + notification.fromuser.username + '</strong> ' + notification.msg;
    const idHidden = document.createElement('input');
    idHidden.type = 'hidden';
    idHidden.value = notification.id;
    const readHidden = document.createElement('input');
    readHidden.type = 'hidden';
    readHidden.value = notification.is_read;
    authorInline.append(profilePic);
    authorInline.append(notifMsg);
    authorInline.append(idHidden);
    authorInline.append(readHidden);
    notif.append(authorInline);

    notif.addEventListener('click', function() {
        $.ajax({
            url: '/action/user/' + session_encoded + '/notifications/' + idHidden.value + '/read',
            method: 'PUT',
            success: function() {
                window.location.href = '/post/' +  postGuid;
            },
            error: function(xhr) {
                const errorMessage = xhr.responseText;
                window.location.href = '/error?message=' + encodeURIComponent(errorMessage);
            }
        });
    });
    return notif;
}

export function makeFollowAcceptNotif(notification) {
    const notif = document.createElement('div');
    notif.classList.add('notif-item');
    if (!notification.is_read) {
        notif.classList.add('notif-item-unread');
    }
    const authorInline = document.createElement('div');
    authorInline.classList.add('author-info-inline');
    const profilePic = document.createElement('img');
    profilePic.classList.add('profile-pic-mini');
    if (notification.fromuser.profile_pic === '') {
        profilePic.src = '/static/img/no-profile-pic.jpg';
    } else {
        profilePic.src = '/action/profile-pic?user-encoded=' + notification.fromuser.encoded
    }
    const notifMsg = document.createElement('div');
    notifMsg.innerHTML = '<strong>' + notification.fromuser.username + '</strong> ' + notification.msg;
    const idHidden = document.createElement('input');
    idHidden.type = 'hidden';
    idHidden.value = notification.id;
    const readHidden = document.createElement('input');
    readHidden.type = 'hidden';
    readHidden.value = notification.is_read;
    authorInline.append(profilePic);
    authorInline.append(notifMsg);
    authorInline.append(idHidden);
    authorInline.append(readHidden); 
    notif.append(authorInline);

    notif.addEventListener('click', function() {
        $.ajax({
            url: '/action/user/' + session_encoded + '/notifications/' + idHidden.value + '/read',
            method: 'PUT',
            success: function() {
                window.location.href = '/user/' +  notification.fromuser.encoded;
            },
            error: function(xhr) {
                const errorMessage = xhr.responseText;
                window.location.href = '/error?message=' + encodeURIComponent(errorMessage);
            }
        });
    });
    return notif;
}

export function makeFollowRequestNotif(notification) {
    const fromUser = notification.fromuser.username;
    const notif = document.createElement('div');
    notif.classList.add('notif-item');
    
    const authorInline = document.createElement('div');
    authorInline.classList.add('author-info-inline');

    const profilePicLink = document.createElement('a');
    profilePicLink.href = '/user/' + notification.fromuser.encoded;
    
    const profilePic = document.createElement('img');
    profilePic.classList.add('profile-pic-mini');
    if (notification.fromuser.profile_pic === '') {
        profilePic.src = '/static/img/no-profile-pic.jpg';
    } else {
        profilePic.src = '/action/profile-pic?user-encoded=' + notification.fromuser.encoded
    }
    profilePicLink.append(profilePic);

    const notifMsg = document.createElement('div');
    notifMsg.innerHTML = '<a href="/user/' + notification.fromuser.encoded + '"><strong>' + notification.fromuser.username + '</strong> ' + notification.msg;
    
    const idHidden = document.createElement('input');
    idHidden.type = 'hidden';
    idHidden.value = notification.id;
    
    const readHidden = document.createElement('input');
    readHidden.type = 'hidden';
    readHidden.value = notification.is_read;

    const acceptRequestButton = document.createElement('button');
    acceptRequestButton.innerText = 'Accept';

    const refuseRequestButton = document.createElement('button');
    refuseRequestButton.innerText = 'Refuse';

    const buttonsDiv = document.createElement('div');
    buttonsDiv.classList.add('request-buttons');
    buttonsDiv.append(refuseRequestButton);
    buttonsDiv.append(acceptRequestButton);
    
    authorInline.append(profilePicLink);
    authorInline.append(notifMsg);
    authorInline.append(idHidden); 
    authorInline.append(readHidden); 
    notif.append(authorInline);
    notif.append(buttonsDiv);

    acceptRequestButton.addEventListener('click', function() {
        $.ajax({
            url: '/action/user/' + session_encoded + '/accept',
            method: 'PUT',
            data: {
                requester: fromUser
            },
            success: function() {
                const message = {
                    from_username: session_username,
                    type: wsoc.TYPE_FOLLOW_ACCEPT,
                    msg: wsoc.MSG_FOLLOW_ACCEPT,
                    resource_id: notification.fromuser.encoded,
                    parent_id: ''
                };
                wsoc.sendWSmsg(message);
                location.reload();
            },
            error: function(xhr) {
                const errorMessage = xhr.responseText;
                window.location.href = '/error?message=' + encodeURIComponent(errorMessage);
            }
        });
    });
    
    refuseRequestButton.addEventListener('click', function() {
        $.ajax({
            url: '/action/user/' + session_encoded + '/unfollow' + (fromUser ? '?requester=' + encodeURIComponent(fromUser) : ''),
            method: 'DELETE',
            success: function() {
                location.reload();
            },
            error: function(xhr) {
                const errorMessage = xhr.responseText;
                window.location.href = '/error?message=' + encodeURIComponent(errorMessage);
            }
        });
    });

    return notif;
}
