package atproto

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// mustDocValue marshals a documentValue to json.RawMessage for test fixtures.
func mustDocValue(t *testing.T, v documentValue) json.RawMessage {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshalling documentValue: %v", err)
	}
	return json.RawMessage(b)
}

func TestFetchPublicationDocuments_basic(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/xrpc/com.atproto.repo.listRecords" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("collection") != "site.standard.document" {
			t.Errorf("unexpected collection: %s", r.URL.Query().Get("collection"))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(listRecordsResponse{
			Records: []recordEntry{
				{
					URI:   "at://did:plc:abc/site.standard.document/rkey1",
					Value: mustDocValue(t, documentValue{Title: "Hello World", PublishedAt: "2026-01-15T00:00:00Z", Path: "/rkey1", Description: "A post"}),
				},
				{
					URI:   "at://did:plc:abc/site.standard.document/rkey2",
					Value: mustDocValue(t, documentValue{Title: "Second Post", PublishedAt: "2026-02-01T00:00:00Z", Path: "/rkey2"}),
				},
			},
		})
	}))
	defer srv.Close()

	docs, err := FetchPublicationDocuments(t.Context(), srv.URL, "did:plc:abc", "https://example.pub")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(docs) != 2 {
		t.Fatalf("got %d docs, want 2", len(docs))
	}

	// Sorted descending by date — rkey2 (Feb) first
	if docs[0].Title != "Second Post" {
		t.Errorf("docs[0].Title = %q, want %q", docs[0].Title, "Second Post")
	}
	if docs[0].Link != "https://example.pub/rkey2" {
		t.Errorf("docs[0].Link = %q, want %q", docs[0].Link, "https://example.pub/rkey2")
	}
	wantDate := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)
	if !docs[0].Date.Equal(wantDate) {
		t.Errorf("docs[0].Date = %v, want %v", docs[0].Date, wantDate)
	}
	if docs[1].Title != "Hello World" {
		t.Errorf("docs[1].Title = %q, want %q", docs[1].Title, "Hello World")
	}
	if docs[1].Description != "A post" {
		t.Errorf("docs[1].Description = %q, want %q", docs[1].Description, "A post")
	}
}

func TestFetchPublicationDocuments_emptyResponse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(listRecordsResponse{Records: []recordEntry{}})
	}))
	defer srv.Close()

	docs, err := FetchPublicationDocuments(t.Context(), srv.URL, "did:plc:abc", "https://example.pub")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(docs) != 0 {
		t.Fatalf("got %d docs, want 0", len(docs))
	}
}

func TestFetchPublicationDocuments_httpError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	_, err := FetchPublicationDocuments(t.Context(), srv.URL, "did:plc:abc", "https://example.pub")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestFetchPublicationDocuments_skipsRecordsWithNoPath(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(listRecordsResponse{
			Records: []recordEntry{
				{
					URI:   "at://did:plc:abc/site.standard.document/r1",
					Value: mustDocValue(t, documentValue{Title: "Has Path", PublishedAt: "2026-01-01T00:00:00Z", Path: "/r1"}),
				},
				{
					URI:   "at://did:plc:abc/site.standard.document/r2",
					Value: mustDocValue(t, documentValue{Title: "No Path", PublishedAt: "2026-01-02T00:00:00Z", Path: ""}),
				},
			},
		})
	}))
	defer srv.Close()

	docs, err := FetchPublicationDocuments(t.Context(), srv.URL, "did:plc:abc", "https://example.pub")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(docs) != 1 {
		t.Fatalf("got %d docs, want 1 (record without path skipped)", len(docs))
	}
	if docs[0].Title != "Has Path" {
		t.Errorf("docs[0].Title = %q, want %q", docs[0].Title, "Has Path")
	}
}

func TestFetchPublicationDocuments_trailingSlashOnPubURL(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(listRecordsResponse{
			Records: []recordEntry{
				{
					URI:   "at://did:plc:abc/site.standard.document/r1",
					Value: mustDocValue(t, documentValue{Title: "Post", PublishedAt: "2026-01-01T00:00:00Z", Path: "/r1"}),
				},
			},
		})
	}))
	defer srv.Close()

	// pub URL with trailing slash — link should not have double slash
	docs, err := FetchPublicationDocuments(t.Context(), srv.URL, "did:plc:abc", "https://example.pub/")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if docs[0].Link != "https://example.pub/r1" {
		t.Errorf("docs[0].Link = %q, want no double slash", docs[0].Link)
	}
}
