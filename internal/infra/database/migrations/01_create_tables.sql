-- +goose up

-- ACCOUNTS

create table if not exists accounts
(
    account_id      bigserial primary key,
    document_number varchar not null unique
);

alter table if exists accounts
    owner to pismo;

-- OPERATION_TYPES

create table if not exists operation_types
(
    operation_type_id bigint primary key,
    description       varchar not null
);

alter table if exists operation_types
    owner to pismo;

-- TRANSACTIONS

create table if not exists transactions
(
    transaction_id    bigserial primary key,
    account_id        bigint                   not null references accounts (account_id),
    operation_type_id bigint                   not null references operation_types (operation_type_id),
    amount            double precision         not null,
    event_date        timestamp with time zone not null
);

alter table transactions
    owner to pismo;


-- +goose down
drop table if exists accounts;
drop table if exists transactions;
