
// Logout
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

// Add a post
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

// Edit comment button behaviour
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

function editComment(el) {
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
        },
        error: function(err) {
            console.error("Error:", err);
        }
    })
    return false;
}


// Delete comment button behaviour
function deleteComment(el) {
    let id = el.find('.comment_id').val();
    let guid = $(el).closest('.post_guid').val();
    $.ajax({
        url: '/action/post/' + guid + '/comment/' + id,
        method: 'DELETE',
        success: function(res) {
            location.reload()
        },
        error: function(err) {
            console.error("Error:", err);
        }
    })
    return false;
}

const commentModal = document.getElementById("modal-container-comment");
let commentIdToDelete = null;
document.querySelectorAll('.comment-deleter-button').forEach(function(button) {
    button.onclick = function() {
        commentIdToDelete = this.getAttribute('data-id');
        commentModal.style.display = "block";
    }
});

const modalCancelBtn = document.getElementById("delete-comment-sure-no");
modalCancelBtn.onclick = function() {
    commentModal.style.display = "none";
    commentIdToDelete = null;
}

document.getElementById('delete-comment-sure-yes').onclick = function() {
    if (commentIdToDelete !== null) {
        var commentElement = $('.comment[data-id="' + commentIdToDelete + '"]');
        if (commentElement) {
            deleteComment(commentElement);
        }
        commentModal.style.display = "none"; 
        commentIdToDelete = null;
    }
}

window.onclick = function(event) {
    if (event.target == commentModal) {
        commentModal.style.display = "none";
        commentIdToDelete = null;
    }
}

// Delete post button behaviour
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

// Edit post button behaviour
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