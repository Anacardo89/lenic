

// Edit comment textarea hide/show
let edit_comment_buttons = document.getElementsByClassName('comment-editor-button');
for (let i = 0; i < edit_comment_buttons.length; i++) {
    const button = edit_comment_buttons[i];
    button.addEventListener('click', function() {
        const id = edit_comment_buttons[i].getAttribute('data-id');
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

// Rate comment buttons behaviour
let rate_comment_up_buttons = document.getElementsByClassName('comment-rate-up-button');
for (let i = 0; i < rate_comment_up_buttons.length; i++) {
    const button = rate_comment_up_buttons[i];
    button.addEventListener('click', function() {
        id = rate_comment_up_buttons[i].getAttribute('data-id');
        let commentElement = $('.comment[data-id="' + id + '"]');
        if (commentElement) {
            rateCommentUp(commentElement);
        }
    })
}

let rate_comment_down_buttons = document.getElementsByClassName('comment-rate-up-button');
for (let i = 0; i < rate_comment_down_buttons.length; i++) {
    const button = rate_comment_down_buttons[i];
    button.addEventListener('click', function() {
        id = rate_comment_down_buttons[i].getAttribute('data-id');
        let commentElement = $('.comment[data-id="' + id + '"]');
        console.log(commentElement);
        if (commentElement) {
            rateCommentDown(commentElement);
        }
    })
}



// Delete comment modal behaviour
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
        let commentElement = $('.comment[data-id="' + commentIdToDelete + '"]');
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


// AJAX calls

// Edit Comment
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

// Delete comment
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

// Rate comment Up
function rateCommentUp(el) {
    let id = $(el).find('.comment_id').val();
    let guid = $(el).find('.post_guid').val();
    $.ajax({
        url: '/action/post/' + guid + '/comment/' + id + '/up',
        method: 'POST',
        data: ({
            rating: 1
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

// Rate comment Down
function rateCommentDown(el) {
    let id = $(el).find('.comment_id').val();
    let guid = $(el).find('.post_guid').val();
    $.ajax({
        url: '/action/post/' + guid + '/comment/' + id + '/down',
        method: 'POST',
        data: ({
            rating: -1
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