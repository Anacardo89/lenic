import * as wsoc from './wsManager.js';
import * as notifs from './notifs.js';
import * as dms from './dms.js';


export const session_username = $('#session-username').val();
export const session_encoded = $('#session-encoded').val();



$(document).ready(function() {
    const userName = localStorage.getItem('user_name');
    if (userName) {
        wsoc.connectWS(userName);
    }
});

$(document).ready(function() {
    window.addEventListener('beforeunload', function() {
        if (wsoc.ws && wsoc.ws.readyState === WebSocket.OPEN) {
            wsoc.ws.close(1000, "Page unload");
        }
    });
});

$(document).ready(function() {
    $('.notif-readAll').on('click', readAllNotifs);
});

function readAllNotifs() {
    let notifs = $('.notif-item');
    notifs.each(function() {
        let notif = $(this);
        let idHidden = notif.find('input[type="hidden"]').eq(0);
        let notifId = idHidden.val();
        $.ajax({
            url: '/action/user/' + session_encoded + '/notifications/' + notifId + '/read',
            method: 'PUT',
            success: function() {
                console.log('Notification marked as read.');
            },
            error: function(xhr) {
                localStorage.removeItem('user_name');
                wsoc.closeWS();
                const errorMessage = xhr.responseText;
                window.location.href = '/error?message=' + encodeURIComponent(errorMessage);
            }
        });
    });
    setTimeout(function() {
        location.reload();
    }, 1000);
}


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
        error: function(xhr) {
            localStorage.removeItem('user_name');
            wsoc.closeWS();
            const errorMessage = xhr.responseText;
            window.location.href = '/error?message=' + encodeURIComponent(errorMessage);
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
                notif = notifs.makeCommentRateNotif(notification);
                if (!notification.is_read) {
                    notifButton.css('--notif-display', 'block');
                }
                break;
            case wsoc.TYPE_POST_RATE:
                notif = notifs.makePostRateNotif(notification);
                if (!notification.is_read) {
                    notifButton.css('--notif-display', 'block');
                }
                break;
            case wsoc.TYPE_COMMENT_TAG:
                notif = notifs.makeCommentTagNotif(notification);
                if (!notification.is_read) {
                    notifButton.css('--notif-display', 'block');
                }
                break;
            case wsoc.TYPE_POST_TAG:
                notif = notifs.makePostTagNotif(notification);
                if (!notification.is_read) {
                    notifButton.css('--notif-display', 'block');
                }
                break;
            case wsoc.TYPE_COMMENT_ON_POST:
                notif = notifs.makeCommentOnPostNotif(notification);
                if (!notification.is_read) {
                    notifButton.css('--notif-display', 'block');
                }
                break;
            case wsoc.TYPE_FOLLOW_ACCEPT:
                notif = notifs.makeFollowAcceptNotif(notification);
                if (!notification.is_read) {
                    notifButton.css('--notif-display', 'block');
                }
                break;
            case wsoc.TYPE_FOLLOW_REQUEST:
                notif = notifs.makeFollowRequestNotif(notification);
                if (!notification.is_read) {
                    notifButton.css('--notif-display', 'block');
                }
                break;
            default:
                console.warn('Unknown message type:', message.type);
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

// DMs
$(document).ready(function() {
    const $dmButton = $('.dm-button');
    const $dmDropdown = $('.dm-dropdown');

    // Toggle dropdown visibility on button click
    $dmButton.on('click', function(event) {
        event.stopPropagation(); // Prevent the click event from propagating to the document
        $dmDropdown.toggle();
    });

    $(document).on('click', function(event) {
    if (!$dmButton.is(event.target) && $dmButton.has(event.target).length === 0 &&
        !$dmDropdown.is(event.target) && $dmDropdown.has(event.target).length === 0) {
        $dmDropdown.hide();
    }
    });

    // Example condition to show the dot (replace with your actual condition)
    const hasUnreadMessage = true; // Example condition
    if (hasUnreadMessage) {
        $dmButton.addClass('show-dot');
    } else {
        $dmButton.removeClass('show-dot');
    }
});

// DM Module
const DMModule = (function() {
    const $dmcontainer = $('.dm-body');
    let offset = 0;
    const limit = 50;
    let loading = false;
    let hasMore = true;

    // Function to fetch conversations
    function fetchConversations() {
        if (!session_encoded) {
            console.error('session_encoded is not defined');
            return;
        }
        if (loading || !hasMore) return;

        loading = true;

        $.ajax({
            url: '/action/user/' + session_encoded + '/conversations?offset=' + offset + '&limit=' + limit,
            method: 'GET',
            dataType: 'json',
            success: function(data) {
                if (data !== null && data.length > 0) {
                    console.log(data);
                    appendConversations(data);
                    hasMore = data.hasMore;
                    offset += limit;
                } else {
                    hasMore = false;
                }
            },
            error: function(textStatus, errorThrown) {
                console.error('Error fetching conversations:', textStatus, errorThrown);
            },
            complete: function() {
                loading = false;
            }
        });
    }

    // Function to append conversations to the container
    function appendConversations(conversations) {
            const $dmButton = $('.dm-button');
            conversations.forEach(function(conversation) {
                let convo = dms.makeConversation(conversation);
                if (!conversation.is_read) {
                    $dmButton.css('--dm-display', 'block');
                }
                if (convo) {
                    $dmcontainer.append(convo);
                } else {
                    console.warn('Failed to create conversation element for:', conversation);
                }
            });
    }

    // Function to clear and fetch new conversations
    function clearAndFetchConversations() {
        offset = 0; // Reset offset
        hasMore = true; // Reset hasMore
        $dmcontainer.empty(); // Clear current conversations
        fetchConversations(); // Fetch new conversations
    }

    // Scroll event handler
    function handleScroll() {
        const scrollHeight = $dmcontainer[0].scrollHeight;
        const scrollTop = $dmcontainer.scrollTop();
        const clientHeight = $dmcontainer.height();

        if (scrollHeight - scrollTop === clientHeight) {
            fetchConversations();
        }
    }

    // Attach scroll event listener
    $dmcontainer.on('scroll', handleScroll);

    // Initial fetch of conversations
    fetchConversations();

    // Expose the clearAndFetchConversations method
    return {
        clearAndFetchConversations: clearAndFetchConversations
    };
})();

// Open DM buttons and DM Window populate
const DMChatModule = (function () {
    let offset = 0;
    const limit = 50;
    
    const $dmWindow = $('#dm-window');
    const $dmTitle = $('#dm-title');
    const $dmContent = $('#dm-content');
    const $sendButton = $('#send-message-btn');
    const $inputField = $('#dm-input-field');
    
    // Open DM window and fetch conversation
    function openDM(conversationId, fromUser) {
        $dmWindow.data('conversation-id', conversationId);
        $dmWindow.data('from', fromUser);
        $dmTitle.text(fromUser);
        $dmWindow.removeClass('hidden');
        fetchConversation(conversationId);
        readConversation(conversationId);
    }

    // Close the DM window
    function closeDM() {
        $dmWindow.addClass('hidden');
        $dmContent.empty();
    }

    // Mark conversation as read
    function readConversation(conversationId) {
        $.ajax({
            url: '/action/user/' + session_encoded + '/conversations/' + conversationId + '/read',
            method: 'PUT',
            success: function () {
                console.log('conversation.is_read updated');
            },
            error: function (jqXHR, textStatus, errorThrown) {
                console.error('Error updating conversation.is_read:', textStatus, errorThrown);
            }
        });
    }

    // Fetch conversation history
    function fetchConversation(conversationId) {
        $.ajax({
            url: '/action/user/' + session_encoded + '/conversations/' + conversationId + '/dms?offset=' + offset + '&limit=' + limit,
            method: 'GET',
            dataType: 'json',
            success: function (data) {
                if (data !== null) {
                    console.log('Fetched conversation data:', data);
                    $dmContent.empty();
                    data.forEach(function (message) {
                        if (message.sender.username === session_username) {
                            appendMessage(message.content, 'received');
                        } else {
                            appendMessage(message.content, 'sent');
                        }
                    });
                    $dmContent.scrollTop($dmContent[0].scrollHeight);
                }
            },
            error: function (jqXHR, textStatus, errorThrown) {
                console.error('Error fetching conversation:', textStatus, errorThrown);
            }
        });
    }

    // Append message to DM window
    function appendMessage(message, type) {
        const messageClass = type === 'sent' ? 'message-sent' : 'message-received';
        const $messageElement = $('<div></div>').addClass(messageClass).text(message);
        $dmContent.append($messageElement);
    }

    // Send a message
    function sendMessage(conversationId, message) {
        $.ajax({
            url: '/action/user/' + session_encoded + '/conversations/' + conversationId + '/dms',
            method: 'POST',
            contentType: 'application/json',
            data: JSON.stringify({ text: message }),
            success: function (response) {
                appendMessage(message, 'sent');
                $inputField.val('');
                $dmContent.scrollTop($dmContent[0].scrollHeight);
                const notif = {
                    from_username: session_username,
                    type: wsoc.TYPE_DM,
                    msg: message.content,
                    resource_id: $dmTitle.text(),
                    parent_id: ""
                };
                wsoc.sendWSmsg(notif);
            },
            error: function (jqXHR, textStatus, errorThrown) {
                console.error('Error sending message:', textStatus, errorThrown);
            }
        });
    }

    // Start a new conversation
    function startConversation(toUser) {
        $.ajax({
            url: '/action/user/' + session_encoded + '/conversations',
            method: 'POST',
            contentType: 'application/json',
            data: JSON.stringify({ to_user: toUser }),
            success: function (conversation) {
                const newConversationId = conversation.id;
                openDM(newConversationId, toUser);
                fetchConversation(newConversationId);
            },
            error: function (jqXHR, textStatus, errorThrown) {
                console.error('Error creating conversation:', textStatus, errorThrown);
            }
        });
    }

    // Bind events
    function bindEvents() {
        $(document).on('click', '.open-dm-button', function () {
            const conversationId = $(this).data('conversation-id');
            const fromUser = $(this).data('from');
            openDM(conversationId, fromUser);
        });

        $sendButton.on('click', function () {
            const message = $inputField.val().trim();
            const conversationId = $dmWindow.data('conversation-id');
            if (message) {
                sendMessage(conversationId, message);
            }
        });

        $('#close-dm-btn').on('click', closeDM);

        $('.start-dm-button').on('click', function () {
            const toUser = $('#profile-username').val();
            startConversation(toUser);
        });
    }

    // Public methods
    return {
        init: bindEvents,
        openDM: openDM,
        closeDM: closeDM,
        appendMessage: appendMessage,
        startConversation: startConversation,
        readConversation: readConversation // Expose the readConversation method
    };
})();

export { DMModule, DMChatModule };

$(document).ready(function() {
    DMChatModule.init(); 
});

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