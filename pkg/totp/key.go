package totp

import (
	"bytes"
	"encoding/base64"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"image/png"
	"net/url"
	"strconv"
	"strings"
)

type Key struct {
	orig string
	url  *url.URL
}

func NewKeyFromURL(orig string) (*Key, error) {
	s := strings.TrimSpace(orig)

	u, err := url.Parse(s)
	if err != nil {
		return nil, err
	}

	return &Key{
		orig: s,
		url:  u,
	}, nil
}

func (k *Key) String() string {
	return k.orig
}

func (k *Key) Image(width int, height int) (string, error) {
	b, err := qr.Encode(k.orig, qr.M, qr.Auto)
	if err != nil {
		return "", err
	}

	b, err = barcode.Scale(b, width, height)

	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err = png.Encode(&buf, b); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func (k *Key) Type() string {
	return k.url.Host
}

func (k *Key) Issuer() string {
	q := k.url.Query()

	issuer := q.Get("issuer")

	if issuer != "" {
		return issuer
	}

	p := strings.TrimPrefix(k.url.Path, "/")
	i := strings.Index(p, ":")

	if i == -1 {
		return ""
	}

	return p[:i]
}

func (k *Key) Account() string {
	p := strings.TrimPrefix(k.url.Path, "/")
	i := strings.Index(p, ":")

	if i == -1 {
		return p
	}

	return p[i+1:]
}

func (k *Key) Secret() string {
	q := k.url.Query()

	return q.Get("secret")
}

func (k *Key) Period() uint64 {
	q := k.url.Query()

	if u, err := strconv.ParseUint(q.Get("period"), 10, 64); err == nil {
		return u
	}

	return 30
}

func (k *Key) Digits() Digits {
	q := k.url.Query()

	if u, err := strconv.ParseUint(q.Get("digits"), 10, 64); err == nil {
		switch u {
		case 8:
			return DigitsEight
		default:
			return DigitsSix
		}
	}

	// Six is the most common value.
	return DigitsSix
}

func (k *Key) Algorithm() Algorithm {
	q := k.url.Query()

	a := strings.ToLower(q.Get("algorithm"))
	switch a {
	case "md5":
		return AlgorithmMD5
	case "sha256":
		return AlgorithmSHA256
	case "sha512":
		return AlgorithmSHA512
	default:
		return AlgorithmSHA1
	}
}

// URL returns the OTP URL as a string
func (k *Key) URL() string {
	return k.url.String()
}
