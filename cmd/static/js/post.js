import {session_username, session_encoded} from "./auth.js";
import  * as wsoc from './wsManager.js';

export let guid = $('#post-guid').val();

// Edit post textarea hide/show
$(document).ready(function() {
    const post_edit_button = $('#post-editor-button');
    post_edit_button.on('click', function() {
        const edit_form = $('#post-edit-container');
        const post_text = $('#post-text');
        if (edit_form.is(':hidden')) {
            edit_form.show();
            post_text.hide();
        } else {
            edit_form.hide();
            post_text.show();
        }
    });
});

// Delete post modal behaviour
const postModal = $('#modal-container-post');
let  deletePostBtn = $('#post-deleter-button')
if (deletePostBtn !== null) {
    deletePostBtn.on('click', function() {
        postModal.show();
    });
}
let modalPostCancelBtn = $('#delete-post-sure-no');
if (modalPostCancelBtn !== null) {
    modalPostCancelBtn.on('click', function() {
        postModal.hide();
    });
}
let modalPostDeleteBtn = $('#delete-post-sure-yes');
if (modalPostDeleteBtn !== null) {
    modalPostDeleteBtn.on('click', function() {
        deletePost(document);
        postModal.hide();
    });
}

$(window).on('click', function(event) {
    if ($(event.target).is('#postModal')) { // Use jQuery to compare the event target
        $('#postModal').hide(); // Hide the modal
    }
});

// Rate post buttons behaviour
let rate_post_up_button = $('#post-rate-up-button');
if (rate_post_up_button !== null) {
    rate_post_up_button.on('click', function() {
        ratePostUp();
    });
}
let rate_post_down_button = $('#post-rate-down-button');
if (rate_post_down_button !== null) {
    rate_post_down_button.on('click', function() {
        ratePostDown();
    });
}

let rate_post_hidden = $('#post-rating-hidden');
if (rate_post_hidden !== null) {
let postUserRating = rate_post_hidden.val();
    if (postUserRating > 0) {
        let rate_up_button = rate_post_hidden.prev();
        rate_up_button.css('color', 'orange');
    } else if (postUserRating < 0) {
        let rate_down_button = rate_post_hidden.next();
        rate_down_button.css('color', 'orange');
    }
}


// AJAX calls

// Edit Post
$(document).ready(function() {
    $('#post-edit-form').on('submit', editPost);  
});

function editPost(el) {
    el.preventDefault();
    let form = $(el.currentTarget);
    let edited_title = form.find('#edit-post-title').val();
    let edited_post = form.find('#edit-post').val();
    let edited_visibility = form.find('input[name="post-visibility"]:checked').val();
    $.ajax({
        url: '/action/post/' + guid,
        method: 'PUT',
        data: ({
            title: edited_title,
            content: edited_post,
            visibility: edited_visibility
        }),
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

// Delete Post
function deletePost() {
    $.ajax({
        url: '/action/post/' + guid,
        method: 'DELETE',
        success: function(res) {
            window.location.href = '/user/' + session_encoded + '/feed'
        },
        error: function(xhr) {
            const errorMessage = xhr.responseText;
            window.location.href = '/error?message=' + encodeURIComponent(errorMessage);
        }
    })
    return false;
}

// Rate post Up
function ratePostUp() {
    $.ajax({
        url: '/action/post/' + guid + '/up',
        method: 'POST',
        success: function() {
            const message = {
                from_username: session_username,
                type: wsoc.TYPE_POST_RATE,
                msg: wsoc.MSG_POST_RATE,
                resource_id: guid
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

// Rate post Down
function ratePostDown() {
    $.ajax({
        url: '/action/post/' + guid + '/down',
        method: 'POST',
        success: function() {
            const message = {
                from_username: session_username,
                type: wsoc.TYPE_POST_RATE,
                msg: wsoc.MSG_POST_RATE,
                resource_id: guid
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