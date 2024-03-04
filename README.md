# gin-i18n
Customized international multilingual support for gin, built on **golangorg/x/text/language** and supporting GRPC server calls

```shell
go get -u https://github.com/hinet/gin-i18n
```

Using for GRPC

```go
import (
    "google.golang.org/grpc"
    grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
)

grpc.NewServer(
    grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
        i18n.LanguageInterceptor,
    )),
)
//services

func login(ctx context.Context) {
    t := i18n.GetTranslator()
    println(t.Translate("home.welcome!", nil))
}
```