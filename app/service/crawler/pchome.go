package crawler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"orange-backstage-api/infra/util/convert"
	"regexp"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"golang.org/x/net/html"
)

type Pchome struct {
	client *http.Client
}

func NewPchome() *Pchome {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.MaxIdleConns = 10
	transport.MaxIdleConnsPerHost = 10

	return &Pchome{
		client: &http.Client{
			Transport: transport,
			Timeout:   10 * time.Second,
		},
	}
}

func (p *Pchome) SellerType() SellerType {
	return SellerTypePCHome
}

type PCHOMENode struct {
	ID  uint64 `json:"Id"`
	Img struct {
		Src string `json:"Src"`
	}
	Link struct {
		URL string `json:"Url"`

		// Text is the product name
		Text string `json:"Text"`

		// Text1 is the price, but it may not a valid json string.
		// Convert it to string first, then unmarshal it.
		Text1 json.RawMessage `json:"Text1"`
	}
}

type PCHomeResp struct {
	Nodes []PCHOMENode `json:"Nodes"`
}

const (
	pchomeNotebookURL = "https://24h.pchome.com.tw/cdn/region/DSAA/data&28368754"
	pchomeTabletURL   = "https://24h.pchome.com.tw/cdn/region/DYAM/data&28368773"
)

func (p *Pchome) Fetch(ctx context.Context) ([]Product, error) {
	type fetchType []struct {
		url string
		pt  ProductType
	}
	productTypes := fetchType{
		{
			url: pchomeNotebookURL,
			pt:  ProductTypeLaptop,
		},
		{
			url: pchomeTabletURL,
			pt:  ProductTypeTablet,
		},
	}

	var products []Product
	for _, ft := range productTypes {
		body, err := p.fetch(ctx, ft.url)
		if err != nil {
			return nil, fmt.Errorf("fetch: %w", err)
		}
		defer body.Close()

		bs, err := io.ReadAll(body)
		if err != nil {
			return nil, fmt.Errorf("read body: %w", err)
		}

		jsonString, err := extractPchomeJSON(convert.BytesToStr(bs))
		if err != nil {
			return nil, fmt.Errorf("extract json: %w", err)
		}

		var resps []PCHomeResp
		if err := json.Unmarshal(convert.StrToBytes(jsonString), &resps); err != nil {
			return nil, fmt.Errorf("unmarshal: %w", err)
		}

		for _, r := range resps {
			for _, node := range r.Nodes {
				if node.Link.Text == "" {
					continue
				}

				name, err := extractTextFromHTML(node.Link.Text)
				if err != nil {
					log.Println("failed to extract product name:", err)
					continue
				}

				var text1 string
				if err := json.Unmarshal(node.Link.Text1, &text1); err != nil {
					continue
				}

				price, err := decimal.NewFromString(convert.BytesToStr(convert.StrToBytes(text1)))
				if err != nil {
					continue
				}

				products = append(products, Product{
					ProductPrice: ProductPrice{
						Price:      price,
						SourceURL:  node.Link.URL,
						SellerType: SellerTypePCHome,
					},
					Name: name,
					Type: ft.pt,
				})
			}
		}
	}

	return products, nil
}

func (p *Pchome) fetch(ctx context.Context, url string) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(
		ctx, http.MethodGet, url, nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	return resp.Body, nil
}

var pchomeJSONRegex = regexp.MustCompile(`=\s*([^;]+);`)

func extractPchomeJSON(input string) (string, error) {
	match := pchomeJSONRegex.FindStringSubmatch(input)

	if len(match) < 2 {
		return "", fmt.Errorf("no match found")
	}

	return match[1], nil
}

func extractTextFromHTML(htmlString string) (string, error) {
	doc, err := html.Parse(strings.NewReader(htmlString))
	if err != nil {
		return "", err
	}

	var textContent string
	var extractText func(*html.Node)
	extractText = func(n *html.Node) {
		if n.Type == html.TextNode {
			textContent += n.Data
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extractText(c)
		}
	}

	extractText(doc)
	return textContent, nil
}

var _ Crawler = (*Pchome)(nil)
