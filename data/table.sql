create table if not exists users (
    id varchar(255) primary key,
    avatar_url varchar(255),
    user_name varchar(255),
    first_name varchar(255),
    last_name varchar(255),
    middle_name varchar(255),
    birth_date TIMESTAMPTZ,
    created_at TIMESTAMPTZ default CURRENT_TIMESTAMP,
    created_by varchar(255),
    updated_at TIMESTAMPTZ default CURRENT_TIMESTAMP,
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
    created_at TIMESTAMPTZ default CURRENT_TIMESTAMP,
    updated_by varchar(255),
    updated_at TIMESTAMPTZ default CURRENT_TIMESTAMP,
    status varchar(255),
    provider varchar(255),
    version int4 not null default 1
);

create table if not exists rooms(
    id varchar(255) primary key,
    name varchar(255),
    members jsonb[],
    created_at TIMESTAMPTZ  default  CURRENT_TIMESTAMP,
    created_by varchar(255),
    updated_at TIMESTAMPTZ default  CURRENT_TIMESTAMP,
    updated_by varchar(255),
    version int4 not null default 1


);

create table if not exists sequences (
    name varchar(255) primary key,
    sequence_no int4 default 1
);

-- status: friend | none | blocked
create table if not exists friends (
    user_id1 varchar(255),
    user_id2 varchar(255),
    status varchar(255),
    primary key (user_id1, user_id2)
);

create table if not exists friend_requests (
    id varchar(255) primary key,
    requester_id varchar(255),
    requestee_id varchar(255),
    status VARCHAR(255),
    created_at TIMESTAMPTZ,
    created_by varchar(255),
    updated_at TIMESTAMPTZ,
    updated_by varchar(255)
);

create table if not exists refresh_tokens (
    user_id varchar(255),
    device_id varchar(255),
    user_agent varchar(255),
    ip_address varchar(255),
    token varchar(255),
    expiry int8,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    primary key (
        user_id,
        ip_address,
        user_agent,
        device_id
    )
);

create table if not exists notifications(
    id varchar(255),
    requester jsonb,
    subscribers jsonb[],
    title varchar(255),
    content VARCHAR(255),
    type VARCHAR(255),
    url VARCHAR(255),
    created_at TIMESTAMPTZ,
    created_by VARCHAR(255),
    updated_at TIMESTAMPTZ,
    updated_by VARCHAR(255),
    visible BOOLEAN,
    primary key(id) 
);

create table if not exists carriers (
    id VARCHAR(255),
    user_id VARCHAR(255),
    company VARCHAR(255),
    position varchar(255),
    address varchar(255),
    description VARCHAR(255),
    start_time TIMESTAMPTZ,
    end_time TIMESTAMPTZ,
    show_with varchar(255)[],
    not_show_with varchar(255)[],
    primary key(id, user_id)
);

create table if not exists user_follow_count (
    user_id VARCHAR(255),
    count VARCHAR(255)
);

create table if not exists user_followers (
    user_id VARCHAR(255),
    follower_id VARCHAR(255)
);