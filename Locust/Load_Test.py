'''
Accessing the website by locust -f .\Load_Test.py --host=http://192.168.2.15
'''
from locust import HttpLocust, TaskSet, between
from random import randrange
print(randrange(10))

# Opening the website http://example.com/Hausarbeit
def index(l):
    l.client.get(f'/customerdata/{randrange(1, 4500)}')
    # print("Accessing the website.")

class UserBehavior(TaskSet):
    tasks = {index: 1}

    # Things to do before doing anything else.


class WebsiteUser(HttpLocust):
    task_set = UserBehavior
    # Definition of the user behavior: Wait at least 5 seconds and maximum 9 seconds. 
    wait_time = between(5.0, 9.0)
