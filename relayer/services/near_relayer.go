package services

import (
"bytes"
"encoding/base64"
"encoding/json"
"fmt"
"io"
"net/http"
"os/exec"
"agrovault-backend/config"
)

type NearService struct {
cfg config.Config
}

func NewNearService(cfg config.Config) *NearService {
return &NearService{cfg: cfg}
}

type FarmerCreditResponse struct {
IsActive   bool   `json:"is_active"`
TotalLimit string `json:"total_limit"`
UsedCredit string `json:"used_credit"`
}

// Struct RPC Request & Response NEAR
type rpcRequest struct {
Jsonrpc string      `json:"jsonrpc"`
Id      string      `json:"id"`
Method  string      `json:"method"`
Params  interface{} `json:"params"`
}

type rpcViewFunctionParams struct {
RequestType string `json:"request_type"`
Finality    string `json:"finality"`
AccountId   string `json:"account_id"`
MethodName  string `json:"method_name"`
ArgsBase64  string `json:"args_base64"`
}

type rpcResponse struct {
Result struct {
Result []byte `json:"result"`
Error  string `json:"error"`
} `json:"result"`
Error *struct {
Message string `json:"message"`
Data    string `json:"data"`
} `json:"error"`
}

// GetFarmerCredit memanggil View Function langsung via HTTP JSON-RPC Resmi NEAR
func (s *NearService) GetFarmerCredit(farmerID string) (*FarmerCreditResponse, error) {
argsJSON := fmt.Sprintf(`{"farmer_id":"%s"}`, farmerID)
argsBase64 := base64.StdEncoding.EncodeToString([]byte(argsJSON))

reqBody := rpcRequest{
Jsonrpc: "2.0",
Id:      "dontcare",
Method:  "query",
Params: rpcViewFunctionParams{
RequestType: "call_function",
Finality:    "final",
AccountId:   s.cfg.ContractID,
MethodName:  "get_farmer_credit",
ArgsBase64:  argsBase64,
},
}

payload, err := json.Marshal(reqBody)
if err != nil {
return nil, err
}

resp, err := http.Post(s.cfg.RPCEndpoint, "application/json", bytes.NewBuffer(payload))
if err != nil {
return nil, fmt.Errorf("gagal koneksi ke RPC NEAR: %v", err)
}
defer resp.Body.Close()

body, err := io.ReadAll(resp.Body)
if err != nil {
return nil, err
}

var rpcRes rpcResponse
if err := json.Unmarshal(body, &rpcRes); err != nil {
return nil, fmt.Errorf("gagal unmarshal RPC response: %v", err)
}

if rpcRes.Error != nil {
return nil, fmt.Errorf("RPC Error: %s - %s", rpcRes.Error.Message, rpcRes.Error.Data)
}

if len(rpcRes.Result.Result) == 0 {
return nil, fmt.Errorf("data kredit petani tidak ditemukan")
}

var credit FarmerCreditResponse
if err := json.Unmarshal(rpcRes.Result.Result, &credit); err != nil {
return nil, fmt.Errorf("gagal unmarshal struct data kredit: %v", err)
}

return &credit, nil
}

// SpendCreditRelay mengeksekusi transaksi via near CLI non-interaktif
func (s *NearService) SpendCreditRelay(farmerID string, vendorID string, amountYocto string) (string, error) {
// Memanggil transaksi dengan non-interactive flag yang kompatibel dengan near-cli-rs
cmd := exec.Command(
"near", "contract", "call-function", "as-transaction",
s.cfg.ContractID, "spend_credit_at_vendor",
"json-args", fmt.Sprintf(`{"vendor_id":"%s","amount":"%s"}`, vendorID, amountYocto),
"prepaid-gas", "100.0 Tgas",
"attached-deposit", "0 NEAR",
"sign-as", farmerID,
"network-config", "testnet",
"sign-with-keychain", "send",
)

out, err := cmd.CombinedOutput()
outputStr := string(out)

// Jika flag keychain gagal menemukan keyring, jalankan fallback dengan legacy syntax / credentials flag
if err != nil {
cmdFallback := exec.Command(
"near", "contract", "call-function", "as-transaction",
s.cfg.ContractID, "spend_credit_at_vendor",
"json-args", fmt.Sprintf(`{"vendor_id":"%s","amount":"%s"}`, vendorID, amountYocto),
"prepaid-gas", "100.0 Tgas",
"attached-deposit", "0 NEAR",
"sign-as", farmerID,
"network-config", "testnet",
"sign-with-legacy-keychain", "send",
)
outFallback, errFallback := cmdFallback.CombinedOutput()
if errFallback == nil {
return string(outFallback), nil
}
return "", fmt.Errorf("output: %s", outputStr)
}

return outputStr, nil
}
