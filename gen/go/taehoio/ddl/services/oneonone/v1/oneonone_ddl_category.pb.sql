-- Code generated by protoc-gen-go-ddl. DO NOT EDIT.
-- versions:
--  protoc-gen-go-ddl v0.0.1-alpha
--  protoc            (unknown)
-- source: taehoio/ddl/services/oneonone/v1/oneonone.proto

CREATE TABLE `category` (
	`id` BIGINT UNSIGNED,
	`created_at` TIMESTAMP(6) NULL DEFAULT NULL,
	`updated_at` TIMESTAMP(6) NULL DEFAULT NULL,
	`deleted_at` TIMESTAMP(6) NULL DEFAULT NULL,
	`name` VARCHAR(255),
	PRIMARY KEY (`id`)
);
