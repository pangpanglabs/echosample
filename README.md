# echosample

## Getting Started

Get source
```
$ go get github.com/pangpanglabs/echosample
```

Test
```
$ go test github.com/pangpanglabs/echosample/...
```

Run
```
$ cd $GOPATH/src/github.com/pangpanglabs/echosample
$ go run main.go
```

Visit http://127.0.0.1:8080/

## Tips

### Live reload utility

[github.com/codegangsta/gin](https://github.com/codegangsta/gin)

```
$ gin -a 8080  -i --all r
```

Visit http://127.0.0.1:3000/

## Prerequisite

Install packages
```
$ go get github.com/pangpanglabs/goutils...
```

## References

- web framework: [echo framework](https://echo.labstack.com/)
- orm tool: [xorm](http://xorm.io/)
- logger : [logrus](https://github.com/sirupsen/logrus)
- configuration tool: [viper](https://github.com/spf13/viper)
- validator: [govalidator](github.com/asaskevich/govalidator)
- utils: https://github.com/pangpanglabs/goutils