create table if not exists users (
    id varchar(255) primary key,
    user_name varchar(255),
    first_name varchar(255),
    last_name varchar(255),
    middle_name varchar(255),
    birth_date timestamp,
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
    phone varchar(255),
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
);

create table if not exists relations(
    user_id1 varchar(255),
    user_id2 varchar(255),
    status varchar(255),
    primary key(user_id1, user_id2)
);

create table if not exists refresh_tokens(
    user_id varchar(255),
    device_id varchar(255),
    browser varchar(255),
    ip_address varchar(255),
    refresh_token varchar(255),
    expiry int8,
    created_at timestamp,
    updated_at timestamp,
    primary key(user_id, ip_address, browser, device_id)
);

create table if not exists notifications{
    id varchar(255),
    requestor_id varchar(255),
    subscribers string[],
    is_read boolean,
    content string,
    primary key(id) 
}