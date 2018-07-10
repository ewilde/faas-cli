// Copyright (c) Alex Ellis 2017. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package commands

import (
	"regexp"
	"testing"

	"fmt"
	"net/http"

	"github.com/openfaas/faas-cli/test"
	"github.com/openfaas/faas-cli/version"
	)

func Test_addVersionDev(t *testing.T) {
	version.GitCommit = "sha-test"

	stdOut := test.CaptureStdout(func() {
		faasCmd.SetArgs([]string{"version"})
		faasCmd.Execute()
	})

	expected := "CLI commit: sha-test"
	if found, err := regexp.MatchString(fmt.Sprintf(`(?m:%s)`, expected), stdOut); err != nil || !found {
		t.Fatalf("Commit is not as expected - want: %s, got: %s", expected, stdOut)
	}

	expected = "CLI version: dev"
	if found, err := regexp.MatchString(fmt.Sprintf(`(?m:%s)`, expected), stdOut); err != nil || !found {
		t.Fatalf("Version is not as expected - want: %s, got: %s", expected, stdOut)
	}
}

func Test_addVersion(t *testing.T) {
	version.GitCommit = "sha-test"
	version.Version = "version.tag"

	stdOut := test.CaptureStdout(func() {
		faasCmd.SetArgs([]string{"version"})
		faasCmd.Execute()
	})

	expected := "CLI commit: sha-test"
	if found, err := regexp.MatchString(fmt.Sprintf(`(?m:%s)`, expected), stdOut); err != nil || !found {
		t.Fatalf("Commit is not as expected:\n%s", stdOut)
	}

	expected = "CLI version: version.tag"
	if found, err := regexp.MatchString(fmt.Sprintf(`(?m:%s)`, expected), stdOut); err != nil || !found {
		t.Fatalf("Version is not as expected - want: %s, got: %s", expected, stdOut)
	}
}

func Test_addVersion_short_version(t *testing.T) {
	version.Version = "version.tag"

	stdOut := test.CaptureStdout(func() {
		faasCmd.SetArgs([]string{"version", "--short-version"})
		faasCmd.Execute()
	})

	if found, err := regexp.MatchString("^version\\.tag", stdOut); err != nil || !found {
		t.Fatalf("Version is not as expected - want: %s, got: %s", version.Version, stdOut)
	}
}

func Test_gateway_and_provider_information(t *testing.T) {
	var testCases =
		[]struct{
		 responseBody string
		 params []struct {
			name       string
			value      string
		 }
		}{
			{
				responseBody: gateway_response_0_8_4_onwards,
				params: []struct {
					name string
					value string
				}{
					{ "gateway version",         "version: gateway-0.4.3" },
					{ "gateway sha"            ,"sha: 999a6669148c30adeb64400609953cf59db2fb64"},
					{ "gateway commit"         ,"commit: Bump faas-swarm to latest"},
					{ "provider name"          ,"name:          faas-swarm"},
					{ "provider orchestration" ,"orchestration: swarm"},
					{ "provider version"       ,"version:       provider-0.3.3"},
					{ "provider sha"           ,"sha:           c890cba302d059de8edbef3f3de7fe15444b1ecf"},

				},
			},
			{
				responseBody:gateway_response_prior_to_0_8_4,
				params: []struct {
					name string
					value string
				}{
					{  "provider name"          ,"name:          faas-swarm"},
					{  "provider orchestration" ,"orchestration: swarm"},
					{  "provider version"       ,"version:       provider-0.3.3"},
					{  "provider sha"           ,"sha:           c890cba302d059de8edbef3f3de7fe15444b1ecf"},

				},
			},
		}


	for _, testCase := range testCases {


		for _, param := range testCase.params {
			t.Run(param.name, func(t *testing.T) {
				resetForTest()
				s := test.MockHttpServer(t, []test.Request{
					{
						Method:             http.MethodGet,
						Uri:                "/system/info",
						ResponseStatusCode: http.StatusOK,
						ResponseBody:       testCase.responseBody,
					},
				})
				defer s.Close()

				stdOut := test.CaptureStdout(func() {
					faasCmd.SetArgs([]string{
						"version",
						"--gateway=" + s.URL,
					})
					faasCmd.Execute()
				})
				if found, err := regexp.MatchString(fmt.Sprintf(`(?m:%s)`, param.value), stdOut); err != nil || !found {
					t.Fatalf("%s is not as expected - want: `%s` got: `%s`", param.name, param.value, stdOut)
				}
			})
		}
	}
}

