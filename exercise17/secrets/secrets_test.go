package secrets

import (
	"path/filepath"
	"testing"
)

func Test_CanCreateFileSecrets(t *testing.T) {
	s := New("encrypt-key", filepath.Join(t.TempDir(), ".secrets"))
	if s == nil {
		t.Error("nil secrets returned")
	}
}

func Test_CanStoreSecretToFile(t *testing.T) {
	s := New("abc", filepath.Join(t.TempDir(), ".secrets"))
	err := s.Store("some-key", "some-value")
	if err != nil {
		t.Errorf("problem storing secret")
	}
}

func Test_CanLoadSecretFromFile(t *testing.T) {
	s := New("abc", filepath.Join(t.TempDir(), ".secrets"))
	_ = s.Store("some-key", "some-value")

	v := s.Load("some-key")
	if v != "some-value" {
		t.Errorf("Load() wanted=\"some-value\", got=%s", v)
	}
}

func Test_CanSaveMultipleSecrets(t *testing.T) {
	s := New("abc", filepath.Join(t.TempDir(), ".secrets"))
	_ = s.Store("some-key", "some-value")
	_ = s.Store("some-key-2", "some-value-2")

	v := s.Load("some-key")
	if v != "some-value" {
		t.Errorf("Load() wanted=\"some-value\", got=%s", v)
	}
	v = s.Load("some-key-2")
	if v != "some-value-2" {
		t.Errorf("Load() wanted=\"some-value-2\", got=%s", v)
	}
}

func Test_CanUpdateExistingSecret(t *testing.T) {
	s := New("abc", filepath.Join(t.TempDir(), ".secrets"))
	_ = s.Store("some-key", "some-value")
	_ = s.Store("some-key", "some-value-2")

	v := s.Load("some-key")
	if v != "some-value-2" {
		t.Errorf("Load() wanted=\"some-value-2\", got=%s", v)
	}
}

func Test_CanUpdateExistingSecretToShorterValue(t *testing.T) {
	s := New("abc", filepath.Join(t.TempDir(), ".secrets"))
	_ = s.Store("some-key", "some-value")
	_ = s.Store("some-key", "value-2")

	v := s.Load("some-key")
	if v != "value-2" {
		t.Errorf("Load() wanted=\"value-2\", got=%s", v)
	}
}
