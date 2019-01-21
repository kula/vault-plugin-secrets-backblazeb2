# vault-plugin-secrets-backblazeb2

**NOTE**: This is a very initial release of this code, please test
in a non-critical environment before you use it.

This is a plugin for [HashiCorp Vault][vault] which will provision
API keys for the [Backblaze B2 Cloud Storage][b2] service. A good
deal of help was gleaned from the [vault-plugin-secrets-helloworld][helloworld]
plugin from @daveadams, and the Vault builtin AWS and database secrets
engines.

## Usage

Once the plugin is registered with your vault instance, you can enable it
on a particular path:

    $ vault secrets enable \
		-path=b2 \
		-plugin-name=backblazeb2-plugin \
		-description="Instance of the Backblaze B2 plugin" \
		plugin

### Configuration

In order to configure the plugin instance, you must supply it with either your
Backblaze B2 key id and value. This can be either your account master key, or
a key with the `writeKeys` capability.

    $ vault write b2/config \
		account_id=<account id> \
		key_id=<key id or account id> \
		key=<key value> \

If `key_id` is not supplied, it's assumed that you are using your account 
master key, and the `key_id` is the same as your `account_id`.

You can read the current configuration:

    $ vault read b2/config

This returns the `account_id`, `key_id` and `key_name` of the currently
configured key.

You can also rotate the key stored in the plugin configuration. This will 
cause the plugin to call the `b2_create_key` call to create a new key
with the `writeKeys` capability and store it, discarding the previously
used key. This key cannot be extracted.

    $ vault write b2/config/rotate \
		key_name=<optional key name>

The `<optional key name>` is B2 key name. If not supplied, it defaults
to `vault-plugin-secrets-backblazeb2`.

### Roles

Before you can issue keys, you must define a role. A role defines the 
set of capabilities the key will have, the name prefix for the key,
and any optional bucket or path restrictions.

    $ vault write b2/roles/example-role \
		capabilities=<comma separated list of capabilities> \
		name_prefix=<name prefix> \
		bucket=<optional bucket restriction> \
		path=<optional path restriction> \

`<name prefix>` is prefixed to the Vault request id for a key request,
and defaults to an empty string. Having the Vault request id as the 
latter part of the name allows you to trace the key issuer via the Vault
audit log. If you set a prefix, please note the limitations for
the key name, and that the Vault request id is (currently) 36 characters
in length.

    $ vault read b2/roles/example-role

Returns the capabilities, name_prefix and optional bucket and path restrictions
for the role.

    $ vault list b2/roles

Lists all configured roles.

### Provisioning keys

    $ vault read b2/keys/example-role

Returns the keyName, applicationKeyId, applicationKey, capabilities, accountId,
expirationTimestamp, bucketId and namePrefix of the key as outlined in the
`b2_create_key` API call.


[vault]: https://www.vaultproject.io
[b2]: https://www.backblaze.com/b2/
[helloworld]: https://github.com/daveadams/vault-plugin-secrets-helloworld
