package models

//easyjson:json
type Status struct{
	Forum int `json:"forum"`
	Post int `json:"post"`
	Thread int `json:"thread"`
	User int `json:"user"`
}

const InitScript = `ALTER SYSTEM SET checkpoint_completion_target = '0.9';
ALTER SYSTEM SET wal_buffers = '6912kB';
ALTER SYSTEM SET default_statistics_target = '100';
ALTER SYSTEM SET random_page_cost = '1.1';
ALTER SYSTEM SET effective_io_concurrency = '200';
ALTER SYSTEM SET seq_page_cost = '0.1';
ALTER SYSTEM SET random_page_cost = '0.1';


CREATE EXTENSION IF NOT EXISTS CITEXT;

DROP TABLE IF EXISTS users, forum, thread, post, vote, forum_users CASCADE;
DROP FUNCTION IF EXISTS thread_insert();


CREATE TABLE users (
  id       SERIAL,

  nickname CITEXT COLLATE ucs_basic NOT NULL,
  email    CITEXT NOT NULL,

  about    TEXT DEFAULT NULL,
  fullname TEXT   NOT NULL
);

CREATE INDEX users_cover_idx ON users (nickname, email, about, fullname);

CREATE UNIQUE INDEX users_nickname_idx ON users (nickname);

CREATE UNIQUE INDEX users_email_idx ON users (email);

CREATE INDEX ON users (nickname, email);


CREATE TABLE forum (
  id        SERIAL PRIMARY KEY,
  slug      CITEXT  NOT NULL,

  title     TEXT    NOT NULL,
  moderator CITEXT  NOT NULL,

  threads   INTEGER NOT NULL DEFAULT 0,
  posts     BIGINT  NOT NULL DEFAULT 0
);

CREATE UNIQUE INDEX forum_slug_idx ON forum (slug);

CREATE INDEX forum_slug_id_idx ON forum (slug, id);

CREATE INDEX on forum (slug, id, title, moderator, threads, posts);

CREATE TABLE thread (
  id          SERIAL PRIMARY KEY,

  slug        CITEXT  DEFAULT NULL,
  title       TEXT    NOT NULL,
  message     TEXT    NOT NULL,

  forum_id    INTEGER NOT NULL,
  forum_slug  CITEXT  NOT NULL,

  user_id     INTEGER,
  user_nick   CITEXT  NOT NULL,

  created     TIMESTAMP(3) WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
  votes_count INTEGER DEFAULT 0
);


CREATE FUNCTION thread_insert()
  RETURNS TRIGGER AS
$BODY$
BEGIN
  UPDATE forum
  SET
    threads = forum.threads + 1
  WHERE slug = NEW.forum_slug;
  RETURN NULL;
END;
$BODY$
LANGUAGE plpgsql;

CREATE TRIGGER on_thread_insert
  AFTER INSERT
  ON thread
  FOR EACH ROW EXECUTE PROCEDURE thread_insert();

CREATE UNIQUE INDEX thread_slug_idx ON thread (slug);

CREATE INDEX thread_slug_id_idx ON thread (slug, id);

CREATE INDEX thread_forum_id_created_idx ON thread (forum_id, created);

CREATE INDEX thread_forum_id_created_desc_idx
  ON thread (forum_id, created DESC);

CREATE UNIQUE INDEX thread_id_forum_slug_idx
  ON thread (id, forum_slug);

CREATE UNIQUE INDEX thread_slug_forum_slug_idx
  ON thread (slug, forum_slug);

CREATE UNIQUE INDEX thread_cover_idx
  ON thread (forum_id, created, id, slug, title, message, forum_slug, user_nick, created, votes_count);


CREATE TABLE post (
  id          SERIAL primary key,

  user_nick   TEXT      NOT NULL,

  message     TEXT      NOT NULL,
  created     TIMESTAMP(3) WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,

  forum_slug  TEXT      NOT NULL,
  thread_id   INTEGER   NOT NULL,

  parent      INTEGER            DEFAULT 0,
  parents     INT [] NOT NULL,
  main_parent INT    NOT NULL,

  is_edited   BOOLEAN   NOT NULL DEFAULT FALSE
);



CREATE INDEX posts_thread_id_id_idx ON post (thread_id, id);

CREATE INDEX posts_thread_id_idx ON post (thread_id);

CREATE INDEX posts_thread_id_parents_idx ON post (thread_id, parents);

CREATE INDEX ON post (thread_id, id, parent, main_parent) WHERE parent = 0;

CREATE INDEX parent_tree_3_1_idx ON post (main_parent, parents DESC, id);

CREATE INDEX parent_tree_4_idx ON post (id, main_parent);

CREATE TABLE vote (
  id         SERIAL,

  user_id    INTEGER NOT NULL,
  thread_id  INTEGER NOT NULL REFERENCES thread,

  voice      INTEGER,
  prev_voice INTEGER DEFAULT 0,
  CONSTRAINT unique_user_and_thread UNIQUE (user_id, thread_id)
);


CREATE TABLE forum_users (
  forumId  INTEGER,
  nickname CITEXT COLLATE ucs_basic NOT NULL,
  email    TEXT,
  about    TEXT,
  fullname TEXT
);

CREATE UNIQUE INDEX forum_users_forum_id_nickname_idx2 ON forum_users (forumId, lower(nickname));

CREATE INDEX forum_users_cover_idx2 ON forum_users (forumId, lower(nickname), nickname, email, about, fullname);`
