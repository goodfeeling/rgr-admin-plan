package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

// CalculateFileMD5
func CalculateFileMD5(filePath string) (string, error) {
	// open file
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("can't open file: %v", err)
	}
	defer file.Close()

	// create md5 hash object
	hash := md5.New()

	// copy hash object to file
	_, err = io.Copy(hash, file)
	if err != nil {
		return "", fmt.Errorf("read file context error: %v", err)
	}

	// generate hash
	md5Sum := hex.EncodeToString(hash.Sum(nil))
	return md5Sum, nil
}
