Create TYPE shinfo.user(
   id uuid,
   name text,
   avatar text
);

CREATE TABLE shinfo.rooms_by_user(
  user_id uuid,
  room_id timeuuid,
  room_name text,
  archived boolean,
  public boolean,
  "type" int,
  created_at timestamp,
  joined_at timestamp,
  PRIMARY KEY(user_id, room_id),
) WITH CLUSTERING ORDER BY (room_id DESC);

CREATE TABLE shinfo.users_by_room(
  room_id timeuuid,
  joined_at timeuuid,
  user_id uuid,
  user_name text,
  user_avatar text,
  PRIMARY KEY(room_id, joined_at),
)WITH CLUSTERING ORDER BY (joined_at ASC);
CREATE INDEX user_idx ON shinfo.users_by_room ( user_id );

CREATE TABLE shinfo.messages(
  id timeuuid,
  room_id timeuuid,
  "user" user,
  message_type int,
  entry text,
  metadata map<text, text>,
  PRIMARY KEY(room_id, id),
) WITH CLUSTERING ORDER BY (id DESC);

CREATE TABLE shinfo.rooms(
  id timeuuid,
  "name" text,
  users list<frozen<user>>,
  archived boolean,
  public boolean,
  "type" int,
  created_at timestamp,
  PRIMARY KEY(id),
);
