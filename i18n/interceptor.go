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

	// Extracting "Accept-Language" information from HTTP headers from metadata
	acceptLanguages := md.Get("accept-language")
	if len(acceptLanguages) == 0 {
		return nil, fmt.Errorf("accept-language header not found")
	}
	language := acceptLanguages[0]
	//Initialize translator
	translator := NewTranslator()
	tag := GetPreferredLanguage(language)
	translator.SetLanguage(tag)
	return handler(ctx, req)
}
