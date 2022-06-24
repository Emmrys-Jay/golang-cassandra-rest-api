package repository

import (
	"time"

	"github.com/gocql/gocql"
	"github.com/google/uuid"
	"github.com/windswept321/smartest-city-roadmap-go/proposals-and-comments/entity"
)

func StoreProposal(session *gocql.Session, title string, proposalText string, userID uuid.UUID, username, firstname, lastname string) error {

	updateTime := time.Now()
	id := gocql.UUIDFromTime(time.Now())

	err := session.Query(`INSERT INTO proposals_by_id(user_id, id, username, title, proposal_text, created_at, last_updated, upvotes, downvotes, no_of_comments, firstname, lastname) VALUES 
					(?, ?, ?, ?, ?, ?, ?, 0, 0, 0, ?, ?);`, gocql.UUID(userID), id, username, title, proposalText, updateTime, updateTime, firstname, lastname).Exec()

	if err != nil {
		return err
	}

	err = session.Query(`INSERT INTO proposals_by_user_id(user_id, id, username, title, proposal_text, created_at, last_updated, upvotes, downvotes, no_of_comments, firstname, lastname) VALUES 
					(?, ?, ?, ?, ?, ?, ?, 0, 0, 0, ?, ?);`, gocql.UUID(userID), id, username, title, proposalText, updateTime, updateTime, firstname, lastname).Exec()

	if err != nil {
		return err
	}

	err = session.Query(`INSERT INTO proposals_by_created_at(user_id, id, username, title, proposal_text, created_at, last_updated, upvotes, downvotes, no_of_comments, firstname, lastname) VALUES 
					(?, ?, ?, ?, ?, ?, ?, 0, 0, 0, ?, ?);`, gocql.UUID(userID), id, username, title, proposalText, updateTime, updateTime, firstname, lastname).Exec()

	return err
}

// GetAllProposals returns all stored proposals starting with the most recently created
func GetAllProposals(session *gocql.Session) ([]entity.Proposal, error) {
	var proposals []entity.Proposal
	var inversedProposals []entity.Proposal

	var m = map[string]interface{}{}

	iter := session.Query(`SELECT * FROM proposals_by_created_at;`).Iter()

	for iter.MapScan(m) {
		proposals = append(proposals, entity.Proposal{
			ID:           uuid.UUID(m["id"].(gocql.UUID)),
			Title:        m["title"].(string),
			ProposalText: m["proposal_text"].(string),
			UserID:       uuid.UUID(m["user_id"].(gocql.UUID)),
			Username:     m["username"].(string),
			FirstName:    m["firstname"].(string),
			LastName:     m["lastname"].(string),
			UpVotes:      m["upvotes"].(int),
			DownVotes:    m["downvotes"].(int),
			NoOfComments: m["no_of_comments"].(int),
			CreatedAt:    m["created_at"].(time.Time),
			LastUpdated:  m["last_updated"].(time.Time),
		})
		m = map[string]interface{}{}
	}

	err := iter.Close()
	if err != nil {
		return proposals, err
	}

	for i := len(proposals) - 1; i >= 0; i-- {
		inversedProposals = append(inversedProposals, proposals[i])
	}
	return inversedProposals, err
}

func GetProposalsByUserID(session *gocql.Session, userID uuid.UUID) ([]entity.Proposal, error) {
	var proposals []entity.Proposal

	var m = map[string]interface{}{}

	iter := session.Query(`SELECT * FROM proposals_by_user_id
							WHERE user_id = ?
							ORDER BY created_at DESC;`, gocql.UUID(userID)).Iter()

	for iter.MapScan(m) {
		proposals = append(proposals, entity.Proposal{
			ID:           uuid.UUID(m["id"].(gocql.UUID)),
			Title:        m["title"].(string),
			ProposalText: m["proposal_text"].(string),
			UserID:       uuid.UUID(m["user_id"].(gocql.UUID)),
			Username:     m["username"].(string),
			FirstName:    m["firstname"].(string),
			LastName:     m["lastname"].(string),
			UpVotes:      m["upvotes"].(int),
			DownVotes:    m["downvotes"].(int),
			NoOfComments: m["no_of_comments"].(int),
			CreatedAt:    m["created_at"].(time.Time),
			LastUpdated:  m["last_updated"].(time.Time),
		})
		m = map[string]interface{}{}
	}

	err := iter.Close()

	return proposals, err
}

