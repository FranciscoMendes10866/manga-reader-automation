package helpers

import (
	"time"

	"github.com/gocolly/colly"
)

var ScrapLimitOptions = &colly.LimitRule{
	RandomDelay: 5 * time.Second,
}
