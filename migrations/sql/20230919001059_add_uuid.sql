-- +goose Up
-- disable the enforcement of foreign-keys constraints
PRAGMA foreign_keys = off;
-- create "new_messages" table
CREATE TABLE `new_messages` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `uuid` text NOT NULL DEFAULT '', `from` text NOT NULL, `to` json NOT NULL, `subject` text NOT NULL, `text` text NOT NULL, `html` text NOT NULL, `date` datetime NOT NULL, `created_at` datetime NOT NULL);
-- copy rows from old table "messages" to new temporary table "new_messages"
INSERT INTO `new_messages` (`id`, `from`, `to`, `subject`, `text`, `html`, `date`, `created_at`) SELECT `id`, `from`, `to`, `subject`, `text`, `html`, `date`, `created_at` FROM `messages`;
-- drop "messages" table after copying rows
DROP TABLE `messages`;
-- rename temporary table "new_messages" to "messages"
ALTER TABLE `new_messages` RENAME TO `messages`;
-- create "new_attachments" table
CREATE TABLE `new_attachments` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `uuid` text NOT NULL DEFAULT '', `message_id` integer NULL, `name` text NOT NULL, `mime` text NOT NULL, `extension` text NOT NULL, CONSTRAINT `message_id` FOREIGN KEY (`message_id`) REFERENCES `messages` (`id`) ON UPDATE CASCADE ON DELETE SET NULL);
-- copy rows from old table "attachments" to new temporary table "new_attachments"
INSERT INTO `new_attachments` (`id`, `message_id`, `name`, `mime`, `extension`) SELECT `id`, `message_id`, `name`, `mime`, `extension` FROM `attachments`;
-- drop "attachments" table after copying rows
DROP TABLE `attachments`;
-- rename temporary table "new_attachments" to "attachments"
ALTER TABLE `new_attachments` RENAME TO `attachments`;
-- enable back the enforcement of foreign-keys constraints
PRAGMA foreign_keys = on;

-- +goose Down
-- reverse: create "new_attachments" table
DROP TABLE `new_attachments`;
-- reverse: create "new_messages" table
DROP TABLE `new_messages`;
