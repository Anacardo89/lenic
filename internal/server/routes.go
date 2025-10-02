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

	// Static
	staticDir := "./static"
	r.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))),
	)

	// Page
	r.HandleFunc("/home", ph.Home).Schemes("http")
	r.HandleFunc("/login", ph.Login).Schemes("http")
	r.HandleFunc("/register", ph.Register).Schemes("http")
	r.HandleFunc("/error", ph.Error).Schemes("http")
	r.HandleFunc("/newPost", ph.NewPost).Schemes("http")
	r.HandleFunc("/post/{post_id}", ph.Post).Schemes("http")
	r.HandleFunc("/user/{encoded_username}", ph.UserProfile).Schemes("http")
	r.HandleFunc("/user/{encoded_username}/feed", ph.Feed).Schemes("http")
	r.HandleFunc("/user/{encoded_username}/followers", ph.UserFollowers).Schemes("http")
	r.HandleFunc("/user/{encoded_username}/following", ph.UserFollowing).Schemes("http")
	r.HandleFunc("/forgot-password", ph.ForgotPassword).Schemes("http")
	r.HandleFunc("/recover-password/{encoded_username}", ph.RecoverPassword).Schemes("http")
	r.HandleFunc("/change-password/{encoded_username}", ph.ChangePassword).Schemes("http")

	// API
	r.HandleFunc("/action/register", ah.RegisterUser).Methods("POST").Schemes("http")
	r.HandleFunc("/action/activate/{encoded_username}", ah.ActivateUser).Schemes("http")
	r.HandleFunc("/action/login", ah.Login).Methods("POST").Schemes("http")
	r.HandleFunc("/action/logout", ah.Logout).Methods("POST").Schemes("http")
	// User
	r.HandleFunc("/action/search/user", ah.SearchUsers).Methods("GET").Schemes("http")
	r.HandleFunc("/action/user/{encoded_username}/follow", ah.RequestFollowUser).Methods("POST").Schemes("http")
	r.HandleFunc("/action/user/{encoded_username}/accept", ah.AcceptFollowRequest).Methods("PUT").Schemes("http")
	r.HandleFunc("/action/user/{encoded_username}/unfollow", ah.UnfollowUser).Methods("DELETE").Schemes("http")
	r.HandleFunc("/action/user/{encoded_username}/profile-pic", ah.PostProfilePic).Methods("POST").Schemes("http")
	// DM
	r.HandleFunc("/action/user/{encoded_username}/conversations", ah.GetConversations).Methods("GET").Schemes("http")
	r.HandleFunc("/action/user/{encoded_username}/conversations", ah.StartConversation).Methods("POST").Schemes("http")
	r.HandleFunc("/action/user/{encoded_username}/conversations/{conversation_id}/read", ah.ReadConversation).Methods("PUT").Schemes("http")
	r.HandleFunc("/action/user/{encoded_username}/conversations/{conversation_id}/dms", ah.GetDMs).Methods("GET").Schemes("http")
	r.HandleFunc("/action/user/{encoded_username}/conversations/{conversation_id}/dms", ah.SendDM).Methods("POST").Schemes("http")
	// Notif
	r.HandleFunc("/action/user/{encoded_username}/notifications", ah.GetNotifs).Methods("GET").Schemes("http")
	r.HandleFunc("/action/user/{encoded_username}/notifications/{notif_id}/read", ah.UpdateNotif).Methods("PUT").Schemes("http")
	// Post
	r.HandleFunc("/action/post", ah.AddPost).Methods("POST").Schemes("http")
	r.HandleFunc("/action/post/{post_id}", ah.EditPost).Methods("PUT").Schemes("http")
	r.HandleFunc("/action/post/{post_id}", ah.DeletePost).Methods("DELETE").Schemes("http")
	r.HandleFunc("/action/post/{post_id}/up", ah.RatePostUp).Methods("POST").Schemes("http")
	r.HandleFunc("/action/post/{post_id}/down", ah.RatePostDown).Methods("POST").Schemes("http")
	// Comment
	r.HandleFunc("/action/post/{post_id}/comment", ah.AddComment).Methods("POST").Schemes("http")
	r.HandleFunc("/action/post/{post_id}/comment/{comment_id}", ah.EditComment).Methods("PUT").Schemes("http")
	r.HandleFunc("/action/post/{post_id}/comment/{comment_id}", ah.DeleteComment).Methods("DELETE").Schemes("http")
	r.HandleFunc("/action/post/{post_id}/comment/{comment_id}/up", ah.RateCommentUp).Methods("POST").Schemes("http")
	r.HandleFunc("/action/post/{post_id}/comment/{comment_id}/down", ah.RateCommentDown).Methods("POST").Schemes("http")
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
