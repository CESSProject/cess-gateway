# <h1 align="center">CESS-GATEWAY &middot; [![GitHub license](https://img.shields.io/badge/license-Apache2-blue)](#LICENSE) <a href=""><img src="https://img.shields.io/badge/golang-%3E%3D1.16-blue.svg" /></a></h1>

cess-gateway is a service that using REST API specification for accessing CESS cloud storage.

# Reporting a Vulnerability

If you find out any system bugs or you have a better suggestions, Please send an email to frode@cess.one,
we are happy to communicate with you.

# Service address

* The block explorer address is http://139.224.19.104:3000/?rpc=ws%3A%2F%2F106.15.44.155%3A9948%2F#/accounts
* The address use for obtain free TCESS tokens on the testnet is http://47.243.82.77:9708/transfer

# Register a wallet

**Before using, you should refer to the following process to get the token and storage space**

* First apply for a wallet account in the block explorer, and then buy CESS tokens. If it is a test network, you can get it for free.
![createAccount](https://github.com/CESSProject/W3F-illustration/blob/main/httpservice/createAccount.PNG)

* Operate the purchase space interface in the block explorer to purchase space
![purchaseSpace](https://github.com/CESSProject/W3F-illustration/blob/main/httpservice/purchaseSpace.PNG)

# Authentication

The CESS-Gateway API uses bearer tokens to authenticate requests. 

Your tokens carry many privileges, so be sure to keep them secure! Do not share your *secret tokens* in publicly accessible locations such as a GitHub repository, client-side code, and so forth.

The bearer token is a cryptic string, usually generated by the server in response to a login request. The client must send this token in the `Authorization` header when making requests to protected resources:

| Authorization:"your token" |
| --------------------- |

## token information

| field           | description                              |
| :-------------- | ---------------------------------------- |
| UserId          | User proprietary ID                      |
| CreateUserTime  | The unix time when the user was created  |
| CreateTokenTime | The unix time when the token was created |
| ExpirationTime  | The unix time when token expires         |
| Mailbox         | User's email address                     |
| RandomCode      | Random code                              |

## Token generation

The token is generated by the CESS-Gateway service. Each CESS-Gateway service has its own pair of public and private keys. The public key is used for encrypt the token information, and then base64 is used for encode the ciphertext to obtain the final token.

# Configuration file

Before using cess-gateway, prepare a configuration file named "conf.toml", put it in the same directory as cess-gateway, and its contents are as follows:
```
#Cess chain address
ChainAddr     = ""
#The ip address that the cess-gateway service listens to
ServiceAddr   = ""
#The port number on which the cess-gateway service listens
ServicePort   = ""
#The address of the wallet account
AccountAddr   = ""
#Seed for wallet account
AccountSeed   = ""
#Email address
EmailAddress  = ""
#Email password
EmailPassword = ""
#Outgoing server address of SMTP service
EmailHost     = ""
#Outgoing server port number of SMTP service
EmailHostPort = 0
```
_You need to fill in your own information into the configuration file._

# CESS-Gateway HTTP API

The public API endpoint URL of CESS-Gateway is the server you deploy, All endpoints described in this document should be made relative to this root URL.

## User authorization

| **POST** /auth |
| -------------- |

The authorization interface is used to generate user tokens.

The user uploads his own email address and a captcha (the captcha should be written as 0 when uploading for the first time). If the email address is not authorized, CESS-Gateway will send a captcha to the mailbox. After the user gets the captcha in the mailbox, upload it again. CESS-Gateway generates a token for it and sends it to the mailbox，the validity of the token is 30 days.

If the mailbox uploaded by the user has been authorized, CESS-Gateway will check whether its token has expired. If it expires, CESS-Gateway will generate a new token and send it to its mailbox.

**Request Header**

| field        | value            |
| ------------ | ---------------- |
| Content-Type | application/json |

**Request Body**

| field   | value                         |
| ------- | ----------------------------- |
| mailbox | your mailbox address          |
| captcha | captcha received by mailbox   |

**Responses**

Response Schema: `application/json`

| status code               | structure                                 | description                                                  |
| ------------------------- | ----------------------------------------- | ------------------------------------------------------------ |
| 200 OK                    | code:200<br />msg:string<br />data:string | `msg` Enum:["success","captcha has expired and a new captcha has been sent to your mailbox","A new token has been sent to your mailbox"]<br />`data` Enum:["","token"] |
| 400 Bad Request           | code:400<br />msg:string                  | `msg` Enum:["HTTP error"，"Email format error"，"captcha error"，"Please check your email address and whether to enable SMTP service"] |
| 500 Internal Server Error | code:500<br />msg:string                  | `msg` Enum: ["Server internal data error"，"Server unexpected error"] |

**Request example**
```
curl "url"/auth -X POST -d '{"mailbox": "", "captcha": 0}' -H "Content-Type: application/json"
```

## Upload a file

| **PUT** /"file name" |
| -------------------- |

The put file interface is used for allow users to store files in the cess system.

You need to submit the file as form data and use provide the specific field.

**Request Header**

| field         | value   |
| ------------- | ------- |
| Authorization | token   |

**Request Body**

| field | value        |
| ----- | ------------ |
| file  | file[binary] |

**Responses**

Response Schema: `application/json`

| status code               | structure                | description                                                  |
| ------------------------- | ------------------------ | ------------------------------------------------------------ |
| 200 OK                    | code:200<br />msg:string | `msg` Default:"success                                       |
| 400 Bad Request           | code:400<br />msg:string | `msg` Default:"HTTP error"                                   |
| 401 Unauthorized          | code:401<br />msg:string | `msg` Enum:["Unauthorized"，"token expired"]                 |
| 403 Forbidden             | code:403<br />msg:string | `msg` Enum: ["duplicate filename"，"not enough space"，"The file is in hot backup, please try again later."] |
| 500 Internal Server Error | code:500<br />msg:string | `msg` Enum: ["Server internal data error"，"Server internal chain data error"，"Server unexpected error"] |

**Request example**
```
curl -X PUT "url"/test.log -F 'file=@test.log' -H "Authorization: "token""
```

## Download a file

| **GET** /"file name" |
| -------------------- |

The get file interface is used for get files in the CESS storage system. Currently, service only supported get files that upload by yourself.

**Request Header**

| field         | value   |
| ------------- | ------- |
| Authorization |  token  |

**Responses**

The response schema for the normal return status is: `application/octet-stream`

The response schema for the exception return status is: `application/json`, The message returned by the exception is as follows:

| status code               | structure                | description                                                  |
| ------------------------- | ------------------------ | ------------------------------------------------------------ |
| 400 Bad Request           | code:400<br />msg:string | `msg` Enum:["HTTP error"，"This file has not been uploaded"] |
| 401 Unauthorized          | code:401<br />msg:string | `msg` Enum:["Unauthorized"，"token expired"]                 |
| 403 Forbidden             | code:403<br />msg:string | `msg` Enum: ["Token is not valid","Token expired"]           |
| 500 Internal Server Error | code:500<br />msg:string | `msg` Enum: ["Server internal data error"，"Server internal chain data error"，"Server unexpected error"] |

**Request example**
```
curl -X GET "url"/test.log -H "Authorization: "token""
```

## Delete a file

The delete file interface is used for delete a put file.

| **DELETE** /"file name" |
| ----------------------- |

**Request Header**

| field         | value   |
| ------------- | ------- |
| Authorization |  token  |

**Responses**

Response Schema: `application/json`

| status code               | structure                | description                                                  |
| ------------------------- | ------------------------ | ------------------------------------------------------------ |
| 200 OK                    | code:200<br />msg:string | `msg` Default: "success"                                     |
| 400 Bad Request           | code:400<br />msg:string | `msg` Default: "HTTP error"                                  |
| 401 Unauthorized          | code:401<br />msg:string | `msg` Enum:["Unauthorized"，"token expired"]                 |
| 404 Not Found             | code:404<br />msg:string | `msg`Default: "This file has not been uploaded"              |
| 500 Internal Server Error | code:500<br />msg:string | `msg` Enum: ["Server internal data error"，"Server unexpected error"] |

**Request example**
```
curl -X DELETE "url"/test.log -H "Authorization: "token""
```

## List previous operation

| **GET** /files |
| -------------- |

List the previously put files, and display the 30 files closest to the current time by default. It also supports searching by page.

**Request Header**

| field         | value   |
| ------------- | ------- |
| Authorization |  token  |

**Query Parameters**

| field | description                                                  |
| ----- | ------------------------------------------------------------ |
| size  | type:<int32><br />default:30<br />Specifies the maximum number of uploads to return，up to 1000. |
| page  | type:<int32><br />default:0<br />Specifies the number of puts on which page to return. |

**Responses**

Response Schema: `application/json`

| status code               | structure                                   | description                                                  |
| ------------------------- | ------------------------------------------- | ------------------------------------------------------------ |
| 200 OK                    | code:200<br />msg:string<br />data:[]string | `msg` Default:"success"<br />`data`:[file name list]         |
| 400 Bad Request           | code:400<br />msg:string                    | `msg` Default: "HTTP error"                                  |
| 401 Unauthorized          | code:401<br />msg:string                    | `msg` Enum:["Unauthorized"，"token expired"]                 |
| 500 Internal Server Error | code:500<br />msg:string                    | `msg` Enum: ["Server internal data error"，"Server unexpected error"] |

**Request example**
```
curl -X GET "url"/files -H "Authorization: "token""
```
  
## License
Licensed under [Apache 2.0](https://github.com/CESSProject/cess-gateway/blob/main/LICENSE)
