CREATE TYPE "PermType" AS ENUM (
  'ManagePost',
  'CreatePost',
  'ManageComments',
  'CreateComment',
  'ManageUsers'
);

CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "permgroup_id" int NOT NULL DEFAULT 1,
  "name" text NOT NULL,
  "email" text UNIQUE NOT NULL,
  "password_hash" text NOT NULL,
  "bio" text,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "posts" (
  "id" bigserial PRIMARY KEY,
  "book_title" text NOT NULL,
  "vote" smallint NOT NULL,
  "user_id" bigint NOT NULL,
  "title" text,
  "content" text NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "bookshelf" (
  "book_id" text NOT NULL,
  "user_id" bigint NOT NULL,
  "in_bookshelf" bool NOT NULL DEFAULT true,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "comments" (
  "id" bigserial PRIMARY KEY,
  "post_id" bigint NOT NULL,
  "user_id" bigint NOT NULL,
  "content" text NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "post_likes" (
  "post_id" bigint NOT NULL,
  "user_id" bigint NOT NULL,
  "is_liked" bool NOT NULL DEFAULT true,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "permgroups" (
  "id" serial PRIMARY KEY,
  "name" varchar(128) UNIQUE NOT NULL,
  "permissions" "PermType"[] NOT NULL
);

CREATE INDEX ON "users" ("email");

CREATE UNIQUE INDEX ON "bookshelf" ("book_id", "user_id");

CREATE INDEX ON "comments" ("user_id");

CREATE INDEX ON "comments" ("post_id");

CREATE UNIQUE INDEX ON "post_likes" ("post_id", "user_id");

CREATE INDEX ON "permgroups" ("name");

ALTER TABLE "users" ADD FOREIGN KEY ("permgroup_id") REFERENCES "permgroups" ("id");

ALTER TABLE "posts" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;

ALTER TABLE "bookshelf" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;

ALTER TABLE "comments" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;

ALTER TABLE "comments" ADD FOREIGN KEY ("post_id") REFERENCES "posts" ("id") ON DELETE CASCADE;

ALTER TABLE "post_likes" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;

ALTER TABLE "post_likes" ADD FOREIGN KEY ("post_id") REFERENCES "posts" ("id") ON DELETE CASCADE;