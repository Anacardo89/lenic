import { connectWS } from "./wsManager.js";

// Login
$(document).ready(function() {
    $('#login-button').on('click', login);  
});

function login(el) {
    el.preventDefault();
    const userName = $('#login-user').val();
    const encoded = btoa(userName);
    const userPassword = $('#login-password').val();
    $.ajax({
        url: '/action/login',
        method: 'POST',
        contentType: 'application/json',
        data: JSON.stringify({
            user_name: userName,
            user_password: userPassword
        }),
        success: function(xhr) {
            localStorage.setItem('user_name', userName);
            connectWS(localStorage.getItem('user_name'));
            window.location.href = '/user/' + encoded + '/feed';
        },
        error: function(xhr) {
            const errorMessage = xhr.responseText;
            alert(errorMessage);
        }
    });
    return false;
}