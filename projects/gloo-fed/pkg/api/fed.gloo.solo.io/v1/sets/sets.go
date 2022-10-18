// Code generated by skv2. DO NOT EDIT.

//go:generate mockgen -source ./sets.go -destination mocks/sets.go

package v1sets

import (
	fed_gloo_solo_io_v1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.gloo.solo.io/v1"

	"github.com/rotisserie/eris"
	sksets "github.com/solo-io/skv2/contrib/pkg/sets"
	"github.com/solo-io/skv2/pkg/ezkube"
	"k8s.io/apimachinery/pkg/util/sets"
)

type FederatedUpstreamSet interface {
	// Get the set stored keys
	Keys() sets.String
	// List of resources stored in the set. Pass an optional filter function to filter on the list.
	List(filterResource ...func(*fed_gloo_solo_io_v1.FederatedUpstream) bool) []*fed_gloo_solo_io_v1.FederatedUpstream
	// Unsorted list of resources stored in the set. Pass an optional filter function to filter on the list.
	UnsortedList(filterResource ...func(*fed_gloo_solo_io_v1.FederatedUpstream) bool) []*fed_gloo_solo_io_v1.FederatedUpstream
	// Return the Set as a map of key to resource.
	Map() map[string]*fed_gloo_solo_io_v1.FederatedUpstream
	// Insert a resource into the set.
	Insert(federatedUpstream ...*fed_gloo_solo_io_v1.FederatedUpstream)
	// Compare the equality of the keys in two sets (not the resources themselves)
	Equal(federatedUpstreamSet FederatedUpstreamSet) bool
	// Check if the set contains a key matching the resource (not the resource itself)
	Has(federatedUpstream ezkube.ResourceId) bool
	// Delete the key matching the resource
	Delete(federatedUpstream ezkube.ResourceId)
	// Return the union with the provided set
	Union(set FederatedUpstreamSet) FederatedUpstreamSet
	// Return the difference with the provided set
	Difference(set FederatedUpstreamSet) FederatedUpstreamSet
	// Return the intersection with the provided set
	Intersection(set FederatedUpstreamSet) FederatedUpstreamSet
	// Find the resource with the given ID
	Find(id ezkube.ResourceId) (*fed_gloo_solo_io_v1.FederatedUpstream, error)
	// Get the length of the set
	Length() int
	// returns the generic implementation of the set
	Generic() sksets.ResourceSet
	// returns the delta between this and and another FederatedUpstreamSet
	Delta(newSet FederatedUpstreamSet) sksets.ResourceDelta
	// Create a deep copy of the current FederatedUpstreamSet
	Clone() FederatedUpstreamSet
}

func makeGenericFederatedUpstreamSet(federatedUpstreamList []*fed_gloo_solo_io_v1.FederatedUpstream) sksets.ResourceSet {
	var genericResources []ezkube.ResourceId
	for _, obj := range federatedUpstreamList {
		genericResources = append(genericResources, obj)
	}
	return sksets.NewResourceSet(genericResources...)
}

type federatedUpstreamSet struct {
	set sksets.ResourceSet
}

func NewFederatedUpstreamSet(federatedUpstreamList ...*fed_gloo_solo_io_v1.FederatedUpstream) FederatedUpstreamSet {
	return &federatedUpstreamSet{set: makeGenericFederatedUpstreamSet(federatedUpstreamList)}
}

func NewFederatedUpstreamSetFromList(federatedUpstreamList *fed_gloo_solo_io_v1.FederatedUpstreamList) FederatedUpstreamSet {
	list := make([]*fed_gloo_solo_io_v1.FederatedUpstream, 0, len(federatedUpstreamList.Items))
	for idx := range federatedUpstreamList.Items {
		list = append(list, &federatedUpstreamList.Items[idx])
	}
	return &federatedUpstreamSet{set: makeGenericFederatedUpstreamSet(list)}
}

func (s *federatedUpstreamSet) Keys() sets.String {
	if s == nil {
		return sets.String{}
	}
	return s.Generic().Keys()
}

