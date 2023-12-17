package crawler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

type Momo struct {
	client *http.Client
}

// TODO: remove useless fields after confirming the data structure.

type MOMOBlock struct {
	ColumnType string      `json:"columnType"`
	Group      []MOMOGroup `json:"group"`
	Img        string      `json:"img"`
	MoreImg    string      `json:"moreImg"`
	MoreURL    string      `json:"moreUrl"`
	RecomdID   string      `json:"recomdId"`
	Title      string      `json:"title"`
	Type       string      `json:"type"`
	URL        string      `json:"url"`
}

type MOMOGroup struct {
	Goods []MOMOGood `json:"goods"`
}

type MOMOGood struct {
	GoodsCode       string `json:"goodsCode"`
	GoodsName       string `json:"goodsName"`
	GoodsStock      string `json:"goodsStock"`
	Img             string `json:"img"`
	ImgBottomTagURL string `json:"imgBottomTagUrl"`
	ImgLongTagURL   string `json:"imgLongTagUrl"`
	ImgPrice        string `json:"imgPrice"`
	ImgTag          string `json:"imgTag"`
	Instant         bool   `json:"instant"`
	MarketPrice     string `json:"marketPrice"`
	PrefixPrice     string `json:"prefixPrice"`
	ShowPrice       string `json:"showPrice"`
	URL             string `json:"url"`
}

type MOMOResp struct {
	Test    json.RawMessage
	RtnData struct {
		ResultCode string               `json:"resultCode"`
		ResultMsg  string               `json:"resultMsg"`
		BlockData  map[string]MOMOBlock `json:"blockData"`
	} `json:"rtnData"`
	RtnMsg  string `json:"rtnMsg"`
	RtnCode string `json:"rtnCode"`
}

func NewMomo() *Momo {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.MaxIdleConns = 10
	transport.MaxIdleConnsPerHost = 10

	return &Momo{
		client: &http.Client{
			Transport: transport,
			Timeout:   10 * time.Second,
		},
	}
}

var _ Crawler = (*Momo)(nil)

const momoBaseURL = "https://m.momoshop.com.tw/ajax/ajaxTool.jsp"

func (m *Momo) Fetch(ctx context.Context) ([]Product, error) {
	body, err := m.fetch(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetch: %w", err)
	}
	defer body.Close()

	resp, err := readToMOMOResp(body)
	if err != nil {
		return nil, fmt.Errorf("to resp: %w", err)
	}

	var products []Product
	// TODO: refactor, too deep loops in one function.
	for _, block := range resp.RtnData.BlockData {
		for _, group := range block.Group {
			for _, good := range group.Goods {
				price, err := momoPriceToDecimal(good.ShowPrice)
				if err != nil {
					if errors.Is(err, errNoPrice) {
						// ignore the product without price because of it is useless for us.
						continue
					}
					log.Println("failed to parse price:", err)
				}

				var discount decimal.Decimal
				marketPrice, err := momoPriceToDecimal(good.MarketPrice)
				if err == nil { // Market price available.
					discount = marketPrice.Sub(price)
				}

				products = append(products, Product{
					ProductImg: ProductImg{
						URL: good.Img,
					},

					ProductPrice: ProductPrice{
						Price:      price,
						Discount:   discount,
						SourceURL:  good.URL,
						SellerType: SellerTypeMOMO,
					},
					Name: good.GoodsName,
				})
			}
		}
	}

	return products, err
}
func (m *Momo) SellerType() SellerType {
	return SellerTypeMOMO
}

func (m *Momo) fetch(ctx context.Context) (io.ReadCloser, error) {
	// TODO: need to check dataValue for meeting requirements.
	dataValue := `{"flag":"getRecomdBlock","data":{"arrBlockId":["bt_7_602_01","bt_7_603_01","bt_7_605_01"],"categoryCode":"4202700000"}}`

	values := url.Values{}
	values.Set("data", dataValue)
	payload := values.Encode()

	req, err := http.NewRequestWithContext(
		ctx, http.MethodPost, momoBaseURL, strings.NewReader(payload),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// TODO: refactor to const after confirming the UA.
	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36 Edg/118.0.2088.61")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	nowEpoch := time.Now().UnixNano()
	query := req.URL.Query()
	query.Set("n", "getRecomdBlock")
	query.Set("t", strconv.FormatUint(uint64(nowEpoch), 10))
	req.URL.RawQuery = query.Encode()

	resp, err := m.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	return resp.Body, nil
}

func readToMOMOResp(reader io.Reader) (*MOMOResp, error) {
	var resp MOMOResp
	if err := json.NewDecoder(reader).Decode(&resp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return &resp, nil
}

var errNoPrice = fmt.Errorf("no price")

func momoPriceToDecimal(price string) (decimal.Decimal, error) {
	price = strings.ReplaceAll(price, ",", "")
	if price == "" {
		return decimal.Zero, errNoPrice
	}

	return decimal.NewFromString(price)
}
