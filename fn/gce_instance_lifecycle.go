package function

/*
gcloud asset search-all-resources \
  --scope=organizations/${ORG_ID} \
  --asset-types=compute.googleapis.com/Instance
*/
import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	asset "cloud.google.com/go/asset/apiv1"
	"google.golang.org/api/iterator"
	assetpb "google.golang.org/genproto/googleapis/cloud/asset/v1"
)

func StartPubSub(ctx context.Context, m PubSubMessage) error {
	logger.Info(ctx, "Got PubSub message")
	searchRequest := &assetpb.SearchAllResourcesRequest{}
	log.Println(string(m.Data))
	err := json.Unmarshal(m.Data, &searchRequest)
	if err != nil {
		logger.Info(ctx, fmt.Sprintf("Error: could not unmarshall to search request %v\n", err))
	}
	return SearchStart(ctx, searchRequest)
}

func StopPubSub(ctx context.Context, m PubSubMessage) error {
	logger.Info(ctx, "Got PubSub message")
	searchRequest := &assetpb.SearchAllResourcesRequest{}
	log.Println(string(m.Data))
	err := json.Unmarshal(m.Data, &searchRequest)
	if err != nil {
		logger.Info(ctx, fmt.Sprintf("Error: could not unmarshall to search request %v\n", err))
	}
	return SearchStop(ctx, searchRequest)
}

func SearchStart(ctx context.Context, req *assetpb.SearchAllResourcesRequest) error {
	return Search(ctx, req, Start)
}

func SearchStop(ctx context.Context, req *assetpb.SearchAllResourcesRequest) error {
	return Search(ctx, req, Stop)
}

func Stop(ctx context.Context, instance *assetpb.ResourceSearchResult) error {
	return stopInstance(ctx, instance.Project[9:], instance.Location, instance.DisplayName, false)
}

func Start(ctx context.Context, instance *assetpb.ResourceSearchResult) error {
	return startInstance(ctx, instance.Project[9:], instance.Location, instance.DisplayName, false)
}

func Search(ctx context.Context, req *assetpb.SearchAllResourcesRequest, fn func(ctx context.Context, instance *assetpb.ResourceSearchResult) error) error {
	// Unsupported field: 'assetType'. Supported fields include: 'name', 'displayName', 'description', 'location', 'networkTags', 'project', 'folders', 'organization', 'parentAssetType', 'parentFullResourceName', 'labels', 'labels.<key>', 'kmsKey', 'createTime', 'updateTime', 'state', 'tagKeys', 'tagValues', 'tagValueIds'. For more details on how to construct a query, please read: https://cloud.google.com/asset-inventory/docs/searching-resources#how_to_construct_a_query.
	log.Printf("Got search request %v", req)
	client, err := asset.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("asset.NewClient: %v", err)
	}
	defer client.Close()
	it := client.SearchAllResources(ctx, req)
	for {
		resource, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		err = fn(ctx, resource)
		if err != nil {
			log.Printf("%v", err)
		}
		// fmt.Println(resource)
	}
	return nil
}
