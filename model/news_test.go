package model_test

import (
	"testing"

	"github.com/jmattson4/go-sample-api/model"
)

var news *model.NewsData

func TestCreate(t *testing.T) {
	news = model.NewsDataInit("TestArticleLink4", "TestArticleText", "tst/art/img/link", "Test Paragraphs", "testWeb")
	createErr := news.Create()
	if createErr != nil {
		t.Errorf("Error Create did not occur because: %v", createErr)
	}
}
func TestGet(t *testing.T) {
	getErr := news.Get()
	if getErr != nil {
		t.Errorf("Error Get did not occur because: %v.", getErr)
	}
}
func TestGetByArticleLink(t *testing.T) {
	getErr := news.GetNewsByArticleLink()
	if getErr != nil {
		t.Errorf("Error Get did not occur because: %v.", getErr)
	}
}
func TestUpdate(t *testing.T) {
	news.ImageURL = "test/test"
	getErr := news.Update()
	if getErr != nil {
		t.Errorf("Error Update did not occur because: %v.", getErr)
	}
	testNews := &model.NewsData{}
	testNews.ID = news.ID
	if testErr := testNews.Get(); testErr != nil {
		t.Errorf("Error could not get newly Updated Value because: %v.", testErr)
	}
	if testNews.ArticleLink != news.ArticleLink {
		t.Error("Error Updated value does not match")
	}
}
func TestDelete(t *testing.T) {
	err := news.Delete()
	if err != nil {
		t.Errorf("Error delete failed because: %v.", err)
	}
}
func TestHardDelete(t *testing.T) {
	err := news.HardDelete()
	if err != nil {
		t.Errorf("Error delete failed because: %v.", err)
	}
}
