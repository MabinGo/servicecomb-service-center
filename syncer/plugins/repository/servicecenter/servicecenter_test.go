/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package servicecenter

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/apache/servicecomb-service-center/syncer/plugins"
	"github.com/apache/servicecomb-service-center/syncer/plugins/repository"
	"github.com/apache/servicecomb-service-center/syncer/test/datacenter/servicecenter"
)

func TestClient_GetAll(t *testing.T) {
	svr, repo := newServiceCenterRepo(t)
	_, err := repo.GetAll(context.Background())
	if err != nil {
		t.Errorf("get all from %s server failed, error: %s", PluginName, err)
	}

	svr.Close()
	_, err = repo.GetAll(context.Background())
	if err != nil {
		t.Logf("get all from %s server failed, error: %s", PluginName, err)
	}
}

func newServiceCenterRepo(t *testing.T) (*httptest.Server, repository.Repository) {
	plugins.SetPluginConfig(plugins.PluginRepository.String(), PluginName)
	adaptor := plugins.Plugins().Repository()
	if adaptor == nil {
		t.Errorf("get repository adaptor %s failed", PluginName)
	}
	svr := servicecenter.NewMockServer()
	if svr == nil {
		t.Error("new httptest server failed")
	}
	repo, err := adaptor.New([]string{svr.URL})
	if err != nil {
		t.Errorf("new repository %s failed, error: %s", PluginName, err)
	}
	return svr, repo
}
