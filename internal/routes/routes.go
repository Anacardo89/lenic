package routes

import (
	"github.com/Anacardo89/tpsi25_blog/internal/handlers/actions"
	"github.com/Anacardo89/tpsi25_blog/internal/handlers/pages"
	"github.com/gorilla/mux"
)

func DeclareRoutes(r *mux.Router) {
	r.HandleFunc("/", pages.RedirIndex).Schemes("https")
	r.HandleFunc("/home", pages.Index).Schemes("https")
	r.HandleFunc("/login", pages.Login).Schemes("https")
	r.HandleFunc("/register", pages.Register).Schemes("https")
	r.HandleFunc("/activate/{encoded_user_name}", pages.ActivateUser).Schemes("https")
	r.HandleFunc("/error", pages.Error).Schemes("https")
	r.HandleFunc("/newPost", pages.NewPost).Schemes("https")
	r.HandleFunc("/post/{post_guid}", pages.Post).Schemes("https")
	r.HandleFunc("/forgot-password", pages.ForgotPassword).Schemes("https")

	r.HandleFunc("/action/register", actions.RegisterUser).Methods("POST").Schemes("https")
	r.HandleFunc("/action/login", actions.Login).Methods("POST").Schemes("https")
	r.HandleFunc("/action/logout", actions.Logout).Methods("POST").Schemes("https")
	r.HandleFunc("/action/post", actions.AddPost).Methods("POST").Schemes("https")
	r.HandleFunc("/action/post/{post_guid}/comment", actions.AddComment).Methods("POST").Schemes("https")
	r.HandleFunc("/action/post/{post_guid}/comment/{comment_id}", actions.EditComment).Methods("PUT").Schemes("https")
	r.HandleFunc("/action/forgot-password", actions.ForgotPassword).Methods("POST").Schemes("https")
	r.HandleFunc("/action/image", actions.Image).Schemes("https")
}
