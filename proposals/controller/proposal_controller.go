package controller

import (
	"strings"
	"time"

	"github.com/gocql/gocql"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/windswept321/smartest-city-roadmap-go/infrastructure/controller"
	"github.com/windswept321/smartest-city-roadmap-go/infrastructure/response"
	TokenSessionsRepository "github.com/windswept321/smartest-city-roadmap-go/module/tokensession/repository"
	commentsRepository "github.com/windswept321/smartest-city-roadmap-go/proposals-and-comments/comments/repository"
	"github.com/windswept321/smartest-city-roadmap-go/proposals-and-comments/proposals/repository"
)

type ProposalController struct {
	TokenSessionRepository TokenSessionsRepository.TokenSessionRepository
	controller.BaseController
	*gocql.Session
}

func NewProposalController(tokenSessionRepository TokenSessionsRepository.TokenSessionRepository, session *gocql.Session) *ProposalController {
	return &ProposalController{
		TokenSessionRepository: tokenSessionRepository,
		Session:                session,
	}
}

type WriteProposalRequest struct {
	Title        string `json:"title" form:"title"`
	ProposalText string `json:"proposal_text" form:"proposal_text"`
}

// WriteProposal
// @Summary Create a new proposal
// @Tags proposal
// @Accept json
// @Produce json
// @Description API create new proposal
// @Param write_proposal_request body WriteProposalRequest true "req with title and proposal"
// @Success 200 {object} response.Response{Data=string}
// @Failure 400 {object} response.Response{Data=response.ErrorResponse}
// @Failure 500 {object} response.Response{Data=response.ErrorResponse}
// @Router /proposal/create [post]
// @Security JWTToken
func (p *ProposalController) WriteProposal(c echo.Context) error {
	var req WriteProposalRequest

	if err := c.Bind(&req); err != nil {
		resp := response.ErrorResponse{
			ErrorCode: 400,
			Message:   "Please check your request again for errors",
		}
		message := "false"
		return p.WriteBadRequest(c, message, resp)
	}

	if req.ProposalText == "" || req.Title == "" {
		resp := response.ErrorResponse{
			ErrorCode: 400,
			Message:   "No proposal or title found, A proposal is required",
		}
		message := "false"
		return p.WriteBadRequest(c, message, resp)
	}

	token := c.Request().Header.Get("Authorization")
	tokenSession, err := p.TokenSessionRepository.GetOneFlexible("token", token)
	if err != nil {
		resp := response.ErrorResponse{
			ErrorCode: 500,
			Message:   "Something went wrong.",
		}
		message := "false"
		return p.WriteInternalServerError(c, message, resp, "")
	}

	err = repository.StoreProposal(p.Session, req.Title, req.ProposalText, tokenSession.UserID, tokenSession.User.Username, tokenSession.User.FirstName, tokenSession.User.LastName)
	if err != nil {
		resp := response.ErrorResponse{
			ErrorCode: 500,
			Message:   "Something went wrong",
		}
		message := "false"
		return p.WriteInternalServerError(c, message, resp, "")
	}

	return p.WriteSuccess(c, "inserted proposal")
}

// GetAllProposals
// @Summary Get all proposals
// @Description API Get all proposals ordered by time
// @Tags proposal
// @Produce json
// @Success 200 {object} response.Response{Data=entity.Proposal}
// @Failure 400 {object} response.Response{Data=response.ErrorResponse}
// @Failure 500 {object} response.Response{Data=response.ErrorResponse}
// @Router /proposal/getAll [get]
// @Security JWTToken
// @Security APIKey
func (p *ProposalController) GetAllProposals(c echo.Context) error {
	proposals, err := repository.GetAllProposals(p.Session)
	if err != nil {
		resp := response.ErrorResponse{
			ErrorCode: 500,
			Message:   "Something went wrong",
		}
		message := "false"
		return p.WriteInternalServerError(c, message, resp, "")
	}

	return p.WriteSuccess(c, proposals)
}

