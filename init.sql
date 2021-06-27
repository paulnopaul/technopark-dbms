create extension if not exists citext;

drop table if exists users cascade;
create unlogged table users
(
    id       bigserial                      not null,
    nickname citext collate "C" primary key not null,
    fullname text                           not null,
    about    text,
    email    citext unique                  not null
);

drop table if exists forums cascade;
create unlogged table forums
(
    id       bigserial not null,
    slug     citext primary key,
    title    text      not null,
    username citext    not null,
    posts    integer   not null default 0,
    threads  integer   not null default 0,
    foreign key (username) references users (nickname)
);

drop table if exists f_u cascade;
create unlogged table f_u
(
    f citext             not null,
    u citext collate "C" not null,
    unique (f, u),
    foreign key (u) references users (nickname),
    foreign key (f) references forums (slug)
);

drop table if exists threads cascade;
create unlogged table threads
(
    id      bigserial primary key,
    title   text    not null,
    author  citext  not null,
    message text    not null,
    votes   integer not null         default 0,
    slug    citext unique,
    created timestamp with time zone default now(),
    forum   citext  not null,
    foreign key (author) references users (nickname),
    foreign key (forum) references forums (slug)
);

drop table if exists votes cascade;
create unlogged table votes
(
    thread   bigint not null,
    username citext not null,
    voice    int    not null,
    unique (thread, username),
    foreign key (thread) references threads (id),
    foreign key (username) references users (nickname)
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
    way       bigint[],
    foreign key (author) references users (nickname),
    foreign key (forum) references forums (slug),
    foreign key (thread) references threads (id)
);

create index user_nickname_index on users using hash (nickname);
create index user_email_index on users using hash (email);

create index forum_slug_index on forums using hash (slug);

create index thread_slug_index on threads using hash (slug);
create index thread_forum_index on threads using hash (forum);
create index thread_fcreated_index on threads (forum, created);

create index fu_user_index on f_u using hash (u);

create index votes_index on votes (thread, username);

create index post_forum_index on posts (forum);
create index post_user_index on posts (author);
create index posts_way_index on posts (way);
create index posts_way_second_index on posts ((way[2]));

--- NEW THREAD
create or replace function new_thread_update()
    returns trigger as
$$
begin
    insert into f_u(f, u)
    select new.forum, new.author
    on conflict do nothing;
    update forums
    set threads = threads + 1
    where slug = new.forum;
    return null;
end;
$$
    language 'plpgsql';

drop trigger if exists new_thread_created on threads;
create trigger new_thread_created
    after insert
    on threads
    for each row
execute procedure new_thread_update();

--- NEW POST

create or replace function new_post_update()
    returns trigger as
$$
begin
    insert into f_u(f, u)
    select new.forum, new.author
    on conflict do nothing;
    update forums
    set posts = posts + 1
    where slug = new.forum;
    return null;
end;
$$
    language 'plpgsql';

drop trigger if exists new_post_created on posts;
create trigger new_post_created
    after insert
    on posts
    for each row
execute procedure new_post_update();

create or replace function update_post_ways()
    returns trigger as
$$
declare
    parent_way    bigint[];
    parent_thread bigint;
begin
    if (new.parent = 0) then
        new.way = array [0,new.id];
    else
        select p.way, p.thread
        from posts p
        where p.id = new.parent
        into parent_way, parent_thread;
        if parent_thread != new.thread or parent_thread is null then
            raise exception using errcode = '66666';
        end if;
        new.way := parent_way || new.id;
    end if;
    return new;
end;
$$
    language 'plpgsql';

drop trigger if exists add_path on posts;
create trigger add_path
    before insert
    on posts
    for each row
execute procedure update_post_ways();

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

drop trigger if exists new_vote_set on votes;
create trigger new_vote_set
    after insert
    on votes
    for each row
execute procedure new_vote_update_thread();


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

drop trigger if exists vote_updated on votes;
create trigger vote_updated
    after update
    on votes
    for each row
execute procedure updated_vote_update_thread();
