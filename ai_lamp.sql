/*
 Navicat Premium Data Transfer

 Source Server         : 121.42.143.130
 Source Server Type    : MySQL
 Source Server Version : 100117
 Source Host           : 121.42.143.130:3306
 Source Schema         : ai_lamp

 Target Server Type    : MySQL
 Target Server Version : 100117
 File Encoding         : 65001

 Date: 11/08/2019 08:31:54
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for device_group_info
-- ----------------------------
DROP TABLE IF EXISTS `device_group_info`;
CREATE TABLE `device_group_info` (
  `group_id` int(11) NOT NULL AUTO_INCREMENT,
  `group_name` varchar(255) DEFAULT NULL,
  `create_by` varchar(50) DEFAULT NULL,
  `create_date` date DEFAULT NULL,
  `update_by` varchar(50) DEFAULT NULL,
  `update_date` date DEFAULT NULL,
  `delete_status` tinyint(4) DEFAULT NULL,
  PRIMARY KEY (`group_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for device_install_info
-- ----------------------------
DROP TABLE IF EXISTS `device_install_info`;
CREATE TABLE `device_install_info` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `serial_number` varchar(100) DEFAULT NULL,
  `device_address` varchar(255) DEFAULT NULL,
  `lat` varchar(255) DEFAULT NULL,
  `lon` varchar(255) DEFAULT NULL,
  `create_date` datetime DEFAULT NULL,
  `create_by` varchar(100) DEFAULT NULL,
  `update_date` datetime DEFAULT NULL,
  `update_by` varchar(100) DEFAULT NULL,
  `delete_status` tinyint(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for device_manage_info
-- ----------------------------
DROP TABLE IF EXISTS `device_manage_info`;
CREATE TABLE `device_manage_info` (
  `device_id` int(11) NOT NULL AUTO_INCREMENT COMMENT '设备序列号',
  `serial_number` varchar(80) NOT NULL,
  `device_name` varchar(300) DEFAULT NULL COMMENT '设备名称',
  `dhcp_flag` varchar(255) DEFAULT NULL COMMENT '是否使用DHCP',
  `device_ip` varchar(255) DEFAULT NULL COMMENT '设备IP',
  `device_port` varchar(255) DEFAULT NULL COMMENT '设备监听端口',
  `device_status` varchar(255) DEFAULT NULL COMMENT '设备状态',
  `device_attribute` varchar(255) DEFAULT NULL COMMENT '设备属性',
  `device_power` varchar(255) DEFAULT NULL COMMENT '设备功率',
  `device_light_threshold` varchar(255) DEFAULT NULL COMMENT '光敏阀值',
  `device_brightness` varchar(255) DEFAULT NULL COMMENT '亮度',
  `device_temperature` varchar(255) DEFAULT NULL COMMENT '设备温度',
  `auto_start_time` varchar(20) DEFAULT NULL,
  `auto_end_time` varchar(20) DEFAULT NULL,
  `light_mode` varchar(255) DEFAULT NULL COMMENT '补光模式',
  `flash_duration` int(255) DEFAULT NULL COMMENT '爆闪时间（ms）',
  `mac_address` varchar(255) DEFAULT NULL COMMENT 'MAC地址',
  `firmware_version` varchar(255) DEFAULT NULL,
  `longitude` varchar(255) DEFAULT NULL,
  `latitude` varchar(255) DEFAULT NULL,
  `mask` varchar(255) DEFAULT NULL,
  `gateway` varchar(255) DEFAULT NULL COMMENT '网关',
  `pin` varchar(255) DEFAULT NULL,
  `create_by` varchar(50) DEFAULT NULL,
  `create_date` datetime DEFAULT NULL,
  `update_by` varchar(50) DEFAULT NULL,
  `update_data` datetime DEFAULT NULL,
  `delete_status` varchar(255) DEFAULT NULL COMMENT '删除标识。0:正常，1:删除',
  `power_total` int(255) DEFAULT NULL,
  `strobe_count` int(255) DEFAULT NULL,
  PRIMARY KEY (`device_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for device_scan_info
-- ----------------------------
DROP TABLE IF EXISTS `device_scan_info`;
CREATE TABLE `device_scan_info` (
  `serial_number` varchar(200) NOT NULL,
  `refresh_time` int(11) DEFAULT NULL,
  `online_status` tinyint(4) DEFAULT NULL COMMENT '1:在线，0:不在线',
  `firmware_version` varchar(255) DEFAULT NULL,
  `device_ip` varchar(255) DEFAULT NULL,
  `mask` varchar(255) DEFAULT NULL,
  `gateway_addr` varchar(255) DEFAULT NULL,
  `device_port` varchar(255) DEFAULT NULL,
  `mac_addr` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`serial_number`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for event_alarm_info
-- ----------------------------
DROP TABLE IF EXISTS `event_alarm_info`;
CREATE TABLE `event_alarm_info` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `event_type_cd` varchar(255) DEFAULT NULL COMMENT '1:温度异常，2:设备离线异常，3:灯珠异常',
  `occurrence_time` varchar(30) DEFAULT NULL,
  `serial_number` varchar(100) DEFAULT NULL,
  `device_name` varchar(255) DEFAULT NULL,
  `device_ip` varchar(255) DEFAULT NULL,
  `device_attribute` varchar(255) DEFAULT NULL,
  `device_brightness` varchar(255) DEFAULT NULL,
  `device_temperature` varchar(255) DEFAULT NULL,
  `handle_status` varchar(255) DEFAULT NULL,
  `delete_status` tinyint(4) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for group_device_bind
-- ----------------------------
DROP TABLE IF EXISTS `group_device_bind`;
CREATE TABLE `group_device_bind` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `group_id` int(11) DEFAULT NULL,
  `device_id` int(11) DEFAULT NULL,
  `create_by` varchar(50) DEFAULT NULL,
  `create_date` datetime DEFAULT NULL,
  `update_by` varchar(50) DEFAULT NULL,
  `update_date` datetime DEFAULT NULL,
  `delete_status` tinyint(4) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_group_device_bind_device_group_info_1` (`group_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for plat_device_logger
-- ----------------------------
DROP TABLE IF EXISTS `plat_device_logger`;
CREATE TABLE `plat_device_logger` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `direction` varchar(255) DEFAULT NULL,
  `biz_type` varchar(255) DEFAULT NULL,
  `message` text,
  `ret_code` varchar(255) DEFAULT NULL,
  `ret_msg` varchar(255) DEFAULT NULL,
  `serial_number` varchar(200) DEFAULT NULL,
  `handle_time` varchar(20) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for plat_user_logger
-- ----------------------------
DROP TABLE IF EXISTS `plat_user_logger`;
CREATE TABLE `plat_user_logger` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` varchar(50) DEFAULT NULL,
  `handle_time` datetime DEFAULT NULL,
  `req_method` varchar(255) DEFAULT NULL,
  `req_url` varchar(255) DEFAULT NULL,
  `req_param` varchar(255) DEFAULT NULL,
  `ret_msg` varchar(255) DEFAULT NULL,
  `ret_code` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for sys_global_config
-- ----------------------------
DROP TABLE IF EXISTS `sys_global_config`;
CREATE TABLE `sys_global_config` (
  `item_id` int(11) NOT NULL,
  `item_name` varchar(255) DEFAULT NULL,
  `item_value` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`item_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of sys_global_config
-- ----------------------------
BEGIN;
INSERT INTO `sys_global_config` VALUES (1, '广播地址', '0.0.0.0');
INSERT INTO `sys_global_config` VALUES (2, '广播端口', '8900');
INSERT INTO `sys_global_config` VALUES (3, '温度阀值', '80');
COMMIT;

-- ----------------------------
-- Table structure for sys_user_info
-- ----------------------------
DROP TABLE IF EXISTS `sys_user_info`;
CREATE TABLE `sys_user_info` (
  `user_id` varchar(60) NOT NULL,
  `password` varchar(255) DEFAULT NULL,
  `nickname` varchar(255) DEFAULT NULL,
  `remark` varchar(255) DEFAULT NULL,
  `mobile_phone` varchar(20) DEFAULT NULL,
  `avatar` varchar(300) DEFAULT NULL,
  `weixin` varchar(255) DEFAULT NULL,
  `qq` varchar(255) DEFAULT NULL,
  `email` varchar(255) DEFAULT NULL,
  `delete_status` tinyint(4) DEFAULT NULL,
  PRIMARY KEY (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of sys_user_info
-- ----------------------------
BEGIN;
INSERT INTO `sys_user_info` VALUES ('admin', 'e10adc3949ba59abbe56e057f20f883e', '超级管理员', '最高权限', '10000000000', '/ui/avatar.png', 'ai', '260979971', 'hzwy23@163.com', 0);
INSERT INTO `sys_user_info` VALUES ('inspector', 'e10adc3949ba59abbe56e057f20f883e', '巡视人员', '巡视人员', '10000000000', '/ui/avatar.png', 'ai-lamp', '260979971', 'hzwy23@163.com', 0);
INSERT INTO `sys_user_info` VALUES ('operation', 'e10adc3949ba59abbe56e057f20f883e', '运维管理员', '日常运维', '10000000000', '/ui/avatar.png', 'ai-lamp', '260979971', 'hzwy23@163.com', 0);
COMMIT;

SET FOREIGN_KEY_CHECKS = 1;
