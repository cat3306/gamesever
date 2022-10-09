# 验签
生成公私钥
``` 
$ openssl genrsa -out private_key.pem 3072
```
``` 
$ openssl rsa -in private_key.pem -pubout -out public_key.pem
```