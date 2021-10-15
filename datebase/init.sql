/*
 Navicat Premium Data Transfer

 Source Server         : 郫县
 Source Server Type    : MySQL
 Source Server Version : 50733
 Source Host           : 30.10.0.18:3306
 Source Schema         : cmdb

 Target Server Type    : MySQL
 Target Server Version : 50733
 File Encoding         : 65001

 Date: 15/10/2021 13:47:12
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for cabinet
-- ----------------------------
DROP TABLE IF EXISTS `cabinet`;
CREATE TABLE `cabinet`  (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `idc_id` bigint(20) NOT NULL,
  `cabinet_number_id` bigint(20) NOT NULL,
  `cabinet_number` varchar(30) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_cabinet_deleted_at`(`deleted_at`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 85 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for cloud_instance
-- ----------------------------
DROP TABLE IF EXISTS `cloud_instance`;
CREATE TABLE `cloud_instance`  (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `instance_id` varchar(30) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `host_name` varchar(30) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `status` varchar(30) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `cpu` bigint(20) NOT NULL,
  `memory` bigint(20) NOT NULL,
  `os_name` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `region_id` varchar(30) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `instance_type` varchar(30) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `os_type` varchar(30) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `internet_max_bandwidth_in` varchar(30) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `start_time` datetime(3) DEFAULT NULL,
  `expired_time` datetime(3) DEFAULT NULL,
  `instance_creation_time` datetime(3) DEFAULT NULL,
  `local_storage_capacity` varchar(30) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `private_ip_address` varchar(30) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `public_ip_address` varchar(30) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `cloud` varchar(10) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_cloud_instance_deleted_at`(`deleted_at`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for idc
-- ----------------------------
DROP TABLE IF EXISTS `idc`;
CREATE TABLE `idc`  (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `city` varchar(30) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `idc_name` varchar(30) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `idc_id` bigint(20) NOT NULL,
  `cabinet_number_id` bigint(20) NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_idc_deleted_at`(`deleted_at`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 79 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for monitor_prometheus
-- ----------------------------
DROP TABLE IF EXISTS `monitor_prometheus`;
CREATE TABLE `monitor_prometheus`  (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `server_id` bigint(20) NOT NULL,
  `node_export_port` bigint(20) DEFAULT 9100,
  `process_export_port` bigint(20) DEFAULT 9256,
  `script_export_port` bigint(20) DEFAULT 9172,
  `node_export_status` bigint(20) DEFAULT 0,
  `process_export_status` bigint(20) DEFAULT 0,
  `script_export_status` bigint(20) DEFAULT 0,
  `disable_node_export` bigint(20) DEFAULT 1,
  `disable_process_export` bigint(20) DEFAULT 1,
  `disable_script_export` bigint(20) DEFAULT 1,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `server_id`(`server_id`) USING BTREE,
  INDEX `idx_monitor_prometheus_deleted_at`(`deleted_at`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 951 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for ops_records
-- ----------------------------
DROP TABLE IF EXISTS `ops_records`;
CREATE TABLE `ops_records`  (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `user` varchar(30) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `object` varchar(500) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `action` varchar(500) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `state` bigint(20) NOT NULL DEFAULT 2,
  `success` longtext CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `error` longtext CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_ops_records_deleted_at`(`deleted_at`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 29 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for server
-- ----------------------------
DROP TABLE IF EXISTS `server`;
CREATE TABLE `server`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(42) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `models` varchar(30) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `location` varchar(30) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `private_ip_address` varchar(30) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `public_ip_address` varchar(30) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `label` varchar(30) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `cluster` varchar(30) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `label_ip_address` varchar(30) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `cpu` varchar(500) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `memory` varchar(200) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `disk` varchar(500) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `user` varchar(30) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `state` varchar(10) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `server_id` bigint(20) NOT NULL,
  `idc_id` bigint(20) NOT NULL,
  `cabinet_number_id` bigint(20) NOT NULL,
  `gpu` varchar(500) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `mac` varchar(500) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `private_ip_address`(`private_ip_address`) USING BTREE,
  UNIQUE INDEX `private_ip_address_2`(`private_ip_address`) USING BTREE,
  UNIQUE INDEX `label_ip_address`(`label_ip_address`) USING BTREE,
  INDEX `idx_server_deleted_at`(`deleted_at`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 496 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for term_log
-- ----------------------------
DROP TABLE IF EXISTS `term_log`;
CREATE TABLE `term_log`  (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `protocol` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `term_user` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `client_ip` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `log` text CHARACTER SET utf8 COLLATE utf8_general_ci,
  `user` varchar(30) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `private_ip_address` varchar(30) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_term_log_deleted_at`(`deleted_at`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 22 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for term_user
-- ----------------------------
DROP TABLE IF EXISTS `term_user`;
CREATE TABLE `term_user`  (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `port` bigint(20) NOT NULL,
  `username` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `password` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `identity_file` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `protocol` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `username`(`username`) USING BTREE,
  UNIQUE INDEX `password`(`password`) USING BTREE,
  INDEX `idx_term_user_deleted_at`(`deleted_at`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 17 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user`  (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `username` varchar(20) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `password` varchar(20) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `role` bigint(20) DEFAULT 2,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_user_deleted_at`(`deleted_at`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 13 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for user_permissions
-- ----------------------------
DROP TABLE IF EXISTS `user_permissions`;
CREATE TABLE `user_permissions`  (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `server_id` bigint(20) NOT NULL,
  `user_id` bigint(20) NOT NULL,
  `term_user_id` bigint(20) NOT NULL,
  `group` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_user_permissions_deleted_at`(`deleted_at`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1931 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

SET FOREIGN_KEY_CHECKS = 1;
INSERT INTO `user`
VALUES (4, '2021-01-26 15:22:12.588', '2021-01-26 15:22:12.588', NULL, 'luofeng', 'rytBNnf2J2vciy+8', 1);
