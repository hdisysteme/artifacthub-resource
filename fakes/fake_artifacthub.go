// Code generated by counterfeiter. DO NOT EDIT.
package fakes

import (
	"github.com/hdisysteme/artifacthub-resource/internal/pkg/resource"
	"sync"
)

type FakeArtifactHub struct {
	ListHelmVersionStub        func(resource.Package, string) (*resource.HelmVersion, error)
	listHelmVersionMutex       sync.RWMutex
	listHelmVersionArgsForCall []struct {
		arg1 resource.Package
		arg2 string
	}
	listHelmVersionReturns struct {
		result1 *resource.HelmVersion
		result2 error
	}
	listHelmVersionReturnsOnCall map[int]struct {
		result1 *resource.HelmVersion
		result2 error
	}
	ListHelmVersionsStub        func(resource.Package) ([]resource.Version, error)
	listHelmVersionsMutex       sync.RWMutex
	listHelmVersionsArgsForCall []struct {
		arg1 resource.Package
	}
	listHelmVersionsReturns struct {
		result1 []resource.Version
		result2 error
	}
	listHelmVersionsReturnsOnCall map[int]struct {
		result1 []resource.Version
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeArtifactHub) ListHelmVersion(arg1 resource.Package, arg2 string) (*resource.HelmVersion, error) {
	fake.listHelmVersionMutex.Lock()
	ret, specificReturn := fake.listHelmVersionReturnsOnCall[len(fake.listHelmVersionArgsForCall)]
	fake.listHelmVersionArgsForCall = append(fake.listHelmVersionArgsForCall, struct {
		arg1 resource.Package
		arg2 string
	}{arg1, arg2})
	stub := fake.ListHelmVersionStub
	fakeReturns := fake.listHelmVersionReturns
	fake.recordInvocation("ListHelmVersion", []interface{}{arg1, arg2})
	fake.listHelmVersionMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeArtifactHub) ListHelmVersionCallCount() int {
	fake.listHelmVersionMutex.RLock()
	defer fake.listHelmVersionMutex.RUnlock()
	return len(fake.listHelmVersionArgsForCall)
}

func (fake *FakeArtifactHub) ListHelmVersionCalls(stub func(resource.Package, string) (*resource.HelmVersion, error)) {
	fake.listHelmVersionMutex.Lock()
	defer fake.listHelmVersionMutex.Unlock()
	fake.ListHelmVersionStub = stub
}

func (fake *FakeArtifactHub) ListHelmVersionArgsForCall(i int) (resource.Package, string) {
	fake.listHelmVersionMutex.RLock()
	defer fake.listHelmVersionMutex.RUnlock()
	argsForCall := fake.listHelmVersionArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeArtifactHub) ListHelmVersionReturns(result1 *resource.HelmVersion, result2 error) {
	fake.listHelmVersionMutex.Lock()
	defer fake.listHelmVersionMutex.Unlock()
	fake.ListHelmVersionStub = nil
	fake.listHelmVersionReturns = struct {
		result1 *resource.HelmVersion
		result2 error
	}{result1, result2}
}

func (fake *FakeArtifactHub) ListHelmVersionReturnsOnCall(i int, result1 *resource.HelmVersion, result2 error) {
	fake.listHelmVersionMutex.Lock()
	defer fake.listHelmVersionMutex.Unlock()
	fake.ListHelmVersionStub = nil
	if fake.listHelmVersionReturnsOnCall == nil {
		fake.listHelmVersionReturnsOnCall = make(map[int]struct {
			result1 *resource.HelmVersion
			result2 error
		})
	}
	fake.listHelmVersionReturnsOnCall[i] = struct {
		result1 *resource.HelmVersion
		result2 error
	}{result1, result2}
}

func (fake *FakeArtifactHub) ListHelmVersions(arg1 resource.Package) ([]resource.Version, error) {
	fake.listHelmVersionsMutex.Lock()
	ret, specificReturn := fake.listHelmVersionsReturnsOnCall[len(fake.listHelmVersionsArgsForCall)]
	fake.listHelmVersionsArgsForCall = append(fake.listHelmVersionsArgsForCall, struct {
		arg1 resource.Package
	}{arg1})
	stub := fake.ListHelmVersionsStub
	fakeReturns := fake.listHelmVersionsReturns
	fake.recordInvocation("ListHelmVersions", []interface{}{arg1})
	fake.listHelmVersionsMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeArtifactHub) ListHelmVersionsCallCount() int {
	fake.listHelmVersionsMutex.RLock()
	defer fake.listHelmVersionsMutex.RUnlock()
	return len(fake.listHelmVersionsArgsForCall)
}

func (fake *FakeArtifactHub) ListHelmVersionsCalls(stub func(resource.Package) ([]resource.Version, error)) {
	fake.listHelmVersionsMutex.Lock()
	defer fake.listHelmVersionsMutex.Unlock()
	fake.ListHelmVersionsStub = stub
}

func (fake *FakeArtifactHub) ListHelmVersionsArgsForCall(i int) resource.Package {
	fake.listHelmVersionsMutex.RLock()
	defer fake.listHelmVersionsMutex.RUnlock()
	argsForCall := fake.listHelmVersionsArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeArtifactHub) ListHelmVersionsReturns(result1 []resource.Version, result2 error) {
	fake.listHelmVersionsMutex.Lock()
	defer fake.listHelmVersionsMutex.Unlock()
	fake.ListHelmVersionsStub = nil
	fake.listHelmVersionsReturns = struct {
		result1 []resource.Version
		result2 error
	}{result1, result2}
}

func (fake *FakeArtifactHub) ListHelmVersionsReturnsOnCall(i int, result1 []resource.Version, result2 error) {
	fake.listHelmVersionsMutex.Lock()
	defer fake.listHelmVersionsMutex.Unlock()
	fake.ListHelmVersionsStub = nil
	if fake.listHelmVersionsReturnsOnCall == nil {
		fake.listHelmVersionsReturnsOnCall = make(map[int]struct {
			result1 []resource.Version
			result2 error
		})
	}
	fake.listHelmVersionsReturnsOnCall[i] = struct {
		result1 []resource.Version
		result2 error
	}{result1, result2}
}

func (fake *FakeArtifactHub) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.listHelmVersionMutex.RLock()
	defer fake.listHelmVersionMutex.RUnlock()
	fake.listHelmVersionsMutex.RLock()
	defer fake.listHelmVersionsMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeArtifactHub) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ resource.ArtifactHub = new(FakeArtifactHub)
