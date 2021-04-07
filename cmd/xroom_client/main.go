package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gocql/gocql"
	"github.com/jklaw90/shinfo/internal/room"
)

func main() {
	roomService, err := room.NewRoomClient("localhost:9000")
	if err != nil {
		os.Exit(1)
	}
	s := time.Now()

	// a := true
	// t := model.Classic
	// r := model.RoomCreate{
	// 	Name:     "Really cool room",
	// 	Type:     &t,
	// 	Public:   &a,
	// 	Archived: &a,
	// }

	// room, err := roomService.Create(context.Background(), r)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// for i := 0; i < 100; i++ {
	// 	id, _ := gocql.RandomUUID()
	// 	fmt.Println(roomService.AddUser(context.Background(), room.ID, model.User{
	// 		ID:     id,
	// 		Name:   fmt.Sprintf("User %d", i),
	// 		Avatar: "avatar here",
	// 	}))
	// }
	id, _ := gocql.ParseUUID("21b2901e-4e38-11eb-814e-acde48001122")
	lastSeenID, _ := gocql.ParseUUID("21c22f56-4e38-11eb-8158-acde48001122")
	var next *gocql.UUID
	next = &lastSeenID
	for next != nil {
		resp, err := roomService.GetUsers(context.Background(), id, 1, next)
		fmt.Println(resp, err)
		next = resp.NextID
	}

	userID, _ := gocql.ParseUUID("e349baac-f1fc-4cb2-ae66-d79c06d5396f")
	id, _ = gocql.ParseUUID("21b2901e-4e38-11eb-814e-acde48001122")
	newResp, err := roomService.GetByUserID(context.Background(), userID, 1, nil)
	fmt.Println("here", newResp, err)
	fmt.Println(time.Since(s))
}
