package memory

import (
	"testing"

	scpb "github.com/apache/servicecomb-service-center/server/core/proto"
	"github.com/apache/servicecomb-service-center/syncer/plugins"
	"github.com/apache/servicecomb-service-center/syncer/plugins/storage"
	pb "github.com/apache/servicecomb-service-center/syncer/proto"
)

func TestSyncData(t *testing.T) {
	store := getmemoryStorage(t)
	data := store.GetSyncData()
	if len(data.Services) > 0 {
		t.Error("default sync data was wrong")
	}
	data = getSyncData()
	store.SaveSyncData(data)
	nd := store.GetSyncData()
	if nd == nil {
		t.Error("save sync data failed!")
	}

	store.Stop()
}

func TestSyncMapping(t *testing.T) {
	store := getmemoryStorage(t)
	defer store.Stop()
	nodeName := "testnode"
	data := store.GetSyncMapping(nodeName)
	if len(data) > 0 {
		t.Error("default sync mapping was wrong")
	}
	data = getSyncMapping(nodeName)
	store.SaveSyncMapping(nodeName, data)
	nd := store.GetSyncMapping(nodeName)
	_, ok := nd[nodeName]
	if !ok {
		t.Error("save sync mapping failed!")
	}

	all := store.GetAllMapping()
	for key := range all {
		if key == nodeName {
			return
		}
	}

	t.Errorf("all mapping has not node name %s!", nodeName)
}

func getmemoryStorage(t *testing.T) storage.Repository {
	plugins.SetPluginConfig(plugins.PluginStorage.String(), PluginName)
	store := plugins.Plugins().Storage()
	if store == nil {
		t.Errorf("get storage repository %s failed", PluginName)
	}
	return store
}

func getSyncMapping(nodeName string) pb.SyncMapping {
	return pb.SyncMapping{nodeName: &pb.SyncServiceKey{
		DomainProject: "default/default",
		ServiceID:     "5db1b794aa6f8a875d6e68110260b5491ee7e223",
		InstanceID:    "4d41a637471f11e9888cfa163eca30e0",
	}}
}

func getSyncData() *pb.SyncData {
	return &pb.SyncData{
		Services: []*pb.SyncService{
			{
				DomainProject: "default/default",
				Service: &scpb.MicroService{
					ServiceId:   "5db1b794aa6f8a875d6e68110260b5491ee7e223",
					AppId:       "default",
					ServiceName: "SERVICECENTER",
					Version:     "1.1.0",
					Level:       "BACK",
					Schemas: []string{
						"servicecenter.grpc.api.ServiceCtrl",
						"servicecenter.grpc.api.ServiceInstanceCtrl",
					},
					Status: "UP",
					Properties: map[string]string{
						"allowCrossApp": "true",
					},
					Timestamp:    "1552626180",
					ModTimestamp: "1552626180",
					Environment:  "production",
				},
				Instances: []*scpb.MicroServiceInstance{
					{
						InstanceId: "4d41a637471f11e9888cfa163eca30e0",
						ServiceId:  "5db1b794aa6f8a875d6e68110260b5491ee7e223",
						Endpoints: []string{
							"rest://127.0.0.1:30100/",
						},
						HostName: "testmock",
						Status:   "UP",
						HealthCheck: &scpb.HealthCheck{
							Mode:     "push",
							Interval: 30,
							Times:    3,
						},
						Timestamp:    "1552653537",
						ModTimestamp: "1552653537",
						Version:      "1.1.0",
					},
				},
			},
		},
	}
}
