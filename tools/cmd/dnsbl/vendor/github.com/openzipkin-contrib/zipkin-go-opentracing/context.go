// Copyright 2022 The OpenZipkin Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package zipkintracer

import (
	"github.com/openzipkin/zipkin-go/model"
)

// SpanContext holds the basic Span metadata.
type SpanContext model.SpanContext

// ForeachBaggageItem belongs to the opentracing.SpanContext interface
func (c SpanContext) ForeachBaggageItem(handler func(k, v string) bool) {}
