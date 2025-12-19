-- ----------------------------
-- Table structure for system_sensitive_word
-- ----------------------------
DROP TABLE IF EXISTS `system_sensitive_word`;
CREATE TABLE `system_sensitive_word` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '编号',
  `name` varchar(255) NOT NULL COMMENT '敏感词',
  `tags` json NOT NULL COMMENT '标签',
  `status` int NOT NULL DEFAULT '0' COMMENT '状态: 0-开启, 1-关闭',
  `description` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '描述',
  `creator` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '创建者',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updater` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '更新者',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted` bit(1) NOT NULL DEFAULT b'0' COMMENT '是否删除',
  PRIMARY KEY (`id`),
  KEY `idx_name` (`name`),
  KEY `idx_create_time` (`create_time`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='系统敏感词';

-- ----------------------------
-- Migration: Add business hours and employee binding fields to trade_delivery_pick_up_store
-- Purpose: Align with Java version for pickup store management
-- Date: 2025-12-19
-- ----------------------------

-- Step 1: Add business hours fields (opening_time and closing_time)
ALTER TABLE `trade_delivery_pick_up_store`
ADD COLUMN `opening_time` TIME NULL COMMENT '营业开始时间' AFTER `logo`,
ADD COLUMN `closing_time` TIME NULL COMMENT '营业结束时间' AFTER `opening_time`;

-- Step 2: Add employee binding field (verify_user_ids)
-- Note: This field stores employee IDs as JSON array (e.g., [10,11,12])
ALTER TABLE `trade_delivery_pick_up_store`
ADD COLUMN `verify_user_ids` VARCHAR(500) NULL COMMENT '核销员工用户编号数组' AFTER `closing_time`;

-- Step 3: Add index for status field (optional performance optimization)
ALTER TABLE `trade_delivery_pick_up_store`
ADD INDEX `idx_status` (`status`);

-- ----------------------------
-- End of Migration
-- ----------------------------
