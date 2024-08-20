// Copyright (c) 2022 Alibaba Group Holding Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tests

import (
	"testing"

	"github.com/alibaba/higress/test/e2e/conformance/utils/http"
	"github.com/alibaba/higress/test/e2e/conformance/utils/suite"
)

func init() {
	Register(WasmPluginsAiProxy)
}

var WasmPluginsAiProxy = suite.ConformanceTest{
	ShortName:   "WasmPluginsAIProxy",
	Description: "The Ingress in the higress-conformance-infra namespace test the ai-proxy wasmplugins.",
	Manifests:   []string{"tests/go-wasm-ai-proxy.yaml"},
	Features:    []suite.SupportedFeature{suite.WASMGoConformanceFeature},
	Test: func(t *testing.T, suite *suite.ConformanceTestSuite) {
		testcases := []http.Assertion{
			{
				Meta: http.AssertionMeta{
					TestCaseName:    "openai: basic case 1",
					TargetBackend:   "infra-backend-v1",
					TargetNamespace: "higress-conformance-infra",
				},
				Request: http.AssertionRequest{
					ActualRequest: http.Request{
						Host:             "foo.com",
						Path:             "/api/openai/v1/chat/completions",
						Headers:          map[string]string{"Accept": "application/json, text/event-stream"},
						UnfollowRedirect: true,
						Body: []byte(`
						{
							"model": "qwen-long",
							"input":{
								"messages":[      
									{
										"role": "system",
										"content": "You are a helpful assistant."
									},
									{
										"role": "user",
										"content": "你是谁？"
									}
								]
							},
							"parameters": {
								"temperature":0.7
							}
						}
						`),
						ContentType: http.ContentTypeApplicationJson,
					},

					ExpectedRequest: &http.ExpectedRequest{
						Request: http.Request{
							Host:        "foo.com",
							Path:        "/api/openai/v1/chat/completions",
							Method:      "GET",
							ContentType: http.ContentTypeApplicationJson,
							Body: []byte(`
						{
							"model": "qwen-long",
							"input":{
								"messages":[      
									{
										"role": "system",
										"content": "You are a helpful assistant."
									},
									{
										"role": "user",
										"content": "你是谁？"
									}
								]
							},
							"parameters": {
								"temperature":0.66
							}
						}
						`),
						},
					},
				},
				Response: http.AssertionResponse{
					ExpectedResponse: http.Response{
						StatusCode: 200,
					},
				},
			},
		}
		t.Run("WasmPlugins ai-proxy", func(t *testing.T) {
			for _, testcase := range testcases {
				http.MakeRequestAndExpectEventuallyConsistentResponse(t, suite.RoundTripper, suite.TimeoutConfig, suite.GatewayAddress, testcase)
			}
		})
	},
}
