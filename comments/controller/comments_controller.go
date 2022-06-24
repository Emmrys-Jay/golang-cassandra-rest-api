package controller

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/windswept321/smartest-city-roadmap-go/infrastructure/response"
	"github.com/windswept321/smartest-city-roadmap-go/proposals-and-comments/comments/repository"
	"github.com/windswept321/smartest-city-roadmap-go/proposals-and-comments/proposals/controller"
	proposalRepository "github.com/windswept321/smartest-city-roadmap-go/proposals-and-comments/proposals/repository"
)

type CommentsController struct {
	*controller.ProposalController
}

func NewCommentsController(proposalController *controller.ProposalController) *CommentsController {
	return &CommentsController{
		ProposalController: proposalController,
	}
}

type WriteCommentRequest struct {
	ProposalID string `json:"proposal_id" form:"proposal_id"`
	Comment    string `json:"comment" form:"comment"`
}

// WriteComment
// @Summary Create a comment for a proposal
// @Description Create a comment under a proposal
// @Tags proposal comment
// @Accept json
// @Produce json
// @Param write_comment_request body WriteCommentRequest true "json request with proposal id and comment"
// @Success 200 {object} response.Response{Data=string}
// @Failure 400 {object} response.Response{Data=response.ErrorResponse}
// @Failure 500 {object} response.Response{Data=response.ErrorResponse}
// @Router /proposal/comment/create [post]
// @Security JWTToken
func (p *CommentsController) WriteComment(c echo.Context) error {
	var req WriteCommentRequest

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

	if req.ProposalID == "" || req.Comment == "" {
		resp := response.ErrorResponse{
			ErrorCode: 400,
			Message:   "Please check your request again for errors",
		}
		message := "false"
		return p.WriteBadRequest(c, message, resp)
	}

	if err := c.Bind(&req); err != nil {
		resp := response.ErrorResponse{
			ErrorCode: 400,
			Message:   "Please check your request again for errors",
		}
		message := "false"
		return p.WriteBadRequest(c, message, resp)
	}

	proposalID, err := uuid.Parse(req.ProposalID)
	if err != nil {
		resp := response.ErrorResponse{
			ErrorCode: 400,
			Message:   "Please check your request UUID for errors",
		}
		message := "false"
		return p.WriteBadRequest(c, message, resp)
	}

	err = repository.StoreComment(p.Session, proposalID, req.Comment, tokenSession.UserID, tokenSession.User.Username)
	if err != nil {
		resp := response.ErrorResponse{
			ErrorCode: 500,
			Message:   "Something went wrong",
		}
		message := "false"
		return p.WriteInternalServerError(c, message, resp, "")
	}

	err = proposalRepository.AddToNumberOfComments(p.Session, proposalID)
	if err != nil {
		resp := response.ErrorResponse{
			ErrorCode: 500,
			Message:   "Something went wrong",
		}
		message := "false"
		return p.WriteInternalServerError(c, message, resp, "")
	}

	return p.WriteSuccess(c, "commented")
}

// GetCommentsByProposalID
// @Summary Get all comments under a proposal
// @Description Gets all comments with the same proposal id
// @Tags proposal comment
// @Accept plain
// @Produce json
// @Param proposal_id path string true "get all comments by proposal id"
// @Success 200 {object} response.Response{Data=[]entity.Comment}
// @Failure 400 {object} response.Response{Data=response.ErrorResponse}
// @Failure 500 {object} response.Response{Data=response.ErrorResponse}
// @Router /proposal/comment/getAll/:proposal-id [get]
// @Security JWTToken
// @Security APIKey
func (p *CommentsController) GetCommentsByProposalID(c echo.Context) error {
	proposalIDString := c.Param("proposal-id")
	proposalID, err := uuid.Parse(proposalIDString)
	if err != nil {
		resp := response.ErrorResponse{
			ErrorCode: 400,
			Message:   "Please check your request again for errors",
		}
		message := "false"
		return p.WriteBadRequest(c, message, resp)
	}

	comments, err := repository.GetCommentsByProposalID(p.Session, proposalID)
	if err != nil {
		resp := response.ErrorResponse{
			ErrorCode: 500,
			Message:   "Something went wrong",
		}
		message := "false"
		return p.WriteInternalServerError(c, message, resp, "")
	}

	return p.WriteSuccess(c, comments)
}

