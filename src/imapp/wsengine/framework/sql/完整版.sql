create table ab_key
(
    id      bigint unsigned auto_increment
        primary key,
    jid     varchar(50)  not null,
    content varchar(128) not null
);

create index idx_ab_key_jid
    on ab_key (jid);

create table account_info
(
    id                 bigint auto_increment comment 'id'
        primary key,
    jid                varchar(15)        not null comment '国家区号和手机号码',
    signature          varchar(139)       not null comment '签名',
    have_avatar        tinyint(1) not null comment '是否有头像',
    status             smallint default 0 not null comment '账号状态',
    logout_time        int                not null comment '登出时间',
    send_event_0_time  int                not null comment '发送渠道0日志时间',
    send_event_0_count int                not null comment '发送渠道0事件次数',
    send_event_2_time  int                not null comment '发送渠道2日志时间',
    send_event_2_count int                not null comment '发送渠道2事件次数',
    login_count        int                not null comment '登录次数',
    app_page           tinyint(1) not null comment '当前所在页面'
) comment '账号信息';

create index idx_account_info_phone_number
    on account_info (jid);

create index idx_account_info_status_logout_time
    on account_info (status, logout_time);

create table aes_key
(
    id      int auto_increment comment 'id'
        primary key,
    jid     varchar(50) not null comment '国家区号和手机号码',
    aes_key varbinary(32) not null comment 'aes_key',
    pub_key varbinary(32) not null comment 'pub_key',
    pri_key varbinary(32) not null comment 'pri_key'
) comment '账号AesKey';

create index idx_aes_key_jid
    on aes_key (jid);

create table client_persistence
(
    id           bigint unsigned auto_increment comment 'id'
        primary key,
    phone_number varchar(15)  not null comment '国际区号和手机号',
    socks_info   varchar(128) not null comment '代理连接信息',
    push_target  tinyint(1) not null comment '推送设置'
) comment '客户端持久化信息';

create index idx_client_persistence_phone_number
    on client_persistence (phone_number);

create table contacts
(
    id               bigint unsigned auto_increment comment 'id'
        primary key,
    src_number       varchar(15) not null comment '国际区号和手机号',
    dst_phone_number varchar(15) not null,
    dst_jid_number   varchar(15) not null comment 'jid号码',
    trusted_contact  tinyint(1) not null,
    add_time         int         not null comment '添加时间',
    have_avatar      tinyint(1) not null comment '是否有头像',
    chat_with        tinyint(1) not null comment '发送过消息',
    receive_chat     tinyint(1) not null comment '接收过消息',
    in_address_book  tinyint(1) not null comment '加过通讯录'
) comment '联系人映射';

create index idx_contacts_src_number_dst_jid_number_index
    on contacts (src_number, dst_jid_number);

create index idx_contacts_src_number_dst_number
    on contacts (src_number, dst_phone_number);

create table device_list
(
    id                       bigint unsigned auto_increment
        primary key,
    our_jid                  varchar(15) not null comment 'JID',
    their_jid                varchar(15) not null comment 'JID',
    device_id                smallint    not null comment '设备ID',
    registration_id          int unsigned not null comment '设备注册的ID',
    identity                 varbinary(32) not null comment '设备身份',
    initialization           tinyint(1) not null comment '初始化标志',
    unacknowledged_state     tinyint(1) not null comment '未确认的会话;pkmsg',
    pending_prekey_id        int         not null comment '密钥ID',
    pending_signed_prekey_id int unsigned not null comment '密钥ID',
    pending_base_key         varbinary(32) not null comment '密钥',
    previous_counter         int         not null comment '计数器',
    receiver_public_key      varbinary(32) not null comment '密钥',
    receiver_private_key     varbinary(32) not null comment '密钥',
    receiver_chain_index     int         not null comment '计数器',
    receiver_chain_key       varbinary(32) not null comment '密钥',
    sender_public_key        varbinary(32) not null comment '密钥',
    sender_private_key       varbinary(32) not null comment '密钥',
    sender_chain_index       int         not null comment '计数器',
    sender_chain_key         varbinary(32) not null comment '密钥',
    sender_base_key          varbinary(32) not null comment '密钥',
    kdf_root_key             varbinary(32) not null comment '密钥',
    session_version          int         not null comment '会话版本'
);

create index idx_device_list_our_jid_their_jid_device_id_initialization
    on device_list (our_jid, their_jid, device_id, initialization);

