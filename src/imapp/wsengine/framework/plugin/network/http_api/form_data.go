package httpApi

import (
	"bytes"
	"fmt"
	hertzProtocol "github.com/cloudwego/hertz/pkg/protocol"
	"io"
	"mime/multipart"
	"net/textproto"
	"strings"
	functionTools "ws/framework/utils/function_tools"
)

// formDataContentType .
type formDataContentType uint8

const (
	// TextPlain .
	TextPlain formDataContentType = iota
	// ApplicationJson .
	ApplicationJson
	// ApplicationOctetStream .
	ApplicationOctetStream
)

// MultipartField .
type MultipartField struct {
	Name        string
	FileName    string
	ContentType formDataContentType
	Body        interface{}
}

func createMIMEHeader(name, header string) textproto.MIMEHeader {
	quoteEscape := strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"`, quoteEscape.Replace(name)))
	h.Set("Content-Type", header)

	return h
}

// FormData .
func FormData(boundary string, fields []MultipartField) RequestOptionsFn {
	return func(req *hertzProtocol.Request) (err error) {
		buf := bytes.NewBuffer(make([]byte, 0))

		multipartWriter := multipart.NewWriter(buf)
		err = multipartWriter.SetBoundary(boundary)
		if err != nil {
			return
		}

		var fieldWriter io.Writer

		for i := range fields {
			f := fields[i]

			switch f.ContentType {
			case TextPlain:
				fieldHeader := createMIMEHeader(f.Name, "text/plain")
				fieldWriter, err = multipartWriter.CreatePart(fieldHeader)
			case ApplicationJson:
				fieldHeader := createMIMEHeader(f.Name, "application/json")
				fieldWriter, err = multipartWriter.CreatePart(fieldHeader)
			case ApplicationOctetStream:
				fieldWriter, err = multipartWriter.CreateFormFile(f.Name, f.FileName)
			default:
				err = fmt.Errorf("unsupported %v", f.ContentType)
			}

			if err != nil {
				return
			}

			switch f.Body.(type) {
			case string:
				_, err = fieldWriter.Write(functionTools.S2B(f.Body.(string)))
			case []byte:
				_, err = fieldWriter.Write(f.Body.([]byte))
			default:
				err = fmt.Errorf("unsupported %s body type", f.Name)
			}

			if err != nil {
				return
			}
		}

		err = multipartWriter.Close()
		if err != nil {
			return
		}

		req.Header.SetContentTypeBytes([]byte(multipartWriter.FormDataContentType()))
		req.SetBody(buf.Bytes())
		return
	}
}
