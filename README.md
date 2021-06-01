# websiteGo
website with using  Go lang

#How To Run
*Create a new database with a users table*


###employee
```mysql
CREATE TABLE `employee` (
  `id` int(6) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(30) NOT NULL,
  `city` varchar(30) NOT NULL,
  PRIMARY KEY (`id`)
)
```
###users
```mysql
CREATE TABLE users(
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50),
    password VARCHAR(120)
);
```

##Go get required packages

```gitexclude
go get golang.org/x/crypto/bcrypt

go get github.com/go-sql-driver/mysql
```

##Run the following command

```go
go run main.go
```

##Load the following URL
```go
192.168.1.9:80
```





