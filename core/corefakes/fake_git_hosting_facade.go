// Code generated by counterfeiter. DO NOT EDIT.
package corefakes

import (
	"sync"

	"github.com/ccremer/greposync/core"
)

type FakeGitHostingFacade struct {
	CreateOrUpdateLabelsForRepoStub        func(*core.GitURL, []core.GitRepositoryLabel) error
	createOrUpdateLabelsForRepoMutex       sync.RWMutex
	createOrUpdateLabelsForRepoArgsForCall []struct {
		arg1 *core.GitURL
		arg2 []core.GitRepositoryLabel
	}
	createOrUpdateLabelsForRepoReturns struct {
		result1 error
	}
	createOrUpdateLabelsForRepoReturnsOnCall map[int]struct {
		result1 error
	}
	DeleteLabelsForRepoStub        func(*core.GitURL, []core.GitRepositoryLabel) error
	deleteLabelsForRepoMutex       sync.RWMutex
	deleteLabelsForRepoArgsForCall []struct {
		arg1 *core.GitURL
		arg2 []core.GitRepositoryLabel
	}
	deleteLabelsForRepoReturns struct {
		result1 error
	}
	deleteLabelsForRepoReturnsOnCall map[int]struct {
		result1 error
	}
	InitializeStub        func() error
	initializeMutex       sync.RWMutex
	initializeArgsForCall []struct {
	}
	initializeReturns struct {
		result1 error
	}
	initializeReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeGitHostingFacade) CreateOrUpdateLabelsForRepo(arg1 *core.GitURL, arg2 []core.GitRepositoryLabel) error {
	var arg2Copy []core.GitRepositoryLabel
	if arg2 != nil {
		arg2Copy = make([]core.GitRepositoryLabel, len(arg2))
		copy(arg2Copy, arg2)
	}
	fake.createOrUpdateLabelsForRepoMutex.Lock()
	ret, specificReturn := fake.createOrUpdateLabelsForRepoReturnsOnCall[len(fake.createOrUpdateLabelsForRepoArgsForCall)]
	fake.createOrUpdateLabelsForRepoArgsForCall = append(fake.createOrUpdateLabelsForRepoArgsForCall, struct {
		arg1 *core.GitURL
		arg2 []core.GitRepositoryLabel
	}{arg1, arg2Copy})
	stub := fake.CreateOrUpdateLabelsForRepoStub
	fakeReturns := fake.createOrUpdateLabelsForRepoReturns
	fake.recordInvocation("CreateOrUpdateLabelsForRepo", []interface{}{arg1, arg2Copy})
	fake.createOrUpdateLabelsForRepoMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeGitHostingFacade) CreateOrUpdateLabelsForRepoCallCount() int {
	fake.createOrUpdateLabelsForRepoMutex.RLock()
	defer fake.createOrUpdateLabelsForRepoMutex.RUnlock()
	return len(fake.createOrUpdateLabelsForRepoArgsForCall)
}

func (fake *FakeGitHostingFacade) CreateOrUpdateLabelsForRepoCalls(stub func(*core.GitURL, []core.GitRepositoryLabel) error) {
	fake.createOrUpdateLabelsForRepoMutex.Lock()
	defer fake.createOrUpdateLabelsForRepoMutex.Unlock()
	fake.CreateOrUpdateLabelsForRepoStub = stub
}

