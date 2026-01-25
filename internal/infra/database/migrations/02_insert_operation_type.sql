-- +goose up

insert into operation_types(operation_type_id, description) values(1, 'PURCHASE');
insert into operation_types(operation_type_id, description) values(2, 'INSTALLMENT PURCHASE');
insert into operation_types(operation_type_id, description) values(3, 'WITHDRAWAL');
insert into operation_types(operation_type_id, description) values(4, 'PAYMENT');

-- +goose down

delete from operation_types where operation_type_id = 1;
delete from operation_types where operation_type_id = 2;
delete from operation_types where operation_type_id = 3;
delete from operation_types where operation_type_id = 4;