func (s *federatedUpstreamSet) List(filterResource ...func(*fed_gloo_solo_io_v1.FederatedUpstream) bool) []*fed_gloo_solo_io_v1.FederatedUpstream {
	if s == nil {
		return nil
	}
	var genericFilters []func(ezkube.ResourceId) bool
	for _, filter := range filterResource {
		filter := filter
		genericFilters = append(genericFilters, func(obj ezkube.ResourceId) bool {
			return filter(obj.(*fed_gloo_solo_io_v1.FederatedUpstream))
		})
	}

	objs := s.Generic().List(genericFilters...)
	federatedUpstreamList := make([]*fed_gloo_solo_io_v1.FederatedUpstream, 0, len(objs))
	for _, obj := range objs {
		federatedUpstreamList = append(federatedUpstreamList, obj.(*fed_gloo_solo_io_v1.FederatedUpstream))
	}
	return federatedUpstreamList
}

func (s *federatedUpstreamSet) UnsortedList(filterResource ...func(*fed_gloo_solo_io_v1.FederatedUpstream) bool) []*fed_gloo_solo_io_v1.FederatedUpstream {
	if s == nil {
		return nil
	}
	var genericFilters []func(ezkube.ResourceId) bool
	for _, filter := range filterResource {
		filter := filter
		genericFilters = append(genericFilters, func(obj ezkube.ResourceId) bool {
			return filter(obj.(*fed_gloo_solo_io_v1.FederatedUpstream))
		})
	}

	var federatedUpstreamList []*fed_gloo_solo_io_v1.FederatedUpstream
	for _, obj := range s.Generic().UnsortedList(genericFilters...) {
		federatedUpstreamList = append(federatedUpstreamList, obj.(*fed_gloo_solo_io_v1.FederatedUpstream))
	}
	return federatedUpstreamList
}

func (s *federatedUpstreamSet) Map() map[string]*fed_gloo_solo_io_v1.FederatedUpstream {
	if s == nil {
		return nil
	}

	newMap := map[string]*fed_gloo_solo_io_v1.FederatedUpstream{}
	for k, v := range s.Generic().Map() {
		newMap[k] = v.(*fed_gloo_solo_io_v1.FederatedUpstream)
	}
	return newMap
}

func (s *federatedUpstreamSet) Insert(
	federatedUpstreamList ...*fed_gloo_solo_io_v1.FederatedUpstream,
) {
	if s == nil {
		panic("cannot insert into nil set")
	}

	for _, obj := range federatedUpstreamList {
		s.Generic().Insert(obj)
	}
}

func (s *federatedUpstreamSet) Has(federatedUpstream ezkube.ResourceId) bool {
	if s == nil {
		return false
	}
	return s.Generic().Has(federatedUpstream)
}

func (s *federatedUpstreamSet) Equal(
	federatedUpstreamSet FederatedUpstreamSet,
) bool {
	if s == nil {
		return federatedUpstreamSet == nil
	}
	return s.Generic().Equal(federatedUpstreamSet.Generic())
}

func (s *federatedUpstreamSet) Delete(FederatedUpstream ezkube.ResourceId) {
	if s == nil {
		return
	}
	s.Generic().Delete(FederatedUpstream)
}

func (s *federatedUpstreamSet) Union(set FederatedUpstreamSet) FederatedUpstreamSet {
	if s == nil {
		return set
	}
	return NewFederatedUpstreamSet(append(s.List(), set.List()...)...)
}

func (s *federatedUpstreamSet) Difference(set FederatedUpstreamSet) FederatedUpstreamSet {
	if s == nil {
		return set
	}
	newSet := s.Generic().Difference(set.Generic())
	return &federatedUpstreamSet{set: newSet}
}

func (s *federatedUpstreamSet) Intersection(set FederatedUpstreamSet) FederatedUpstreamSet {
	if s == nil {
		return nil
	}
	newSet := s.Generic().Intersection(set.Generic())
	var federatedUpstreamList []*fed_gloo_solo_io_v1.FederatedUpstream
	for _, obj := range newSet.List() {
		federatedUpstreamList = append(federatedUpstreamList, obj.(*fed_gloo_solo_io_v1.FederatedUpstream))
	}
	return NewFederatedUpstreamSet(federatedUpstreamList...)
}

