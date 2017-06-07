FROM golang

RUN go get github.com/pangpanglabs/echosample \
    go test github.com/pangpanglabs/echosample/...

WORKDIR $GOPATH/src/github.com/pangpanglabs/echosample

CMD go run main.go