func Test_gateway_uri(t *testing.T) {
	resetForTest()
	s := test.MockHttpServer(t, []test.Request{
		{
			Method:             http.MethodGet,
			Uri:                "/system/info",
			ResponseStatusCode: http.StatusOK,
			ResponseBody:       gateway_response_0_8_4_onwards,
		},
	})
	defer s.Close()

	stdOut := test.CaptureStdout(func() {
		faasCmd.SetArgs([]string{
			"version",
			"--gateway=" + s.URL,
		})
		faasCmd.Execute()
	})

	if found, err := regexp.MatchString(fmt.Sprintf(`(?m:uri: %s)`, s.URL), stdOut); err != nil || !found {
		t.Fatalf("Output is not as expected:\n%s", stdOut)
	}
}

func Test_gateway_version(t *testing.T) {
	resetForTest()
	s := test.MockHttpServer(t, []test.Request{
		{
			Method:             http.MethodGet,
			Uri:                "/system/info",
			ResponseStatusCode: http.StatusOK,
			ResponseBody:       gateway_response_0_8_4_onwards,
		},
	})
	defer s.Close()

	stdOut := test.CaptureStdout(func() {
		faasCmd.SetArgs([]string{
			"version",
			"--gateway=" + s.URL,
		})
		faasCmd.Execute()
	})

	if found, err := regexp.MatchString(`(?m:version: gateway-0.4.3)`, stdOut); err != nil || !found {
		t.Fatalf("Output is not as expected:\n%s", stdOut)
	}
}

func Test_gateway_sha(t *testing.T) {
	resetForTest()
	s := test.MockHttpServer(t, []test.Request{
		{
			Method:             http.MethodGet,
			Uri:                "/system/info",
			ResponseStatusCode: http.StatusOK,
			ResponseBody:       gateway_response_0_8_4_onwards,
		},
	})
	defer s.Close()

	stdOut := test.CaptureStdout(func() {
		faasCmd.SetArgs([]string{
			"version",
			"--gateway=" + s.URL,
		})
		faasCmd.Execute()
	})

	if found, err := regexp.MatchString(`(?m:sha: 999a6669148c30adeb64400609953cf59db2fb64)`, stdOut); err != nil || !found {
		t.Fatalf("Output is not as expected:\n%s", stdOut)
	}
}

func Test_gateway_commit(t *testing.T) {
	resetForTest()
	s := test.MockHttpServer(t, []test.Request{
		{
			Method:             http.MethodGet,
			Uri:                "/system/info",
			ResponseStatusCode: http.StatusOK,
			ResponseBody:       gateway_response_0_8_4_onwards,
		},
	})
	defer s.Close()

	stdOut := test.CaptureStdout(func() {
		faasCmd.SetArgs([]string{
			"version",
			"--gateway=" + s.URL,
		})
		faasCmd.Execute()
	})

	if found, err := regexp.MatchString(`(?m:commit: Bump faas-swarm to latest)`, stdOut); err != nil || !found {
		t.Fatalf("Output is not as expected:\n%s", stdOut)
	}
}

func Test_provider_name(t *testing.T) {
	resetForTest()
	s := test.MockHttpServer(t, []test.Request{
		{
			Method:             http.MethodGet,
			Uri:                "/system/info",
			ResponseStatusCode: http.StatusOK,
			ResponseBody:       gateway_response_0_8_4_onwards,
		},
	})
	defer s.Close()

	stdOut := test.CaptureStdout(func() {
		faasCmd.SetArgs([]string{
			"version",
			"--gateway=" + s.URL,
		})
		faasCmd.Execute()
	})

	if found, err := regexp.MatchString(`(?m:name: faas-swarm)`, stdOut); err != nil || !found {
		t.Fatalf("Output is not as expected:\n%s", stdOut)
	}
}

