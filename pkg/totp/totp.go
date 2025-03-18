package totp

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base32"
	"encoding/binary"
	"math"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	b32NoPadding = base32.StdEncoding.WithPadding(base32.NoPadding)
	totp         *TOTP
	once         sync.Once
)

type TOTP struct {
	cfg *Config
}

type Params struct {
	TOTPEnabled bool   `json:"totp_enabled"`
	TOTPSecret  string `json:"totp_secret,omitempty"`
	TOTPURL     string `json:"totp_url,omitempty"`
}

func NewTOTP(cfg *Config) *TOTP {
	once.Do(func() {
		totp = &TOTP{
			cfg: cfg,
		}
	})

	return totp
}

func (totp *TOTP) Generate(opts *GenerateOptions) (*Key, error) {
	if opts.username == "" {
		return nil, ErrGenerateMissingUsername
	}

	v := url.Values{}
	secret := make([]byte, totp.cfg.SecretSize)
	if _, err := rand.Reader.Read(secret); err != nil {
		return nil, err
	}

	v.Set("secret", b32NoPadding.EncodeToString(secret))
	v.Set("issuer", totp.cfg.Issuer)
	v.Set("period", strconv.FormatUint(uint64(totp.cfg.Period), 10))
	v.Set("algorithm", totp.cfg.Algorithm.String())
	v.Set("digits", totp.cfg.Digits.String())

	u := url.URL{
		Scheme:   "otpauth",
		Host:     "totp",
		Path:     "/" + totp.cfg.Issuer + ":" + opts.username,
		RawQuery: encodeQuery(v),
	}

	return NewKeyFromURL(u.String())
}

func (totp *TOTP) Validate(passcode string, secret string) (bool, error) {
	counter := uint64(float64(time.Now().Unix()) / float64(totp.cfg.Period))
	passcode = strings.TrimSpace(passcode)

	if len(passcode) != totp.cfg.Digits.Length() {
		return false, ErrValidateInputInvalidLength
	}

	otp, err := totp.generateCode(secret, counter)
	if err != nil {
		return false, err
	}

	if subtle.ConstantTimeCompare([]byte(otp), []byte(passcode)) == 1 {
		return true, nil
	}

	return false, nil
}

func T() *TOTP {
	return totp
}

func (totp *TOTP) generateCode(secret string, counter uint64) (passcode string, err error) {
	secret = strings.TrimSpace(secret)
	if n := len(secret) % 8; n != 0 {
		secret = secret + strings.Repeat("=", 8-n)
	}

	secret = strings.ToUpper(secret)
	secretBytes, err := base32.StdEncoding.DecodeString(secret)
	if err != nil {
		return "", ErrValidateSecretInvalidBase32
	}

	buf := make([]byte, 8)
	mac := hmac.New(totp.cfg.Algorithm.Hash, secretBytes)
	binary.BigEndian.PutUint64(buf, counter)

	mac.Write(buf)
	sum := mac.Sum(nil)
	offset := sum[len(sum)-1] & 0xf
	value := int64(((int(sum[offset]) & 0x7f) << 24) |
		((int(sum[offset+1] & 0xff)) << 16) |
		((int(sum[offset+2] & 0xff)) << 8) |
		(int(sum[offset+3]) & 0xff))

	l := totp.cfg.Digits.Length()
	mod := int32(value % int64(math.Pow10(l)))

	return totp.cfg.Digits.Format(mod), nil
}