func (s *federatedUpstreamSet) Find(id ezkube.ResourceId) (*fed_gloo_solo_io_v1.FederatedUpstream, error) {
	if s == nil {
		return nil, eris.Errorf("empty set, cannot find FederatedUpstream %v", sksets.Key(id))
	}
	obj, err := s.Generic().Find(&fed_gloo_solo_io_v1.FederatedUpstream{}, id)
	if err != nil {
		return nil, err
	}

	return obj.(*fed_gloo_solo_io_v1.FederatedUpstream), nil
}

func (s *federatedUpstreamSet) Length() int {
	if s == nil {
		return 0
	}
	return s.Generic().Length()
}

func (s *federatedUpstreamSet) Generic() sksets.ResourceSet {
	if s == nil {
		return nil
	}
	return s.set
}

func (s *federatedUpstreamSet) Delta(newSet FederatedUpstreamSet) sksets.ResourceDelta {
	if s == nil {
		return sksets.ResourceDelta{
			Inserted: newSet.Generic(),
		}
	}
	return s.Generic().Delta(newSet.Generic())
}

func (s *federatedUpstreamSet) Clone() FederatedUpstreamSet {
	if s == nil {
		return nil
	}
	return &federatedUpstreamSet{set: sksets.NewResourceSet(s.Generic().Clone().List()...)}
}

type FederatedUpstreamGroupSet interface {
	// Get the set stored keys
	Keys() sets.String
	// List of resources stored in the set. Pass an optional filter function to filter on the list.
	List(filterResource ...func(*fed_gloo_solo_io_v1.FederatedUpstreamGroup) bool) []*fed_gloo_solo_io_v1.FederatedUpstreamGroup
	// Unsorted list of resources stored in the set. Pass an optional filter function to filter on the list.
	UnsortedList(filterResource ...func(*fed_gloo_solo_io_v1.FederatedUpstreamGroup) bool) []*fed_gloo_solo_io_v1.FederatedUpstreamGroup
	// Return the Set as a map of key to resource.
	Map() map[string]*fed_gloo_solo_io_v1.FederatedUpstreamGroup
	// Insert a resource into the set.
	Insert(federatedUpstreamGroup ...*fed_gloo_solo_io_v1.FederatedUpstreamGroup)
	// Compare the equality of the keys in two sets (not the resources themselves)
	Equal(federatedUpstreamGroupSet FederatedUpstreamGroupSet) bool
	// Check if the set contains a key matching the resource (not the resource itself)
	Has(federatedUpstreamGroup ezkube.ResourceId) bool
	// Delete the key matching the resource
	Delete(federatedUpstreamGroup ezkube.ResourceId)
	// Return the union with the provided set
	Union(set FederatedUpstreamGroupSet) FederatedUpstreamGroupSet
	// Return the difference with the provided set
	Difference(set FederatedUpstreamGroupSet) FederatedUpstreamGroupSet
	// Return the intersection with the provided set
	Intersection(set FederatedUpstreamGroupSet) FederatedUpstreamGroupSet
	// Find the resource with the given ID
	Find(id ezkube.ResourceId) (*fed_gloo_solo_io_v1.FederatedUpstreamGroup, error)
	// Get the length of the set
	Length() int
	// returns the generic implementation of the set
	Generic() sksets.ResourceSet
	// returns the delta between this and and another FederatedUpstreamGroupSet
	Delta(newSet FederatedUpstreamGroupSet) sksets.ResourceDelta
	// Create a deep copy of the current FederatedUpstreamGroupSet
	Clone() FederatedUpstreamGroupSet
}

func makeGenericFederatedUpstreamGroupSet(federatedUpstreamGroupList []*fed_gloo_solo_io_v1.FederatedUpstreamGroup) sksets.ResourceSet {
	var genericResources []ezkube.ResourceId
	for _, obj := range federatedUpstreamGroupList {
		genericResources = append(genericResources, obj)
	}
	return sksets.NewResourceSet(genericResources...)
}

