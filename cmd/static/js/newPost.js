import { session_encoded } from './auth.js';

// Add post
$(document).ready(function() {
    $('#post-button').on('click', addPost);
});

function addPost() {
    let form = $('.post-form form')[0];
    let formData = new FormData(form)
    $.ajax({
        url: '/action/post',
        method: 'POST',
        data: formData,
        processData: false,
        contentType: false, 
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