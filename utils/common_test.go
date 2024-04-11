package utils

import (
	"testing"
)

func TestGetTimeDeltaSeconds(t *testing.T) {
	type args struct {
		startTime string
		endTime   string
	}

	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			"case1",
			args{
				"2019-01-01T00:00:00Z",
				"2019-01-01T00:00:01Z",
			},
			1,
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetTimeDeltaSeconds(tt.args.startTime, tt.args.endTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTimeDeltaSeconds() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetTimeDeltaSeconds() = %v, want %v", got, tt.want)
			}
		})
	}
}
