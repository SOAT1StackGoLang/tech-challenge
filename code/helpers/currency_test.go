package helpers

import (
	"github.com/shopspring/decimal"
	"testing"
)

func TestParseDecimalFromString(t *testing.T) {
	type args struct {
		pS string
	}
	tests := []struct {
		name    string
		args    args
		want    decimal.Decimal
		wantErr bool
	}{
		{
			name:    "001_should_parse_correctly",
			args:    args{pS: "R$ 2,00"},
			want:    decimal.NewFromInt(2),
			wantErr: false,
		},
		{
			name:    "002_should_parse_correctly",
			args:    args{pS: "R$2,00"},
			want:    decimal.NewFromInt(2),
			wantErr: false,
		},
		{
			name:    "003_should_parse_correctly",
			args:    args{pS: "2,00"},
			want:    decimal.NewFromInt(2),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDecimalFromString(tt.args.pS)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDecimalFromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.String() != tt.want.String() {
				t.Errorf("ParseDecimalFromString() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseDecimalToString(t *testing.T) {
	type args struct {
		d decimal.Decimal
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "001_should_parse_decimal_to_string",
			args: args{d: decimal.NewFromInt(3)},
			want: "R$ 3,00",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseDecimalToString(tt.args.d); got != tt.want {
				t.Errorf("ParseDecimalToString() = %v, want %v", got, tt.want)
			}
		})
	}
}
