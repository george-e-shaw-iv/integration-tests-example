package db

// schema is the constant that contains the postgres database schema
const schema = `
CREATE TABLE IF NOT EXISTS list (
	list_id SERIAL PRIMARY KEY,
	name varchar(255) NOT NULL UNIQUE,
	created timestamp NOT NULL DEFAULT NOW(),
	modified timestamp NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS item (
	item_id SERIAL PRIMARY KEY,
	list_id int NOT NULL,
	name varchar(255) NOT NULL,
	quantity int NOT NULL,
	created timestamp NOT NULL DEFAULT NOW(),
	modified timestamp NOT NULL DEFAULT NOW(),
	FOREIGN KEY(list_id) REFERENCES list(list_id)
);`
