CREATE EXTENSION IF NOT EXISTS CITEXT;

CREATE UNLOGGED TABLE users (
    nickname CITEXT  PRIMARY KEY,
    about TEXT, 
    email CITEXT UNIQUE, 
    fullname VARCHAR(100)
);

CREATE UNLOGGED TABLE forums (
    slug     CITEXT        PRIMARY KEY, 
    title    VARCHAR(100) NOT NULL, 
    "user"   CITEXT        REFERENCES users(nickname) NOT NULL,
    posts    BIGINT        DEFAULT 0,
    threads  BIGINT        DEFAULT 0
);

CREATE UNLOGGED TABLE threads (
    id       SERIAL         PRIMARY KEY, 
    author   CITEXT         REFERENCES users(nickname), 
    created  TIMESTAMP WITH TIME ZONE, 
    forum    CITEXT         REFERENCES forums(slug) NOT NULL, 
    message  TEXT, 
    slug     CITEXT         UNIQUE, 
    title    VARCHAR(100), 
    votes    INT            DEFAULT 0
);

CREATE SEQUENCE IF NOT EXISTS posts_id_seq START 1;

CREATE UNLOGGED TABLE posts (
    id          SERIAL      PRIMARY KEY, 
    author      CITEXT      REFERENCES users(nickname), 
    created     TIMESTAMP WITH TIME ZONE   DEFAULT NOW(), 
    forum       CITEXT      REFERENCES forums(slug), 
    isEdited    BOOLEAN     DEFAULT false, 
    message     TEXT, 
    parent      INT         DEFAULT NULL, 
    thread      INT         NOT NULL REFERENCES threads(id),
    pathtopost  INT         ARRAY
);

ALTER SEQUENCE posts_id_seq OWNED BY posts.id;

CREATE UNLOGGED TABLE votes (
    nickname CITEXT REFERENCES users(nickname) NOT NULL, 
    voice    INT                               NOT NULL, 
    thread   INT    REFERENCES threads(id)     NOT NULL
);


CREATE UNLOGGED TABLE IF NOT EXISTS forumusers (
	forum            CITEXT       NOT NULL,
	nickname         CITEXT       NOT NULL
);

ALTER TABLE forumusers
ADD CONSTRAINT unique_forum_user_pair UNIQUE (forum, nickname);

DROP INDEX IF EXISTS idx_users_nickname;
DROP INDEX IF EXISTS idx_users_nickname_email;
DROP INDEX IF EXISTS idx_forums_slug;
DROP INDEX IF EXISTS idx_threads_id;
DROP INDEX IF EXISTS idx_threads_slug;
DROP INDEX IF EXISTS idx_threads_created_forum;
DROP INDEX IF EXISTS idx_posts_id;
DROP INDEX IF EXISTS idx_posts_thread_id;
DROP INDEX IF EXISTS idx_posts_thread_id0;
DROP INDEX IF EXISTS idx_posts_thread_path1_id;
DROP INDEX IF EXISTS idx_posts_thread_path_parent;
DROP INDEX IF EXISTS idx_posts_thread;
DROP INDEX IF EXISTS idx_posts_path_AA;
DROP INDEX IF EXISTS idx_posts_path_AD;
DROP INDEX IF EXISTS idx_posts_path_DA;
DROP INDEX IF EXISTS idx_posts_path_DD;
DROP INDEX IF EXISTS idx_posts_path_desc;
DROP INDEX IF EXISTS idx_posts_paths;
DROP INDEX IF EXISTS idx_posts_thread_path;
DROP INDEX IF EXISTS idx_posts_thread_id_created;
DROP INDEX IF EXISTS idx_votes_thread_nickname;

DROP INDEX IF EXISTS idx_fu_user;
DROP INDEX IF EXISTS idx_fu_forum;


CREATE INDEX IF NOT EXISTS idx_fu_user ON forumusers (forum, nickname);
CREATE INDEX IF NOT EXISTS idx_fu_forum ON forumusers (forum);

CREATE INDEX IF NOT EXISTS idx_users_nickname ON users (nickname);

CREATE INDEX IF NOT EXISTS idx_forums_slug ON forums (slug);

CREATE INDEX IF NOT EXISTS idx_threads_id ON threads (id);
CREATE INDEX IF NOT EXISTS idx_threads_slug ON threads (slug);
CREATE INDEX IF NOT EXISTS idx_threads_forum ON threads (forum);

CREATE INDEX IF NOT EXISTS idx_posts_forum ON posts (forum);
CREATE INDEX IF NOT EXISTS idx_posts_id ON posts (id);
CREATE INDEX IF NOT EXISTS idx_posts_thread_path ON posts (thread, pathtopost);
CREATE INDEX IF NOT EXISTS idx_posts_thread_id ON posts (thread, id);
CREATE INDEX IF NOT EXISTS idx_posts_thread_id0 ON posts (thread, id) WHERE parent = 0;
CREATE INDEX IF NOT EXISTS idx_posts_thread_id_created ON posts (id, created, thread);
CREATE INDEX IF NOT EXISTS idx_posts_thread_path1_id ON posts (thread, (pathtopost[1]), id);

CREATE UNIQUE INDEX IF NOT EXISTS idx_votes_thread_nickname ON votes (thread, nickname);