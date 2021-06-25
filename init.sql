create extension if not exists citext;

drop table if exists users cascade;
create table users
(
    nickname citext collate "C" primary key not null,
    fullname text                           not null,
    about    text,
    email    citext unique                  not null
);

drop table if exists forums cascade;
create table forums
(
    slug     citext primary key,
    title    text    not null,
    username citext  not null,
    posts    integer not null default 0,
    threads  integer not null default 0,
    foreign key (username) references users (nickname)
);

drop table if exists f_u cascade;
create table f_u
(
    f citext not null,
    u citext not null,
    unique (f, u),
    foreign key (u) references users (nickname),
    foreign key (f) references forums (slug)
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

drop table if exists votes cascade;
create table votes
(
    thread   bigint not null,
    username citext not null,
    voice    int not null,
    unique (thread, username),
    foreign key (thread) references threads(id),
    foreign key (username) references users(nickname)
);


drop table if exists posts cascade;
create unlogged table if not exists posts
(
    id        bigserial primary key,
    parent    bigint,
    author    citext  not null,
    message   text    not null,
    is_edited boolean not null         default false,
    forum     citext  not null,
    thread    bigint,
    created   timestamp with time zone default now(),
    foreign key (author) references users (nickname),
    foreign key (forum) references forums (slug),
    foreign key (thread) references threads (id)
);

create index user_nickname_index on users using hash (nickname);
create index user_email_index on users using hash (email);
create index forum_slug_index on forums using hash (slug);
create index thread_slug_index on threads using hash (slug);
create index thread_forum_index on threads using hash (forum);
create index post_forum_index on posts using hash (forum);
create index fu_forum_index on f_u (f);
create index fu_user_index on f_u (u);
create index votes_index on votes (thread, username);

--- NEW THREAD
create or replace function new_thread_update_count()
    returns trigger as
$$
begin
    --- update thread count
    update forums
    set threads = threads + 1
    where slug = new.forum;
    return null;
end;
$$
    language 'plpgsql';

create or replace function new_thread_add_relation()
    returns trigger as
$$
begin
    --- update forum users
    insert into f_u(f, u)
    values (new.forum, new.author)
    on conflict do nothing;
    return null;
end;
$$
    language 'plpgsql';

drop trigger if exists new_thread_created_count on threads;
create trigger new_thread_created_count
    after insert
    on threads
    for each row
execute procedure new_thread_update_count();

drop trigger if exists new_thread_created_u on threads;
create trigger new_thread_created_u
    after insert
    on threads
    for each row
execute procedure new_thread_add_relation();


--- NEW POST
create or replace function new_post_update_count()
    returns trigger as
$$
begin
    --- update post count
    update forums
    set posts = posts + 1
    where slug = new.forum;
    return null;
end;
$$
    language 'plpgsql';

create or replace function new_post_add_relation()
    returns trigger as
$$
begin
    --- update forum users
    insert into f_u(f, u)
    values (new.forum, new.author)
    on conflict do nothing;
    return null;
end;
$$
    language 'plpgsql';

drop trigger if exists new_post_created_count on posts;
create trigger new_post_created_count
    after insert
    on posts
    for each row
execute procedure new_post_update_count();

drop trigger if exists new_post_created_u on posts;
create trigger new_post_created_u
    after insert
    on posts
    for each row
execute procedure new_post_add_relation();

--- VOTES
create or replace function new_vote_update_thread()
    returns trigger as
$$
begin
    update threads
    set votes = votes + new.voice
    where id = new.thread;
    return null;
end;
$$
    language 'plpgsql';

create or replace function updated_vote_update_thread()
    returns trigger as
$$
begin
    update threads
    set votes = (votes + new.voice - old.voice)
    where id = new.thread;
    return null;
end;
$$
    language 'plpgsql';


CREATE TRIGGER new_vote_set
    AFTER INSERT
    ON votes
    FOR EACH ROW
EXECUTE PROCEDURE new_vote_update_thread();

CREATE TRIGGER vote_updated
    AFTER update
    ON votes
    FOR EACH ROW
EXECUTE PROCEDURE updated_vote_update_thread();


