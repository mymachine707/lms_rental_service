BEGIN;

CREATE TYPE "my_type" AS ENUM ('RENTAL', 'RETURNED', 'NOT RETURNED');

CREATE TABLE IF NOT EXISTS "rental_book" (
	"id" CHAR(36) PRIMARY KEY,
	"book_id" CHAR(36)  NOT NULL,
	"book_name" VARCHAR(255)  NOT NULL,
	"user_id" CHAR(36) NOT NULL,
	"rental_date_time" TIMESTAMP DEFAULT now(),
    "expected_return_date" TIMESTAMP,
	"return_date" TIMESTAMP,
	"book_fines" VARCHAR(255),
	"rental_status" my_type  DEFAULT 'RENTAL',
	"deleted_at" TIMESTAMP
);

COMMIT;