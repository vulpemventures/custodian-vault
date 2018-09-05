package btc

// BackendHelp ..
const BackendHelp = "The bitcoin custodian plugin lets you to store your private keys securely.\nWith this plugin you can create wallet or multisig wallet and generate receiving addresses or sign transactions."

// PathSecrets ..
const PathSecrets = "secrets/"

// PathAddress ..
const PathAddress = "address/"

// PathCreds ..
const PathCreds = "creds/"

// PathWallet ..
const PathWallet = "wallet/"

// PathMultiSigAddress ..
const PathMultiSigAddress = PathAddress + "multisig/"

// PathMultiSigCreds ..
const PathMultiSigCreds = PathCreds + "multisig/"

// PathMultiSigWallet ..
const PathMultiSigWallet = PathWallet + "multisig/"

// MultiSigPrefix ..
const MultiSigPrefix = "multisig_"

// MissingTokenError ..
const MissingTokenError = "Missing auth token"

// MissingWalletNameError ..
const MissingWalletNameError = "Missing wallet name"

// MissingNetworkError ..
const MissingNetworkError = "Missing network"

// MissingPubKeysError ..
const MissingPubKeysError = "Missing public keys"

// MissingRawTxError ..
const MissingRawTxError = "Missing raw transaction to sign"

// MissingInternalDataError ..
const MissingInternalDataError = "Secret is missing internal data"

// InvalidTokenError ..
const InvalidTokenError = "Invalid auth token"

// InvalidNetworkError ..
const InvalidNetworkError = "Invalid network"

// InvalidMError ..
const InvalidMError = "Missing or invalid m param: it must be a positive number"

// InvalidNError ..
const InvalidNError = "Missing or invalid n param: it must be a positive number"

// MBiggerThanNError ..
const MBiggerThanNError = "Invalid m param: it must be minor or equal to n"

// NOutOfRangeError ..
const NOutOfRangeError = "Invalid N param: it must be a value between 1 and 7 (inclusive)"

// MOutOfRangeError ..
const MOutOfRangeError = "Invalid M param: it must be between 1 and N (inclusive)"

// PathAddressHelpSyn ..
const PathAddressHelpSyn = "Returns a new receiving address for selected wallet. You must provide an authentication token"

// PathAddressHelpDesc ..
const PathAddressHelpDesc = ""

// PathCredsHelpSyn ..
const PathCredsHelpSyn = "Creates authorization tokens for a wallet"

// PathCredsHelpDesc ..
const PathCredsHelpDesc = ""

// PathMultiSigAddressHelpSyn ..
const PathMultiSigAddressHelpSyn = "Returns the receiving address of the selected multisig wallet"

// PathMultiSigAddressHelpDesc ..
const PathMultiSigAddressHelpDesc = ""

// PathMultiSigCredsHelpSyn ..
const PathMultiSigCredsHelpSyn = "Creates access tokens for a multisig wallet"

// PathMultiSigCredsHelpDesc ..
const PathMultiSigCredsHelpDesc = ""

// PathMultiSigWalletsHelpSyn ..
const PathMultiSigWalletsHelpSyn = "Creates a new wallet that is used as nth key of the m/n multisig wallet"

// PathMultiSigWalletsHelpDesc ..
const PathMultiSigWalletsHelpDesc = ""

// PathTransactionHelpSyn ..
const PathTransactionHelpSyn = "Sign bitcoin raw transaction"

// PathTransactionHelpDesc ..
const PathTransactionHelpDesc = ""

// PathWalletsHelpSyn ..
const PathWalletsHelpSyn = "Creates a new bbip44 wallet by specifying network and name"

// PathWalletsHelpDesc ..
const PathWalletsHelpDesc = ""

// SecretCredsType ..
const SecretCredsType = "creds"

// MultiSigSecretCredsType ..
const MultiSigSecretCredsType = "multisig_creds"

// EntropyBitSize ..
const EntropyBitSize = 256

// HardenedKeyStart ..
const HardenedKeyStart = uint32(0x80000000)

// MainNet ..
const MainNet = "mainnet"

// TestNet ..
const TestNet = "testnet"

// Purpose ..
const Purpose = HardenedKeyStart + 44

// CoinTypeMainNet ..
const CoinTypeMainNet = HardenedKeyStart

// CoinTypeTestNet ..
const CoinTypeTestNet = HardenedKeyStart + 1

// Account ..
const Account = HardenedKeyStart

// Change ..
const Change = HardenedKeyStart

// MultiSigDefaultAddressIndex ..
const MultiSigDefaultAddressIndex = 0

// MinMultiSigN ..
const MinMultiSigN = 1

// MaxMultiSigN ..
const MaxMultiSigN = 7
