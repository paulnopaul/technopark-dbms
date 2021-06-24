-- drop table if exists users cascade;
create table users
(
    nickname text unique not null,
    fullname text        not null,
    about    text,
    email    text unique not null
);

-- drop table if exists forums cascade;
create table forums
(
    title    text        not null,
    username text        not null,
    slug     text unique not null,
    posts    integer     not null default 0,
    threads  integer     not null default 0,
    constraint user_fk foreign key (username) references users (nickname)
);

-- drop table if exists threads cascade;
create table threads
(
    id      serial unique not null,
    title   text          not null,
    author  text          not null,
    message text          not null,
    votes   integer       not null default 0,
    slug    text,
    created timestamp     not null default now(),
    constraint unique_thread_slug unique (slug),
    constraint user_fk foreign key (author) references users (nickname)
);

-- drop table if exists f_t cascade;
create table f_t
(
    f_slug text    not null,
    t_id   integer not null,
    constraint forum_fk foreign key (f_slug) references forums (slug),
    constraint thread_fk foreign key (t_id) references threads (id),
    constraint unique_forum_thread_pair unique (f_slug, t_id)
);

-- drop table if exists f_u cascade;
create table f_u
(
    f_slug text not null,
    u_nick text not null,
    constraint forum_fk foreign key (f_slug) references forums (slug),
    constraint user_fk foreign key (u_nick) references users (nickname),
    constraint unique_user_forum_pair unique (f_slug, u_nick)
);