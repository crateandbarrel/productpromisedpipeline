go fmt *.go
go fmt ./clients/*.go
go fmt ./oauth/*.go
go fmt ./i18n/*.go
go fmt ./oauth2/*.go
$GOPATH/bin/golint *.go
$GOPATH/bin/golint ./clients/*.go
$GOPATH/bin/golint ./oauth/*.go
$GOPATH/bin/golint ./i18n/*.go
$GOPATH/bin/golint ./oauth2/*.go
$GOPATH/bin/golint ./events/*.go
go vet github.com/cfrye2000/productPromisedEventMS
go install
cp ./productPromisedEventMS.cfg $GOPATH/bin
go test -cover ./oauth/
go test -cover ./i18n/
go test -cover ./oauth2/
go test -cover ./clients/
go test -cover ./events/
