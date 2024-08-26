
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