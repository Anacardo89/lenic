import * as wsoc from './wsManager.js';
import { session_username } from './auth.js';
import * as dms from './dms.js';

$(document).ready(function() {
    let follow_button = $('#follow-button')?.on('click', followUser);
    let unfollow_button = $('#unfollow-button')?.on('click', unfollowUser);
    let profilePicInput = $('#profile-pic-input')?.on('change', checkFileSelected);
    let confirmButton = $('#profile-pic-confirm')?.on('click', addProfilePic);

    function checkFileSelected() {
        if (profilePicInput[0].files.length > 0) {
            confirmButton.show();
        } else {
            confirmButton.hide();
        }
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
            alert(errorMessage);
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
            alert(errorMessage);
        }
    });
    return false;
}

function addProfilePic() {
    let formData = new FormData();

    const imageFile = $('#profile-image')[0].files[0];

    if (imageFile) {
        formData.append('profile-image', imageFile);
    }
    
    $.ajax({
        url: '/action/user/' + encoded + '/profile-pic',
        method: 'POST',
        data: formData,
        processData: false, 
        contentType: false,  
        success: function(res) {
            location.reload()
        },
        error: function(xhr) {
            const errorMessage = xhr.responseText;
            alert(errorMessage);
        }
    })
    return false;
}

const fileInput = $('#profile-image');
const imageLabel = $('#profile-image-label').find('i');
const confirmButton = $('#profile-pic-confirm').find('button');

fileInput.on('change', () => {
    const rawInput = fileInput[0]; 

    if (rawInput?.files && rawInput.files.length > 0) {
        imageLabel.css('background-color', 'green');
        confirmButton.css('display', 'block');
    } else {
        imageLabel.css('background-color', '#333');
        confirmButton.css('display', 'none');
    }
});