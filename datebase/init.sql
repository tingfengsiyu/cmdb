--
-- Host: 172.22.0.20    Database: cmdb
-- ------------------------------------------------------
-- Server version	8.0.23

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `user`
--
create
    database cmdb ;
use
    cmdb;
DROP TABLE IF EXISTS `user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `user`
(
    `id`         bigint unsigned NOT NULL AUTO_INCREMENT,
    `created_at` datetime(3) DEFAULT NULL,
    `updated_at` datetime(3) DEFAULT NULL,
    `deleted_at` datetime(3) DEFAULT NULL,
    `username`   varchar(20) NOT NULL,
    `password`   varchar(20) NOT NULL,
    `role`       bigint DEFAULT '2',
    PRIMARY KEY (`id`),
    KEY          `idx_user_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user`
--

LOCK
    TABLES `user` WRITE;
/*!40000 ALTER TABLE `user` DISABLE KEYS */;
INSERT INTO `user`
VALUES (4, '2021-01-26 15:22:12.588', '2021-01-26 15:22:12.588', NULL, 'luofeng', 'rytBNnf2J2vciy+8', 1);
/*!40000 ALTER TABLE `user` ENABLE KEYS */;
UNLOCK
    TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2021-03-25 13:13:10

CREATE TABLE `ops_records` (
                               `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
                               `created_at` datetime(3) DEFAULT NULL,
                               `updated_at` datetime(3) DEFAULT NULL,
                               `deleted_at` datetime(3) DEFAULT NULL,
                               `user` varchar(30) NOT NULL,
                               `object` varchar(1000) DEFAULT NULL,
                               `action` varchar(500) NOT NULL,
                               `state` bigint(20) NOT NULL DEFAULT '2',
                               `success` longtext NOT NULL DEFAULT 'success' ,
                               `error` longtext NOT NULL DEFAULT 'error' ,
                               PRIMARY KEY (`id`),
                               KEY `idx_ops_records_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;

CREATE TABLE `user_permissions` (   `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,   `created_at` datetime(3) DEFAULT NULL,
                                    `updated_at` datetime(3) DEFAULT NULL,   `deleted_at` datetime(3) DEFAULT NULL,
                                    `server_id` bigint(20) NOT NULL,
                                    `user_id` bigint(20) NOT NULL,
                                    `term_user_id` bigint(20) NOT NULL,
                                    `group` varchar(50) NOT NULL,
                                    PRIMARY KEY (`id`)) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `term_user` (
                             `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
                             `created_at` datetime(3) DEFAULT NULL,
                             `updated_at` datetime(3) DEFAULT NULL,
                             `deleted_at` datetime(3) DEFAULT NULL,
                             `port` bigint(20) NOT NULL,
                             `username` varchar(50) NOT NULL,
                             `password` varchar(50) NOT NULL,
                             `identity_file` varchar(50) NOT NULL,
                             `protocol` varchar(50) NOT NULL,
                             PRIMARY KEY (`id`),
                             UNIQUE KEY `username` (`username`),
                             UNIQUE KEY `password` (`password`),
                             KEY `idx_term_user_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8

