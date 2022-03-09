package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/DuC-cnZj/geekbang2md/cache"
)

var c = &cache.Cache{}

type Product struct {
	ID        int `json:"id"`
	Spu       int `json:"spu"`
	Ctime     int `json:"ctime"`
	Utime     int `json:"utime"`
	BeginTime int `json:"begin_time"`
	EndTime   int `json:"end_time"`
	Price     struct {
		Market    int `json:"market"`
		Sale      int `json:"sale"`
		SaleType  int `json:"sale_type"`
		StartTime int `json:"start_time"`
		EndTime   int `json:"end_time"`
	} `json:"price"`
	IsOnborad     bool   `json:"is_onborad"`
	IsSale        bool   `json:"is_sale"`
	IsGroupbuy    bool   `json:"is_groupbuy"`
	IsPromo       bool   `json:"is_promo"`
	IsShareget    bool   `json:"is_shareget"`
	IsSharesale   bool   `json:"is_sharesale"`
	Type          string `json:"type"`
	IsColumn      bool   `json:"is_column"`
	IsCore        bool   `json:"is_core"`
	IsVideo       bool   `json:"is_video"`
	IsAudio       bool   `json:"is_audio"`
	IsDailylesson bool   `json:"is_dailylesson"`
	IsUniversity  bool   `json:"is_university"`
	IsOpencourse  bool   `json:"is_opencourse"`
	IsQconp       bool   `json:"is_qconp"`
	NavID         int    `json:"nav_id"`
	TimeNotSale   int    `json:"time_not_sale"`
	Title         string `json:"title"`
	Subtitle      string `json:"subtitle"`
	Intro         string `json:"intro"`
	IntroHTML     string `json:"intro_html"`
	Ucode         string `json:"ucode"`
	IsFinish      bool   `json:"is_finish"`
	Author        struct {
		Name      string `json:"name"`
		Intro     string `json:"intro"`
		Info      string `json:"info"`
		Avatar    string `json:"avatar"`
		BriefHTML string `json:"brief_html"`
		Brief     string `json:"brief"`
	} `json:"author"`
	Cover struct {
		Square            string `json:"square"`
		Rectangle         string `json:"rectangle"`
		Horizontal        string `json:"horizontal"`
		LectureHorizontal string `json:"lecture_horizontal"`
		LearnHorizontal   string `json:"learn_horizontal"`
		Transparent       string `json:"transparent"`
		Color             string `json:"color"`
	} `json:"cover"`
	Article struct {
		ID                int    `json:"id"`
		Count             int    `json:"count"`
		CountReq          int    `json:"count_req"`
		CountPub          int    `json:"count_pub"`
		TotalLength       int    `json:"total_length"`
		FirstArticleID    int    `json:"first_article_id"`
		FirstArticleTitle string `json:"first_article_title"`
	} `json:"article"`
	Seo struct {
		Keywords []string `json:"keywords"`
	} `json:"seo"`
	Share struct {
		Title   string `json:"title"`
		Content string `json:"content"`
		Cover   string `json:"cover"`
		Poster  string `json:"poster"`
	} `json:"share"`
	Labels     []int  `json:"labels"`
	Unit       string `json:"unit"`
	ColumnType int    `json:"column_type"`
	Column     struct {
		Unit             string `json:"unit"`
		Bgcolor          string `json:"bgcolor"`
		UpdateFrequency  string `json:"update_frequency"`
		IsPreorder       bool   `json:"is_preorder"`
		IsFinish         bool   `json:"is_finish"`
		IsIncludePreview bool   `json:"is_include_preview"`
		ShowChapter      bool   `json:"show_chapter"`
		IsSaleProduct    bool   `json:"is_sale_product"`
		StudyService     []int  `json:"study_service"`
		Path             struct {
			Desc     string `json:"desc"`
			DescHTML string `json:"desc_html"`
		} `json:"path"`
		IsCamp                   bool        `json:"is_camp"`
		CatalogPicURL            string      `json:"catalog_pic_url"`
		RecommendArticles        []int       `json:"recommend_articles"`
		RecommendComments        []int       `json:"recommend_comments"`
		Ranks                    interface{} `json:"ranks"`
		HotComments              interface{} `json:"hot_comments"`
		HotLines                 interface{} `json:"hot_lines"`
		DisplayType              int         `json:"display_type"`
		IntroBgStyle             int         `json:"intro_bg_style"`
		CommentTopAds            string      `json:"comment_top_ads"`
		ArticleFloatQrcodeURL    string      `json:"article_float_qrcode_url"`
		ArticleFloatAppQrcodeURL string      `json:"article_float_app_qrcode_url"`
		ArticleFloatQrcodeJump   string      `json:"article_float_qrcode_jump"`
		InRank                   bool        `json:"in_rank"`
	} `json:"column"`
	Dl struct {
		Article struct {
			ID            int    `json:"id"`
			VideoDuration string `json:"video_duration"`
			VideoHot      int    `json:"video_hot"`
			CouldPreview  bool   `json:"could_preview"`
		} `json:"article"`
		TopicIds []interface{} `json:"topic_ids"`
	} `json:"dl"`
	University struct {
		TotalHour       int    `json:"total_hour"`
		Term            int    `json:"term"`
		RedirectType    string `json:"redirect_type"`
		RedirectParam   string `json:"redirect_param"`
		WxQrcode        string `json:"wx_qrcode"`
		WxRule          string `json:"wx_rule"`
		ServerStartTime int    `json:"server_start_time"`
		LecturerHCover  string `json:"lecturer_h_cover"`
	} `json:"university"`
	Opencourse struct {
		VideoBg string `json:"video_bg"`
		Ad      struct {
			Cover         string `json:"cover"`
			CoverWeb      string `json:"cover_web"`
			RedirectType  string `json:"redirect_type"`
			RedirectParam string `json:"redirect_param"`
		} `json:"ad"`
		ArticleFav struct {
			Aid     int  `json:"aid"`
			HadDone bool `json:"had_done"`
			Count   int  `json:"count"`
		} `json:"article_fav"`
		AuthorHCover string `json:"author_h_cover"`
	} `json:"opencourse"`
	Qconp struct {
		TopicID      int    `json:"topic_id"`
		CoverAppoint string `json:"cover_appoint"`
		Article      struct {
			ID            int    `json:"id"`
			Cover         string `json:"cover"`
			VideoDuration string `json:"video_duration"`
			VideoHot      int    `json:"video_hot"`
		} `json:"article"`
	} `json:"qconp"`
	FavQrcode string `json:"fav_qrcode"`
	Extra     struct {
		Sub struct {
			Count      int  `json:"count"`
			HadDone    bool `json:"had_done"`
			CouldOrder bool `json:"could_order"`
			AccessMask int  `json:"access_mask"`
		} `json:"sub"`
		Fav struct {
			Count   int  `json:"count"`
			HadDone bool `json:"had_done"`
		} `json:"fav"`
		Rate struct {
			ArticleCount    int  `json:"article_count"`
			ArticleCountReq int  `json:"article_count_req"`
			IsFinished      bool `json:"is_finished"`
			RatePercent     int  `json:"rate_percent"`
			VideoSeconds    int  `json:"video_seconds"`
			LastArticleID   int  `json:"last_article_id"`
			LastChapterID   int  `json:"last_chapter_id"`
			HasLearn        bool `json:"has_learn"`
		} `json:"rate"`
		Cert struct {
			ID   string `json:"id"`
			Type int    `json:"type"`
		} `json:"cert"`
		Nps struct {
			Min    int    `json:"min"`
			Status int    `json:"status"`
			URL    string `json:"url"`
		} `json:"nps"`
		AnyRead struct {
			Total int `json:"total"`
			Count int `json:"count"`
		} `json:"any_read"`
		University struct {
			Status               int           `json:"status"`
			ViewStatus           int           `json:"view_status"`
			ChargeStatus         int           `json:"charge_status"`
			ShareRenewalStatus   int           `json:"share_renewal_status"`
			UnlockedStatus       int           `json:"unlocked_status"`
			UnlockedChapterIds   []interface{} `json:"unlocked_chapter_ids"`
			UnlockedChapterID    int           `json:"unlocked_chapter_id"`
			UnlockedChapterTitle string        `json:"unlocked_chapter_title"`
			UnlockedArticleCount int           `json:"unlocked_article_count"`
			UnlockedNextTime     int           `json:"unlocked_next_time"`
			ExpireTime           int           `json:"expire_time"`
			IsExpired            bool          `json:"is_expired"`
			IsGraduated          bool          `json:"is_graduated"`
			HadSub               bool          `json:"had_sub"`
			Timeline             []interface{} `json:"timeline"`
			HasWxFriend          bool          `json:"has_wx_friend"`
			SubTermTitle         string        `json:"sub_term_title"`
		} `json:"university"`
		Vip struct {
			IsYearCard bool `json:"is_year_card"`
			Show       bool `json:"show"`
			EndTime    int  `json:"end_time"`
		} `json:"vip"`
		Appoint struct {
			CouldDo bool `json:"could_do"`
			HadDone bool `json:"had_done"`
			Count   int  `json:"count"`
		} `json:"appoint"`
		GroupBuy struct {
			SuccessUcount int           `json:"success_ucount"`
			JoinCode      string        `json:"join_code"`
			CouldGroupbuy bool          `json:"could_groupbuy"`
			HadJoin       bool          `json:"had_join"`
			Price         int           `json:"price"`
			List          []interface{} `json:"list"`
		} `json:"group_buy"`
		Sharesale struct {
			OriginalPicColor    string `json:"original_pic_color"`
			OriginalPicURL      string `json:"original_pic_url"`
			PromoPicColor       string `json:"promo_pic_color"`
			PromoPicURL         string `json:"promo_pic_url"`
			ShareSalePrice      int    `json:"share_sale_price"`
			ShareSaleGuestPrice int    `json:"share_sale_guest_price"`
		} `json:"sharesale"`
		Promo struct {
			EntTime int `json:"ent_time"`
		} `json:"promo"`
		Channel struct {
			Is         bool `json:"is"`
			BackAmount int  `json:"back_amount"`
		} `json:"channel"`
		FirstPromo struct {
			Price     int  `json:"price"`
			CouldJoin bool `json:"could_join"`
		} `json:"first_promo"`
		CouponPromo struct {
			CouldJoin bool `json:"could_join"`
			Price     int  `json:"price"`
		} `json:"coupon_promo"`
		Helper []interface{} `json:"helper"`
		Tab    struct {
			Comment bool `json:"comment"`
		} `json:"tab"`
		Modules []struct {
			Name    string `json:"name"`
			Title   string `json:"title"`
			Content string `json:"content"`
			Type    string `json:"type"`
			IsTop   bool   `json:"is_top"`
		} `json:"modules"`
		Cid       int           `json:"cid"`
		FirstAids []interface{} `json:"first_aids"`
		StudyPlan struct {
			ID              int `json:"id"`
			DayNums         int `json:"day_nums"`
			ArticleNums     int `json:"article_nums"`
			LearnedWeekNums int `json:"learned_week_nums"`
			Status          int `json:"status"`
		} `json:"study_plan"`
		CateID   int    `json:"cate_id"`
		CateName string `json:"cate_name"`
		GroupTag struct {
			IsRecommend     bool `json:"is_recommend"`
			IsRecentlyLearn bool `json:"is_recently_learn"`
		} `json:"group_tag"`
		FirstAward struct {
			Show          bool   `json:"show"`
			Talks         int    `json:"talks"`
			Reads         int    `json:"reads"`
			Amount        int    `json:"amount"`
			ExpireTime    int    `json:"expire_time"`
			RedirectType  string `json:"redirect_type"`
			RedirectParam string `json:"redirect_param"`
		} `json:"first_award"`
	} `json:"extra"`
	AvailableCoupons []int `json:"available_coupons"`
	InPvip           int   `json:"in_pvip"`
}
type Data struct {
	HasExpiringProduct bool `json:"has_expiring_product"`
	LearnCount         struct {
		Total int `json:"total"`
	} `json:"learn_count"`
	List []struct {
		Pid             int    `json:"pid"`
		Ptype           string `json:"ptype"`
		Aid             int    `json:"aid"`
		Ctime           int64  `json:"ctime"`
		Score           int    `json:"score"`
		IsExpire        bool   `json:"is_expire"`
		ExpireTime      int    `json:"expire_time"`
		LastLearnedTime int    `json:"last_learned_time"`
	} `json:"list"`
	Articles []interface{} `json:"articles"`
	Products []Product     `json:"products"`
	Page     struct {
		More   bool `json:"more"`
		Count  int  `json:"count"`
		Score  int  `json:"score"`
		Score0 int  `json:"score0"`
	} `json:"page"`
}
type ApiProjectResponse struct {
	Code  int   `json:"code"`
	Data  *Data `json:"data"`
	Extra struct {
		Cost      float64 `json:"cost"`
		RequestID string  `json:"request-id"`
	} `json:"extra"`
}

