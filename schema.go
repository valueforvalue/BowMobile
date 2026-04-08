package bow

const Schema = `
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
`
