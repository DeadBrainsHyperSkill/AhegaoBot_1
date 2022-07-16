package discordmarket

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Database struct {
	Conn *pgxpool.Pool
}

//CRoleAdd adds new custom role to the prodcuts table.
func (db Database) CustomRoleAdd(roleID, name, customerID string, expires int64) error {
	query := `INSERT INTO products(product_id, product_dial, product_name, product_type, role_id, price)
    VALUES((SELECT MAX(product_dial) FROM products)+1, (SELECT MAX(product_dial) FROM products)+1, $1, 'customrole', $2, 0) RETURNING product_dial`

	ctx := context.Background()

	var productID int

	err := db.Conn.QueryRow(ctx, query, name, roleID).Scan(&productID)
	if err != nil {
		return err
	}

	query = "INSERT INTO orders(customer_id, product_id, expires) VALUES ($1, $2, $3)"
	_, err = db.Conn.Exec(ctx, query, customerID, productID, expires)
	return err
}

// CustomerAdd adds a new customer to the customers table.
func (db Database) CustomerAdd(ID string) error {
	query := "INSERT INTO customers(customer_id, balance, spent, bio, lootbox_amount, next_farm, voice_time, text_channel_id, voice_channel_id, role_id, channels_expires) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) ON CONFLICT DO NOTHING"

	ctx := context.Background()

	_, err := db.Conn.Exec(ctx, query, ID, 0, 0, "", 0, 0, 0, nil, nil, nil, 0)
	return err
}

// CustomerProfileData returns customer's data for profile.
func (db Database) CustomerProfileData(ID string) (bio string, voiceT int64, spent, lbxs int, err error) {
	query := "SELECT COALESCE(bio, ''), voice_time, spent, lootbox_amount FROM customers WHERE customer_id=$1"

	ctx := context.Background()

	err = db.Conn.QueryRow(ctx, query, ID).Scan(&bio, &voiceT, &spent, &lbxs)
	return
}

// CustomerStatusUpdate updates stats	 in the customers table .
func (db Database) CustomerStatusUpdate(ID, bio string, price int) error {
	query := "UPDATE customers SET bio=$1, balance=balance-$2, spent=spent+$2 WHERE customer_id=$3"

	ctx := context.Background()

	_, err := db.Conn.Exec(ctx, query, bio, price, ID)
	return err
}

// CustomerTotalVoiceUpdate updates voice_time in the customers table.
func (db Database) CustomerTotalVoiceUpdate(ID string, deltaM int64) error {
	query := "UPDATE customers SET voice_time=voice_time+$1 WHERE customer_id=$2"

	ctx := context.Background()

	_, err := db.Conn.Exec(ctx, query, deltaM, ID)
	return err
}

// TopTenSpenders returns spender slice with top 10 customers by spent field.
func (db Database) TopTenSpenders() ([]spender, error) {

	query := "SELECT customer_id, spent FROM customers ORDER BY spent DESC LIMIT 10"

	ctx := context.Background()

	rows, err := db.Conn.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	var sps []spender
	for rows.Next() {
		var s spender
		err = rows.Scan(&s.UserID, &s.Spent)
		if err != nil {
			return sps, err
		}
		sps = append(sps, s)
	}

	return sps, rows.Err()
}

// TopTenVoice returns voicer slice with top 10 customers by voice_time field.
func (db Database) TopTenVoice() ([]voicer, error) {
	query := "SELECT customer_id, voice_time FROM customers ORDER BY spent DESC LIMIT 10"

	ctx := context.Background()

	rows, err := db.Conn.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	var vs []voicer
	for rows.Next() {
		var v voicer
		err = rows.Scan(&v.UserID, &v.Spent)
		if err != nil {
			return vs, err
		}
		vs = append(vs, v)
	}

	return vs, rows.Err()
}

