package database

import (
	"context"
	"time"

	"github.com/gocql/gocql"
	"github.com/jklaw90/shinfo/pkg/model"
)

type RoomRepo struct {
	session *gocql.Session
}

func NewDB(session *gocql.Session) *RoomRepo {
	return &RoomRepo{
		session: session,
	}
}

func (db *RoomRepo) Get(ctx context.Context, id gocql.UUID) (model.Room, error) {
	var room model.Room
	err := db.session.
		Query(`SELECT id, "name", users, archived, public, "type", created_at FROM shinfo.rooms WHERE id = ?`, id).
		Scan(&room.ID, &room.Name, &room.Users, &room.Archived, &room.Public, &room.Type, &room.CreatedAt)
	return room, err
}

func (db *RoomRepo) GetByIDs(ctx context.Context, ids ...gocql.UUID) ([]model.Room, error) {
	var rooms []model.Room
	var room model.Room

	iter := db.session.Query(`SELECT id, "name", users, archived, public, "type", created_at FROM shinfo.rooms WHERE user_id in ?`, ids).Iter()
	for iter.Scan(&room.ID, &room.Name, &room.Users, &room.Archived, &room.Public, &room.Type, &room.CreatedAt) {
		rooms = append(rooms, room)
	}
	iter.Close()

	return rooms, nil
}

func (db *RoomRepo) GetByUserID(ctx context.Context, id gocql.UUID, limit int64, lastSeenID *gocql.UUID) ([]model.RoomByUser, error) {
	var rooms []model.RoomByUser
	var room model.RoomByUser

	iter := db.session.Query(`SELECT room_id, room_name, archived, created_at, joined_at, "type", public FROM shinfo.rooms_by_user
								WHERE user_id = ? and room_id < ?
								ORDER BY room_id DESC
								Limit ?`, id, lastSeenID, limit).Iter()
	for iter.Scan(&room.RoomID, &room.RoomName, &room.Archived, &room.CreatedAt, &room.JoinedAt, &room.Type, &room.Public) {
		rooms = append(rooms, room)
	}
	iter.Close()

	return rooms, nil
}

func (db *RoomRepo) GetUsersByRoomID(ctx context.Context, id gocql.UUID, limit int64, lastSeenID *gocql.UUID) ([]model.UserByRoom, error) {
	var rooms []model.UserByRoom
	var room model.UserByRoom

	iter := db.session.Query(`SELECT room_id, joined_at, user_id, user_name, user_avatar FROM shinfo.users_by_room
								WHERE room_id = ? AND joined_at > ?
								ORDER BY joined_at ASC
								LIMIT ?`, id, lastSeenID, limit).Iter()
	for iter.Scan(&room.RoomID, &room.JoinedAt, &room.User.ID, &room.User.Name, &room.User.Avatar) {
		rooms = append(rooms, room)
	}
	iter.Close()

	return rooms, nil
}

// func (db *RoomRepo) GetLastMsgByIDs(ctx context.Context, ids ...gocql.UUID) ([]msgmodel.Message, error) {
// 	var msgs []msgmodel.Message
// 	var msg msgmodel.Message

// 	iter := db.session.
// 		Query(
// 			`SELECT id, room_id, user_id, message_type, entry, metadata FROM shinfo.messages
// 				WHERE room_id in ?
// 				PER PARTITION LIMIT 1`, ids).
// 		Iter()
// 	for iter.Scan(&msg.ID, &msg.RoomID, &msg.User, &msg.Type, &msg.Entry, &msg.Metadata) {
// 		msgs = append(msgs, msg)
// 	}

// 	return msgs, nil
// }

func (db *RoomRepo) Create(ctx context.Context, room model.Room) error {
	return db.session.
		Query(
			`INSERT INTO shinfo.rooms (id, "name", users, "type", public, archived, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)`,
			room.ID, room.Name, room.Users, room.Type, room.Public, room.Archived, room.CreatedAt,
		).
		Exec()
}

func (db *RoomRepo) Delete(ctx context.Context, roomID gocql.UUID) error {
	room, err := db.Get(ctx, roomID)
	if err != nil {
		return err
	}

	for _, userID := range room.Users {
		if err := db.session.
			Query(`DELETE FROM shinfo.rooms_by_user WHERE user_id = ? and room_id = ?`, userID, roomID).
			Exec(); err != nil {
			return err
		}
	}

	if err := db.session.
		Query(`DELETE FROM shinfo.rooms WHERE id=?`, roomID).
		Exec(); err != nil {
		return err
	}

	return nil
}

func (db *RoomRepo) AddUser(ctx context.Context, room model.Room, user model.User) error {
	if err := db.session.
		Query(`UPDATE shinfo.rooms SET users = users + ? WHERE id = ?`, []model.User{user}, room.ID).
		Exec(); err != nil {
		return err
	}

	now := time.Now().UTC()
	nowUUID := gocql.UUIDFromTime(now)
	if err := db.session.
		Query(`INSERT INTO shinfo.rooms_by_user (user_id, room_id, room_name, archived, public, "type", created_at, joined_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			user.ID, room.ID, room.Name, room.Archived, room.Public, room.Type, room.CreatedAt, now).Exec(); err != nil {
		return err
	}

	if err := db.session.
		Query(`INSERT INTO shinfo.users_by_room (room_id, joined_at, user_id, user_name, user_avatar) VALUES (?, ?, ?, ?, ?)`,
			room.ID, nowUUID, user.ID, user.Name, user.Avatar).Exec(); err != nil {
		return err
	}

	return nil
}

func (db *RoomRepo) RemoveUser(ctx context.Context, roomUUID model.Room, user model.User) error {
	if err := db.session.
		Query(`UPDATE shinfo.rooms SET users = users - {?} WHERE id = ?`, user, roomUUID).
		Exec(); err != nil {
		return err
	}

	if err := db.session.
		Query(`DELETE FROM shinfo.rooms_by_user WHERE user_id=? AND room_id=?`, user.ID, roomUUID).
		Exec(); err != nil {
		return err
	}

	if err := db.session.
		Query(`DELETE FROM shinfo.users_by_room WHERE room_id=? AND user_id=?`, roomUUID, user.ID).
		Exec(); err != nil {
		return err
	}

	return nil
}
