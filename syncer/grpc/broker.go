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
package grpc

import (
	"context"
	"sync"

	"github.com/apache/servicecomb-service-center/syncer/datacenter"
	"github.com/apache/servicecomb-service-center/syncer/grpc/client"
	"github.com/apache/servicecomb-service-center/syncer/grpc/server"
	pb "github.com/apache/servicecomb-service-center/syncer/proto"
)

type Broker struct {
	svr     *server.Server
	clients map[string]*client.Client
	lock    sync.RWMutex
}

func NewBroker(addr string, store datacenter.Store) *Broker {
	return &Broker{
		svr:     server.NewServer(addr, store),
		clients: map[string]*client.Client{},
	}
}

func (b *Broker) Run() {
	b.svr.Run()
}

func (b *Broker) Stop() {
	b.svr.Stop()
}

func (b *Broker) Pull(ctx context.Context, addr string) (*pb.SyncData, error) {
	cli := b.getClient(addr)
	return cli.Pull(ctx)
}

func (b *Broker) getClient(addr string) *client.Client {
	b.lock.RLock()
	cli, ok := b.clients[addr]
	b.lock.RUnlock()
	if !ok {
		nc, err := client.NewClient(addr)
		if err != nil {
			return nil
		}
		cli = nc
		b.lock.Lock()
		b.clients[addr] = cli
		b.lock.Unlock()
	}
	return cli
}
