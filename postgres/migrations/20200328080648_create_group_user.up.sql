CREATE TABLE IF NOT EXISTS  group_user(
    group_id BIGSERIAL REFERENCES groups (id) ON DELETE CASCADE NOT NULL,
    user_id BIGSERIAL REFERENCES users (id) ON DELETE CASCADE NOT NULL
)