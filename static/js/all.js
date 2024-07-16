
title = document.getElementById('home-link');
title.addEventListener('click', function() {
    location.href = "/home";
})

let visiv = false;
coment_editor = document.getElementById('comment-editor-button');
coment_editor.addEventListener('click', function() {
    edit_form = document.getElementById('comment-editor')
    if (!visiv){
        edit_form.style.display = 'block';
    } else{
        edit_form.style.display = 'none';
    }
    visiv = !visiv;
})

function putComment(el) {

    var id = $(el).find('.comment_id').val();
    var guid = $(el).find('.post_guid').val();
    var edited_comment = $(el).find('.edit_comment').val();
    $.ajax({
        url: '/api/post/' + guid + '/comment/' + id,
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