package tools

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/bwmarrin/snowflake"
)

// Get a random integer in a specified range
func RandomInRange(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

//Get unique identifier
func GetGuid(num int64) (int64, error) {
	node, err := snowflake.NewNode(num)
	if err != nil {
		return 0, err
	}

	id := node.Generate()
	return id.Int64(), nil
}

//  ----------------------- Random key -----------------------
const baseStr = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()[]{}+-*/_=."

// Generate random password
func GetRandomcode(length uint8) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano() + rand.Int63()))
	bytes := make([]byte, length)
	l := len(baseStr)
	for i := uint8(0); i < length; i++ {
		bytes[i] = baseStr[r.Intn(l)]
	}
	return string(bytes)
}

// Calculate the file hash value
func CalcFileHash(fpath string) (string, error) {
	f, err := os.Open(fpath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// Calculate MD5
func CalcMD5(s string) ([]byte, error) {
	h := md5.New()
	_, err := h.Write([]byte(s))
	if err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

// Int64 to Bytes
func Int64ToBytes(n int64) []byte {
	bytebuf := bytes.NewBuffer([]byte{})
	binary.Write(bytebuf, binary.BigEndian, n)
	return bytebuf.Bytes()
}

// Bytes to Int64
func BytesToInt64(bys []byte) int64 {
	bytebuff := bytes.NewBuffer(bys)
	var data int64
	binary.Read(bytebuff, binary.BigEndian, &data)
	return data
}

var reg_mail = regexp.MustCompile(`^[0-9a-z][_,0-9a-z-]{0,31}@([0-9a-z][0-9a-z-]{0,30}[0-9a-z]\.){1,4}[a-z]{2,4}$`)

//
func VerifyMailboxFormat(mailbox string) bool {
	return reg_mail.MatchString(mailbox)
}

// Get all files in dir
func WalkDir(dir string) ([]string, error) {
	files := make([]string, 0)
	fs, err := ioutil.ReadDir(dir)
	if err != nil {
		return files, err
	} else {
		for _, v := range fs {
			if !v.IsDir() {
				files = append(files, v.Name())
			}
		}
	}
	return files, nil
}

//Get the number of non-blank lines in a file
func GetFileNonblankLine(path string) (int, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	count := 0
	defer file.Close()
	buffer := bufio.NewReader(file)
	for {
		ctx, _, err := buffer.ReadLine()
		if err != nil {
			return count, nil
		}
		if strings.TrimSpace(string(ctx)) == "" {
			continue
		}
		count++
	}
}

// Create a directory
func CreatDirIfNotExist(dir string) error {
	_, err := os.Stat(dir)
	if err != nil {
		return os.MkdirAll(dir, os.ModeDir)
	}
	return nil
}

func CalcHash(data []byte) (string, error) {
	if len(data) <= 0 {
		return "", errors.New("data is nil")
	}
	h := sha256.New()
	_, err := h.Write(data)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// Write string content to file
func WriteStringtoFile(content, fileName string) error {
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(content)
	if err != nil {
		return err
	}
	return nil
}
