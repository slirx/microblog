create table if not exists feed
(
    id      bigserial not null
        constraint feed_pk primary key,
    user_id int       not null,
    post_id bigint    not null
);

create index if not exists feed_user_id_index on feed (user_id);
