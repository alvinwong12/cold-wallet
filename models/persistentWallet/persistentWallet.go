package persistentWallet

// Refactor with generic
type PersistentWallet interface {
	ExportWalletToFile(file string)
}
