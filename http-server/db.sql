use testdb;
delete from data;
alter table data auto_increment = 1;
#insert into data(sender,receiver,message) values(1,1,1);
#SELECT message from data WHERE (receiver = "b" AND sender = "a") OR (sender = "a" AND receiver = "b");
select * from data

    