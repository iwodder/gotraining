package secrets

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"strings"
)

type Secrets struct {
	key, path string
}

func New(key, path string) *Secrets {
	_, err := os.Stat(path)
	if errors.Is(err, fs.ErrNotExist) {
		log.Printf("file doesn't exists at %s...creating", path)
		f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0o600)
		if err != nil {
			panic(err)
		}
		f.Close()
	}
	return &Secrets{
		key: key, path: path,
	}
}

func (s *Secrets) Store(key, value string) error {
	f, err := os.OpenFile(s.path, os.O_RDWR, 0o600)
	if err != nil {
		return err
	}
	defer f.Close()
	contents, err := decrypt(s.key, f)
	if err != nil {
		return err
	}
	contents = s.update(contents, key, value)
	err = f.Truncate(0)
	if err != nil {
		return err
	}
	_, err = f.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	err = encrypt(s.key, contents, f)
	if err != nil {
		return err
	}
	return nil
}

func (s *Secrets) Load(key string) string {
	f, err := os.Open(s.path)
	if err != nil {
		return ""
	}
	defer f.Close()
	contents, err := decrypt(s.key, f)
	if err != nil {
		panic(err)
	}
	contents = strings.Trim(contents, "\n")
	for _, line := range strings.Split(contents, "\n") {
		parts := strings.Split(line, ":")
		if parts[0] == key {
			return parts[1]
		}
	}
	return ""
}

func (s *Secrets) update(lines string, key, value string) string {
	found := false
	var sb strings.Builder
	lines = strings.Trim(lines, "\n")
	for _, line := range strings.Split(lines, "\n") {
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, key) {
			sb.WriteString(fmt.Sprintf("%s:%s\n", key, value))
			found = true
		} else {
			sb.WriteString(line)
			sb.WriteByte('\n')
		}
	}
	if !found {
		sb.WriteString(fmt.Sprintf("%s:%s\n", key, value))
	}
	return sb.String()
}

func encrypt(key string, value string, w io.Writer) error {
	hash := md5.Sum([]byte(key))
	c, err := aes.NewCipher(hash[:])
	if err != nil {
		return err
	}
	var iv [aes.BlockSize]byte
	stream := cipher.NewOFB(c, iv[:])

	sw := cipher.StreamWriter{S: stream, W: w}
	_, err = sw.Write([]byte(value))
	if err != nil {
		return err
	}
	return nil
}

func decrypt(key string, r io.Reader) (string, error) {
	hash := md5.Sum([]byte(key))
	c, err := aes.NewCipher(hash[:])
	if err != nil {
		return "", err
	}
	var iv [aes.BlockSize]byte
	stream := cipher.NewOFB(c, iv[:])

	sw := cipher.StreamReader{S: stream, R: r}

	var sb strings.Builder
	var buff [1024]byte
	for n, err := sw.Read(buff[:]); ; n, err = sw.Read(buff[:]) {
		if err != nil {
			if errors.Is(err, io.EOF) {
				sb.Write(buff[:n])
				return sb.String(), nil
			}
			return "", err
		}
		sb.Write(buff[:n])
	}
}
