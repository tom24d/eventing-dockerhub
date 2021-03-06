/*
Copyright 2020 The Knative Authors

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

// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/tom24d/eventing-dockerhub/pkg/apis/sources/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// DockerHubSourceLister helps list DockerHubSources.
// All objects returned here must be treated as read-only.
type DockerHubSourceLister interface {
	// List lists all DockerHubSources in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.DockerHubSource, err error)
	// DockerHubSources returns an object that can list and get DockerHubSources.
	DockerHubSources(namespace string) DockerHubSourceNamespaceLister
	DockerHubSourceListerExpansion
}

// dockerHubSourceLister implements the DockerHubSourceLister interface.
type dockerHubSourceLister struct {
	indexer cache.Indexer
}

// NewDockerHubSourceLister returns a new DockerHubSourceLister.
func NewDockerHubSourceLister(indexer cache.Indexer) DockerHubSourceLister {
	return &dockerHubSourceLister{indexer: indexer}
}

// List lists all DockerHubSources in the indexer.
func (s *dockerHubSourceLister) List(selector labels.Selector) (ret []*v1alpha1.DockerHubSource, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.DockerHubSource))
	})
	return ret, err
}

// DockerHubSources returns an object that can list and get DockerHubSources.
func (s *dockerHubSourceLister) DockerHubSources(namespace string) DockerHubSourceNamespaceLister {
	return dockerHubSourceNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// DockerHubSourceNamespaceLister helps list and get DockerHubSources.
// All objects returned here must be treated as read-only.
type DockerHubSourceNamespaceLister interface {
	// List lists all DockerHubSources in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.DockerHubSource, err error)
	// Get retrieves the DockerHubSource from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha1.DockerHubSource, error)
	DockerHubSourceNamespaceListerExpansion
}

// dockerHubSourceNamespaceLister implements the DockerHubSourceNamespaceLister
// interface.
type dockerHubSourceNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all DockerHubSources in the indexer for a given namespace.
func (s dockerHubSourceNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.DockerHubSource, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.DockerHubSource))
	})
	return ret, err
}

// Get retrieves the DockerHubSource from the indexer for a given namespace and name.
func (s dockerHubSourceNamespaceLister) Get(name string) (*v1alpha1.DockerHubSource, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("dockerhubsource"), name)
	}
	return obj.(*v1alpha1.DockerHubSource), nil
}
