package wallet

// Refactor with generic
type PersistentWallet interface {
	ExportWalletToFile(file string)
}
