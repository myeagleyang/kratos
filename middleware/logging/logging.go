package logging

import (
	"context"
	"fmt"
	"gitlab.wwgame.com/wwgame/kratos/v2/errors"
	"gitlab.wwgame.com/wwgame/kratos/v2/log"
	"gitlab.wwgame.com/wwgame/kratos/v2/middleware"
	"gitlab.wwgame.com/wwgame/kratos/v2/transport"
	"time"
)

// Redacter defines how to log an object
type Redacter interface {
	Redact() string
}

// Server is an server logging middleware.
func Server(logger log.Logger) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			var (
				code      int32
				reason    string
				kind      string
				operation string
			)
			startTime := time.Now()
			if info, ok := transport.FromServerContext(ctx); ok {
				kind = info.Kind().String()
				operation = info.Operation()
			}
			reply, err = handler(ctx, req)
			if se := errors.FromError(err); se != nil {
				code = se.Code
				reason = se.Reason
			}
			level, stack := extractError(err)
			_ = log.WithContext(ctx, logger).Log(level,
				"kind", "server",
				"component", kind,
				"operation", operation,
				"args", extractArgs(req),
				"code", code,
				"reason", reason,
				"stack", stack,
				"latency", time.Since(startTime).Seconds(),
			)
			return
		}
	}
}

// Client is a client logging middleware.
func Client(logger log.Logger) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			var (
				code      int32
				reason    string
				kind      string
				operation string
			)
			startTime := time.Now()
			if info, ok := transport.FromClientContext(ctx); ok {
				kind = info.Kind().String()
				operation = info.Operation()
			}
			reply, err = handler(ctx, req)
			if se := errors.FromError(err); se != nil {
				code = se.Code
				reason = se.Reason
			}
			level, stack := extractError(err)
			_ = log.WithContext(ctx, logger).Log(level,
				"kind", "client",
				"component", kind,
				"operation", operation,
				"args", extractArgs(req),
				"code", code,
				"reason", reason,
				"stack", stack,
				"latency", time.Since(startTime).Seconds(),
			)
			return
		}
	}
}

// extractArgs returns the string of the req
func extractArgs(req interface{}) string {
	result := ""
	if redacter, ok := req.(Redacter); ok {
		result = redacter.Redact()
	}
	if stringer, ok := req.(fmt.Stringer); ok {
		result = stringer.String()
	}
	if result == "" {
		result = fmt.Sprintf("%+v", req)
	}
	if len(result) > 1024 {
		return ""
	}
	return result
}

// extractError returns the string of the error
func extractError(err error) (log.Level, string) {
	if err != nil {
		return log.LevelError, fmt.Sprintf("%+v", err)
	}
	return log.LevelInfo, ""
}
