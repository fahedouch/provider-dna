/*
Copyright 2019 The Crossplane Authors.

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

package dna

import (
	"context"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/crossplane/crossplane-runtime/pkg/event"
	"github.com/crossplane/crossplane-runtime/pkg/logging"
	"github.com/crossplane/provider-dna/pkg/runtime/resource"
)

const (
	reconcileGracePeriod = 30 * time.Second
	reconcileTimeout     = 1 * time.Minute

	defaultpollInterval = 1 * time.Minute
	defaultGracePeriod  = 30 * time.Second
)

// Error strings.
//TODO

// Event reasons.
//TODO

// ControllerName returns the recommended name for controllers that use this
// package to reconcile a particular kind of dna resource.
func ControllerName(kind string) string {
	return "managed/" + strings.ToLower(kind)
}

// A Client manages the lifecycle of an crossplane resource.
// None of the calls here should be blocking. All of the calls should be
// idempotent. For example, Create call should not return AlreadyExists error
// if it's called again with the same parameters or Delete call should not
// return error if there is an ongoing deletion or resource does not exist.
type Client interface {
	// Observe the crossplane resource the supplied Dna resource
	// represents, if any. Observe implementations must not modify the
	// crossplane resource, but may update the supplied Dna resource to
	// reflect the state of the crossplane resource. Status modifications are
	// automatically persisted unless ResourceLateInitialized is true - see
	// ResourceLateInitialized for more detail.
	Observe(ctx context.Context, mg resource.Dna) (CrossplaneObservation, error)

	// Create an crossplane resource per the specifications of the supplied
	// Dna resource. Called when Observe reports that the associated
	// crossplane resource does not exist. Create implementations may update
	// dna resource annotations, and those updates will be persisted.
	// All other updates will be discarded.
	Create(ctx context.Context, mg resource.Dna) (CrossplaneCreation, error)

	// Update the crossplane resource represented by the supplied Dna
	// resource, if necessary. Called unless Observe reports that the
	// associated crossplane resource is up to date.
	Update(ctx context.Context, mg resource.Dna) (CrossplaneUpdate, error)

	// Delete the crossplane resource upon deletion of its associated Dna
	// resource. Called when the managed resource has been deleted.
	Delete(ctx context.Context, mg resource.Dna) error
}

// An CrossplaneObservation is the result of an observation of an crossplane
// resource.
type CrossplaneObservation struct {
	// ResourceExists must be true if a corresponding crossplane resource exists
	// for the dna resource. Typically this is proven by the presence of
	// crossplane resource of the expected kind whose unique identifier matches
	// the dna resource's crossplane name. Crossplane uses this information to
	// determine whether it needs to create or delete the crossplane resource.
	ResourceExists bool

	// ResourceUpToDate should be true if the corresponding external resource
	// appears to be up-to-date - i.e. updating the external resource to match
	// the desired state of the managed resource would be a no-op. Keep in mind
	// that often only a subset of external resource fields can be updated.
	// Crossplane uses this information to determine whether it needs to update
	// the external resource.
	ResourceUpToDate bool

	// ResourceLateInitialized should be true if the managed resource's spec was
	// updated during its observation. A Crossplane provider may update a
	// managed resource's spec fields after it is created or updated, as long as
	// the updates are limited to setting previously unset fields, and adding
	// keys to maps. Crossplane uses this information to determine whether
	// changes to the spec were made during observation that must be persisted.
	// Note that changes to the spec will be persisted before changes to the
	// status, and that pending changes to the status may be lost when the spec
	// is persisted. Status changes will be persisted by the first subsequent
	// observation that _does not_ late initialize the managed resource, so it
	// is important that Observe implementations do not late initialize the
	// resource every time they are called.
	ResourceLateInitialized bool

	// Diff is a Debug level message that is sent to the reconciler when
	// there is a change in the observed Managed Resource. It is useful for
	// finding where the observed diverges from the desired state.
	// The string should be a cmp.Diff that details the difference.
	Diff string
}

// An CrossplaneCreation is the result of the creation of an crossplane resource.
type CrossplaneCreation struct {
	// ExternalNameAssigned should be true if the Create operation resulted
	// in a change in the resource's external name. This is typically only
	// needed for external resource's with unpredictable external names that
	// are returned from the API at create time.
	//
	// Deprecated: The managed.Reconciler no longer needs to be informed
	// when an external name is assigned by the Create operation. It will
	// automatically detect and handle external name assignment.
	CrossplaneNameAssigned bool
}

// An CrossplaneUpdate is the result of an update to an crossplane resource.
type CrossplaneUpdate struct {
}

// A Reconciler reconciles managed resources by creating and managing the
// lifecycle of an external resource, i.e. a resource in an external system such
// as a cloud provider API. Each controller must watch the managed resource kind
// for which it is responsible.
type Reconciler struct {
	client client.Client
	newDna func() resource.Dna

	pollInterval        time.Duration
	timeout             time.Duration
	creationGracePeriod time.Duration

	log    logging.Logger
	record event.Recorder
}

// A ReconcilerOption configures a Reconciler.
type ReconcilerOption func(*Reconciler)

// WithTimeout specifies the timeout duration cumulatively for all the calls happen
// in the reconciliation function. In case the deadline exceeds, reconciler will
// still have some time to make the necessary calls to report the error such as
// status update.
func WithTimeout(duration time.Duration) ReconcilerOption {
	return func(r *Reconciler) {
		r.timeout = duration
	}
}

// WithPollInterval specifies how long the Reconciler should wait before queueing
// a new reconciliation after a successful reconcile. The Reconciler requeues
// after a specified duration when it is not actively waiting for an external
// operation, but wishes to check whether an existing external resource needs to
// be synced to its Crossplane Managed resource.
func WithPollInterval(after time.Duration) ReconcilerOption {
	return func(r *Reconciler) {
		r.pollInterval = after
	}
}

// WithCreationGracePeriod configures an optional period during which we will
// wait for the external API to report that a newly created external resource
// exists. This allows us to tolerate eventually consistent APIs that do not
// immediately report that newly created resources exist when queried. All
// resources have a 30 second grace period by default.
func WithCreationGracePeriod(d time.Duration) ReconcilerOption {
	return func(r *Reconciler) {
		r.creationGracePeriod = d
	}
}

// WithLogger specifies how the Reconciler should log messages.
func WithLogger(l logging.Logger) ReconcilerOption {
	return func(r *Reconciler) {
		r.log = l
	}
}

// WithRecorder specifies how the Reconciler should record events.
func WithRecorder(er event.Recorder) ReconcilerOption {
	return func(r *Reconciler) {
		r.record = er
	}
}

// NewReconciler returns a Reconciler that reconciles dna resources of the
// supplied DnaKind with resources in an external system such as a crossplane API.
// It panics if asked to reconcile a dna resource kind that is
// not registered with the supplied manager's runtime.Scheme. The returned
// Reconciler reconciles with a dummy, no-op 'external system' by default;
// callers should supply a connector that returns a client
// capable of managing resources in a real system.
func NewReconciler(m manager.Manager, of resource.DnaKind, o ...ReconcilerOption) *Reconciler {
	nm := func() resource.Dna {
		return resource.MustCreateObject(schema.GroupVersionKind(of), m.GetScheme()).(resource.Dna)
	}

	// Panic early if we've been asked to reconcile a resource kind that has not
	// been registered with our controller manager's scheme.
	_ = nm()

	r := &Reconciler{
		client:              m.GetClient(),
		newDna:              nm,
		pollInterval:        defaultpollInterval,
		creationGracePeriod: defaultGracePeriod,
		timeout:             reconcileTimeout,
		log:                 logging.NewNopLogger(),
		record:              event.NewNopRecorder(),
	}

	for _, ro := range o {
		ro(r)
	}

	return r
}

// Reconcile a managed resource with an crossplane resource.
func (r *Reconciler) Reconcile(_ context.Context, req reconcile.Request) (reconcile.Result, error) { // nolint:gocyclo
	// NOTE(negz): This method is a well over our cyclomatic complexity goal.
	// Be wary of adding additional complexity.

	log := r.log.WithValues("request", req)
	log.Debug("Reconciling")

	//TODO GENERIC Reconcile logic here
	_, cancel := context.WithTimeout(context.Background(), r.timeout+reconcileGracePeriod)
	defer cancel()

	return reconcile.Result{}, nil
}
