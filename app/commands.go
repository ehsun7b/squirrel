package app

type Command func(...string)
type Printer func(string, ...interface{})
type Encryptor func(string) (string, error)
