package atproto

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/bluesky-social/indigo/atproto/auth/oauth"
	"github.com/bluesky-social/indigo/atproto/syntax"
)

const (
	authStatePath    = ".cedar-auth.json"
	publishStatePath = ".cedar-state.json"
)

// fileAuthStore implements oauth.ClientAuthStore backed by a single JSON file.
// It holds one session (the CLI user) and keeps auth request info in memory.
type fileAuthStore struct {
	mu       sync.Mutex
	requests map[string]oauth.AuthRequestData
}

func newFileAuthStore() *fileAuthStore {
	return &fileAuthStore{requests: make(map[string]oauth.AuthRequestData)}
}

func (s *fileAuthStore) GetSession(_ context.Context, _ syntax.DID, _ string) (*oauth.ClientSessionData, error) {
	data, err := os.ReadFile(authStatePath)
	if err != nil {
		return nil, err
	}
	var sess oauth.ClientSessionData
	if err := json.Unmarshal(data, &sess); err != nil {
		return nil, fmt.Errorf("parsing auth state: %w", err)
	}
	return &sess, nil
}

func (s *fileAuthStore) SaveSession(_ context.Context, sess oauth.ClientSessionData) error {
	data, err := json.MarshalIndent(sess, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(authStatePath, data, 0600)
}

func (s *fileAuthStore) DeleteSession(_ context.Context, _ syntax.DID, _ string) error {
	return os.Remove(authStatePath)
}

func (s *fileAuthStore) GetAuthRequestInfo(_ context.Context, state string) (*oauth.AuthRequestData, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	info, ok := s.requests[state]
	if !ok {
		return nil, fmt.Errorf("auth request not found: %s", state)
	}
	return &info, nil
}

func (s *fileAuthStore) SaveAuthRequestInfo(_ context.Context, info oauth.AuthRequestData) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.requests[info.State] = info
	return nil
}

func (s *fileAuthStore) DeleteAuthRequestInfo(_ context.Context, state string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.requests, state)
	return nil
}

// LoadSession loads the saved session from disk.
func LoadSession() (*oauth.ClientSessionData, error) {
	store := newFileAuthStore()
	return store.GetSession(context.Background(), "", "")
}

// PublicationState tracks a single publication's ATProto record.
type PublicationState struct {
	ATURI string `json:"at_uri"`
	RKey  string `json:"rkey"`
}

// DocumentRecord tracks a single document's ATProto record.
type DocumentRecord struct {
	Publication string `json:"publication"`
	ATURI       string `json:"at_uri"`
	RKey        string `json:"rkey"`
	UpdatedAt   string `json:"updated_at"`
	ContentHash string `json:"content_hash"`
}

// PublishState tracks all published ATProto records across publications.
type PublishState struct {
	Publications map[string]PublicationState `json:"publications"`
	Documents    map[string]DocumentRecord   `json:"documents"`
}

func LoadPublishState() (*PublishState, error) {
	data, err := os.ReadFile(publishStatePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &PublishState{
				Publications: make(map[string]PublicationState),
				Documents:    make(map[string]DocumentRecord),
			}, nil
		}
		return nil, err
	}
	var state PublishState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, err
	}
	if state.Publications == nil {
		state.Publications = make(map[string]PublicationState)
	}
	if state.Documents == nil {
		state.Documents = make(map[string]DocumentRecord)
	}
	return &state, nil
}

func (s *PublishState) Save() error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(publishStatePath, data, 0644)
}
