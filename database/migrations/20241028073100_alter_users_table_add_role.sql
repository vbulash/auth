-- +goose Up
-- +goose StatementBegin
ALTER TABLE public.users ADD "role" int DEFAULT 1 NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE public.chats DROP COLUMN role;
-- +goose StatementEnd
