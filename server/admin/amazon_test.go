package admin

import "testing"

func Test_amazonShortUrlToAsin(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "valid url",
			args: args{
				url: "https://amzn.to/42051PN",
			},
			want:    "B01M2BOZDL",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := amazonShortUrlToAsin(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("amazonShortUrlToAsin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("amazonShortUrlToAsin() got = %v, want %v", got, tt.want)
			}
		})
	}
}
