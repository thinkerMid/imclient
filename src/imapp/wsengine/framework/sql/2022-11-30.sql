create table business
(
    id bigint unsigned auto_increment comment 'id'
        primary key,
    jid varchar(15) not null comment 'jid',
    profile_tag varchar(24) not null comment 'profile_tag',
    catalog_session_id varchar(32) not null comment 'catalog_session_id'
) comment '商业信息';

create index idx_business_jid
    on business (jid);

