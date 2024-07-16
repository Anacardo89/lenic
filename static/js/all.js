
title = document.getElementById('home-link');
title.addEventListener('click', function() {
    location.href = "/home";
})

function putComment(el) {

    var id = $(el).find('.edit_id').val();
    var comments = $(el).find('.edit_comments').val();
    $.ajax({
        url: '/api/page/{{.GUID}}/comments/' + id,
        method: 'PUT',
        data: ({
            comments: comments
        }),
        success: function(res) {
            window.location.replace('/page/{{.GUID}}')
        }
    })
    return false;
}