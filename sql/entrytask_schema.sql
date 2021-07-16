-- database 'entrytask'
CREATE TABLE `t_user` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `username` varchar(64) NOT NULL DEFAULT '' COMMENT '用户名',
  `nickname` varchar(128) NOT NULL DEFAULT '' COMMENT '昵称',
  `passwd` varchar(256) NOT NULL DEFAULT '*6BB4837EB74329105EE4568DDA7DC67ED2CA2AD9' COMMENT '登录密码,默认CONCAT(''*'', UPPER(SHA1(UNHEX(SHA1(''123456'')))))',
  `profile_path` varchar(256) NOT NULL DEFAULT '/profile/default.jpg' COMMENT '用户头像图片地址',
  `created_timestamp` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `modified_timestamp` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '最近更新时间',
  `state` tinyint NOT NULL DEFAULT '1' COMMENT '逻辑状态(1正常，2冻结，0已删除)',
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_idx_username` (`username`)
) ENGINE=InnoDB AUTO_INCREMENT=28 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户信息表' ;