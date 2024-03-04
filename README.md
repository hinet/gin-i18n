# gin-i18n
Customized international multilingual support for gin, built on **golangorg/x/text/language** and supporting GRPC server calls

```shell
go get -u https://github.com/hinet/gin-i18n
```

Using for GRPC

```go
import (
    "google.golang.org/grpc"
)

grpc.NewServer(
    grpc.UnaryInterceptor(i18n.LanguageInterceptor),
)
//services

func login(ctx context.Context) {
    t := i18n.GetTranslator()
    println(t.Translate("welcome {name}!!", map[string]string{"name":"liming"}))
}
```
