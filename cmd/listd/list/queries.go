package list

// PostgreSQL queries for the list table and tables related to the list table through
// foreign keys, all used in the list package.
const (
	// selectAll is a query that selects all rows from the list table.
	selectAll = "SELECT * FROM list;"

	// selectByID is a query that selects a row from the list table based off of
	// the given list_id.
	selectByID = "SELECT * FROM list WHERE list_id = $1;"

	// insert is a query that inserts a new row in the list table using the values
	// given in order for name, created, and modified.
	insert = "INSERT INTO list (name, created, modified) VALUES ($1, $2, $3) RETURNING list_id;"

	// update is a query that updates a row in the list table based off of list_id.
	// The values able to be updated are name and modified.
	update = "UPDATE list SET name = $1, modified = $2 WHERE list_id = $3;"

	// delRelatedItems deletes rows in the item table that are related to a list by
	// a given list_id.
	delRelatedItems = "DELETE FROM item WHERE list_id = $1"

	// del is a query that deletes a row in the list table given a list_id.
	del = "DELETE FROM list WHERE list_id = $1;"
)
