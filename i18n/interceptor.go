package i18n

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func LanguageInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("failed to get metadata from context")
	}

	// 从元数据中提取HTTP标头中的"accept-language"信息
	acceptLanguages := md.Get("accept-language")
	if len(acceptLanguages) == 0 {
		return nil, fmt.Errorf("accept-language header not found")
	}
	language := acceptLanguages[0]
	//初始化语言转换器
	translator := NewTranslator()
	tag := GetPreferredLanguage(language)
	translator.SetLanguage(tag)
	// 调用下一个处理程序
	return handler(ctx, req)
}
