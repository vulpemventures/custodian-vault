package btc

// PathWallet ..
const PathWallet = "wallet/"

// PathAddress ..
const PathAddress = "address/"

// PathCreds ..
const PathCreds = "creds/"

// PathSecrets ..
const PathSecrets = "secrets/"

// PathMultiSigWallet ..
const PathMultiSigWallet = PathWallet + "multisig/"

// PathMultiSigAddress ..
const PathMultiSigAddress = PathAddress + "multisig/"

// PathMultiSigCreds ..
const PathMultiSigCreds = PathCreds + "multisig/"

// PathSegWitWallet ..
const PathSegWitWallet = PathWallet + "segwit/"

// PathSegWitAddress ..
const PathSegWitAddress = PathAddress + "segwit/"

// PathSegWitCreds ..
const PathSegWitCreds = PathCreds + "segwit/"

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

// InvalidModeError ..
const InvalidModeError = "Invalid mode"

// WalletNotFoundError ..
const WalletNotFoundError = "Wallet not found"

// MultiSigWalletNotFoundError ..
const MultiSigWalletNotFoundError = "Multisig wallet not found"

// SegWitWalletNotFoundError ..
const SegWitWalletNotFoundError = "Native segwit wallet not found"

// WalletAlreadyExistsError ..
const WalletAlreadyExistsError = "BIP44 wallet with same name already exists"

// MultiSigWalletAlreadyExistsError ..
const MultiSigWalletAlreadyExistsError = "Multisig wallet with same name already exists"

// SegWitWalletAlreadyExistsError ..
const SegWitWalletAlreadyExistsError = "Native segwit wallet with same name already exists"

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

// UnknownWalletTypeError ..
const UnknownWalletTypeError = "Unknown wallet type"

// PathWalletsHelpSyn ..
const PathWalletsHelpSyn = "Creates a new BIP44 wallet by specifying network and name"

// PathWalletsHelpDesc ..
const PathWalletsHelpDesc = ""

// PathAddressHelpSyn ..
const PathAddressHelpSyn = "Returns a new receiving address for selected wallet. You must provide an authentication token"

// PathAddressHelpDesc ..
const PathAddressHelpDesc = ""

// PathCredsHelpSyn ..
const PathCredsHelpSyn = "Creates authorization tokens for a wallet"

// PathCredsHelpDesc ..
const PathCredsHelpDesc = ""

// PathMultiSigWalletsHelpSyn ..
const PathMultiSigWalletsHelpSyn = "Creates a new wallet that is used as nth key of the m/n multisig wallet"

// PathMultiSigWalletsHelpDesc ..
const PathMultiSigWalletsHelpDesc = ""

// PathMultiSigAddressHelpSyn ..
const PathMultiSigAddressHelpSyn = "Returns the receiving address of the selected multisig wallet"

// PathMultiSigAddressHelpDesc ..
const PathMultiSigAddressHelpDesc = ""

// PathMultiSigCredsHelpSyn ..
const PathMultiSigCredsHelpSyn = "Creates access tokens for a multisig wallet"

// PathMultiSigCredsHelpDesc ..
const PathMultiSigCredsHelpDesc = ""

// PathSegWitWalletsHelpSyn ..
const PathSegWitWalletsHelpSyn = "Creates a new BIP84 wallet by specifying network and name"

// PathSegWitWalletsHelpDesc ..
const PathSegWitWalletsHelpDesc = ""

// PathSegWitAddressHelpSyn ..
const PathSegWitAddressHelpSyn = "Returns a new receiving address for selected native segwit wallet. You must provide an authentication token"

// PathSegWitAddressHelpDesc ..
const PathSegWitAddressHelpDesc = ""

// PathSegWitCredsHelpSyn ..
const PathSegWitCredsHelpSyn = "Creates access tokens for a native segwit wallet"

// PathSegWitCredsHelpDesc ..
const PathSegWitCredsHelpDesc = ""

// PathTransactionHelpSyn ..
const PathTransactionHelpSyn = "Sign bitcoin raw transaction"

// PathTransactionHelpDesc ..
const PathTransactionHelpDesc = ""

// BackendHelp ..
const BackendHelp = "The bitcoin custodian plugin lets you to store your private keys securely.\nWith this plugin you can create wallet or multisig wallet and generate receiving addresses or sign transactions."

// SecretCredsType ..
const SecretCredsType = "creds"

// MultiSigSecretCredsType ..
const MultiSigSecretCredsType = "multisig_" + SecretCredsType

// SegWitSecretCredsType ..
const SegWitSecretCredsType = "segwit_" + SecretCredsType

// StandardType ..
const StandardType = 0

// MultiSigType ..
const MultiSigType = 1

// SegWitType ..
const SegWitType = 2

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

// SegwitPurpose ..
const SegwitPurpose = HardenedKeyStart + 49

// NativeSegwitPurpose ..
const NativeSegwitPurpose = HardenedKeyStart + 84

// CoinType ..
var CoinType = map[string]uint32{
	MainNet: HardenedKeyStart,
	TestNet: HardenedKeyStart + 1,
}

// Account ..
const Account = HardenedKeyStart

// Change ..
const Change = 0

// MultiSigDefaultAddressIndex ..
const MultiSigDefaultAddressIndex = 0

// MinMultiSigN ..
const MinMultiSigN = 1

// MaxMultiSigN ..
const MaxMultiSigN = 7
