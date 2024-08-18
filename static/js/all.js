
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

function postPost(el) {
    let post_title = $(el).find('[name="post-title"]').val();
    let post_content = $(el).find('[name="post-content"]').val();
    $.ajax({
        url: '/api/post/',
        method: 'POST',
        data: ({
            post_title: post_title,
            post_content: post_content
        }),
        success: function(res) {
            location.reload()
        }
    })
    return false;
}

function logout(el) {
    $.ajax({
        url: '/api/logout/',
        method: 'POST',
        success: function(res) {
            window.location.href = '/home';
        }
    })
    return false;
}