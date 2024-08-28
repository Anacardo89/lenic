
let guidElem;
let guidElems = document.getElementsByClassName('post_id');
if (guidElems.length > 0) {
    guidElem = guidElems[0];
    guid = guidElem.getAttribute('value');
}

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
let  deletePostBtn = document.querySelector('#post-deleter-button')
if (deletePostBtn !== null) {
    deletePostBtn.addEventListener('click', function() {
        postModal.style.display = "block";
    });
}
let modalPostCancelBtn = document.getElementById('delete-post-sure-no');
if (modalPostCancelBtn !== null) {
    modalPostCancelBtn.addEventListener('click', function() {
        postModal.style.display = 'none';
    });
}
let modalPostDeleteBtn = document.getElementById('delete-post-sure-yes');
if (modalPostDeleteBtn !== null) {
    modalPostDeleteBtn.addEventListener('click', function() {
        deletePost(document);
            postModal.style.display = 'none'; 
    });
}

window.addEventListener('click', function(event) {
    if (event.target === postModal) {
        postModal.style.display = 'none';
    }
});

// Rate comment buttons behaviour
let rate_post_up_button = document.getElementById('post_rate_up_button');
if (rate_post_up_button !== null) {
    rate_post_up_button.addEventListener('click', function() {
        ratePostUp();
    });
}
let rate_post_down_button = document.getElementById('post_rate_down_button');
if (rate_post_down_button !== null) {
    rate_post_down_button.addEventListener('click', function() {
        ratePostDown();
    });
}

let rate_post_hidden = document.getElementById('post_rating_hidden');
if (rate_post_hidden !== null) {
let postUserRating = rate_post_hidden.getAttribute('value');
    if (postUserRating > 0) {
        let rate_up_button = rate_post_hidden.previousElementSibling;
        rate_up_button.style.color = 'orange';
    } else if (postUserRating < 0) {
        let rate_down_button = rate_post_hidden.nextElementSibling;
        rate_down_button.style.color = 'orange';
    }
}


// AJAX calls

// Add post
function addPost() {
    const user_encoded = $('#user-encoded').val();
    let form = $('.post-form form')[0];
    let formData = new FormData(form)
    $.ajax({
        url: '/action/post',
        method: 'POST',
        data: formData,
        processData: false,
        contentType: false, 
        success: function(res) {
            window.location.href = '/user/' + user_encoded + '/feed'
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

// Rate post Up
function ratePostUp() {
    $.ajax({
        url: '/action/post/' + guid + '/up',
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

// Rate post Down
function ratePostDown() {
    $.ajax({
        url: '/action/post/' + guid + '/down',
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