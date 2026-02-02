package route

import (
	"time"

	"github.com/Pro100-Almaz/trading-chat/api/controller"
	"github.com/Pro100-Almaz/trading-chat/bootstrap"
	"github.com/Pro100-Almaz/trading-chat/repository"
	"github.com/Pro100-Almaz/trading-chat/usecase"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

func NewPostRouter(env *bootstrap.Env, timeout time.Duration, db *sqlx.DB, r *mux.Router) {
	postRepo := repository.NewPostRepository(db)
	userRepo := repository.NewUserRepository(db)
	likeRepo := repository.NewLikeRepository(db)
	commentRepo := repository.NewCommentRepository(db)

	postUseCase := usecase.NewPostUseCase(postRepo, userRepo, likeRepo, commentRepo, timeout)
	likeUseCase := usecase.NewLikeUseCase(likeRepo, postRepo, timeout)
	commentUseCase := usecase.NewCommentUseCase(commentRepo, postRepo, userRepo, timeout)

	postController := &controller.PostController{
		PostUseCase: postUseCase,
		Env:         env,
	}

	likeController := &controller.LikeController{
		LikeUseCase: likeUseCase,
		Env:         env,
	}

	commentController := &controller.CommentController{
		CommentUseCase: commentUseCase,
		Env:            env,
	}

	// Posts routes
	postsGroup := r.PathPrefix("/posts").Subrouter()
	postsGroup.HandleFunc("", postController.GetGlobalFeed).Methods("GET")
	postsGroup.HandleFunc("", postController.CreatePost).Methods("POST")
	postsGroup.HandleFunc("/following", postController.GetFollowingFeed).Methods("GET")
	postsGroup.HandleFunc("/user/{id}", postController.GetUserPosts).Methods("GET")
	postsGroup.HandleFunc("/{id}", postController.GetPost).Methods("GET")
	postsGroup.HandleFunc("/{id}", postController.DeletePost).Methods("DELETE")

	// Likes routes
	postsGroup.HandleFunc("/{id}/like", likeController.LikePost).Methods("POST")
	postsGroup.HandleFunc("/{id}/like", likeController.UnlikePost).Methods("DELETE")

	// Comments routes
	postsGroup.HandleFunc("/{id}/comments", commentController.GetComments).Methods("GET")
	postsGroup.HandleFunc("/{id}/comments", commentController.CreateComment).Methods("POST")

	// Delete comment route (under /comments prefix)
	commentsGroup := r.PathPrefix("/comments").Subrouter()
	commentsGroup.HandleFunc("/{id}", commentController.DeleteComment).Methods("DELETE")
}
