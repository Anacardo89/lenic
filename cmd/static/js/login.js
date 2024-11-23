import { connectWS } from "./wsManager.js";


$(document).ready(function() {
    $('#register-button')?.on('click', register);
    $('#login-button')?.on('click', login);
    $('#forgotpasswd-button')?.on('click', forgotPasswd);
    $('#changePasswd-button')?.on('click', changePasswd);
});

// Register
function register() {
    el.preventDefault();
    const username = $('#register-user').val();
    const email = $('#register-email').val();
    const passwd = $('#register-passwd').val();
    const passwd2 = $('#register-passwd2').val();
    $.ajax({
        url: '/action/register',
        method: 'POST',
        contentType: 'application/json',
        data: JSON.stringify({
            user_name: username,
		    user_email: email,
		    user_password: passwd,
		    user_password2: passwd2
        }),
        success: function(xhr) {
            window.location.href = '/home';
        },
        error: function(xhr) {
            const errorMessage = xhr.responseText;
            alert(errorMessage);
        }
    });
}

// Login
function login(el) {
    el.preventDefault();
    const username = $('#login-user').val();
    const encoded = btoa(username);
    const userPassword = $('#login-password').val();
    $.ajax({
        url: '/action/login',
        method: 'POST',
        contentType: 'application/json',
        data: JSON.stringify({
            user_name: username,
            user_password: userPassword
        }),
        success: function(xhr) {
            localStorage.setItem('user_name', username);
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

// Forgot Passwd
function forgotPasswd(el) {
    el.preventDefault();
    const email = $('#forgotpasswd-email').val();
    $.ajax({
        url: '/action/forgot-password',
        method: 'POST',
        contentType: 'application/json',
        data: JSON.stringify({
            user_email: email
        }),
        success: function(xhr) {
            window.location.href = '/home';
        },
        error: function(xhr) {
            const errorMessage = xhr.responseText;
            alert(errorMessage);
        }
    });
    return false;
}

// Change Passwd
function changePasswd() {

    let formData = new FormData();

    const user = $('#session-username').val();
    console.log(user);
    const oldPasswd = $('#old-passwd').val();
    const newPasswd = $('#new-passwd').val();
    const newPasswd2 = $('#new-passwd2').val();

    formData.append('user_name', user);
    formData.append('old_password', oldPasswd);
    formData.append('password', newPasswd);
    formData.append('password2', newPasswd2);

    $.ajax({
        url: '/action/change-password',
        method: 'POST',
        data: formData,
        processData: false,
        contentType: false, 
        success: function(xhr) {
            window.location.href = '/home';
        },
        error: function(xhr) {
            const errorMessage = xhr.responseText;
            alert(errorMessage);
        }
    });
    return false;
}