type federatedUpstreamGroupSet struct {
	set sksets.ResourceSet
}

func NewFederatedUpstreamGroupSet(federatedUpstreamGroupList ...*fed_gloo_solo_io_v1.FederatedUpstreamGroup) FederatedUpstreamGroupSet {
	return &federatedUpstreamGroupSet{set: makeGenericFederatedUpstreamGroupSet(federatedUpstreamGroupList)}
}

func NewFederatedUpstreamGroupSetFromList(federatedUpstreamGroupList *fed_gloo_solo_io_v1.FederatedUpstreamGroupList) FederatedUpstreamGroupSet {
	list := make([]*fed_gloo_solo_io_v1.FederatedUpstreamGroup, 0, len(federatedUpstreamGroupList.Items))
	for idx := range federatedUpstreamGroupList.Items {
		list = append(list, &federatedUpstreamGroupList.Items[idx])
	}
	return &federatedUpstreamGroupSet{set: makeGenericFederatedUpstreamGroupSet(list)}
}

func (s *federatedUpstreamGroupSet) Keys() sets.String {
	if s == nil {
		return sets.String{}
	}
	return s.Generic().Keys()
}

func (s *federatedUpstreamGroupSet) List(filterResource ...func(*fed_gloo_solo_io_v1.FederatedUpstreamGroup) bool) []*fed_gloo_solo_io_v1.FederatedUpstreamGroup {
	if s == nil {
		return nil
	}
	var genericFilters []func(ezkube.ResourceId) bool
	for _, filter := range filterResource {
		filter := filter
		genericFilters = append(genericFilters, func(obj ezkube.ResourceId) bool {
			return filter(obj.(*fed_gloo_solo_io_v1.FederatedUpstreamGroup))
		})
	}

	objs := s.Generic().List(genericFilters...)
	federatedUpstreamGroupList := make([]*fed_gloo_solo_io_v1.FederatedUpstreamGroup, 0, len(objs))
	for _, obj := range objs {
		federatedUpstreamGroupList = append(federatedUpstreamGroupList, obj.(*fed_gloo_solo_io_v1.FederatedUpstreamGroup))
	}
	return federatedUpstreamGroupList
}

func (s *federatedUpstreamGroupSet) UnsortedList(filterResource ...func(*fed_gloo_solo_io_v1.FederatedUpstreamGroup) bool) []*fed_gloo_solo_io_v1.FederatedUpstreamGroup {
	if s == nil {
		return nil
	}
	var genericFilters []func(ezkube.ResourceId) bool
	for _, filter := range filterResource {
		filter := filter
		genericFilters = append(genericFilters, func(obj ezkube.ResourceId) bool {
			return filter(obj.(*fed_gloo_solo_io_v1.FederatedUpstreamGroup))
		})
	}

	var federatedUpstreamGroupList []*fed_gloo_solo_io_v1.FederatedUpstreamGroup
	for _, obj := range s.Generic().UnsortedList(genericFilters...) {
		federatedUpstreamGroupList = append(federatedUpstreamGroupList, obj.(*fed_gloo_solo_io_v1.FederatedUpstreamGroup))
	}
	return federatedUpstreamGroupList
}

func (s *federatedUpstreamGroupSet) Map() map[string]*fed_gloo_solo_io_v1.FederatedUpstreamGroup {
	if s == nil {
		return nil
	}

	newMap := map[string]*fed_gloo_solo_io_v1.FederatedUpstreamGroup{}
	for k, v := range s.Generic().Map() {
		newMap[k] = v.(*fed_gloo_solo_io_v1.FederatedUpstreamGroup)
	}
	return newMap
}

func (s *federatedUpstreamGroupSet) Insert(
	federatedUpstreamGroupList ...*fed_gloo_solo_io_v1.FederatedUpstreamGroup,
) {
	if s == nil {
		panic("cannot insert into nil set")
	}

	for _, obj := range federatedUpstreamGroupList {
		s.Generic().Insert(obj)
	}
}

