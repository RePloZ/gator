-- +goose Up
CREATE TABLE users
(
    "id" UUID PRIMARY KEY,
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL,
    "updated_at" TIMESTAMP WITH TIME ZONE NOT NULL,
    "name" TEXT UNIQUE NOT NULL
);

-- +goose Down
DROP TABLE users;