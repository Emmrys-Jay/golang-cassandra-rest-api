package repository

import (
	"fmt"
	"time"

	"github.com/gocql/gocql"
	"github.com/google/uuid"
	"github.com/windswept321/smartest-city-roadmap-go/proposals-and-comments/entity"
	"github.com/windswept321/smartest-city-roadmap-go/proposals-and-comments/proposals/repository"
)

func StoreComment(session *gocql.Session, proposalID uuid.UUID, comment string, userID uuid.UUID, username string) error {
	uID := gocql.UUID(userID)
	if uID == gocql.UUID(uuid.Nil) {
		return fmt.Errorf("something went wrong")
	}

	if comment == "" {
		return fmt.Errorf("invalid request")
	}

	time := time.Now()
	proposal, err := repository.GetProposalByProposalID(session, proposalID)
	if err != nil {
		return err
	}

	commentID := gocql.UUIDFromTime(time)

	err = session.Query(`INSERT INTO comments_by_proposal_id(proposal_id, id, comment, user_posted_id, user_posted_username, 
		user_commented_id, user_commented_username, created_at, last_updated, upvotes) VALUES 
		(?, ?, ?, ?, ?, ?, ?, ?, ?, 0);`, gocql.UUID(proposalID), commentID, comment, gocql.UUID(proposal[0].UserID),
		proposal[0].Username, gocql.UUID(userID), username, time, time).Exec()

	if err != nil {
		return err
	}

	err = session.Query(`INSERT INTO comments_by_proposal_and_comment_id(proposal_id, id, comment, user_posted_id, user_posted_username, 
		user_commented_id, user_commented_username, created_at, last_updated, upvotes) VALUES 
		(?, ?, ?, ?, ?, ?, ?, ?, ?, 0);`, gocql.UUID(proposalID), commentID, comment, gocql.UUID(proposal[0].UserID),
		proposal[0].Username, gocql.UUID(userID), username, time, time).Exec()

	return err

}

func GetCommentsByProposalID(session *gocql.Session, proposalID uuid.UUID) ([]entity.Comment, error) {
	var comments []entity.Comment

	var m = map[string]interface{}{}

	iter := session.Query(`SELECT * FROM comments_by_proposal_id
							WHERE proposal_id=?
							ORDER BY created_at DESC;`, gocql.UUID(proposalID)).Iter()

	for iter.MapScan(m) {
		comments = append(comments, entity.Comment{
			ProposalID:            uuid.UUID(m["proposal_id"].(gocql.UUID)),
			CommentID:             uuid.UUID(m["id"].(gocql.UUID)),
			CommentText:           m["comment"].(string),
			UserPostedProposalID:  uuid.UUID(m["user_posted_id"].(gocql.UUID)),
			UserPostedUsername:    m["user_posted_username"].(string),
			UserCommentedID:       uuid.UUID(m["user_commented_id"].(gocql.UUID)),
			UserCommentedUsername: m["user_commented_username"].(string),
			UpVotes:               m["upvotes"].(int),
			CreatedAt:             m["created_at"].(time.Time),
			LastUpdated:           m["last_updated"].(time.Time),
		})
		m = map[string]interface{}{}
	}

	err := iter.Close()

	return comments, err
}

func GetCommentByIDAndProposalID(session *gocql.Session, proposalID uuid.UUID, commentID uuid.UUID) (*entity.Comment, error) {
	var comment *entity.Comment

	var m = map[string]interface{}{}

	iter := session.Query(`SELECT * FROM comments_by_proposal_and_comment_id
							WHERE proposal_id=? AND id=? LIMIT 1;`, gocql.UUID(proposalID), gocql.UUID(commentID)).Iter()

	for iter.MapScan(m) {
		comment = &entity.Comment{
			ProposalID:            uuid.UUID(m["proposal_id"].(gocql.UUID)),
			CommentID:             uuid.UUID(m["id"].(gocql.UUID)),
			CommentText:           m["comment"].(string),
			UserPostedProposalID:  uuid.UUID(m["user_posted_id"].(gocql.UUID)),
			UserPostedUsername:    m["user_posted_username"].(string),
			UserCommentedID:       uuid.UUID(m["user_commented_id"].(gocql.UUID)),
			UserCommentedUsername: m["user_commented_username"].(string),
			UpVotes:               m["upvotes"].(int),
			CreatedAt:             m["created_at"].(time.Time),
			LastUpdated:           m["last_updated"].(time.Time),
		}
	}

	err := iter.Close()

	return comment, err
}

