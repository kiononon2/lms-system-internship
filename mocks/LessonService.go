// Code generated by mockery v2.53.3. DO NOT EDIT.

package mocks

import (
	context "context"
	entities "lms-system-internship/entities"

	mock "github.com/stretchr/testify/mock"
)

// LessonService is an autogenerated mock type for the LessonService type
type LessonService struct {
	mock.Mock
}

// AddLessonToChapter provides a mock function with given fields: ctx, chapterID, lesson
func (_m *LessonService) AddLessonToChapter(ctx context.Context, chapterID uint, lesson *entities.Lesson) error {
	ret := _m.Called(ctx, chapterID, lesson)

	if len(ret) == 0 {
		panic("no return value specified for AddLessonToChapter")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint, *entities.Lesson) error); ok {
		r0 = rf(ctx, chapterID, lesson)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteLesson provides a mock function with given fields: ctx, lessonID
func (_m *LessonService) DeleteLesson(ctx context.Context, lessonID uint) error {
	ret := _m.Called(ctx, lessonID)

	if len(ret) == 0 {
		panic("no return value specified for DeleteLesson")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint) error); ok {
		r0 = rf(ctx, lessonID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAllLessons provides a mock function with given fields: ctx
func (_m *LessonService) GetAllLessons(ctx context.Context) ([]*entities.Lesson, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetAllLessons")
	}

	var r0 []*entities.Lesson
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]*entities.Lesson, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []*entities.Lesson); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*entities.Lesson)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetLesson provides a mock function with given fields: ctx, lessonID
func (_m *LessonService) GetLesson(ctx context.Context, lessonID uint) (*entities.Lesson, error) {
	ret := _m.Called(ctx, lessonID)

	if len(ret) == 0 {
		panic("no return value specified for GetLesson")
	}

	var r0 *entities.Lesson
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint) (*entities.Lesson, error)); ok {
		return rf(ctx, lessonID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint) *entities.Lesson); ok {
		r0 = rf(ctx, lessonID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entities.Lesson)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint) error); ok {
		r1 = rf(ctx, lessonID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ReorderLessons provides a mock function with given fields: ctx, chapterID, orderedLessonIDs
func (_m *LessonService) ReorderLessons(ctx context.Context, chapterID uint, orderedLessonIDs []uint) error {
	ret := _m.Called(ctx, chapterID, orderedLessonIDs)

	if len(ret) == 0 {
		panic("no return value specified for ReorderLessons")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint, []uint) error); ok {
		r0 = rf(ctx, chapterID, orderedLessonIDs)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateLessonContent provides a mock function with given fields: ctx, lessonID, content
func (_m *LessonService) UpdateLessonContent(ctx context.Context, lessonID uint, content string) error {
	ret := _m.Called(ctx, lessonID, content)

	if len(ret) == 0 {
		panic("no return value specified for UpdateLessonContent")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint, string) error); ok {
		r0 = rf(ctx, lessonID, content)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewLessonService creates a new instance of LessonService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewLessonService(t interface {
	mock.TestingT
	Cleanup(func())
}) *LessonService {
	mock := &LessonService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
