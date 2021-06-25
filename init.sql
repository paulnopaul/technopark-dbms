create extension if not exists citext;

drop table if exists users cascade;
create table users
(
    nickname citext collate "POSIX" primary key not null,
    fullname text                               not null,
    about    text,
    email    citext unique                      not null
);

drop table if exists forums cascade;
create table forums
(
    title    text          not null,
    username citext        not null,
    slug     citext unique not null,
    posts    integer       not null default 0,
    threads  integer       not null default 0,
    foreign key (username) references users (nickname)
);

drop table if exists threads cascade;
create table threads
(
    id      bigserial primary key,
    title   text    not null,
    author  citext  not null,
    message text    not null,
    votes   integer not null default 0,
    slug    citext unique,
    created timestamp with time zone,
    forum   citext  not null,
    foreign key (author) references users (nickname),
    foreign key (forum) references forums (slug)
);

drop table if exists posts cascade;
create unlogged table if not exists posts
(
    id        bigserial primary key,
    parent    bigint  not null,
    author    citext  not null,
    message   text    not null,
    is_edited boolean not null,
    forum     citext  not null,
    thread    bigint,
    created   timestamp with time zone default now(),
    foreign key (author) references users (nickname),
    foreign key (forum) references forums (slug),
    foreign key (thread) references threads (id)
);

drop table if exists f_t cascade;
create table f_t
(
    f_slug citext  not null,
    t_id   integer not null,
    foreign key (f_slug) references forums (slug),
    foreign key (t_id) references threads (id),
    unique (f_slug, t_id)
);

drop table if exists f_u cascade;
create table f_u
(
    f_slug citext not null,
    u_nick citext not null,
    foreign key (f_slug) references forums (slug),
    foreign key (u_nick) references users (nickname),
    unique (f_slug, u_nick)
);