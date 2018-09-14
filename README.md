# echosample

[![Build Status](https://travis-ci.org/pangpanglabs/echosample.svg?branch=master)](https://travis-ci.org/pangpanglabs/echosample)
[![codecov](https://codecov.io/gh/pangpanglabs/echosample/branch/master/graph/badge.svg)](https://codecov.io/gh/pangpanglabs/echosample)


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

Install
```
$ go get github.com/codegangsta/gin
```

Run
```
$ gin -a 8080  -i --all r
```

Visit http://127.0.0.1:3000/


## References

- web framework: [echo framework](https://echo.labstack.com/)
- orm tool: [xorm](http://xorm.io/)
- logger : [logrus](https://github.com/sirupsen/logrus)
- configuration tool: [viper](https://github.com/spf13/viper)
- validator: [govalidator](github.com/asaskevich/govalidator)
- utils: https://github.com/pangpanglabs/goutils
