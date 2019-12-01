package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/labbcb/wf/models"
)

type Client struct {
	Host string
}

func (c *Client) Submit(workflow string, inputs []string, imports []string, options string) (*models.IdAndStatus, error) {
	var b bytes.Buffer
	bw := multipart.NewWriter(&b)

	// upload workflow file
	if err := uploadFile(bw, "workflowSource", workflow); err != nil {
		return nil, err
	}

	// upload inputs files
	if inputs != nil {
		fieldname := "workflowInputs"
		for idx, filename := range inputs {
			if idx > 0 {
				fieldname = "workflowInputs_" + strconv.Itoa(idx + 1)
			}
			if err := uploadFile(bw, fieldname, filename); err != nil {
				return nil, err
			}
		}
	}

	// zip workflow dependencies and upload
	if imports != nil {
		zip, err := ioutil.TempFile("", "*.zip")
		if err != nil {
			return nil, err
		}
		if err := zipFiles(zip, imports); err != nil {
			return nil, err
		}
		if err := zip.Close(); err != nil {
			return nil, err
		}
		if err := uploadFile(bw, "workflowDependencies", zip.Name()); err != nil {
			return nil, err
		}
		if err := os.Remove(zip.Name()); err != nil {
			return nil, err
		}
	}

	// upload options file
	if options != "" {
		if err := uploadFile(bw, options, "workflowOptions"); err != nil {
			return nil, err
		}
	}

	// create content-type
	contentType := bw.FormDataContentType()
	if err := bw.Close(); err != nil {
		return nil, err
	}

	// submit workflow
	url := c.Host + "/api/workflows/v1"
	resp, err := http.Post(url, contentType, &b)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// check HTTP satus code
	if resp.StatusCode != http.StatusCreated {
		return nil, raiseHTTPError(resp)
	}

	// unmarshal JSON
	var res models.IdAndStatus
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) Abort(id string) (*models.IdAndStatus, error) {
	url := c.Host + "/api/workflows/v1/" + id + "/abort"
	resp, err := http.Post(url, "application/json", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, raiseHTTPError(resp)
	}

	// unmarshal JSON
	var res models.IdAndStatus
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) Status(id string) (*models.IdAndStatus, error) {
	url := c.Host + "/api/workflows/v1/" + id + "/status"
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, raiseHTTPError(resp)
	}

	var res models.IdAndStatus
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) Query(params []*models.WorkflowQueryParameter) (*models.WorkflowQueryResponse, error) {
	var b bytes.Buffer
	if err := json.NewEncoder(&b).Encode(params); err != nil {
		return nil, err
	}

	url := c.Host + "/api/workflows/v1/query"
	resp, err := http.Post(url, "application/json", &b)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, raiseHTTPError(resp)
	}

	var res models.WorkflowQueryResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) Describe(workflow, inputs string) (*models.WorkflowDescription, error) {
	var b bytes.Buffer
	bw := multipart.NewWriter(&b)

	// upload workflow file
	if err := uploadFile(bw, "workflowSource", workflow); err != nil {
		return nil, err
	}

	// upload inputs files
	if inputs != "" {
		if err := uploadFile(bw, "workflowInputs", inputs); err != nil {
			return nil, err
		}
	}

	// create content-type
	contentType := bw.FormDataContentType()
	if err := bw.Close(); err != nil {
		return nil, err
	}

	url := c.Host + "/api/womtool/v1/describe"
	resp, err := http.Post(url, contentType, &b)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, raiseHTTPError(resp)
	}

	var res models.WorkflowDescription
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) Logs(id string) (string, error) {
	url := c.Host + "/api/workflows/v1/" + id + "/logs"
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", raiseHTTPError(resp)
	}

	var b bytes.Buffer
	_, err = b.ReadFrom(resp.Body)
	if err != nil {
		return "", err
	}
	return b.String(), nil
}

func (c *Client) Outputs(id string) (string, error) {
	url := c.Host + "/api/workflows/v1/" + id + "/outputs"
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", raiseHTTPError(resp)
	}

	var b bytes.Buffer
	_, err = b.ReadFrom(resp.Body)
	if err != nil {
		return "", err
	}
	return b.String(), nil
}

// write file content to multipart writer
func uploadFile(w *multipart.Writer, fieldname, filename string) error {
	fh, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fh.Close()

	fw, err := w.CreateFormFile(fieldname, filepath.Base(filename))
	if err != nil {
		return err
	}

	_, err = io.Copy(fw, fh)
	return err
}

// new error with 'HTTP Status (Status Code): Body'
// it doesn't close resp.Body reader
func raiseHTTPError(resp *http.Response) error {
	var b bytes.Buffer
	_, err := b.ReadFrom(resp.Body)
	if err != nil {
		return err
	}
	return errors.New(fmt.Sprintf("%s: %s", resp.Status, b.String()))
}