func (s *federatedUpstreamGroupSet) Has(federatedUpstreamGroup ezkube.ResourceId) bool {
	if s == nil {
		return false
	}
	return s.Generic().Has(federatedUpstreamGroup)
}

func (s *federatedUpstreamGroupSet) Equal(
	federatedUpstreamGroupSet FederatedUpstreamGroupSet,
) bool {
	if s == nil {
		return federatedUpstreamGroupSet == nil
	}
	return s.Generic().Equal(federatedUpstreamGroupSet.Generic())
}

func (s *federatedUpstreamGroupSet) Delete(FederatedUpstreamGroup ezkube.ResourceId) {
	if s == nil {
		return
	}
	s.Generic().Delete(FederatedUpstreamGroup)
}

func (s *federatedUpstreamGroupSet) Union(set FederatedUpstreamGroupSet) FederatedUpstreamGroupSet {
	if s == nil {
		return set
	}
	return NewFederatedUpstreamGroupSet(append(s.List(), set.List()...)...)
}

func (s *federatedUpstreamGroupSet) Difference(set FederatedUpstreamGroupSet) FederatedUpstreamGroupSet {
	if s == nil {
		return set
	}
	newSet := s.Generic().Difference(set.Generic())
	return &federatedUpstreamGroupSet{set: newSet}
}

func (s *federatedUpstreamGroupSet) Intersection(set FederatedUpstreamGroupSet) FederatedUpstreamGroupSet {
	if s == nil {
		return nil
	}
	newSet := s.Generic().Intersection(set.Generic())
	var federatedUpstreamGroupList []*fed_gloo_solo_io_v1.FederatedUpstreamGroup
	for _, obj := range newSet.List() {
		federatedUpstreamGroupList = append(federatedUpstreamGroupList, obj.(*fed_gloo_solo_io_v1.FederatedUpstreamGroup))
	}
	return NewFederatedUpstreamGroupSet(federatedUpstreamGroupList...)
}

func (s *federatedUpstreamGroupSet) Find(id ezkube.ResourceId) (*fed_gloo_solo_io_v1.FederatedUpstreamGroup, error) {
	if s == nil {
		return nil, eris.Errorf("empty set, cannot find FederatedUpstreamGroup %v", sksets.Key(id))
	}
	obj, err := s.Generic().Find(&fed_gloo_solo_io_v1.FederatedUpstreamGroup{}, id)
	if err != nil {
		return nil, err
	}

	return obj.(*fed_gloo_solo_io_v1.FederatedUpstreamGroup), nil
}

func (s *federatedUpstreamGroupSet) Length() int {
	if s == nil {
		return 0
	}
	return s.Generic().Length()
}

func (s *federatedUpstreamGroupSet) Generic() sksets.ResourceSet {
	if s == nil {
		return nil
	}
	return s.set
}

func (s *federatedUpstreamGroupSet) Delta(newSet FederatedUpstreamGroupSet) sksets.ResourceDelta {
	if s == nil {
		return sksets.ResourceDelta{
			Inserted: newSet.Generic(),
		}
	}
	return s.Generic().Delta(newSet.Generic())
}

func (s *federatedUpstreamGroupSet) Clone() FederatedUpstreamGroupSet {
	if s == nil {
		return nil
	}
	return &federatedUpstreamGroupSet{set: sksets.NewResourceSet(s.Generic().Clone().List()...)}
}

