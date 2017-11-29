# hamster
Backend contents server with golang

## Setup
### 1. Install & run MongoDB
#### MacOS
Install
```
$ brew install mongodb
$ sudo mkdir /var/lib/mongodb
$ sudo touch /var/log/mongodb.log
```
Run
```
$ sudo mongod --dbpath /var/lib/mongodb --logpath /var/log/mongodb.log
```
### 2. Install Go dependencies
```
$ dep ensure
```

## Run
```
$ go run main.go
```