// Balance returns balance of the customer from the customers table.
func (db Database) Balance(ID string) (int, error) {

	query := "SELECT balance FROM customers WHERE customer_id=$1"
	var balance int

	ctx := context.Background()

	row := db.Conn.QueryRow(ctx, query, ID)
	err := row.Scan(&balance)

	return balance, err
}

func (db Database) HasCustomerChannels(ID string) (textChannelID string, voiceChannelID string, roleID string, channelsExpires int64, err error) {

	query := "SELECT COALESCE(text_channel_id, ''), COALESCE(voice_channel_id, ''), COALESCE(role_id, ''), channels_expires FROM customers WHERE customer_id=$1"

	ctx := context.Background()

	err = db.Conn.QueryRow(ctx, query, ID).Scan(&textChannelID, &voiceChannelID, &roleID, &channelsExpires)
	return
}

func (db Database) CustomerChannelsAdd(ID string, textChannelID string, voiceChannelID string, roleID string, hasCustomerChannels bool) error {

	ctx := context.Background()

	var err error
	if hasCustomerChannels {
		query := "UPDATE customers SET channels_expires=channels_expires+$1 WHERE customer_id=$2"
		_, err = db.Conn.Exec(ctx, query, dayAsSec*180, ID)
	} else {
		query := "UPDATE customers SET text_channel_id=$1, voice_channel_id=$2, role_id=$3, channels_expires=channels_expires+$4 WHERE customer_id=$5"
		_, err = db.Conn.Exec(ctx, query, textChannelID, voiceChannelID, roleID, time.Now().Unix()+dayAsSec*180, ID)
	}
	return err
}

func (db Database) IsHiddenUpdate(ID string, isHide bool, roleDial int) error {
	query := `UPDATE orders SET is_hidden=$1 FROM products 
	WHERE orders.product_id = products.product_id AND orders.customer_id=$2 AND products.product_dial=$3`

	ctx := context.Background()

	_, err := db.Conn.Exec(ctx, query, isHide, ID, roleDial)
	return err
}