func Test_provider_orchestration(t *testing.T) {
	resetForTest()
	s := test.MockHttpServer(t, []test.Request{
		{
			Method:             http.MethodGet,
			Uri:                "/system/info",
			ResponseStatusCode: http.StatusOK,
			ResponseBody:       gateway_response_0_8_4_onwards,
		},
	})
	defer s.Close()

	stdOut := test.CaptureStdout(func() {
		faasCmd.SetArgs([]string{
			"version",
			"--gateway=" + s.URL,
		})
		faasCmd.Execute()
	})

	if found, err := regexp.MatchString(`(?m:orchestration: swarm)`, stdOut); err != nil || !found {
		t.Fatalf("Output is not as expected:\n%s", stdOut)
	}
}

func Test_provider_version(t *testing.T) {
	resetForTest()
	s := test.MockHttpServer(t, []test.Request{
		{
			Method:             http.MethodGet,
			Uri:                "/system/info",
			ResponseStatusCode: http.StatusOK,
			ResponseBody:       gateway_response_0_8_4_onwards,
		},
	})
	defer s.Close()

	stdOut := test.CaptureStdout(func() {
		faasCmd.SetArgs([]string{
			"version",
			"--gateway=" + s.URL,
		})
		faasCmd.Execute()
	})

	if found, err := regexp.MatchString(`(?m:version: provider-0.3.3)`, stdOut); err != nil || !found {
		t.Fatalf("Output is not as expected:\n%s", stdOut)
	}
}

func Test_provider_sha(t *testing.T) {
	resetForTest()
	s := test.MockHttpServer(t, []test.Request{
		{
			Method:             http.MethodGet,
			Uri:                "/system/info",
			ResponseStatusCode: http.StatusOK,
			ResponseBody:       gateway_response_0_8_4_onwards,
		},
	})
	defer s.Close()

	stdOut := test.CaptureStdout(func() {
		faasCmd.SetArgs([]string{
			"version",
			"--gateway=" + s.URL,
		})
		faasCmd.Execute()
	})

	if found, err := regexp.MatchString(`(?m:sha: c890cba302d059de8edbef3f3de7fe15444b1ecf)`, stdOut); err != nil || !found {
		t.Fatalf("Output is not as expected:\n%s", stdOut)
	}
}

func Test_gateway_uri_prior_to_0_8_4(t *testing.T) {
	resetForTest()
	s := test.MockHttpServer(t, []test.Request{
		{
			Method:             http.MethodGet,
			Uri:                "/system/info",
			ResponseStatusCode: http.StatusOK,
			ResponseBody:       gateway_response_prior_to_0_8_4,
		},
	})
	defer s.Close()

	stdOut := test.CaptureStdout(func() {
		faasCmd.SetArgs([]string{
			"version",
			"--gateway=" + s.URL,
		})
		faasCmd.Execute()
	})

	if found, err := regexp.MatchString(fmt.Sprintf(`(?m:uri: %s)`, s.URL), stdOut); err != nil || !found {
		t.Fatalf("Output is not as expected:\n%s", stdOut)
	}
}

func Test_gateway_details_prior_to_0_8_4_should_not_be_displayed(t *testing.T) {
	resetForTest()
	s := test.MockHttpServer(t, []test.Request{
		{
			Method:             http.MethodGet,
			Uri:                "/system/info",
			ResponseStatusCode: http.StatusOK,
			ResponseBody:       gateway_response_prior_to_0_8_4,
		},
	})
	defer s.Close()

	stdOut := test.CaptureStdout(func() {
		faasCmd.SetArgs([]string{
			"version",
			"--gateway=" + s.URL,
		})
		faasCmd.Execute()
	})

	if found, err := regexp.MatchString(`(?m:\tversion: $)`, stdOut); err != nil || found {
		t.Fatalf("Output is not as expected for version:\n%s", stdOut)
	}

	if found, err := regexp.MatchString(`(?m:\tsha: $)`, stdOut); err != nil || found {
		t.Fatalf("Output is not as expected for sha:\n%s", stdOut)
	}

	if found, err := regexp.MatchString(`(?m:\tcommit: $)`, stdOut); err != nil || found {
		t.Fatalf("Output is not as expected for commit:\n%s", stdOut)
	}
}

