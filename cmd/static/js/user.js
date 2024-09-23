import * as wsoc from './wsManager.js';
import { session_username } from './auth.js';
import * as dms from './dms.js';

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