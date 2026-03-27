// Copyright 2026 Google LLC
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

package api

import (
	"slices"
	"testing"

	"github.com/gorilla/mux"
	"google.golang.org/adk/cmd/launcher"
)

func TestSetupSubrouters_TriggerSourcesValidation(t *testing.T) {
	tests := []struct {
		name           string
		triggerSources string
		wantErr        bool
		wantSources    []string
	}{
		{
			name:           "empty trigger sources",
			triggerSources: "",
			wantErr:        false,
			wantSources:    nil,
		},
		{
			name:           "valid trigger sources single",
			triggerSources: "pubsub",
			wantErr:        false,
			wantSources:    []string{"pubsub"},
		},
		{
			name:           "valid trigger sources multiple",
			triggerSources: "pubsub,bq",
			wantErr:        false,
			wantSources:    []string{"pubsub", "bq"},
		},
		{
			name:           "valid trigger sources all",
			triggerSources: "pubsub,bq,eventarc",
			wantErr:        false,
			wantSources:    []string{"pubsub", "bq", "eventarc"},
		},
		{
			name:           "invalid trigger source",
			triggerSources: "invalid",
			wantErr:        true,
			wantSources:    nil,
		},
		{
			name:           "mixed valid and invalid",
			triggerSources: "pubsub,invalid,bq",
			wantErr:        true,
			wantSources:    nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			a := &apiLauncher{
				config: &apiConfig{
					triggerSources: tc.triggerSources,
				},
			}
			router := mux.NewRouter()
			config := &launcher.Config{}

			err := a.SetupSubrouters(router, config)
			if tc.wantErr {
				if err == nil {
					t.Errorf("SetupSubrouters() error = nil, wantErr %v", tc.wantErr)
				}
			} else {
				if err != nil {
					t.Errorf("SetupSubrouters() error = %v, wantErr %v", err, tc.wantErr)
				}
				if !slices.Equal(config.TriggerSources, tc.wantSources) {
					t.Errorf("SetupSubrouters() config.TriggerSources = %v, want %v", config.TriggerSources, tc.wantSources)
				}
			}
		})
	}
}
