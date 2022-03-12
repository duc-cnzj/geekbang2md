package constant

import (
	"time"

	"golang.org/x/time/rate"
)

var (
	// VideoDownloadParallel 视频并发下载数量
	VideoDownloadParallel int64 = 5
	// ImageDownloadParallel 图片并发下载数量
	ImageDownloadParallel int64 = 30

	// RequestLimit http 请求速率
	RequestLimit = rate.Every(5 * time.Second)
	// BucketSize 令牌桶初始大小
	BucketSize = 10
)
