package database

import (
	"context"

	"github.com/gocql/gocql"
	"github.com/jklaw90/shinfo/pkg/model"
)

type MessageRepo struct {
	session *gocql.Session
}

func NewDB(session *gocql.Session) *MessageRepo {
	return &MessageRepo{
		session: session,
	}
}

func (db *MessageRepo) Get(ctx context.Context, roomID, msgID gocql.UUID) (model.Message, error) {
	var msg model.Message

	if err := db.session.
		Query(
			`SELECT id, room_id, user, message_type, entry, metadata FROM shinfo.messages
				WHERE room_id = ? AND id = ?`, roomID, msgID).
		Scan(&msg.ID, &msg.RoomID, &msg.User, &msg.Type, &msg.Entry, &msg.Metadata); err != nil {
		return msg, err
	}
	return msg, nil
}

func (db *MessageRepo) GetMessages(ctx context.Context, roomID gocql.UUID, limit int64, lastSeenID *gocql.UUID) ([]model.Message, error) {
	var msgs []model.Message
	var msg model.Message

	iter := db.session.
		Query(
			`SELECT id, room_id, user, message_type, entry, metadata FROM shinfo.messages
				WHERE room_id = ? AND id < ?
				LIMIT ?`, roomID, lastSeenID, limit).Iter()
	for iter.Scan(&msg.ID, &msg.RoomID, &msg.User, &msg.Type, &msg.Entry, &msg.Metadata) {
		msgs = append(msgs, msg)
	}
	return msgs, nil
}

func (db *MessageRepo) Create(ctx context.Context, msg model.Message) error {
	return db.session.
		Query(
			`INSERT INTO shinfo.messages (id, room_id, user, message_type, entry, metadata)
				VALUES (?, ?, ?, ?, ?, ?)`,
			msg.ID, msg.RoomID, msg.User, msg.Type, msg.Entry, msg.Metadata,
		).
		Exec()
}

func (db *MessageRepo) Delete(ctx context.Context, roomID, messageID gocql.UUID) error {
	if err := db.session.
		Query(`DELETE FROM shinfo.messages WHERE room_id = ? and id = ?`, roomID, messageID).
		Exec(); err != nil {
		return err
	}
	return nil
}
