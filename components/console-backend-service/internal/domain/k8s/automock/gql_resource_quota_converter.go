// Code generated by mockery v1.0.0. DO NOT EDIT.

package automock

import gqlschema "github.com/kyma-project/kyma/components/console-backend-service/internal/gqlschema"

import mock "github.com/stretchr/testify/mock"
import v1 "k8s.io/api/core/v1"

// gqlResourceQuotaConverter is an autogenerated mock type for the gqlResourceQuotaConverter type
type gqlResourceQuotaConverter struct {
	mock.Mock
}

// ToGQL provides a mock function with given fields: in
func (_m *gqlResourceQuotaConverter) ToGQL(in *v1.ResourceQuota) *gqlschema.ResourceQuota {
	ret := _m.Called(in)

	var r0 *gqlschema.ResourceQuota
	if rf, ok := ret.Get(0).(func(*v1.ResourceQuota) *gqlschema.ResourceQuota); ok {
		r0 = rf(in)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gqlschema.ResourceQuota)
		}
	}

	return r0
}

// ToGQLs provides a mock function with given fields: in
func (_m *gqlResourceQuotaConverter) ToGQLs(in []*v1.ResourceQuota) []gqlschema.ResourceQuota {
	ret := _m.Called(in)

	var r0 []gqlschema.ResourceQuota
	if rf, ok := ret.Get(0).(func([]*v1.ResourceQuota) []gqlschema.ResourceQuota); ok {
		r0 = rf(in)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]gqlschema.ResourceQuota)
		}
	}

	return r0
}