func UpdateCommentByID(session *gocql.Session, proposalID uuid.UUID, commentID uuid.UUID, updatedComment string) error {
	comment, err := GetCommentByIDAndProposalID(session, proposalID, commentID)
	if err != nil {
		return err
	}

	updateTime := time.Now()

	err = session.Query(`UPDATE comments_by_proposal_id SET comment=?, last_updated=?
							WHERE proposal_id=? AND id=? AND created_at=?;`, updatedComment, updateTime, gocql.UUID(proposalID), gocql.UUID(commentID), comment.CreatedAt).Exec()

	if err != nil {
		return err
	}

	err = session.Query(`UPDATE comments_by_proposal_and_comment_id SET comment=?, last_updated=?
							WHERE proposal_id=? AND id=? AND created_at=?;`, updatedComment, updateTime, gocql.UUID(proposalID), gocql.UUID(commentID), comment.CreatedAt).Exec()

	return err
}

func DeleteCommentByID(session *gocql.Session, proposalID uuid.UUID, commentID uuid.UUID) error {
	comment, err := GetCommentByIDAndProposalID(session, proposalID, commentID)
	if err != nil {
		return err
	}

	err = session.Query(`DELETE FROM comments_by_proposal_id
							WHERE proposal_id=? AND id=? AND created_at=?`, gocql.UUID(proposalID), gocql.UUID(commentID), comment.CreatedAt).Exec()

	if err != nil {
		return err
	}

	err = session.Query(`DELETE FROM comments_by_proposal_and_comment_id
							WHERE proposal_id=? AND id=? AND created_at=?`, gocql.UUID(proposalID), gocql.UUID(commentID), comment.CreatedAt).Exec()

	return err
}

func DeleteAllProposalComments(session *gocql.Session, proposalID uuid.UUID) error {
	err := session.Query(`DELETE FROM comments_by_proposal_id
							WHERE proposal_id=?`, gocql.UUID(proposalID)).Exec()

	if err != nil {
		return err
	}

	err = session.Query(`DELETE FROM comments_by_proposal_and_comment_id
							WHERE proposal_id=?`, gocql.UUID(proposalID)).Exec()

	return err
}

func DeleteAllComments(session *gocql.Session) error {
	err := session.Query(`TRUNCATE TABLE user_proposals_and_comments.comments_by_proposal_id`).Exec()

	if err != nil {
		return err
	}

	err = session.Query(`TRUNCATE TABLE user_proposals_and_comments.comments_by_proposal_and_comment_id`).Exec()

	return err
}

func UpvoteComment(session *gocql.Session, proposalID uuid.UUID, commentID uuid.UUID) error {
	comment, err := GetCommentByIDAndProposalID(session, proposalID, commentID)
	if err != nil {
		return err
	}

	err = session.Query(`UPDATE comments_by_proposal_id SET upvotes=?
							WHERE proposal_id=? AND id=? AND created_at=?;`, comment.UpVotes+1, gocql.UUID(proposalID), gocql.UUID(commentID), comment.CreatedAt).Exec()

	if err != nil {
		return err
	}

	err = session.Query(`UPDATE comments_by_proposal_and_comment_id SET upvotes=?
							WHERE proposal_id=? AND id=? AND created_at=?;`, comment.UpVotes+1, gocql.UUID(proposalID), gocql.UUID(commentID), comment.CreatedAt).Exec()

	return err
}
