CREATE VIEW vu_version AS
SELECT 1 AS version;

CREATE TABLE IF NOT EXISTS configs (
	key TEXT NOT NULL,
	value TEXT NOT NULL,
    PRIMARY KEY (key)
);

CREATE TABLE IF NOT EXISTS items (
	id INTEGER NOT NULL,
	pos INTEGER NOT NULL UNIQUE,
	color TEXT NOT NULL,
	item TEXT NOT NULL,
	subitem TEXT NULL,
    active INTEGER NOT NULL DEFAULT 1,
    PRIMARY KEY (id)
);

CREATE VIEW vu_items_lvl_1 AS
SELECT item, MIN(pos) AS pos
  FROM items 
 GROUP BY item;

CREATE TABLE IF NOT EXISTS checkouts (
    id TEXT NOT NULL,
    pos INTEGER NOT NULL UNIQUE,
    can_access_order_list_page INTEGER NOT NULL,
    can_access_stats_page INTEGER NOT NULL,
    can_access_stock_page INTEGER NOT NULL,
    can_access_console_page INTEGER NOT NULL,
    active INTEGER NOT NULL DEFAULT 1,
    password_hash INTEGER NULL,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS rules (
    profile TEXT NOT NULL,
    item TEXT NOT NULL,
    quantity_w1 INTEGER NULL,
    quantity_w2 INTEGER NULL,
    quantity_w3 INTEGER NULL,
    quantity_w4 INTEGER NULL,
    PRIMARY KEY (profile, item)
    --FOREIGN KEY(item) REFERENCES items(item)
);

CREATE VIEW vu_enabled_weeks AS
SELECT r.profile,
       EXISTS(SELECT 1 FROM rules r2 WHERE r.profile = r2.profile AND r2.quantity_w1 > 0) AS enabled_w1,
       EXISTS(SELECT 1 FROM rules r2 WHERE r.profile = r2.profile AND r2.quantity_w2 > 0) AS enabled_w2,
       EXISTS(SELECT 1 FROM rules r2 WHERE r.profile = r2.profile AND r2.quantity_w3 > 0) AS enabled_w3,
       EXISTS(SELECT 1 FROM rules r2 WHERE r.profile = r2.profile AND r2.quantity_w4 > 0) AS enabled_w4
  FROM rules r
 GROUP BY r.profile;

CREATE TABLE IF NOT EXISTS beneficiaries (
    id TEXT NOT NULL,
    profile TEXT NOT NULL,
    active INTEGER NOT NULL DEFAULT 1,
    PRIMARY KEY (id)
    -- FOREIGN KEY(profile) REFERENCES rules(profile)
);

CREATE TABLE IF NOT EXISTS stock (
	item_id INTEGER NOT NULL,
	quantity INTEGER NOT NULL,
    PRIMARY KEY (item_id),
    FOREIGN KEY(item_id) REFERENCES items(id)
);

CREATE TABLE IF NOT EXISTS orders (
	id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
	checkout_id TEXT NOT NULL,
	operator TEXT NOT NULL,
	beneficiary_id TEXT NOT NULL,
	datetime TEXT NOT NULL DEFAULT (DATETIME('now', 'localtime')),
	note TEXT NULL,
    active INTEGER NOT NULL DEFAULT 1,
    FOREIGN KEY(checkout_id) REFERENCES checkouts(id),
    FOREIGN KEY(beneficiary_id) REFERENCES beneficiaries(id)
);

CREATE INDEX idx__orders__active__datetime ON orders (active, datetime);

CREATE TABLE IF NOT EXISTS order_rows (
	order_id INTEGER NOT NULL,
	item_id INTEGER NOT NULL,
	quantity INTEGER NOT NULL CHECK (quantity > 0),
	PRIMARY KEY (order_id, item_id),
    FOREIGN KEY(order_id) REFERENCES orders(id),
    FOREIGN KEY(item_id) REFERENCES items(id)
);

CREATE TABLE IF NOT EXISTS sessions (
	checkout_id TEXT NOT NULL,
    operator TEXT NOT NULL,
	datetime TEXT NOT NULL DEFAULT (DATETIME('now', 'localtime')),
    PRIMARY KEY (checkout_id),
    FOREIGN KEY(checkout_id) REFERENCES checkouts(id)
);

CREATE VIEW vu_active_sessions AS
SELECT checkout_id, operator, datetime, 
       datetime > DATETIME('now', 'localtime', '-10 minutes') AS active
  FROM sessions
 WHERE datetime > DATETIME('now', 'localtime', '-24 hours');
