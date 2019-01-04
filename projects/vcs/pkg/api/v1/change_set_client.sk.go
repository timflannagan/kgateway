// Code generated by solo-kit. DO NOT EDIT.

package v1

import (
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/errors"
)

type ChangeSetClient interface {
	BaseClient() clients.ResourceClient
	Register() error
	Read(namespace, name string, opts clients.ReadOpts) (*ChangeSet, error)
	Write(resource *ChangeSet, opts clients.WriteOpts) (*ChangeSet, error)
	Delete(namespace, name string, opts clients.DeleteOpts) error
	List(namespace string, opts clients.ListOpts) (ChangeSetList, error)
	Watch(namespace string, opts clients.WatchOpts) (<-chan ChangeSetList, <-chan error, error)
}

type changeSetClient struct {
	rc clients.ResourceClient
}

func NewChangeSetClient(rcFactory factory.ResourceClientFactory) (ChangeSetClient, error) {
	return NewChangeSetClientWithToken(rcFactory, "")
}

func NewChangeSetClientWithToken(rcFactory factory.ResourceClientFactory, token string) (ChangeSetClient, error) {
	rc, err := rcFactory.NewResourceClient(factory.NewResourceClientParams{
		ResourceType: &ChangeSet{},
		Token:        token,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "creating base ChangeSet resource client")
	}
	return NewChangeSetClientWithBase(rc), nil
}

func NewChangeSetClientWithBase(rc clients.ResourceClient) ChangeSetClient {
	return &changeSetClient{
		rc: rc,
	}
}

func (client *changeSetClient) BaseClient() clients.ResourceClient {
	return client.rc
}

func (client *changeSetClient) Register() error {
	return client.rc.Register()
}

func (client *changeSetClient) Read(namespace, name string, opts clients.ReadOpts) (*ChangeSet, error) {
	opts = opts.WithDefaults()
	resource, err := client.rc.Read(namespace, name, opts)
	if err != nil {
		return nil, err
	}
	return resource.(*ChangeSet), nil
}

func (client *changeSetClient) Write(changeSet *ChangeSet, opts clients.WriteOpts) (*ChangeSet, error) {
	opts = opts.WithDefaults()
	resource, err := client.rc.Write(changeSet, opts)
	if err != nil {
		return nil, err
	}
	return resource.(*ChangeSet), nil
}

func (client *changeSetClient) Delete(namespace, name string, opts clients.DeleteOpts) error {
	opts = opts.WithDefaults()
	return client.rc.Delete(namespace, name, opts)
}

func (client *changeSetClient) List(namespace string, opts clients.ListOpts) (ChangeSetList, error) {
	opts = opts.WithDefaults()
	resourceList, err := client.rc.List(namespace, opts)
	if err != nil {
		return nil, err
	}
	return convertToChangeSet(resourceList), nil
}

func (client *changeSetClient) Watch(namespace string, opts clients.WatchOpts) (<-chan ChangeSetList, <-chan error, error) {
	opts = opts.WithDefaults()
	resourcesChan, errs, initErr := client.rc.Watch(namespace, opts)
	if initErr != nil {
		return nil, nil, initErr
	}
	changesetsChan := make(chan ChangeSetList)
	go func() {
		for {
			select {
			case resourceList := <-resourcesChan:
				changesetsChan <- convertToChangeSet(resourceList)
			case <-opts.Ctx.Done():
				close(changesetsChan)
				return
			}
		}
	}()
	return changesetsChan, errs, nil
}

func convertToChangeSet(resources resources.ResourceList) ChangeSetList {
	var changeSetList ChangeSetList
	for _, resource := range resources {
		changeSetList = append(changeSetList, resource.(*ChangeSet))
	}
	return changeSetList
}
