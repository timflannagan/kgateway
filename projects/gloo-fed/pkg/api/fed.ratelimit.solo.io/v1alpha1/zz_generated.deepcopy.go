// Code generated by skv2. DO NOT EDIT.

// This file contains generated Deepcopy methods for fed.ratelimit.solo.io/v1alpha1 resources

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// Generated Deepcopy methods for FederatedRateLimitConfig

func (in *FederatedRateLimitConfig) DeepCopyInto(out *FederatedRateLimitConfig) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)

	// deepcopy spec
	in.Spec.DeepCopyInto(&out.Spec)
	// deepcopy status
	in.Status.DeepCopyInto(&out.Status)

	return
}

func (in *FederatedRateLimitConfig) DeepCopy() *FederatedRateLimitConfig {
	if in == nil {
		return nil
	}
	out := new(FederatedRateLimitConfig)
	in.DeepCopyInto(out)
	return out
}

func (in *FederatedRateLimitConfig) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (in *FederatedRateLimitConfigList) DeepCopyInto(out *FederatedRateLimitConfigList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]FederatedRateLimitConfig, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

func (in *FederatedRateLimitConfigList) DeepCopy() *FederatedRateLimitConfigList {
	if in == nil {
		return nil
	}
	out := new(FederatedRateLimitConfigList)
	in.DeepCopyInto(out)
	return out
}

func (in *FederatedRateLimitConfigList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}
