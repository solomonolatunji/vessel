package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	nethttp "net/http"
	"path/filepath"
)

type ComposeAnalyzeRequest struct {
	ProjectID      string `json:"projectId"`
	ComposeContent string `json:"composeContent"`
}

type ComposeAnalyzeResponse struct {
	AppServices []map[string]interface{} `json:"appServices"`
	Databases   []map[string]interface{} `json:"databases"`
}

func (c *Client) AnalyzeCompose(req *ComposeAnalyzeRequest) (*ComposeAnalyzeResponse, error) {
	resp, err := c.sendRequest("POST", "/compose/analyze", req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != nethttp.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to analyze compose (status %d): %s", resp.StatusCode, string(body))
	}

	var wrapper struct {
		Data *ComposeAnalyzeResponse `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return wrapper.Data, nil
}

func (c *Client) DeployCompose(projectID string, composeContent []byte, filename string) (int, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	if projectID != "" {
		if err := writer.WriteField("projectId", projectID); err != nil {
			return 0, err
		}
	}

	part, err := writer.CreateFormFile("file", filepath.Base(filename))
	if err != nil {
		return 0, err
	}
	if _, err := io.Copy(part, bytes.NewReader(composeContent)); err != nil {
		return 0, err
	}

	if err := writer.Close(); err != nil {
		return 0, err
	}

	req, err := nethttp.NewRequest("POST", fmt.Sprintf("%s/api/v1/compose/deploy", c.BaseURL), &body)
	if err != nil {
		return 0, err
	}

	if c.Token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != nethttp.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return 0, fmt.Errorf("failed to deploy compose (status %d): %s", resp.StatusCode, string(body))
	}

	var wrapper struct {
		Data struct {
			Count int `json:"count"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
		return 0, fmt.Errorf("failed to decode response: %w", err)
	}

	return wrapper.Data.Count, nil
}
