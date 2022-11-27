CREATE TABLE IF NOT EXISTS `iam_user`
(
    `id`              bigint AUTO_INCREMENT COMMENT '自增主键' PRIMARY KEY,
    `user_id`         varchar(64)                           NOT NULL COMMENT '用户ID',
    `tenant_id`       varchar(64)                           NOT NULL COMMENT '租户ID',
    `username`        varchar(360)                          NOT NULL COMMENT '登录账号',
    `nick_name`       varchar(360)                          NOT NULL COMMENT '昵称别名',
    `pwd`             varchar(200)                          NOT NULL COMMENT '密码',
    `avatar`          varchar(1000)                         NULL COMMENT '头像',
    `birthday`        datetime                              NULL COMMENT '生日',
    `gender`          tinyint                               NULL COMMENT '性别',
    `email`           varchar(64)                           NULL COMMENT '邮箱',
    `phone`           varchar(64)                           NULL COMMENT '电话',
    `user_type`       tinyint                               NOT NULL COMMENT '用户类型',
    `user_state`      tinyint                               NOT NULL COMMENT '用户状态',
    `def_pwd`         tinyint                               NOT NULL COMMENT '默认密码',
    `country`         varchar(64) DEFAULT ''                NOT NULL COMMENT '国家',
    `last_login_ip`   varchar(64)                           NULL COMMENT '最后登录IP',
    `last_login_time` datetime                              NULL COMMENT '最后登录时间',
    `infos`           varchar(2048)                         NULL COMMENT '用户扩展信息',
    `create_time`     datetime    DEFAULT CURRENT_TIMESTAMP NOT NULL COMMENT '创建时间',
    `created_by`      varchar(64)                           NULL COMMENT '创建者',
    `update_time`     datetime    DEFAULT CURRENT_TIMESTAMP NOT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    `updated_by`      varchar(64)                           NULL COMMENT '修改者',
    constraint iam_user_tenant_id_username_uindex unique (`tenant_id`, `username`),
    constraint iam_user_user_id_uindex unique (`user_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT ='用户表';

create index iam_user_username_index
    on iam_user (`username`);

CREATE TABLE IF NOT EXISTS `iam_group`
(
    `id`          bigint AUTO_INCREMENT COMMENT '自增主键'
        PRIMARY KEY,
    `group_id`    varchar(64)                        NOT NULL COMMENT '用户ID',
    `tenant_id`   varchar(64)                        NOT NULL COMMENT '租户ID',
    `group_name`  varchar(320)                       NOT NULL COMMENT '组名称',
    `create_time` datetime DEFAULT CURRENT_TIMESTAMP NOT NULL COMMENT '创建时间',
    `created_by`  varchar(64)                        NULL COMMENT '创建者',
    `update_time` datetime DEFAULT CURRENT_TIMESTAMP NOT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    `updated_by`  varchar(64)                        NULL COMMENT '修改者',
    constraint iam_group_tenant_id_group_name_uindex
        unique (`tenant_id`, `group_name`),
    constraint iam_group_group_id_uindex
        unique (`group_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT ='用户组表';

create index iam_group_group_name_index
    on iam_group (`group_name`);

CREATE TABLE IF NOT EXISTS `iam_user_group`
(
    `id`          bigint AUTO_INCREMENT COMMENT '自增主键'
        PRIMARY KEY,
    `user_id`     varchar(64)                        NOT NULL COMMENT '用户ID',
    `group_id`    varchar(64)                        NOT NULL COMMENT '组ID',
    `tenant_id`   varchar(64)                        NOT NULL COMMENT '租户ID',
    `create_time` datetime DEFAULT CURRENT_TIMESTAMP NOT NULL COMMENT '创建时间',
    `created_by`  varchar(64)                        NULL COMMENT '创建者',
    `update_time` datetime DEFAULT CURRENT_TIMESTAMP NOT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    `updated_by`  varchar(64)                        NULL COMMENT '修改者',
    constraint iam_user_group_user_id_group_id_uindex
        unique (`user_id`, `group_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT ='用户用户组关联表';

create index iam_user_group_tenant_id_user_id_index
    on iam_user_group (`tenant_id`, `user_id`);
create index iam_user_group_tenant_id_group_id_index
    on iam_user_group (`tenant_id`, `group_id`);

CREATE TABLE IF NOT EXISTS `iam_tenant`
(
    `id`          bigint AUTO_INCREMENT COMMENT '自增主键'
        PRIMARY KEY,
    `tenant_id`   varchar(64)                        NOT NULL COMMENT '租户ID(全局唯一)',
    `tenant_name` varchar(380)                       NOT NULL COMMENT '租户名',
    `owner_id`    varchar(64)                        NOT NULL COMMENT '租户所有者ID',
    `username`    varchar(360)                       NOT NULL COMMENT '租户所有者名字/邮箱',
    `annotate`    varchar(500)                       NULL COMMENT '描述',
    `create_time` datetime DEFAULT CURRENT_TIMESTAMP NOT NULL COMMENT '创建时间',
    `created_by`  varchar(64)                        NULL COMMENT '创建者',
    `update_time` datetime DEFAULT CURRENT_TIMESTAMP NOT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    `updated_by`  varchar(64)                        NULL COMMENT '修改者',
    constraint iam_tenant_tenant_id_owner_id_uindex
        unique (`tenant_id`),
    constraint iam_tenant_tenant_name_uindex
        unique (`tenant_name`),
    constraint iam_tenant_username_uindex
        unique (`username`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT ='租户信息表';
