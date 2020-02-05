CREATE TABLE IF NOT EXISTS "Readings" (
	"Customers_ID_FK"	INTEGER,
	"Measure_Date"	TEXT NOT NULL,
	"Value"	INTEGER NOT NULL,
	PRIMARY KEY("Customers_ID_FK"),
	FOREIGN KEY("Customers_ID_FK") REFERENCES "Customers"("Customers_ID")
);
CREATE TABLE IF NOT EXISTS "Customers" (
	"Customers_ID"	INTEGER PRIMARY KEY AUTOINCREMENT,
	"Surname"	TEXT NOT NULL,
	"Givenname"	TEXT NOT NULL
);
