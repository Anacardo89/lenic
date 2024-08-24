

// Edit post textarea hide/show
post_edit_button = $('#post-editor-button');
post_edit_button.on('click', function() {
    const edit_form = $('#post-edit');
    const post_text = $('#post-text');
    if (edit_form.css('display') === 'none' || edit_form.css('display') === '') {
        edit_form.css('display', 'block');
        post_text.css('display', 'none');
    } else {
        edit_form.css('display', 'none');
        post_text.css('display', 'block');
    }
})

// Delete post modal behaviour
const postModal = document.getElementById('modal-container-post');
const modalPostCancelBtn = document.getElementById('delete-post-sure-no');
deletePostBtn = document.querySelector('#post-deleter-button')
deletePostBtn.addEventListener('click', function() {
    postModal.style.display = "block";
});
modalPostCancelBtn.addEventListener('click', function() {
    postModal.style.display = 'none';
});
document.getElementById('delete-post-sure-yes').addEventListener('click', function() {
    deletePost(document);
    postModal.style.display = 'none'; 
});
window.addEventListener('click', function(event) {
    if (event.target === postModal) {
        postModal.style.display = 'none';
    }
});


// AJAX calls

// Add post
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
            window.location.href = '/'
        },
        error: function(err) {
            console.error("Error:", err);
        }
    })
    return false;
}

// Edit Post
function editPost(el) {
    let guid = $(el).find('.post_id').val();
    let edited_post = $(el).find('.edit_post').val();
    $.ajax({
        url: '/action/post/' + guid,
        method: 'PUT',
        data: ({
            post: edited_post
        }),
        success: function(res) {
            location.reload()
        },
        error: function(err) {
            console.error("Error:", err);
        }
    })
    return false;
}

// Delete Post
function deletePost(el) {
    let guid = $(el).find('.post_id').val();
    $.ajax({
        url: '/action/post/' + guid,
        method: 'DELETE',
        success: function(res) {
            window.location.href = '/'
        },
        error: function(err) {
            console.error("Error:", err);
        }
    })
    return false;
}