func Test_provider_name_prior_to_0_8_4(t *testing.T) {
	resetForTest()
	s := test.MockHttpServer(t, []test.Request{
		{
			Method:             http.MethodGet,
			Uri:                "/system/info",
			ResponseStatusCode: http.StatusOK,
			ResponseBody:       gateway_response_prior_to_0_8_4,
		},
	})
	defer s.Close()

	stdOut := test.CaptureStdout(func() {
		faasCmd.SetArgs([]string{
			"version",
			"--gateway=" + s.URL,
		})
		faasCmd.Execute()
	})

	if found, err := regexp.MatchString(`(?m:name: faas-swarm)`, stdOut); err != nil || !found {
		t.Fatalf("Output is not as expected:\n%s", stdOut)
	}
}

func Test_provider_sha_prior_to_0_8_4(t *testing.T) {
	resetForTest()
	s := test.MockHttpServer(t, []test.Request{
		{
			Method:             http.MethodGet,
			Uri:                "/system/info",
			ResponseStatusCode: http.StatusOK,
			ResponseBody:       gateway_response_prior_to_0_8_4,
		},
	})
	defer s.Close()

	stdOut := test.CaptureStdout(func() {
		faasCmd.SetArgs([]string{
			"version",
			"--gateway=" + s.URL,
		})
		faasCmd.Execute()
	})

	if found, err := regexp.MatchString(`(?m:sha: c890cba302d059de8edbef3f3de7fe15444b1ecf)`, stdOut); err != nil || !found {
		t.Fatalf("Output is not as expected:\n%s", stdOut)
	}
}

func Test_provider_version_prior_to_0_8_4(t *testing.T) {
	resetForTest()
	s := test.MockHttpServer(t, []test.Request{
		{
			Method:             http.MethodGet,
			Uri:                "/system/info",
			ResponseStatusCode: http.StatusOK,
			ResponseBody:       gateway_response_prior_to_0_8_4,
		},
	})
	defer s.Close()

	stdOut := test.CaptureStdout(func() {
		faasCmd.SetArgs([]string{
			"version",
			"--gateway=" + s.URL,
		})
		faasCmd.Execute()
	})

	if found, err := regexp.MatchString(`(?m:version: provider-0.3.3)`, stdOut); err != nil || !found {
		t.Fatalf("Output is not as expected:\n%s", stdOut)
	}
}

func Test_provider_orchestration_prior_to_0_8_4(t *testing.T) {
	resetForTest()
	s := test.MockHttpServer(t, []test.Request{
		{
			Method:             http.MethodGet,
			Uri:                "/system/info",
			ResponseStatusCode: http.StatusOK,
			ResponseBody:       gateway_response_prior_to_0_8_4,
		},
	})
	defer s.Close()

	stdOut := test.CaptureStdout(func() {
		faasCmd.SetArgs([]string{
			"version",
			"--gateway=" + s.URL,
		})
		faasCmd.Execute()
	})

	if found, err := regexp.MatchString(`(?m:orchestration: swarm)`, stdOut); err != nil || !found {
		t.Fatalf("Output is not as expected:\n%s", stdOut)
	}
}

const gateway_response_0_8_4_onwards = `{
  "provider": {
    "provider": "faas-swarm",
    "orchestration": "swarm",
    "version": {
      "sha": "c890cba302d059de8edbef3f3de7fe15444b1ecf",
      "release": "provider-0.3.3"
    }
  },
  "version": {
    "sha": "999a6669148c30adeb64400609953cf59db2fb64",
    "release": "gateway-0.4.3",
    "commit_message": "Bump faas-swarm to latest"
  } 
}`

const gateway_response_prior_to_0_8_4 = `{
  "provider": "faas-swarm",
  "version": {
    "sha": "c890cba302d059de8edbef3f3de7fe15444b1ecf",
    "release": "provider-0.3.3"
  },
  "orchestration": "swarm"
}`
