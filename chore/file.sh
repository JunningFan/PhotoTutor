#！/in/bash
# This is curl for uploading file, which has same effect on the api spec
curl -X POST http://localhost:8000/upload/ \
 -F "upload=@/home/tecty/Pictures/upload.png" \
 -H "Content-Type: multipart/form-data"\
 -H "Authorization: ${1}"
