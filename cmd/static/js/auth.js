let ws;

// Login
function login(el) {
    const userName = $('.login-field input[name="user_name"]').val();
    const encoded = btoa(userName);
    const userPassword = $('.password-field input[name="user_password"]').val();
    $.ajax({
        url: '/action/login',
        method: 'POST',
        contentType: 'application/json',
        data: JSON.stringify({
            user_name: userName,
            user_password: userPassword
        }),
        success: function() {
            localStorage.setItem('user_name', userName);
            connectWebSocket(localStorage.getItem('user_name'));
            window.location.href = '/user/' + encoded + '/feed';
        },
        error: function() {
            console.error('Login failed');
        }
    });
    return false;
}

// Logout
function logout() {
    $.ajax({
        url: '/action/logout',
        method: 'POST',
        success: function() {
            console.log('Logout successful'); 
            localStorage.removeItem('user_name');
            closeWebSocket();
            window.location.href = '/home';
        },
        error: function(status, error) {
            console.error('Logout failed:', status, error);
            localStorage.removeItem('user_name');
            closeWebSocket();
            window.location.href = '/home';
        }
    });
    return false;
}

// WebSocket connection
function connectWebSocket(user_name) {
    if (ws && ws.readyState === WebSocket.OPEN) {
        console.log('WebSocket connection already open');
        return;
    }

    const wsUrl = `wss://${window.location.host}/ws?user_id=${user_name}`;
    ws = new WebSocket(wsUrl);

    ws.onopen = function() {
        console.log('WebSocket connection established');
    };

    ws.onmessage = function(event) {
        console.log('Message from server:', event.data);
    };

    ws.onerror = function(error) {
        console.error('WebSocket error:', error);
    };

    ws.onclose = function(event) {
        console.log('WebSocket connection closed:', event);
    };
}

// Close WebSocket connection
function closeWebSocket() {
    if (ws && ws.readyState === WebSocket.OPEN) {
        ws.close();
    }
    ws = null;
}

$(document).ready(function() {
    const userName = localStorage.getItem('user_name');
    if (userName) {
        connectWebSocket(userName);
    }
});

window.addEventListener('beforeunload', function(event) {
    closeWebSocket();
});