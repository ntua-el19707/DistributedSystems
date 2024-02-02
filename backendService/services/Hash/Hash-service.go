package Hasher

import (
	"Logger"
	"Service"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

type HashService interface {
	Service.Service
	Hash(previous, currentHash string) (error, string)
	Valid(previous, currentHash, expected string) error
	Seed(currentHash string) (int64, error)
	ParrentOFall() string
	InstantHash(seed int64) string
}
type HashImpl struct {
	loggerService Logger.LoggerService
}

func (service *HashImpl) Construct() error {
	if service.loggerService == nil {
		service.loggerService = &Logger.Logger{ServiceName: "HasherService"}
		err := service.loggerService.Construct()
		if err != nil {
			return err
		}
	}
	service.loggerService.Log("Service  created")
	return nil
}
func (service HashImpl) Hash(previous, currentHash string) (error, string) {
	service.loggerService.Log("Start  hashing next  block")
	xor, err := performXorString(previous, currentHash, service.loggerService)
	if err != nil {
		return err, ""
	}
	times, err := perfomAndAndGetSum(previous, currentHash, service.loggerService)
	if err != nil {
		return err, ""
	}
	hashed := hasher(xor, times, service.loggerService)
	service.loggerService.Log("Commit  hashing next  block")
	return nil, hashed
}
func (service HashImpl) Valid(previous, currentHash, expected string) error {
	err, hash := service.Hash(previous, currentHash)
	if err != nil {
		return err
	}
	if hash == expected {
		return nil
	}
	return errors.New("Hashes  do not match")
}
func (service HashImpl) Seed(hash string) (int64, error) {
	var seed int64
	bytes := []byte(hash)
	if len(bytes) >= 8 {
		seed = int64(binary.LittleEndian.Uint64(bytes[:8]))
	} else {
		return 0, errors.New("Error: Not enough bytes to generate a seed.")

	}
	return seed, nil
}
func (service HashImpl) ParrentOFall() string {
	var parent string
	for i := 0; i < 64; i++ {
		parent += "1"
	}
	return parent
}
func (service HashImpl) InstantHash(seed int64) string {
	timeInt := time.Now().Unix() + seed
	hash := hasher(fmt.Sprint(timeInt), 1, service.loggerService)
	return hash
}

//-- USE - FULL

func performXorString(str1, str2 string, loggerService Logger.LoggerService) (string, error) {
	loggerService.Log(fmt.Sprintf("Start performing XOR %s ,  %s ", str1, str2))
	if len(str1) != len(str2) {
		errmsg := fmt.Sprintf("Abbort performing XOR %s ,  %s string  does not  have  the same size ", str1, str2)
		loggerService.Error(errmsg)
		return "", errors.New(errmsg)
	}
	result := make([]byte, len(str1))

	for i := 0; i < len(str1); i++ {
		result[i] = str1[i] ^ str2[i]
	}
	strResult := string(result)
	loggerService.Log(fmt.Sprintf("Commit performing XOR %s ,  %s ", str1, str2))
	return strResult, nil
}

func perfomAndAndGetSum(str1, str2 string, loggerService Logger.LoggerService) (int, error) {
	loggerService.Log(fmt.Sprintf("Start performing and  to get sum %s ,  %s ", str1, str2))
	if len(str1) != len(str2) {
		errmsg := fmt.Sprintf("Abbort performing And  to get sum %s ,  %s string  does not  have  the same size ", str1, str2)
		loggerService.Error(errmsg)
		return 0, errors.New(errmsg)
	}
	result := make([]byte, len(str1))
	for i := 0; i < len(str1); i++ {
		result[i] = str1[i] & str2[i]
	}
	strResult := string(result)
	sum := 1 // for at lest on iteration unlikly to be 1
	for _, charcter := range strResult {
		sum += int(charcter)
	}
	loggerService.Log(fmt.Sprintf("Commit performing and  and get sum %s ,  %s ", str1, str2))
	return sum, nil

}
func hasher(str string, times int, logger Logger.LoggerService) string {
	logger.Log(fmt.Sprintf("Start  loop hasshing for %d", times))
	hash := sha256.New()
	hash.Write([]byte(str))
	for i := 0; i < times-1; i++ {
		bytes := hash.Sum(nil)
		hash = sha256.New()
		hash.Write(bytes)
	}
	logger.Log(fmt.Sprintf("Commit  loop hasshing for %d", times))
	return hex.EncodeToString(hash.Sum(nil))
}

// -- Mock Hasher --
type MockHasher struct {
	Hashvalue        string
	HashFailed       bool
	Invalid          bool
	InvalidError     string
	SeedVal          int64
	SeedFailed       bool
	InstantHashValue string
	CallHash         int
	CallParentOfAll  int
	CallValid        int
	CallSeed         int
	CallInstand      int
}

func (m *MockHasher) Construct() error {

	return nil
}
func (m *MockHasher) Hash(previous, currentHash string) (error, string) {

	m.CallHash++
	if m.HashFailed {
		return errors.New("has  faield "), ""
	}
	return nil, m.Hashvalue

}

func (m *MockHasher) Valid(previous, currentHash, expected string) error {
	m.CallValid++
	if m.Invalid {
		return errors.New(m.InvalidError)
	}
	return nil
}

func (m *MockHasher) Seed(currentHash string) (int64, error) {
	m.CallSeed++
	if m.SeedFailed {
		return 0, errors.New("seed failed ")
	}
	return m.SeedVal, nil
}

func (m *MockHasher) ParrentOFall() string {
	m.CallParentOfAll++
	return "1"
}

func (m *MockHasher) InstantHash(seed int64) string {
	m.CallInstand++
	return m.InstantHashValue
}
