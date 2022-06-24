package repository

import (
	"github.com/gocql/gocql"
)

func InitializeCassandraDB(cassandraHost string) (*gocql.Session, error) {
	var err error
	cluster := gocql.NewCluster(cassandraHost)
	cluster.Keyspace = "user_proposals_and_comments"
	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}

	return session, nil
}

func SchemaMigration(session *gocql.Session) error {
	// err := CreateProposalTables(session)
	// if err != nil && err != gocql.ErrTimeoutNoResponse {
	// 	return err
	// }

	err := CreateCommentsTable(session)
	if err != nil && err != gocql.ErrTimeoutNoResponse {
		return err
	}

	return nil
}

func CreateProposalTables(session *gocql.Session) error {

	// Create Proposal Table By ID
	err := session.Query(`CREATE TABLE proposals_by_id(
			id timeuuid, title text, proposal_text text, user_id uuid, username text,
			firstname text, lastname text, upvotes int, downvotes int, no_of_comments int,
			created_at timestamp, last_updated timestamp,
			PRIMARY KEY (id, created_at, user_id, username)
			); `).Exec()

	if err != nil && err != gocql.ErrTimeoutNoResponse {
		return err
	}

	// Create Proposal Table By UserID
	err = session.Query(`CREATE TABLE proposals_by_user_id(
		id timeuuid, title text, proposal_text text, user_id uuid, username text,
		firstname text, lastname text, upvotes int, downvotes int, no_of_comments int,
		created_at timestamp, last_updated timestamp,
		PRIMARY KEY (user_id, created_at, id, username)
		); `).Exec()

	if err != nil && err != gocql.ErrTimeoutNoResponse {
		return err
	}

	// Create Proposal Table By time created
	err = session.Query(`CREATE TABLE proposals_by_created_at(
		id timeuuid, title text, proposal_text text, user_id uuid, username text,
		firstname text, lastname text, upvotes int, downvotes int, no_of_comments int,
		created_at timestamp, last_updated timestamp,
		PRIMARY KEY (created_at, id, user_id, username)
		); `).Exec()

	return err
}

func CreateCommentsTable(session *gocql.Session) error {

	// Create Comment Table
	err := session.Query(`CREATE TABLE comments_by_proposal_id(
			proposal_id uuid, id timeuuid, comment text, user_posted_id uuid, user_posted_username text,
			user_commented_id uuid, user_commented_username text, upvotes int,
			created_at timestamp, last_updated timestamp,
			PRIMARY KEY (proposal_id, created_at, id)
			); `).Exec()

	if err != nil && err != gocql.ErrTimeoutNoResponse {
		return err
	}

	err = session.Query(`CREATE TABLE comments_by_proposal_and_comment_id(
			proposal_id uuid, id timeuuid, comment text, user_posted_id uuid, user_posted_username text,
			user_commented_id uuid, user_commented_username text, upvotes int,
			created_at timestamp, last_updated timestamp,
			PRIMARY KEY (proposal_id, id, created_at)
			); `).Exec()

	return err
}