// GetProposalsByUserID
// @Summary Get proposals
// @Description Get all proposals by a single user
// @Tags proposal
// @Accept plain
// @Produce json
// @Param user_id path string true "path string with id"
// @Success 200 {object} response.Response{Data=[]entity.Proposal}
// @Failure 400 {object} response.Response{Data=response.ErrorResponse}
// @Failure 500 {object} response.Response{Data=response.ErrorResponse}
// @Router /proposal/get/user-id/:id [get]
// @Security JWTToken
// @Security APIKey
func (p *ProposalController) GetProposalsByUserID(c echo.Context) error {
	userIDString := c.Param("id")
	if userIDString == "" {
		resp := response.ErrorResponse{
			ErrorCode: 400,
			Message:   "We cant find your ID, please check the URL again",
		}
		message := "false"
		return p.WriteBadRequest(c, message, resp)
	}

	userID, err := uuid.Parse(userIDString)
	if err != nil {
		resp := response.ErrorResponse{
			ErrorCode: 400,
			Message:   "Please check your request again for errors",
		}
		message := "false"
		return p.WriteBadRequest(c, message, resp)
	}

	proposals, err := repository.GetProposalsByUserID(p.Session, userID)
	if err != nil {
		resp := response.ErrorResponse{
			ErrorCode: 500,
			Message:   "Something went wrong",
		}
		message := "false"
		return p.WriteInternalServerError(c, message, resp, "")
	}

	return p.WriteSuccess(c, proposals)
}

// GetProposalsByTimeCreated
// @Summary Get proposals
// @Description Get all proposals created within a duration
// @Tags proposal
// @Accept plain
// @Produce json
// @Param date-from query string true "format: 2022-06-23-14:00"
// @Param date-to query string true "format: 2022-06-23-14:00"
// @Success 200 {object} response.Response{Data=[]entity.Proposal}
// @Failure 400 {object} response.Response{Data=response.ErrorResponse}
// @Failure 500 {object} response.Response{Data=response.ErrorResponse}
// @Router /proposal/get/time [get]
// @Security JWTToken
// @Security APIKey
func (p *ProposalController) GetProposalByTimeCreated(c echo.Context) error {
	dateFromString := c.QueryParam("date-from")
	dateToString := c.QueryParam("date-to")

	dateFromSlice := strings.Split(dateFromString, "-")
	dateToSlice := strings.Split(dateToString, "-")

	if len(dateFromSlice) != 4 || len(dateToSlice) != 4 {
		resp := response.ErrorResponse{
			ErrorCode: 400,
			Message:   "Time format error",
		}
		message := "false"
		return p.WriteBadRequest(c, message, resp)

	}

	dateFromString = dateFromSlice[0] + "-" + dateFromSlice[1] + "-" + dateFromSlice[2] + "T" + dateFromSlice[3] + ":00Z"
	dateToString = dateToSlice[0] + "-" + dateToSlice[1] + "-" + dateToSlice[2] + "T" + dateToSlice[3] + ":00Z"

	const timeForm = time.RFC3339

	dateFrom, err := time.Parse(timeForm, dateFromString)
	if err != nil {
		resp := response.ErrorResponse{
			ErrorCode: 400,
			Message:   "Time format error",
		}
		message := "false"
		return p.WriteBadRequest(c, message, resp)
	}

	dateTo, err := time.Parse(timeForm, dateToString)
	if err != nil {
		resp := response.ErrorResponse{
			ErrorCode: 400,
			Message:   "Time format error",
		}
		message := "false"
		return p.WriteBadRequest(c, message, resp)
	}

	proposals, err := repository.GetProposalsByTimeCreated(p.Session, dateFrom, dateTo)
	if err != nil {
		resp := response.ErrorResponse{
			ErrorCode: 500,
			Message:   err.Error(),
		}
		message := "false"
		return p.WriteInternalServerError(c, message, resp, "")
	}

	return p.WriteSuccess(c, proposals)
}

