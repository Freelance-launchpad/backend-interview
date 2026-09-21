package gotenberg

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"

	"github.com/Freelance-launchpad/backend-interview/common/jerror/v2"
	"github.com/Freelance-launchpad/backend-interview/common/jhttp/v2"
)

type Client interface {
	// HTMLToPDF returns a pdf file based on html files.
	HTMLToPDF(ctx context.Context, content io.Reader, footer io.Reader) (int64, io.ReadSeeker, error)
	HTMLToFacturXPDF(ctx context.Context, content io.Reader, footer io.Reader, facturx io.Reader) (int64, io.ReadSeeker, error)
	MergePDFs(ctx context.Context, files []io.Reader) (file io.ReadSeeker, length int64, err error)
	// Add a stamp on top of the provided file.
	// Stamp must be a PNG file.
	StampPDF(ctx context.Context, file io.Reader, stamp io.Reader) (int64, io.ReadSeeker, error)
}

var ErrNilContent = errors.New("content cannot be nil")

type client struct {
	jhttp.Doer
	url url.URL
}

func New(httpClient jhttp.Doer, url url.URL) *client {
	return &client{
		Doer: httpClient,
		url:  url,
	}
}

func prepareParts(content io.Reader, footer io.Reader, facturx io.Reader, body *bytes.Buffer) (*multipart.Writer, error) {
	writer := multipart.NewWriter(body)

	if content == nil {
		return nil, ErrNilContent
	}
	files := map[string]io.Reader{"index.html": content}
	if footer != nil {
		files["footer.html"] = footer
	}

	for name, data := range files {
		part, err := writer.CreateFormFile("files", name)
		if err != nil {
			return nil, err
		}
		content, err := io.ReadAll(data)
		if err != nil {
			return nil, err
		}
		_, err = part.Write(content)
		if err != nil {
			return nil, err
		}
	}

	if facturx != nil {
		if err := prepareFacturX(writer, facturx); err != nil {
			return nil, err
		}
	}

	widthPart, err := writer.CreateFormField("paperWidth")
	if err != nil {
		return nil, err
	}
	if _, err := widthPart.Write([]byte("210mm")); err != nil {
		return nil, err
	}

	heightPart, err := writer.CreateFormField("paperHeight")
	if err != nil {
		return nil, err
	}
	if _, err := heightPart.Write([]byte("297mm")); err != nil {
		return nil, err
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}

	return writer, nil
}

func prepareFacturX(multi *multipart.Writer, facturx io.Reader) error {
	xml, err := multi.CreateFormFile("facturxXml", "factur-x.xml")
	if err != nil {
		return err
	}
	content, err := io.ReadAll(facturx)
	if err != nil {
		return err
	}
	_, err = xml.Write(content)
	if err != nil {
		return err
	}

	profile, err := multi.CreateFormField("facturxConformanceLevel")
	if err != nil {
		return err
	}

	_, err = profile.Write([]byte("EN 16931"))
	if err != nil {
		return err
	}
	return nil
}

const pathHTMLToPDF = "/forms/chromium/convert/html"

func (c *client) HTMLToPDF(ctx context.Context, content io.Reader, footer io.Reader) (int64, io.ReadSeeker, error) {
	body := new(bytes.Buffer)

	writer, err := prepareParts(content, footer, nil, body)
	if err != nil {
		return 0, nil, err
	}

	return c.generatePDF(ctx, body, writer)
}

func (c *client) HTMLToFacturXPDF(ctx context.Context, content io.Reader, footer io.Reader, facturx io.Reader) (int64, io.ReadSeeker, error) {
	body := new(bytes.Buffer)

	writer, err := prepareParts(content, footer, facturx, body)
	if err != nil {
		return 0, nil, err
	}

	return c.generatePDF(ctx, body, writer)
}

