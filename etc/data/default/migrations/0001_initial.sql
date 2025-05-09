-- +goose Up
CREATE TABLE IF NOT EXISTS "profile" (
  "id" CHAR(26) NOT NULL PRIMARY KEY,
  "slug" TEXT NOT NULL CONSTRAINT "profile_slug_unique" UNIQUE,
  "kind" TEXT NOT NULL,
  "custom_domain" TEXT,
  "profile_picture_uri" TEXT,
  "pronouns" TEXT,
  "title" TEXT NOT NULL,
  "description" TEXT NOT NULL,
  "properties" JSONB,
  "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
  "updated_at" TIMESTAMP WITH TIME ZONE,
  "deleted_at" TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS "user" (
  "id" CHAR(26) NOT NULL PRIMARY KEY,
  "kind" TEXT NOT NULL,
  "name" TEXT NOT NULL,
  "email" TEXT CONSTRAINT "user_email_unique" UNIQUE,
  "phone" TEXT,
  "github_handle" TEXT,
  "github_remote_id" TEXT CONSTRAINT "user_github_remote_id_unique" UNIQUE,
  "bsky_handle" TEXT,
  "bsky_remote_id" TEXT,
  "x_handle" TEXT,
  "x_remote_id" TEXT,
  "individual_profile_id" CHAR(26) CONSTRAINT "user_individual_profile_id_fk" REFERENCES "profile",
  "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
  "updated_at" TIMESTAMP WITH TIME ZONE,
  "deleted_at" TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS "profile_membership" (
  "id" CHAR(26) NOT NULL PRIMARY KEY,
  "profile_id" CHAR(26) NOT NULL CONSTRAINT "profile_membership_profile_id_fk" REFERENCES "profile",
  "user_id" CHAR(26) NOT NULL CONSTRAINT "profile_membership_user_id_fk" REFERENCES "user",
  "kind" TEXT NOT NULL,
  "properties" JSONB,
  "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
  "updated_at" TIMESTAMP WITH TIME ZONE,
  "deleted_at" TIMESTAMP WITH TIME ZONE,
  CONSTRAINT "profile_membership_profile_id_user_id_unique" UNIQUE ("profile_id", "user_id")
);

CREATE TABLE IF NOT EXISTS "profile_link" (
  "id" CHAR(26) NOT NULL PRIMARY KEY,
  "profile_id" CHAR(26) NOT NULL CONSTRAINT "profile_link_profile_id_fk" REFERENCES "profile",
  "kind" TEXT NOT NULL,
  "order" INTEGER NOT NULL,
  "is_managed" BOOLEAN DEFAULT FALSE NOT NULL,
  "is_verified" BOOLEAN DEFAULT FALSE NOT NULL,
  "remote_id" TEXT,
  "public_id" TEXT,
  "uri" TEXT,
  "title" TEXT NOT NULL,
  "auth_provider" TEXT NOT NULL,
  "auth_access_token_scope" TEXT NOT NULL,
  "auth_access_token" TEXT NOT NULL,
  "auth_access_token_expires_at" TIMESTAMP WITH TIME ZONE,
  "auth_refresh_token" TEXT,
  "auth_refresh_token_expires_at" TIMESTAMP WITH TIME ZONE,
  "properties" JSONB,
  "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
  "updated_at" TIMESTAMP WITH TIME ZONE,
  "deleted_at" TIMESTAMP WITH TIME ZONE,
  CONSTRAINT "profile_link_profile_id_kind_remote_id_unique" UNIQUE ("profile_id", "kind", "remote_id")
);

CREATE TABLE IF NOT EXISTS "profile_link_import" (
  "id" CHAR(26) NOT NULL PRIMARY KEY,
  "profile_link_id" CHAR(26) NOT NULL CONSTRAINT "profile_link_import_profile_link_id_fk" REFERENCES "profile_link",
  "remote_id" TEXT,
  "properties" JSONB,
  "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
  "updated_at" TIMESTAMP WITH TIME ZONE,
  "deleted_at" TIMESTAMP WITH TIME ZONE,
  CONSTRAINT "profile_link_import_profile_link_id_remote_id_unique" UNIQUE ("profile_link_id", "remote_id")
);

CREATE TABLE IF NOT EXISTS "profile_page" (
  "id" CHAR(26) NOT NULL PRIMARY KEY,
  "profile_id" CHAR(26) NOT NULL CONSTRAINT "profile_page_profile_id_fk" REFERENCES "profile",
  "slug" TEXT NOT NULL,
  "order" INTEGER NOT NULL,
  "cover_picture_uri" TEXT,
  "title" TEXT NOT NULL,
  "summary" TEXT NOT NULL,
  "content" TEXT NOT NULL,
  "published_at" TIMESTAMP WITH TIME ZONE,
  "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
  "updated_at" TIMESTAMP WITH TIME ZONE,
  "deleted_at" TIMESTAMP WITH TIME ZONE,
  CONSTRAINT "profile_page_profile_id_slug_unique" UNIQUE ("profile_id", "slug")
);

CREATE TABLE IF NOT EXISTS "session" (
  "id" CHAR(26) NOT NULL PRIMARY KEY,
  "status" TEXT NOT NULL,
  "oauth_request_state" TEXT NOT NULL,
  "oauth_request_code_verifier" TEXT NOT NULL,
  "oauth_redirect_uri" TEXT,
  "logged_in_user_id" CHAR(26) CONSTRAINT "session_logged_in_user_id_fk" REFERENCES "user",
  "logged_in_at" TIMESTAMP WITH TIME ZONE,
  "expires_at" TIMESTAMP WITH TIME ZONE,
  "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
  "updated_at" TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS "session_logged_in_user_id_index" ON "session" ("logged_in_user_id");

CREATE TABLE IF NOT EXISTS "question" (
  "id" CHAR(26) NOT NULL PRIMARY KEY,
  "user_id" CHAR(26) NOT NULL CONSTRAINT "question_user_id_fk" REFERENCES "user",
  "profile_id" CHAR(26) CONSTRAINT "question_profile_id_fk" REFERENCES "profile",
  "content" TEXT NOT NULL,
  "is_hidden" BOOLEAN DEFAULT FALSE NOT NULL,
  "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
  "updated_at" TIMESTAMP WITH TIME ZONE,
  "deleted_at" TIMESTAMP WITH TIME ZONE,
  "answered_at" TIMESTAMP WITH TIME ZONE,
  "answer_uri" TEXT,
  "is_anonymous" BOOLEAN DEFAULT FALSE NOT NULL,
  "answer_kind" TEXT,
  "answer_content" TEXT
);

CREATE TABLE IF NOT EXISTS "question_vote" (
  "id" CHAR(26) NOT NULL PRIMARY KEY,
  "question_id" CHAR(26) NOT NULL CONSTRAINT "question_vote_question_id_fk" REFERENCES "question",
  "user_id" CHAR(26) NOT NULL CONSTRAINT "question_vote_user_id_fk" REFERENCES "user",
  "score" INTEGER NOT NULL,
  "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
  CONSTRAINT "question_vote_question_id_user_id_unique" UNIQUE ("question_id", "user_id")
);

CREATE TABLE IF NOT EXISTS "event_series" (
  "id" CHAR(26) NOT NULL PRIMARY KEY,
  "slug" TEXT NOT NULL CONSTRAINT "event_series_slug_unique" UNIQUE,
  "event_picture_uri" TEXT,
  "title" TEXT NOT NULL,
  "description" TEXT NOT NULL,
  "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
  "updated_at" TIMESTAMP WITH TIME ZONE,
  "deleted_at" TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS "event" (
  "id" CHAR(26) NOT NULL PRIMARY KEY,
  "slug" TEXT NOT NULL CONSTRAINT "event_slug_unique" UNIQUE,
  "kind" TEXT NOT NULL,
  "event_picture_uri" TEXT,
  "title" TEXT NOT NULL,
  "description" TEXT NOT NULL,
  "time_start" TIMESTAMP WITH TIME ZONE NOT NULL,
  "time_end" TIMESTAMP WITH TIME ZONE NOT NULL,
  "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
  "updated_at" TIMESTAMP WITH TIME ZONE,
  "deleted_at" TIMESTAMP WITH TIME ZONE,
  "series_id" CHAR(26) CONSTRAINT "event_series_id_fk" REFERENCES "event_series",
  "status" TEXT DEFAULT 'draft'::TEXT NOT NULL,
  "attendance_uri" TEXT,
  "published_at" TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS "event_attendance" (
  "id" CHAR(26) NOT NULL PRIMARY KEY,
  "event_id" CHAR(26) NOT NULL CONSTRAINT "event_attendance_event_id_fk" REFERENCES "event",
  "profile_id" CHAR(26) NOT NULL CONSTRAINT "event_attendance_profile_id_fk" REFERENCES "profile",
  "kind" TEXT NOT NULL,
  "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
  "updated_at" TIMESTAMP WITH TIME ZONE,
  "deleted_at" TIMESTAMP WITH TIME ZONE,
  CONSTRAINT "event_attendance_event_id_profile_id_unique" UNIQUE ("event_id", "profile_id")
);

CREATE TABLE IF NOT EXISTS "story" (
  "id" CHAR(26) NOT NULL PRIMARY KEY,
  "author_profile_id" CHAR(26) CONSTRAINT "story_author_profile_id_fk" REFERENCES "profile",
  "slug" TEXT NOT NULL,
  "kind" TEXT NOT NULL,
  "status" TEXT NOT NULL,
  "is_featured" BOOLEAN NOT NULL DEFAULT FALSE,
  "story_picture_uri" TEXT,
  "title" TEXT NOT NULL,
  "summary" TEXT NOT NULL,
  "content" TEXT NOT NULL,
  "properties" JSONB,
  "published_at" TIMESTAMP WITH TIME ZONE,
  "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
  "updated_at" TIMESTAMP WITH TIME ZONE,
  "deleted_at" TIMESTAMP WITH TIME ZONE,
  CONSTRAINT "story_author_profile_id_slug_unique" UNIQUE ("author_profile_id", "slug")
);

CREATE TABLE IF NOT EXISTS "custom_domain" (
  "id" CHAR(26) NOT NULL PRIMARY KEY,
  "domain" TEXT NOT NULL CONSTRAINT "custom_domain_domain_unique" UNIQUE,
  "profile_id" CHAR(26) NOT NULL CONSTRAINT "custom_domain_profile_id_fk" REFERENCES "profile",
  "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
  "updated_at" TIMESTAMP WITH TIME ZONE,
  "deleted_at" TIMESTAMP WITH TIME ZONE
);

-- +goose Down
DROP TABLE IF EXISTS "custom_domain";

DROP TABLE IF EXISTS "story";

DROP TABLE IF EXISTS "event_attendance";

DROP TABLE IF EXISTS "event";

DROP TABLE IF EXISTS "event_series";

DROP TABLE IF EXISTS "question_vote";

DROP TABLE IF EXISTS "question";

DROP INDEX IF EXISTS "session_logged_in_user_id_index";

DROP TABLE IF EXISTS "session";

DROP TABLE IF EXISTS "profile_page";

DROP TABLE IF EXISTS "profile_link_import";

DROP TABLE IF EXISTS "profile_link";

DROP TABLE IF EXISTS "profile_membership";

DROP TABLE IF EXISTS "user";

DROP TABLE IF EXISTS "profile";