type FederatedSettingsSet interface {
	// Get the set stored keys
	Keys() sets.String
	// List of resources stored in the set. Pass an optional filter function to filter on the list.
	List(filterResource ...func(*fed_gloo_solo_io_v1.FederatedSettings) bool) []*fed_gloo_solo_io_v1.FederatedSettings
	// Unsorted list of resources stored in the set. Pass an optional filter function to filter on the list.
	UnsortedList(filterResource ...func(*fed_gloo_solo_io_v1.FederatedSettings) bool) []*fed_gloo_solo_io_v1.FederatedSettings
	// Return the Set as a map of key to resource.
	Map() map[string]*fed_gloo_solo_io_v1.FederatedSettings
	// Insert a resource into the set.
	Insert(federatedSettings ...*fed_gloo_solo_io_v1.FederatedSettings)
	// Compare the equality of the keys in two sets (not the resources themselves)
	Equal(federatedSettingsSet FederatedSettingsSet) bool
	// Check if the set contains a key matching the resource (not the resource itself)
	Has(federatedSettings ezkube.ResourceId) bool
	// Delete the key matching the resource
	Delete(federatedSettings ezkube.ResourceId)
	// Return the union with the provided set
	Union(set FederatedSettingsSet) FederatedSettingsSet
	// Return the difference with the provided set
	Difference(set FederatedSettingsSet) FederatedSettingsSet
	// Return the intersection with the provided set
	Intersection(set FederatedSettingsSet) FederatedSettingsSet
	// Find the resource with the given ID
	Find(id ezkube.ResourceId) (*fed_gloo_solo_io_v1.FederatedSettings, error)
	// Get the length of the set
	Length() int
	// returns the generic implementation of the set
	Generic() sksets.ResourceSet
	// returns the delta between this and and another FederatedSettingsSet
	Delta(newSet FederatedSettingsSet) sksets.ResourceDelta
	// Create a deep copy of the current FederatedSettingsSet
	Clone() FederatedSettingsSet
}

func makeGenericFederatedSettingsSet(federatedSettingsList []*fed_gloo_solo_io_v1.FederatedSettings) sksets.ResourceSet {
	var genericResources []ezkube.ResourceId
	for _, obj := range federatedSettingsList {
		genericResources = append(genericResources, obj)
	}
	return sksets.NewResourceSet(genericResources...)
}

type federatedSettingsSet struct {
	set sksets.ResourceSet
}

func NewFederatedSettingsSet(federatedSettingsList ...*fed_gloo_solo_io_v1.FederatedSettings) FederatedSettingsSet {
	return &federatedSettingsSet{set: makeGenericFederatedSettingsSet(federatedSettingsList)}
}

func NewFederatedSettingsSetFromList(federatedSettingsList *fed_gloo_solo_io_v1.FederatedSettingsList) FederatedSettingsSet {
	list := make([]*fed_gloo_solo_io_v1.FederatedSettings, 0, len(federatedSettingsList.Items))
	for idx := range federatedSettingsList.Items {
		list = append(list, &federatedSettingsList.Items[idx])
	}
	return &federatedSettingsSet{set: makeGenericFederatedSettingsSet(list)}
}

func (s *federatedSettingsSet) Keys() sets.String {
	if s == nil {
		return sets.String{}
	}
	return s.Generic().Keys()
}

func (s *federatedSettingsSet) List(filterResource ...func(*fed_gloo_solo_io_v1.FederatedSettings) bool) []*fed_gloo_solo_io_v1.FederatedSettings {
	if s == nil {
		return nil
	}
	var genericFilters []func(ezkube.ResourceId) bool
	for _, filter := range filterResource {
		filter := filter
		genericFilters = append(genericFilters, func(obj ezkube.ResourceId) bool {
			return filter(obj.(*fed_gloo_solo_io_v1.FederatedSettings))
		})
	}

	objs := s.Generic().List(genericFilters...)
	federatedSettingsList := make([]*fed_gloo_solo_io_v1.FederatedSettings, 0, len(objs))
	for _, obj := range objs {
		federatedSettingsList = append(federatedSettingsList, obj.(*fed_gloo_solo_io_v1.FederatedSettings))
	}
	return federatedSettingsList
}

func (s *federatedSettingsSet) UnsortedList(filterResource ...func(*fed_gloo_solo_io_v1.FederatedSettings) bool) []*fed_gloo_solo_io_v1.FederatedSettings {
	if s == nil {
		return nil
	}
	var genericFilters []func(ezkube.ResourceId) bool
	for _, filter := range filterResource {
		filter := filter
		genericFilters = append(genericFilters, func(obj ezkube.ResourceId) bool {
			return filter(obj.(*fed_gloo_solo_io_v1.FederatedSettings))
		})
	}

	var federatedSettingsList []*fed_gloo_solo_io_v1.FederatedSettings
	for _, obj := range s.Generic().UnsortedList(genericFilters...) {
		federatedSettingsList = append(federatedSettingsList, obj.(*fed_gloo_solo_io_v1.FederatedSettings))
	}
	return federatedSettingsList
}

