package entity

import (
	"time"

	"github.com/google/uuid"
)

type Proposal struct {
	ID           uuid.UUID `json:"id,omitempty"  form:"id"`
	Title        string    `json:"title,omitempty" form:"title"`
	ProposalText string    `json:"proposal_text,omitempty"  form:"proposal_text"`
	UserID       uuid.UUID `json:"user_id,omitempty"  form:"user_id"`
	Username     string    `json:"username,omitempty"  form:"username"`
	FirstName    string    `json:"firstname,omitempty"  form:"firstname"`
	LastName     string    `json:"lastname,omitempty"  form:"lastname"`
	UpVotes      int       `json:"upvotes,omitempty"`
	DownVotes    int       `json:"downvotes,omitempty"`
	NoOfComments int       `json:"no_of_comments,omitempty" form:"no_of_comments"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
	LastUpdated  time.Time `json:"last_updated,omitempty"`
}
