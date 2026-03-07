package atproto

import (
	"context"
	"fmt"

	"github.com/bluesky-social/indigo/atproto/atclient"
	"github.com/bluesky-social/indigo/atproto/auth/oauth"
	"github.com/bluesky-social/indigo/atproto/syntax"
)

// Client is an authenticated ATProto XRPC client.
type Client struct {
	api  *atclient.APIClient
	sess *oauth.ClientSession
}

// NewClient creates an authenticated XRPC client from a saved session.
// The caller is responsible for providing the auth store implementation.
func NewClient(sess *oauth.ClientSessionData, store oauth.ClientAuthStore) (*Client, error) {
	oauthCfg := oauth.NewLocalhostConfig("", oauthScopes)
	app := oauth.NewClientApp(&oauthCfg, store)

	clientSess, err := app.ResumeSession(context.Background(), sess.AccountDID, sess.SessionID)
	if err != nil {
		return nil, fmt.Errorf("resuming session: %w", err)
	}

	return &Client{
		api:  clientSess.APIClient(),
		sess: clientSess,
	}, nil
}

type createRecordRequest struct {
	Repo       string `json:"repo"`
	Collection string `json:"collection"`
	Validate   bool   `json:"validate"`
	Record     any    `json:"record"`
}

type createRecordResponse struct {
	URI string `json:"uri"`
	CID string `json:"cid"`
}

// CreateRecord creates a new record in the given collection.
// Returns the AT-URI and rkey.
func (c *Client) CreateRecord(ctx context.Context, collection string, record any) (string, string, error) {
	input := createRecordRequest{
		Repo:       c.sess.Data.AccountDID.String(),
		Collection: collection,
		Validate:   false,
		Record:     record,
	}

	var output createRecordResponse
	if err := c.api.Post(ctx, "com.atproto.repo.createRecord", input, &output); err != nil {
		return "", "", err
	}

	rkey := ""
	if aturi, err := syntax.ParseATURI(output.URI); err == nil {
		rkey = aturi.RecordKey().String()
	}
	return output.URI, rkey, nil
}

type putRecordRequest struct {
	Repo       string `json:"repo"`
	Collection string `json:"collection"`
	Rkey       string `json:"rkey"`
	Validate   bool   `json:"validate"`
	Record     any    `json:"record"`
}

// PutRecord creates or updates a record with a specific rkey.
func (c *Client) PutRecord(ctx context.Context, collection, rkey string, record any) error {
	input := putRecordRequest{
		Repo:       c.sess.Data.AccountDID.String(),
		Collection: collection,
		Rkey:       rkey,
		Validate:   false,
		Record:     record,
	}
	return c.api.Post(ctx, "com.atproto.repo.putRecord", input, nil)
}

type deleteRecordRequest struct {
	Repo       string `json:"repo"`
	Collection string `json:"collection"`
	Rkey       string `json:"rkey"`
}

// DeleteRecord deletes a record by collection and rkey.
func (c *Client) DeleteRecord(ctx context.Context, collection, rkey string) error {
	input := deleteRecordRequest{
		Repo:       c.sess.Data.AccountDID.String(),
		Collection: collection,
		Rkey:       rkey,
	}
	return c.api.Post(ctx, "com.atproto.repo.deleteRecord", input, nil)
}