func Products() (ApiProjectResponse, error) {
	var result ApiProjectResponse

	file, err := c.Get("products")
	if err == nil && len(file) > 0 {
		err = json.NewDecoder(bytes.NewReader(file)).Decode(&result)
		if err == nil {
			return result, err
		}
	}

	res, err := HttpClient.Post("https://time.geekbang.org/serv/v3/learn/product", `{"desc":true,"expire":1,"last_learn":0,"learn_status":0,"prev":0,"size":100,"sort":1,"type":"c1","with_learn_count":1}`, false)
	if err != nil {
		return ApiProjectResponse{}, err
	}
	defer func() {
		io.Copy(io.Discard, res.Body)
		res.Body.Close()
	}()
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return ApiProjectResponse{}, err
	}
	if res.StatusCode < 400 {
		c.Set("products", result)
	}

	return result, nil
}

type ArticleResponse struct {
	Data struct {
		TextReadVersion int           `json:"text_read_version"`
		AudioSize       int           `json:"audio_size"`
		ArticleCover    string        `json:"article_cover"`
		Subtitles       []interface{} `json:"subtitles"`
		ProductType     string        `json:"product_type"`
		AudioDubber     string        `json:"audio_dubber"`
		IsFinished      bool          `json:"is_finished"`
		Like            struct {
			HadDone bool `json:"had_done"`
			Count   int  `json:"count"`
		} `json:"like"`
		AudioTime string `json:"audio_time"`
		Share     struct {
			Content string `json:"content"`
			Title   string `json:"title"`
			Poster  string `json:"poster"`
			Cover   string `json:"cover"`
		} `json:"share"`
		ArticleContent     string `json:"article_content"`
		FloatQrcode        string `json:"float_qrcode"`
		ArticleCoverHidden bool   `json:"article_cover_hidden"`
		IsRequired         bool   `json:"is_required"`
		Score              string `json:"score"`
		LikeCount          int    `json:"like_count"`
		ArticleSubtitle    string `json:"article_subtitle"`
		VideoTime          string `json:"video_time"`
		HadViewed          bool   `json:"had_viewed"`
		ArticleTitle       string `json:"article_title"`
		ColumnBgcolor      string `json:"column_bgcolor"`
		OfflinePackage     string `json:"offline_package"`
		AudioTitle         string `json:"audio_title"`
		AudioTimeArr       struct {
			M string `json:"m"`
			S string `json:"s"`
			H string `json:"h"`
		} `json:"audio_time_arr"`
		TextReadPercent int64  `json:"text_read_percent"`
		Cid             int    `json:"cid"`
		ArticleCshort   string `json:"article_cshort"`
		VideoWidth      int    `json:"video_width"`
		ColumnCouldSub  bool   `json:"column_could_sub"`
		VideoID         string `json:"video_id"`
		Sku             string `json:"sku"`
		VideoCover      string `json:"video_cover"`
		AuthorName      string `json:"author_name"`
		ColumnIsOnboard bool   `json:"column_is_onboard"`
		AudioURL        string `json:"audio_url"`
		ChapterID       string `json:"chapter_id"`
		ColumnHadSub    bool   `json:"column_had_sub"`
		ColumnCover     string `json:"column_cover"`
		RatePercent     int    `json:"rate_percent"`
		FooterCoverData struct {
			ImgURL  string `json:"img_url"`
			LinkURL string `json:"link_url"`
			MpURL   string `json:"mp_url"`
		} `json:"footer_cover_data"`
		FloatAppQrcode     string `json:"float_app_qrcode"`
		ColumnIsExperience bool   `json:"column_is_experience"`
		Rate               struct {
			Num1 struct {
				CurVersion     int  `json:"cur_version"`
				MaxRate        int  `json:"max_rate"`
				CurRate        int  `json:"cur_rate"`
				IsFinished     bool `json:"is_finished"`
				TotalRate      int  `json:"total_rate"`
				LearnedSeconds int  `json:"learned_seconds"`
			} `json:"1"`
			Num2 struct {
				CurVersion     int  `json:"cur_version"`
				MaxRate        int  `json:"max_rate"`
				CurRate        int  `json:"cur_rate"`
				IsFinished     bool `json:"is_finished"`
				TotalRate      int  `json:"total_rate"`
				LearnedSeconds int  `json:"learned_seconds"`
			} `json:"2"`
			Num3 struct {
				CurVersion     int  `json:"cur_version"`
				MaxRate        int  `json:"max_rate"`
				CurRate        int  `json:"cur_rate"`
				IsFinished     bool `json:"is_finished"`
				TotalRate      int  `json:"total_rate"`
				LearnedSeconds int  `json:"learned_seconds"`
			} `json:"3"`
		} `json:"rate"`
		ProductID           int    `json:"product_id"`
		HadLiked            bool   `json:"had_liked"`
		ID                  int    `json:"id"`
		FreeGet             bool   `json:"free_get"`
		IsVideoPreview      bool   `json:"is_video_preview"`
		ArticleSummary      string `json:"article_summary"`
		ColumnSaleType      int    `json:"column_sale_type"`
		FloatQrcodeJump     string `json:"float_qrcode_jump"`
		ColumnID            int    `json:"column_id"`
		VideoHeight         int    `json:"video_height"`
		ArticleFeatures     int    `json:"article_features"`
		ArticlePosterWxlite string `json:"article_poster_wxlite"`
		AudioMd5            string `json:"audio_md5"`
		CommentCount        int    `json:"comment_count"`
		VideoSize           int    `json:"video_size"`
		Offline             struct {
			Size        int    `json:"size"`
			FileName    string `json:"file_name"`
			DownloadURL string `json:"download_url"`
		} `json:"offline"`
		ArticleCouldPreview bool          `json:"article_could_preview"`
		HlsVideos           []interface{} `json:"hls_videos"`
		InPvip              int           `json:"in_pvip"`
		AudioDownloadURL    string        `json:"audio_download_url"`
		ArticleCtime        int           `json:"article_ctime"`
		ArticleSharetitle   string        `json:"article_sharetitle"`
	} `json:"data"`
	Code int `json:"code"`
}

