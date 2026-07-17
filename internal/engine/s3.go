package engine

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"vessl.dev/vessl/internal/models"
)

type S3Error struct {
	Status  int
	Code    string
	Message string
}

func (e *S3Error) Error() string {
	return fmt.Sprintf("S3 request failed: %d %s", e.Status, e.Message)
}

func sha256Hex(data []byte) string {
	h := sha256.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

func hmacSha256(key []byte, data string) []byte {
	h := hmac.New(sha256.New, key)
	h.Write([]byte(data))
	return h.Sum(nil)
}

func hmacHex(key []byte, data string) string {
	return hex.EncodeToString(hmacSha256(key, data))
}

func getSigningKey(secret, date, region, service string) []byte {
	kDate := hmacSha256([]byte("AWS4"+secret), date)
	kRegion := hmacSha256(kDate, region)
	kService := hmacSha256(kRegion, service)
	kSigning := hmacSha256(kService, "aws4_request")
	return kSigning
}

func encodePathSegment(value string) string {
	var buf strings.Builder
	for i := 0; i < len(value); i++ {
		c := value[i]
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-' || c == '_' || c == '.' || c == '~' {
			buf.WriteByte(c)
		} else {
			buf.WriteString(fmt.Sprintf("%%%02X", c))
		}
	}
	return buf.String()
}

func canonicalPath(bucket string, key string) string {
	bucketPath := encodePathSegment(bucket)
	if key == "" {
		return "/" + bucketPath
	}
	var parts []string
	for _, p := range strings.Split(key, "/") {
		parts = append(parts, encodePathSegment(p))
	}
	return "/" + bucketPath + "/" + strings.Join(parts, "/")
}

// EnsureS3Bucket implements HEAD bucket and PUT bucket if missing
func EnsureS3Bucket(ctx context.Context, dest *models.S3Destination) error {
	resp, err := signedS3Request(ctx, dest, http.MethodHead, "", nil, "")
	if err != nil {
		s3Err, ok := err.(*S3Error)
		if ok && s3Err.Status == 404 {
			// Bucket not found, create it
			respPut, errPut := signedS3Request(ctx, dest, http.MethodPut, "", nil, "")
			if errPut != nil {
				return fmt.Errorf("failed to create bucket: %w", errPut)
			}
			respPut.Body.Close()
			return nil
		}
		return fmt.Errorf("failed to check if bucket exists: %w", err)
	}
	resp.Body.Close()
	return nil
}

func signedS3Request(ctx context.Context, dest *models.S3Destination, method, key string, body []byte, contentType string) (*http.Response, error) {
	if body == nil {
		body = []byte{}
	}
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	payloadHash := sha256Hex(body)
	now := time.Now().UTC()
	amzDateFull := now.Format("20060102T150405Z")
	amzDateShort := now.Format("20060102")

	endpoint := strings.TrimPrefix(dest.Endpoint, "https://")
	endpoint = strings.TrimPrefix(endpoint, "http://")
	endpoint = strings.TrimRight(endpoint, "/")

	region := dest.Region
	if region == "" {
		region = "auto"
	}

	host := endpoint
	pathStr := canonicalPath(dest.Bucket, key)

	signedHeaders := "content-type;host;x-amz-content-sha256;x-amz-date"
	canonicalHeaders := fmt.Sprintf("content-type:%s\nhost:%s\nx-amz-content-sha256:%s\nx-amz-date:%s\n", contentType, host, payloadHash, amzDateFull)

	canonicalRequest := fmt.Sprintf("%s\n%s\n\n%s\n%s\n%s", method, pathStr, canonicalHeaders, signedHeaders, payloadHash)

	scope := fmt.Sprintf("%s/%s/s3/aws4_request", amzDateShort, region)
	stringToSign := fmt.Sprintf("AWS4-HMAC-SHA256\n%s\n%s\n%s", amzDateFull, scope, sha256Hex([]byte(canonicalRequest)))

	signingKey := getSigningKey(dest.SecretAccessKey, amzDateShort, region, "s3")
	signature := hmacHex(signingKey, stringToSign)

	urlStr := dest.Endpoint
	if !strings.HasPrefix(urlStr, "http") {
		urlStr = "https://" + urlStr
	}
	urlStr = strings.TrimRight(urlStr, "/") + pathStr

	req, err := http.NewRequestWithContext(ctx, method, urlStr, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("x-amz-content-sha256", payloadHash)
	req.Header.Set("x-amz-date", amzDateFull)

	authHeader := fmt.Sprintf("AWS4-HMAC-SHA256 Credential=%s/%s, SignedHeaders=%s, Signature=%s", dest.AccessKeyID, scope, signedHeaders, signature)
	req.Header.Set("Authorization", authHeader)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		b, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, &S3Error{Status: resp.StatusCode, Message: string(b)}
	}

	return resp, nil
}
