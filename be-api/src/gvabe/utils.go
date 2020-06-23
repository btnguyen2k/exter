package gvabe

import (
	"bytes"
	"compress/zlib"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"io"
	"math"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/btnguyen2k/consu/reddo"
	"github.com/btnguyen2k/consu/semita"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"

	"main/src/gvabe/bo/app"
	"main/src/gvabe/bo/user"
	"main/src/mico"
)

const (
	frontendAppId       = "exter_fe"
	frontendAppIdPrefix = frontendAppId + ":"

	systemAppId   = "exter"
	systemAppDesc = "Exter"
)

const (
	apiResultExtraAccessToken = "_access_token_"

	loginSessionTtl        = 3600 * 8
	loginSessionNearExpiry = 3600 * 3
)

var (
	systemAdminId        string
	enabledLoginChannels = make(map[string]bool)

	appDao  app.AppDao
	userDao user.UserDao

	rsaPrivKey *rsa.PrivateKey
	rsaPubKey  *rsa.PublicKey

	sessionCache         mico.ICache
	preLoginSessionCache mico.ICache
)

const (
	loginChannelGoogle   = "google"
	loginChannelFacebook = "facebook"
)

func genRsaKey(numBits int) (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, numBits)
}

func parseRsaPublicKeyFromPem(pemStr string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		return nil, errors.New("failed to parse PEM block")
	}

	if pub, err := x509.ParsePKIXPublicKey(block.Bytes); err != nil {
		return nil, err
	} else {
		switch pub := pub.(type) {
		case *rsa.PublicKey:
			return pub, nil
		default:
			return nil, errors.New("not RSA public key")
		}
	}
}

func encryptPassword(username, rawPassword string) string {
	saltAndPwd := username + "." + rawPassword
	out := sha1.Sum([]byte(saltAndPwd))
	return strings.ToLower(hex.EncodeToString(out[:]))
}

// padRight adds "0" right right of a string until its length reach a specific value.
func padRight(str string, l int) string {
	for len(str) < l {
		str += "0"
	}
	return str
}

