-- +goose Up
CREATE TABLE task_completions (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    task_id UUID NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    
    UNIQUE(user_id, task_id) 
);

-- +goose Down
DROP TABLE task_completions;
