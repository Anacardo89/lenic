package server

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/Anacardo89/lenic/internal/middleware"
	"github.com/Anacardo89/lenic/internal/server/httphandle/api"
	"github.com/Anacardo89/lenic/internal/server/httphandle/page"
	"github.com/Anacardo89/lenic/internal/server/httphandle/redirect"
	"github.com/Anacardo89/lenic/internal/server/wshandle"
)

func NewRouter(ah *api.APIHandler, ph *page.PageHandler, wsh *wshandle.WSHandler, mw *middleware.MiddlewareHandler) http.Handler {

	r := mux.NewRouter()

	r.HandleFunc("/", redirect.RedirIndex).Schemes("http")

	// Page
	r.HandleFunc("/home", ph.Home).Schemes("http")
	r.HandleFunc("/login", ph.Login).Schemes("http")
	r.HandleFunc("/register", ph.Register).Schemes("http")
	r.HandleFunc("/error", ph.Error).Schemes("http")
	r.HandleFunc("/newPost", ph.NewPost).Schemes("http")
	r.HandleFunc("/post/{post_guid}", ph.Post).Schemes("http")
	r.HandleFunc("/user/{encoded_user_name}", ph.UserProfile).Schemes("http")
	r.HandleFunc("/user/{encoded_user_name}/feed", ph.Feed).Schemes("http")
	r.HandleFunc("/user/{encoded_user_name}/followers", ph.UserFollowers).Schemes("http")
	r.HandleFunc("/user/{encoded_user_name}/following", ph.UserFollowing).Schemes("http")
	r.HandleFunc("/forgot-password", ph.ForgotPassword).Schemes("http")
	r.HandleFunc("/recover-password/{encoded_user_name}", ph.RecoverPassword).Schemes("http")
	r.HandleFunc("/change-password/{encoded_user_name}", ph.ChangePassword).Schemes("http")

	// API
	r.HandleFunc("/action/register", ah.RegisterUser).Methods("POST").Schemes("http")
	r.HandleFunc("/action/activate/{encoded_user_name}", ah.ActivateUser).Schemes("http")
	r.HandleFunc("/action/login", ah.Login).Methods("POST").Schemes("http")
	r.HandleFunc("/action/logout", ah.Logout).Methods("POST").Schemes("http")
	// User
	r.HandleFunc("/action/search/user", ah.SearchUsers).Methods("GET").Schemes("http")
	r.HandleFunc("/action/user/{encoded_user_name}/follow", ah.RequestFollowUser).Methods("POST").Schemes("http")
	r.HandleFunc("/action/user/{encoded_user_name}/accept", ah.AcceptFollowRequest).Methods("PUT").Schemes("http")
	r.HandleFunc("/action/user/{encoded_user_name}/unfollow", ah.UnfollowUser).Methods("DELETE").Schemes("http")
	r.HandleFunc("/action/user/{encoded_user_name}/profile-pic", ah.PostProfilePic).Methods("POST").Schemes("http")
	// DM
	r.HandleFunc("/action/user/{encoded_user_name}/conversations", ah.GetConversations).Methods("GET").Schemes("http")
	r.HandleFunc("/action/user/{encoded_user_name}/conversations", ah.StartConversation).Methods("POST").Schemes("http")
	r.HandleFunc("/action/user/{encoded_user_name}/conversations/{conversation_id}/read", ah.ReadConversation).Methods("PUT").Schemes("http")
	r.HandleFunc("/action/user/{encoded_user_name}/conversations/{conversation_id}/dms", ah.GetDMs).Methods("GET").Schemes("http")
	r.HandleFunc("/action/user/{encoded_user_name}/conversations/{conversation_id}/dms", ah.SendDM).Methods("POST").Schemes("http")
	// Notif
	r.HandleFunc("/action/user/{encoded_user_name}/notifications", ah.GetNotifs).Methods("GET").Schemes("http")
	r.HandleFunc("/action/user/{encoded_user_name}/notifications/{notif_id}/read", ah.UpdateNotif).Methods("PUT").Schemes("http")
	// Post
	r.HandleFunc("/action/post", ah.AddPost).Methods("POST").Schemes("http")
	r.HandleFunc("/action/post/{post_guid}", ah.EditPost).Methods("PUT").Schemes("http")
	r.HandleFunc("/action/post/{post_guid}", ah.DeletePost).Methods("DELETE").Schemes("http")
	r.HandleFunc("/action/post/{post_guid}/up", ah.RatePostUp).Methods("POST").Schemes("http")
	r.HandleFunc("/action/post/{post_guid}/down", ah.RatePostDown).Methods("POST").Schemes("http")
	// Comment
	r.HandleFunc("/action/post/{post_guid}/comment", ah.AddComment).Methods("POST").Schemes("http")
	r.HandleFunc("/action/post/{post_guid}/comment/{comment_id}", ah.EditComment).Methods("PUT").Schemes("http")
	r.HandleFunc("/action/post/{post_guid}/comment/{comment_id}", ah.DeleteComment).Methods("DELETE").Schemes("http")
	r.HandleFunc("/action/post/{post_guid}/comment/{comment_id}/up", ah.RateCommentUp).Methods("POST").Schemes("http")
	r.HandleFunc("/action/post/{post_guid}/comment/{comment_id}/down", ah.RateCommentDown).Methods("POST").Schemes("http")
	// Password
	r.HandleFunc("/action/forgot-password", ah.ForgotPassword).Methods("POST").Schemes("http")
	r.HandleFunc("/action/recover-password", ah.RecoverPassword).Methods("POST").Schemes("http")
	r.HandleFunc("/action/change-password", ah.ChangePassword).Methods("POST").Schemes("http")
	// Image
	r.HandleFunc("/action/image", ah.GetPostImage).Schemes("http")
	r.HandleFunc("/action/profile-pic", ah.GetProfilePic).Schemes("http")

	// Websocket
	r.HandleFunc("/ws", wsh.HandleWSMsg)

	return mw.Wrap(r)

}