// GetProposalsByProposalID
// @Summary Get a single proposal
// @Description Get a single proposal by its unique id
// @Tags proposal
// @Accept plain
// @Produce json
// @Param proposal_id path string true "unique proposal id"
// @Success 200 {object} response.Response{Data=[]entity.Proposal}
// @Failure 400 {object} response.Response{Data=response.ErrorResponse}
// @Failure 500 {object} response.Response{Data=response.ErrorResponse}
// @Router /proposal/get/:id [get]
// @Security JWTToken
// @Security APIKey
func (p *ProposalController) GetProposalByProposalID(c echo.Context) error {
	proposalIDString := c.Param("id")
	if proposalIDString == "" {
		resp := response.ErrorResponse{
			ErrorCode: 400,
			Message:   "We cant find your ID, please check the URL again",
		}
		message := "false"
		return p.WriteBadRequest(c, message, resp)
	}

	proposalID, err := uuid.Parse(proposalIDString)
	if err != nil {
		resp := response.ErrorResponse{
			ErrorCode: 400,
			Message:   "Please check your request again for errors",
		}
		message := "false"
		return p.WriteBadRequest(c, message, resp)
	}

	proposal, err := repository.GetProposalByProposalID(p.Session, proposalID)
	if err != nil {
		resp := response.ErrorResponse{
			ErrorCode: 500,
			Message:   "Something went wrong",
		}
		message := "false"
		return p.WriteInternalServerError(c, message, resp, "")
	}

	return p.WriteSuccess(c, proposal)
}

type UpdateProposalRequest struct {
	ID           string `json:"id" form:"id"`
	Title        string `json:"title" form:"title"`
	ProposalText string `json:"proposal" form:"proposal"`
}

// UpdateProposal
// @Summary Update an existing proposal
// @Description Update an existing proposal
// @Tags proposal
// @Accept json
// @Produce json
// @Param update_proposal_request body UpdateProposalRequest true "json req with ID, updated title, and text "
// @Success 200 {object} response.Response{Data=string}
// @Failure 400 {object} response.Response{Data=response.ErrorResponse}
// @Failure 500 {object} response.Response{Data=response.ErrorResponse}
// @Router /proposal/update [put]
// @Security JWTToken
func (p *ProposalController) UpdateProposal(c echo.Context) error {
	var req UpdateProposalRequest

	if err := c.Bind(&req); err != nil {
		resp := response.ErrorResponse{
			ErrorCode: 400,
			Message:   "Please check your request again for errors",
		}
		message := "false"
		return p.WriteBadRequest(c, message, resp)
	}

	proposalID, err := uuid.Parse(req.ID)
	if err != nil {
		resp := response.ErrorResponse{
			ErrorCode: 400,
			Message:   "Please check your UUID string for errors",
		}
		message := "false"
		return p.WriteBadRequest(c, message, resp)
	}

	if req.ProposalText == "" || req.Title == "" {
		resp := response.ErrorResponse{
			ErrorCode: 400,
			Message:   "No proposal title or text found",
		}
		message := "false"
		return p.WriteBadRequest(c, message, resp)
	}

	err = repository.UpdateProposal(p.Session, proposalID, req.Title, req.ProposalText)
	if err != nil {
		resp := response.ErrorResponse{
			ErrorCode: 500,
			Message:   "Something went wrong",
		}
		message := "false"
		return p.WriteInternalServerError(c, message, resp, "")
	}

	return p.WriteSuccess(c, "updated")
}

// DeleteProposal
// @Summary Delete a single proposal
// @Description Delete a proposal using its unique id
// @Tags proposal
// @Accept plain
// @Produce json
// @Param proposal_id path string true "unique proposal id"
// @Success 200 {object} response.Response{Data=string}
// @Failure 400 {object} response.Response{Data=response.ErrorResponse}
// @Failure 500 {object} response.Response{Data=response.ErrorResponse}
// @Router /proposal/delete/:id [delete]
// @Security JWTToken
func (p *ProposalController) DeleteProposal(c echo.Context) error {
	proposalIDString := c.Param("id")
	proposalID, err := uuid.Parse(proposalIDString)
	if err != nil {
		resp := response.ErrorResponse{
			ErrorCode: 400,
			Message:   "Please check your request again for errors",
		}
		message := "false"
		return p.WriteBadRequest(c, message, resp)
	}

	err = repository.DeleteProposal(p.Session, proposalID)
	if err != nil {
		resp := response.ErrorResponse{
			ErrorCode: 500,
			Message:   "Something went wrong",
		}
		message := "false"
		return p.WriteInternalServerError(c, message, resp, "")
	}

	return p.WriteSuccess(c, "deleted")
}

