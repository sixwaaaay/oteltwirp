package oteltwirp

import (
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
)

const (
	// instrumentationName is the name of this instrumentation package.
	instrumentationName = "github.com/sixwaaaay/twirp-tracing"
	Version             = "v0.1.0"
	// RPCStatusCodeKey is convention for numeric status code of a twirp request.
	RPCStatusCodeKey = attribute.Key("rpc.twirp.status_code")
)

// Semantic conventions for common RPC attributes.
var (
	// RPCSystemTwirp Semantic convention for twirp as the remoting system.
	RPCSystemTwirp = semconv.RPCSystemKey.String("twirp")
)
