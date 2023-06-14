use testdb;

CREATE TABLE data(
	ID INT PRIMARY KEY AUTO_INCREMENT,
    sender VARCHAR(100),
    receiver VARCHAR(100),
    message VARCHAR(10000),
    time_stamp INT
);

select * from data

    