// GetCommentsByIDAndProposalID
// @Summary Get a single comment
// @Description Gets a single comment with a single proposal id and the unique comment id
// @Tags proposal comment
// @Accept plain
// @Produce json
// @Param proposal_id path string true "a common proposal id "
// @Param comment_id path string true "a unique comment id"
// @Success 200 {object} response.Response{Data=[]entity.Comment}
// @Failure 400 {object} response.Response{Data=response.ErrorResponse}
// @Failure 500 {object} response.Response{Data=response.ErrorResponse}
// @Router /proposal/comment/get [get]
// @Security JWTToken
// @Security APIKey
func (p *CommentsController) GetCommentByIDAndProposalID(c echo.Context) error {
	proposalIDString := c.QueryParam("proposal-id")
	commentIDString := c.QueryParam("comment-id")

	proposalID, err := uuid.Parse(proposalIDString)
	if err != nil {
		resp := response.ErrorResponse{
			ErrorCode: 400,
			Message:   "Please check your proposal ID in request again for errors",
		}
		message := "false"
		return p.WriteBadRequest(c, message, resp)
	}

	commentID, err := uuid.Parse(commentIDString)
	if err != nil {
		resp := response.ErrorResponse{
			ErrorCode: 400,
			Message:   "Please check your comment ID in request again for errors",
		}
		message := "false"
		return p.WriteBadRequest(c, message, resp)
	}

	comment, err := repository.GetCommentByIDAndProposalID(p.Session, proposalID, commentID)
	if err != nil {
		resp := response.ErrorResponse{
			ErrorCode: 500,
			Message:   "Something went wrong",
		}
		message := "false"
		return p.WriteInternalServerError(c, message, resp, "")
	}

	return p.WriteSuccess(c, comment)
}

type UpdateCommentRequest struct {
	ProposalID     string `json:"proposal_id" form:"proposal_id"`
	CommentID      string `json:"comment_id" form:"comment_id"`
	UpdatedComment string `json:"updated_comment" form:"updated_comment"`
}

