create database cs_s_update charset utf8mb4;

use cs_s_update;

CREATE TABLE `t_sqllog` (
  `id` int NOT NULL AUTO_INCREMENT,
  `sqlid` int DEFAULT NULL,
  `dbname` varchar(50) NOT NULL,
  `isdone` tinyint DEFAULT NULL,
  `create_date` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `t_sql` (
  `id` int NOT NULL AUTO_INCREMENT,
  `esql` varchar(8000) DEFAULT NULL,
  `create_date` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
