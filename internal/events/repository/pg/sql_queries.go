package db

const insert string = `INSERT INTO events_%d (id, user_id, active, title, daily, isSent, subject, author, letter_id, "timestamp") VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);`
