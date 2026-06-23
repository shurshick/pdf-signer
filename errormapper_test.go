package main

import (
	"errors"
	"testing"
)

func TestFriendlyErrorMessage(t *testing.T) {
	tests := []struct {
		err      error
		contains string
	}{
		{nil, ""},
		{errors.New("certmgr not found"), "certmgr"},
		{errors.New("csptest not found"), "csptest"},
		{errors.New("no certificates found"), "Сертификаты"},
		{errors.New("certificate CN is empty"), "CN"},
		{errors.New("some other error"), "Ошибка"},
	}
	for _, tt := range tests {
		result := FriendlyErrorMessage(tt.err)
		if tt.err == nil {
			if result != "" {
				t.Errorf("FriendlyErrorMessage(nil) = %q, want empty", result)
			}
			continue
		}
		if result == "" {
			t.Errorf("FriendlyErrorMessage(%q) returned empty", tt.err.Error())
		}
	}
}
