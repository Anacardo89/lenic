

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