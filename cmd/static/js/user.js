import * as wsoc from './wsManager.js';
import { session_username } from './auth.js';

$(document).ready(function() {
    let follow_button = $('#follow-button');
    let unfollow_button = $('#unfollow-button');
    let profilePicInput = $('#profile-pic-input');
    let confirmButton = $('#profile-pic-confirm-button');

    follow_button.on('click', function() {
        followUser();
    })

    unfollow_button.on('click', function() {
        unfollowUser();
    })

    profilePicInput.on('change', function() {
        checkFileSelected();
    });

    confirmButton.on('click', function() {
        addProfilePic();
    });

    function checkFileSelected() {
        if (profilePicInput[0].files.length > 0) {
            confirmButton.show();
        } else {
            confirmButton.hide();
        }
    }
});

document.addEventListener('DOMContentLoaded', function() {
    const openButton = document.getElementsByClassName('start-dm-button')[0];
    const dmWindow = document.getElementById('dm-window');
    const closeButton = document.getElementById('close-dm-btn');
    const dmContent = document.getElementById('dm-content');
    const inputField = document.getElementById('dm-input-field');
    const sendButton = document.getElementById('send-message-btn');

    openButton.addEventListener('click', function() {
        dmWindow.classList.toggle('hidden');
        // Optionally, populate the DM window with the active conversation
        populateConversation();
    });

    closeButton.addEventListener('click', function() {
        dmWindow.classList.add('hidden');
    });

    sendButton.addEventListener('click', function() {
        const message = inputField.value.trim();
        if (message) {
            addMessage(message, 'outgoing');
            inputField.value = '';
        }
    });
    function populateConversation() {
        // Clear existing content
        dmContent.innerHTML = '';

        // Load messages dynamically here
        // Example messages
        const messages = [
            { text: 'Hello! How can I help you?', type: 'incoming' },
            { text: 'I need information about your services.', type: 'outgoing' }
        ];

        messages.forEach(msg => addMessage(msg.text, msg.type));
    }

    function addMessage(text, type) {
        const messageElement = document.createElement('div');
        messageElement.classList.add(type === 'incoming' ? 'dm-message-incoming' : 'dm-message-outgoing');
        messageElement.textContent = text;
        dmContent.appendChild(messageElement);
        dmContent.scrollTop = dmContent.scrollHeight; // Scroll to bottom
    }
});


// AJAX calls

let encoded = $('#encoded-user').val();

function followUser() {
    $.ajax({
        url: '/action/user/' + encoded + '/follow',
        method: 'POST',
        success: function() {
            const message = {
                from_username: session_username,
                type: wsoc.TYPE_FOLLOW_REQUEST,
                msg: wsoc.MSG_FOLLOW_REQUEST,
                resource_id: encoded,
                parent_id: ''
            };
            wsoc.sendWSmsg(message);
            location.reload()
        },
        error: function(xhr) {
            const errorMessage = xhr.responseText;
            window.location.href = '/error?message=' + encodeURIComponent(errorMessage);
        }
    })
    return false;
}

function unfollowUser() {
    $.ajax({
        url: '/action/user/' + encoded + '/unfollow' + (session_username ? '?requester=' + encodeURIComponent(session_username) : ''),
        method: 'DELETE',
        success: function() {
            location.reload();
        },
        error: function(xhr) {
            const errorMessage = xhr.responseText;
            window.location.href = '/error?message=' + encodeURIComponent(errorMessage);
        }
    });
    return false;
}

function addProfilePic() {
    let form = $('#profile-pic-form')[0];
    let formData = new FormData(form);
    $.ajax({
        url: '/action/user/' + encoded + '/profile-pic',
        method: 'POST',
        data: formData,
        processData: false,  // Prevent jQuery from automatically transforming the data into a query string
        contentType: false,  // Let the browser set the content type, including boundary
        success: function(res) {
            location.reload()
        },
        error: function(xhr) {
            const errorMessage = xhr.responseText;
            window.location.href = '/error?message=' + encodeURIComponent(errorMessage);
        }
    })
    return false;
}