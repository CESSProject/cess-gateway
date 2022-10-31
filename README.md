# <h1 align="center">CESS-GATEWAY &middot; [![GitHub license](https://img.shields.io/badge/license-Apache2-blue)](#LICENSE) <a href=""><img src="https://img.shields.io/badge/golang-%3E%3D1.19-blue.svg"/></a> [![Go Reference](https://pkg.go.dev/badge/github.com/CESSProject/cess-gateway.svg)](https://pkg.go.dev/github.com/CESSProject/cess-gateway)</h1>

CESS-Gateway is a service that using REST API specification for accessing CESS cloud storage.


## Reporting a Vulnerability
If you find out any system bugs or you have a better suggestions, Please send an email to frode@cess.one,
we are happy to communicate with you.


## System Requirements
- Linux-amd64


## Build from source

### Step 1: Install common libraries

Take the ubuntu distribution as an example:

```shell
sudo apt update && sudo upgrade
sudo apt install make gcc git curl wget vim util-linux -y
```

### Step 2: Install go locale

CESS-Gateway requires [Go1.19](https://golang.org/dl/) or higher.
> See the [official Golang installation instructions](https://golang.org/doc/install) If you get stuck in the following process.

- Download go1.19 compress the package and extract it to the /use/local directory:

```shell
sudo wget -c https://golang.org/dl/go1.19.linux-amd64.tar.gz -O - | sudo tar -xz -C /usr/local
```

- You'll need to add `/usr/local/go/bin` to your path. For most Linux distributions you can run something like:

```shell
echo "export PATH=$PATH:/usr/local/go/bin" >> ~/.bashrc && source ~/.bashrc
```

- View your go version:

```shell
go version
```

### Step 3: Build a gateway

```
git clone https://github.com/CESSProject/cess-gateway.git
cd cess-gateway/
go build -o gateway cmd/main.go
```
If all goes well, you will get a program called `gateway`.

## Get started with gateway

### Step 1: Register a polka wallet

Browser access: [App](https://testnet-rpc.cess.cloud/explorer) implemented by [CESS Explorer](https://github.com/CESSProject/cess-explorer), and [add an account](https://github.com/CESSProject/W3F-illustration/blob/main/gateway/createAccount.PNG).

### Step 2: Recharge your polka wallet

If you are using the test network, Please join the [CESS discord](https://discord.gg/mYHTMfBwNS) to get it for free. If you are using the official network, please buy CESS tokens.

### Step 3: Prepare configuration file

Use `gateway` to generate configuration file templates directly in the current directory:
```shell
sudo chmod +x gateway
./gateway default
```
The content of the configuration file template is as follows. You need to fill in your own information into the file. By default, the `gateway` uses conf.toml in the current directory as the runtime configuration file. You can use `-c` or `--config` to specify the configuration file Location.

```toml
#The rpc address of the chain node
RpcAddr           = ""
#The port number on which the cess-gateway service listens
ServicePort       = "8081"
#Phrase or seed for wallet account
AccountSeed       = ""
#Email address
EmailAddress      = ""
#Email authorization code
AuthorizationCode = ""
#Outgoing server address of SMTP service
SMTPHost          = ""
#Outgoing server port number of SMTP service
SMTPPort          = 0
```
*Our testnet rpc address is as follows:*<br>
`wss://testnet-rpc0.cess.cloud/ws/`<br>
`wss://testnet-rpc1.cess.cloud/ws/`

### Step 4: Buy space package for your account
There are five types of space packages, represented by 1 to 5.
- Package 1 means purchasing 10GiB space and the cost is 0;
- Package 2 means purchasing 500GiB space;
- Package 3 means purchasing 1TiB space;
- Package 4 means purchasing 5TiB space;
- Package 5 means that the purchased space exceeds 5TiB, and an integer greater than 5 needs to be specified;

**Buy space operation:**

```shell
./gateway buy [1,2,3,4,5] [6~]
```

### Step 5: Start the gateway service

```shell
sudo nohup ./gateway run 2>&1 &
```

## Other usage guidelines for gateway
### Upgrade space package
The space package can only be upgraded from low-level to high-level, and cannot be downgraded.
Take the upgrade of package 1 to package 2 as an example:
```
./gateway upgrade 2
```

### Space Package Renewal
By default, the space package is only valid for 1 month, and each renewal will add a month of validity.
```
./gateway renewal
```

### View space package details
```
./gateway space
```

# Usage for gateway API

The public API endpoint URL of CESS-Gateway is the server you deploy, All endpoints described in this document should be made relative to this root URL,The following example uses URL instead.

## Authentication

The CESS-Gateway API uses bearer tokens to authenticate requests. 

Your tokens carry many privileges, so be sure to keep them secure! Do not share your *secret tokens* in publicly accessible locations such as a GitHub repository, client-side code, and so forth.

The bearer token is a cryptic string, usually generated by the server in response to a login request. The client must send this token in the `Authorization` header when making requests to protected resources:

| Authorization: token  |
| --------------------- |

- Token information

| field           | description                              |
| :-------------- | ---------------------------------------- |
| UserId          | User ID                                  |
| CreateUserTime  | The unix time when the user was created  |
| CreateTokenTime | The unix time when the token was created |
| ExpirationTime  | The unix time when token expires         |
| Mailbox         | User's email address                     |
| RandomCode      | Random code                              |

- Token generation

The token is generated by the CESS-Gateway service. Each CESS-Gateway service has its own pair of public and private keys. The public key is used to encrypt the token information, and then base64 is used to encode the ciphertext to obtain the final token.

## Get token

| **POST** /auth |
| -------------- |

The authorization interface is used to generate user tokens.

The user uploads his own email address and a captcha (the captcha should be written as 0 when uploading for the first time). If the email address is not authorized, CESS-Gateway will send a captcha to the mailbox. After the user gets the captcha in the mailbox, upload it again. CESS-Gateway generates a token for it and sends it to the mailbox，the validity of the token is 30 days.

If the mailbox uploaded by the user has been authorized, CESS-Gateway will check whether its token has expired. If it expires, CESS-Gateway will generate a new token and send it to its mailbox.

- Request Header

| field        | value            |
| ------------ | ---------------- |
| Content-Type | application/json |

- Request Body

| field   | value                         |
| ------- | ----------------------------- |
| mailbox | your mailbox address          |
| captcha | captcha received by mailbox   |

- Responses

Response Schema: `application/json`

| status code               | structure                                 | description                                                  |
| ------------------------- | ----------------------------------------- | ------------------------------------------------------------ |
| 200 OK                    | code:200<br />msg:string<br />data:string | `msg` Enum:["success","captcha has expired and a new captcha has been sent to your mailbox","A new token has been sent to your mailbox"]<br />`data` Enum:["","token"] |
| 400 Bad Request           | code:400<br />msg:string                  | `msg` Enum:["HTTP error"，"Email format error"，"captcha error"，"Please check your email address and whether to enable SMTP service"] |
| 500 Internal Server Error | code:500<br />msg:string                  | `msg` Enum: ["Server internal data error"，"Server unexpected error"] |

- Request example

```
curl URL/auth -X POST -d '{"mailbox": "", "captcha": 0}' -H "Content-Type: application/json"
```


## Upload a file

| **PUT** /{filename} |
| ------------------- |

The put file interface is used to upload files to the cess system. You need to submit the file as form data and use provide the specific field.
If the upload is successful, you will get the fid of the file.

- Request Header

| field         | value   |
| ------------- | ------- |
| Authorization | token   |

- Request Body

| field | value        |
| ----- | ------------ |
| file  | file[binary] |

- Responses

Response Schema: `application/json`

| status code               | structure                | description                                                  |
| ------------------------- | ------------------------ | ------------------------------------------------------------ |
| 200 OK                    | code:200<br />msg:string | `msg` Default:"success                                       |
| 400 Bad Request           | code:400<br />msg:string | `msg` Default:"HTTP error"                                   |
| 401 Unauthorized          | code:401<br />msg:string | `msg` Enum:["Unauthorized"，"token expired"]                 |
| 403 Forbidden             | code:403<br />msg:string | `msg` Enum: ["duplicate filename"，"not enough space"，"The file is in hot backup, please try again later."] |
| 500 Internal Server Error | code:500<br />msg:string | `msg` Enum: ["Server internal data error"，"Server internal chain data error"，"Server unexpected error"] |

- Request example

```
# curl -X PUT URL/test.log -F 'file=@test.log' -H "Authorization: token"
{"code":200,"msg":"success","data":"fid"}
```

## Download a file

| **GET** /{fid} |
| -------------- |

The get file interface downloads the file in the CESS storage system according to the fid.

- Responses

The response schema for the normal return status is: `application/octet-stream`

The response schema for the exception return status is: `application/json`, The message returned by the exception is as follows:

| status code               | structure                | description                                                  |
| ------------------------- | ------------------------ | ------------------------------------------------------------ |
| 400 Bad Request           | code:400<br />msg:string | `msg` Enum:["HTTP error"，"This file has not been uploaded"] |
| 401 Unauthorized          | code:401<br />msg:string | `msg` Enum:["Unauthorized"，"token expired"]                 |
| 403 Forbidden             | code:403<br />msg:string | `msg` Enum: ["Token is not valid","Token expired"]           |
| 500 Internal Server Error | code:500<br />msg:string | `msg` Enum: ["Server internal data error"，"Server internal chain data error"，"Server unexpected error"] |

- Request example

```
curl -o {filename} -X GET -L URL/{fid}
```


## Delete a file

The delete file interface is used for delete a put file.

| **DELETE** /{fid} |
| ----------------- |

- Request Header

| field         | value   |
| ------------- | ------- |
| Authorization |  token  |

- Responses

Response Schema: `application/json`

| status code               | structure                | description                                                  |
| ------------------------- | ------------------------ | ------------------------------------------------------------ |
| 200 OK                    | code:200<br />msg:string | `msg` Default: "success"                                     |
| 400 Bad Request           | code:400<br />msg:string | `msg` Default: "HTTP error"                                  |
| 401 Unauthorized          | code:401<br />msg:string | `msg` Enum:["Unauthorized"，"token expired"]                 |
| 404 Not Found             | code:404<br />msg:string | `msg`Default: "This file has not been uploaded"              |
| 500 Internal Server Error | code:500<br />msg:string | `msg` Enum: ["Server internal data error"，"Server unexpected error"] |

- Request example

```
curl -X DELETE URL/{fid} -H "Authorization: token"
```

## List previous operation

| **GET** /files |
| -------------- |

List the previously put files, and display the 30 files closest to the current time by default. It also supports searching by page.

- Request Header

| field         | value   |
| ------------- | ------- |
| Authorization |  token  |

- Query Parameters

| field | description                                                  |
| ----- | ------------------------------------------------------------ |
| size  | type:<int32><br />default:30<br />Specifies the maximum number of uploads to return，up to 1000. |
| page  | type:<int32><br />default:0<br />Specifies the number of puts on which page to return. |

- Responses

Response Schema: `application/json`

| status code               | structure                                   | description                                                  |
| ------------------------- | ------------------------------------------- | ------------------------------------------------------------ |
| 200 OK                    | code:200<br />msg:string<br />data:[]string | `msg` Default:"success"<br />`data`:[file name list]         |
| 400 Bad Request           | code:400<br />msg:string                    | `msg` Default: "HTTP error"                                  |
| 401 Unauthorized          | code:401<br />msg:string                    | `msg` Enum:["Unauthorized"，"token expired"]                 |
| 500 Internal Server Error | code:500<br />msg:string                    | `msg` Enum: ["Server internal data error"，"Server unexpected error"] |

- Request example

```
curl -X GET URL/files -H "Authorization: token"
```

## View file status

| **GET** /state/{fid} |
| -------------------- |

View file status (size, status, name).

- Responses

Response Schema: `application/json`

| status code               | structure                                   | description                                                  |
| ------------------------- | ------------------------------------------- | ------------------------------------------------------------ |
| 200 OK                    | code:200<br />msg:string<br />data:{"Size","State","Names"} | `msg` `data`:{"Size","State","Names"}        |
| 400 Bad Request           | code:400<br />msg:string                    | `msg` Default: "HTTP error"                                  |
| 404 Not Found             | code:404<br />msg:string                    | `msg`Default: "Empty"                                        |
| 500 Internal Server Error | code:500<br />msg:string                    | `msg` Enum: ["Server internal data error"，"Server unexpected error"] |

- Request example

```
curl -X GET URL/state/{fid}
```

## License

Licensed under [Apache 2.0](https://github.com/CESSProject/cess-gateway/blob/main/LICENSE)
