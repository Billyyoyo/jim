
CREATE TABLE `ack` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `msg_id` bigint(20) NOT NULL COMMENT '消息id',
  `send_count` int(11) NOT NULL DEFAULT '0' COMMENT '发送目标数量, 在线链接端数',
  `arrive_count` int(11) NOT NULL DEFAULT '0' COMMENT '送达数量',
  PRIMARY KEY (`id`),
  KEY `ack_message_id_IDX` (`msg_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=129 DEFAULT CHARSET=utf8 COMMENT='所有发往客户端的消息记录发送和收到ack作为统计，由于没有离线消息的单独存储，所以ack并不具备更多的功能'


CREATE TABLE `device` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `user_id` bigint(20) NOT NULL COMMENT '用户id',
  `serial_no` varchar(100) NOT NULL COMMENT '设备号',
  `create_time` bigint(20) NOT NULL COMMENT '创建时间',
  `last_address` varchar(100) NOT NULL COMMENT '接入地址',
  `last_conn_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '最后接入时间',
  `last_discon_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '最后断开时间',
  `last_sequence` bigint(20) NOT NULL DEFAULT '0' COMMENT '最后送达的消息序列号',
  PRIMARY KEY (`id`),
  KEY `device_user_id_IDX` (`user_id`) USING BTREE,
  KEY `device_last_conn_time_IDX` (`last_conn_time`) USING BTREE,
  KEY `device_serial_IDX` (`serial_no`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=21 DEFAULT CHARSET=utf8


CREATE TABLE `member` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `session_id` bigint(20) NOT NULL COMMENT '会话id',
  `user_id` bigint(20) NOT NULL COMMENT '用户id',
  `create_time` bigint(20) NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `user_session_session_id_IDX` (`session_id`) USING BTREE,
  KEY `user_session_user_id_IDX` (`user_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=31 DEFAULT CHARSET=utf8


CREATE TABLE `message` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `sender_id` bigint(20) NOT NULL COMMENT '发送者id',
  `session_id` bigint(20) NOT NULL COMMENT '会话id',
  `type` int(11) NOT NULL COMMENT '类型',
  `status` int(11) NOT NULL DEFAULT '1' COMMENT '状态   1正常  2撤回',
  `create_time` bigint(20) NOT NULL COMMENT '发送时间',
  `sequence` bigint(20) NOT NULL COMMENT '按接收者Id的消息顺序号',
  `receptor_id` bigint(20) NOT NULL COMMENT '接收者',
  `body` varbinary(1024) NOT NULL COMMENT '消息体',
  `send_no` bigint(20) NOT NULL DEFAULT '0' COMMENT '按发送者id的消息顺序号',
  PRIMARY KEY (`id`),
  KEY `message_session_id_IDX` (`session_id`,`create_time`) USING BTREE,
  KEY `message_sequnce_IDX` (`sequence`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=151 DEFAULT CHARSET=utf8 COMMENT='采用写扩散，每个用户保存自己收到的消息   包括群消息，当退出群聊，消息保留但置空，保证sequnce的连续性'


CREATE TABLE `session` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `name` varchar(256) NOT NULL COMMENT '会话名',
  `type` int(11) NOT NULL COMMENT '1对话, 2群聊',
  `create_time` bigint(20) NOT NULL COMMENT '创建时间',
  `owner` bigint(20) NOT NULL COMMENT '创建者',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=17 DEFAULT CHARSET=utf8


CREATE TABLE `user` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `name` varchar(100) NOT NULL COMMENT '用户名',
  `create_time` bigint(20) NOT NULL COMMENT '创建时间',
  `password` varchar(100) NOT NULL COMMENT '密码',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8

