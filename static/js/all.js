
title = document.getElementById('home-link');
title.addEventListener('click', function() {
    location.href = "/home";
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