create table deviceinfo
(
    id                 bigint unsigned auto_increment comment 'id'
        primary key,
    jid                varchar(50) not null comment '国家区号和手机号码',
    registrationId     bigint      not null comment 'registrationId',
    identityKey        varbinary(32) not null comment 'identityKey',
    pushName           varchar(30) not null comment '消息推送名字',
    userName           varchar(30) not null comment '昵称',
    fbuuid             varchar(50) not null comment 'uuid',
    fbuuidCreateTime   int         not null comment 'uuid创建时间',
    uuid               varchar(50) not null comment 'uuid',
    osVersion          varchar(20) not null comment '系统版本',
    mcc                varchar(6)  not null comment 'mcc',
    mnc                varchar(6)  not null comment 'mnc',
    manufacturer       varchar(20) not null comment 'manufacturer',
    device             varchar(30) not null comment 'device',
    language           varchar(5)  not null comment '语言',
    country            varchar(10) not null comment '城市',
    buildNumber        varchar(20) not null comment '系统编译版本号',
    privateStatsId     varchar(50) not null comment 'privateStatsId',
    area               varchar(10) not null comment '地区',
    phone              varchar(13) not null comment '手机号',
    securityCodeSet    tinyint(1) not null comment 'securityCodeSet',
    clientStaticPriKey varbinary(32) not null comment 'clientStaticPriKey',
    clientStaticPubKey varbinary(32) not null comment 'clientStaticPubKey',
    serverStaticKey    varbinary(32) not null comment 'serverStaticKey',
    businessName       varchar(30) not null comment '商家名称',
    platform           varchar(20) not null comment '平台'
) comment '设备信息';

create index idx_deviceinfo_jid
    on deviceinfo (jid);

create table event_cache
(
    id                bigint unsigned auto_increment comment 'id'
        primary key,
    jid               varchar(15) not null comment '国家区号和手机号码',
    serial_number     bigint      not null comment '日志序号',
    auto_increment_id int         not null comment '事件序号',
    channel_id        int         not null comment '渠道号',
    event_log         varbinary(2048) not null comment '日志内容'
) comment '客户端日志缓存';

create index idx_account_event_cache_phone_number
    on event_cache (jid, channel_id, serial_number);

create table `group`
(
    id                   bigint unsigned auto_increment
        primary key,
    jid                  varchar(15) not null comment 'JID',
    group_id             varchar(20) not null comment '群组ID',
    is_admin             tinyint(1) not null comment '自己是不是管理员',
    have_group_icon      tinyint(1) not null comment '是否有头像',
    last_edit_desc_key   varchar(24) not null comment '上次编辑描述的KEY',
    chat_permission      tinyint(1) not null comment '聊天权限',
    edit_desc_permission tinyint(1) not null comment '编辑权限'
);

create index idx_group_jid_group_id
    on `group` (jid, group_id);

create table identity
(
    id       bigint unsigned auto_increment comment 'id'
        primary key,
    ourJid   varchar(50) not null comment 'ourJid',
    theirJid varchar(50) not null comment 'theirJid',
    identity varbinary(32) not null comment 'identity'
) comment 'identity';

create index idx_identity_ourJid_theirJid
    on identity (ourJid, theirJid);

create table prekey
(
    id       bigint unsigned auto_increment comment 'id'
        primary key,
    jid      varchar(50)  not null comment 'jid',
    keyId    int          not null comment 'keyId',
    keyBuff  varchar(300) not null comment 'keyBuff',
    isUpload tinyint(1) not null comment 'isUpload'
) comment 'prekey';

create index idx_prekey_jid
    on prekey (jid);

create table push_config
(
    id         bigint auto_increment
        primary key,
    jid        varchar(15) not null comment '国际区号和手机号',
    voip_token varchar(64) not null comment '推送token',
    apns_token varchar(64) not null comment '推送token',
    pkey       varbinary(32) not null comment 'pkey'
);

create index idx_client_push_config_phone_number
    on push_config (jid);

create table revertoken
(
    id           bigint unsigned auto_increment comment 'id'
        primary key,
    jid          varchar(50) not null comment 'jid',
    backupToken  varbinary(20) not null comment 'backupToken',
    recoverToken varbinary(16) not null comment 'recoverToken',
    backupKey    varbinary(32) not null comment 'backupKey',
    backup_key_2 varbinary(32) not null comment 'backup_key_2',
    isUpload     tinyint(1) not null comment 'isUpload'
) comment 'revertoken';

create index idx_revertoken_jid
    on revertoken (jid);

create table routing_info
(
    id      bigint unsigned auto_increment
        primary key,
    jid     varchar(50) not null,
    content varbinary(4) not null
);

create index idx_routing_info_jid
    on routing_info (jid);

create table senderkey
(
    id               bigint unsigned auto_increment comment 'id'
        primary key,
    our_jid          varchar(50) not null comment '自身JID',
    their_jid        varchar(50) not null comment '对方JID',
    group_id         varchar(50) not null comment '群组ID',
    device_id        smallint    not null comment '设备索引',
    chat_with        tinyint(1) not null comment '发送过消息',
    key_id           bigint unsigned not null comment '密钥',
    iteration        int         not null comment '密钥',
    chain_key        varbinary(32) not null comment '密钥',
    public_sign_key  varbinary(32) not null comment '密钥',
    private_sign_key varbinary(32) not null comment '密钥'
) comment 'senderkey';

create index idx_senderkey_our_jid_group_id_chat_with
    on senderkey (our_jid, group_id, chat_with);

create index idx_senderkey_our_jid_their_jid_group_id_device_id
    on senderkey (our_jid, their_jid, group_id, device_id);

create table signedprekey
(
    id      bigint unsigned auto_increment comment 'id'
        primary key,
    jid     varchar(50)  not null comment 'jid',
    keyId   int          not null comment 'keyId',
    keyBuff varchar(500) not null comment 'keyBuff'
) comment 'signedprekey';

create index idx_signedprekey_jid
    on signedprekey (jid);

