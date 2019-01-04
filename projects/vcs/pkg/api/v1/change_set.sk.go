// Code generated by solo-kit. DO NOT EDIT.

package v1

import (
	"sort"

	"github.com/gogo/protobuf/proto"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube/crd"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/utils/hashutils"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// TODO: modify as needed to populate additional fields
func NewChangeSet(namespace, name string) *ChangeSet {
	return &ChangeSet{
		Metadata: core.Metadata{
			Name:      name,
			Namespace: namespace,
		},
	}
}

func (r *ChangeSet) SetStatus(status core.Status) {
	r.Status = status
}

func (r *ChangeSet) SetMetadata(meta core.Metadata) {
	r.Metadata = meta
}

func (r *ChangeSet) Hash() uint64 {
	metaCopy := r.GetMetadata()
	metaCopy.ResourceVersion = ""
	return hashutils.HashAll(
		metaCopy,
		r.Branch,
		r.PendingAction,
		r.Description,
		r.EditCount,
		r.UserId,
		r.RootCommit,
		r.RootDescription,
		r.ErrorMsg,
		r.Data,
	)
}

type ChangeSetList []*ChangeSet
type ChangesetsByNamespace map[string]ChangeSetList

// namespace is optional, if left empty, names can collide if the list contains more than one with the same name
func (list ChangeSetList) Find(namespace, name string) (*ChangeSet, error) {
	for _, changeSet := range list {
		if changeSet.Metadata.Name == name {
			if namespace == "" || changeSet.Metadata.Namespace == namespace {
				return changeSet, nil
			}
		}
	}
	return nil, errors.Errorf("list did not find changeSet %v.%v", namespace, name)
}

func (list ChangeSetList) AsResources() resources.ResourceList {
	var ress resources.ResourceList
	for _, changeSet := range list {
		ress = append(ress, changeSet)
	}
	return ress
}

func (list ChangeSetList) AsInputResources() resources.InputResourceList {
	var ress resources.InputResourceList
	for _, changeSet := range list {
		ress = append(ress, changeSet)
	}
	return ress
}

func (list ChangeSetList) Names() []string {
	var names []string
	for _, changeSet := range list {
		names = append(names, changeSet.Metadata.Name)
	}
	return names
}

func (list ChangeSetList) NamespacesDotNames() []string {
	var names []string
	for _, changeSet := range list {
		names = append(names, changeSet.Metadata.Namespace+"."+changeSet.Metadata.Name)
	}
	return names
}

func (list ChangeSetList) Sort() ChangeSetList {
	sort.SliceStable(list, func(i, j int) bool {
		return list[i].Metadata.Less(list[j].Metadata)
	})
	return list
}

func (list ChangeSetList) Clone() ChangeSetList {
	var changeSetList ChangeSetList
	for _, changeSet := range list {
		changeSetList = append(changeSetList, proto.Clone(changeSet).(*ChangeSet))
	}
	return changeSetList
}

func (list ChangeSetList) Each(f func(element *ChangeSet)) {
	for _, changeSet := range list {
		f(changeSet)
	}
}

func (list ChangeSetList) AsInterfaces() []interface{} {
	var asInterfaces []interface{}
	list.Each(func(element *ChangeSet) {
		asInterfaces = append(asInterfaces, element)
	})
	return asInterfaces
}

func (list ChangeSetList) ByNamespace() ChangesetsByNamespace {
	byNamespace := make(ChangesetsByNamespace)
	for _, changeSet := range list {
		byNamespace.Add(changeSet)
	}
	return byNamespace
}

func (byNamespace ChangesetsByNamespace) Add(changeSet ...*ChangeSet) {
	for _, item := range changeSet {
		byNamespace[item.Metadata.Namespace] = append(byNamespace[item.Metadata.Namespace], item)
	}
}

func (byNamespace ChangesetsByNamespace) Clear(namespace string) {
	delete(byNamespace, namespace)
}

func (byNamespace ChangesetsByNamespace) List() ChangeSetList {
	var list ChangeSetList
	for _, changeSetList := range byNamespace {
		list = append(list, changeSetList...)
	}
	return list.Sort()
}

func (byNamespace ChangesetsByNamespace) Clone() ChangesetsByNamespace {
	return byNamespace.List().Clone().ByNamespace()
}

var _ resources.Resource = &ChangeSet{}

// Kubernetes Adapter for ChangeSet

func (o *ChangeSet) GetObjectKind() schema.ObjectKind {
	t := ChangeSetCrd.TypeMeta()
	return &t
}

func (o *ChangeSet) DeepCopyObject() runtime.Object {
	return resources.Clone(o).(*ChangeSet)
}

var ChangeSetCrd = crd.NewCrd("vcs.solo.io",
	"changesets",
	"vcs.solo.io",
	"v1",
	"ChangeSet",
	"chg",
	&ChangeSet{})
