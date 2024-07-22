package db

const insert string = `INSERT INTO users (id, name, last_name, email, active, created_at, deactivated_at) VALUES ($1, $2, $3, $4, $5, $6, $7);`
