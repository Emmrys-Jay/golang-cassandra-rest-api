package user

import (
	"github.com/gocql/gocql"
	tokenSessionRepository "github.com/windswept321/smartest-city-roadmap-go/module/tokensession/repository"
	"github.com/windswept321/smartest-city-roadmap-go/proposals-and-comments/proposals/controller"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

//
//
func Initialize(e *echo.Echo, db *gorm.DB, session *gocql.Session, casbinMdw echo.MiddlewareFunc, apiKeyMdw echo.MiddlewareFunc) {
	tokenSessionRepository := tokenSessionRepository.NewTokenSessionRepository(db)
	proposalController := controller.NewProposalController(tokenSessionRepository, session)

	proposal := e.Group("api/v1/user/proposal")
	proposal.POST("/create", proposalController.WriteProposal, casbinMdw)
	proposal.GET("/getAll", proposalController.GetAllProposals, apiKeyMdw)
	proposal.GET("/get/:id", proposalController.GetProposalByProposalID, apiKeyMdw)
	proposal.GET("/get/time", proposalController.GetProposalByTimeCreated, apiKeyMdw)
	proposal.GET("/get/user-id/:id", proposalController.GetProposalsByUserID, apiKeyMdw)
	proposal.PUT("/update", proposalController.UpdateProposal, casbinMdw)
	proposal.DELETE("/delete/:id", proposalController.DeleteProposal, casbinMdw)
	proposal.DELETE("/deleteAll", proposalController.DeleteAllProposals, casbinMdw)
	proposal.PUT("/upvote/:id", proposalController.UpvoteProposal, casbinMdw)
	proposal.PUT("/downvote/:id", proposalController.DownvoteProposal, casbinMdw)
}
