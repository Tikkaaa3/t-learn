-- +goose Up
-- Create the new table for multiple steps
CREATE TABLE task_steps (
    id UUID PRIMARY KEY,
    task_id UUID NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    position INT NOT NULL,
    command TEXT NOT NULL,
    expected_output TEXT NOT NULL,
    
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- Remove the old single-step columns from tasks
ALTER TABLE tasks DROP COLUMN command;
ALTER TABLE tasks DROP COLUMN expected_output;

-- +goose Down
ALTER TABLE tasks ADD COLUMN command TEXT NOT NULL DEFAULT '';
ALTER TABLE tasks ADD COLUMN expected_output TEXT NOT NULL DEFAULT '';
DROP TABLE task_steps;
