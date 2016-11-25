SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for demo1
-- ----------------------------
DROP TABLE IF EXISTS `demo1`;
CREATE TABLE `demo1` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `int_data` int(11) DEFAULT NULL,
  `str_data` varchar(255) DEFAULT NULL,
  `float_data` float DEFAULT NULL,
  `uint_data` int(10) unsigned zerofill DEFAULT NULL,
  `bool_data` tinyint(1) DEFAULT NULL,
  `time_date` time DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for demo_item
-- ----------------------------
DROP TABLE IF EXISTS `demo_item`;
CREATE TABLE `demo_item` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `item_id` int(11) NOT NULL,
  `num` int(10) unsigned zerofill NOT NULL,
  PRIMARY KEY (`id`,`item_id`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for tt
-- ----------------------------
DROP TABLE IF EXISTS `tt`;
CREATE TABLE `tt` (
  `ii` int(11) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
