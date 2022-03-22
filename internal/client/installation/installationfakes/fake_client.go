// Code generated by counterfeiter. DO NOT EDIT.
package installationfakes

import (
	"context"
	"net/http"
	"sync"

	"github.com/google/go-github/v41/github"
	"github.com/morphysm/famed-github-backend/internal/client/installation"
)

type FakeClient struct {
	AddInstallationStub        func(string, int64) error
	addInstallationMutex       sync.RWMutex
	addInstallationArgsForCall []struct {
		arg1 string
		arg2 int64
	}
	addInstallationReturns struct {
		result1 error
	}
	addInstallationReturnsOnCall map[int]struct {
		result1 error
	}
	CheckInstallationStub        func(string) bool
	checkInstallationMutex       sync.RWMutex
	checkInstallationArgsForCall []struct {
		arg1 string
	}
	checkInstallationReturns struct {
		result1 bool
	}
	checkInstallationReturnsOnCall map[int]struct {
		result1 bool
	}
	GetCommentsStub        func(context.Context, string, string, int) ([]*github.IssueComment, error)
	getCommentsMutex       sync.RWMutex
	getCommentsArgsForCall []struct {
		arg1 context.Context
		arg2 string
		arg3 string
		arg4 int
	}
	getCommentsReturns struct {
		result1 []*github.IssueComment
		result2 error
	}
	getCommentsReturnsOnCall map[int]struct {
		result1 []*github.IssueComment
		result2 error
	}
	GetIssueEventsStub        func(context.Context, string, string, int) ([]installation.IssueEvent, error)
	getIssueEventsMutex       sync.RWMutex
	getIssueEventsArgsForCall []struct {
		arg1 context.Context
		arg2 string
		arg3 string
		arg4 int
	}
	getIssueEventsReturns struct {
		result1 []installation.IssueEvent
		result2 error
	}
	getIssueEventsReturnsOnCall map[int]struct {
		result1 []installation.IssueEvent
		result2 error
	}
	GetIssuePullRequestStub        func(context.Context, string, string, int) (*installation.PullRequest, error)
	getIssuePullRequestMutex       sync.RWMutex
	getIssuePullRequestArgsForCall []struct {
		arg1 context.Context
		arg2 string
		arg3 string
		arg4 int
	}
	getIssuePullRequestReturns struct {
		result1 *installation.PullRequest
		result2 error
	}
	getIssuePullRequestReturnsOnCall map[int]struct {
		result1 *installation.PullRequest
		result2 error
	}
	GetIssuesByRepoStub        func(context.Context, string, string, []string, installation.IssueState) ([]installation.Issue, error)
	getIssuesByRepoMutex       sync.RWMutex
	getIssuesByRepoArgsForCall []struct {
		arg1 context.Context
		arg2 string
		arg3 string
		arg4 []string
		arg5 installation.IssueState
	}
	getIssuesByRepoReturns struct {
		result1 []installation.Issue
		result2 error
	}
	getIssuesByRepoReturnsOnCall map[int]struct {
		result1 []installation.Issue
		result2 error
	}
	PostCommentStub        func(context.Context, string, string, int, string) error
	postCommentMutex       sync.RWMutex
	postCommentArgsForCall []struct {
		arg1 context.Context
		arg2 string
		arg3 string
		arg4 int
		arg5 string
	}
	postCommentReturns struct {
		result1 error
	}
	postCommentReturnsOnCall map[int]struct {
		result1 error
	}
	PostLabelStub        func(context.Context, string, string, installation.Label) error
	postLabelMutex       sync.RWMutex
	postLabelArgsForCall []struct {
		arg1 context.Context
		arg2 string
		arg3 string
		arg4 installation.Label
	}
	postLabelReturns struct {
		result1 error
	}
	postLabelReturnsOnCall map[int]struct {
		result1 error
	}
	PostLabelsStub        func(context.Context, string, []installation.Repository, map[string]installation.Label) []error
	postLabelsMutex       sync.RWMutex
	postLabelsArgsForCall []struct {
		arg1 context.Context
		arg2 string
		arg3 []installation.Repository
		arg4 map[string]installation.Label
	}
	postLabelsReturns struct {
		result1 []error
	}
	postLabelsReturnsOnCall map[int]struct {
		result1 []error
	}
	UpdateCommentStub        func(context.Context, string, string, int64, string) error
	updateCommentMutex       sync.RWMutex
	updateCommentArgsForCall []struct {
		arg1 context.Context
		arg2 string
		arg3 string
		arg4 int64
		arg5 string
	}
	updateCommentReturns struct {
		result1 error
	}
	updateCommentReturnsOnCall map[int]struct {
		result1 error
	}
	ValidateWebHookEventStub        func(*http.Request) (interface{}, error)
	validateWebHookEventMutex       sync.RWMutex
	validateWebHookEventArgsForCall []struct {
		arg1 *http.Request
	}
	validateWebHookEventReturns struct {
		result1 interface{}
		result2 error
	}
	validateWebHookEventReturnsOnCall map[int]struct {
		result1 interface{}
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeClient) AddInstallation(arg1 string, arg2 int64) error {
	fake.addInstallationMutex.Lock()
	ret, specificReturn := fake.addInstallationReturnsOnCall[len(fake.addInstallationArgsForCall)]
	fake.addInstallationArgsForCall = append(fake.addInstallationArgsForCall, struct {
		arg1 string
		arg2 int64
	}{arg1, arg2})
	stub := fake.AddInstallationStub
	fakeReturns := fake.addInstallationReturns
	fake.recordInvocation("AddInstallation", []interface{}{arg1, arg2})
	fake.addInstallationMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeClient) AddInstallationCallCount() int {
	fake.addInstallationMutex.RLock()
	defer fake.addInstallationMutex.RUnlock()
	return len(fake.addInstallationArgsForCall)
}

func (fake *FakeClient) AddInstallationCalls(stub func(string, int64) error) {
	fake.addInstallationMutex.Lock()
	defer fake.addInstallationMutex.Unlock()
	fake.AddInstallationStub = stub
}

func (fake *FakeClient) AddInstallationArgsForCall(i int) (string, int64) {
	fake.addInstallationMutex.RLock()
	defer fake.addInstallationMutex.RUnlock()
	argsForCall := fake.addInstallationArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeClient) AddInstallationReturns(result1 error) {
	fake.addInstallationMutex.Lock()
	defer fake.addInstallationMutex.Unlock()
	fake.AddInstallationStub = nil
	fake.addInstallationReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeClient) AddInstallationReturnsOnCall(i int, result1 error) {
	fake.addInstallationMutex.Lock()
	defer fake.addInstallationMutex.Unlock()
	fake.AddInstallationStub = nil
	if fake.addInstallationReturnsOnCall == nil {
		fake.addInstallationReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.addInstallationReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeClient) CheckInstallation(arg1 string) bool {
	fake.checkInstallationMutex.Lock()
	ret, specificReturn := fake.checkInstallationReturnsOnCall[len(fake.checkInstallationArgsForCall)]
	fake.checkInstallationArgsForCall = append(fake.checkInstallationArgsForCall, struct {
		arg1 string
	}{arg1})
	stub := fake.CheckInstallationStub
	fakeReturns := fake.checkInstallationReturns
	fake.recordInvocation("CheckInstallation", []interface{}{arg1})
	fake.checkInstallationMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeClient) CheckInstallationCallCount() int {
	fake.checkInstallationMutex.RLock()
	defer fake.checkInstallationMutex.RUnlock()
	return len(fake.checkInstallationArgsForCall)
}

func (fake *FakeClient) CheckInstallationCalls(stub func(string) bool) {
	fake.checkInstallationMutex.Lock()
	defer fake.checkInstallationMutex.Unlock()
	fake.CheckInstallationStub = stub
}

func (fake *FakeClient) CheckInstallationArgsForCall(i int) string {
	fake.checkInstallationMutex.RLock()
	defer fake.checkInstallationMutex.RUnlock()
	argsForCall := fake.checkInstallationArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeClient) CheckInstallationReturns(result1 bool) {
	fake.checkInstallationMutex.Lock()
	defer fake.checkInstallationMutex.Unlock()
	fake.CheckInstallationStub = nil
	fake.checkInstallationReturns = struct {
		result1 bool
	}{result1}
}

func (fake *FakeClient) CheckInstallationReturnsOnCall(i int, result1 bool) {
	fake.checkInstallationMutex.Lock()
	defer fake.checkInstallationMutex.Unlock()
	fake.CheckInstallationStub = nil
	if fake.checkInstallationReturnsOnCall == nil {
		fake.checkInstallationReturnsOnCall = make(map[int]struct {
			result1 bool
		})
	}
	fake.checkInstallationReturnsOnCall[i] = struct {
		result1 bool
	}{result1}
}

func (fake *FakeClient) GetComments(arg1 context.Context, arg2 string, arg3 string, arg4 int) ([]*github.IssueComment, error) {
	fake.getCommentsMutex.Lock()
	ret, specificReturn := fake.getCommentsReturnsOnCall[len(fake.getCommentsArgsForCall)]
	fake.getCommentsArgsForCall = append(fake.getCommentsArgsForCall, struct {
		arg1 context.Context
		arg2 string
		arg3 string
		arg4 int
	}{arg1, arg2, arg3, arg4})
	stub := fake.GetCommentsStub
	fakeReturns := fake.getCommentsReturns
	fake.recordInvocation("GetComments", []interface{}{arg1, arg2, arg3, arg4})
	fake.getCommentsMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3, arg4)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeClient) GetCommentsCallCount() int {
	fake.getCommentsMutex.RLock()
	defer fake.getCommentsMutex.RUnlock()
	return len(fake.getCommentsArgsForCall)
}

func (fake *FakeClient) GetCommentsCalls(stub func(context.Context, string, string, int) ([]*github.IssueComment, error)) {
	fake.getCommentsMutex.Lock()
	defer fake.getCommentsMutex.Unlock()
	fake.GetCommentsStub = stub
}

func (fake *FakeClient) GetCommentsArgsForCall(i int) (context.Context, string, string, int) {
	fake.getCommentsMutex.RLock()
	defer fake.getCommentsMutex.RUnlock()
	argsForCall := fake.getCommentsArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4
}

func (fake *FakeClient) GetCommentsReturns(result1 []*github.IssueComment, result2 error) {
	fake.getCommentsMutex.Lock()
	defer fake.getCommentsMutex.Unlock()
	fake.GetCommentsStub = nil
	fake.getCommentsReturns = struct {
		result1 []*github.IssueComment
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) GetCommentsReturnsOnCall(i int, result1 []*github.IssueComment, result2 error) {
	fake.getCommentsMutex.Lock()
	defer fake.getCommentsMutex.Unlock()
	fake.GetCommentsStub = nil
	if fake.getCommentsReturnsOnCall == nil {
		fake.getCommentsReturnsOnCall = make(map[int]struct {
			result1 []*github.IssueComment
			result2 error
		})
	}
	fake.getCommentsReturnsOnCall[i] = struct {
		result1 []*github.IssueComment
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) GetIssueEvents(arg1 context.Context, arg2 string, arg3 string, arg4 int) ([]installation.IssueEvent, error) {
	fake.getIssueEventsMutex.Lock()
	ret, specificReturn := fake.getIssueEventsReturnsOnCall[len(fake.getIssueEventsArgsForCall)]
	fake.getIssueEventsArgsForCall = append(fake.getIssueEventsArgsForCall, struct {
		arg1 context.Context
		arg2 string
		arg3 string
		arg4 int
	}{arg1, arg2, arg3, arg4})
	stub := fake.GetIssueEventsStub
	fakeReturns := fake.getIssueEventsReturns
	fake.recordInvocation("GetIssueEvents", []interface{}{arg1, arg2, arg3, arg4})
	fake.getIssueEventsMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3, arg4)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeClient) GetIssueEventsCallCount() int {
	fake.getIssueEventsMutex.RLock()
	defer fake.getIssueEventsMutex.RUnlock()
	return len(fake.getIssueEventsArgsForCall)
}

func (fake *FakeClient) GetIssueEventsCalls(stub func(context.Context, string, string, int) ([]installation.IssueEvent, error)) {
	fake.getIssueEventsMutex.Lock()
	defer fake.getIssueEventsMutex.Unlock()
	fake.GetIssueEventsStub = stub
}

func (fake *FakeClient) GetIssueEventsArgsForCall(i int) (context.Context, string, string, int) {
	fake.getIssueEventsMutex.RLock()
	defer fake.getIssueEventsMutex.RUnlock()
	argsForCall := fake.getIssueEventsArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4
}

func (fake *FakeClient) GetIssueEventsReturns(result1 []installation.IssueEvent, result2 error) {
	fake.getIssueEventsMutex.Lock()
	defer fake.getIssueEventsMutex.Unlock()
	fake.GetIssueEventsStub = nil
	fake.getIssueEventsReturns = struct {
		result1 []installation.IssueEvent
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) GetIssueEventsReturnsOnCall(i int, result1 []installation.IssueEvent, result2 error) {
	fake.getIssueEventsMutex.Lock()
	defer fake.getIssueEventsMutex.Unlock()
	fake.GetIssueEventsStub = nil
	if fake.getIssueEventsReturnsOnCall == nil {
		fake.getIssueEventsReturnsOnCall = make(map[int]struct {
			result1 []installation.IssueEvent
			result2 error
		})
	}
	fake.getIssueEventsReturnsOnCall[i] = struct {
		result1 []installation.IssueEvent
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) GetIssuePullRequest(arg1 context.Context, arg2 string, arg3 string, arg4 int) (*installation.PullRequest, error) {
	fake.getIssuePullRequestMutex.Lock()
	ret, specificReturn := fake.getIssuePullRequestReturnsOnCall[len(fake.getIssuePullRequestArgsForCall)]
	fake.getIssuePullRequestArgsForCall = append(fake.getIssuePullRequestArgsForCall, struct {
		arg1 context.Context
		arg2 string
		arg3 string
		arg4 int
	}{arg1, arg2, arg3, arg4})
	stub := fake.GetIssuePullRequestStub
	fakeReturns := fake.getIssuePullRequestReturns
	fake.recordInvocation("GetIssuePullRequest", []interface{}{arg1, arg2, arg3, arg4})
	fake.getIssuePullRequestMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3, arg4)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeClient) GetIssuePullRequestCallCount() int {
	fake.getIssuePullRequestMutex.RLock()
	defer fake.getIssuePullRequestMutex.RUnlock()
	return len(fake.getIssuePullRequestArgsForCall)
}

func (fake *FakeClient) GetIssuePullRequestCalls(stub func(context.Context, string, string, int) (*installation.PullRequest, error)) {
	fake.getIssuePullRequestMutex.Lock()
	defer fake.getIssuePullRequestMutex.Unlock()
	fake.GetIssuePullRequestStub = stub
}

func (fake *FakeClient) GetIssuePullRequestArgsForCall(i int) (context.Context, string, string, int) {
	fake.getIssuePullRequestMutex.RLock()
	defer fake.getIssuePullRequestMutex.RUnlock()
	argsForCall := fake.getIssuePullRequestArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4
}

func (fake *FakeClient) GetIssuePullRequestReturns(result1 *installation.PullRequest, result2 error) {
	fake.getIssuePullRequestMutex.Lock()
	defer fake.getIssuePullRequestMutex.Unlock()
	fake.GetIssuePullRequestStub = nil
	fake.getIssuePullRequestReturns = struct {
		result1 *installation.PullRequest
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) GetIssuePullRequestReturnsOnCall(i int, result1 *installation.PullRequest, result2 error) {
	fake.getIssuePullRequestMutex.Lock()
	defer fake.getIssuePullRequestMutex.Unlock()
	fake.GetIssuePullRequestStub = nil
	if fake.getIssuePullRequestReturnsOnCall == nil {
		fake.getIssuePullRequestReturnsOnCall = make(map[int]struct {
			result1 *installation.PullRequest
			result2 error
		})
	}
	fake.getIssuePullRequestReturnsOnCall[i] = struct {
		result1 *installation.PullRequest
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) GetIssuesByRepo(arg1 context.Context, arg2 string, arg3 string, arg4 []string, arg5 installation.IssueState) ([]installation.Issue, error) {
	var arg4Copy []string
	if arg4 != nil {
		arg4Copy = make([]string, len(arg4))
		copy(arg4Copy, arg4)
	}
	fake.getIssuesByRepoMutex.Lock()
	ret, specificReturn := fake.getIssuesByRepoReturnsOnCall[len(fake.getIssuesByRepoArgsForCall)]
	fake.getIssuesByRepoArgsForCall = append(fake.getIssuesByRepoArgsForCall, struct {
		arg1 context.Context
		arg2 string
		arg3 string
		arg4 []string
		arg5 installation.IssueState
	}{arg1, arg2, arg3, arg4Copy, arg5})
	stub := fake.GetIssuesByRepoStub
	fakeReturns := fake.getIssuesByRepoReturns
	fake.recordInvocation("GetIssuesByRepo", []interface{}{arg1, arg2, arg3, arg4Copy, arg5})
	fake.getIssuesByRepoMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3, arg4, arg5)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeClient) GetIssuesByRepoCallCount() int {
	fake.getIssuesByRepoMutex.RLock()
	defer fake.getIssuesByRepoMutex.RUnlock()
	return len(fake.getIssuesByRepoArgsForCall)
}

func (fake *FakeClient) GetIssuesByRepoCalls(stub func(context.Context, string, string, []string, installation.IssueState) ([]installation.Issue, error)) {
	fake.getIssuesByRepoMutex.Lock()
	defer fake.getIssuesByRepoMutex.Unlock()
	fake.GetIssuesByRepoStub = stub
}

func (fake *FakeClient) GetIssuesByRepoArgsForCall(i int) (context.Context, string, string, []string, installation.IssueState) {
	fake.getIssuesByRepoMutex.RLock()
	defer fake.getIssuesByRepoMutex.RUnlock()
	argsForCall := fake.getIssuesByRepoArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4, argsForCall.arg5
}

func (fake *FakeClient) GetIssuesByRepoReturns(result1 []installation.Issue, result2 error) {
	fake.getIssuesByRepoMutex.Lock()
	defer fake.getIssuesByRepoMutex.Unlock()
	fake.GetIssuesByRepoStub = nil
	fake.getIssuesByRepoReturns = struct {
		result1 []installation.Issue
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) GetIssuesByRepoReturnsOnCall(i int, result1 []installation.Issue, result2 error) {
	fake.getIssuesByRepoMutex.Lock()
	defer fake.getIssuesByRepoMutex.Unlock()
	fake.GetIssuesByRepoStub = nil
	if fake.getIssuesByRepoReturnsOnCall == nil {
		fake.getIssuesByRepoReturnsOnCall = make(map[int]struct {
			result1 []installation.Issue
			result2 error
		})
	}
	fake.getIssuesByRepoReturnsOnCall[i] = struct {
		result1 []installation.Issue
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) PostComment(arg1 context.Context, arg2 string, arg3 string, arg4 int, arg5 string) error {
	fake.postCommentMutex.Lock()
	ret, specificReturn := fake.postCommentReturnsOnCall[len(fake.postCommentArgsForCall)]
	fake.postCommentArgsForCall = append(fake.postCommentArgsForCall, struct {
		arg1 context.Context
		arg2 string
		arg3 string
		arg4 int
		arg5 string
	}{arg1, arg2, arg3, arg4, arg5})
	stub := fake.PostCommentStub
	fakeReturns := fake.postCommentReturns
	fake.recordInvocation("PostComment", []interface{}{arg1, arg2, arg3, arg4, arg5})
	fake.postCommentMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3, arg4, arg5)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeClient) PostCommentCallCount() int {
	fake.postCommentMutex.RLock()
	defer fake.postCommentMutex.RUnlock()
	return len(fake.postCommentArgsForCall)
}

func (fake *FakeClient) PostCommentCalls(stub func(context.Context, string, string, int, string) error) {
	fake.postCommentMutex.Lock()
	defer fake.postCommentMutex.Unlock()
	fake.PostCommentStub = stub
}

func (fake *FakeClient) PostCommentArgsForCall(i int) (context.Context, string, string, int, string) {
	fake.postCommentMutex.RLock()
	defer fake.postCommentMutex.RUnlock()
	argsForCall := fake.postCommentArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4, argsForCall.arg5
}

func (fake *FakeClient) PostCommentReturns(result1 error) {
	fake.postCommentMutex.Lock()
	defer fake.postCommentMutex.Unlock()
	fake.PostCommentStub = nil
	fake.postCommentReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeClient) PostCommentReturnsOnCall(i int, result1 error) {
	fake.postCommentMutex.Lock()
	defer fake.postCommentMutex.Unlock()
	fake.PostCommentStub = nil
	if fake.postCommentReturnsOnCall == nil {
		fake.postCommentReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.postCommentReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeClient) PostLabel(arg1 context.Context, arg2 string, arg3 string, arg4 installation.Label) error {
	fake.postLabelMutex.Lock()
	ret, specificReturn := fake.postLabelReturnsOnCall[len(fake.postLabelArgsForCall)]
	fake.postLabelArgsForCall = append(fake.postLabelArgsForCall, struct {
		arg1 context.Context
		arg2 string
		arg3 string
		arg4 installation.Label
	}{arg1, arg2, arg3, arg4})
	stub := fake.PostLabelStub
	fakeReturns := fake.postLabelReturns
	fake.recordInvocation("PostLabel", []interface{}{arg1, arg2, arg3, arg4})
	fake.postLabelMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3, arg4)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeClient) PostLabelCallCount() int {
	fake.postLabelMutex.RLock()
	defer fake.postLabelMutex.RUnlock()
	return len(fake.postLabelArgsForCall)
}

func (fake *FakeClient) PostLabelCalls(stub func(context.Context, string, string, installation.Label) error) {
	fake.postLabelMutex.Lock()
	defer fake.postLabelMutex.Unlock()
	fake.PostLabelStub = stub
}

func (fake *FakeClient) PostLabelArgsForCall(i int) (context.Context, string, string, installation.Label) {
	fake.postLabelMutex.RLock()
	defer fake.postLabelMutex.RUnlock()
	argsForCall := fake.postLabelArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4
}

func (fake *FakeClient) PostLabelReturns(result1 error) {
	fake.postLabelMutex.Lock()
	defer fake.postLabelMutex.Unlock()
	fake.PostLabelStub = nil
	fake.postLabelReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeClient) PostLabelReturnsOnCall(i int, result1 error) {
	fake.postLabelMutex.Lock()
	defer fake.postLabelMutex.Unlock()
	fake.PostLabelStub = nil
	if fake.postLabelReturnsOnCall == nil {
		fake.postLabelReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.postLabelReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeClient) PostLabels(arg1 context.Context, arg2 string, arg3 []installation.Repository, arg4 map[string]installation.Label) []error {
	var arg3Copy []installation.Repository
	if arg3 != nil {
		arg3Copy = make([]installation.Repository, len(arg3))
		copy(arg3Copy, arg3)
	}
	fake.postLabelsMutex.Lock()
	ret, specificReturn := fake.postLabelsReturnsOnCall[len(fake.postLabelsArgsForCall)]
	fake.postLabelsArgsForCall = append(fake.postLabelsArgsForCall, struct {
		arg1 context.Context
		arg2 string
		arg3 []installation.Repository
		arg4 map[string]installation.Label
	}{arg1, arg2, arg3Copy, arg4})
	stub := fake.PostLabelsStub
	fakeReturns := fake.postLabelsReturns
	fake.recordInvocation("PostLabels", []interface{}{arg1, arg2, arg3Copy, arg4})
	fake.postLabelsMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3, arg4)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeClient) PostLabelsCallCount() int {
	fake.postLabelsMutex.RLock()
	defer fake.postLabelsMutex.RUnlock()
	return len(fake.postLabelsArgsForCall)
}

func (fake *FakeClient) PostLabelsCalls(stub func(context.Context, string, []installation.Repository, map[string]installation.Label) []error) {
	fake.postLabelsMutex.Lock()
	defer fake.postLabelsMutex.Unlock()
	fake.PostLabelsStub = stub
}

func (fake *FakeClient) PostLabelsArgsForCall(i int) (context.Context, string, []installation.Repository, map[string]installation.Label) {
	fake.postLabelsMutex.RLock()
	defer fake.postLabelsMutex.RUnlock()
	argsForCall := fake.postLabelsArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4
}

func (fake *FakeClient) PostLabelsReturns(result1 []error) {
	fake.postLabelsMutex.Lock()
	defer fake.postLabelsMutex.Unlock()
	fake.PostLabelsStub = nil
	fake.postLabelsReturns = struct {
		result1 []error
	}{result1}
}

func (fake *FakeClient) PostLabelsReturnsOnCall(i int, result1 []error) {
	fake.postLabelsMutex.Lock()
	defer fake.postLabelsMutex.Unlock()
	fake.PostLabelsStub = nil
	if fake.postLabelsReturnsOnCall == nil {
		fake.postLabelsReturnsOnCall = make(map[int]struct {
			result1 []error
		})
	}
	fake.postLabelsReturnsOnCall[i] = struct {
		result1 []error
	}{result1}
}

func (fake *FakeClient) UpdateComment(arg1 context.Context, arg2 string, arg3 string, arg4 int64, arg5 string) error {
	fake.updateCommentMutex.Lock()
	ret, specificReturn := fake.updateCommentReturnsOnCall[len(fake.updateCommentArgsForCall)]
	fake.updateCommentArgsForCall = append(fake.updateCommentArgsForCall, struct {
		arg1 context.Context
		arg2 string
		arg3 string
		arg4 int64
		arg5 string
	}{arg1, arg2, arg3, arg4, arg5})
	stub := fake.UpdateCommentStub
	fakeReturns := fake.updateCommentReturns
	fake.recordInvocation("UpdateComment", []interface{}{arg1, arg2, arg3, arg4, arg5})
	fake.updateCommentMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3, arg4, arg5)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeClient) UpdateCommentCallCount() int {
	fake.updateCommentMutex.RLock()
	defer fake.updateCommentMutex.RUnlock()
	return len(fake.updateCommentArgsForCall)
}

func (fake *FakeClient) UpdateCommentCalls(stub func(context.Context, string, string, int64, string) error) {
	fake.updateCommentMutex.Lock()
	defer fake.updateCommentMutex.Unlock()
	fake.UpdateCommentStub = stub
}

func (fake *FakeClient) UpdateCommentArgsForCall(i int) (context.Context, string, string, int64, string) {
	fake.updateCommentMutex.RLock()
	defer fake.updateCommentMutex.RUnlock()
	argsForCall := fake.updateCommentArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4, argsForCall.arg5
}

func (fake *FakeClient) UpdateCommentReturns(result1 error) {
	fake.updateCommentMutex.Lock()
	defer fake.updateCommentMutex.Unlock()
	fake.UpdateCommentStub = nil
	fake.updateCommentReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeClient) UpdateCommentReturnsOnCall(i int, result1 error) {
	fake.updateCommentMutex.Lock()
	defer fake.updateCommentMutex.Unlock()
	fake.UpdateCommentStub = nil
	if fake.updateCommentReturnsOnCall == nil {
		fake.updateCommentReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.updateCommentReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeClient) ValidateWebHookEvent(arg1 *http.Request) (interface{}, error) {
	fake.validateWebHookEventMutex.Lock()
	ret, specificReturn := fake.validateWebHookEventReturnsOnCall[len(fake.validateWebHookEventArgsForCall)]
	fake.validateWebHookEventArgsForCall = append(fake.validateWebHookEventArgsForCall, struct {
		arg1 *http.Request
	}{arg1})
	stub := fake.ValidateWebHookEventStub
	fakeReturns := fake.validateWebHookEventReturns
	fake.recordInvocation("ValidateWebHookEvent", []interface{}{arg1})
	fake.validateWebHookEventMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeClient) ValidateWebHookEventCallCount() int {
	fake.validateWebHookEventMutex.RLock()
	defer fake.validateWebHookEventMutex.RUnlock()
	return len(fake.validateWebHookEventArgsForCall)
}

func (fake *FakeClient) ValidateWebHookEventCalls(stub func(*http.Request) (interface{}, error)) {
	fake.validateWebHookEventMutex.Lock()
	defer fake.validateWebHookEventMutex.Unlock()
	fake.ValidateWebHookEventStub = stub
}

func (fake *FakeClient) ValidateWebHookEventArgsForCall(i int) *http.Request {
	fake.validateWebHookEventMutex.RLock()
	defer fake.validateWebHookEventMutex.RUnlock()
	argsForCall := fake.validateWebHookEventArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeClient) ValidateWebHookEventReturns(result1 interface{}, result2 error) {
	fake.validateWebHookEventMutex.Lock()
	defer fake.validateWebHookEventMutex.Unlock()
	fake.ValidateWebHookEventStub = nil
	fake.validateWebHookEventReturns = struct {
		result1 interface{}
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) ValidateWebHookEventReturnsOnCall(i int, result1 interface{}, result2 error) {
	fake.validateWebHookEventMutex.Lock()
	defer fake.validateWebHookEventMutex.Unlock()
	fake.ValidateWebHookEventStub = nil
	if fake.validateWebHookEventReturnsOnCall == nil {
		fake.validateWebHookEventReturnsOnCall = make(map[int]struct {
			result1 interface{}
			result2 error
		})
	}
	fake.validateWebHookEventReturnsOnCall[i] = struct {
		result1 interface{}
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.addInstallationMutex.RLock()
	defer fake.addInstallationMutex.RUnlock()
	fake.checkInstallationMutex.RLock()
	defer fake.checkInstallationMutex.RUnlock()
	fake.getCommentsMutex.RLock()
	defer fake.getCommentsMutex.RUnlock()
	fake.getIssueEventsMutex.RLock()
	defer fake.getIssueEventsMutex.RUnlock()
	fake.getIssuePullRequestMutex.RLock()
	defer fake.getIssuePullRequestMutex.RUnlock()
	fake.getIssuesByRepoMutex.RLock()
	defer fake.getIssuesByRepoMutex.RUnlock()
	fake.postCommentMutex.RLock()
	defer fake.postCommentMutex.RUnlock()
	fake.postLabelMutex.RLock()
	defer fake.postLabelMutex.RUnlock()
	fake.postLabelsMutex.RLock()
	defer fake.postLabelsMutex.RUnlock()
	fake.updateCommentMutex.RLock()
	defer fake.updateCommentMutex.RUnlock()
	fake.validateWebHookEventMutex.RLock()
	defer fake.validateWebHookEventMutex.RUnlock()
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

var _ installation.Client = new(FakeClient)
