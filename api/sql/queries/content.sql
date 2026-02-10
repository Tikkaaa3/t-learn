-- name: CreateCourse :one
INSERT INTO courses (id, created_at, updated_at, title, description)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: GetCourses :many
SELECT * FROM courses
ORDER BY created_at DESC;

-- name: GetCourse :one
SELECT * FROM courses WHERE id = $1;

-- name: CreateLesson :one
INSERT INTO lessons (id, created_at, updated_at, course_id, title, content, "position")
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2,
    $3,
    $4
)
RETURNING *;

-- name: GetLessonsByCourseID :many
SELECT * FROM lessons 
WHERE course_id = $1 
ORDER BY "position" ASC;

-- name: GetLessonsWithStatus :many
SELECT 
    l.id, 
    l.title, 
    l.position,
    l.course_id,
    CASE 
        WHEN tc.task_id IS NOT NULL THEN true 
        ELSE false 
    END as is_completed
FROM lessons l
LEFT JOIN tasks t ON t.lesson_id = l.id  -- <--- CHANGE TO 'LEFT JOIN'
LEFT JOIN task_completions tc 
    ON tc.task_id = t.id 
    AND tc.user_id = $2
WHERE l.course_id = $1
ORDER BY l.position;

-- name: GetLesson :one
SELECT * FROM lessons WHERE id = $1;

-- name: CreateTask :one
INSERT INTO tasks (id, created_at, updated_at, lesson_id, description)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: CreateTaskStep :one
INSERT INTO task_steps (id, created_at, updated_at, task_id, position, command, expected_output)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2,
    $3,
    $4
)
RETURNING *;

-- name: GetTaskByLessonID :one
SELECT * FROM tasks WHERE lesson_id = $1;

-- name: GetStepsByTaskID :many
SELECT * FROM task_steps 
WHERE task_id = $1 
ORDER BY position ASC;

-- name: CompleteTask :exec
INSERT INTO task_completions (id, created_at, updated_at, user_id, task_id)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1, -- User ID
    $2  -- Task ID
)
ON CONFLICT (user_id, task_id) DO NOTHING;

-- name: DeleteCourse :exec
DELETE FROM courses WHERE id = $1;

-- name: DeleteLesson :exec
DELETE FROM lessons WHERE id = $1;

-- name: DeleteTask :exec
DELETE FROM tasks WHERE id = $1;
