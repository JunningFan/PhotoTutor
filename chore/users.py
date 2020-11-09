#!/usr/bin/python3
import requests
import random
from requests_toolbelt.multipart.encoder import MultipartEncoder
# BASE_URL = "http://localhost:8080/"
BASE_URL = "http://localhost:8000/user/"
# BASE_URL = "http://whiteboard.house:8000/"
access = []

# for i in range(10):
#     name = "ttt%d" % i
#     response = requests.post(
#         BASE_URL + "",
#         json={
#             "username": name,
#             "password": "apple123",
#             "nickname": name
#         }
#     )
#     access.append(response.json()['access'])
for i in range(10):
    name = "ttt%d" % i
    response = requests.post(
        BASE_URL + "login/",
        json={
            "username": name,
            "password": "apple123",
            "nickname": name
        }
    )
    access.append(response.json()['access'])


for i in range(3):
    for j in range(10):

        to = random.randint(1, 10)
        print("%d -> %d " % (j, to))
        res = requests.post(
            BASE_URL + "follow/%d" % to,
            headers={
                "Authorization": access[j]
            }
        )
        if not res.ok:
            print(res)
            raise res
