-- Disable the enforcement of foreign-keys constraints
PRAGMA foreign_keys = off;
-- Create "categories" table
CREATE TABLE `categories` (
  `id` integer NULL PRIMARY KEY AUTOINCREMENT,
  `name` text NULL
);
-- Create "new_todos" table
CREATE TABLE `new_todos` (
  `id` integer NULL PRIMARY KEY AUTOINCREMENT,
  `title` text NULL,
  `completed` numeric NULL,
  `category_id` integer NULL,
  CONSTRAINT `fk_todos_category` FOREIGN KEY (`category_id`) REFERENCES `categories` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Copy rows from old table "todos" to new temporary table "new_todos"
INSERT INTO `new_todos` (`id`, `title`, `completed`) SELECT `id`, `title`, `completed` FROM `todos`;
-- Drop "todos" table after copying rows
DROP TABLE `todos`;
-- Rename temporary table "new_todos" to "todos"
ALTER TABLE `new_todos` RENAME TO `todos`;
-- Enable back the enforcement of foreign-keys constraints
PRAGMA foreign_keys = on;
