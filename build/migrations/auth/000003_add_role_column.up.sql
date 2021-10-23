CREATE TYPE auth_role AS ENUM ('user', 'admin');
alter table "auth" add role auth_role default 'user';
