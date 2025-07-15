package nodes

import (
    "bytes"
    "encoding/json"
    "fmt"
    "go-workflow/pkg/framework"
    retryablehttp "github.com/hashicorp/go-retryablehttp"
)

// HTTPRequest performs templated HTTP calls
type HTTPRequest struct {
    URLKey      string
    MethodKey   string
    HeadersKey  string
    BodyKey     string
}

func NewHTTPRequest(urlKey, methodKey, headersKey, bodyKey string) *HTTPRequest {
    return &HTTPRequest{
        URLKey:      urlKey,
        MethodKey:   methodKey,
        HeadersKey:  headersKey,
        BodyKey:     bodyKey,
    }
}

func (n *HTTPRequest) Execute(ctx *framework.Context, inputs []map[string]interface{}) ([]map[string]interface{}, error) {
    var out []map[string]interface{}
    for _, rec := range inputs {
        urlStr, ok := rec[n.URLKey].(string)
        if !ok {
            return nil, fmt.Errorf("URL not found or not a string in input record for key %s", n.URLKey)
        }
        methodStr, ok := rec[n.MethodKey].(string)
        if !ok {
            return nil, fmt.Errorf("Method not found or not a string in input record for key %s", n.MethodKey)
        }

        var bodyBytes []byte
        if n.BodyKey != "" {
            bodyContent, ok := rec[n.BodyKey]
            if !ok {
                return nil, fmt.Errorf("Body not found in input record for key %s", n.BodyKey)
            }
            var err error
            bodyBytes, err = json.Marshal(bodyContent)
            if err != nil {
                return nil, fmt.Errorf("failed to marshal body content: %w", err)
            }
        }

        req, err := retryablehttp.NewRequest(methodStr, urlStr, bytes.NewBuffer(bodyBytes))
        if err != nil {
            return nil, fmt.Errorf("failed to create request: %w", err)
        }

        req.Header.Set("Content-Type", "application/json")

        if n.HeadersKey != "" {
            headers, ok := rec[n.HeadersKey].(map[string]interface{})
            if !ok {
                return nil, fmt.Errorf("Headers not found or not a map in input record for key %s", n.HeadersKey)
            }
            for k, v := range headers {
                if headerVal, isString := v.(string); isString {
                    req.Header.Set(k, headerVal)
                }
            }
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