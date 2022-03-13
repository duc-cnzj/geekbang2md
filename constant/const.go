package constant

import (
	"time"

	"golang.org/x/time/rate"
)

var (
	// VideoDownloadParallelNum 视频并发下载数量
	VideoDownloadParallelNum int64 = 5
	// ImageDownloadParallelNum 图片并发下载数量
	ImageDownloadParallelNum int64 = 30

	// RequestLimit http 请求速率
	RequestLimit = rate.Every(5 * time.Second)
	// BucketSize 令牌桶初始大小
	BucketSize = 10
)
