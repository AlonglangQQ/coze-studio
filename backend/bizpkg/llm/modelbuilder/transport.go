/*
 * Copyright 2025 coze-dev Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package modelbuilder

import (
	"net/http"
)

type headerTransport struct {
	transport http.RoundTripper
	headers   map[string]string
}

func (t *headerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	for k, v := range t.headers {
		req.Header.Set(k, v)
	}
	if t.transport == nil {
		return http.DefaultTransport.RoundTrip(req)
	}
	return t.transport.RoundTrip(req)
}

func newHeaderClient(headers map[string]string) *http.Client {
	return &http.Client{
		Transport: &headerTransport{
			headers:   headers,
			transport: http.DefaultTransport,
		},
	}
}