func (s *federatedSettingsSet) Map() map[string]*fed_gloo_solo_io_v1.FederatedSettings {
	if s == nil {
		return nil
	}

	newMap := map[string]*fed_gloo_solo_io_v1.FederatedSettings{}
	for k, v := range s.Generic().Map() {
		newMap[k] = v.(*fed_gloo_solo_io_v1.FederatedSettings)
	}
	return newMap
}

func (s *federatedSettingsSet) Insert(
	federatedSettingsList ...*fed_gloo_solo_io_v1.FederatedSettings,
) {
	if s == nil {
		panic("cannot insert into nil set")
	}

	for _, obj := range federatedSettingsList {
		s.Generic().Insert(obj)
	}
}

func (s *federatedSettingsSet) Has(federatedSettings ezkube.ResourceId) bool {
	if s == nil {
		return false
	}
	return s.Generic().Has(federatedSettings)
}

func (s *federatedSettingsSet) Equal(
	federatedSettingsSet FederatedSettingsSet,
) bool {
	if s == nil {
		return federatedSettingsSet == nil
	}
	return s.Generic().Equal(federatedSettingsSet.Generic())
}

func (s *federatedSettingsSet) Delete(FederatedSettings ezkube.ResourceId) {
	if s == nil {
		return
	}
	s.Generic().Delete(FederatedSettings)
}

func (s *federatedSettingsSet) Union(set FederatedSettingsSet) FederatedSettingsSet {
	if s == nil {
		return set
	}
	return NewFederatedSettingsSet(append(s.List(), set.List()...)...)
}

func (s *federatedSettingsSet) Difference(set FederatedSettingsSet) FederatedSettingsSet {
	if s == nil {
		return set
	}
	newSet := s.Generic().Difference(set.Generic())
	return &federatedSettingsSet{set: newSet}
}

func (s *federatedSettingsSet) Intersection(set FederatedSettingsSet) FederatedSettingsSet {
	if s == nil {
		return nil
	}
	newSet := s.Generic().Intersection(set.Generic())
	var federatedSettingsList []*fed_gloo_solo_io_v1.FederatedSettings
	for _, obj := range newSet.List() {
		federatedSettingsList = append(federatedSettingsList, obj.(*fed_gloo_solo_io_v1.FederatedSettings))
	}
	return NewFederatedSettingsSet(federatedSettingsList...)
}

func (s *federatedSettingsSet) Find(id ezkube.ResourceId) (*fed_gloo_solo_io_v1.FederatedSettings, error) {
	if s == nil {
		return nil, eris.Errorf("empty set, cannot find FederatedSettings %v", sksets.Key(id))
	}
	obj, err := s.Generic().Find(&fed_gloo_solo_io_v1.FederatedSettings{}, id)
	if err != nil {
		return nil, err
	}

	return obj.(*fed_gloo_solo_io_v1.FederatedSettings), nil
}

func (s *federatedSettingsSet) Length() int {
	if s == nil {
		return 0
	}
	return s.Generic().Length()
}

func (s *federatedSettingsSet) Generic() sksets.ResourceSet {
	if s == nil {
		return nil
	}
	return s.set
}

func (s *federatedSettingsSet) Delta(newSet FederatedSettingsSet) sksets.ResourceDelta {
	if s == nil {
		return sksets.ResourceDelta{
			Inserted: newSet.Generic(),
		}
	}
	return s.Generic().Delta(newSet.Generic())
}

func (s *federatedSettingsSet) Clone() FederatedSettingsSet {
	if s == nil {
		return nil
	}
	return &federatedSettingsSet{set: sksets.NewResourceSet(s.Generic().Clone().List()...)}
}
