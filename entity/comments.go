package entity

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	ProposalID            uuid.UUID `json:"proposal_id,omitempty" form:"proposal_id"` //Partition key
	CommentID             uuid.UUID `json:"id,omitempty" form:"id"`
	CommentText           string    `json:"comment,omitempty" form:"id"`
	UserPostedProposalID  uuid.UUID `json:"user_posted_id,omitempty" form:"posted_user_id"`
	UserPostedUsername    string    `json:"user_posted,omitempty" form:"user_posted"`
	UserCommentedID       uuid.UUID `json:"user_commented_id,omitempty" form:"user_commented_id"`
	UserCommentedUsername string    `json:"user_commented,omitempty" form:"user_commented"`
	UpVotes               int       `json:"upvotes,omitempty" form:"upvotes"`
	CreatedAt             time.Time `json:"created_at,omitempty"`
	LastUpdated           time.Time `json:"last_updated,omitempty"`
}
