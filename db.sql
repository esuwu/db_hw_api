CREATE EXTENSION IF NOT EXISTS CITEXT;

DROP TABLE IF EXISTS users, forum, thread, post, vote, forum_users CASCADE;

DROP FUNCTION IF EXISTS thread_insert();

CREATE TABLE users (
  id       SERIAL,

  nickname CITEXT NOT NULL,
  email    CITEXT NOT NULL,

  about    TEXT DEFAULT NULL,
  fullname TEXT   NOT NULL
);

CREATE TABLE forum (
  id        SERIAL PRIMARY KEY,
  slug      CITEXT  NOT NULL,

  title     TEXT    NOT NULL,
  moderator CITEXT  NOT NULL,

  threads   INTEGER NOT NULL DEFAULT 0,
  posts     BIGINT  NOT NULL DEFAULT 0
);

CREATE TABLE thread (
  id          SERIAL PRIMARY KEY,

  slug        CITEXT  DEFAULT NULL,
  title       TEXT    NOT NULL,
  message     TEXT    NOT NULL,

  forum_id    INTEGER NOT NULL,
  forum_slug  CITEXT  NOT NULL,

  user_id     INTEGER,
  user_nick   CITEXT  NOT NULL,

  created     TIMESTAMPTZ,
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

CREATE TABLE post (
  id          SERIAL primary key,

  user_nick   TEXT      NOT NULL,

  message     TEXT      NOT NULL,
  created     TIMESTAMPTZ,

  forum_slug  TEXT      NOT NULL,
  thread_id   INTEGER   NOT NULL,

  parent      INTEGER            DEFAULT 0,
  parents     INT [] NOT NULL,
  main_parent INT    NOT NULL,

  is_edited   BOOLEAN   NOT NULL DEFAULT FALSE
);


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
  nickname TEXT,
  email    TEXT,
  about    TEXT,
  fullname TEXT
);
