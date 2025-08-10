package server

import (
	"github.com/Anacardo89/tpsi25_blog/internal/handlers/wsoc"
	"github.com/gorilla/mux"
)

func (s *Server) DeclareRoutes(r *mux.Router) {
	r.HandleFunc("/", s.RedirIndex).Schemes("https")
	r.HandleFunc("/home", s.Home).Schemes("https")
	r.HandleFunc("/login", s.PageLogin).Schemes("https")
	r.HandleFunc("/register", s.PageRegister).Schemes("https")
	r.HandleFunc("/error", s.Error).Schemes("https")
	r.HandleFunc("/newPost", s.NewPost).Schemes("https")
	r.HandleFunc("/post/{post_guid}", s.Post).Schemes("https")
	r.HandleFunc("/user/{encoded_user_name}", s.UserProfile).Schemes("https")
	r.HandleFunc("/user/{encoded_user_name}/feed", s.Feed).Schemes("https")
	r.HandleFunc("/user/{encoded_user_name}/followers", s.UserFollowers).Schemes("https")
	r.HandleFunc("/user/{encoded_user_name}/following", s.UserFollowing).Schemes("https")
	r.HandleFunc("/forgot-password", s.PageForgotPassword).Schemes("https")
	r.HandleFunc("/recover-password/{encoded_user_name}", s.PageRecoverPassword).Schemes("https")
	r.HandleFunc("/change-password/{encoded_user_name}", s.PageChangePassword).Schemes("https")

	r.HandleFunc("/action/register", s.RegisterUser).Methods("POST").Schemes("https")
	r.HandleFunc("/action/activate/{encoded_user_name}", s.ActivateUser).Schemes("https")
	r.HandleFunc("/action/login", s.Login).Methods("POST").Schemes("https")
	r.HandleFunc("/action/logout", s.Logout).Methods("POST").Schemes("https")
	r.HandleFunc("/action/search/user", s.SearchUsers).Methods("GET").Schemes("https")
	r.HandleFunc("/action/user/{encoded_user_name}/follow", s.RequestFollowUser).Methods("POST").Schemes("https")
	r.HandleFunc("/action/user/{encoded_user_name}/accept", s.AcceptFollowRequest).Methods("PUT").Schemes("https")
	r.HandleFunc("/action/user/{encoded_user_name}/unfollow", s.UnfollowUser).Methods("DELETE").Schemes("https")
	r.HandleFunc("/action/user/{encoded_user_name}/profile-pic", s.PostProfilePic).Methods("POST").Schemes("https")
	r.HandleFunc("/action/user/{encoded_user_name}/conversations", s.GetConversations).Methods("GET").Schemes("https")
	r.HandleFunc("/action/user/{encoded_user_name}/conversations", s.StartConversation).Methods("POST").Schemes("https")
	r.HandleFunc("/action/user/{encoded_user_name}/conversations/{conversation_id}/read", s.ReadConversation).Methods("PUT").Schemes("https")
	r.HandleFunc("/action/user/{encoded_user_name}/conversations/{conversation_id}/dms", s.GetDMs).Methods("GET").Schemes("https")
	r.HandleFunc("/action/user/{encoded_user_name}/conversations/{conversation_id}/dms", s.SendDM).Methods("POST").Schemes("https")
	r.HandleFunc("/action/user/{encoded_user_name}/notifications", s.GetNotifs).Methods("GET").Schemes("https")
	r.HandleFunc("/action/user/{encoded_user_name}/notifications/{notif_id}/read", s.UpdateNotif).Methods("PUT").Schemes("https")
	r.HandleFunc("/action/post", s.AddPost).Methods("POST").Schemes("https")
	r.HandleFunc("/action/post/{post_guid}", s.EditPost).Methods("PUT").Schemes("https")
	r.HandleFunc("/action/post/{post_guid}", s.DeletePost).Methods("DELETE").Schemes("https")
	r.HandleFunc("/action/post/{post_guid}/up", s.RatePostUp).Methods("POST").Schemes("https")
	r.HandleFunc("/action/post/{post_guid}/down", s.RatePostDown).Methods("POST").Schemes("https")
	r.HandleFunc("/action/post/{post_guid}/comment", s.AddComment).Methods("POST").Schemes("https")
	r.HandleFunc("/action/post/{post_guid}/comment/{comment_id}", s.EditComment).Methods("PUT").Schemes("https")
	r.HandleFunc("/action/post/{post_guid}/comment/{comment_id}", s.DeleteComment).Methods("DELETE").Schemes("https")
	r.HandleFunc("/action/post/{post_guid}/comment/{comment_id}/up", s.RateCommentUp).Methods("POST").Schemes("https")
	r.HandleFunc("/action/post/{post_guid}/comment/{comment_id}/down", s.RateCommentDown).Methods("POST").Schemes("https")
	r.HandleFunc("/action/forgot-password", s.ForgotPassword).Methods("POST").Schemes("https")
	r.HandleFunc("/action/recover-password", s.RecoverPassword).Methods("POST").Schemes("https")
	r.HandleFunc("/action/change-password", s.ChangePassword).Methods("POST").Schemes("https")

	r.HandleFunc("/ws", wsoc.HandleWebSocket)

	r.HandleFunc("/action/image", s.PostImage).Schemes("https")
	r.HandleFunc("/action/profile-pic", s.ProfilePic).Schemes("https")
}
