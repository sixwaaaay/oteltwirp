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
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type config struct {
	// Propagator is used for extracting and injecting
	// trace context from requests.
	Propagator propagation.TextMapPropagator

	// TracerProvider is the TracerProvider to use for creating a Tracer.
	TracerProvider trace.TracerProvider
}

// Option is a function that applies an option to the config.
type Option interface {
	apply(*config)
}

type optionFunc func(*config)

func (o optionFunc) apply(c *config) {
	o(c)
}

func newConfig(opts []Option) *config {
	c := &config{
		Propagator:     otel.GetTextMapPropagator(),
		TracerProvider: otel.GetTracerProvider(),
	}
	for _, opt := range opts {
		opt.apply(c)
	}
	return c
}

func WithPropagators(p propagation.TextMapPropagator) Option {
	return optionFunc(func(c *config) {
		if p != nil {
			c.Propagator = p
		}
	})
}

func WithTracerProvider(tp trace.TracerProvider) Option {
	return optionFunc(func(c *config) {
		if tp != nil {
			c.TracerProvider = tp
		}
	})
}
