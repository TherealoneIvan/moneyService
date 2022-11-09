CREATE TABLE users_balance
(
    id int  not null unique ,
    balance numeric not null CONSTRAINT positive_balance CHECK (balance >= 0),
    reserved_balance numeric not null DEFAULT 0 CONSTRAINT positive_reserve_balance CHECK (reserved_balance >= 0)
);

CREATE TABLE users_orders
(
    id int not null ,
    service_id int not null ,
    order_id   int not null unique ,
    order_cost numeric,
    date timestamp
);
CREATE TABLE success_orders
(
    id int not null,
    service_id int not null ,
    order_id   int not null unique ,
    order_cost numeric ,
    date timestamp
);
CREATE TABLE "all_events"
(
    id int not null ,
    event char(40),
    amount numeric,
    date timestamp
)
