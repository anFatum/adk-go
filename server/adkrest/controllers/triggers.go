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

package controllers

import "net/http"

// TriggersAPIController handles the ADK triggered endpoints.
type TriggersAPIController struct{}

// NewTriggersAPIController creates a new TriggersAPIController.
func NewTriggersAPIController() *TriggersAPIController {
	return &TriggersAPIController{}
}

// BQTriggerHandler handles the BigQuery trigger endpoint.
func (c *TriggersAPIController) BQTriggerHandler(w http.ResponseWriter, r *http.Request) {
	Unimplemented(w, r)
}

// PubSubTriggerHandler handles the PubSub trigger endpoint.
func (c *TriggersAPIController) PubSubTriggerHandler(w http.ResponseWriter, r *http.Request) {
	Unimplemented(w, r)
}

// EventarcTriggerHandler handles the Eventarc trigger endpoint.
func (c *TriggersAPIController) EventarcTriggerHandler(w http.ResponseWriter, r *http.Request) {
	Unimplemented(w, r)
}
