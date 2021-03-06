package main

import (
    "os"

    b2 "github.com/kula/vault-plugin-secrets-backblazeb2/plugin"
    hclog "github.com/hashicorp/go-hclog"
    "github.com/hashicorp/vault/helper/pluginutil"
    "github.com/hashicorp/vault/logical/plugin"
)

func main() {
    logger := hclog.New(&hclog.LoggerOptions{})

    defer func() {
	if r:= recover(); r != nil {
	    logger.Error("panic in the plugin", "error", r)
	    os.Exit(1)
	}
    }()

    meta := &pluginutil.APIClientMeta{}

    flags := meta.FlagSet()
    flags.Parse(os.Args[1:])

    tlsConfig := meta.GetTLSConfig()
    tlsProviderFunc := pluginutil.VaultPluginTLSProvider(tlsConfig)

    if err := plugin.Serve(&plugin.ServeOpts{
	BackendFactoryFunc: b2.Factory,
	TLSProviderFunc: tlsProviderFunc,
    }); err != nil {
	logger.Error("plugin shutting down", "error", err);
	os.Exit(1)
    }
}
