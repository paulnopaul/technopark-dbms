-- \set strLen 100;

drop table if exists forums cascade;
create table forums
(
    title    varchar(:strLen) not null,
    username varchar(:strLen) not null,
    slug     varchar(:strLen) not null,
    posts    integer          not null default 0,
    threads  integer          not null default 0,
    constraint unique_forum_slug unique (slug)
);

drop table if exists threads cascade;
create table threads
(
    id      serial           not null,
    title   varchar(:strLen) not null,
    author  varchar(:strLen) not null,
    message text             not null,
    votes   integer          not null default 0,
    slug    varchar(:strLen),
    created timestamp        not null default now(),
    constraint unique_thread_slug unique (slug),
    constraint thread_pk_id primary key (id)
);

drop table if exists f_t cascade;
create table f_t
(
    f_slug varchar(:strLen) not null,
    t_id   integer          not null,
    constraint forum_fk foreign key (f_slug) references forums (slug),
    constraint thread_fk foreign key (t_id) references threads (id),
    constraint unique_forum_thread_pair unique (f_slug, t_id)
);

drop table if exists users cascade;
create table users
(
    nickname varchar(:strLen) not null,
    fullname varchar(:strLen) not null,
    about    text,
    email    varchar(:strLen) not null,
    constraint user_pk_nickname primary key (nickname);
);

drop table if exists f_u cascade;
create table f_u
(
    f_slug varchar(:strLen) not null,
    u_nick varchar(:strLen) not null,
    constraint forum_fk foreign key (f_slug) references forums (slug),
    constraint user_fk foreign key (u_nick) references users (nickname),
    constraint unique_user_forum_pair unique (f_slug, u_nick)
);

-- fill tables
-- \i fill.sql
