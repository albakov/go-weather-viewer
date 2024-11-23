-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS sessions
(
    id         VARCHAR(255) PRIMARY KEY,
    user_id    BIGINT UNSIGNED NOT NULL,
    expires_at DATETIME        NOT NULL,
    CONSTRAINT `sessions_user_id_fn`
        FOREIGN KEY (user_id) REFERENCES users (id)
            ON DELETE CASCADE
            ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS sessions;
-- +goose StatementEnd
