package btc

// SignRawTransactionResult 签名结果
type SignRawTransactionResult struct {
	Hex      string
	Changed  bool //is raw tx changed
	Complete bool //multi input,multisig require meetted
	Errors   *string
}
