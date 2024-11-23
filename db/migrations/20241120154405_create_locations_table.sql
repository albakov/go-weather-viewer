-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS locations
(
    id        BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    name      VARCHAR(255)    NOT NULL,
    user_id   BIGINT UNSIGNED NOT NULL,
    latitude  DECIMAL(9, 6)   NOT NULL,
    longitude DECIMAL(10, 6)  NOT NULL,
    CONSTRAINT `locations_user_id_fn`
        FOREIGN KEY (user_id) REFERENCES users (id)
            ON DELETE CASCADE
            ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS locations;
-- +goose StatementEnd
