// Code generated by mockery v2.4.0-beta. DO NOT EDIT.

package services

import (
	context "context"

	github "github.com/google/go-github/v39/github"
	mock "github.com/stretchr/testify/mock"
)

// MockRepositoriesService is an autogenerated mock type for the RepositoriesService type
type MockRepositoriesService struct {
	mock.Mock
}

// CreateFile provides a mock function with given fields: ctx, owner, repo, path, opts
func (_m *MockRepositoriesService) CreateFile(ctx context.Context, owner string, repo string, path string, opts *github.RepositoryContentFileOptions) (*github.RepositoryContentResponse, *github.Response, error) {
	ret := _m.Called(ctx, owner, repo, path, opts)

	var r0 *github.RepositoryContentResponse
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string, *github.RepositoryContentFileOptions) *github.RepositoryContentResponse); ok {
		r0 = rf(ctx, owner, repo, path, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*github.RepositoryContentResponse)
		}
	}

	var r1 *github.Response
	if rf, ok := ret.Get(1).(func(context.Context, string, string, string, *github.RepositoryContentFileOptions) *github.Response); ok {
		r1 = rf(ctx, owner, repo, path, opts)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*github.Response)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, string, string, string, *github.RepositoryContentFileOptions) error); ok {
		r2 = rf(ctx, owner, repo, path, opts)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetContents provides a mock function with given fields: ctx, owner, repo, path, opts
func (_m *MockRepositoriesService) GetContents(ctx context.Context, owner string, repo string, path string, opts *github.RepositoryContentGetOptions) (*github.RepositoryContent, []*github.RepositoryContent, *github.Response, error) {
	ret := _m.Called(ctx, owner, repo, path, opts)

	var r0 *github.RepositoryContent
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string, *github.RepositoryContentGetOptions) *github.RepositoryContent); ok {
		r0 = rf(ctx, owner, repo, path, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*github.RepositoryContent)
		}
	}

	var r1 []*github.RepositoryContent
	if rf, ok := ret.Get(1).(func(context.Context, string, string, string, *github.RepositoryContentGetOptions) []*github.RepositoryContent); ok {
		r1 = rf(ctx, owner, repo, path, opts)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).([]*github.RepositoryContent)
		}
	}

	var r2 *github.Response
	if rf, ok := ret.Get(2).(func(context.Context, string, string, string, *github.RepositoryContentGetOptions) *github.Response); ok {
		r2 = rf(ctx, owner, repo, path, opts)
	} else {
		if ret.Get(2) != nil {
			r2 = ret.Get(2).(*github.Response)
		}
	}

	var r3 error
	if rf, ok := ret.Get(3).(func(context.Context, string, string, string, *github.RepositoryContentGetOptions) error); ok {
		r3 = rf(ctx, owner, repo, path, opts)
	} else {
		r3 = ret.Error(3)
	}

	return r0, r1, r2, r3
}