// DeleteAllProposals
// @Summary Delete all proposals
// @Description Delete all proposal - for only admin
// @Tags proposal
// @Produce json
// @Success 200 {object} response.Response{Data=string}
// @Failure 400 {object} response.Response{Data=response.ErrorResponse}
// @Failure 500 {object} response.Response{Data=response.ErrorResponse}
// @Router /proposal/deleteAll [delete]
// @Security JWTToken
func (p *ProposalController) DeleteAllProposals(c echo.Context) error {
	err := repository.DeleteAllProposals(p.Session)
	if err != nil && err != gocql.ErrTimeoutNoResponse {
		resp := response.ErrorResponse{
			ErrorCode: 500,
			Message:   err.Error(),
		}
		message := "false"
		return p.WriteInternalServerError(c, message, resp, "")
	}

	err = commentsRepository.DeleteAllComments(p.Session)
	if err != nil && err != gocql.ErrTimeoutNoResponse {
		resp := response.ErrorResponse{
			ErrorCode: 500,
			Message:   err.Error(),
		}
		message := "false"
		return p.WriteInternalServerError(c, message, resp, "")
	}

	return p.WriteSuccess(c, "deleted all proposals")
}

// UpvoteProposal
// @Summary Upvote a single proposal
// @Description Upvote a proposal using its unique id
// @Tags proposal
// @Accept plain
// @Produce json
// @Param proposal_id path string true "upvote propoosal by unique id"
// @Success 200 {object} response.Response{Data=string}
// @Failure 400 {object} response.Response{Data=response.ErrorResponse}
// @Failure 500 {object} response.Response{Data=response.ErrorResponse}
// @Router /proposal/upvote/:id [put]
// @Security JWTToken
func (p *ProposalController) UpvoteProposal(c echo.Context) error {
	proposalIDString := c.Param("id")
	proposalID, err := uuid.Parse(proposalIDString)
	if err != nil {
		resp := response.ErrorResponse{
			ErrorCode: 400,
			Message:   "Please check your request again for errors",
		}
		message := "false"
		return p.WriteBadRequest(c, message, resp)
	}

	err = repository.UpvoteProposal(p.Session, proposalID)
	if err != nil {
		resp := response.ErrorResponse{
			ErrorCode: 500,
			Message:   err.Error(),
		}
		message := "false"
		return p.WriteInternalServerError(c, message, resp, "")
	}

	return p.WriteSuccess(c, "upvoted")
}

// DownvoteProposal
// @Summary Downvote a single proposal
// @Description Downvote a proposal using its unique id
// @Tags proposal
// @Accept plain
// @Produce json
// @Param proposal_id path string true "downvote propoosal by unique id"
// @Success 200 {object} response.Response{Data=string}
// @Failure 400 {object} response.Response{Data=response.ErrorResponse}
// @Failure 500 {object} response.Response{Data=response.ErrorResponse}
// @Router /proposal/downvote/:id [put]
// @Security JWTToken
func (p *ProposalController) DownvoteProposal(c echo.Context) error {
	proposalIDString := c.Param("id")
	proposalID, err := uuid.Parse(proposalIDString)
	if err != nil {
		resp := response.ErrorResponse{
			ErrorCode: 400,
			Message:   "Please check your request again for errors",
		}
		message := "false"
		return p.WriteBadRequest(c, message, resp)
	}

	err = repository.DownvoteProposal(p.Session, proposalID)
	if err != nil {
		resp := response.ErrorResponse{
			ErrorCode: 500,
			Message:   "Something went wrong",
		}
		message := "false"
		return p.WriteInternalServerError(c, message, resp, "")
	}

	return p.WriteSuccess(c, "downvoted")
}
