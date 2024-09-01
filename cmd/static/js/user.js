import * as wsoc from './wsManager.js';
import { session_username } from './auth.js';

let encodedElem = document.getElementById('encoded-user');
let encoded = encodedElem.getAttribute('value');
console.log(encoded);

let follow_button = $('#follow-button');
console.log(follow_button);
if (follow_button !== null) {
    follow_button.on('click', function() {
        followUser();
    })
}

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
        error: function(err) {
            console.error("Error:", err);
        }
    })
    return false;
}


let unfollow_button = $('#unfollow-button');
if (unfollow_button !== null) {
    unfollow_button.on('click', function() {
        unfollowUser();
    })
}

function unfollowUser() {
    $.ajax({
        url: '/action/user/' + encoded + '/unfollow' + (session_username ? '?requester=' + encodeURIComponent(session_username) : ''),
        method: 'DELETE',
        success: function() {
            location.reload();
        },
        error: function(err) {
            console.error("Error:", err);
        }
    });
    return false;
}


let profilePicInput = $('#cprofile-pic-input');
if (profilePicInput !== null) {
    profilePicInput.on('change', function() {
        checkFileSelected();
    })
}

let confirm_button = $('#profile-pic-confirm-button');
if (confirm_button !== null) {
    confirm_button.on('click', function() {
        addProfilePic();
    })
}

function checkFileSelected() {
    let fileInput = $('#profile-pic-input');
    if (fileInput.files.length > 0) {
        confirm_button.show();
    } else {
        confirm_button.hide();
    }
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
        error: function(err) {
            console.error("Error:", err);
        }
    })
    return false;
}