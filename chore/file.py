#!/usr/bin/python3
import requests
from requests_toolbelt.multipart.encoder import MultipartEncoder

response = requests.post(
    "http://localhost:3000/user/login/",
    json={
        "username": "ttt1",
        "password": "apple123"
    }
)

access = response.json()['access']

encoder = MultipartEncoder(fields={
    'upload': ('upload.png', open('/home/tecty/Pictures/upload.png', 'rb'))
})

response = requests.post(
    "http://localhost:3000/upload/",
    data=encoder,
    headers={
        "Authorization": access,
        'Content-Type': encoder.content_type
    },
)

# print(response.request.body)
print(response.json())
