create table if not exists post
(
    id         bigserial                           not null
        constraint post_pk primary key,
    user_id    int                                 not null,
    text       varchar(140)                        not null,
    created_at timestamp default current_timestamp not null
);

create index if not exists post_user_id_index on post (user_id);