func (c *client) generatePDF(ctx context.Context, body io.Reader, writer *multipart.Writer) (_ int64, _ io.ReadSeeker, err error) {
	defer jerror.Wrap(&err)

	uri := c.url
	uri.Path = pathHTMLToPDF

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, uri.String(), body)
	if err != nil {
		return 0, nil, err
	}

	request.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.Do(request)
	if err != nil {
		return 0, nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			responseBody = fmt.Appendf(nil, "error reading response body: %v", err)
		}
		return 0, nil, jhttp.APIError{
			Err:        fmt.Errorf("bad status, expected: %d, received: %d - with body %s", http.StatusOK, resp.StatusCode, string(responseBody)),
			StatusCode: resp.StatusCode,
		}
	}

	pdf, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, err
	}

	return resp.ContentLength, bytes.NewReader(pdf), nil
}

func (c *client) MergePDFs(ctx context.Context, files []io.Reader) (file io.ReadSeeker, length int64, err error) {
	defer jerror.Wrap(&err)

	uri := c.url
	uri.Path = "/forms/pdfengines/merge"

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for i, data := range files {
		part, err := writer.CreateFormFile("files", fmt.Sprintf("file-%02d.pdf", i))
		if err != nil {
			return nil, 0, err
		}

		content, err := io.ReadAll(data)
		if err != nil {
			return nil, 0, err
		}

		if _, err := part.Write(content); err != nil {
			return nil, 0, err
		}
	}

	if err := writer.Close(); err != nil {
		return nil, 0, err
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, uri.String(), bytes.NewReader(body.Bytes()))
	if err != nil {
		return nil, 0, err
	}
	request.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.Do(request)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			responseBody = fmt.Appendf(nil, "error reading response body: %v", err)
		}

		return nil, 0, jhttp.APIError{
			Err:        fmt.Errorf("bad status, expected: %d, received: %d - with body %s", http.StatusOK, resp.StatusCode, string(responseBody)),
			StatusCode: resp.StatusCode,
		}
	}

	pdf, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}

	return bytes.NewReader(pdf), resp.ContentLength, nil
}

func (c *client) StampPDF(ctx context.Context, file io.Reader, stamp io.Reader) (_ int64, _ io.ReadSeeker, err error) {
	defer jerror.Wrap(&err)

	if file == nil {
		return 0, nil, ErrNilContent
	}
	if stamp == nil {
		return 0, nil, fmt.Errorf("stamp cannot be nil")
	}

	uri := c.url
	uri.Path = "/forms/pdfengines/stamp"

	body := &bytes.Buffer{}

	writer := multipart.NewWriter(body)

	content, err := writer.CreateFormFile("files", "file.pdf")
	if err != nil {
		return 0, nil, err
	}
	_, err = io.Copy(content, file)
	if err != nil {
		return 0, nil, err
	}

	stampFile, err := writer.CreateFormFile("stamp", "stamp.png")
	if err != nil {
		return 0, nil, err
	}
	_, err = io.Copy(stampFile, stamp)
	if err != nil {
		return 0, nil, err
	}

	stampSource, err := writer.CreateFormField("stampSource")
	if err != nil {
		return 0, nil, err
	}
	if _, err := stampSource.Write([]byte("image")); err != nil {
		return 0, nil, err
	}

	stampOptions, err := writer.CreateFormField("stampOptions")
	if err != nil {
		return 0, nil, err
	}
	if _, err := stampOptions.Write([]byte(`{"rotation": "0"}`)); err != nil {
		return 0, nil, err
	}

	stampExpression, err := writer.CreateFormField("stampExpression")
	if err != nil {
		return 0, nil, err
	}
	if _, err := stampExpression.Write([]byte("stamp.png")); err != nil {
		return 0, nil, err
	}

	if err := writer.Close(); err != nil {
		return 0, nil, err
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, uri.String(), bytes.NewReader(body.Bytes()))
	if err != nil {
		return 0, nil, err
	}
	request.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.Do(request)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			responseBody = fmt.Appendf(nil, "error reading response body: %v", err)
		}

		return 0, nil, jhttp.APIError{
			Err:        fmt.Errorf("bad status, expected: %d, received: %d - with body %s", http.StatusOK, resp.StatusCode, string(responseBody)),
			StatusCode: resp.StatusCode,
		}
	}

	pdf, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, err
	}
	return resp.ContentLength, bytes.NewReader(pdf), nil
}
