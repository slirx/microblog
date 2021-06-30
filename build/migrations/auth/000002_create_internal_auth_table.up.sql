create table if not exists internal_auth
(
    id           serial                              not null
        constraint internal_auth_pk
            primary key,
    service_name varchar(100)                        not null,
    password     varchar(60)                         not null,
    created_at   timestamp default current_timestamp not null
);
create unique index if not exists internal_auth_service_name_uindex on internal_auth (service_name);
