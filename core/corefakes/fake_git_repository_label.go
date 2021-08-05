// Code generated by counterfeiter. DO NOT EDIT.
package corefakes

import (
	"sync"

	"github.com/ccremer/greposync/core"
)

type FakeGitRepositoryLabel struct {
	IsBoundForDeletionStub        func() bool
	isBoundForDeletionMutex       sync.RWMutex
	isBoundForDeletionArgsForCall []struct {
	}
	isBoundForDeletionReturns struct {
		result1 bool
	}
	isBoundForDeletionReturnsOnCall map[int]struct {
		result1 bool
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeGitRepositoryLabel) IsBoundForDeletion() bool {
	fake.isBoundForDeletionMutex.Lock()
	ret, specificReturn := fake.isBoundForDeletionReturnsOnCall[len(fake.isBoundForDeletionArgsForCall)]
	fake.isBoundForDeletionArgsForCall = append(fake.isBoundForDeletionArgsForCall, struct {
	}{})
	stub := fake.IsBoundForDeletionStub
	fakeReturns := fake.isBoundForDeletionReturns
	fake.recordInvocation("IsBoundForDeletion", []interface{}{})
	fake.isBoundForDeletionMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeGitRepositoryLabel) IsBoundForDeletionCallCount() int {
	fake.isBoundForDeletionMutex.RLock()
	defer fake.isBoundForDeletionMutex.RUnlock()
	return len(fake.isBoundForDeletionArgsForCall)
}

func (fake *FakeGitRepositoryLabel) IsBoundForDeletionCalls(stub func() bool) {
	fake.isBoundForDeletionMutex.Lock()
	defer fake.isBoundForDeletionMutex.Unlock()
	fake.IsBoundForDeletionStub = stub
}

func (fake *FakeGitRepositoryLabel) IsBoundForDeletionReturns(result1 bool) {
	fake.isBoundForDeletionMutex.Lock()
	defer fake.isBoundForDeletionMutex.Unlock()
	fake.IsBoundForDeletionStub = nil
	fake.isBoundForDeletionReturns = struct {
		result1 bool
	}{result1}
}

func (fake *FakeGitRepositoryLabel) IsBoundForDeletionReturnsOnCall(i int, result1 bool) {
	fake.isBoundForDeletionMutex.Lock()
	defer fake.isBoundForDeletionMutex.Unlock()
	fake.IsBoundForDeletionStub = nil
	if fake.isBoundForDeletionReturnsOnCall == nil {
		fake.isBoundForDeletionReturnsOnCall = make(map[int]struct {
			result1 bool
		})
	}
	fake.isBoundForDeletionReturnsOnCall[i] = struct {
		result1 bool
	}{result1}
}

func (fake *FakeGitRepositoryLabel) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.isBoundForDeletionMutex.RLock()
	defer fake.isBoundForDeletionMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeGitRepositoryLabel) recordInvocation(key string, args []interface{}) {
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

var _ core.GitRepositoryLabel = new(FakeGitRepositoryLabel)