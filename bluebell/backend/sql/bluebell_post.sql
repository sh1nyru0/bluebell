create table post
(
    id           bigint auto_increment
        primary key,
    post_id      bigint                              not null comment '帖子id',
    title        varchar(128)                        not null comment '标题',
    content      varchar(8192)                       not null comment '内容',
    author_id    bigint                              not null comment '作者的用户id',
    community_id bigint                              not null comment '所属社区',
    status       tinyint   default 1                 not null comment '帖子状态',
    create_time  timestamp default CURRENT_TIMESTAMP null comment '创建时间',
    update_time  timestamp default CURRENT_TIMESTAMP null on update CURRENT_TIMESTAMP comment '更新时间',
    constraint idx_post_id
        unique (post_id)
)
    collate = utf8mb4_general_ci;

create index idx_author_id
    on post (author_id);

create index idx_community_id
    on post (community_id);