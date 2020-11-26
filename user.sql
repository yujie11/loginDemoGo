create table users
(
    mobile_number varchar(15) not null
        constraint users_pkey
            primary key,
    user_name     varchar(20) not null,
    password      varchar(50) not null
);

comment on table users is 'The user table';

comment on column users.mobile_number is '手机号';

comment on column users.user_name is '用户名';

comment on column users.password is '密码';