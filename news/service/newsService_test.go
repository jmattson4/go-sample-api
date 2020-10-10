package service_test

import (
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/assert"

	"github.com/jmattson4/go-sample-api/domain"
	"github.com/jmattson4/go-sample-api/news/service"

	"github.com/jmattson4/go-sample-api/domain/mocks"
)

func TestGet(t *testing.T) {

	mockNews := domain.NewsDataInit("testArticleLink", "testArticleText", "testImgLink", "testParagraph", "test.com", 1)

	t.Run("success with-cache", func(t *testing.T) {
		mockNewsCacheRepo := new(mocks.NewsCacheRepository)
		mockNewsDbRepo := new(mocks.NewsDBRepository)
		mockNewsCacheRepo.On("Get", mock.AnythingOfType("*domain.NewsData")).Return(nil).Once()
		u := service.ConstructNewsService(mockNewsDbRepo, mockNewsCacheRepo)

		news := domain.NewsDataBasicInit()
		news.ID = mockNews.ID

		err := u.Get(news)

		assert.NoError(t, err)
		assert.NotNil(t, news)

		mockNewsCacheRepo.AssertExpectations(t)
		mockNewsDbRepo.AssertExpectations(t)

	})
	t.Run("success no-cache", func(t *testing.T) {
		mockNewsCacheRepo := new(mocks.NewsCacheRepository)
		mockNewsDbRepo := new(mocks.NewsDBRepository)
		mockNewsCacheRepo.On("Get", mock.AnythingOfType("*domain.NewsData")).Return(domain.NEWS_STANDARD_TESTING_ERROR).Once()
		mockNewsCacheRepo.On("Create", mock.AnythingOfType("*domain.NewsData")).Return(nil).Once()
		mockNewsDbRepo.On("Get", mock.AnythingOfType("*domain.NewsData")).Return(nil).Once()

		u := service.ConstructNewsService(mockNewsDbRepo, mockNewsCacheRepo)

		news := domain.NewsDataBasicInit()
		news.ID = mockNews.ID

		err := u.Get(news)

		assert.NoError(t, err)
		assert.NotNil(t, news)

		mockNewsCacheRepo.AssertExpectations(t)
		mockNewsDbRepo.AssertExpectations(t)
	})
	t.Run("error failed", func(t *testing.T) {
		mockNewsCacheRepo := new(mocks.NewsCacheRepository)
		mockNewsDbRepo := new(mocks.NewsDBRepository)
		mockNewsCacheRepo.On("Get", mock.AnythingOfType("*domain.NewsData")).Return(domain.NEWS_STANDARD_TESTING_ERROR).Once()
		mockNewsDbRepo.On("Get", mock.AnythingOfType("*domain.NewsData")).Return(domain.NEWS_STANDARD_TESTING_ERROR).Once()

		u := service.ConstructNewsService(mockNewsDbRepo, mockNewsCacheRepo)

		news := domain.NewsDataBasicInit()
		news.ID = mockNews.ID
		expected := news
		err := u.Get(news)

		assert.Error(t, err)
		assert.Equal(t, expected, news)

		mockNewsCacheRepo.AssertExpectations(t)
		mockNewsDbRepo.AssertExpectations(t)

	})
}

func TestCreate(t *testing.T) {

	mockNews := domain.NewsDataInit("testArticleLink", "testArticleText", "testImgLink", "testParagraph", "test.com", 1)
	t.Run("success DB, success Cache", func(t *testing.T) {
		mockNewsCacheRepo := new(mocks.NewsCacheRepository)
		mockNewsDbRepo := new(mocks.NewsDBRepository)
		mockNewsCacheRepo.On("Create", mock.AnythingOfType("*domain.NewsData")).Return(nil).Once()
		mockNewsDbRepo.On("Create", mock.AnythingOfType("*domain.NewsData")).Return(nil).Once()

		u := service.ConstructNewsService(mockNewsDbRepo, mockNewsCacheRepo)

		err := u.Create(mockNews)

		assert.NoError(t, err)
		assert.Nil(t, err)

		mockNewsCacheRepo.AssertExpectations(t)
		mockNewsDbRepo.AssertExpectations(t)
	})

	t.Run("success DB fail cache", func(t *testing.T) {
		mockNewsCacheRepo := new(mocks.NewsCacheRepository)
		mockNewsDbRepo := new(mocks.NewsDBRepository)
		mockNewsCacheRepo.On("Create", mock.AnythingOfType("*domain.NewsData")).Return(domain.NEWS_CACHE_CREATION_ERROR).Once()
		mockNewsDbRepo.On("Create", mock.AnythingOfType("*domain.NewsData")).Return(nil).Once()

		u := service.ConstructNewsService(mockNewsDbRepo, mockNewsCacheRepo)

		err := u.Create(mockNews)

		assert.Error(t, err)
		assert.Equal(t, domain.NEWS_CACHE_CREATION_ERROR, err)

		mockNewsCacheRepo.AssertExpectations(t)
		mockNewsDbRepo.AssertExpectations(t)

	})
	t.Run("fail DB fail Cache", func(t *testing.T) {
		mockNewsCacheRepo := new(mocks.NewsCacheRepository)
		mockNewsDbRepo := new(mocks.NewsDBRepository)
		mockNewsDbRepo.On("Create", mock.AnythingOfType("*domain.NewsData")).Return(domain.NEWS_DB_CREATION_ERROR).Once()

		u := service.ConstructNewsService(mockNewsDbRepo, mockNewsCacheRepo)

		err := u.Create(mockNews)

		assert.Error(t, err)
		assert.Equal(t, domain.NEWS_DB_CREATION_ERROR, err)

		mockNewsCacheRepo.AssertExpectations(t)
		mockNewsDbRepo.AssertExpectations(t)

	})
}

