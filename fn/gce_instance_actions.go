package function

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	asset "cloud.google.com/go/asset/apiv1"
	"github.com/cloudevents/sdk-go/v2/event"
	"google.golang.org/api/iterator"
	assetpb "google.golang.org/genproto/googleapis/cloud/asset/v1"
)

var (
	actions map[string]func(context.Context, *assetpb.SearchAllResourcesRequest) error = map[string]func(context.Context, *assetpb.SearchAllResourcesRequest) error{
		"start": SearchStart,
		"stop":  SearchStop,
	}
	/* commands map[string]interface{} = map[string]interface{}{}*/
)

type ActionsSearch struct {
	SearchRequest *assetpb.SearchAllResourcesRequest `json:"search"`
	Action        string                             `json:"action"`
}

func processAction(ctx context.Context, actionsSearch ActionsSearch) error {
	action := actions[actionsSearch.Action]
	Logger.Debug(ctx, fmt.Sprintf("Got action %s, search %v", actionsSearch.Action, actionsSearch))
	if action != nil {
		return action(ctx, actionsSearch.SearchRequest)
	} else {
		return fmt.Errorf("got no action %s", actionsSearch.Action)
	}
}

func ActionsPubSub(ctx context.Context, m PubSubMessage) error {
	Logger.Info(ctx, "Got PubSub message")
	actionsSearch := &ActionsSearch{}
	Logger.Info(ctx, string(m.Data))
	if err := json.Unmarshal(m.Data, &actionsSearch); err != nil {
		Logger.Info(ctx, fmt.Sprintf("Error: could not unmarshall to search request %v\n", err))
	}
	if err := processAction(ctx, *actionsSearch); err != nil {
		Logger.Error(ctx, err.Error())
	}
	// We don't want retries atm
	return nil
}

func ActionsEvent(ctx context.Context, ev event.Event) error {
	Logger.Info(ctx, fmt.Sprintf("Got CloudEvent : id = %s", ev.Context.GetID())) // %+v", ev))
	// TODO
	eventData := &PubSubEventData{}
	if err := ev.DataAs(eventData); err != nil {
		return fmt.Errorf("error parsing event payload : %w", err)
	}
	actionsSearch := &ActionsSearch{}
	json.Unmarshal(eventData.Message.Data, &actionsSearch)
	// This works if we call the function directly
	/*
		if err := ev.DataAs(actionsSearch); err != nil {
			return fmt.Errorf("error parsing event payload : %w", err)
		}
	*/
	Logger.Info(ctx, fmt.Sprintf("actionsSearch = %+v", actionsSearch))
	if err := processAction(ctx, *actionsSearch); err != nil {
		Logger.Error(ctx, err.Error())
	}
	// We don't want retries atm
	return nil
	// return processAction(ctx, *actionsSearch)
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
			return err
		}
		err = fn(ctx, resource)
		if err != nil {
			log.Printf("%v", err)
		}
		// fmt.Println(resource)
	}
	return nil
}
