package b2

import (
    "context"
    "fmt"
    "time"

    "github.com/hashicorp/vault/logical"
    "github.com/hashicorp/vault/logical/framework"

    b2client "github.com/kurin/blazer/b2"
)

const b2KeyType = "b2_application_key"

func (b *backend) b2ApplicationsKeys() *framework.Secret {
    return &framework.Secret{
	Type: b2KeyType,
	Fields: map[string]*framework.FieldSchema{
	    "keyName": &framework.FieldSchema{
		Type: framework.TypeString,
		Description: "Application key name",
	    },
	    "applicationKeyId": &framework.FieldSchema{
		Type: framework.TypeString,
		Description: "Application Key ID",
	    },
	    "capabilities": &framework.FieldSchema{
		Type: framework.TypeCommaStringSlice,
		Description: "List of key capabilities",
	    },
	    "accountId": &framework.FieldSchema{
		Type: framework.TypeString,
		Description: "Account ID the key is associated with",
	    },
	    "bucketName": &framework.FieldSchema{
		Type: framework.TypeString,
		Description: "(Optional) Bucket name to which the key is restricted",
	    },
	    "namePrefix": &framework.FieldSchema{
		Type: framework.TypeString,
		Description: "(Optional) Object name prefix to which the key is restricted",
	    },
	    // We skip expirationTimestamp because that will be set to
	    // the Vault TTL for this secret
	},

	Revoke: b.b2ApplicationKeyRevoke,
    }
}

func (b *backend) b2ApplicationKeyCreate(ctx context.Context, s logical.Storage,
    keyName, accountId string, capabilities []string,
    bucketName, namePrefix string,
    lifetime time.Duration) (*b2client.Key, error) {

    client, err := b.getB2Client(ctx, s)
    if err != nil {
	return nil, err
    }

    // Set key options
    var keyOpts []b2client.KeyOption
    
    // Set TTL, if lifetime is > 0
    if lifetime > 0 {
	keyOpts = append(keyOpts, b2client.Lifetime(lifetime))
    }

    // Set capabilities
    keyOpts = append(keyOpts, b2client.Capabilities(capabilities...))

    // If bucketName is set, look up the bucket and create that way.
    // Else, create the key directly. This is how Blazer does it.

    if bucketName != "" {
	bucket, err := client.Bucket(ctx, bucketName)
	if err != nil {
	    return nil, err
	}
	// Set prefix if asked for
	if namePrefix != "" {
	    keyOpts = append(keyOpts, b2client.Prefix(namePrefix))
	}

	newKey, err := bucket.CreateKey(ctx, keyName, keyOpts...)
	if err != nil {
	    return nil, err
	}
	return newKey, nil

    } else {
	newKey, err := client.CreateKey(ctx, keyName, keyOpts...)
	if err != nil {
	    return nil, err
	}
	return newKey, nil
    }

}


func (b *backend) b2ApplicationKeyRevoke(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {

    client, err := b.getB2Client(ctx, req.Storage)
    if err != nil {
	return nil, err
    }

    // Get applicationKeyId from secret internal data
    applicationKeyIdRaw, ok := req.Secret.InternalData["applicationKeyId"]
    if !ok {
	return nil, fmt.Errorf("secret is missing internal applicationKeyId")
    }

    applicationKeyId, ok := applicationKeyIdRaw.(string)
    if !ok {
	return nil, fmt.Errorf("secret is missing internal applicationKeyId")
    }

    // Find key
    var applicationKey *b2client.Key
    keys, _, err := client.ListKeys(ctx, 1, applicationKeyId)
    if err != nil {
	return nil, err
    }

    // We should only get one, but verify
    for _, key := range keys {
	if key.ID() == applicationKeyId {
	    applicationKey = key
	    break
	}
    }

    if applicationKey == nil {
	return nil, fmt.Errorf("cannot find key in b2")
    }

    if err := applicationKey.Delete(ctx); err != nil {
	return nil, err
    }

    return nil, nil
}
