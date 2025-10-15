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

	r.Handle("/", http.HandlerFunc(redirect.RedirIndex))

	// Page
	r.HandleFunc("/home", ph.Home).Schemes("http")
	r.HandleFunc("/login", ph.Login).Schemes("http")
	r.HandleFunc("/register", ph.Register).Schemes("http")
	r.HandleFunc("/error", ph.Error).Schemes("http")
	r.Handle("/newPost", mw.Auth(http.HandlerFunc(ph.NewPost))).Schemes("http")
	r.Handle("/post/{post_id}", mw.Auth(http.HandlerFunc(ph.Post))).Schemes("http")
	r.Handle("/user/{encoded_username}", mw.Auth(http.HandlerFunc(ph.UserProfile))).Schemes("http")
	r.Handle("/user/{encoded_username}/feed", mw.Auth(http.HandlerFunc(ph.Feed))).Schemes("http")
	r.Handle("/user/{encoded_username}/followers", mw.Auth(http.HandlerFunc(ph.UserFollowers))).Schemes("http")
	r.Handle("/user/{encoded_username}/following", mw.Auth(http.HandlerFunc(ph.UserFollowing))).Schemes("http")
	r.Handle("/change-password/{encoded_username}", mw.Auth(http.HandlerFunc(ph.ChangePassword))).Schemes("http")
	r.Handle("/forgot-password", http.HandlerFunc(ph.ForgotPassword)).Schemes("http")
	r.Handle("/recover-password/{encoded_username}", http.HandlerFunc(ph.RecoverPassword)).Schemes("http")

	// API
	r.HandleFunc("/action/register", ah.RegisterUser).Methods("POST").Schemes("http")
	r.HandleFunc("/action/activate/{encoded_username}", ah.ActivateUser).Schemes("http")
	r.HandleFunc("/action/login", ah.Login).Methods("POST").Schemes("http")
	r.HandleFunc("/action/logout", ah.Logout).Methods("POST").Schemes("http")
	r.HandleFunc("/action/forgot-password", ah.ForgotPassword).Methods("POST").Schemes("http")
	r.HandleFunc("/action/recover-password", ah.RecoverPassword).Methods("POST").Schemes("http")

	authRoutes := r.PathPrefix("/action").Subrouter()
	authRoutes.Use(mw.Auth)
	// User
	authRoutes.HandleFunc("/search/user", ah.SearchUsers).Methods("GET").Schemes("http")
	authRoutes.HandleFunc("/user/{encoded_username}/follow", ah.RequestFollowUser).Methods("POST").Schemes("http")
	authRoutes.HandleFunc("/user/{encoded_username}/accept", ah.AcceptFollowRequest).Methods("PUT").Schemes("http")
	authRoutes.HandleFunc("/user/{encoded_username}/unfollow", ah.UnfollowUser).Methods("DELETE").Schemes("http")
	authRoutes.HandleFunc("/user/{encoded_username}/profile-pic", ah.PostProfilePic).Methods("POST").Schemes("http")
	// DM
	authRoutes.HandleFunc("/user/{encoded_username}/conversations", ah.GetConversations).Methods("GET").Schemes("http")
	authRoutes.HandleFunc("/user/{encoded_username}/conversations", ah.StartConversation).Methods("POST").Schemes("http")
	authRoutes.HandleFunc("/user/{encoded_username}/conversations/{conversation_id}/read", ah.ReadConversation).Methods("PUT").Schemes("http")
	authRoutes.HandleFunc("/user/{encoded_username}/conversations/{conversation_id}/dms", ah.GetDMs).Methods("GET").Schemes("http")
	authRoutes.HandleFunc("/user/{encoded_username}/conversations/{conversation_id}/dms", ah.SendDM).Methods("POST").Schemes("http")
	// Notif
	authRoutes.HandleFunc("/user/{encoded_username}/notifications", ah.GetNotifs).Methods("GET").Schemes("http")
	authRoutes.HandleFunc("/user/{encoded_username}/notifications/{notif_id}/read", ah.UpdateNotif).Methods("PUT").Schemes("http")
	// Post
	authRoutes.HandleFunc("/post", ah.AddPost).Methods("POST").Schemes("http")
	authRoutes.HandleFunc("/post/{post_id}", ah.EditPost).Methods("PUT").Schemes("http")
	authRoutes.HandleFunc("/post/{post_id}", ah.DeletePost).Methods("DELETE").Schemes("http")
	authRoutes.HandleFunc("/post/{post_id}/up", ah.RatePostUp).Methods("POST").Schemes("http")
	authRoutes.HandleFunc("/post/{post_id}/down", ah.RatePostDown).Methods("POST").Schemes("http")
	// Comment
	authRoutes.HandleFunc("/post/{post_id}/comment", ah.AddComment).Methods("POST").Schemes("http")
	authRoutes.HandleFunc("/post/{post_id}/comment/{comment_id}", ah.EditComment).Methods("PUT").Schemes("http")
	authRoutes.HandleFunc("/post/{post_id}/comment/{comment_id}", ah.DeleteComment).Methods("DELETE").Schemes("http")
	authRoutes.HandleFunc("/post/{post_id}/comment/{comment_id}/up", ah.RateCommentUp).Methods("POST").Schemes("http")
	authRoutes.HandleFunc("/post/{post_id}/comment/{comment_id}/down", ah.RateCommentDown).Methods("POST").Schemes("http")
	// Password

	authRoutes.HandleFunc("/change-password", ah.ChangePassword).Methods("POST").Schemes("http")
	// Image
	authRoutes.HandleFunc("/image", ah.GetPostImage).Schemes("http")
	authRoutes.HandleFunc("/profile-pic", ah.GetProfilePic).Schemes("http")

	// Websocket
	r.Handle("/ws", mw.Auth(http.HandlerFunc(wsh.HandleWSMsg)))

	// Static
	staticDir := "../frontend/static"
	r.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))),
	)

	return mw.Wrap(r)

}
