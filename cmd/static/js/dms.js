
export function makeConversation(conversation) {
    const fromUser = conversation.fromuser.username;
    const notif = document.createElement('div');
    notif.classList.add('dm-item');
    
    const authorInline = document.createElement('div');
    authorInline.classList.add('author-info-inline');
    const profilePic = document.createElement('img');
    profilePic.classList.add('profile-pic-mini');
    if (conversation.fromuser.profile_pic === '') {
        profilePic.src = '/static/img/no-profile-pic.jpg';
    } else {
        profilePic.src = '/action/profile-pic?user-encoded=' + conversation.fromuser.encoded
    }

    const notifMsg = document.createElement('div');
    notifMsg.innerHTML = '<strong>' + conversation.fromuser.username + '</strong> sent you a message';

    const openDMButton = document.createElement('button');
    openDMButton.innerText = 'Open';
    openDMButton.classList.add('open-dm-button');
    openDMButton.setAttribute('data-conversation-id', conversation.id);
    openDMButton.setAttribute('data-from', fromUser);

    authorInline.append(profilePic);
    authorInline.append(notifMsg);
    authorInline.append(idHidden); 
    authorInline.append(openDMButton);
    notif.append(authorInline);

    return notif;
}