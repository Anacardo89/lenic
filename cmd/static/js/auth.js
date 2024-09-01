import * as wsoc from './wsManager.js';


export const session_username = $('#session-username').val();
export const session_encoded = $('#session-encoded').val();



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


// Notifs
$(document).ready(function() {
    const $notifButton = $('.notif-button');
    const $notifDropdown = $('.notif-dropdown');

    // Toggle dropdown visibility on button click
    $notifButton.on('click', function(event) {
        event.stopPropagation(); // Prevent the click event from propagating to the document
        $notifDropdown.toggle();
    });

    $(document).on('click', function(event) {
    if (!$notifButton.is(event.target) && $notifButton.has(event.target).length === 0 &&
        !$notifDropdown.is(event.target) && $notifDropdown.has(event.target).length === 0) {
        $notifDropdown.hide();
    }
    });

    // Example condition to show the dot (replace with your actual condition)
    const hasUnreadNotifications = true; // Example condition
    if (hasUnreadNotifications) {
        $notifButton.addClass('show-dot');
    } else {
        $notifButton.removeClass('show-dot');
    }
});

// DMs
$(document).ready(function() {
    const $dmButton = $('.dm-button');
    const $dmDropdown = $('.dm-dropdown');

    // Toggle dropdown visibility on button click
    $dmButton.on('click', function(event) {
        event.stopPropagation();
        $dmDropdown.toggle();
    });

    $(document).on('click', function(event) {
    if (!$dmButton.is(event.target) && $dmButton.has(event.target).length === 0 &&
        !$dmDropdown.is(event.target) && $dmDropdown.has(event.target).length === 0) {
        $dmDropdown.hide();
    }
    });

    // Example condition to show the dot (replace with your actual condition)
    const hasUnreadDMs = true; // Example condition
    if (hasUnreadDMs) {
        $dmButton.addClass('show-dot');
    } else {
        $dmButton.removeClass('show-dot');
    }
});


// Fetch notifications
$(document).ready(function() {
    const $container = $('.notif-body');
    let offset = 0;
    const limit = 50;
    let loading = false;
    let hasMore = true;

    // Function to fetch notifications
    function fetchNotifications() {

        if (!session_encoded) {
            console.error('session_encoded is not defined');
            return;
        }
        if (loading || !hasMore) return;

        loading = true;

        $.ajax({
            url: '/action/user/'+ session_encoded +'/notifications?offset=' + offset + '&limit=' + limit,
            method: 'GET',
            dataType: 'json',
            success: function(data) {
                if (data !== null) {
                    if (data.length > 0) {
                        console.log(data);
                        appendNotifications(data);
                        hasMore = data.hasMore;
                        offset += limit;
                    } else {
                        hasMore = false;
                    }
                }
            },
            error: function(textStatus, errorThrown) {
                console.error('Error fetching notifications:', textStatus, errorThrown);
            },
            complete: function() {
                loading = false;
            }
        });
    }

    // Function to append notifications to the container
    function appendNotifications(notifications) {
        const notifContainer = $('.notif-body');
        const notifButton = $('.notif-button');
        notifications.forEach(function(notification) {
            let notif = null;
            switch (notification.type) {
            case wsoc.TYPE_COMMENT_RATE:
                notif = makeCommentRateNotif(notification);
                if (!notification.is_read) {
                    notifButton.css('--notif-display', 'block');
                }
                break;
            case wsoc.TYPE_POST_RATE:
                notif = makePostRateNotif(notification);
                if (!notification.is_read) {
                    notifButton.css('--notif-display', 'block');
                }
                break;
            case wsoc.TYPE_FOLLOW_REQUEST:
                notif = makeFollowRequestNotif(notification);
                if (!notification.is_read) {
                    notifButton.css('--notif-display', 'block');
                }
                break;
            }
            notifContainer.append(notif);
        });
    }

    // Scroll event handler
    function handleScroll() {
      const scrollHeight = $container[0].scrollHeight;
      const scrollTop = $container.scrollTop();
      const clientHeight = $container.height();

      if (scrollHeight - scrollTop === clientHeight) {
        fetchNotifications();
      }
    }

    // Attach scroll event listener
    $container.on('scroll', handleScroll);

    fetchNotifications();
});

export function makeCommentRateNotif(notification) {
    const postGuid = notification.parent_id;
    const notif = document.createElement('div');
    notif.classList.add('notif-item');
    const authorInline = document.createElement('div');
    authorInline.classList.add('author-info-inline');
    const profilePic = document.createElement('img');
    profilePic.classList.add('profile-pic-mini');
    if (notification.fromuser.profile_pic === '') {
        profilePic.src = '/static/img/no-profile-pic.jpg';
    } else {
        profilePic.src = '/action/profile-pic?user-encoded=' + notification.fromuser.encoded
    }
    const notifMsg = document.createElement('div');
    notifMsg.innerHTML = '<strong>' + notification.fromuser.username + '</strong> ' + notification.msg;
    const idHidden = document.createElement('input');
    idHidden.type = 'hidden';
    idHidden.value = notification.id;
    const readHidden = document.createElement('input');
    readHidden.type = 'hidden';
    readHidden.value = notification.is_read;
    authorInline.append(profilePic);
    authorInline.append(notifMsg);
    authorInline.append(idHidden);
    authorInline.append(readHidden);
    notif.append(authorInline);

    notif.addEventListener('click', function() {
        $.ajax({
            url: '/action/user/' + session_encoded + '/notifications/' + idHidden.value + '/read',
            method: 'PUT',
            success: function() {
                window.location.href = '/post/' +  postGuid;
            },
            error: function(err) {
                console.error("Error:", err);
            }
        });
    });
    return notif;
}

