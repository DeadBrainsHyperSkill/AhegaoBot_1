CREATE TABLE customers(
	customer_id char(18) PRIMARY KEY,
	balance integer NOT NULL,
    spent integer NOT NULL,
    bio varchar(100),
    lootbox_amount int NOT NULL,
    next_farm bigint NOT NULL,
    voice_time bigint NOT NULL,
    text_channel_id char(18),
    voice_channel_id char(18),
    role_id char(18),
    channels_expires bigint NOT NULL
);

CREATE TABLE products(
    product_id SERIAL PRIMARY KEY,
    product_dial int UNIQUE NOT NULL,
    product_name varchar(100) NOT NULL,
    product_type varchar(100) NOT NULL,
	role_id char(18) UNIQUE,
	price integer NOT NULL CHECK (price >= 0)
);

CREATE TABLE orders(
	order_id SERIAL PRIMARY KEY,	
	customer_id char(18) NOT NULL,
    product_id SERIAL NOT NULL,
    is_hidden boolean NOT NULL DEFAULT FALSE,
	expires bigint NOT NULL
);


CREATE TABLE voicelog(
    voicelog_id serial PRIMARY KEY,
    customer_id char(18) NOT NULL,
    joined_at bigint NOT NULL
);

INSERT into customers VALUES('903076261558108221', 0, 0, '', 0, 0, 0, NULL, NULL, NULL, 0);
INSERT into customers VALUES('474149886040997908', 0, 0, '', 0, 0, 0, NULL, NULL, NULL, 0);
INSERT into customers VALUES('948749513210855434', 0, 0, '', 0, 0, 0, NULL, NULL, NULL, 0);
INSERT into customers VALUES('150580008912420865', 0, 0, '', 0, 0, 0, NULL, NULL, NULL, 0); 
INSERT into customers VALUES('876179504387739649', 0, 0, '', 0, 0, 0, NULL, NULL, NULL, 0);
INSERT into customers VALUES('257126441446014977', 0, 0, '', 0, 0, 0, NULL, NULL, NULL, 0);