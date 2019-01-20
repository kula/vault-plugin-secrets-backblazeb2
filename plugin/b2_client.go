package b2

import (
    "context"
	
	b2client "github.com/kurin/blazer/b2"
	"github.com/hashicorp/vault/logical"
)

// Call this to set a new b2client in the backend.
func (b *backend) newB2Client(ctx context.Context, id string, key string) error {

	b.Logger().Debug("newB2Client", "id", id)

	client, err := b2client.NewClient(ctx, id, key)

	if err != nil {
		b.Logger().Error("Error getting new b2 client", "error", err)
		return err
	}

	b.Logger().Debug("Getting clientMutex.Lock")
	b.clientMutex.Lock()
	defer b.clientMutex.Unlock()

	b.client = client
	b.Logger().Debug("Set new b.client, unlocking and returning")
	return nil
}

// Convenience function to get the b2client
func (b *backend) getB2Client(ctx context.Context, s logical.Storage) (*b2client.Client, error) {
	b.Logger().Debug("getB2Client, getting clientMutex.RLock")
	b.clientMutex.RLock()
	if b.client != nil {
		b.Logger().Debug("have client already, unlocking and returning")
		b.clientMutex.RUnlock()
		return b.client, nil
	}
	b.clientMutex.RUnlock()

	// We don't have a current client, look up the id and key
	// from the current configuration and create a new client

	b.Logger().Info("Getting new b2 client, fetching config")
	c, err := b.GetConfig(ctx, s)
	if err != nil {
		b.Logger().Error("Error fetching configuration to make new b2client", "error", err)
		return nil, err
	}

	if c.KeyId == "" {
	    b.Logger().Error("KeyID not set when trying to create new client")
	    return nil, err
	}

	if c.Key == "" {
	    b.Logger().Error("Key not set when trying to create new client")
	    return nil, err
	}

	if err := b.newB2Client(ctx, c.KeyId, c.Key); err != nil {
		return nil, err
	}

	return b.client, nil
}
