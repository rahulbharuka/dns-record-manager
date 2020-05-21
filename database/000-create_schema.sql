-- +migrate Up
CREATE TABLE `cluster` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(256) DEFAULT NULL,
  `subdomain` varchar(45) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `index_subdomain` (`subdomain`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


CREATE TABLE `server` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(256) NOT NULL,
  `cluster_id` int(11) NOT NULL,
  `ip` varchar(45) NOT NULL,
  `added_to_rotation` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `ip_UNIQUE` (`ip`),
  KEY `clutser_id_index` (`cluster_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- +migrate Down
DROP TABLE IF EXISTS `cluster`;
DROP TABLE IF EXISTS `server`;