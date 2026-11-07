# 🌾 AgroVault

> **Programmable Agricultural Credit & Autonomous Harvest Settlement on NEAR Protocol**  
> *Target Tracks: Real-World Money & Programmable Payments — NEAR Hackathon 2026*

---

## 📌 Executive Summary
AgroVault eliminates fund diversion in smallholder farmer financing through **programmable, closed-loop micro-credit**. Instead of disbursing cash loans that can be spent on consumer goods, credit limits are enforced on-chain via smart contracts and can only be routed to whitelisted agricultural input vendors (*Saprotan*). Upon harvest sale, settlements from off-takers are automatically split by NEAR cross-contract promises to settle outstanding debt before streaming the remaining net margin directly to the farmer.

---

## ⚡ Verified On-Chain Deployments (NEAR Testnet)

| Parameter | Details |
| :--- | :--- |
| **Network** | NEAR Testnet |
| **Smart Contract Account** | [`agrovault-1792886368.testnet`](https://explorer.testnet.near.org/accounts/agrovault-1792886368.testnet) |
| **Contract Deployment Tx** | [`AceXwMVxvSbnPSd1TP3Kym2A57GoFoZHn1VEiaVGmYHZ`](https://explorer.testnet.near.org/transactions/AceXwMVxvSbnPSd1TP3Kym2A57GoFoZHn1VEiaVGmYHZ) |
| **Harvest Settlement Tx** | [`Dj6b62cfGQEXRxo8rwTHqaHYwQFCD77fhvX3HcXWbKUd`](https://explorer.testnet.near.org/transactions/Dj6b62cfGQEXRxo8rwTHqaHYwQFCD77fhvX3HcXWbKUd) |
| **Near SDK Version** | `near-sdk v5.29.1` (using optimized `near_sdk::store`) |

---

## 🏗️ Core Architecture & Flow

```text
   [ Lender / Cooperative ]
              │
              ▼ (deposit_liquidity)
     ┌─────────────────┐
     │ AgroVault Pool  │ ── (grant_credit) ──► [ Farmer Credit Limit ]
     └─────────────────┘                               │
              │                                        ▼ (spend_credit_at_vendor)
              ├──────────────────────────────► [ Whitelisted Kiosk / Merchant ]
              │                                 (Input: Fertilizer / Seeds)
              ▼ (settle_harvest)
     [ Off-Taker / Buyer ]
              │
      ┌───────┴────────────────────────┐
      ▼                                ▼
[ Settle Debt Pool ]         [ Direct Farmer Profit ]
  (used_credit: 0)            (Immediate Native Transfer)
1. Restricted Spending (spend_credit_at_vendor)
Farmers cannot withdraw raw NEAR/token cash.

Credit can strictly be redeemed against whitelisted vendors stored in IterableSet<AccountId>.

Guarantees 100% of capital allocation is channeled into verified production inputs.

2. Autonomous Settlement (settle_harvest)
Off-takers purchase crops by depositing funds directly into the contract method.

The contract programmatically deducts the principal loan from the deposited amount.

The remaining net profit is instantly transferred to the farmer's account via non-blocking NEAR Promises.

📂 Repository Structure
contracts/: Near-SDK Rust smart contract managing the credit vault, vendor whitelist, and repayment logic.

relayer/: Go (Gin) service connecting clients to NEAR JSON-RPC, enabling gasless meta-transactions and merchant integrations.

relayer/static/: Mobile-first dashboard for farmers to inspect credit limits and simulate kiosk QR payments.

docs/: Pitch deck slides and architecture diagrams.

🚀 Quickstart & Local Setup
1. Smart Contract Build
Bash
cd contracts
cargo near build non-reproducible-wasm --no-abi
2. Run Backend Relayer & Dashboard
Bash
cd relayer
PORT=8890 go run main.go
Access the dashboard at http://localhost:8890.

🗺️ Roadmap Towards Mainnet
[x] Core Rust smart contract deployed & verified on NEAR Testnet.

[x] Native JSON-RPC integration via Go relayer backend.

[x] Interactive web dashboard for farmer spending simulation.

[ ] Integration with NEP-141 stablecoins (USDT/USDC) to avoid market volatility.

[ ] Account abstraction with FastAuth for seamless seedless onboarding.

[ ] NEAR Intents integration for automated local fiat off-ramping.

📄 License
MIT License. Built for the NEAR Hackathon 2026.
