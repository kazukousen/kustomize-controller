/*
Copyright 2020 The Flux authors

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

package controllers

import (
	"context"
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	kustomizev1 "github.com/fluxcd/kustomize-controller/api/v1beta1"
	"github.com/fluxcd/pkg/runtime/dependency"
	sourcev1 "github.com/fluxcd/source-controller/api/v1beta1"
)

func (r *KustomizationReconciler) requestsForGitRepositoryRevisionChange(obj client.Object) []reconcile.Request {
	repo, ok := obj.(*sourcev1.GitRepository)
	if !ok {
		panic(fmt.Sprintf("Expected a GitRepository but got a %T", obj))
	}
	// If we do not have an artifact, we have no requests to make
	if repo.GetArtifact() == nil {
		return nil
	}

	ctx := context.Background()
	var list kustomizev1.KustomizationList
	if err := r.List(ctx, &list, client.MatchingFields{
		kustomizev1.GitRepositoryIndexKey: ObjectKey(obj).String(),
	}); err != nil {
		return nil
	}
	var dd []dependency.Dependent
	for _, d := range list.Items {
		// If the revision of the artifact equals to the last attempted revision,
		// we should not make a request for this Kustomization
		if repo.GetArtifact().Revision == d.Status.LastAttemptedRevision {
			continue
		}
		dd = append(dd, d)
	}
	sorted, err := dependency.Sort(dd)
	if err != nil {
		return nil
	}
	reqs := make([]reconcile.Request, len(sorted), len(sorted))
	for i := range sorted {
		reqs[i].NamespacedName.Name = sorted[i].Name
		reqs[i].NamespacedName.Namespace = sorted[i].Namespace
	}
	return reqs
}

func (r *KustomizationReconciler) indexByGitRepository(o client.Object) []string {
	k, ok := o.(*kustomizev1.Kustomization)
	if !ok {
		panic(fmt.Sprintf("Expected a Kustomization, got %T", o))
	}

	if k.Spec.SourceRef.Kind == sourcev1.GitRepositoryKind {
		namespace := k.GetNamespace()
		if k.Spec.SourceRef.Namespace != "" {
			namespace = k.Spec.SourceRef.Namespace
		}
		return []string{fmt.Sprintf("%s/%s", namespace, k.Spec.SourceRef.Name)}
	}

	return nil
}

func (r *KustomizationReconciler) requestsForBucketRevisionChange(obj client.Object) []reconcile.Request {
	bucket, ok := obj.(*sourcev1.Bucket)
	if !ok {
		panic(fmt.Sprintf("Expected a Bucket but got a %T", obj))
	}
	// If we do not have an artifact, we have no requests to make
	if bucket.GetArtifact() == nil {
		return nil
	}

	ctx := context.Background()
	var list kustomizev1.KustomizationList
	if err := r.List(ctx, &list, client.MatchingFields{
		kustomizev1.BucketIndexKey: ObjectKey(obj).String(),
	}); err != nil {
		return nil
	}
	var dd []dependency.Dependent
	for _, d := range list.Items {
		// If the revision of the artifact equals to the last attempted revision,
		// we should not make a request for this Kustomization
		if bucket.GetArtifact().Revision == d.Status.LastAttemptedRevision {
			continue
		}
		dd = append(dd, d)
	}
	sorted, err := dependency.Sort(dd)
	if err != nil {
		return nil
	}
	reqs := make([]reconcile.Request, len(sorted), len(sorted))
	for i := range sorted {
		reqs[i].NamespacedName.Name = sorted[i].Name
		reqs[i].NamespacedName.Namespace = sorted[i].Namespace
	}
	return reqs
}

func (r *KustomizationReconciler) indexByBucket(o client.Object) []string {
	k, ok := o.(*kustomizev1.Kustomization)
	if !ok {
		panic(fmt.Sprintf("Expected a Kustomization, got %T", o))
	}

	if k.Spec.SourceRef.Kind == sourcev1.BucketKind {
		namespace := k.GetNamespace()
		if k.Spec.SourceRef.Namespace != "" {
			namespace = k.Spec.SourceRef.Namespace
		}
		return []string{fmt.Sprintf("%s/%s", namespace, k.Spec.SourceRef.Name)}
	}

	return nil
}
