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
            if (data.length > 0) {
            appendNotifications(data);
            hasMore = data.hasMore;
            offset += limit;
            } else {
            hasMore = false;
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
        notifications.forEach(function(notification) {
            let notif = null;
            switch (notification.type) {
            case 'rate_comment':
                notif = makeCommentNotif(notification);
                break;
            case 'rate_post':
                notif = makePostNotif(notification);
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

function makeCommentNotif(notification) {
    const notif = document.createElement('div');
    notif.classList.add('notif-item');
    let postGuid = notification.parent_id;
    const notifLink = document.createElement('a');
    notifLink.href = '/post/' + postGuid;
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
    authorInline.append(profilePic);
    authorInline.append(notifMsg);
    authorInline.style.color = '#333';
    notifLink.append(authorInline);
    notif.append(notifLink);
    return notif;
}

function makePostNotif(notification) {
    const notif = document.createElement('div');
    notif.classList.add('notif-item');
    let postGuid = notification.resource_id;
    const notifLink = document.createElement('a');
    notifLink.href = '/post/' + postGuid;
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
    authorInline.append(profilePic);
    authorInline.append(notifMsg);
    authorInline.style.color = '#333';
    notifLink.append(authorInline);
    notif.append(notifLink);
    return notif;
}
