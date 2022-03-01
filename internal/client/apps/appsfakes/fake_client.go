// Code generated by counterfeiter. DO NOT EDIT.
package appsfakes

import (
	"context"
	"sync"

	"github.com/google/go-github/v41/github"
	"github.com/morphysm/famed-github-backend/internal/client/apps"
)

type FakeClient struct {
	GetAccessTokenStub        func(context.Context, int64) (*github.InstallationToken, error)
	getAccessTokenMutex       sync.RWMutex
	getAccessTokenArgsForCall []struct {
		arg1 context.Context
		arg2 int64
	}
	getAccessTokenReturns struct {
		result1 *github.InstallationToken
		result2 error
	}
	getAccessTokenReturnsOnCall map[int]struct {
		result1 *github.InstallationToken
		result2 error
	}
	GetInstallationsStub        func(context.Context) ([]*github.Installation, error)
	getInstallationsMutex       sync.RWMutex
	getInstallationsArgsForCall []struct {
		arg1 context.Context
	}
	getInstallationsReturns struct {
		result1 []*github.Installation
		result2 error
	}
	getInstallationsReturnsOnCall map[int]struct {
		result1 []*github.Installation
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeClient) GetAccessToken(arg1 context.Context, arg2 int64) (*github.InstallationToken, error) {
	fake.getAccessTokenMutex.Lock()
	ret, specificReturn := fake.getAccessTokenReturnsOnCall[len(fake.getAccessTokenArgsForCall)]
	fake.getAccessTokenArgsForCall = append(fake.getAccessTokenArgsForCall, struct {
		arg1 context.Context
		arg2 int64
	}{arg1, arg2})
	stub := fake.GetAccessTokenStub
	fakeReturns := fake.getAccessTokenReturns
	fake.recordInvocation("GetAccessToken", []interface{}{arg1, arg2})
	fake.getAccessTokenMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeClient) GetAccessTokenCallCount() int {
	fake.getAccessTokenMutex.RLock()
	defer fake.getAccessTokenMutex.RUnlock()
	return len(fake.getAccessTokenArgsForCall)
}

func (fake *FakeClient) GetAccessTokenCalls(stub func(context.Context, int64) (*github.InstallationToken, error)) {
	fake.getAccessTokenMutex.Lock()
	defer fake.getAccessTokenMutex.Unlock()
	fake.GetAccessTokenStub = stub
}

func (fake *FakeClient) GetAccessTokenArgsForCall(i int) (context.Context, int64) {
	fake.getAccessTokenMutex.RLock()
	defer fake.getAccessTokenMutex.RUnlock()
	argsForCall := fake.getAccessTokenArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeClient) GetAccessTokenReturns(result1 *github.InstallationToken, result2 error) {
	fake.getAccessTokenMutex.Lock()
	defer fake.getAccessTokenMutex.Unlock()
	fake.GetAccessTokenStub = nil
	fake.getAccessTokenReturns = struct {
		result1 *github.InstallationToken
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) GetAccessTokenReturnsOnCall(i int, result1 *github.InstallationToken, result2 error) {
	fake.getAccessTokenMutex.Lock()
	defer fake.getAccessTokenMutex.Unlock()
	fake.GetAccessTokenStub = nil
	if fake.getAccessTokenReturnsOnCall == nil {
		fake.getAccessTokenReturnsOnCall = make(map[int]struct {
			result1 *github.InstallationToken
			result2 error
		})
	}
	fake.getAccessTokenReturnsOnCall[i] = struct {
		result1 *github.InstallationToken
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) GetInstallations(arg1 context.Context) ([]*github.Installation, error) {
	fake.getInstallationsMutex.Lock()
	ret, specificReturn := fake.getInstallationsReturnsOnCall[len(fake.getInstallationsArgsForCall)]
	fake.getInstallationsArgsForCall = append(fake.getInstallationsArgsForCall, struct {
		arg1 context.Context
	}{arg1})
	stub := fake.GetInstallationsStub
	fakeReturns := fake.getInstallationsReturns
	fake.recordInvocation("GetInstallations", []interface{}{arg1})
	fake.getInstallationsMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeClient) GetInstallationsCallCount() int {
	fake.getInstallationsMutex.RLock()
	defer fake.getInstallationsMutex.RUnlock()
	return len(fake.getInstallationsArgsForCall)
}

func (fake *FakeClient) GetInstallationsCalls(stub func(context.Context) ([]*github.Installation, error)) {
	fake.getInstallationsMutex.Lock()
	defer fake.getInstallationsMutex.Unlock()
	fake.GetInstallationsStub = stub
}

func (fake *FakeClient) GetInstallationsArgsForCall(i int) context.Context {
	fake.getInstallationsMutex.RLock()
	defer fake.getInstallationsMutex.RUnlock()
	argsForCall := fake.getInstallationsArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeClient) GetInstallationsReturns(result1 []*github.Installation, result2 error) {
	fake.getInstallationsMutex.Lock()
	defer fake.getInstallationsMutex.Unlock()
	fake.GetInstallationsStub = nil
	fake.getInstallationsReturns = struct {
		result1 []*github.Installation
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) GetInstallationsReturnsOnCall(i int, result1 []*github.Installation, result2 error) {
	fake.getInstallationsMutex.Lock()
	defer fake.getInstallationsMutex.Unlock()
	fake.GetInstallationsStub = nil
	if fake.getInstallationsReturnsOnCall == nil {
		fake.getInstallationsReturnsOnCall = make(map[int]struct {
			result1 []*github.Installation
			result2 error
		})
	}
	fake.getInstallationsReturnsOnCall[i] = struct {
		result1 []*github.Installation
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.getAccessTokenMutex.RLock()
	defer fake.getAccessTokenMutex.RUnlock()
	fake.getInstallationsMutex.RLock()
	defer fake.getInstallationsMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeClient) recordInvocation(key string, args []interface{}) {
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

var _ apps.Client = new(FakeClient)
