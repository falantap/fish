# Auth & Fetching

There are 2 programs in this repository

-Auth which is JWT Program using non persistent registry for user implemented in GoLang

-Fetching which is a program that calling HTTP request and manipulating JSON response as desired with some JWT Authentication, implemented in NodeJs Express Server

(Tested on Windows)
## Installation

Clone with HTTPS

```bash
git clone https://github.com/falantap/fish.git
```
Fetching, install [nodejs](https://nodejs.org/en/download/)

Auth, install [GO](https://golang.org/doc/install)

## Running The Program
Configure Server Port for Fetching Program

```bash
nodejs\fetching\startup_port.json
```
Configure Server Port for Auth Program

```bash
go\src\go-auth\server_url.txt
```

Fetching Server Startup Command
```python
nodejs\fetching\node index.js
```
Auth Server Startup Command
```python
go\src\go-auth\go run main.go
```
## Usage
Fetching Endpoint, Documentation click [Here](https://documenter.getpostman.com/view/9400071/SzzheyHA)
```python
/jwt
/fetch
/privateClaims
/aggregate
```
Auth Endpoing, Documentation click [Here](https://documenter.getpostman.com/view/9400071/SzzheyHB)
```python
/auth/register
/auth/getData
/auth/verify
```
## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.
