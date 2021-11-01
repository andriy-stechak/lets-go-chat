package hasher

import "testing"

func TestHashPassword(t *testing.T) {
	want := "9b71d224bd62f3785d96d46ad3ea3d73319bfbc2890caadae2dff72519673ca72323c3d99ba5c11d7c7acc6e14b8c5da0c4663475c2e5c3adef46f73bcdec043"
	if got := HashPassword("hello"); want != got {
		t.Errorf("HashPassword() = %q, want %q", got, want)
	}
}

func TestCheckPasswordHash(t *testing.T) {
	want := "9b71d224bd62f3785d96d46ad3ea3d73319bfbc2890caadae2dff72519673ca72323c3d99ba5c11d7c7acc6e14b8c5da0c4663475c2e5c3adef46f73bcdec043"
	if got := CheckPasswordHash("hello", want); !got {
		t.Errorf("CheckPasswordHash() = %t, want %t", got, true)
	}
}
