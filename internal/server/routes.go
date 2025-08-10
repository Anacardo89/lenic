package server

import (
	"github.com/gorilla/mux"
)

func (s *Server) DeclareRoutes(r *mux.Router) {

	r.HandleFunc("/", s.redirectHandler.RedirIndex).Schemes("https")

	// Page
	r.HandleFunc("/home", s.pageHandler.Home).Schemes("https")
	r.HandleFunc("/login", s.pageHandler.Login).Schemes("https")
	r.HandleFunc("/register", s.pageHandler.Register).Schemes("https")
	r.HandleFunc("/error", s.pageHandler.Error).Schemes("https")
	r.HandleFunc("/newPost", s.pageHandler.NewPost).Schemes("https")
	r.HandleFunc("/post/{post_guid}", s.pageHandler.Post).Schemes("https")
	r.HandleFunc("/user/{encoded_user_name}", s.pageHandler.UserProfile).Schemes("https")
	r.HandleFunc("/user/{encoded_user_name}/feed", s.pageHandler.Feed).Schemes("https")
	r.HandleFunc("/user/{encoded_user_name}/followers", s.pageHandler.UserFollowers).Schemes("https")
	r.HandleFunc("/user/{encoded_user_name}/following", s.pageHandler.UserFollowing).Schemes("https")
	r.HandleFunc("/forgot-password", s.pageHandler.ForgotPassword).Schemes("https")
	r.HandleFunc("/recover-password/{encoded_user_name}", s.pageHandler.RecoverPassword).Schemes("https")
	r.HandleFunc("/change-password/{encoded_user_name}", s.pageHandler.ChangePassword).Schemes("https")

	// API
	r.HandleFunc("/action/register", s.apiHandler.RegisterUser).Methods("POST").Schemes("https")
	r.HandleFunc("/action/activate/{encoded_user_name}", s.apiHandler.ActivateUser).Schemes("https")
	r.HandleFunc("/action/login", s.apiHandler.Login).Methods("POST").Schemes("https")
	r.HandleFunc("/action/logout", s.apiHandler.Logout).Methods("POST").Schemes("https")
	// User
	r.HandleFunc("/action/search/user", s.apiHandler.SearchUsers).Methods("GET").Schemes("https")
	r.HandleFunc("/action/user/{encoded_user_name}/follow", s.apiHandler.RequestFollowUser).Methods("POST").Schemes("https")
	r.HandleFunc("/action/user/{encoded_user_name}/accept", s.apiHandler.AcceptFollowRequest).Methods("PUT").Schemes("https")
	r.HandleFunc("/action/user/{encoded_user_name}/unfollow", s.apiHandler.UnfollowUser).Methods("DELETE").Schemes("https")
	r.HandleFunc("/action/user/{encoded_user_name}/profile-pic", s.apiHandler.PostProfilePic).Methods("POST").Schemes("https")
	r.HandleFunc("/action/user/{encoded_user_name}/conversations", s.apiHandler.GetConversations).Methods("GET").Schemes("https")
	r.HandleFunc("/action/user/{encoded_user_name}/conversations", s.apiHandler.StartConversation).Methods("POST").Schemes("https")
	r.HandleFunc("/action/user/{encoded_user_name}/conversations/{conversation_id}/read", s.apiHandler.ReadConversation).Methods("PUT").Schemes("https")
	r.HandleFunc("/action/user/{encoded_user_name}/conversations/{conversation_id}/dms", s.apiHandler.GetDMs).Methods("GET").Schemes("https")
	r.HandleFunc("/action/user/{encoded_user_name}/conversations/{conversation_id}/dms", s.apiHandler.SendDM).Methods("POST").Schemes("https")
	r.HandleFunc("/action/user/{encoded_user_name}/notifications", s.apiHandler.GetNotifs).Methods("GET").Schemes("https")
	r.HandleFunc("/action/user/{encoded_user_name}/notifications/{notif_id}/read", s.apiHandler.UpdateNotif).Methods("PUT").Schemes("https")
	// Post
	r.HandleFunc("/action/post", s.apiHandler.AddPost).Methods("POST").Schemes("https")
	r.HandleFunc("/action/post/{post_guid}", s.apiHandler.EditPost).Methods("PUT").Schemes("https")
	r.HandleFunc("/action/post/{post_guid}", s.apiHandler.DeletePost).Methods("DELETE").Schemes("https")
	r.HandleFunc("/action/post/{post_guid}/up", s.apiHandler.RatePostUp).Methods("POST").Schemes("https")
	r.HandleFunc("/action/post/{post_guid}/down", s.apiHandler.RatePostDown).Methods("POST").Schemes("https")
	// Comment
	r.HandleFunc("/action/post/{post_guid}/comment", s.apiHandler.AddComment).Methods("POST").Schemes("https")
	r.HandleFunc("/action/post/{post_guid}/comment/{comment_id}", s.apiHandler.EditComment).Methods("PUT").Schemes("https")
	r.HandleFunc("/action/post/{post_guid}/comment/{comment_id}", s.apiHandler.DeleteComment).Methods("DELETE").Schemes("https")
	r.HandleFunc("/action/post/{post_guid}/comment/{comment_id}/up", s.apiHandler.RateCommentUp).Methods("POST").Schemes("https")
	r.HandleFunc("/action/post/{post_guid}/comment/{comment_id}/down", s.apiHandler.RateCommentDown).Methods("POST").Schemes("https")
	// Password
	r.HandleFunc("/action/forgot-password", s.apiHandler.ForgotPassword).Methods("POST").Schemes("https")
	r.HandleFunc("/action/recover-password", s.apiHandler.RecoverPassword).Methods("POST").Schemes("https")
	r.HandleFunc("/action/change-password", s.apiHandler.ChangePassword).Methods("POST").Schemes("https")
	// Image
	r.HandleFunc("/action/image", s.apiHandler.PostImage).Schemes("https")
	r.HandleFunc("/action/profile-pic", s.apiHandler.ProfilePic).Schemes("https")

	// Websocket
	r.HandleFunc("/ws", s.websocketHanler.HandleWSMsg)

}