// Article 获取cid
func Article(id string) (ArticleResponse, error) {
	var result ArticleResponse
	file, err := c.Get("article-" + id)
	if err == nil && len(file) > 0 {
		err = json.NewDecoder(bytes.NewReader(file)).Decode(&result)
		if err == nil {
			return result, err
		}
	}

	res, err := HttpClient.Post("https://time.geekbang.org/serv/v1/article", fmt.Sprintf(`{"id":"%s","include_neighbors":true,"is_freelyread":true}`, id), false)
	if err != nil {
		return ArticleResponse{}, err
	}
	defer func() {
		io.Copy(io.Discard, res.Body)
		res.Body.Close()
	}()
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return ArticleResponse{}, err
	}

	if res.StatusCode < 400 {
		c.Set("article-"+id, result)
	}

	return result, nil
}

type ArticlesResponseItem struct {
	AuthorName   string        `json:"author_name,omitempty"`
	AudioSize    int           `json:"audio_size,omitempty"`
	IncludeAudio bool          `json:"include_audio"`
	Subtitles    []interface{} `json:"subtitles"`
	AudioURL     string        `json:"audio_url,omitempty"`
	ChapterID    string        `json:"chapter_id"`
	ColumnHadSub bool          `json:"column_had_sub"`
	AudioDubber  string        `json:"audio_dubber,omitempty"`
	IsFinished   bool          `json:"is_finished"`
	AudioTime    string        `json:"audio_time,omitempty"`
	RatePercent  int           `json:"rate_percent"`
	ColumnSku    int           `json:"column_sku"`
	IsRequired   bool          `json:"is_required"`
	Rate         struct {
		Num1 struct {
			CurVersion     int  `json:"cur_version"`
			MaxRate        int  `json:"max_rate"`
			CurRate        int  `json:"cur_rate"`
			IsFinished     bool `json:"is_finished"`
			TotalRate      int  `json:"total_rate"`
			LearnedSeconds int  `json:"learned_seconds"`
		} `json:"1"`
		Num2 struct {
			CurVersion     int  `json:"cur_version"`
			MaxRate        int  `json:"max_rate"`
			CurRate        int  `json:"cur_rate"`
			IsFinished     bool `json:"is_finished"`
			TotalRate      int  `json:"total_rate"`
			LearnedSeconds int  `json:"learned_seconds"`
		} `json:"2"`
		Num3 struct {
			CurVersion     int  `json:"cur_version"`
			MaxRate        int  `json:"max_rate"`
			CurRate        int  `json:"cur_rate"`
			IsFinished     bool `json:"is_finished"`
			TotalRate      int  `json:"total_rate"`
			LearnedSeconds int  `json:"learned_seconds"`
		} `json:"3"`
	} `json:"rate"`
	Score            int64  `json:"score"`
	ArticleSubtitle  string `json:"article_subtitle"`
	AudioDownloadURL string `json:"audio_download_url,omitempty"`
	ID               int    `json:"id"`
	HadViewed        bool   `json:"had_viewed"`
	ArticleTitle     string `json:"article_title"`
	ColumnBgcolor    string `json:"column_bgcolor,omitempty"`
	IsVideoPreview   bool   `json:"is_video_preview"`
	ArticleSummary   string `json:"article_summary"`
	ColumnID         int    `json:"column_id,omitempty"`
	AudioTitle       string `json:"audio_title,omitempty"`
	AudioTimeArr     struct {
		M string `json:"m"`
		S string `json:"s"`
		H string `json:"h"`
	} `json:"audio_time_arr,omitempty"`
	AuthorIntro string `json:"author_intro,omitempty"`
	Offline     struct {
		FileName    string `json:"file_name"`
		DownloadURL string `json:"download_url"`
	} `json:"offline"`
	ArticleCover        string `json:"article_cover"`
	ArticleSharetitle   string `json:"article_sharetitle,omitempty"`
	ArticleCouldPreview bool   `json:"article_could_preview"`
	AudioMd5            string `json:"audio_md5,omitempty"`
	ArticleCtime        int    `json:"article_ctime"`
	ColumnCover         string `json:"column_cover,omitempty"`
}