func TestUpdate(t *testing.T) {
	mockNews := domain.NewsDataInit("testArticleLink", "testArticleText", "testImgLink", "testParagraph", "test.com", 1)
	t.Run("PASS DB PASS Cache", func(t *testing.T) {
		mockNewsCacheRepo := new(mocks.NewsCacheRepository)
		mockNewsDbRepo := new(mocks.NewsDBRepository)
		mockNewsCacheRepo.On("Update", mock.AnythingOfType("*domain.NewsData")).Return(nil).Once()
		mockNewsDbRepo.On("Update", mock.AnythingOfType("*domain.NewsData")).Return(nil).Once()

		u := service.ConstructNewsService(mockNewsDbRepo, mockNewsCacheRepo)

		err := u.Update(mockNews)

		assert.NoError(t, err)
		assert.Nil(t, err)

		mockNewsCacheRepo.AssertExpectations(t)
		mockNewsDbRepo.AssertExpectations(t)

	})
	t.Run("PASS DB fail Cache", func(t *testing.T) {
		mockNewsCacheRepo := new(mocks.NewsCacheRepository)
		mockNewsDbRepo := new(mocks.NewsDBRepository)
		mockNewsCacheRepo.On("Update", mock.AnythingOfType("*domain.NewsData")).Return(domain.NEWS_CACHE_UPDATE_ERROR).Once()
		mockNewsDbRepo.On("Update", mock.AnythingOfType("*domain.NewsData")).Return(nil).Once()

		u := service.ConstructNewsService(mockNewsDbRepo, mockNewsCacheRepo)

		err := u.Update(mockNews)

		assert.Error(t, err)
		assert.Equal(t, domain.NEWS_CACHE_UPDATE_ERROR, err)

		mockNewsCacheRepo.AssertExpectations(t)
		mockNewsDbRepo.AssertExpectations(t)

	})
	t.Run("fail DB fail Cache", func(t *testing.T) {
		mockNewsCacheRepo := new(mocks.NewsCacheRepository)
		mockNewsDbRepo := new(mocks.NewsDBRepository)
		mockNewsDbRepo.On("Update", mock.AnythingOfType("*domain.NewsData")).Return(domain.NEWS_DB_UPDATE_ERROR).Once()

		u := service.ConstructNewsService(mockNewsDbRepo, mockNewsCacheRepo)

		err := u.Update(mockNews)

		assert.Error(t, err)
		assert.Equal(t, domain.NEWS_DB_UPDATE_ERROR, err)

		mockNewsCacheRepo.AssertExpectations(t)
		mockNewsDbRepo.AssertExpectations(t)

	})
}
func TestDelete(t *testing.T) {
	mockNews := domain.NewsDataInit("testArticleLink", "testArticleText", "testImgLink", "testParagraph", "test.com", 1)
	t.Run("PASS DB PASS Cache", func(t *testing.T) {
		mockNewsCacheRepo := new(mocks.NewsCacheRepository)
		mockNewsDbRepo := new(mocks.NewsDBRepository)
		mockNewsCacheRepo.On("Delete", mock.AnythingOfType("*domain.NewsData")).Return(nil).Once()
		mockNewsDbRepo.On("Delete", mock.AnythingOfType("*domain.NewsData")).Return(nil).Once()

		u := service.ConstructNewsService(mockNewsDbRepo, mockNewsCacheRepo)

		err := u.Delete(mockNews)

		assert.NoError(t, err)
		assert.Nil(t, err)

		mockNewsCacheRepo.AssertExpectations(t)
		mockNewsDbRepo.AssertExpectations(t)

	})
	t.Run("PASS DB fail Cache", func(t *testing.T) {
		mockNewsCacheRepo := new(mocks.NewsCacheRepository)
		mockNewsDbRepo := new(mocks.NewsDBRepository)
		mockNewsCacheRepo.On("Delete", mock.AnythingOfType("*domain.NewsData")).Return(domain.NEWS_CACHE_DELETE_ERROR).Once()
		mockNewsDbRepo.On("Delete", mock.AnythingOfType("*domain.NewsData")).Return(nil).Once()

		u := service.ConstructNewsService(mockNewsDbRepo, mockNewsCacheRepo)

		err := u.Delete(mockNews)

		assert.Error(t, err)
		assert.Equal(t, domain.NEWS_CACHE_DELETE_ERROR, err)

		mockNewsCacheRepo.AssertExpectations(t)
		mockNewsDbRepo.AssertExpectations(t)

	})
	t.Run("fail DB fail Cache", func(t *testing.T) {
		mockNewsCacheRepo := new(mocks.NewsCacheRepository)
		mockNewsDbRepo := new(mocks.NewsDBRepository)
		mockNewsDbRepo.On("Delete", mock.AnythingOfType("*domain.NewsData")).Return(domain.NEWS_DB_DELETE_ERROR).Once()

		u := service.ConstructNewsService(mockNewsDbRepo, mockNewsCacheRepo)

		err := u.Delete(mockNews)

		assert.Error(t, err)
		assert.Equal(t, domain.NEWS_DB_DELETE_ERROR, err)

		mockNewsCacheRepo.AssertExpectations(t)
		mockNewsDbRepo.AssertExpectations(t)

	})
}
