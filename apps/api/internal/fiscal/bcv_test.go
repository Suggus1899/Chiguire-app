package fiscal

import "testing"

func TestParseBCVRate(t *testing.T) {
	cases := []struct {
		name    string
		html    string
		want    float64
		wantErr bool
	}{
		{
			name: "standard format",
			html: `<div id="dolar" class="r1"><strong>36.500,00</strong></div>`,
			want: 36500.0,
		},
		{
			name: "no thousands separator",
			html: `<div id="dolar"><strong>125,50</strong></div>`,
			want: 125.5,
		},
		{
			name: "whitespace around value",
			html: `<div id="dolar"><strong>  1.234,56  </strong></div>`,
			want: 1234.56,
		},
		{
			name:    "missing dolar div",
			html:    `<div id="euro"><strong>1,00</strong></div>`,
			wantErr: true,
		},
		{
			name:    "missing strong tag",
			html:    `<div id="dolar">no rate here</div>`,
			wantErr: true,
		},
		{
			name:    "non-numeric value",
			html:    `<div id="dolar"><strong>N/D</strong></div>`,
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseBCVRate(tc.html)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error, got %.4f", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Fatalf("got %.4f, want %.4f", got, tc.want)
			}
		})
	}
}
