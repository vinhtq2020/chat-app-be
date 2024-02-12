create table if not exists users (
    id varchar(255) primary key,
    user_name varchar(255),
    first_name varchar(255),
    last_name varchar(255),
    midde_name varchar(255),
    created_at timestamp default  CURRENT_TIMESTAMP,
    created_by varchar(255),
    updated_at timestamp default  CURRENT_TIMESTAMP,
    updated_by varchar(255),
    version int4 not null default 1
);

create table if not exists users_login_data (
    id varchar(255) primary key, 
    user_name varchar(255),
    email varchar(255),
    password_hash varchar(255),
    created_by varchar(255),
    created_at timestamp default CURRENT_TIMESTAMP,
    updated_by varchar(255),
    updated_at timestamp default CURRENT_TIMESTAMP,
    status varchar(255),
    provider varchar(255),
    version int4 not null default 1
);

create table if not exists rooms(
    id varchar(255) primary key,
    name varchar(255),
    members jsonb[],
    created_at timestamp  default  CURRENT_TIMESTAMP,
    created_by varchar(255),
    updated_at timestamp default  CURRENT_TIMESTAMP,
    updated_by varchar(255),
    version int4 not null default 1

);

create table if not exists sequences(
    name varchar(255) primary key,
    sequence_no int4  default 1
)

create table if not exists relations(
    userId1 varchar(255) primary key,
    userId2 varchar(255) primary key,
    status varchar(255) primary key
)