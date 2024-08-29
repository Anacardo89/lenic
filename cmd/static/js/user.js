
let encodedElem = document.getElementById('encoded-user');
let encoded = encodedElem.getAttribute('value');

let follow_button = $('#follow-button');
if (follow_button !== null) {
    follow_button.on('click', function() {
        followUser();
    })
}

let unfollow_button = $('#unfollow-button');
if (unfollow_button !== null) {
    unfollow_button.on('click', function() {
        unfollowUser();
    })
}

function followUser() {
    $.ajax({
        url: '/action/user/' + encoded + '/follow',
        method: 'POST',
        success: function(res) {
            location.reload()
        },
        error: function(err) {
            console.error("Error:", err);
        }
    })
    return false;
}


function unfollowUser() {
    $.ajax({
        url: '/action/user/' + encoded + '/unfollow',
        method: 'POST',
        success: function(res) {
            location.reload()
        },
        error: function(err) {
            console.error("Error:", err);
        }
    })
    return false;
}

function checkFileSelected() {
    var fileInput = $('#profile-pic-input')[0];
    var confirmButton = $('#confirm-button');
    if (fileInput.files.length > 0) {
        confirmButton.css('display', 'inline-block');
    } else {
        confirmButton.hide();
    }
}

function addProfilePic() {
    let form = $('#profile-pic-form')[0];
    let formData = new FormData(form);
    $.ajax({
        url: '/action/user/' + encoded + '/profile-pic',
        method: 'POST',
        data: formData,
        processData: false,  // Prevent jQuery from automatically transforming the data into a query string
        contentType: false,  // Let the browser set the content type, including boundary
        success: function(res) {
            location.reload()
        },
        error: function(err) {
            console.error("Error:", err);
        }
    })
    return false;
}