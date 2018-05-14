package manager

import (
	"reflect"
	"testing"

	"github.com/LepikovStan/bCrawler2/crawler"
	"github.com/LepikovStan/bCrawler2/parser"
)

func TestComposeCrawlerJob(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
		want *crawler.Job
	}{
		// TODO: Add test cases.
		{
			"n1",
			args{"http://example.com"},
			&crawler.Job{
				Id:             genId(),
				Url:            "http://example.com",
				Headers:        make(map[string]string),
				RequestTimeout: defaultTimeout,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := composeCrawlerJob(tt.args.url); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ComposeCrawlerJob() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestComposeParserJob(t *testing.T) {
	type args struct {
		result *crawler.Result
	}
	tests := []struct {
		name string
		args args
		want *parser.Job
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := composeParserJob(tt.args.result); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ComposeParserJob() = %v, want %v", got, tt.want)
			}
		})
	}
}