func GetProposalsByTimeCreated(session *gocql.Session, dateFrom time.Time, dateTo time.Time) ([]entity.Proposal, error) {
	var proposals []entity.Proposal
	var inversedProposals []entity.Proposal

	var m = map[string]interface{}{}

	iter := session.Query(`SELECT * FROM proposals_by_created_at
							WHERE created_at>=? AND created_at<=?
							ALLOW FILTERING;`, dateFrom, dateTo).Iter()

	for iter.MapScan(m) {
		proposals = append(proposals, entity.Proposal{
			ID:           uuid.UUID(m["id"].(gocql.UUID)),
			Title:        m["title"].(string),
			ProposalText: m["proposal_text"].(string),
			UserID:       uuid.UUID(m["user_id"].(gocql.UUID)),
			Username:     m["username"].(string),
			FirstName:    m["firstname"].(string),
			LastName:     m["lastname"].(string),
			UpVotes:      m["upvotes"].(int),
			DownVotes:    m["downvotes"].(int),
			NoOfComments: m["no_of_comments"].(int),
			CreatedAt:    m["created_at"].(time.Time),
			LastUpdated:  m["last_updated"].(time.Time),
		})
		m = map[string]interface{}{}
	}

	err := iter.Close()

	if err != nil {
		return proposals, err
	}

	for i := len(proposals) - 1; i >= 0; i-- {
		inversedProposals = append(inversedProposals, proposals[i])
	}
	return inversedProposals, err
}

func GetProposalByProposalID(session *gocql.Session, proposalID uuid.UUID) ([]entity.Proposal, error) {
	var proposals []entity.Proposal
	var m = map[string]interface{}{}

	iter := session.Query(`SELECT * FROM proposals_by_id WHERE id=? LIMIT 1;`, gocql.UUID(proposalID)).Iter()

	for iter.MapScan(m) {
		proposals = append(proposals, entity.Proposal{
			ID:           uuid.UUID(m["id"].(gocql.UUID)),
			Title:        m["title"].(string),
			ProposalText: m["proposal_text"].(string),
			UserID:       uuid.UUID(m["user_id"].(gocql.UUID)),
			Username:     m["username"].(string),
			FirstName:    m["firstname"].(string),
			LastName:     m["lastname"].(string),
			UpVotes:      m["upvotes"].(int),
			DownVotes:    m["downvotes"].(int),
			NoOfComments: m["no_of_comments"].(int),
			CreatedAt:    m["created_at"].(time.Time),
			LastUpdated:  m["last_updated"].(time.Time),
		})
		m = map[string]interface{}{}
	}

	err := iter.Close()

	return proposals, err
}

func UpdateProposal(session *gocql.Session, proposalID uuid.UUID, title, proposalText string) error {

	proposal, err := GetProposalByProposalID(session, proposalID)
	if err != nil {
		return err
	}
	updateTime := time.Now()

	err = session.Query(`UPDATE proposals_by_id SET title=?, proposal_text=?, last_updated=?
							WHERE id=? AND user_id=? AND created_at=? AND username=?`, title, proposalText, updateTime, gocql.UUID(proposalID), gocql.UUID(proposal[0].UserID), proposal[0].CreatedAt, proposal[0].Username).Exec()

	if err != nil {
		return err
	}

	err = session.Query(`UPDATE proposals_by_created_at SET title=?, proposal_text=?, last_updated=?
							WHERE id=? AND user_id=? AND created_at=? AND username=?`, title, proposalText, updateTime, gocql.UUID(proposalID), gocql.UUID(proposal[0].UserID), proposal[0].CreatedAt, proposal[0].Username).Exec()

	if err != nil {
		return err
	}

	err = session.Query(`UPDATE proposals_by_user_id SET title=?, proposal_text=?, last_updated=?
							WHERE id=? AND user_id=? AND created_at=? AND username=?`, title, proposalText, updateTime, gocql.UUID(proposalID), gocql.UUID(proposal[0].UserID), proposal[0].CreatedAt, proposal[0].Username).Exec()

	return err
}

