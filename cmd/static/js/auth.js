import * as wsoc from './wsManager.js';


export let session_username = $('#session-username').val();
export let session_encoded = $('#session-encoded').val();


$(document).ready(function() {
    const userName = localStorage.getItem('user_name');
    if (userName) {
        wsoc.connectWS(userName);
    }
});

window.addEventListener('beforeunload', function() {
    if (ws && ws.readyState === WebSocket.OPEN) {
        wsoc.ws.close(1000, "Page unload");
    }
});


// Logout
$(document).ready(function() {
    $('#logout-button').on('click', logout);
});

function logout() {
    $.ajax({
        url: '/action/logout',
        method: 'POST',
        success: function() {
            console.log('Logout successful'); 
            localStorage.removeItem('user_name');
            wsoc.closeWS();
            window.location.href = '/home';
        },
        error: function(status, error) {
            console.error('Logout failed:', status, error);
            localStorage.removeItem('user_name');
            wsoc.closeWS();
            window.location.href = '/home';
        }
    });
    return false;
}