-- "name" is stored as varchar to support variable length string with 
-- no upper limit. "password" is stored as bytea to support the bcrypt implementation. 
-- The Go implementation outputs the hash in base64 format which makes varchar an alternative
-- to bytea in our case. I still choose to go with bytea to be verbose about the fact that we
-- are dealing with binary data here. "email" is stored in accordance with RFC 5321 and Errata
-- 1690.
create table if not exists users(
    id uuid primary key,
    name varchar not null,
    email varchar(319) not null,
    password bytea not null
);

-- "startTime" and "endTime" are supposed to be 64-bit epoch time values.
-- Choices for other data types are self explanatory.
create table if not exists logs(
    id uuid primary key,
    userId uuid references users(id),
    latitude double precision,
    longitude double precision,
    activity varchar not null,
    startTime bigint not null,
    endTime bigint not null,
    notes varchar
);
