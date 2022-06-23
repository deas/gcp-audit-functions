package function

import (
	"context"
	"fmt"
	"log"

	compute "cloud.google.com/go/compute/apiv1"
	"google.golang.org/api/option"
	computepb "google.golang.org/genproto/googleapis/cloud/compute/v1"
)

var (
	Opts []option.ClientOption
)

func logAction(action string, projectID, zone, instance string) {
	log.Printf("%s : project = %s, zone = %s, instance = %s", action, projectID, zone, instance)
}

func startInstance(ctx context.Context, projectID string, zone string, instance string, wait bool) error {
	logAction("Start", projectID, zone, instance)
	client, err := compute.NewInstancesRESTClient(ctx, Opts...)
	if err != nil {
		return fmt.Errorf("NewInstancesRESTClient: %v", err)
	}
	defer client.Close()

	req := &computepb.StartInstanceRequest{
		Project:  projectID,
		Zone:     zone,
		Instance: instance,
	}

	op, err := client.Start(ctx, req)
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

	client, err := compute.NewInstancesRESTClient(ctx, Opts...)
	if err != nil {
		return fmt.Errorf("NewInstancesRESTClient: %v", err)
	}
	defer client.Close()

	req := &computepb.StopInstanceRequest{
		Project:  projectID,
		Zone:     zone,
		Instance: instance,
	}

	op, err := client.Stop(ctx, req)
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
