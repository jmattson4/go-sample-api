// Code generated by mockery v2.2.1. DO NOT EDIT.

package mocks

import (
	domain "github.com/jmattson4/go-sample-api/domain"
	mock "github.com/stretchr/testify/mock"
)

// NewsDBRepository is an autogenerated mock type for the NewsDBRepository type
type NewsDBRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: news
func (_m *NewsDBRepository) Create(news *domain.NewsData) error {
	ret := _m.Called(news)

	var r0 error
	if rf, ok := ret.Get(0).(func(*domain.NewsData) error); ok {
		r0 = rf(news)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: news
func (_m *NewsDBRepository) Delete(news *domain.NewsData) error {
	ret := _m.Called(news)

	var r0 error
	if rf, ok := ret.Get(0).(func(*domain.NewsData) error); ok {
		r0 = rf(news)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields: news
func (_m *NewsDBRepository) Get(news *domain.NewsData) error {
	ret := _m.Called(news)

	var r0 error
	if rf, ok := ret.Get(0).(func(*domain.NewsData) error); ok {
		r0 = rf(news)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetByArticleLink provides a mock function with given fields: news
func (_m *NewsDBRepository) GetByArticleLink(news *domain.NewsData) error {
	ret := _m.Called(news)

	var r0 error
	if rf, ok := ret.Get(0).(func(*domain.NewsData) error); ok {
		r0 = rf(news)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetMultipleNews provides a mock function with given fields: start, count, news
func (_m *NewsDBRepository) GetMultipleNews(start int, count int, news []*domain.NewsData) error {
	ret := _m.Called(start, count, news)

	var r0 error
	if rf, ok := ret.Get(0).(func(int, int, []*domain.NewsData) error); ok {
		r0 = rf(start, count, news)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetMultipleNewsByWebName provides a mock function with given fields: start, count, news, webName
func (_m *NewsDBRepository) GetMultipleNewsByWebName(start int, count int, news []*domain.NewsData, webName string) error {
	ret := _m.Called(start, count, news, webName)

	var r0 error
	if rf, ok := ret.Get(0).(func(int, int, []*domain.NewsData, string) error); ok {
		r0 = rf(start, count, news, webName)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetNewsByWebNameAndID provides a mock function with given fields: news
func (_m *NewsDBRepository) GetNewsByWebNameAndID(news *domain.NewsData) error {
	ret := _m.Called(news)

	var r0 error
	if rf, ok := ret.Get(0).(func(*domain.NewsData) error); ok {
		r0 = rf(news)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// HardDelete provides a mock function with given fields: news
func (_m *NewsDBRepository) HardDelete(news *domain.NewsData) error {
	ret := _m.Called(news)

	var r0 error
	if rf, ok := ret.Get(0).(func(*domain.NewsData) error); ok {
		r0 = rf(news)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RollBack provides a mock function with given fields:
func (_m *NewsDBRepository) RollBack() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Update provides a mock function with given fields: news
func (_m *NewsDBRepository) Update(news *domain.NewsData) error {
	ret := _m.Called(news)

	var r0 error
	if rf, ok := ret.Get(0).(func(*domain.NewsData) error); ok {
		r0 = rf(news)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
