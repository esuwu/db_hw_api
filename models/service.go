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

ALTER SYSTEM SET max_worker_processes = '4';
ALTER SYSTEM SET max_parallel_workers_per_gather = '2';
ALTER SYSTEM SET max_parallel_workers = '4';
ALTER SYSTEM SET max_parallel_maintenance_workers = '2';


CREATE EXTENSION IF NOT EXISTS CITEXT;

DROP TABLE IF EXISTS users CASCADE;
CREATE TABLE users
(
    ID       BIGSERIAL NOT NULL PRIMARY KEY,
    nickname CITEXT,
    about    TEXT,
    email    CITEXT UNIQUE,
    fullname TEXT
);

CREATE INDEX users_covering_index ON users (nickname, email, about, fullname);

CREATE UNIQUE INDEX users_nickname_index ON users (nickname);

CREATE UNIQUE INDEX users_email_index ON users (email);

CREATE INDEX ON users (nickname, email);



DROP TABLE IF EXISTS forums CASCADE;
CREATE TABLE forums
(
    ID       BIGSERIAL NOT NULL PRIMARY KEY,
    slug     CITEXT    NOT NULL UNIQUE,
    title    TEXT      NOT NULL,
    authorID BIGINT    NOT NULL,
    FOREIGN KEY (authorID) REFERENCES users (ID) ON DELETE CASCADE
);

CREATE UNIQUE INDEX forum_slug_index ON forums (slug);

CREATE INDEX forum_slug_id_index ON forums (slug, ID);

CREATE INDEX on forums (slug, ID, title, authorID);


DROP TABLE IF EXISTS threads CASCADE;
CREATE TABLE threads
(
    ID       BIGSERIAL NOT NULL PRIMARY KEY,
    created  TIMESTAMP WITH TIME ZONE,
    forumID  BIGINT    NOT NULL,
    message  TEXT,

    slug     CITEXT UNIQUE DEFAULT NULL,
    title    TEXT,
    vote     INTEGER       DEFAULT 0,

    authorID BIGINT    NOT NULL,
    FOREIGN KEY (authorID) REFERENCES users (ID) ON DELETE CASCADE,
    FOREIGN KEY (forumID) REFERENCES forums (ID) ON DELETE CASCADE
);

CREATE UNIQUE INDEX thread_slug_index
  ON threads (slug);

CREATE INDEX thread_slug_id_index
  ON threads (slug, ID);

CREATE INDEX thread_forum_id_created_index
  ON threads (forumID, created);

CREATE INDEX thread_forum_id_created_index2
  ON threads (forumID, created DESC);

CREATE UNIQUE INDEX thread_covering_index
  ON threads (forumID, created, ID, slug, title, message, created, vote);


DROP TABLE IF EXISTS posts CASCADE;
CREATE TABLE posts
(
    ID       BIGSERIAL NOT NULL PRIMARY KEY,
    created  TIMESTAMP WITH TIME ZONE,
    forumID  BIGINT    NOT NULL,
    isEdited BOOLEAN,
    message  TEXT,
    parentID BIGINT DEFAULT 0,
    parents BIGINT[] NOT NULL,

    authorID BIGINT    NOT NULL,
    threadID BIGINT    NOT NULL,
    FOREIGN KEY (authorID) REFERENCES users (ID) ON DELETE CASCADE,
    FOREIGN KEY (threadID) REFERENCES threads (ID) ON DELETE CASCADE,
    FOREIGN KEY (forumID) REFERENCES forums (ID) ON DELETE CASCADE
    --FOREIGN KEY (parentID) REFERENCES posts (ID) ON DELETE CASCADE
);

CREATE INDEX idx_messages_tid_mid ON posts (threadID, ID);
CREATE INDEX idx_messages_parent_tree_tid_parent ON posts (threadID, ID) WHERE parentID = 0;
CREATE INDEX idx_messages_all ON posts (ID, created, message, isEdited, parentID, threadID);

CREATE INDEX posts_thread_id_index2 ON posts (threadID);

CREATE INDEX posts_thread_id_parents_index ON posts (threadID, parents);



DROP TABLE IF EXISTS votes CASCADE;
CREATE TABLE votes
(
    ID       BIGSERIAL NOT NULL PRIMARY KEY,
    voice    BOOLEAN,
    threadID BIGINT    NOT NULL,
    authorID BIGINT    NOT NULL,
    FOREIGN KEY (authorID) REFERENCES users (ID) ON DELETE CASCADE,
    FOREIGN KEY (threadID) REFERENCES threads (ID) ON DELETE CASCADE
);`
