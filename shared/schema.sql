CREATE TABLE IF NOT EXISTS manuals (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	filename TEXT UNIQUE,
	model_series TEXT,
	revision TEXT
);
CREATE TABLE IF NOT EXISTS figures (
	manual_id INTEGER,
	id TEXT,
	PRIMARY KEY (manual_id, id),
	FOREIGN KEY(manual_id) REFERENCES manuals(id)
);
CREATE TABLE IF NOT EXISTS metadata (
	id INTEGER PRIMARY KEY,
	last_updated TEXT
);
CREATE TABLE IF NOT EXISTS parts (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	manual_id INTEGER,
	figure_id TEXT,
	key_no TEXT,
	part_number TEXT,
	base_part TEXT,
	revision TEXT,
	qty TEXT,
	description TEXT,
	remarks TEXT,
	FOREIGN KEY(manual_id) REFERENCES manuals(id)
);
CREATE INDEX IF NOT EXISTS idx_base_part ON parts(base_part);
CREATE TABLE IF NOT EXISTS staging_parts (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	manual_id INTEGER NOT NULL DEFAULT 0,
	figure_id TEXT NOT NULL DEFAULT '',
	key_no TEXT NOT NULL DEFAULT '',
	part_number TEXT NOT NULL DEFAULT '',
	base_part TEXT NOT NULL DEFAULT '',
	revision TEXT NOT NULL DEFAULT '',
	qty TEXT NOT NULL DEFAULT '',
	description TEXT NOT NULL DEFAULT '',
	remarks TEXT NOT NULL DEFAULT '',
	status TEXT NOT NULL DEFAULT 'pending',
	error_reason TEXT,
	staged_at TEXT NOT NULL DEFAULT (datetime('now'))
);
CREATE INDEX IF NOT EXISTS idx_staging_status ON staging_parts(status);
