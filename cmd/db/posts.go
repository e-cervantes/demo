package db

import (
	"database/sql"
	"fmt"

	"demo.com/demo-1/model"
)

var (
	ErrNoRecord = fmt.Errorf("no matching record found")
	insertOp    = "insert"
	deleteOp    = "delete"
	updateOp    = "update"
)

// function to save to database
func (db Database) SavePost(post *model.Post) error {
	var id int
	query := "INSERT INTO posts(title, body) VALUES ($1, $2) RETURNING id"
	err := db.Conn.QueryRow(query, post.Title, post.Body).Scan(&id)

	if err != nil {
		return err
	}

	logQuery := "INSERT INTO post_logs(post_id, operation) VALUES ($1,$2)"
	post.ID = id
	_, err = db.Conn.Exec(logQuery, post.ID, insertOp)

	if err != nil {
		return err
	}

	return nil

}

func (db Database) UpdatePost(postId int, post model.Post) error {
	query := "UPDATE posts SET title = $1, BODY = $2, WHERE id = $3"
	_, err := db.Conn.Exec(query, post.Title, post.Body, post.ID)

	if err != nil {
		return err
	}

	post.ID = postId

	logQuery := "INSERT INTO post_logs(post_id, operation) VALUES ($1, $2)"
	_, err = db.Conn.Exec(logQuery, post.ID, updateOp)

	if err != nil {
		db.Logger.Err(err).Msg("could not log operation for logstash")
	}

	return nil

}

func (db Database) DeletePost(postId int) error {
	query := "DELETE FROM Posts WHERE id = $1"
	_, err := db.Conn.Exec(query, postId)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrNoRecord
		}
		return err
	}

	logQuery := "INSERT INTO post_logs(post_id, operation) VALUES ($1, $2)"
	_, err = db.Conn.Exec(logQuery, postId, deleteOp)
	if err != nil {
		db.Logger.Err(err).Msg("could not log operation for logstash")
	}
	return nil
}
