import { connectWS } from "./wsManager.js";

// Login
$(document).ready(function() {
    $('#login-form').on('submit', login);  
});

function login(el) {
    el.preventDefault();
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
            connectWS(localStorage.getItem('user_name'));
            window.location.href = '/user/' + encoded + '/feed';
        },
        error: function() {
            console.error('Login failed');
        }
    });
    return false;
}