func (db Database) HasRoles(ID string) (bool, error) {
	query := "SELECT customer_id FROM orders WHERE customer_id=$1 AND product_id <= 25"

	ctx := context.Background()

	var customerID string
	err := db.Conn.QueryRow(ctx, query, ID).Scan(&customerID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// BalanceUpdate updates balance of the user in the customers table.
func (db Database) BalanceUpdate(ID string, value int) error {
	query := "UPDATE customers SET balance=balance+$1 WHERE customer_id=$2"
	if value < 0 {
		query = "UPDATE customers SET balance=balance+$1, spent=spent-$1 WHERE customer_id=$2"
	}

	ctx := context.Background()

	_, err := db.Conn.Exec(ctx, query, value, ID)
	return err
}

// LootboxAmount returns lootbox_amount from the customers table.
func (db Database) LootboxAmount(ID string) (int, error) {
	query := "SELECT lootbox_amount FROM customers WHERE customer_id=$1"

	ctx := context.Background()

	var amount int
	err := db.Conn.QueryRow(ctx, query, ID).Scan(&amount)
	return amount, err
}

// LootboxAmount increases lootbox_amount and changes balance and spent values depending on the price in customers table.
func (db Database) LootboxAdd(ID string, amount int, price int) error {

	query := "UPDATE customers SET balance=balance-$1, spent=spent+$1, lootbox_amount=lootbox_amount+$2 WHERE customer_id=$3"

	ctx := context.Background()

	_, err := db.Conn.Exec(ctx, query, price, amount, ID)
	return err
}

// LootboxAmount decreases lootbox_amount by 1 in the customers table.
func (db Database) LootboxRemove(ID string, amount int) error {
	query := "UPDATE customers SET lootbox_amount=lootbox_amount-$1 WHERE customer_id=$2"

	ctx := context.Background()

	_, err := db.Conn.Exec(ctx, query, amount, ID)
	return err
}

// NextFarm retuens next_farm from the customers table.
func (db Database) NextFarm(ID string) (int64, error) {
	query := "SELECT next_farm FROM customers WHERE customer_id=$1"

	ctx := context.Background()

	var nextFarmTime int64
	row := db.Conn.QueryRow(ctx, query, ID)
	err := row.Scan(&nextFarmTime)
	return nextFarmTime, err
}

// Farm increases balance and updates next_farm in in the customers table.
func (db Database) Farm(ID string, value int, nextFarmTime int64) error {
	query := "UPDATE customers SET balance=balance+$1, next_farm=$2 WHERE customer_id=$3"

	ctx := context.Background()

	_, err := db.Conn.Exec(ctx, query, value, nextFarmTime, ID)
	return err
}

// Product returns prodcut struct.
func (db Database) Product(productDial int) (*product, error) {
	query := "SELECT product_id, product_dial, product_name, product_type, COALESCE(role_id, ''), price FROM products WHERE product_dial=$1"

	ctx := context.Background()

	var p product
	row := db.Conn.QueryRow(ctx, query, productDial)
	err := row.Scan(&p.ID, &p.Dial, &p.Name, &p.Type, &p.RoleID, &p.Price)
	return &p, err
}

// Products returns a slice of product struct containing all prodcuts
func (db Database) Products() ([]product, error) {
	query := "SELECT product_id, product_dial, product_name, product_type, COALESCE(role_id, ''), price FROM products ORDER BY product_dial"

	ctx := context.Background()

	rows, err := db.Conn.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	var products []product
	for rows.Next() {
		var p product
		err = rows.Scan(&p.ID, &p.Dial, &p.Name, &p.Type, &p.RoleID, &p.Price)
		if err != nil {
			return products, err
		}
		products = append(products, p)
	}

	if rows.Err() != nil {
		return products, err
	}

	return products, nil
}

// HasProduct checks if user owns a particular product.
func (db Database) HasProduct(userID string, productDial int) (bool, error) {
	query := `SELECT o.order_id FROM orders o
	JOIN products p
	ON o.product_id = p.product_id
	WHERE o.customer_id=$1 AND p.product_dial=$2`

	ctx := context.Background()

	var product_id int
	err := db.Conn.QueryRow(ctx, query, userID, productDial).Scan(&product_id)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return false, nil
		}
		fmt.Println(err.Error())
		return false, err
	}
	return true, nil
}

// UserOrders returns slice of inventoryItem struct containing all items owned by the user.

func (db Database) UserOrders(userID string) ([]order, error) {
	query := `SELECT p.product_id, p.product_dial, p.product_name, p.product_type, COALESCE(p.role_id, ''), p.price, o.order_id, o.customer_id, o.product_id, o.is_hidden, o.expires 
	FROM orders o
	JOIN products p
	ON o.product_id = p.product_id
	WHERE o.customer_id=$1`

	ctx := context.Background()

	rows, err := db.Conn.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []order
	for rows.Next() {
		var p product
		var i order
		err = rows.Scan(&p.ID, &p.Dial, &p.Name, &p.Type, &p.RoleID, &p.Price, &i.OrderID, &i.CustomerID, &p.ID, &i.IsHidden, &i.Expires)
		if err != nil {
			return items, err
		}
		i.Product = &p
		items = append(items, i)
	}

	if rows.Err() != nil {
		return items, err
	}

	return items, nil
}

