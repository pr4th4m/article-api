package content

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/olivere/elastic/v7"
	"github.com/pratz/nine-article-api/logger"
	"github.com/pratz/nine-article-api/mocks"
)

func TestArticleSave(t *testing.T) {
	fields := Fields{}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a test suite for success and failure
	testSuite := []struct {
		name        string
		err         error
		expectedErr error
	}{
		{
			name:        "ArticleSaveSuccess",
			err:         nil,
			expectedErr: nil,
		},
		{
			name:        "ArticleSaveFailure",
			err:         fmt.Errorf("not saved"),
			expectedErr: fmt.Errorf("not saved"),
		},
	}
	for _, ut := range testSuite {
		t.Run(
			fmt.Sprintf("Test for %s", ut.name),
			func(t *testing.T) {

				mockES := mocks.NewMockElasticSearch(ctrl)
				mockES.EXPECT().Index(index, typ, fields.ID, fields).Return(nil, ut.err)

				a := &contentArticle{
					esclient: mockES,
					log:      logger.NewNoop(),
				}
				err := a.Save(fields)

				if err != nil {
					if err.Error() != ut.expectedErr.Error() {
						t.Errorf("Test case should have failed for %s", ut.name)
					}
				}
			},
		)
	}
}

func TestArticleGet(t *testing.T) {

	result := &elastic.GetResult{}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a test suite for success and failure
	testSuite := []struct {
		name           string
		input          string
		expectedOutput Fields
		expectedErr    error
		mockOutput     *elastic.GetResult
		mockErr        error
	}{
		{
			name:           "ArticleGetSuccess",
			input:          "1",
			expectedOutput: Fields{},
			expectedErr:    nil,
			mockOutput:     result,
			mockErr:        nil,
		},
		{
			name:           "ArticleGetFailure",
			input:          "2",
			expectedOutput: Fields{},
			expectedErr:    fmt.Errorf("not found"),
			mockOutput:     nil,
			mockErr:        fmt.Errorf("not found"),
		},
	}

	for _, ut := range testSuite {
		t.Run(
			fmt.Sprintf("Test for %s", ut.name),
			func(t *testing.T) {

				mockES := mocks.NewMockElasticSearch(ctrl)
				mockES.EXPECT().Get(index, typ, ut.input).Return(ut.mockOutput, ut.mockErr)

				a := &contentArticle{
					esclient: mockES,
					log:      logger.NewNoop(),
				}
				fields, err := a.Get(ut.input)

				if fields.ID != ut.expectedOutput.ID {
					t.Errorf("Article get output mismatched for %s", ut.input)
				}

				if err != nil {
					if err.Error() != ut.expectedErr.Error() {
						t.Errorf("Test case should have failed for %s", ut.name)
					}
				}
			},
		)
	}
}

func TestArticleSearch(t *testing.T) {

	result := &elastic.SearchResult{}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a test suite for success and failure
	testSuite := []struct {
		name           string
		input          string
		input2         string
		expectedOutput TagSearchResult
		expectedErr    error
		mockInput      string
		mockOutput     *elastic.SearchResult
		mockErr        error
	}{
		{
			name:           "ArticleSearchSuccess",
			input:          "",
			input2:         "20190101",
			expectedOutput: TagSearchResult{},
			expectedErr:    nil,
			mockInput:      "2019-01-01",
			mockOutput:     result,
			mockErr:        nil,
		},
		{
			name:           "ArticleSearchFailure",
			input:          "",
			input2:         "20190101",
			expectedOutput: TagSearchResult{},
			expectedErr:    fmt.Errorf("not found"),
			mockInput:      "2019-01-01",
			mockOutput:     nil,
			mockErr:        fmt.Errorf("not found"),
		},
	}

	for _, ut := range testSuite {
		t.Run(
			fmt.Sprintf("Test for %s", ut.name),
			func(t *testing.T) {

				mockES := mocks.NewMockElasticSearch(ctrl)
				mockES.EXPECT().SearchByTag(index, typ, ut.input, ut.mockInput).Return(ut.mockOutput, ut.mockErr)

				a := &contentArticle{
					esclient: mockES,
					log:      logger.NewNoop(),
				}
				tagResult, err := a.SearchByTag(ut.input, ut.input2)

				if tagResult.Tag != ut.expectedOutput.Tag {
					t.Errorf("Article search output mismatched for %s", ut.input)
				}

				if err != nil {
					if err.Error() != ut.expectedErr.Error() {
						t.Errorf("Test case should have failed for %s", ut.name)
					}
				}
			},
		)
	}
}