export function makePostRateNotif(notification) {
    const postGuid = notification.resource_id;
    const notif = document.createElement('div');
    notif.classList.add('notif-item');
    const authorInline = document.createElement('div');
    authorInline.classList.add('author-info-inline');
    const profilePic = document.createElement('img');
    profilePic.classList.add('profile-pic-mini');
    if (notification.fromuser.profile_pic === '') {
        profilePic.src = '/static/img/no-profile-pic.jpg';
    } else {
        profilePic.src = '/action/profile-pic?user-encoded=' + notification.fromuser.encoded
    }
    const notifMsg = document.createElement('div');
    notifMsg.innerHTML = '<strong>' + notification.fromuser.username + '</strong> ' + notification.msg;
    const idHidden = document.createElement('input');
    idHidden.type = 'hidden';
    idHidden.value = notification.id;
    const readHidden = document.createElement('input');
    readHidden.type = 'hidden';
    readHidden.value = notification.is_read;
    authorInline.append(profilePic);
    authorInline.append(notifMsg);
    authorInline.append(idHidden);
    authorInline.append(readHidden); 
    notif.append(authorInline);

    notif.addEventListener('click', function() {
        $.ajax({
            url: '/action/user/' + session_encoded + '/notifications/' + idHidden.value + '/read',
            method: 'PUT',
            success: function() {
                window.location.href = '/post/' +  postGuid;
            },
            error: function(err) {
                console.error("Error:", err);
            }
        });
    });
    return notif;
}

export function makeFollowRequestNotif(notification) {
    const encoded = notification.resource_id;
    const fromUser = notification.fromuser.username;
    const notif = document.createElement('div');
    notif.classList.add('notif-item');
    
    const authorInline = document.createElement('div');
    authorInline.classList.add('author-info-inline');

    const profilePicLink = document.createElement('a');
    profilePicLink.href = '/user/' + notification.fromuser.encoded;
    
    const profilePic = document.createElement('img');
    profilePic.classList.add('profile-pic-mini');
    if (notification.fromuser.profile_pic === '') {
        profilePic.src = '/static/img/no-profile-pic.jpg';
    } else {
        profilePic.src = '/action/profile-pic?user-encoded=' + notification.fromuser.encoded
    }
    profilePicLink.append(profilePic);

    const notifMsg = document.createElement('div');
    notifMsg.innerHTML = '<a href="/user/' + notification.fromuser.encoded + '"><strong>' + notification.fromuser.username + '</strong> ' + notification.msg;
    
    const idHidden = document.createElement('input');
    idHidden.type = 'hidden';
    idHidden.value = notification.id;
    
    const readHidden = document.createElement('input');
    readHidden.type = 'hidden';
    readHidden.value = notification.is_read;

    const acceptRequestButton = document.createElement('button');
    acceptRequestButton.innerText = 'Accept';

    const refuseRequestButton = document.createElement('button');
    refuseRequestButton.innerText = 'Refuse';

    const buttonsDiv = document.createElement('div');
    buttonsDiv.classList.add('request-buttons');
    buttonsDiv.append(refuseRequestButton);
    buttonsDiv.append(acceptRequestButton);
    
    authorInline.append(profilePicLink);
    authorInline.append(notifMsg);
    authorInline.append(idHidden); 
    authorInline.append(readHidden); 
    notif.append(authorInline);
    notif.append(buttonsDiv);

    acceptRequestButton.addEventListener('click', function() {
        $.ajax({
            url: '/action/user/' + session_encoded + '/accept',
            method: 'PUT',
            data: {
                requester: fromUser
            },
            success: function() {
                location.reload();
            },
            error: function(err) {
                console.error("Error:", err);
            }
        });
    });
    
    refuseRequestButton.addEventListener('click', function() {
        $.ajax({
            url: '/action/user/' + session_encoded + '/unfollow' + (fromUser ? '?requester=' + encodeURIComponent(fromUser) : ''),
            method: 'DELETE',
            success: function() {
                location.reload();
            },
            error: function(err) {
                console.error("Error:", err);
            }
        });
    });

    return notif;
}

// Search
$(document).ready(function() {
    let timeout;

    const search_button = $('#search-button');
    const search_input = $('#search-input')
    search_button.on('click', function() {
        if (search_input.is(':visible')) {
            search_input.hide();
        } else {
            search_input.show();
        }
    });

    search_input.on('input', function() {
        clearTimeout(timeout);
        const query = $(this).val();

        timeout = setTimeout(() => {
            if (query.length > 0) {
                sendRequest(query);
                $('.search-container').addClass('show');  // Show the dropdown
            } else {
                clearResults();
                $('.search-container').removeClass('show');  // Hide the dropdown if no query
            }
        }, 750);

    });

    function sendRequest(query) {
        $.ajax({
            url: '/action/search/user?username=' + query,
            method: 'GET',
            success: function(data) {
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
                const $resultItem = makeSearchResult(item);
                $('#search-body').append($resultItem);
            });
        } else {
            console.error('Expected an array but received:', data);
        }
    }

    function clearResults() {
        $('#search-body').empty();
    }
});

function makeSearchResult(user) {
    const result = document.createElement('div');
    result.classList.add('search-item');
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

    result.addEventListener('click', function() {
        window.location.href = '/user/' +  user.encoded;
    });
    return result;
}