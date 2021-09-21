/*
Copyright 2020 The Crossplane Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package compute

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"k8s.io/client-go/util/workqueue"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"

	"github.com/crossplane/crossplane-runtime/pkg/event"
	"github.com/crossplane/crossplane-runtime/pkg/logging"
	"github.com/crossplane/crossplane-runtime/pkg/ratelimiter"
	"github.com/crossplane/provider-dna/pkg/runtime/reconciler/dna"
	"github.com/crossplane/provider-dna/pkg/runtime/resource"

	"github.com/crossplane/provider-dna/apis/compute/v1alpha1"
)

const (
	errNotFirewall = "crossplane resource is not a Firewall custom resource"
)

// A NoOpService does nothing.
type NoOpService struct{}

var (
	newNoOpService = func(_ []byte) (interface{}, error) { return &NoOpService{}, nil }
)

// Setup adds a controller that reconciles Firewall crossplane resources.
func SetupFirewall(mgr ctrl.Manager, l logging.Logger, rl workqueue.RateLimiter) error {
	//Set the Firewall controller name
	name := dna.ControllerName(v1alpha1.FirewallGroupKind)

	//Limit how frequently requests may be queued in the controller
	o := controller.Options{
		RateLimiter: ratelimiter.NewDefaultManagedRateLimiter(rl),
	}

	//Initalize the reconciler that will reconcile crossplane resource to DNA Firewall resource
	r := dna.NewReconciler(mgr,
		resource.DnaKind(v1alpha1.FirewallGroupVersionKind),
		dna.WithLogger(l.WithValues("controller", name)),
		dna.WithRecorder(event.NewAPIRecorder(mgr.GetEventRecorderFor(name))))

	return ctrl.NewControllerManagedBy(mgr).
		Named(name).
		WithOptions(o).
		For(&v1alpha1.Firewall{}).
		Complete(r)
}

// A firewallConnector is struct that connect to crossplane firewall resources
// this is struct wrap the kube client for the moment but it can wrap other clients in the furture like externalClient
type firewallConnector struct {
	kube client.Client
}

func (c *firewallConnector) Observe(ctx context.Context, mg resource.Dna) (dna.CrossplaneObservation, error) {
	// dr == dna resource
	dr, ok := mg.(*v1alpha1.Firewall)
	if !ok {
		return dna.CrossplaneObservation{}, errors.New(errNotFirewall)
	}

	//TODO introduce Observing logic here
	// These fmt statements should be removed in the real implementation.
	fmt.Printf("Observing: %+v", dr)

	return dna.CrossplaneObservation{
		// Return false when the crossplane resource does not exist. This lets
		// the managed resource reconciler know that it needs to call Create to
		// (re)create the resource, or that it has successfully been deleted.
		ResourceExists: true,

		// Return false when the crossplane resource exists, but it not up to date
		// with the desired managed resource state. This lets the managed
		// resource reconciler know that it needs to call Update.
		ResourceUpToDate: true,
	}, nil
}

func (c *firewallConnector) Create(ctx context.Context, mg resource.Dna) (dna.CrossplaneCreation, error) {
	// dr == dna resource
	dr, ok := mg.(*v1alpha1.Firewall)
	if !ok {
		return dna.CrossplaneCreation{}, errors.New(errNotFirewall)
	}

	//TODO introduce creating logic here
	fmt.Printf("Creating: %+v", dr)

	return dna.CrossplaneCreation{}, nil
}

func (c *firewallConnector) Update(ctx context.Context, mg resource.Dna) (dna.CrossplaneUpdate, error) {
	// dr == dna resource
	dr, ok := mg.(*v1alpha1.Firewall)
	if !ok {
		return dna.CrossplaneUpdate{}, errors.New(errNotFirewall)
	}
	//TODO introduce updating logic here
	fmt.Printf("Updating: %+v", dr)

	return dna.CrossplaneUpdate{}, nil
}

func (c *firewallConnector) Delete(ctx context.Context, mg resource.Dna) error {
	// dr == dna resource
	dr, ok := mg.(*v1alpha1.Firewall)
	if !ok {
		return errors.New(errNotFirewall)
	}

	fmt.Printf("Deleting: %+v", dr)

	return nil
}
