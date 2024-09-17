package types

type Printer func(string, ...interface{})
type Encryptor func(string) (string, error)
type Decryptor func(string) (string, error)
