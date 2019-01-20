package b2

import (
    "context"
    "fmt"
    "strings"

    "github.com/hashicorp/errwrap"
    "github.com/hashicorp/vault/logical"
    "github.com/hashicorp/vault/logical/framework"

    b2client "github.com/kurin/blazer/b2"
)

// Define the rotate path
func (b *backend) pathConfigRotate() *framework.Path {
    return &framework.Path{
	Pattern: fmt.Sprintf("config/rotate/?$"),
	HelpSynopsis: "Use the existing key to generate a set a new key",
	HelpDescription: "Use this endpoint to use the current key to generate a new key, and use that",

	Fields: map[string]*framework.FieldSchema{
	    "key_name": &framework.FieldSchema{
		Type: framework.TypeString,
		Description: "The name for the newly generated key.",
	    },
	},

	Callbacks: map[logical.Operation]framework.OperationFunc{
	    logical.UpdateOperation: b.pathRotateKey,
	},
    }
}

// Rotate the key
func (b *backend) pathRotateKey(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {

    // Get the current b.client before we blow it away
    client, err := b.getB2Client(ctx, req.Storage)
    if err != nil {
	return nil, err
    }

    // Fetch configuration    
    c, err := b.GetConfig(ctx, req.Storage);
    if err != nil {
	return nil, err
    }

    // Save the old KeyId so we can destroy it
    oldKeyId := c.KeyId

    // Set new key options
    var opts []b2client.KeyOption
    opts = append(opts, b2client.Capabilities("writeKeys"))

    // Set key_name
    var newKeyName string
    if v, ok := d.GetOk("key_name"); ok {
        newKeyName = strings.TrimSpace(v.(string))
    } else {
	newKeyName = fmt.Sprintf("%s-%s", "vault-root-key", req.ID)
    }

    // Create new key

    newKey, err := client.CreateKey(ctx, newKeyName, opts...)
    if err != nil {
	return nil, err
    }

    // Gin up a FieldData to pass to Update
    nd := &framework.FieldData{
	Schema: map[string]*framework.FieldSchema{
	    "key_id": &framework.FieldSchema{Type: framework.TypeString},
	    "key": &framework.FieldSchema{Type: framework.TypeString},
	    "key_name": &framework.FieldSchema{Type: framework.TypeString},
	},
	Raw: map[string]interface{}{
	    "key_id": newKey.ID,
	    "key": newKey.Secret(),
	    "key_name": newKey.Name(),
	},
    }

    // Update the internal configuration
    changed, err := c.Update(nd)
    if err != nil {
	return nil, logical.CodedError(400, err.Error())
    }

    // If we changed the configuration, store it
    if changed {
	// Make a new storage entry
	entry, err := logical.StorageEntryJSON("config", c)
	if err != nil {
	    return nil, errwrap.Wrapf("failed to generate JSON configuration: {{err}}", err)
	}

	// And store it
	if err := req.Storage.Put(ctx, entry); err != nil {
	    return nil, errwrap.Wrapf("Failed to persist configuration: {{err}}", err)
	}

    }

    // Destroy key
    b.Logger().Info("Deleting previous key", "id", oldKeyId)

    // Destroy any old b2client which may exist so we get a new one
    // with the next request

    b.client = nil

    return nil, nil
}
