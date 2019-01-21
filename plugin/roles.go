package b2

import (
    "context"
    "errors"
    "fmt"

    "github.com/hashicorp/errwrap"
    "github.com/hashicorp/vault/logical"
)

var (
    ErrRoleNotFound = errors.New("role not found")
)

// A role stored in the storage backend
type Role struct {

    // Capabilities is a list of strings which reflects
    // the capabilities this key will have in B2
    Capabilities []string `json:"capabilities"`

    // NamePrefix is what we prepend to the key name when we
    // create it, followed by the Vault request ID which asked
    // for the key to be made

    NamePrefix string `json:"name_prefix"`

    // BucketId is an optional restriction to limit this key to 
    // a particular bucket
    BucketId string `json:"bucket_id"`

    // Prefix is an optional restriction to limit which object
    // name prefixes this key can operate on
    Prefix string `json:"prefix"`
}

// List Roles

func (b* backend) ListRoles(ctx context.Context, s logical.Storage) ([]string, error) {
    roles, err := s.List(ctx, "roles/")
    if err != nil {
	return nil, errwrap.Wrapf("Unable to retrieve list of roles: {{err}}", err)
    }

    return roles, nil
}

// Get Role

func (b* backend) GetRole(ctx context.Context, s logical.Storage, role string) (*Role, error) {
    r, err := s.Get(ctx, "roles/"+role)
    if err != nil { 
	return nil, errwrap.Wrapf(fmt.Sprintf("Unable to retrieve role %q: {{err}}", role), err)
    }

    if r == nil {
	return nil, ErrRoleNotFound
    }

    var rv Role
    if err := r.DecodeJSON(&rv); err != nil {
	return nil, errwrap.Wrapf(fmt.Sprintf("Unable to decode role %q: {{err}}", role), err)
    }

    return &rv, nil
}

