package function

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	asset "cloud.google.com/go/asset/apiv1"
	"github.com/cloudevents/sdk-go/v2/event"
	"google.golang.org/api/iterator"
	assetpb "google.golang.org/genproto/googleapis/cloud/asset/v1"
)

var (
	actions map[string]func(context.Context, *assetpb.SearchAllResourcesRequest, map[string]string) error = map[string]func(context.Context, *assetpb.SearchAllResourcesRequest, map[string]string) error{
		"start": func(ctx context.Context, req *assetpb.SearchAllResourcesRequest, params map[string]string) error {
			return Search(ctx, req, Start, params)
		},
		"stop": func(ctx context.Context, req *assetpb.SearchAllResourcesRequest, params map[string]string) error {
			return Search(ctx, req, Stop, params)
		},
		"start-stop": func(ctx context.Context, req *assetpb.SearchAllResourcesRequest, params map[string]string) error {
			return Search(ctx, req, StartStop, params)
		},
	}
	/* commands map[string]interface{} = map[string]interface{}{}*/
)

type ActionsSearch struct {
	SearchRequest *assetpb.SearchAllResourcesRequest `json:"search"`
	Action        string                             `json:"action"`
	Params        map[string]string                  `json:"params"`
}

func Now(timezone string) (*time.Time, error) {
	if loc, err := time.LoadLocation(timezone); err != nil {
		return nil, err
	} else {
		t := time.Now().In(loc)
		return &t, nil
	}
}

// Monday is 0! (not Sunday)
func LabelAction(ctx context.Context, label string, weekDay int, hour int) (string, error) {
	var act string
	ssLabels := strings.Split(label, "_")
	if len(ssLabels) == 7 {
		if ssLabels[weekDay] == "-" {
			act = "stop"
		} else {
			fromTo := strings.Split(ssLabels[weekDay], "-")
			if len(fromTo) == 2 {
				from, fromErr := strconv.Atoi(fromTo[0]) // ParseInt(fromTo[0], 10, 64)
				to, toErr := strconv.Atoi(fromTo[1])
				if fromErr == nil && toErr == nil {
					if hour >= from && hour < to {
						// fmt.Printf("idx = %d, hour = %d, from = %d, to = %d\n", weekDay, hour, from, to /*len(fromTo)*/)
						act = "start"
					} else {
						act = "stop"
					}
				}
				// fmt.Printf("idx = %d, hour = %d, from = %d, to = %d\n", weekDay, hour, from, to /*len(fromTo)*/)
			}
		}

	} else {
		return "", fmt.Errorf("unexpected amount of ranges, got %d, expected 7", len(ssLabels))
	}
	return act, nil

}

func processAction(ctx context.Context, actionsSearch ActionsSearch) error {
	action := actions[actionsSearch.Action]
	Logger.Debug(ctx, fmt.Sprintf("Got action %s, search %v", actionsSearch.Action, actionsSearch))
	if action != nil {
		return action(ctx, actionsSearch.SearchRequest, actionsSearch.Params)
	} else {
		return fmt.Errorf("got no action %s", actionsSearch.Action)
	}
}

func ActionsPubSub(ctx context.Context, m PubSubMessage) error {
	Logger.Info(ctx, "Got PubSub message")
	actionsSearch := &ActionsSearch{}
	Logger.Info(ctx, string(m.Data))
	if err := json.Unmarshal(m.Data, &actionsSearch); err != nil {
		Logger.Error(ctx, fmt.Sprintf("Error: could not unmarshall to search request %v\n", err))
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

func Stop(ctx context.Context, instance *assetpb.ResourceSearchResult, params map[string]string) error {
	return stopInstance(ctx, instance.Project[9:], instance.Location, instance.DisplayName, false)
}

func Start(ctx context.Context, instance *assetpb.ResourceSearchResult, params map[string]string) error {
	return startInstance(ctx, instance.Project[9:], instance.Location, instance.DisplayName, false)
}

func StartStop(ctx context.Context, instance *assetpb.ResourceSearchResult, params map[string]string) error {
	statesLabel := instance.Labels[params["label"]]
	timezone := params["timezone"]
	Logger.Info(ctx, fmt.Sprintf("StartStop : label = %s, timezone = %s, params[\"label\"] = %s", statesLabel, timezone, params["label"]))
	if now, err := Now(timezone); err == nil {
		// Monday is 0! (not Sunday)
		weekday := (int(now.Weekday()) + 6) % 7
		if act, err := LabelAction(ctx, statesLabel, weekday, now.Hour()); err == nil {
			state := instance.State
			Logger.Info(ctx, fmt.Sprintf("labels[\"%s\"] = %s, action = %s, weekday = %d, hour = %d, timezone = %v, state = %s", params["label"], statesLabel, act, weekday, now.Hour(), timezone, state))
			if state == "RUNNING" && act == "stop" {
				return stopInstance(ctx, instance.Project[9:], instance.Location, instance.DisplayName, false)
			} else if state == "TERMINATED" && act == "start" {
				return startInstance(ctx, instance.Project[9:], instance.Location, instance.DisplayName, false)
			} else {
				Logger.Info(ctx, fmt.Sprintf("No action required for project = %s, instance = %s, label = %s, state = %s", instance.Project[9:], instance.DisplayName, statesLabel, state))
			}
		} else {
			return err
		}
	} else {
		return err
	}
	return nil
}

func Search(ctx context.Context,
	req *assetpb.SearchAllResourcesRequest,
	fn func(ctx context.Context, instance *assetpb.ResourceSearchResult, params map[string]string) error,
	params map[string]string,
) error {
	// Unsupported field: 'assetType'. Supported fields include: 'name', 'displayName', 'description', 'location', 'networkTags', 'project', 'folders', 'organization', 'parentAssetType', 'parentFullResourceName', 'labels', 'labels.<key>', 'kmsKey', 'createTime', 'updateTime', 'state', 'tagKeys', 'tagValues', 'tagValueIds'. For more details on how to construct a query, please read: https://cloud.google.com/asset-inventory/docs/searching-resources#how_to_construct_a_query.
	Logger.Info(ctx, fmt.Sprintf("Got search request %v", req))
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
		Logger.Info(ctx, fmt.Sprintf("Got resource %v", resource))

		err = fn(ctx, resource, params)
		if err != nil {
			Logger.Error(ctx, fmt.Sprintf("%v", err))
		}
		// fmt.Println(resource)
	}
	return nil
}
