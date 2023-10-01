package main

import (
	"github.com/google/uuid"
)

// GenerateTransactionIDPlugin is a plugin for generating transaction IDs.
type GenerateTransactionIDPlugin struct{}

// Execute generates a transaction ID.
func (p GenerateTransactionIDPlugin) Execute(args ...interface{}) interface{} {
	id := uuid.New()
	return id.String()
}

// GenerateTransactionIDPluginInstance is a variable that stores an instance of the plugin.
var GenerateTransactionIDPluginInstance GenerateTransactionIDPlugin

//go build -buildmode=plugin -o plugins/generate_transaction_id_plugin.so trx/generateTransactionId.go
