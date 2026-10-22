CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	username TEXT UNIQUE NOT NULL,
	password_hash TEXT NOT NULL,
	permissions INTEGER NOT NULL,
);

CREATE TABLE IF NOT EXISTS servers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    owner_id INTEGER NOT NULL,
    server_name TEXT NOT NULL,
    container_uuid TEXT NOT NULL,
    disk_limit_mb INTEGER NOT NULL,
    cpu_limit_percentage INTEGER NOT NULL CHECK (cpu_limit_percentage >= 0),
    ram_limit_mb INTEGER NOT NULL CHECK (ram_limit_mb >= 0)
)
