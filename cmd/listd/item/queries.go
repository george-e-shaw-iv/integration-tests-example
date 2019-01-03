package item

// Item Queries
const (
	// selectAll is a query that selects all rows filtered by list_id
	selectAll = "SELECT * FROM item WHERE list_id = $1;"

	// selectByIDAndListID is a query that selects a row filtered by item_id
	// and list_id
	selectByIDAndListID = "SELECT * FROM item WHERE item_id = $1 AND list_id = $2;"

	// insert is a query that inserts a row into the item table
	insert = "INSERT INTO item (list_id, name, quantity, created, modified) VALUES ($1, $2, $3, $4, $5) RETURNING item_id;"

	// update is a query that updates a row in the item table
	update = "UPDATE item SET name = $1, quantity = $2, modified = $3 WHERE item_id = $4 AND list_id = $5;"

	// del is a query that deletes a row in the item table
	del = "DELETE FROM item WHERE item_id = $1"
)
