
export function makeConversation(conversation) {
    const fromUser = conversation.user2.username;
    const convo = document.createElement('div');
    convo.classList.add('dm-item');
    if (!conversation.is_read) {
        convo.classList.add('dm-item-unread');
    }
    
    const authorInline = document.createElement('div');
    authorInline.classList.add('author-info-inline');
    const profilePic = document.createElement('img');
    profilePic.classList.add('profile-pic-mini');
    if (conversation.user2.profile_pic === '') {
        profilePic.src = '/static/img/no-profile-pic.jpg';
    } else {
        profilePic.src = '/action/profile-pic?encoded_username=' + conversation.user2.encoded
    }

    const convoMsg = document.createElement('div');
    convoMsg.innerHTML = '<strong>' + conversation.user2.username + '</strong> sent you a message';

    const openDMButton = document.createElement('button');
    openDMButton.innerText = 'Open';
    openDMButton.classList.add('open-dm-button');
    openDMButton.setAttribute('data-conversation-id', conversation.id);
    openDMButton.setAttribute('data-from', fromUser);
    const readHidden = document.createElement('input');
    readHidden.type = 'hidden';
    readHidden.value = conversation.is_read;

    authorInline.append(profilePic);
    authorInline.append(convoMsg);
    convo.append(authorInline);
    convo.append(readHidden);
    convo.append(openDMButton);

    return convo;
}
