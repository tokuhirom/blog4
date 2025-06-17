package admin

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"regexp"
)

// https://amzn.to/42051PN のような記法を B01M2BOZDL に変換する
// https://www.amazon.co.jp/-/en/gp/product/B08WT889V3?ie=UTF8&psc=1&linkCode=sl1&tag=tokuhirom-22&linkId=30e33501a13966a30b888de5b4aef836&language=en_US&ref_=as_li_ss_tl
// https://www.amazon.co.jp/%E3%83%8A%E3%82%AC%E3%82%AA-%E7%87%95%E4%B8%89%E6%9D%A1-%E3%82%B9%E3%83%86%E3%83%B3%E3%83%AC%E3%82%B9-%E3%83%87%E3%82%A3%E3%83%8A%E3%83%BC%E3%83%95%E3%82%A9%E3%83%BC%E3%82%AF-5%E6%9C%AC/dp/B01M2BOZDL?_encoding=UTF8&pd_rd_w=a6agn&content-id=amzn1.sym.bcc66df3-c2cc-4242-967e-174aec86af7a:amzn1.symc.a9cb614c-616d-4684-840d-556cb89e228d&pf_rd_p=bcc66df3-c2cc-4242-967e-174aec86af7a&pf_rd_r=14B098QT0WFV8K1YYMA0&pd_rd_wg=N1Vqd&pd_rd_r=0b99f5ff-8cf1-4ccb-84d0-e958cc8f64f5&th=1&linkCode=sl1&tag=tokuhirom-22&linkId=0bb114c89147942275c86c1dc5efbe4f&language=ja_JP&ref_=as_li_ss_tl

func rewriteAmazonShortUrlInMarkdown(markdown string) string {
	// /https://amzn.to/[a-zA-Z0-9]+ にマッチする文字列を探して､[asin:${asin}:detail] に置換する
	re := regexp.MustCompile(`https://amzn.to/[a-zA-Z0-9]+`)
	return re.ReplaceAllStringFunc(markdown, func(url string) string {
		asin, err := amazonShortUrlToAsin(url)
		if err != nil {
			slog.Error("failed to rewrite amazon short URL", slog.String("url", url), slog.Any("error", err))
			return url
		}
		return fmt.Sprintf("[asin:%s:detail]", asin)
	})
}

// amazonShortUrlToAsin fetches the short URL and extracts the ASIN from the redirect location.
func amazonShortUrlToAsin(url string) (string, error) {
	client := &http.Client{
		// Prevents the client from following the redirect automatically
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	// Perform the GET request
	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			slog.Error("Error closing response body", slog.Any("error", err))
		}
	}(resp.Body)

	// Get the Location header from the response
	location := resp.Header.Get("Location")
	if location == "" {
		return "", fmt.Errorf("location header not found in the response for URL: %s", url)
	}

	// Regular expression to find the ASIN in the URL
	re := regexp.MustCompile(`/((?:dp|product)/([A-Z0-9]{10}))`)
	matches := re.FindStringSubmatch(location)
	if len(matches) < 2 {
		return "", fmt.Errorf("ASIN not found in the Location URL: %s", location)
	}

	// Return the extracted ASIN
	return matches[2], nil
}
