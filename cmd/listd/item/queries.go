package item

// PostgreSQL queries for the item table.
const (
	// selectAll is a query that selects all rows in the item table filtered
	// by list_id.
	selectAll = "SELECT * FROM item WHERE list_id = $1;"

	// selectByIDAndListID is a query that selects a row in the item table
	// filtered by item_id and list_id.
	selectByIDAndListID = "SELECT * FROM item WHERE item_id = $1 AND list_id = $2;"

	// insert is a query that inserts a row into the item table using the
	// values given in order for list_id, name, quantity, created, and
	// modified.
	insert = "INSERT INTO item (list_id, name, quantity, created, modified) VALUES ($1, $2, $3, $4, $5) RETURNING item_id;"

	// update is a query that updates a row in the item table based off of
	// item_id and list_id. The values able to be updated are name,
	// quantity, and modified.
	update = "UPDATE item SET name = $1, quantity = $2, modified = $3 WHERE item_id = $4 AND list_id = $5;"

	// del is a query that deletes a row in the item table given an item_id.
	del = "DELETE FROM item WHERE item_id = $1"
)
