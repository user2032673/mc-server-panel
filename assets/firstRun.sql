CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY AUTOINCREMENT CHECK (id >= 0), -- So as not to accidentally have a UID of -1, triggering sentinel values.
	username TEXT UNIQUE NOT NULL,
	password_hash TEXT NOT NULL,
);

CREATE TABLE IF NOT EXISTS servers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    owner_id INTEGER NOT NULL,
    server_name TEXT NOT NULL,
    container_uuid TEXT NOT NULL,
    disk_limit_mb INTEGER NOT NULL,
    cpu_limit_percentage INTEGER NOT NULL CHECK (cpu_limit_percentage >= 0),
    ram_limit_mb INTEGER NOT NULL CHECK (ram_limit_mb >= 0)
);

CREATE TABLE IF NOT EXISTS sessions (
    id STRING PRIMARY KEY,
    userID INTEGER CHECK (userID >= 0), -- Idk how to use FOREIGN KEY but it's meant to point to the users table's userID

);
