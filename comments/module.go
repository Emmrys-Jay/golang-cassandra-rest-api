package user

import (
	"github.com/gocql/gocql"
	"github.com/labstack/echo/v4"
	tokenSessionsRepository "github.com/windswept321/smartest-city-roadmap-go/module/tokensession/repository"
	"github.com/windswept321/smartest-city-roadmap-go/proposals-and-comments/comments/controller"
	proposalController "github.com/windswept321/smartest-city-roadmap-go/proposals-and-comments/proposals/controller"
	"gorm.io/gorm"
)

func Initialize(e *echo.Echo, db *gorm.DB, session *gocql.Session, casbinMdw echo.MiddlewareFunc, apiKeyMdw echo.MiddlewareFunc) {
	tokenSessionRepository := tokenSessionsRepository.NewTokenSessionRepository(db)
	proposalController := proposalController.NewProposalController(tokenSessionRepository, session)
	commentsController := controller.NewCommentsController(proposalController)

	comment := e.Group("api/v1/user/proposal/comment")
	comment.POST("/create", commentsController.WriteComment, casbinMdw)
	comment.GET("/getAll/:proposal-id", commentsController.GetCommentsByProposalID, apiKeyMdw)
	comment.GET("/get", commentsController.GetCommentByIDAndProposalID, apiKeyMdw)
	comment.PUT("/update", commentsController.UpdateComment, casbinMdw)
	comment.DELETE("/delete", commentsController.DeleteComment, casbinMdw)
	comment.DELETE("/delete/:proposal-id", commentsController.DeleteAllProposalComments, casbinMdw)
	comment.PUT("/upvote", commentsController.UpvoteComment, casbinMdw)
}
