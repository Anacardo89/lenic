
comment_buttons = document.getElementsByClassName('comment-editor-button')
for (let i = 0; i < comment_buttons.length; i++) {
    const button = comment_buttons[i];
    button.addEventListener('click', function() {
    const id = comment_buttons[i].getAttribute('data-id');
    editor_id = 'comment-editor-'+id
    const edit_form = document.getElementById(editor_id);
        if (edit_form.style.display === 'none' || edit_form.style.display === '') {
            edit_form.style.display = 'block';
        } else {
            edit_form.style.display = 'none';
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
            location.reload()
        },
        error: function(err) {
            console.error("Error:", err);  // Log any error for debugging
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