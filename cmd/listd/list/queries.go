package list

// List Queries
const (
	// selectAll is a query that selects all rows
	selectAll = "SELECT * FROM list;"

	// selectByID is a query that selects a row based off of the list_id
	selectByID = "SELECT * FROM list WHERE list_id = $1;"

	// insert is a query that inserts a new list row
	insert = "INSERT INTO list (name, created, modified) VALUES ($1, $2, $3) RETURNING list_id;"

	// update is a query that updates a row based off of list_id
	update = "UPDATE list SET name = $1, modified = $2 WHERE list_id = $3;"

	// del is a query that deletes a row in the list table
	del = "DELETE FROM list WHERE list_id = $1;"
)
