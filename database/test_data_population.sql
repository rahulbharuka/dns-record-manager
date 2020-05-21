INSERT IGNORE INTO `cluster`(`id`, `name`,`subdomain`)
VALUES
(1, 'Los Angeles', 'la'),
(2, 'New York', 'nyc'),
(3, 'Frankfurt', 'fra'),
(4, 'Hong Kong', 'hongkong');

INSERT IGNORE INTO `server`(`id`, `name`, `cluster_id`, `ip`, `added_to_rotation`)
VALUES
(1, 'ubiq-1', 3, '123.123.123.123', 0),
(2, 'ubiq-2', 3, '234.234.234.234', 0),
(3, 'leaseweb-de-1', 4, '34.34.34.34', 0),
(4, 'rackspace-1', 1, '45.45.45.45', 0),
(5, 'rackspace-2', 1, '45.45.45.46', 0);