package attributes

import (
	"reflect"
	"testing"

	"github.com/LepikovStan/bCrawler2/attrStorage"
)

func TestStorage_Retry(t *testing.T) {
	type args struct {
		id int
	}
	tests := []struct {
		name string
		s    Storage
		args args
		want int
	}{
		// TODO: Add test cases.
		{
			"n1",
			attrStorage.New(),
			args{},
			1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Retry(tt.args.id); got != tt.want {
				t.Errorf("Storage.Retry() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_Depth(t *testing.T) {
	type args struct {
		id int
	}
	tests := []struct {
		name string
		s    Storage
		args args
		want int
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Depth(tt.args.id); got != tt.want {
				t.Errorf("Storage.Depth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want Storage
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
