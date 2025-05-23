// Code generated by mockery v2.53.3. DO NOT EDIT.

package mocks

import (
	context "context"
	entities "lms-system-internship/entities"

	mock "github.com/stretchr/testify/mock"
)

// ChapterService is an autogenerated mock type for the ChapterService type
type ChapterService struct {
	mock.Mock
}

// AddChapterToCourse provides a mock function with given fields: ctx, courseID, chapter
func (_m *ChapterService) AddChapterToCourse(ctx context.Context, courseID uint, chapter *entities.Chapter) error {
	ret := _m.Called(ctx, courseID, chapter)

	if len(ret) == 0 {
		panic("no return value specified for AddChapterToCourse")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint, *entities.Chapter) error); ok {
		r0 = rf(ctx, courseID, chapter)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAllChapters provides a mock function with given fields: ctx
func (_m *ChapterService) GetAllChapters(ctx context.Context) ([]*entities.Chapter, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetAllChapters")
	}

	var r0 []*entities.Chapter
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]*entities.Chapter, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []*entities.Chapter); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*entities.Chapter)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetChapter provides a mock function with given fields: ctx, chapterID
func (_m *ChapterService) GetChapter(ctx context.Context, chapterID uint) (*entities.Chapter, error) {
	ret := _m.Called(ctx, chapterID)

	if len(ret) == 0 {
		panic("no return value specified for GetChapter")
	}

	var r0 *entities.Chapter
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint) (*entities.Chapter, error)); ok {
		return rf(ctx, chapterID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint) *entities.Chapter); ok {
		r0 = rf(ctx, chapterID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entities.Chapter)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint) error); ok {
		r1 = rf(ctx, chapterID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RemoveChapter provides a mock function with given fields: ctx, chapterID
func (_m *ChapterService) RemoveChapter(ctx context.Context, chapterID uint) error {
	ret := _m.Called(ctx, chapterID)

	if len(ret) == 0 {
		panic("no return value specified for RemoveChapter")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint) error); ok {
		r0 = rf(ctx, chapterID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateChapterOrder provides a mock function with given fields: ctx, chapterID, newOrder
func (_m *ChapterService) UpdateChapterOrder(ctx context.Context, chapterID uint, newOrder int) error {
	ret := _m.Called(ctx, chapterID, newOrder)

	if len(ret) == 0 {
		panic("no return value specified for UpdateChapterOrder")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint, int) error); ok {
		r0 = rf(ctx, chapterID, newOrder)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewChapterService creates a new instance of ChapterService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewChapterService(t interface {
	mock.TestingT
	Cleanup(func())
}) *ChapterService {
	mock := &ChapterService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
