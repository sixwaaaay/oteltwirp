/*
 * Copyright (c) 2023 sixwaaaay.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package oteltwirp

import (
	"context"
	"github.com/twitchtv/twirp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	"go.opentelemetry.io/otel/trace"
)

func ServerInterceptor(opts ...Option) twirp.Interceptor {
	conf := newConfig(opts)
	tracer := conf.TracerProvider.Tracer(
		instrumentationName,
		trace.WithInstrumentationVersion(Version),
	)
	return func(method twirp.Method) twirp.Method {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			name, attrs := baseAttrs(ctx)

			ctx, span := tracer.Start(
				trace.ContextWithRemoteSpanContext(ctx, trace.SpanContextFromContext(ctx)),
				name,
				trace.WithSpanKind(trace.SpanKindServer), // Server Span
				trace.WithAttributes(attrs...),
			)
			defer span.End()

			resp, err := method(ctx, request)
			if err != nil {
				if status, b := twirp.StatusCode(ctx); b {
					statusCode, msg := serverStatus(twirp.ErrorCode(status))
					span.SetStatus(statusCode, msg)
					span.SetAttributes(RPCStatusCodeKey.String(status))
				}
			} else {
				span.SetStatus(codes.Ok, "OK")
			}
			return resp, err
		}
	}
}

func baseAttrs(ctx context.Context) (string, []attribute.KeyValue) {
	var name string //$package.$service/$method
	if packageName, b := twirp.PackageName(ctx); b && packageName != "" {
		name = packageName + "."
	}
	// get twirp info
	attrs := []attribute.KeyValue{RPCSystemTwirp}
	if serviceName, b := twirp.ServiceName(ctx); b && serviceName != "" {
		attrs = append(attrs, semconv.RPCService(serviceName))
		name += serviceName + "/"
	}
	if methodName, b := twirp.MethodName(ctx); b && methodName != "" {
		attrs = append(attrs, semconv.RPCMethod(methodName))
		name += methodName
	}
	return name, attrs
}

func serverStatus(status twirp.ErrorCode) (codes.Code, string) {
	if twirp.IsValidErrorCode(status) {
		return codes.Error, string(status)
	} else {
		return codes.Unset, ""
	}
}
