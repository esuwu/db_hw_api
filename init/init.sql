ALTER SYSTEM SET checkpoint_completion_target = '0.9';
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
DROP TABLE IF EXISTS forums CASCADE;
DROP TABLE IF EXISTS threads CASCADE;
DROP TABLE IF EXISTS posts CASCADE;
DROP TABLE IF EXISTS votes CASCADE;

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    nickname CITEXT COLLATE "ucs_basic" NOT NULL UNIQUE,
    fullname TEXT NOT NULL,
    email CITEXT NOT NULL UNIQUE,
    about TEXT
);

CLUSTER users USING idx_users_all;

CREATE INDEX users__nickname ON users(id);

CREATE TABLE forums (
    id SERIAL PRIMARY KEY,
    slug CITEXT NOT NULL UNIQUE,
    posts INTEGER DEFAULT 0,
    threads INTEGER DEFAULT 0,
    title TEXT NOT NULL,
    username CITEXT REFERENCES users (nickname) NOT NULL
);

CREATE INDEX forums__slug ON forums(slug);

CREATE TABLE forum_user (
    forum_slug CITEXT,
    user_id INTEGER,
    PRIMARY KEY (forum_slug, user_id)
);

CREATE INDEX forum_user__slug ON forum_user(forum_slug);

CREATE TABLE threads (
    id SERIAL PRIMARY KEY,
    slug CITEXT UNIQUE,
    username CITEXT REFERENCES users (nickname) NOT NULL,
    forum_slug CITEXT REFERENCES forums (slug) NOT NULL,
    created TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    message TEXT,
    title TEXT NOT NULL,
    votes INTEGER DEFAULT 0
);
CLUSTER threads USING idx_threads_fslugdate;
CREATE INDEX IF NOT EXISTS idx_threads_id ON threads (id);

CREATE INDEX threads__forum_created ON threads(forum_slug, created);

CREATE OR REPLACE FUNCTION insertThread()
RETURNS TRIGGER AS
$BODY$
BEGIN
    INSERT INTO forum_user (forum_slug, user_id)
    VALUES (new.forum_slug, (SELECT id FROM users WHERE nickname = new.username))
    ON CONFLICT DO NOTHING;

    UPDATE forums SET threads = threads + 1
    WHERE slug = new.forum_slug;
    RETURN new;
END;
$BODY$
LANGUAGE plpgsql;

CREATE TRIGGER insertThread
AFTER INSERT
ON threads
FOR EACH ROW
EXECUTE PROCEDURE insertThread();

CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    username CITEXT REFERENCES users (nickname) NOT NULL,
    created TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    isEdited BOOLEAN DEFAULT FALSE,
    message TEXT,
    parent_id INTEGER REFERENCES posts (id),
    thread_id INTEGER REFERENCES threads (id) NOT NULL,
    path INTEGER ARRAY
);

CREATE OR REPLACE FUNCTION insertPost()
RETURNS TRIGGER AS
$BODY$
    BEGIN
        IF new.parent_id IS null THEN
            UPDATE posts SET path = ARRAY[new.id]
            WHERE id = new.id;
        ELSE
            UPDATE posts SET path = array_append((SELECT path FROM posts WHERE id = new.parent_id), new.id)
            WHERE id = new.id;
        END IF;
        RETURN new;
    END;
$BODY$
LANGUAGE plpgsql;

CREATE TRIGGER insertPost
AFTER insert
ON posts
FOR EACH ROW
EXECUTE PROCEDURE insertPost();

CREATE INDEX posts__thread_id_created ON posts(thread_id, id, created);

CREATE TABLE votes (
    id SERIAL PRIMARY KEY,
    username CITEXT REFERENCES users (nickname) NOT NULL,
    thread_id INTEGER REFERENCES threads (id) NOT NULL,
    voice INTEGER,
    CONSTRAINT unique_vote UNIQUE (username, thread_id)
);

CREATE OR REPLACE FUNCTION insertVote()
RETURNS TRIGGER AS
$BODY$
    BEGIN
        UPDATE threads SET votes = votes + new.voice
        WHERE id = new.thread_id;
        RETURN new;
    END;
$BODY$
LANGUAGE plpgsql;

CREATE TRIGGER insertVote
AFTER INSERT
ON votes
FOR EACH ROW
EXECUTE PROCEDURE insertVote();


CREATE OR REPLACE FUNCTION updateVote()
RETURNS TRIGGER AS
$BODY$
    BEGIN
        IF old.voice = -1 AND new.voice = 1 THEN
            UPDATE threads SET votes = votes + 2
            WHERE id = new.thread_id;
        END IF;
        IF old.voice = 1 AND new.voice = -1 THEN
            UPDATE threads SET votes = votes - 2
            WHERE id = new.thread_id;
        END IF;
        RETURN NEW;
    END;
$BODY$
LANGUAGE plpgsql;

CREATE TRIGGER updateVote
AFTER UPDATE
ON votes
FOR EACH ROW
EXECUTE PROCEDURE updateVote();