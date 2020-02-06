import random
from datetime import datetime, timedelta
from random import randrange

random_surnames = ['Hans', 'Peter', 'Karl', 'Heinz', 'Klaus', 'Werner', 'Wolfang', 'Michael', 'Helmut', 'Manfred', 'Josef', 'Horst', 
'Maria', 'Ursula', 'Elisabeth', 'Petra', 'Hildegard', 'Monika', 'Sabine', 'Gertrud', 'Helga', 'Gisela', 'Brigitte', 'Renate']
random_givennames = ['Müller', 'Fischer', 'Schmidt', 'Becker', 'Schmitz', 'Hoffmann', 'Klein', 'Meyer', 'Schneider', 'Koch', 'Jansen', 'Peters']

# Generate this number of entries:
number_of_entries = 5


def random_date(start, end):
    """
    This function will return a random datetime between two datetime 
    objects.
    """
    delta = end - start
    int_delta = (delta.days * 24 * 60 * 60) + delta.seconds
    random_second = randrange(int_delta)
    return start + timedelta(seconds=random_second)


with open('example_data.sql', 'w+', encoding="utf-8") as target:
    # First line of the statement for customers
    sql_str_customers = f'insert into Customers\nValues\n'

    # First line of the statement for Readings
    sql_str_readings = f'insert into Readings\nValues\n'

    for x in range(1, number_of_entries + 1):
        
        # Configure the delimiter for new lines 
        delimiter = ''
        if x == number_of_entries:
            delimiter = ';'
        else:
            delimiter = ','

        # Id, Surname, Givenname for the customer Table
        givenname = random.choice(random_givennames)
        surname = random.choice(random_surnames)
        sql_str_customers += f"({x}, '{givenname}', '{surname}'){delimiter}\n"

        # Add the Statement for the readings table, source: https://stackoverflow.com/questions/553303/generate-a-random-date-between-two-other-dates
        d1 = datetime.strptime('1/1/2008 1:30 PM', '%m/%d/%Y %I:%M %p')
        d2 = datetime.strptime('1/1/2020 4:50 AM', '%m/%d/%Y %I:%M %p')

        rnd_date = random_date(d1, d2)
        value = random.randint(0,150)
        # ID, Measure_Date, Value
        sql_str_readings += f"({x}, '{rnd_date}', {value}){delimiter}\n"

    target.write(sql_str_customers)
    target.write(sql_str_readings)
    target.close()