type ArticlesResponse struct {
	Data struct {
		List []*ArticlesResponseItem `json:"list"`
		Page struct {
			Count int  `json:"count"`
			More  bool `json:"more"`
		} `json:"page"`
	} `json:"data"`
	Code int `json:"code"`
}

func Articles(cid int) (ArticlesResponse, error) {
	var result ArticlesResponse
	file, err := c.Get(fmt.Sprintf("articles-%d", cid))
	if err == nil && len(file) > 0 {
		err = json.NewDecoder(bytes.NewReader(file)).Decode(&result)
		if err == nil {
			return result, err
		}
	}
	res, err := HttpClient.Post("https://time.geekbang.org/serv/v1/column/articles",
		fmt.Sprintf(`{"cid":%d,"size":100,"prev":0,"order":"earliest","sample":false}`, cid), false)
	if err != nil {
		return ArticlesResponse{}, err
	}
	defer func() {
		io.Copy(io.Discard, res.Body)
		res.Body.Close()
	}()
	all, _ := io.ReadAll(res.Body)
	err = json.NewDecoder(bytes.NewReader(all)).Decode(&result)
	if err != nil {
		log.Println(string(all), cid)
		return ArticlesResponse{}, err
	}

	if res.StatusCode < 400 {
		c.Set(fmt.Sprintf("articles-%d", cid), result)
	}
	return result, nil
}