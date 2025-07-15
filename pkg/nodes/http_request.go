package nodes

import (
    "bytes"
    "encoding/json"
    "go-workflow/pkg/framework"
    "html/template"
    retryablehttp "github.com/hashicorp/go-retryablehttp"
)

// HTTPRequest performs templated HTTP calls
type HTTPRequest struct {
    Method      string
    URLTemplate *template.Template
    BodyTmpl    *template.Template
}

func NewHTTPRequest(method, urlTmpl, bodyTmpl string) *HTTPRequest {
    return &HTTPRequest{
        Method:      method,
        URLTemplate: template.Must(template.New("url").Parse(urlTmpl)),
        BodyTmpl:    template.Must(template.New("body").Parse(bodyTmpl)),
    }
}

func (n *HTTPRequest) Execute(ctx *framework.Context, inputs []map[string]interface{}) ([]map[string]interface{}, error) {
    var out []map[string]interface{}
    for _, rec := range inputs {
        var urlBuf, bodyBuf bytes.Buffer
        n.URLTemplate.Execute(&urlBuf, rec)
        n.BodyTmpl.Execute(&bodyBuf, rec)

        req, _ := retryablehttp.NewRequest(n.Method, urlBuf.String(), &bodyBuf)
        req.Header.Set("Content-Type", "application/json")
        for k, v := range ctx.Env {
            req.Header.Set(k, v)
        }

        resp, err := ctx.HTTPClient.Do(req)
        if err != nil {
            return nil, err
        }
        defer resp.Body.Close()

        var parsed map[string]interface{}
        json.NewDecoder(resp.Body).Decode(&parsed)
        if items, ok := parsed["items"].([]interface{}); ok {
            for _, it := range items {
                out = append(out, it.(map[string]interface{}))
            }
        } else {
            out = append(out, parsed)
        }
    }
    return out, nil
}