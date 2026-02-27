package db

import (
	"context"
	"testing"
)

// MockStore implements FlatStore for testing.
type MockStore struct {
	IDs      []string
	Written  []string
	ScanErr  error
	WriteErr error
}

func (m *MockStore) GetAllIDs(_ context.Context) ([]string, error) {
	if m.ScanErr != nil {
		return nil, m.ScanErr
	}
	return m.IDs, nil
}

func (m *MockStore) WriteID(_ context.Context, id string) error {
	if m.WriteErr != nil {
		return m.WriteErr
	}
	m.Written = append(m.Written, id)
	return nil
}

func TestMockStore(t *testing.T) {
	store := &MockStore{IDs: []string{"1", "2", "3"}}

	ids, err := store.GetAllIDs(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ids) != 3 {
		t.Errorf("got %d IDs, want 3", len(ids))
	}

	err = store.WriteID(context.Background(), "4")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(store.Written) != 1 || store.Written[0] != "4" {
		t.Errorf("Written = %v, want [4]", store.Written)
	}
}
