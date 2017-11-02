// Copyright (c) 2017 Uber Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mocks

import mock "github.com/stretchr/testify/mock"
import model "github.com/uber/jaeger/model"
import spanstore "github.com/uber/jaeger/storage/spanstore"

// Writer is an autogenerated mock type for the Writer type
type Writer struct {
	mock.Mock
}

// WriteSpan provides a mock function with given fields: span
func (_m *Writer) WriteSpan(span *model.Span) error {
	ret := _m.Called(span)

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.Span) error); ok {
		r0 = rf(span)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

var _ spanstore.Writer = (*Writer)(nil)