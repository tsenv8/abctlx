package config

import "testing"

func TestSealOpen(t *testing.T) {
	s, err := Seal("s3cret", "pw")
	if err != nil {
		t.Fatal(err)
	}
	if got, err := s.Open("pw"); err != nil || got != "s3cret" {
		t.Fatalf("got %q, %v", got, err)
	}
	if _, err := s.Open("nope"); err == nil {
		t.Fatal("wrong passphrase accepted")
	}
}
