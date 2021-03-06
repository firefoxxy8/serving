/*
Copyright 2017 The Knative Authors

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

package v1alpha1

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRouteValidation(t *testing.T) {
	tests := []struct {
		name string
		r    *Route
		want *FieldError
	}{{
		name: "valid",
		r: &Route{
			Spec: RouteSpec{
				Traffic: []TrafficTarget{{
					RevisionName: "foo",
					Percent:      100,
				}},
			},
		},
		want: nil,
	}, {
		name: "valid split",
		r: &Route{
			Spec: RouteSpec{
				Traffic: []TrafficTarget{{
					Name:         "prod",
					RevisionName: "foo",
					Percent:      90,
				}, {
					Name:              "experiment",
					ConfigurationName: "bar",
					Percent:           10,
				}},
			},
		},
		want: nil,
	}, {
		name: "invalid traffic entry",
		r: &Route{
			Spec: RouteSpec{
				Traffic: []TrafficTarget{{
					Name:    "foo",
					Percent: 100,
				}},
			},
		},
		want: &FieldError{
			Message: "Expected exactly one, got neither",
			Paths: []string{
				"spec.traffic[0].revisionName",
				"spec.traffic[0].configurationName",
			},
		},
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.r.Validate()
			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("Validate (-want, +got) = %v", diff)
			}
		})
	}
}

func TestRouteSpecValidation(t *testing.T) {
	tests := []struct {
		name string
		rs   *RouteSpec
		want *FieldError
	}{{
		name: "valid",
		rs: &RouteSpec{
			Traffic: []TrafficTarget{{
				RevisionName: "foo",
				Percent:      100,
			}},
		},
		want: nil,
	}, {
		name: "valid split",
		rs: &RouteSpec{
			Traffic: []TrafficTarget{{
				Name:         "prod",
				RevisionName: "foo",
				Percent:      90,
			}, {
				Name:              "experiment",
				ConfigurationName: "bar",
				Percent:           10,
			}},
		},
		want: nil,
	}, {
		name: "empty spec",
		rs:   &RouteSpec{},
		want: errMissingField(currentField),
	}, {
		name: "invalid traffic entry",
		rs: &RouteSpec{
			Traffic: []TrafficTarget{{
				Name:    "foo",
				Percent: 100,
			}},
		},
		want: &FieldError{
			Message: "Expected exactly one, got neither",
			Paths:   []string{"traffic[0].revisionName", "traffic[0].configurationName"},
		},
	}, {
		name: "invalid name conflict",
		rs: &RouteSpec{
			Traffic: []TrafficTarget{{
				Name:         "foo",
				RevisionName: "bar",
				Percent:      50,
			}, {
				Name:         "foo",
				RevisionName: "baz",
				Percent:      50,
			}},
		},
		want: &FieldError{
			Message: `Multiple definitions for "foo"`,
			Paths:   []string{"traffic[0].name", "traffic[1].name"},
		},
	}, {
		name: "valid name collision (same revision)",
		rs: &RouteSpec{
			Traffic: []TrafficTarget{{
				Name:         "foo",
				RevisionName: "bar",
				Percent:      50,
			}, {
				Name:         "foo",
				RevisionName: "bar",
				Percent:      50,
			}},
		},
		want: nil,
	}, {
		name: "invalid total percentage",
		rs: &RouteSpec{
			Traffic: []TrafficTarget{{
				RevisionName: "bar",
				Percent:      99,
			}, {
				RevisionName: "baz",
				Percent:      99,
			}},
		},
		want: &FieldError{
			Message: "Traffic targets sum to 198, want 100",
			Paths:   []string{"traffic"},
		},
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.rs.Validate()
			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("Validate (-want, +got) = %v", diff)
			}
		})
	}
}

func TestTrafficTargetValidation(t *testing.T) {
	tests := []struct {
		name string
		tt   *TrafficTarget
		want *FieldError
	}{{
		name: "valid with name and revision",
		tt: &TrafficTarget{
			Name:         "foo",
			RevisionName: "bar",
			Percent:      12,
		},
		want: nil,
	}, {
		name: "valid with name and configuration",
		tt: &TrafficTarget{
			Name:              "baz",
			ConfigurationName: "blah",
			Percent:           37,
		},
		want: nil,
	}, {
		name: "valid with no percent",
		tt: &TrafficTarget{
			Name:              "ooga",
			ConfigurationName: "booga",
		},
		want: nil,
	}, {
		name: "valid with no name",
		tt: &TrafficTarget{
			ConfigurationName: "booga",
			Percent:           100,
		},
		want: nil,
	}, {
		name: "invalid with both",
		tt: &TrafficTarget{
			RevisionName:      "foo",
			ConfigurationName: "bar",
		},
		want: &FieldError{
			Message: "Expected exactly one, got both",
			Paths:   []string{"revisionName", "configurationName"},
		},
	}, {
		name: "invalid with neither",
		tt: &TrafficTarget{
			Name:    "foo",
			Percent: 100,
		},
		want: &FieldError{
			Message: "Expected exactly one, got neither",
			Paths:   []string{"revisionName", "configurationName"},
		},
	}, {
		name: "invalid percent too low",
		tt: &TrafficTarget{
			RevisionName: "foo",
			Percent:      -5,
		},
		want: errInvalidValue("-5", "percent"),
	}, {
		name: "invalid percent too high",
		tt: &TrafficTarget{
			RevisionName: "foo",
			Percent:      101,
		},
		want: errInvalidValue("101", "percent"),
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.tt.Validate()
			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("Validate (-want, +got) = %v", diff)
			}
		})
	}
}
