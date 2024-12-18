import {positionSuggestionBox, insertAtCaret} from './tag.js';
import { session_encoded } from './auth.js';

// Add post
$(document).ready(function() {
    $('#post-button').on('click', addPost);
});

function addPost() {
    let formData = new FormData();

    const imageFile = $('#post-image')[0].files[0];
    const postTitle = $('#post-title-input').val();
    const postContent = $('#post-textarea').val();
    const visibility = $('input[name="post-visibility"]:checked').val();

    if (imageFile) {
        formData.append('post-image', imageFile);
    }
    formData.append('post-title', postTitle);
    formData.append('post-content', postContent);
    formData.append('post-visibility', visibility);
    
    $.ajax({
        url: '/action/post',
        method: 'POST',
        data: formData,
        processData: false,
        contentType: false, 
        success: function(res) {
            window.location.href = '/user/' + session_encoded + '/feed'
        },
        error: function(xhr) {
            const errorMessage = xhr.responseText;
            alert(errorMessage);
        }
    })
    return false;
}

// Tag User
$(document).ready(function() {

    const post_textArea = $('#post-textarea');
    const suggestionBox = $('#suggestionBox');
    post_textArea.on('keyup', function(event) {
        const cursorPosition = event.target.selectionStart;
        const textBeforeCursor = event.target.value.slice(0, cursorPosition)
        const mentionMatch = textBeforeCursor.match(/@(\w*)$/);
        if (mentionMatch) {
            const searchText = mentionMatch[1];
            if (searchText.length > 0) {
                fetchUserSuggestions(searchText);
                positionSuggestionBox(post_textArea, suggestionBox);
            }
        } else {
            suggestionBox.css('display', 'none');
        }
    });

    function fetchUserSuggestions(query) {
        $.ajax({
            url: '/action/search/user?username=' + encodeURIComponent(query),
            method: 'GET',
            success: function(data) {
                console.log('making request')
                updateResults(data);
            },
            error: function(xhr, status, error) {
                console.error('Error:', error);
            }
        });
    }

    function updateResults(data) {
        clearResults();

        if (Array.isArray(data)) {
            $.each(data, function(index, item) {
                const $resultItem = makeSuggestionResult(item);
                suggestionBox.append($resultItem);
            });
        } else {
            console.error('Expected an array but received:', data);
        }
    }

    function clearResults() {
        suggestionBox.empty();
    }
});

function makeSuggestionResult(user) {
    const result = document.createElement('div');
    result.classList.add('suggestion-item');
    const authorInline = document.createElement('div');
    authorInline.classList.add('author-info-inline');
    const profilePic = document.createElement('img');
    profilePic.classList.add('profile-pic-mini');
    if (user.profile_pic === '') {
        profilePic.src = '/static/img/no-profile-pic.jpg';
    } else {
        profilePic.src = '/action/profile-pic?user-encoded=' + user.encoded
    }
    const username = document.createElement('div');
    username.innerHTML = '<strong>' + user.username + '</strong>';
    authorInline.append(profilePic);
    authorInline.append(username);
    result.append(authorInline);

    result.addEventListener('click', function(event) {
        if (event.currentTarget === result || event.currentTarget.contains(event.target)) {
            let usernameElement = result.querySelector('.author-info-inline strong');
            if (usernameElement) {
                let selectedUser = usernameElement.textContent;
                insertAtCaret('post-textarea', '@' + selectedUser);
            }
        }
        result.parentElement.style.display = 'none';
    });
    return result;
}

const fileInput = $('#post-image');
const imageLabel = $('#post-image-label').find('i');

fileInput.on('change', () => {
    const rawInput = fileInput[0]; 

    if (rawInput?.files && rawInput.files.length > 0) {
        imageLabel.css('background-color', 'green');
    } else {
        imageLabel.css('background-color', '#333');
    }
});