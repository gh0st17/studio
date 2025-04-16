package database

const connStr = "user=studio " +
	"password=studio dbname=studio " +
	"sslmode=disable search_path=studio"

// SELECT запросы
const (
	loginQuery          = "SELECT access_level FROM users WHERE login = $1"
	fetchCustLoginQuery = "SELECT * FROM customers WHERE login = $1"
	fetchEmplLoginQuery = "SELECT * FROM employees WHERE login = $1"

	fetchOrdersQueryHead = "SELECT o.id, c.first_name || ' ' || c.last_name AS customer_name, " +
		"COALESCE(e.first_name || ' ' || e.last_name, '') AS employee_name, " +
		"o.status, (SELECT SUM(unit_price) FROM order_items WHERE o_id = o.id) AS total_price, " +
		"o.create_date, o.release_date " +
		"FROM orders o " +
		"LEFT JOIN customers c ON o.c_id = c.id " +
		"LEFT JOIN employees e ON o.e_id = e.id "

	fetchOrdersQueryCid  = fetchOrdersQueryHead + " WHERE c_id = $1 ORDER BY o.id"
	fetchOrdersQueryAll  = fetchOrdersQueryHead + " ORDER BY o.id"
	fetchOrderItemsQuery = "SELECT * FROM order_items WHERE c_id = $1 ORDER BY o.id"
	fetchOrderStatus     = "SELECT status FROM orders WHERE id = $1"

	fetchMatQuery      = "SELECT * FROM materials"
	fetchModelsQuery   = "SELECT id, title, price FROM models ORDER BY id"
	fetchModelMatQuery = "SELECT m.id, m.title, mm.leng, m.price " +
		"FROM model_materials mm " +
		"JOIN materials m ON mm.material_id = m.id " +
		"WHERE mm.model_id = $1"
)

// INSERT запросы
const (
	insertUserQuery = "INSERT INTO users (login,access_level) VALUES ($1,$2)"
	insertCustQuery = "INSERT INTO customers (first_name,last_name,login) VALUES ($1,$2,$3)"

	insertOrderQuery      = "INSERT INTO orders (c_id) VALUES ($1) RETURNING id"
	insertOrderItemsQuery = "INSERT INTO order_items (o_id, model, unit_price) VALUES ($1,$2,$3)"
)

// UPDATE запросы
const (
	updateStatusQuery    = "UPDATE orders SET status = $1 WHERE id= $2"
	updateOrderEmplQuery = "UPDATE orders SET e_id = $1 WHERE id = $2"
)