func DeleteProposal(session *gocql.Session, proposalID uuid.UUID) error {
	proposal, err := GetProposalByProposalID(session, proposalID)
	if err != nil {
		return err
	}

	err = session.Query(`DELETE FROM proposals_by_id
							WHERE id=? AND user_id=? AND created_at=? AND username=?`, gocql.UUID(proposalID), gocql.UUID(proposal[0].UserID), proposal[0].CreatedAt, proposal[0].Username).Exec()

	if err != nil {
		return err
	}

	err = session.Query(`DELETE FROM proposals_by_created_at
							WHERE id=? AND user_id=? AND created_at=? AND username=?`, gocql.UUID(proposalID), gocql.UUID(proposal[0].UserID), proposal[0].CreatedAt, proposal[0].Username).Exec()

	if err != nil {
		return err
	}

	err = session.Query(`DELETE FROM proposals_by_user_id
							WHERE id=? AND user_id=? AND created_at=? AND username=?`, gocql.UUID(proposalID), gocql.UUID(proposal[0].UserID), proposal[0].CreatedAt, proposal[0].Username).Exec()

	return err
}

func DeleteAllProposals(session *gocql.Session) error {
	err := session.Query(`TRUNCATE TABLE user_proposals_and_comments.proposals_by_id;`).Exec()

	if err != nil {
		return err
	}

	err = session.Query(`TRUNCATE TABLE user_proposals_and_comments.proposals_by_user_id;`).Exec()

	if err != nil {
		return err
	}

	err = session.Query(`TRUNCATE TABLE user_proposals_and_comments.proposals_by_created_at;`).Exec()

	return err
}

func UpvoteProposal(session *gocql.Session, proposalID uuid.UUID) error {

	proposal, err := GetProposalByProposalID(session, proposalID)
	if err != nil {
		return err
	}

	err = session.Query(`UPDATE proposals_by_id SET upvotes=?
							WHERE id=? AND user_id=? AND created_at=? AND username=?`, proposal[0].UpVotes+1, gocql.UUID(proposalID), gocql.UUID(proposal[0].UserID), proposal[0].CreatedAt, proposal[0].Username).Exec()

	if err != nil {
		return err
	}

	err = session.Query(`UPDATE proposals_by_created_at SET upvotes=?
							WHERE id=? AND user_id=? AND created_at=? AND username=?`, proposal[0].UpVotes+1, gocql.UUID(proposalID), gocql.UUID(proposal[0].UserID), proposal[0].CreatedAt, proposal[0].Username).Exec()

	if err != nil {
		return err
	}

	err = session.Query(`UPDATE proposals_by_user_id SET upvotes=?
							WHERE id=? AND user_id=? AND created_at=? AND username=?`, proposal[0].UpVotes+1, gocql.UUID(proposalID), gocql.UUID(proposal[0].UserID), proposal[0].CreatedAt, proposal[0].Username).Exec()

	return err
}

func DownvoteProposal(session *gocql.Session, proposalID uuid.UUID) error {

	proposal, err := GetProposalByProposalID(session, proposalID)
	if err != nil {
		return err
	}

	err = session.Query(`UPDATE proposals_by_id SET downvotes=?
							WHERE id=? AND user_id=? AND created_at=? AND username=?`, proposal[0].DownVotes+1, gocql.UUID(proposalID), gocql.UUID(proposal[0].UserID), proposal[0].CreatedAt, proposal[0].Username).Exec()

	if err != nil {
		return err
	}

	err = session.Query(`UPDATE proposals_by_created_at SET downvotes=?
							WHERE id=? AND user_id=? AND created_at=? AND username=?`, proposal[0].DownVotes+1, gocql.UUID(proposalID), gocql.UUID(proposal[0].UserID), proposal[0].CreatedAt, proposal[0].Username).Exec()

	if err != nil {
		return err
	}

	err = session.Query(`UPDATE proposals_by_user_id SET downvotes=?
							WHERE id=? AND user_id=? AND created_at=? AND username=?`, proposal[0].DownVotes+1, gocql.UUID(proposalID), gocql.UUID(proposal[0].UserID), proposal[0].CreatedAt, proposal[0].Username).Exec()

	return err
}