// Order is a transaction that updates user's balance and adds a product intp his possession.
// If user already has the prodcut, will update expiry time.
func (db Database) Order(userID string, p *product, duration int64) error {

	query1 := "UPDATE customers SET balance=balance-$1, spent=spent+$1 WHERE customer_id=$2"

	ctx := context.Background()

	tx, err := db.Conn.Begin(context.Background())
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, query1, p.Price, userID)
	if err != nil {
		rlbErr := tx.Rollback(ctx)
		if rlbErr != nil {
			return rlbErr
		}
		return err
	}

	hasProd, err := db.HasProduct(userID, p.Dial)
	if err != nil {
		rlbErr := tx.Rollback(ctx)
		if rlbErr != nil {
			return rlbErr
		}
		return err
	}

	if duration == -1 {
		query2 := "INSERT INTO orders(customer_id, product_id, expires) VALUES ($1, $2, $3)"
		_, err = tx.Exec(ctx, query2, userID, p.ID, duration)
	} else if hasProd {
		query2 := "UPDATE orders SET expires=expires+$1 WHERE customer_id=$2 and product_id=$3"
		_, err = tx.Exec(ctx, query2, duration, userID, p.ID)

	} else {
		query2 := "INSERT INTO orders(customer_id, product_id, expires) VALUES ($1, $2, $3)"
		_, err = tx.Exec(ctx, query2, userID, p.ID, time.Now().Unix()+duration)
	}

	if err != nil {
		rlbErr := tx.Rollback(ctx)
		if rlbErr != nil {
			return rlbErr
		}
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}
func (db Database) OrdersDeleteExpired() ([]expiredRole, error) {
	query := `DELETE FROM orders o 
	USING products p
	WHERE p.product_id=o.product_id AND o.expires > 0 AND o.expires <=$1
	RETURNING o.customer_id, COALESCE(p.role_id, '')`

	ctx := context.Background()

	var es []expiredRole
	rows, err := db.Conn.Query(ctx, query, time.Now().Unix())

	for rows.Next() {
		var e expiredRole
		err = rows.Scan(&e.UserID, &e.RoleID)
		if err != nil {
			return es, err
		}

		if e.RoleID != "" {
			es = append(es, e)
		}
	}

	if rows.Err() != nil {
		return es, err
	}
	return es, nil
}

func (db Database) CustomerChannelsExpired() (textChannels, voiceChannels, roles, usersID []string, err error) {
	timeNow := time.Now().Unix()

	query := "SELECT text_channel_id, voice_channel_id, role_id, customer_id FROM customers WHERE channels_expires > 0 AND channels_expires <=$1"
	query2 := "UPDATE customers SET text_channel_id=NULL, voice_channel_id=NULL, role_id=NULL, channels_expires=0 WHERE channels_expires > 0 AND channels_expires <=$1"

	ctx := context.Background()

	rows, err := db.Conn.Query(ctx, query, timeNow)
	if err != nil {
		return
	}
	for rows.Next() {
		var textChannel string
		var voiceChannel string
		var role string
		var userID string
		err = rows.Scan(&textChannel, &voiceChannel, &role, &userID)
		if err != nil {
			return textChannels, voiceChannels, roles, usersID, err
		}
		textChannels = append(textChannels, textChannel)
		voiceChannels = append(voiceChannels, voiceChannel)
		usersID = append(usersID, userID)
		roles = append(roles, role)
	}
	if len(roles) > 0 {
		_, err := db.Conn.Query(ctx, query2, timeNow)
		if err != nil {
			return textChannels, voiceChannels, roles, usersID, err
		}
	}
	return textChannels, voiceChannels, roles, usersID, nil
}

// Adds entry to the voicelog table.
func (db Database) VoiceLogAdd(userID string, joined_at int64) error {
	query := "INSERT INTO voicelog(customer_id, joined_at) VALUES($1,$2)"

	ctx := context.Background()

	_, err := db.Conn.Exec(ctx, query, userID, joined_at)
	if err != nil {
		return err
	}

	return nil
}

// VoiceLogGetAndClean returns entry from the voicelog table and deletes it.
func (db Database) VoiceLogGetAndClean(userID string) (int64, error) {

	query := "DELETE FROM voicelog WHERE customer_id = $1 RETURNING joined_at"

	ctx := context.Background()

	var t int64
	err := db.Conn.QueryRow(ctx, query, userID).Scan(&t)
	if err != nil {
		return 0, err
	}

	return t, nil
}
