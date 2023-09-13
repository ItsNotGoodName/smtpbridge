-- +goose Up
-- create "messages" table
CREATE TABLE `messages` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `from` text NOT NULL, `to` json NOT NULL, `subject` text NOT NULL, `text` text NOT NULL, `html` text NOT NULL, `date` datetime NOT NULL, `created_at` datetime NOT NULL);
-- create "attachments" table
CREATE TABLE `attachments` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `message_id` integer NULL, `name` text NOT NULL, `mime` text NOT NULL, `extension` text NOT NULL, CONSTRAINT `message_id` FOREIGN KEY (`message_id`) REFERENCES `messages` (`id`) ON UPDATE NO ACTION ON DELETE SET NULL);
-- create "endpoints" table
CREATE TABLE `endpoints` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `internal` boolean NOT NULL, `internal_id` text NULL, `name` text NOT NULL, `attachment_disable` boolean NOT NULL, `text_disable` boolean NOT NULL, `title_template` text NOT NULL, `body_template` text NOT NULL, `kind` text NOT NULL, `config` json NOT NULL, `updated_at` datetime NOT NULL, `created_at` datetime NOT NULL);
-- create index "endpoints_internal_id_idx" to table: "endpoints"
CREATE UNIQUE INDEX `endpoints_internal_id_idx` ON `endpoints` (`internal_id`);
-- create "rules" table
CREATE TABLE `rules` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `internal` boolean NOT NULL, `internal_id` text NULL, `name` text NOT NULL, `expression` text NOT NULL, `enable` boolean NOT NULL, `updated_at` datetime NOT NULL, `created_at` datetime NOT NULL);
-- create index "rules_internal_id_idx" to table: "rules"
CREATE UNIQUE INDEX `rules_internal_id_idx` ON `rules` (`internal_id`);
-- create "rules_to_endpoints" table
CREATE TABLE `rules_to_endpoints` (`internal` boolean NOT NULL, `rule_id` integer NOT NULL, `endpoint_id` integer NOT NULL, CONSTRAINT `endpoint_id` FOREIGN KEY (`endpoint_id`) REFERENCES `endpoints` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE, CONSTRAINT `rule_id` FOREIGN KEY (`rule_id`) REFERENCES `rules` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE);
-- create index "rules_to_endpoints_rule_id_endpoint_id_idx" to table: "rules_to_endpoints"
CREATE UNIQUE INDEX `rules_to_endpoints_rule_id_endpoint_id_idx` ON `rules_to_endpoints` (`rule_id`, `endpoint_id`);
-- create "traces" table
CREATE TABLE `traces` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `request_id` text NOT NULL, `source` text NOT NULL, `seq` integer NOT NULL, `action` text NOT NULL, `level` text NOT NULL, `data` json NOT NULL, `created_at` datetime NOT NULL);
-- create "mailman_queue" table
CREATE TABLE `mailman_queue` (`message_id` integer NOT NULL, `created_at` datetime NOT NULL, CONSTRAINT `message_id` FOREIGN KEY (`message_id`) REFERENCES `messages` (`id`) ON UPDATE CASCADE ON DELETE CASCADE);
-- create index "mailman_queue_message_id_idx" to table: "mailman_queue"
CREATE UNIQUE INDEX `mailman_queue_message_id_idx` ON `mailman_queue` (`message_id`);

-- +goose Down
-- reverse: create index "mailman_queue_message_id_idx" to table: "mailman_queue"
DROP INDEX `mailman_queue_message_id_idx`;
-- reverse: create "mailman_queue" table
DROP TABLE `mailman_queue`;
-- reverse: create "traces" table
DROP TABLE `traces`;
-- reverse: create index "rules_to_endpoints_rule_id_endpoint_id_idx" to table: "rules_to_endpoints"
DROP INDEX `rules_to_endpoints_rule_id_endpoint_id_idx`;
-- reverse: create "rules_to_endpoints" table
DROP TABLE `rules_to_endpoints`;
-- reverse: create index "rules_internal_id_idx" to table: "rules"
DROP INDEX `rules_internal_id_idx`;
-- reverse: create "rules" table
DROP TABLE `rules`;
-- reverse: create index "endpoints_internal_id_idx" to table: "endpoints"
DROP INDEX `endpoints_internal_id_idx`;
-- reverse: create "endpoints" table
DROP TABLE `endpoints`;
-- reverse: create "attachments" table
DROP TABLE `attachments`;
-- reverse: create "messages" table
DROP TABLE `messages`;
