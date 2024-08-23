
comment_buttons = document.getElementsByClassName('comment-editor-button')
for (let i = 0; i < comment_buttons.length; i++) {
    const button = comment_buttons[i];
    button.addEventListener('click', function() {
    const id = comment_buttons[i].getAttribute('data-id');
        const comment_text_id = 'comment-text-'+id;
        const comment_editor_id = 'comment-editor-'+id;
        const edit_form = document.getElementById(comment_editor_id);
        const comment_text = document.getElementById(comment_text_id);
        if (edit_form.style.display === 'none' || edit_form.style.display === '') {
            edit_form.style.display = 'block';
            comment_text.style.display = 'none';
        } else {
            edit_form.style.display = 'none';
            comment_text.style.display = 'block';
        }
    })
 }


function putComment(el) {
    let id = $(el).find('.comment_id').val();
    let guid = $(el).find('.post_guid').val();
    let edited_comment = $(el).find('.edit_comment').val();
    $.ajax({
        url: '/action/post/' + guid + '/comment/' + id,
        method: 'PUT',
        data: ({
            comment: edited_comment
        }),
        success: function(res) {
            location.reload()
        }
    })
    return false;
}

function deleteComment(el) {
    let commentDiv = $(el).closest('.comment');
    let id = commentDiv.find('.comment_id').val();
    let guid = $(el).closest('.post_guid').val();
    $.ajax({
        url: '/action/post/' + guid + '/comment/' + id,
        method: 'DELETE',
        success: function(res) {
            location.reload()
        }
    })
    return false;
}

function postPost() {
    let form = $('.post-form form')[0];
    let formData = new FormData(form)
    $.ajax({
        url: '/action/post',
        method: 'POST',
        data: formData,
        processData: false,
        contentType: false, 
        success: function(res) {
            window.location.href = '/'
        },
        error: function(err) {
            console.error("Error:", err);
        }
    })
    return false;
}

function logout(el) {
    $.ajax({
        url: '/action/logout',
        method: 'POST',
        success: function(res) {
            window.location.href = '/home';
        }
    })
    return false;
}