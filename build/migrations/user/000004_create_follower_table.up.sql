create table follower
(
    id          bigserial not null
        constraint follower_pk primary key,
    user_id     int       not null,
    follower_id int       not null
);
