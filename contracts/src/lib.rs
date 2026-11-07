use near_sdk::store::{IterableSet, LookupMap};
use near_sdk::{env, near, AccountId, NearToken, PanicOnDefault, Promise};

#[near(serializers = [borsh, json])]
#[derive(Clone)]
pub struct CreditLine {
    pub total_limit: NearToken,
    pub used_credit: NearToken,
    pub is_active: bool,
}

#[near(contract_state)]
#[derive(PanicOnDefault)]
pub struct AgroCreditVault {
    pub admin: AccountId,
    pub whitelisted_vendors: IterableSet<AccountId>,
    pub farmer_credits: LookupMap<AccountId, CreditLine>,
}

#[near]
impl AgroCreditVault {
    #[init]
    pub fn new(admin: AccountId) -> Self {
        Self {
            admin,
            whitelisted_vendors: IterableSet::new(b"v"),
            farmer_credits: LookupMap::new(b"f"),
        }
    }

    // --- ADMIN / LENDER LOGIC ---

    pub fn add_vendor(&mut self, vendor_id: AccountId) {
        self.assert_admin();
        self.whitelisted_vendors.insert(vendor_id);
    }

    pub fn grant_credit(&mut self, farmer_id: AccountId, credit_limit: NearToken) {
        self.assert_admin();
        let credit = CreditLine {
            total_limit: credit_limit,
            used_credit: NearToken::from_yoctonear(0),
            is_active: true,
        };
        self.farmer_credits.insert(farmer_id, credit);
    }

    #[payable]
    pub fn deposit_liquidity(&mut self) {
        env::log_str("Liquidity funded to Agro Vault.");
    }

    // --- PETANI LOGIC: HANYA BISA BELANJA KE VENDOR RESMI ---

    pub fn spend_credit_at_vendor(&mut self, vendor_id: AccountId, amount: NearToken) -> Promise {
        let caller = env::predecessor_account_id();

        // 1. Validasi kios terdaftar
        assert!(
            self.whitelisted_vendors.contains(&vendor_id),
            "Vendor toko pupuk/benih tidak terdaftar!"
        );

        // 2. Validasi plafon kredit petani
        let credit = self
            .farmer_credits
            .get_mut(&caller)
            .expect("Petani belum memiliki akun kredit!");
        assert!(credit.is_active, "Kredit dibekukan!");

        let remaining = credit
            .total_limit
            .as_yoctonear()
            .checked_sub(credit.used_credit.as_yoctonear())
            .expect("Plafon kredit tidak mencukupi!");

        assert!(
            amount.as_yoctonear() <= remaining,
            "Pembelian melebihi sisa plafon kredit!"
        );

        // Update saldo pinjaman
        credit.used_credit = credit.used_credit.saturating_add(amount);

        // Transfer dana langsung ke vendor saprotan
        Promise::new(vendor_id).transfer(amount)
    }

    // --- OFF-TAKER LOGIC: PELUNASAN OTOMATIS SAAT PANEN ---

    #[payable]
    pub fn settle_harvest(&mut self, farmer_id: AccountId) -> Promise {
        let payment = env::attached_deposit();
        assert!(payment.as_yoctonear() > 0, "Deposit harus > 0");

        let credit = self
            .farmer_credits
            .get_mut(&farmer_id)
            .expect("Data petani tidak ditemukan");

        let debt = credit.used_credit;

        if payment >= debt {
            // Lunasi hutang di kontrak
            credit.used_credit = NearToken::from_yoctonear(0);

            // Sisa hasil panen langsung diteruskan ke dompet petani
            let net_profit = payment.saturating_sub(debt);
            if net_profit.as_yoctonear() > 0 {
                Promise::new(farmer_id).transfer(net_profit)
            } else {
                Promise::new(self.admin.clone())
            }
        } else {
            // Pembayaran parsial
            credit.used_credit = credit.used_credit.saturating_sub(payment);
            Promise::new(self.admin.clone())
        }
    }

    // --- VIEW METHODS ---

    pub fn get_farmer_credit(&self, farmer_id: AccountId) -> Option<&CreditLine> {
        self.farmer_credits.get(&farmer_id)
    }

    pub fn is_vendor_whitelisted(&self, vendor_id: AccountId) -> bool {
        self.whitelisted_vendors.contains(&vendor_id)
    }

    fn assert_admin(&self) {
        assert_eq!(
            env::predecessor_account_id(),
            self.admin,
            "Akses ditolak: Hanya admin/lender"
        );
    }
}
