package sqlite

const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	date VARCHAR(8) NOT NULL DEFAULT "",
	title VARCHAR(256) NOT NULL DEFAULT "",
	comment TEXT NOT NULL DEFAULT "",
	repeat VARCHAR(128) NOT NULL DEFAULT ""
);
CREATE INDEX IF NOT EXISTS scheduler_date_idx ON scheduler (date);
`
const insertTaskQuery = `
INSERT INTO scheduler (date, title, comment, repeat) 
VALUES (:date, :title, :comment, :repeat)
`
const selectTasksQuery = `
SELECT id, date, title, comment, repeat 
FROM scheduler 
ORDER BY date ASC 
LIMIT :limit
`
const selectTasksByDateQuery = `
SELECT id, date, title, comment, repeat 
FROM scheduler 
WHERE date = :date 
ORDER BY date ASC 
LIMIT :limit
`
const selectTasksBySearchQuery = `
SELECT id, date, title, comment, repeat 
FROM scheduler 
WHERE title LIKE :title_like OR comment LIKE :comment_like
ORDER BY date ASC 
LIMIT :limit
`
const selectTaskByIDQuery = `
SELECT id, date, title, comment, repeat 
FROM scheduler 
WHERE id = :id
`
const updateTaskQuery = `
UPDATE scheduler 
SET date = :date, title = :title, comment = :comment, repeat = :repeat  
WHERE id = :id
`
const deleteTaskQuery = `
DELETE FROM scheduler 
WHERE id = :id
`
