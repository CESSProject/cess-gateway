# <h1 align="center">CESS-HTTPSERVICE &middot; [![GitHub license](https://img.shields.io/badge/license-Apache2-blue)](#LICENSE) <a href=""><img src="https://img.shields.io/badge/golang-%3E%3D1.16-blue.svg" /></a></h1>

cess-httpservice is a service using RESTful API specification for accessing cess cloud storage.

## Reporting a Vulnerability

If you find any bugs or good suggestions, Please send an email to tech@cess.one.
we are happy to communicate with you

## Service address

* The address of cess-httpservice is http://106.15.44.155:8081/, which has been replaced by `IP` below.
* The block explorer address is http://139.224.19.104:3000/?rpc=ws%3A%2F%2F106.15.44.155%3A9948%2F#/accounts
* The address for free access to tCESS tokens on the testnet is http://47.243.82.77:9708/transfer

## Usage for httpservice

**Before using, you should refer to the following process to get the token**

* First apply for a wallet account in the block explorer, and then buy CESS coins. If it is a test network, you can get it for free.
![createAccount](https://github.com/CESSProject/W3F-illustration/blob/main/httpservice/createAccount.PNG)

* Obtain random numbers (2) from the `IP/user/randoms` interface of httpservice through tools such as postman or curl,If everything works fine, you will get something like the following:
```
# curl IP/user/randoms -X POST -d '{"walletaddr": "your wallet address"}' --header "Content-Type: application/json"
{"code":200,"msg":"success","random1":116184,"random2":468019}
```

* Operate the user authentication interface in the block explorer, enter the deposit amount and the first random number, then click Submit transaction, and write down the block number of the transaction.
![userAuth](https://github.com/CESSProject/W3F-illustration/blob/main/httpservice/userAuth.PNG)

* Operate the purchase space interface in the block explorer to purchase space
![purchaseSpace](https://github.com/CESSProject/W3F-illustration/blob/main/httpservice/purchaseSpace.PNG)

* Operate the `IP/user/grant` interface of httpservice, and tell httpservice the block number and the second random number. If all data verification is successful, httpservice will return the user token information (as shown in the data field below), you need to save the token.

```
# curl IP/user/grant -X POST -d '{"blocknumber": your blocknumber, "walletaddr": "your wallet address", "random2": 468019}'  --header "Content-Type: application/json"
{"code":200,"msg":"success","data":"d6sZ2pXKcawXVMko8/JEaZoojuT9Wu7IElcdH9ayoV3neW/AkwGgSd7SZzwTt+Ll+K44DMCH1gOWbOZWu8UjeX0oWq0HoQxfPKguSNa6KMkBzMBYjEnDLZvCwUF7+KzN67zKhE2R9wn6OYsYRrv+KyAsVOGIkJxaP36tZwPAsg67ZsyTIU+O+fO4UXti9cwoWX27tBslbAWMiDyNDtGWKF1ggTueR4GoNSisQmL/jFBz2UhwpD4AH/KWLUaoi8BV+h5OoXbTM/0hRsC+g09z5293Qo+guEKi4fliwQG+0AG9mGQtefilnNkCXXYeuhkhk1NYIbqAVrAmcQt/OCE7kw=="}
```

**upload file**
```
# curl http://localhost:8081/file/upload -X POST --progress-bar  --form file=@test.txt --form token=d6sZ2pXKcawXVMko8/JEaZoojuT9Wu7IElcdH9ayoV3neW/AkwGgSd7SZzwTt+Ll+K44DMCH1gOWbOZWu8UjeX0oWq0HoQxfPKguSNa6KMkBzMBYjEnDLZvCwUF7+KzN67zKhE2R9wn6OYsYRrv+KyAsVOGIkJxaP36tZwPAsg67ZsyTIU+O+fO4UXti9cwoWX27tBslbAWMiDyNDtGWKF1ggTueR4GoNSisQmL/jFBz2UhwpD4AH/KWLUaoi8BV+h5OoXbTM/0hRsC+g09z5293Qo+guEKi4fliwQG+0AG9mGQtefilnNkCXXYeuhkhk1NYIbqAVrAmcQt/OCE7kw== --header "Content-Type: multipart/form-data" | tee /dev/null
############################################################################################# 100.0%
{"code":200,"msg":"success"}
```
**download file**
```
# curl -X GET -o goo.mod -# -G --data-urlencode "token= d6sZ2pXKcawXVMko8/JEaZoojuT9Wu7IElcdH9ayoV3neW/AkwGgSd7SZzwTt+Ll+K44DMCH1gOWbOZWu8UjeX0oWq0HoQxfPKguSNa6KMkBzMBYjEnDLZvCwUF7+KzN67zKhE2R9wn6OYsYRrv+KyAsVOGIkJxaP36tZwPAsg67ZsyTIU+O+fO4UXti9cwoWX27tBslbAWMiDyNDtGWKF1ggTueR4GoNSisQmL/jFBz2UhwpD4AH/KWLUaoi8BV+h5OoXbTM/0hRsC+g09z5293Qo+guEKi4fliwQG+0AG9mGQtefilnNkCXXYeuhkhk1NYIbqAVrAmcQt/OCE7kw==" --data-urlencode "filename=goo.mod" http://localhost:8081/file/download
################################################################################################## 100.0%
```
**View file list**
```
# curl -X GET -G --data-urlencode "token= d6sZ2pXKcawXVMko8/JEaZoojuT9Wu7IElcdH9ayoV3neW/AkwGgSd7SZzwTt+Ll+K44DMCH1gOWbOZWu8UjeX0oWq0HoQxfPKguSNa6KMkBzMBYjEnDLZvCwUF7+KzN67zKhE2R9wn6OYsYRrv+KyAsVOGIkJxaP36tZwPAsg67ZsyTIU+O+fO4UXti9cwoWX27tBslbAWMiDyNDtGWKF1ggTueR4GoNSisQmL/jFBz2UhwpD4AH/KWLUaoi8BV+h5OoXbTM/0hRsC+g09z5293Qo+guEKi4fliwQG+0AG9mGQtefilnNkCXXYeuhkhk1NYIbqAVrAmcQt/OCE7kw==" http://localhost:8081/file/list
{"code":200,"msg":"success","data":["test.log","test1.log"]}
```
**View user status**
```
# curl -X GET -G --data-urlencode "token= d6sZ2pXKcawXVMko8/JEaZoojuT9Wu7IElcdH9ayoV3neW/AkwGgSd7SZzwTt+Ll+K44DMCH1gOWbOZWu8UjeX0oWq0HoQxfPKguSNa6KMkBzMBYjEnDLZvCwUF7+KzN67zKhE2R9wn6OYsYRrv+KyAsVOGIkJxaP36tZwPAsg67ZsyTIU+O+fO4UXti9cwoWX27tBslbAWMiDyNDtGWKF1ggTueR4GoNSisQmL/jFBz2UhwpD4AH/KWLUaoi8BV+h5OoXbTM/0hRsC+g09z5293Qo+guEKi4fliwQG+0AG9mGQtefilnNkCXXYeuhkhk1NYIbqAVrAmcQt/OCE7kw==" http://localhost:8081/user/state
{
"code":200,
"msg":"success",
"data":{
"userId":1514176822592405504,
"deposit":100,
"totalSpace":10485760,
"usedSpace":714726,
"freeSpace":9771034,
"walletaddr":"5EWxDh3Gk5QKZBosoRCSmihagjKS5LVUm31CejeXxxLkdV1y"
}
}
```
**Delete File**
```
# curl http://localhost:8081/file/delete -X POST -d '{"token":" d6sZ2pXKcawXVMko8/JEaZoojuT9Wu7IElcdH9ayoV3neW/AkwGgSd7SZzwTt+Ll+K44DMCH1gOWbOZWu8UjeX0oWq0HoQxfPKguSNa6KMkBzMBYjEnDLZvCwUF7+KzN67zKhE2R9wn6OYsYRrv+KyAsVOGIkJxaP36tZwPAsg67ZsyTIU+O+fO4UXti9cwoWX27tBslbAWMiDyNDtGWKF1ggTueR4GoNSisQmL/jFBz2UhwpD4AH/KWLUaoi8BV+h5OoXbTM/0hRsC+g09z5293Qo+guEKi4fliwQG+0AG9mGQtefilnNkCXXYeuhkhk1NYIbqAVrAmcQt/OCE7kw==", "filename": "test1.log"}'  --header "Content-Type: application/json"
{"code":200,"msg":"success"}
```
**View real-time pricing for space**
```
# curl -X GET http://localhost:8081/space/price
{"code":200,"msg":"success","data":0}
```

**regrant**
```
# curl http://localhost:8081/user/regrant -X POST --form token= d6sZ2pXKcawXVMko8/JEaZoojuT9Wu7IElcdH9ayoV3neW/AkwGgSd7SZzwTt+Ll+K44DMCH1gOWbOZWu8UjeX0oWq0HoQxfPKguSNa6KMkBzMBYjEnDLZvCwUF7+KzN67zKhE2R9wn6OYsYRrv+KyAsVOGIkJxaP36tZwPAsg67ZsyTIU+O+fO4UXti9cwoWX27tBslbAWMiDyNDtGWKF1ggTueR4GoNSisQmL/jFBz2UhwpD4AH/KWLUaoi8BV+h5OoXbTM/0hRsC+g09z5293Qo+guEKi4fliwQG+0AG9mGQtefilnNkCXXYeuhkhk1NYIbqAVrAmcQt/OCE7kw== --header "Content-Type: multipart/form-data"
{
"code":200,
"msg":"success",
"data":"pgPsILT9+YrdccJbOs9XInWHf21cO/E9JdpvW5Xmh79s3YBeKBBFoHWUt10CDspjc2EK5QvQWkk+nQPGVzBY5CWzEI3stsKYDI9YsmnxkFshvAgU8S1bzwjeAJowTQEJalIORgzoTQ3442gj5aXYzdTy10o5iBDru4kWFAg0LS/ajJ43Pc7lo2N9fxiKZ80vrfrQT3mg08Wtn3H2GNzczP9JMgEasjwOW7JgO5K71GkCe/E6ub9YqoQOMXz0XnfqOgrxv3fBha+A66NT3DDxi5fp3kqnIGlkV81hOxlmMolFJ2H/ZTGoFBwFZyKt+UtI6zHXijF1F7+/TwyyUnnAOA=="
}
```

## License
Licensed under [Apache 2.0](https://github.com/CESSProject/cess-httpservice/blob/main/LICENSE)
