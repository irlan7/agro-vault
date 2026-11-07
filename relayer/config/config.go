package config

import "os"

type Config struct {
ContractID  string
RPCEndpoint string
Port        string
}

func LoadConfig() Config {
contractID := os.Getenv("NEAR_CONTRACT_ID")
if contractID == "" {
contractID = "agrovault-1792886368.testnet"
}

rpc := os.Getenv("NEAR_RPC_ENDPOINT")
if rpc == "" {
rpc = "https://rpc.testnet.near.org"
}

port := os.Getenv("PORT")
if port == "" {
port = "8890" // Aman, tidak bertabrakan dengan docker atau microservice lain
}

return Config{
ContractID:  contractID,
RPCEndpoint: rpc,
Port:        port,
}
}
