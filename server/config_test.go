package server

import (
	"reflect"
	"testing"
)

func TestConfig_GetHubUrls(t *testing.T) {
	tests := []struct {
		name    string
		hubUrls string
		want    []string
	}{
		{
			name:    "default hub URLs when HubUrls is empty",
			hubUrls: "",
			want: []string{
				"https://pubsubhubbub.appspot.com/",
				"https://pubsubhubbub.superfeedr.com/",
			},
		},
		{
			name:    "custom hub URLs",
			hubUrls: "https://hub1.example.com,https://hub2.example.com",
			want: []string{
				"https://hub1.example.com",
				"https://hub2.example.com",
			},
		},
		{
			name:    "single custom hub URL",
			hubUrls: "https://hub.example.com",
			want: []string{
				"https://hub.example.com",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				HubUrls: tt.hubUrls,
			}
			if got := c.GetHubUrls(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config.GetHubUrls() = %v, want %v", got, tt.want)
			}
		})
	}
}