/*
Copyright 2026 The Karmada Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package proxy

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apiserver/pkg/endpoints/metrics"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/component-base/metrics/legacyregistry"
	"k8s.io/component-base/metrics/testutil"

	pluginruntime "github.com/karmada-io/karmada/pkg/search/proxy/framework/runtime"
	"github.com/karmada-io/karmada/pkg/util/names"
)

func TestController_Connect_Metrics(t *testing.T) {
	metrics.Register()
	t.Cleanup(metrics.Reset)

	ctl := &Controller{
		proxy:                pluginruntime.NewFramework(convertPluginSlice([]*mockPlugin{{IsSupportRequest: true}})),
		negotiatedSerializer: scheme.Codecs.WithoutConversion(),
	}

	tests := []struct {
		name      string
		proxyPath string
		query     string
		wantVerb  string
		wantScope string
	}{
		{name: "list all namespaces", proxyPath: "/api/v1/pods", wantVerb: "LIST", wantScope: "cluster"},
		{name: "list one namespace", proxyPath: "/api/v1/namespaces/default/pods", wantVerb: "LIST", wantScope: "namespace"},
		{name: "get one pod", proxyPath: "/api/v1/namespaces/default/pods/test", wantVerb: "GET", wantScope: "resource"},
		{name: "watch pods", proxyPath: "/api/v1/pods", query: "?watch=true", wantVerb: "WATCH", wantScope: "cluster"},
	}

	// The apiserver metrics registry is global, so keep cases sequential and reset between them.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metrics.Reset()

			h, err := ctl.Connect(context.TODO(), tt.proxyPath, nil)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodGet, "/prefix"+tt.proxyPath+tt.query, nil)
			recorder := httptest.NewRecorder()
			h.ServeHTTP(recorder, req)
			require.Equal(t, http.StatusOK, recorder.Code)

			want := fmt.Sprintf(`
# HELP apiserver_request_total [STABLE] Counter of apiserver requests broken out for each verb, dry run value, group, version, resource, scope, component, and HTTP response code.
# TYPE apiserver_request_total counter
apiserver_request_total{code="200",component=%q,dry_run="",group="",resource="pods",scope=%q,subresource="",verb=%q,version="v1"} 1
`, names.KarmadaSearchComponentName, tt.wantScope, tt.wantVerb)
			require.NoError(t, testutil.GatherAndCompare(legacyregistry.DefaultGatherer, strings.NewReader(want), "apiserver_request_total"))

			// Histogram buckets and sum depend on wall-clock time, so only check the labels here.
			assert.Equal(t, map[string]string{
				"component":   names.KarmadaSearchComponentName,
				"dry_run":     "",
				"group":       "",
				"version":     "v1",
				"resource":    "pods",
				"subresource": "",
				"scope":       tt.wantScope,
				"verb":        tt.wantVerb,
			}, gatherLabels(t, "apiserver_request_duration_seconds"))
		})
	}
}

// gatherLabels returns the labels of the single sample recorded under name, failing if there is not exactly one.
func gatherLabels(t *testing.T, name string) map[string]string {
	t.Helper()
	families, err := legacyregistry.DefaultGatherer.Gather()
	require.NoError(t, err)
	for _, family := range families {
		if family.GetName() != name {
			continue
		}
		require.Len(t, family.Metric, 1, name)
		labels := make(map[string]string)
		for _, label := range family.Metric[0].GetLabel() {
			labels[label.GetName()] = label.GetValue()
		}
		return labels
	}
	require.Failf(t, "metric not found", "%s", name)
	return nil
}
