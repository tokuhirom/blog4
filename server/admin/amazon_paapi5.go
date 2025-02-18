package admin

import (
	"context"
	"fmt"
	paapi5 "github.com/goark/pa-api"
	"github.com/goark/pa-api/entity"
	"github.com/goark/pa-api/query"
	"log"
)

type PAAPIClient struct {
	AccessKey string
	SecretKey string
}

func NewPAAPIClient(accessKey string, secretKey string) *PAAPIClient {
	return &PAAPIClient{
		AccessKey: accessKey,
		SecretKey: secretKey,
	}
}

type AmazonProductDetail struct {
	ASIN           string
	Title          string
	ImageMediumURL string
	Link           string
}

// FetchAmazonProductDetails fetches Amazon product details by ASINs
// asins: e.g. []string{"B07YCM5K55"}. up to 10 ASINs
// https://webservices.amazon.com/paapi5/documentation/get-items.html
func (c *PAAPIClient) FetchAmazonProductDetails(ctx context.Context, asins []string) ([]AmazonProductDetail, error) {
	if len(asins) > 10 {
		return nil, fmt.Errorf("too many ASINs")
	}

	log.Printf("Fetching Amazon product details by ASINs: %v", asins)

	client := paapi5.New(
		paapi5.WithMarketplace(paapi5.LocaleJapan),
	).CreateClient(
		"tokuhirom-22",
		c.AccessKey,
		c.SecretKey,
	)

	// Make query
	q := query.NewGetItems(
		client.Marketplace(),
		client.PartnerTag(),
		client.PartnerType(),
	).ASINs(asins).EnableImages().EnableItemInfo().EnableParentASIN()

	// Request and response
	body, err := client.RequestContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("failed to request: %w", err)
	}

	res, err := entity.DecodeResponse(body)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	fmt.Println(res.String())

	if len(res.ItemsResult.Items) == 0 {
		return nil, fmt.Errorf("no items found in PA-API response")
	}

	var productDetails []AmazonProductDetail
	for _, item := range res.ItemsResult.Items {
		fmt.Printf("ASIN: %s\n", item.ASIN)
		fmt.Printf("Title: %s\n", item.ItemInfo.Title.DisplayValue)
		fmt.Printf("URL: %s\n", item.DetailPageURL)
		productDetails = append(productDetails, AmazonProductDetail{
			ASIN:           item.ASIN,
			Title:          item.ItemInfo.Title.DisplayValue,
			ImageMediumURL: item.Images.Primary.Medium.URL,
			Link:           item.DetailPageURL,
		})
	}
	log.Printf("Fetched Amazon product details: %v", productDetails)
	return productDetails, nil
}
