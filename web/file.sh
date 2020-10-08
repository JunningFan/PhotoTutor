#！/in/bash

# This is curl for uploading file, which has same effect on the api spec
curl -X POST http://localhost:8080/pictures/upload/ \
 -F "upload=@/home/tecty/Pictures/upload.png" \
 -H "Content-Type: multipart/form-data"\
 -H "Authorization: eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJJRCI6MSwiRXhwaXJlIjoxNjAxNDYwMDk5fQ.H-qhxtZYzUek-QSgr3499jznTsbfaV5iRrKJE7X-7M3_N13t_CVYpqUtB9hskOVUI-lKo-1W8T8hR-y0FIaAHg"