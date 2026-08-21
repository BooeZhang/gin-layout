package apperror

import (
	"errors"
	"fmt"
	"testing"
)

func TestErrorCarriesBusinessMetadata(t *testing.T) {
	err := New(NotFound, 20000, "用户不存在")

	if got := err.Kind(); got != NotFound {
		t.Fatalf("Kind() = %v, want %v", got, NotFound)
	}
	if got := err.Code(); got != 20000 {
		t.Fatalf("Code() = %d, want %d", got, 20000)
	}
	if got := err.Message(); got != "用户不存在" {
		t.Fatalf("Message() = %q, want %q", got, "用户不存在")
	}
	if got := err.Error(); got != "用户不存在" {
		t.Fatalf("Error() = %q, want %q", got, "用户不存在")
	}
}

func TestErrorRemainsIdentifiableWhenWrapped(t *testing.T) {
	target := New(Conflict, 20020, "账号已存在")
	err := fmt.Errorf("create user: %w", target)

	if !errors.Is(err, target) {
		t.Fatal("wrapped error is not identifiable with errors.Is")
	}

	var got *Error
	if !errors.As(err, &got) {
		t.Fatal("wrapped error is not identifiable with errors.As")
	}
	if got != target {
		t.Fatal("errors.As returned a different error value")
	}
}

func TestUnknownKindFallsBackToInternalServerError(t *testing.T) {
	err := New(Unknown, 50002, "未知错误")

	if got := err.Kind(); got != Unknown {
		t.Fatalf("Kind() = %v, want %v", got, Unknown)
	}
}