// aesEncrypt encrypts a block of data using AES/CTR mode.
//
// IV is put at the beginning of the cipher data.
func aesEncrypt(key, data []byte) ([]byte, error) {
	for len(key) < 16 {
		key = append(key, 0)
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	iv := []byte(padRight(strconv.FormatInt(time.Now().UnixNano(), 16), 16))
	cipherData := make([]byte, 16+len(data))
	copy(cipherData, iv)
	ctr := cipher.NewCTR(block, iv)
	ctr.XORKeyStream(cipherData[16:], data)
	return cipherData, nil
}

// aesDecrypt decrypts a block of encrypted data using AES/CTR mode.
//
// Assuming IV is put at the beginning of the cipher data.
func aesDecrypt(key, encryptedData []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	iv := encryptedData[0:16]
	data := make([]byte, len(encryptedData)-16)
	ctr := cipher.NewCTR(block, iv)
	ctr.XORKeyStream(data, encryptedData[16:])
	return data, nil
}

// zlibCompress compresses data using zlib.
func zlibCompress(data []byte) []byte {
	var b bytes.Buffer
	w, _ := zlib.NewWriterLevel(&b, zlib.BestCompression)
	w.Write(data)
	w.Close()
	return b.Bytes()
}

// zlibDecompress decompressed compressed-data using zlib.
func zlibDecompress(compressedData []byte) ([]byte, error) {
	r, err := zlib.NewReader(bytes.NewReader(compressedData))
	if err != nil {
		return nil, err
	}
	var b bytes.Buffer
	_, err = io.Copy(&b, r)
	r.Close()
	return b.Bytes(), err
}

func zipAndEncrypt(data, aesKey []byte) ([]byte, error) {
	zip := zlibCompress(data)
	return aesEncrypt(aesKey, zip)
}
func decryptAndUnzip(encdata, aesKey []byte) ([]byte, error) {
	if zip, err := aesDecrypt(aesKey, encdata); err != nil {
		return nil, err
	} else {
		return zlibDecompress(zip)
	}
}

/*----------------------------------------------------------------------*/

var muxSytemInfo sync.Mutex
var systemInfoArr = make([]map[string]interface{}, 0)

func startUpdateSystemInfo() {
	for {
		go doUpdateSystemInfo()
		<-time.After(10 * time.Second)
	}
}

func lastSystemInfo() map[string]interface{} {
	muxSytemInfo.Lock()
	defer muxSytemInfo.Unlock()
	return systemInfoArr[len(systemInfoArr)-1]
}

func doUpdateSystemInfo() {
	muxSytemInfo.Lock()
	defer muxSytemInfo.Unlock()

	data := make(map[string]interface{})
	{
		load, err := load.Avg()
		cpuLoad := -1.0
		if err == nil && load != nil {
			cpuLoad = math.Floor(load.Load1*100) / 100
		}
		historyLoad := make([]float64, 0)
		for _, data := range systemInfoArr {
			s := semita.NewSemita(data)
			load, _ := s.GetValueOfType("cpu.load", reddo.TypeFloat)
			historyLoad = append(historyLoad, load.(float64))
		}
		data["cpu"] = map[string]interface{}{
			"cores":        runtime.NumCPU(),
			"load":         cpuLoad,
			"history_load": historyLoad,
		}
	}

	{
		mem, err := mem.VirtualMemory()
		memFree := uint64(0)
		if err == nil && mem != nil {
			memFree = mem.Free
		}
		historyFree := make([]uint64, 0)
		historyFreeKb := make([]float64, 0)
		historyFreeMb := make([]float64, 0)
		historyFreeGb := make([]float64, 0)
		for _, data := range systemInfoArr {
			s := semita.NewSemita(data)
			free, _ := s.GetValueOfType("memory.free", reddo.TypeUint)
			historyFree = append(historyFree, free.(uint64))
			freeKb, _ := s.GetValueOfType("memory.freeKb", reddo.TypeFloat)
			historyFreeKb = append(historyFreeKb, freeKb.(float64))
			freeMb, _ := s.GetValueOfType("memory.freeMb", reddo.TypeFloat)
			historyFreeMb = append(historyFreeMb, freeMb.(float64))
			freeGb, _ := s.GetValueOfType("memory.freeGb", reddo.TypeFloat)
			historyFreeGb = append(historyFreeGb, freeGb.(float64))
		}
		data["memory"] = map[string]interface{}{
			"free":           memFree,
			"freeKb":         math.Floor(100.0*(float64(memFree)/1024)) / 100.0,
			"freeMb":         math.Floor(100.0*(float64(memFree)/1024/1024)) / 100.0,
			"freeGb":         math.Floor(100.0*(float64(memFree)/1024/1024/1024)) / 100.0,
			"history_free":   historyFree,
			"history_freeKb": historyFreeKb,
			"history_freeMb": historyFreeMb,
			"history_freeGb": historyFreeGb,
		}
	}

	{
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		historyUsed := make([]uint64, 0)
		historyUsedKb := make([]float64, 0)
		historyUsedMb := make([]float64, 0)
		historyUsedGb := make([]float64, 0)
		for _, data := range systemInfoArr {
			s := semita.NewSemita(data)
			used, _ := s.GetValueOfType("app_memory.used", reddo.TypeUint)
			historyUsed = append(historyUsed, used.(uint64))
			usedKb, _ := s.GetValueOfType("app_memory.usedKb", reddo.TypeFloat)
			historyUsedKb = append(historyUsedKb, usedKb.(float64))
			usedMb, _ := s.GetValueOfType("app_memory.usedMb", reddo.TypeFloat)
			historyUsedMb = append(historyUsedMb, usedMb.(float64))
			usedGb, _ := s.GetValueOfType("app_memory.usedGb", reddo.TypeFloat)
			historyUsedGb = append(historyUsedGb, usedGb.(float64))
		}
		data["app_memory"] = map[string]interface{}{
			"used":           m.Alloc,
			"usedKb":         math.Floor(100.0*(float64(m.Alloc)/1024)) / 100.0,
			"usedMb":         math.Floor(100.0*(float64(m.Alloc)/1024/1024)) / 100.0,
			"usedGb":         math.Floor(100.0*(float64(m.Alloc)/1024/1024/1024)) / 100.0,
			"history_used":   historyUsed,
			"history_usedKb": historyUsedKb,
			"history_usedMb": historyUsedMb,
			"history_usedGb": historyUsedGb,
		}
	}

	{
		history := make([]int, 0)
		for _, data := range systemInfoArr {
			s := semita.NewSemita(data)
			n, _ := s.GetValueOfType("go_routines.num", reddo.TypeInt)
			history = append(history, int(n.(int64)))
		}
		data["go_routines"] = map[string]interface{}{
			"num":     runtime.NumGoroutine(),
			"history": history,
		}
	}

	systemInfoArr = append(systemInfoArr, data)
	if len(systemInfoArr) > 10 {
		systemInfoArr[0] = nil
		systemInfoArr = systemInfoArr[1:]
	}
}

var httpClient = &http.Client{
	Timeout: time.Second * 30,
}