// UpdateComment
// @Summary Update a single comment
// @Description Update a single comment using proposal and comment id
// @Tags proposal comment
// @Accept json
// @Produce json
// @Param update_comment_request body UpdateCommentRequest true "a json body req with proposal id, comment id and the updated comment"
// @Success 200 {object} response.Response{Data=string}
// @Failure 400 {object} response.Response{Data=response.ErrorResponse}
// @Failure 500 {object} response.Response{Data=response.ErrorResponse}
// @Router /proposal/comment/update [put]
// @Security JWTToken
func (p *CommentsController) UpdateComment(c echo.Context) error {
	var req UpdateCommentRequest

	if err := c.Bind(&req); err != nil {
		resp := response.ErrorResponse{
			ErrorCode: 400,
			Message:   "Please check your request again for errors",
		}
		message := "false"
		return p.WriteBadRequest(c, message, resp)
	}

	proposalID, err := uuid.Parse(req.ProposalID)
	if err != nil {
		resp := response.ErrorResponse{
			ErrorCode: 400,
			Message:   "Please check your UUID string for errors",
		}
		message := "false"
		return p.WriteBadRequest(c, message, resp)
	}

	commentID, err := uuid.Parse(req.CommentID)
	if err != nil {
		resp := response.ErrorResponse{
			ErrorCode: 400,
			Message:   "Please check your UUID string for errors",
		}
		message := "false"
		return p.WriteBadRequest(c, message, resp)
	}

	err = repository.UpdateCommentByID(p.Session, proposalID, commentID, req.UpdatedComment)
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

// DeleteComment
// @Summary Delete a single comment
// @Description Delete a single comment using proposal and comment id
// @Tags proposal comment
// @Accept plain
// @Produce json
// @Param proposal_id path string true "a common proposal id "
// @Param comment_id path string true "a unique comment id"
// @Success 200 {object} response.Response{Data=string}
// @Failure 400 {object} response.Response{Data=response.ErrorResponse}
// @Failure 500 {object} response.Response{Data=response.ErrorResponse}
// @Router /proposal/comment/delete [delete]
// @Security JWTToken
func (p *CommentsController) DeleteComment(c echo.Context) error {
	proposalIDString := c.QueryParam("proposal-id")
	commentIDString := c.QueryParam("comment-id")

	proposalID, err := uuid.Parse(proposalIDString)
	if err != nil {
		resp := response.ErrorResponse{
			ErrorCode: 400,
			Message:   "Please check your proposal ID in request again for errors",
		}
		message := "false"
		return p.WriteBadRequest(c, message, resp)
	}

	commentID, err := uuid.Parse(commentIDString)
	if err != nil {
		resp := response.ErrorResponse{
			ErrorCode: 400,
			Message:   "Please check your comment ID in request again for errors",
		}
		message := "false"
		return p.WriteBadRequest(c, message, resp)
	}

	err = repository.DeleteCommentByID(p.Session, proposalID, commentID)
	if err != nil {
		resp := response.ErrorResponse{
			ErrorCode: 500,
			Message:   err.Error(),
		}
		message := "false"
		return p.WriteInternalServerError(c, message, resp, "")
	}

	err = proposalRepository.SubtractFromNumberOfComments(p.Session, proposalID)
	if err != nil {
		resp := response.ErrorResponse{
			ErrorCode: 500,
			Message:   "Something went wrong.",
		}
		message := "false"
		return p.WriteInternalServerError(c, message, resp, "")
	}

	return p.WriteSuccess(c, "deleted")
}

// DeleteAllProposalComments
// @Summary Delete some comments
// @Description Delete all comments under a single proposal
// @Tags proposal comment
// @Accept plain
// @Produce json
// @Param proposal_id path string true "a common proposal id "
// @Success 200 {object} response.Response{Data=string}
// @Failure 400 {object} response.Response{Data=response.ErrorResponse}
// @Failure 500 {object} response.Response{Data=response.ErrorResponse}
// @Router /proposal/comment/delete/:proposal-id [delete]
// @Security JWTToken
func (p *CommentsController) DeleteAllProposalComments(c echo.Context) error {
	proposalIDString := c.Param("proposal-id")

	proposalID, err := uuid.Parse(proposalIDString)
	if err != nil {
		resp := response.ErrorResponse{
			ErrorCode: 400,
			Message:   "Please check your proposal ID in request again for errors",
		}
		message := "false"
		return p.WriteBadRequest(c, message, resp)
	}

	err = repository.DeleteAllProposalComments(p.Session, proposalID)
	if err != nil {
		resp := response.ErrorResponse{
			ErrorCode: 500,
			Message:   "Something went wrong",
		}
		message := "false"
		return p.WriteInternalServerError(c, message, resp, "")
	}

	err = proposalRepository.SetCommentsToZero(p.Session, proposalID)
	if err != nil {
		resp := response.ErrorResponse{
			ErrorCode: 500,
			Message:   "Something went wrong.",
		}
		message := "false"
		return p.WriteInternalServerError(c, message, resp, "")
	}

	return p.WriteSuccess(c, "deleted")
}

// UpvoteComment
// @Summary Upvote a single comment
// @Description Upvote a single document using proposal and comment id
// @Tags proposal comment
// @Accept plain
// @Produce json
// @Param proposal_id path string true "a common proposal id "
// @Param comment_id path string true "a unique comment id"
// @Success 200 {object} response.Response{Data=string}
// @Failure 400 {object} response.Response{Data=response.ErrorResponse}
// @Failure 500 {object} response.Response{Data=response.ErrorResponse}
// @Router /proposal/comment/delete [put]
// @Security JWTToken
func (p *CommentsController) UpvoteComment(c echo.Context) error {
	proposalIDString := c.QueryParam("proposal-id")
	commentIDString := c.QueryParam("comment-id")

	proposalID, err := uuid.Parse(proposalIDString)
	if err != nil {
		resp := response.ErrorResponse{
			ErrorCode: 400,
			Message:   "Please check your proposal ID in request again for errors",
		}
		message := "false"
		return p.WriteBadRequest(c, message, resp)
	}

	commentID, err := uuid.Parse(commentIDString)
	if err != nil {
		resp := response.ErrorResponse{
			ErrorCode: 400,
			Message:   "Please check your comment ID in request again for errors",
		}
		message := "false"
		return p.WriteBadRequest(c, message, resp)
	}
	err = repository.UpvoteComment(p.Session, proposalID, commentID)
	if err != nil {
		resp := response.ErrorResponse{
			ErrorCode: 500,
			Message:   "Something went wrong",
		}
		message := "false"
		return p.WriteInternalServerError(c, message, resp, "")
	}

	return p.WriteSuccess(c, "upvoted")
}
