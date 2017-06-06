# Offer

## Getting Started

Git clone

```
$ cd $GOPATH/src
$ git clone http://git.elandsystems.com.cn:3000/pangpangjan/offer.git
```

Run

```
$ cd offer
$ go test offer/controllers
$ go run main.go
```

Visit http://127.0.0.1:8080/

## Prerequisite

### Install packages

```
$ go get github.com/labstack/echo
$ go get github.com/go-xorm/xorm
$ go get github.com/spf13/viper
$ go get github.com/pangpanglabs/goutils...
$ go get github.com/mattn/go-sqlite3
$ go get github.com/go-sql-driver/mysql
$ go get github.com/asaskevich/govalidator
$ go get github.com/sirupsen/logrus
```

## References

- web framework: [echo framework](https://echo.labstack.com/)
- orm tool: [xorm](http://xorm.io/)
- logger : [logrus](github.com/sirupsen/logrus)
- configuration tool: [viper](https://github.com/spf13/viper)
- validator: [govalidator](github.com/asaskevich/govalidator)
- utils: https://github.com/pangpanglabs/goutils