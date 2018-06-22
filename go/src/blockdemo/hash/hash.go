package hash

import (
	sha256 "crypto/sha256"
	"encoding/hex"
	"log"
)

func calculateHash(toBeHashed string) string  {
	hashInBytes :=sha256.Sum256([]byte(toBeHashed))
	hashInStr := hex.EncodeToString(hashInBytes[:])
	log.Printf("%s","%s",toBeHashed,hashInStr)
	return hashInStr
}

func main()  {
	calculateHash("test1")
}
