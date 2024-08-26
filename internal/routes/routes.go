package routes

import (
	"github.com/Anacardo89/tpsi25_blog/internal/handlers/actions"
	"github.com/Anacardo89/tpsi25_blog/internal/handlers/pages"
	"github.com/Anacardo89/tpsi25_blog/internal/handlers/redirect"
	"github.com/gorilla/mux"
)

func DeclareRoutes(r *mux.Router) {
	r.HandleFunc("/", redirect.RedirIndex).Schemes("https")
	r.HandleFunc("/home", pages.Index).Schemes("https")
	r.HandleFunc("/login", pages.Login).Schemes("https")
	r.HandleFunc("/register", pages.Register).Schemes("https")
	r.HandleFunc("/error", pages.Error).Schemes("https")
	r.HandleFunc("/newPost", pages.NewPost).Schemes("https")
	r.HandleFunc("/post/{post_guid}", pages.Post).Schemes("https")
	r.HandleFunc("/user/{encoded_user_name}", pages.UserProfile).Schemes("https")
	r.HandleFunc("/user/{encoded_user_name}/followers", pages.UserProfile).Schemes("https")
	r.HandleFunc("/user/{encoded_user_name}/following", pages.UserProfile).Schemes("https")
	r.HandleFunc("/forgot-password", pages.ForgotPassword).Schemes("https")
	r.HandleFunc("/recover-password/{encoded_user_name}", pages.RecoverPassword).Schemes("https")

	r.HandleFunc("/action/register", actions.RegisterUser).Methods("POST").Schemes("https")
	r.HandleFunc("/action/activate/{encoded_user_name}", actions.ActivateUser).Schemes("https")
	r.HandleFunc("/action/login", actions.Login).Methods("POST").Schemes("https")
	r.HandleFunc("/action/logout", actions.Logout).Methods("POST").Schemes("https")
	r.HandleFunc("/action/user/{encoded_user_name}/follow", actions.FollowUser).Methods("POST").Schemes("https")
	r.HandleFunc("/action/user/{encoded_user_name}/unfollow", actions.UnfollowUser).Methods("POST").Schemes("https")
	r.HandleFunc("/action/post", actions.AddPost).Methods("POST").Schemes("https")
	r.HandleFunc("/action/post/{post_guid}", actions.EditPost).Methods("PUT").Schemes("https")
	r.HandleFunc("/action/post/{post_guid}", actions.DeletePost).Methods("DELETE").Schemes("https")
	r.HandleFunc("/action/post/{post_guid}/up", actions.RatePostUp).Methods("POST").Schemes("https")
	r.HandleFunc("/action/post/{post_guid}/down", actions.RatePostDown).Methods("POST").Schemes("https")
	r.HandleFunc("/action/post/{post_guid}/comment", actions.AddComment).Methods("POST").Schemes("https")
	r.HandleFunc("/action/post/{post_guid}/comment/{comment_id}", actions.EditComment).Methods("PUT").Schemes("https")
	r.HandleFunc("/action/post/{post_guid}/comment/{comment_id}", actions.DeleteComment).Methods("DELETE").Schemes("https")
	r.HandleFunc("/action/post/{post_guid}/comment/{comment_id}/up", actions.RateCommentUp).Methods("POST").Schemes("https")
	r.HandleFunc("/action/post/{post_guid}/comment/{comment_id}/down", actions.RateCommentDown).Methods("POST").Schemes("https")
	r.HandleFunc("/action/forgot-password", actions.ForgotPassword).Methods("POST").Schemes("https")
	r.HandleFunc("/action/recover-password", actions.RecoverPassword).Methods("POST").Schemes("https")
	r.HandleFunc("/action/image", actions.Image).Schemes("https")
}
