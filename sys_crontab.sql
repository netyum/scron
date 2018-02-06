/*
 Navicat Premium Data Transfer

 Source Server         : localhost
 Source Server Type    : MySQL
 Source Server Version : 50716
 Source Host           : localhost
 Source Database       : sys

 Target Server Type    : MySQL
 Target Server Version : 50716
 File Encoding         : utf-8

 Date: 02/06/2018 10:05:36 AM
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
--  Table structure for `sys_crontab`
-- ----------------------------
DROP TABLE IF EXISTS `sys_crontab`;
CREATE TABLE `sys_crontab` (
  `cronId` int(11) NOT NULL AUTO_INCREMENT,
  `task` varchar(255) COLLATE utf8_unicode_ci NOT NULL COMMENT '任务名',
  `active` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0可用 1不可用',
  `mhdmd` varchar(255) COLLATE utf8_unicode_ci NOT NULL COMMENT '时间',
  `command` varchar(255) COLLATE utf8_unicode_ci NOT NULL COMMENT '命令',
  `params` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '参数',
  `process` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '1' COMMENT '进程数',
  `isQueue` tinyint(2) NOT NULL DEFAULT '0' COMMENT '1是队列 0不是',
  `runAt` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '0' COMMENT '运行时间',
  `host` varchar(15) COLLATE utf8_unicode_ci NOT NULL COMMENT '主机IP',
  `logFile` varchar(255) COLLATE utf8_unicode_ci NOT NULL COMMENT '日志文件名',
  `timeout` int(11) NOT NULL DEFAULT '0' COMMENT '分钟',
  `user` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '用户',
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `errorLogUpdatedSize` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '0' COMMENT '错误日志上次大小',
  `LogUpdatedSize` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '0' COMMENT '日志上次更新大小',
  `receiverPhone` text COLLATE utf8_unicode_ci COMMENT '接收警告的手机',
  `receiverWx` text COLLATE utf8_unicode_ci COMMENT '接收警告的WX',
  PRIMARY KEY (`cronId`)
) ENGINE=InnoDB AUTO_INCREMENT=105 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='定义任务表';

-- ----------------------------
--  Table structure for `sys_crontab_host`
-- ----------------------------
DROP TABLE IF EXISTS `sys_crontab_host`;
CREATE TABLE `sys_crontab_host` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `host_name` varchar(255) COLLATE utf8_unicode_ci NOT NULL COMMENT '主机名',
  `host` varchar(15) COLLATE utf8_unicode_ci NOT NULL COMMENT 'IP',
  `is_enable` tinyint(3) NOT NULL DEFAULT '0' COMMENT '是否配置',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='定时任务主机表';

SET FOREIGN_KEY_CHECKS = 1;
