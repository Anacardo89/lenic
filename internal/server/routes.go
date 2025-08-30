package server

import (
	"net/http"

	"github.com/Anacardo89/lenic/internal/middleware"
	"github.com/Anacardo89/lenic/internal/server/httphandle/api"
	"github.com/Anacardo89/lenic/internal/server/httphandle/page"
	"github.com/Anacardo89/lenic/internal/server/httphandle/redirect"
	"github.com/Anacardo89/lenic/internal/server/wshandle"
	"github.com/gorilla/mux"
)

func NewRouter(ah *api.APIHandler, ph *page.PageHandler, wsh *wshandle.WSHandler, mw *middleware.MiddlewareHandler) http.Handler {

	r := mux.NewRouter()

	r.HandleFunc("/", redirect.RedirIndex).Schemes("https")

	// Page
	r.HandleFunc("/home", ph.Home).Schemes("https")
	r.HandleFunc("/login", ph.Login).Schemes("https")
	r.HandleFunc("/register", ph.Register).Schemes("https")
	r.HandleFunc("/error", ph.Error).Schemes("https")
	r.HandleFunc("/newPost", ph.NewPost).Schemes("https")
	r.HandleFunc("/post/{post_guid}", ph.Post).Schemes("https")
	r.HandleFunc("/user/{encoded_user_name}", ph.UserProfile).Schemes("https")
	r.HandleFunc("/user/{encoded_user_name}/feed", ph.Feed).Schemes("https")
	r.HandleFunc("/user/{encoded_user_name}/followers", ph.UserFollowers).Schemes("https")
	r.HandleFunc("/user/{encoded_user_name}/following", ph.UserFollowing).Schemes("https")
	r.HandleFunc("/forgot-password", ph.ForgotPassword).Schemes("https")
	r.HandleFunc("/recover-password/{encoded_user_name}", ph.RecoverPassword).Schemes("https")
	r.HandleFunc("/change-password/{encoded_user_name}", ph.ChangePassword).Schemes("https")

	// API
	r.HandleFunc("/action/register", ah.RegisterUser).Methods("POST").Schemes("https")
	r.HandleFunc("/action/activate/{encoded_user_name}", ah.ActivateUser).Schemes("https")
	r.HandleFunc("/action/login", ah.Login).Methods("POST").Schemes("https")
	r.HandleFunc("/action/logout", ah.Logout).Methods("POST").Schemes("https")
	// User
	r.HandleFunc("/action/search/user", ah.SearchUsers).Methods("GET").Schemes("https")
	r.HandleFunc("/action/user/{encoded_user_name}/follow", ah.RequestFollowUser).Methods("POST").Schemes("https")
	r.HandleFunc("/action/user/{encoded_user_name}/accept", ah.AcceptFollowRequest).Methods("PUT").Schemes("https")
	r.HandleFunc("/action/user/{encoded_user_name}/unfollow", ah.UnfollowUser).Methods("DELETE").Schemes("https")
	r.HandleFunc("/action/user/{encoded_user_name}/profile-pic", ah.PostProfilePic).Methods("POST").Schemes("https")
	// DM
	r.HandleFunc("/action/user/{encoded_user_name}/conversations", ah.GetConversations).Methods("GET").Schemes("https")
	r.HandleFunc("/action/user/{encoded_user_name}/conversations", ah.StartConversation).Methods("POST").Schemes("https")
	r.HandleFunc("/action/user/{encoded_user_name}/conversations/{conversation_id}/read", ah.ReadConversation).Methods("PUT").Schemes("https")
	r.HandleFunc("/action/user/{encoded_user_name}/conversations/{conversation_id}/dms", ah.GetDMs).Methods("GET").Schemes("https")
	r.HandleFunc("/action/user/{encoded_user_name}/conversations/{conversation_id}/dms", ah.SendDM).Methods("POST").Schemes("https")
	// Notif
	r.HandleFunc("/action/user/{encoded_user_name}/notifications", ah.GetNotifs).Methods("GET").Schemes("https")
	r.HandleFunc("/action/user/{encoded_user_name}/notifications/{notif_id}/read", ah.UpdateNotif).Methods("PUT").Schemes("https")
	// Post
	r.HandleFunc("/action/post", ah.AddPost).Methods("POST").Schemes("https")
	r.HandleFunc("/action/post/{post_guid}", ah.EditPost).Methods("PUT").Schemes("https")
	r.HandleFunc("/action/post/{post_guid}", ah.DeletePost).Methods("DELETE").Schemes("https")
	r.HandleFunc("/action/post/{post_guid}/up", ah.RatePostUp).Methods("POST").Schemes("https")
	r.HandleFunc("/action/post/{post_guid}/down", ah.RatePostDown).Methods("POST").Schemes("https")
	// Comment
	r.HandleFunc("/action/post/{post_guid}/comment", ah.AddComment).Methods("POST").Schemes("https")
	r.HandleFunc("/action/post/{post_guid}/comment/{comment_id}", ah.EditComment).Methods("PUT").Schemes("https")
	r.HandleFunc("/action/post/{post_guid}/comment/{comment_id}", ah.DeleteComment).Methods("DELETE").Schemes("https")
	r.HandleFunc("/action/post/{post_guid}/comment/{comment_id}/up", ah.RateCommentUp).Methods("POST").Schemes("https")
	r.HandleFunc("/action/post/{post_guid}/comment/{comment_id}/down", ah.RateCommentDown).Methods("POST").Schemes("https")
	// Password
	r.HandleFunc("/action/forgot-password", ah.ForgotPassword).Methods("POST").Schemes("https")
	r.HandleFunc("/action/recover-password", ah.RecoverPassword).Methods("POST").Schemes("https")
	r.HandleFunc("/action/change-password", ah.ChangePassword).Methods("POST").Schemes("https")
	// Image
	r.HandleFunc("/action/image", ah.PostImage).Schemes("https")
	r.HandleFunc("/action/profile-pic", ah.ProfilePic).Schemes("https")

	// Websocket
	r.HandleFunc("/ws", wsh.HandleWSMsg)

	return mw.Wrap(r)

}