func AddToNumberOfComments(session *gocql.Session, proposalID uuid.UUID) error {

	proposal, err := GetProposalByProposalID(session, proposalID)
	if err != nil {
		return err
	}

	err = session.Query(`UPDATE proposals_by_id SET no_of_comments=?
							WHERE id=? AND user_id=? AND created_at=? AND username=?`, proposal[0].NoOfComments+1, gocql.UUID(proposalID), gocql.UUID(proposal[0].UserID), proposal[0].CreatedAt, proposal[0].Username).Exec()

	if err != nil {
		return err
	}

	err = session.Query(`UPDATE proposals_by_user_id SET no_of_comments=?
							WHERE id=? AND user_id=? AND created_at=? AND username=?`, proposal[0].NoOfComments+1, gocql.UUID(proposalID), gocql.UUID(proposal[0].UserID), proposal[0].CreatedAt, proposal[0].Username).Exec()

	if err != nil {
		return err
	}

	err = session.Query(`UPDATE proposals_by_created_at SET no_of_comments=?
							WHERE id=? AND user_id=? AND created_at=? AND username=?`, proposal[0].NoOfComments+1, gocql.UUID(proposalID), gocql.UUID(proposal[0].UserID), proposal[0].CreatedAt, proposal[0].Username).Exec()

	return err
}

func SubtractFromNumberOfComments(session *gocql.Session, proposalID uuid.UUID) error {
	proposal, err := GetProposalByProposalID(session, proposalID)
	if err != nil {
		return err
	}

	err = session.Query(`UPDATE proposals_by_id SET no_of_comments=?
							WHERE id=? AND user_id=? AND created_at=? AND username=?`, proposal[0].NoOfComments-1, gocql.UUID(proposalID), gocql.UUID(proposal[0].UserID), proposal[0].CreatedAt, proposal[0].Username).Exec()

	if err != nil {
		return err
	}

	err = session.Query(`UPDATE proposals_by_user_id SET no_of_comments=?
							WHERE id=? AND user_id=? AND created_at=? AND username=?`, proposal[0].NoOfComments-1, gocql.UUID(proposalID), gocql.UUID(proposal[0].UserID), proposal[0].CreatedAt, proposal[0].Username).Exec()

	if err != nil {
		return err
	}

	err = session.Query(`UPDATE proposals_by_created_at SET no_of_comments=?
							WHERE id=? AND user_id=? AND created_at=? AND username=?`, proposal[0].NoOfComments-1, gocql.UUID(proposalID), gocql.UUID(proposal[0].UserID), proposal[0].CreatedAt, proposal[0].Username).Exec()

	return err
}

func SetCommentsToZero(session *gocql.Session, proposalID uuid.UUID) error {
	proposal, err := GetProposalByProposalID(session, proposalID)
	if err != nil {
		return err
	}

	err = session.Query(`UPDATE proposals_by_id SET no_of_comments=?
							WHERE id=? AND user_id=? AND created_at=? AND username=?`, 0, gocql.UUID(proposalID), gocql.UUID(proposal[0].UserID), proposal[0].CreatedAt, proposal[0].Username).Exec()

	if err != nil {
		return err
	}

	err = session.Query(`UPDATE proposals_by_user_id SET no_of_comments=?
							WHERE id=? AND user_id=? AND created_at=? AND username=?`, 0, gocql.UUID(proposalID), gocql.UUID(proposal[0].UserID), proposal[0].CreatedAt, proposal[0].Username).Exec()

	if err != nil {
		return err
	}

	err = session.Query(`UPDATE proposals_by_created_at SET no_of_comments=?
							WHERE id=? AND user_id=? AND created_at=? AND username=?`, 0, gocql.UUID(proposalID), gocql.UUID(proposal[0].UserID), proposal[0].CreatedAt, proposal[0].Username).Exec()

	return err
}
