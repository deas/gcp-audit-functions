package function

import (
	"context"
	"fmt"
	"log"
	"os"

	compute "cloud.google.com/go/compute/apiv1"
	"google.golang.org/api/option"
	computepb "google.golang.org/genproto/googleapis/cloud/compute/v1"
)

func NewOpts() []option.ClientOption {
	var opts []option.ClientOption
	// Don't up in Pipeline : google: could not find default credentials.
	if _, noAuth := os.LookupEnv(fmt.Sprintf("%s_%s", EnvPrefix, "NO_AUTH")); noAuth {
		opts = []option.ClientOption{
			option.WithoutAuthentication(),
		}
	} else {
		opts = []option.ClientOption{}
	}
	return opts
}

func logAction(action string, projectID, zone, instance string) {
	log.Printf("%s : project = %s, zone = %s, instance = %s", action, projectID, zone, instance)
}

func startInstance(ctx context.Context, projectID string, zone string, instance string, wait bool) error {
	logAction("Start", projectID, zone, instance)
	instancesClient, err := compute.NewInstancesRESTClient(ctx) // TODO - move out
	if err != nil {
		return fmt.Errorf("NewInstancesRESTClient: %v", err)
	}
	defer instancesClient.Close()

	req := &computepb.StartInstanceRequest{
		Project:  projectID,
		Zone:     zone,
		Instance: instance,
	}

	op, err := instancesClient.Start(ctx, req)
	if err != nil {
		return fmt.Errorf("unable to start instance: %v", err)
	}
	if wait {
		if err = op.Wait(ctx); err != nil {
			return fmt.Errorf("unable to wait for the operation: %v", err)
		}
	}

	return nil
}

func stopInstance(ctx context.Context, projectID string, zone string, instance string, wait bool) error {
	logAction("Stop", projectID, zone, instance)

	instancesClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewInstancesRESTClient: %v", err)
	}
	defer instancesClient.Close()

	req := &computepb.StopInstanceRequest{
		Project:  projectID,
		Zone:     zone,
		Instance: instance,
	}

	op, err := instancesClient.Stop(ctx, req)
	if err != nil {
		return fmt.Errorf("unable to stop instance: %v", err)
	}

	if wait {
		if err = op.Wait(ctx); err != nil {
			return fmt.Errorf("unable to wait for the operation: %v", err)
		}
	}

	return nil
}
