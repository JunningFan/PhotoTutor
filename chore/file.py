#!/usr/bin/python3
import requests
from requests_toolbelt.multipart.encoder import MultipartEncoder
BASE_URL = "http://localhost:8000/"
# BASE_URL = "http://whiteboard.house:8000/"

response = requests.post(
    BASE_URL + "user/login/",
    json={
        "username": "test",
        "password": "test"
    }
)

access = response.json()['access']

encoder = MultipartEncoder(fields={
    'upload': ('upload.jpeg', open('./avatar.jpeg', 'rb'))
})

response = requests.post(
    BASE_URL + "upload/",
    data=encoder,
    headers={
        "Authorization": access,
        'Content-Type': encoder.content_type
    },
)

# print(response.request.body)
print(response.json())
