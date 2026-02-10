-- +goose Up
ALTER TABLE task_completions ADD CONSTRAINT unique_user_task UNIQUE (user_id, task_id);

-- +goose Down
ALTER TABLE task_completions DROP CONSTRAINT IF EXISTS unique_user_task;
