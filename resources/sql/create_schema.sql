-- create public schema

CREATE SCHEMA IF NOT EXISTS public;

-- Users table

CREATE TABLE IF NOT EXISTS public.users (
    id          serial PRIMARY KEY,
    created_at  timestamp with time zone NOT NULL DEFAULT now(),
    updated_at  timestamp with time zone NOT NULL DEFAULT now(),
    username    varchar(128) NOT NULL UNIQUE CHECK (username <> ''),
    email       text NOT NULL UNIQUE CHECK (email <> ''),
    password    text NOT NULL CHECK(char_length(password) > 0),
    actual_name text NOT NULL CHECK (actual_name <> ''),
    dob         timestamp with time zone NOT NULL
);
GRANT SELECT, INSERT, UPDATE, DELETE ON public.users TO api;
GRANT SELECT, USAGE ON users_id_seq TO api;

-- Followers table

CREATE TABLE IF NOT EXISTS public.followers (
    follower_id integer NOT NULL,
    followee_id integer NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    PRIMARY KEY (follower_id, followee_id),
    FOREIGN KEY (follower_id) REFERENCES users (id)
        ON UPDATE CASCADE
        ON DELETE CASCADE,
    FOREIGN KEY (followee_id) REFERENCES users (id)
        ON UPDATE CASCADE
        ON DELETE CASCADE
);
GRANT SELECT, INSERT, UPDATE, DELETE ON public.followers TO api;

-- Posts table

CREATE TABLE IF NOT EXISTS public.posts (
    id          serial PRIMARY KEY,
    created_at  timestamp with time zone NOT NULL DEFAULT now(),
    author_id   integer NOT NULL,
    contents    text NOT NULL,
    FOREIGN KEY (author_id) REFERENCES users (id)
        ON UPDATE CASCADE
        ON DELETE NO ACTION
);
GRANT SELECT, INSERT, UPDATE, DELETE ON public.posts TO api;
GRANT SELECT, USAGE ON posts_id_seq TO api;

-- Post keks table

CREATE TABLE IF NOT EXISTS public.post_keks (
    id          serial PRIMARY KEY,
    author_id   integer NOT NULL,
    post_id     integer NOT NULL,
    FOREIGN KEY (author_id) REFERENCES users (id)
        ON UPDATE CASCADE
        ON DELETE NO ACTION,
    FOREIGN KEY (post_id) REFERENCES posts (id)
        ON UPDATE CASCADE
        ON DELETE CASCADE
);
GRANT SELECT, INSERT, UPDATE, DELETE ON public.posts TO api;
GRANT SELECT, USAGE ON post_keks_id_seq TO api;

-- Post nos table

CREATE TABLE IF NOT EXISTS public.post_nos (
    id          serial PRIMARY KEY,
    author_id   integer NOT NULL,
    post_id     integer NOT NULL,
    FOREIGN KEY (author_id) REFERENCES users (id)
        ON UPDATE CASCADE
        ON DELETE NO ACTION,
    FOREIGN KEY (post_id) REFERENCES posts (id)
        ON UPDATE CASCADE
        ON DELETE CASCADE
);
GRANT SELECT, INSERT, UPDATE, DELETE ON public.posts TO api;
GRANT SELECT, USAGE ON post_nos_id_seq TO api;

-- Inaccurate but fast counting, unnecessary for now

-- CREATE FUNCTION count_keks(id integer) RETURNS INTEGER AS
-- $func$
-- DECLARE
--     rec record;
--     ROWS INTEGER;
-- BEGIN
--     FOR rec IN EXECUTE 'EXPLAIN SELECT * FROM post_keks WHERE post_id=' || CAST(id AS text) LOOP
--         ROWS := SUBSTRING(rec."QUERY PLAN" FROM ' rows=([[:digit:]]+)');
--         EXIT WHEN ROWS IS NOT NULL;
--     END LOOP;

--     RETURN ROWS;
-- END
-- $func$ LANGUAGE plpgsql;

-- CREATE FUNCTION count_nos(id integer) RETURNS INTEGER AS
-- $func$
-- DECLARE
--     rec record;
--     ROWS INTEGER;
-- BEGIN
--     FOR rec IN EXECUTE 'EXPLAIN SELECT * FROM post_nos WHERE post_id=' || CAST(id AS text) LOOP
--         ROWS := SUBSTRING(rec."QUERY PLAN" FROM ' rows=([[:digit:]]+)');
--         EXIT WHEN ROWS IS NOT NULL;
--     END LOOP;

--     RETURN ROWS;
-- END
-- $func$ LANGUAGE plpgsql;