func (fake *FakeGitHostingFacade) CreateOrUpdateLabelsForRepoArgsForCall(i int) (*core.GitURL, []core.GitRepositoryLabel) {
	fake.createOrUpdateLabelsForRepoMutex.RLock()
	defer fake.createOrUpdateLabelsForRepoMutex.RUnlock()
	argsForCall := fake.createOrUpdateLabelsForRepoArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeGitHostingFacade) CreateOrUpdateLabelsForRepoReturns(result1 error) {
	fake.createOrUpdateLabelsForRepoMutex.Lock()
	defer fake.createOrUpdateLabelsForRepoMutex.Unlock()
	fake.CreateOrUpdateLabelsForRepoStub = nil
	fake.createOrUpdateLabelsForRepoReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeGitHostingFacade) CreateOrUpdateLabelsForRepoReturnsOnCall(i int, result1 error) {
	fake.createOrUpdateLabelsForRepoMutex.Lock()
	defer fake.createOrUpdateLabelsForRepoMutex.Unlock()
	fake.CreateOrUpdateLabelsForRepoStub = nil
	if fake.createOrUpdateLabelsForRepoReturnsOnCall == nil {
		fake.createOrUpdateLabelsForRepoReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.createOrUpdateLabelsForRepoReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeGitHostingFacade) DeleteLabelsForRepo(arg1 *core.GitURL, arg2 []core.GitRepositoryLabel) error {
	var arg2Copy []core.GitRepositoryLabel
	if arg2 != nil {
		arg2Copy = make([]core.GitRepositoryLabel, len(arg2))
		copy(arg2Copy, arg2)
	}
	fake.deleteLabelsForRepoMutex.Lock()
	ret, specificReturn := fake.deleteLabelsForRepoReturnsOnCall[len(fake.deleteLabelsForRepoArgsForCall)]
	fake.deleteLabelsForRepoArgsForCall = append(fake.deleteLabelsForRepoArgsForCall, struct {
		arg1 *core.GitURL
		arg2 []core.GitRepositoryLabel
	}{arg1, arg2Copy})
	stub := fake.DeleteLabelsForRepoStub
	fakeReturns := fake.deleteLabelsForRepoReturns
	fake.recordInvocation("DeleteLabelsForRepo", []interface{}{arg1, arg2Copy})
	fake.deleteLabelsForRepoMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeGitHostingFacade) DeleteLabelsForRepoCallCount() int {
	fake.deleteLabelsForRepoMutex.RLock()
	defer fake.deleteLabelsForRepoMutex.RUnlock()
	return len(fake.deleteLabelsForRepoArgsForCall)
}

func (fake *FakeGitHostingFacade) DeleteLabelsForRepoCalls(stub func(*core.GitURL, []core.GitRepositoryLabel) error) {
	fake.deleteLabelsForRepoMutex.Lock()
	defer fake.deleteLabelsForRepoMutex.Unlock()
	fake.DeleteLabelsForRepoStub = stub
}

func (fake *FakeGitHostingFacade) DeleteLabelsForRepoArgsForCall(i int) (*core.GitURL, []core.GitRepositoryLabel) {
	fake.deleteLabelsForRepoMutex.RLock()
	defer fake.deleteLabelsForRepoMutex.RUnlock()
	argsForCall := fake.deleteLabelsForRepoArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeGitHostingFacade) DeleteLabelsForRepoReturns(result1 error) {
	fake.deleteLabelsForRepoMutex.Lock()
	defer fake.deleteLabelsForRepoMutex.Unlock()
	fake.DeleteLabelsForRepoStub = nil
	fake.deleteLabelsForRepoReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeGitHostingFacade) DeleteLabelsForRepoReturnsOnCall(i int, result1 error) {
	fake.deleteLabelsForRepoMutex.Lock()
	defer fake.deleteLabelsForRepoMutex.Unlock()
	fake.DeleteLabelsForRepoStub = nil
	if fake.deleteLabelsForRepoReturnsOnCall == nil {
		fake.deleteLabelsForRepoReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.deleteLabelsForRepoReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeGitHostingFacade) Initialize() error {
	fake.initializeMutex.Lock()
	ret, specificReturn := fake.initializeReturnsOnCall[len(fake.initializeArgsForCall)]
	fake.initializeArgsForCall = append(fake.initializeArgsForCall, struct {
	}{})
	stub := fake.InitializeStub
	fakeReturns := fake.initializeReturns
	fake.recordInvocation("Initialize", []interface{}{})
	fake.initializeMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeGitHostingFacade) InitializeCallCount() int {
	fake.initializeMutex.RLock()
	defer fake.initializeMutex.RUnlock()
	return len(fake.initializeArgsForCall)
}

func (fake *FakeGitHostingFacade) InitializeCalls(stub func() error) {
	fake.initializeMutex.Lock()
	defer fake.initializeMutex.Unlock()
	fake.InitializeStub = stub
}

func (fake *FakeGitHostingFacade) InitializeReturns(result1 error) {
	fake.initializeMutex.Lock()
	defer fake.initializeMutex.Unlock()
	fake.InitializeStub = nil
	fake.initializeReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeGitHostingFacade) InitializeReturnsOnCall(i int, result1 error) {
	fake.initializeMutex.Lock()
	defer fake.initializeMutex.Unlock()
	fake.InitializeStub = nil
	if fake.initializeReturnsOnCall == nil {
		fake.initializeReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.initializeReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeGitHostingFacade) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.createOrUpdateLabelsForRepoMutex.RLock()
	defer fake.createOrUpdateLabelsForRepoMutex.RUnlock()
	fake.deleteLabelsForRepoMutex.RLock()
	defer fake.deleteLabelsForRepoMutex.RUnlock()
	fake.initializeMutex.RLock()
	defer fake.initializeMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeGitHostingFacade) recordInvocation(key string, args []interface{}) {
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

var _ core.GitHostingFacade = new(FakeGitHostingFacade)
