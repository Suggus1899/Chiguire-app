package invoice

import (
	"testing"
)

func TestInvoiceNumberFormat(t *testing.T) {
	prefix := "01"
	date := "20260807"
	seq := 1
	got := formatInvoiceNumber(prefix, date, seq)
	want := "01-20260807-00001"
	if got != want {
		t.Errorf("formatInvoiceNumber = %q, want %q", got, want)
	}

	seq = 12345
	got = formatInvoiceNumber(prefix, date, seq)
	want = "01-20260807-12345"
	if got != want {
		t.Errorf("formatInvoiceNumber = %q, want %q", got, want)
	}
}
