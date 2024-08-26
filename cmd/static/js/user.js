
let encodedElem = document.getElementById('encoded-user');
let encoded = encodedElem.getAttribute('value');

let follow_button = document.getElementById('follow-button');
if (follow_button !== null) {
    follow_button.addEventListener('click', function() {
        followUser();
    })
}

let unfollow_button = document.getElementById('unfollow-button');
if (unfollow_button !== null) {
    unfollow_button.addEventListener('click', function() {
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
    var fileInput = document.getElementById('profile-pic-input');
    var confirmButton = document.getElementById('confirm-button');
    
    if (fileInput.files.length > 0) {
        // If a file is selected, show the button
        confirmButton.style.display = 'inline-block';
    } else {
        // If no file is selected, hide the button
        confirmButton.style.display = 'none';
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