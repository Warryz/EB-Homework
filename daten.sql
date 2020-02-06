DROP DATABASE IF EXISTS hausarbeit;
CREATE DATABASE hausarbeit; USE hausarbeit;

CREATE TABLE IF NOT EXISTS Customers (
	Customers_ID INT UNSIGNED AUTO_INCREMENT,
	Surname VARCHAR(60) NOT NULL,
	Givenname VARCHAR(60) NOT NULL,
	PRIMARY KEY(Customers_ID)
);

CREATE TABLE IF NOT EXISTS Readings (
	Customers_ID_FK INT UNSIGNED,
	Measure_ID INT UNSIGNED,
	Measure_Date DATETIME NOT NULL,
	Value INT NOT NULL, 
	PRIMARY KEY(Customers_ID_FK, Measure_ID), 
	FOREIGN KEY(Customers_ID_FK) REFERENCES Customers(Customers_ID)
);