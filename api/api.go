package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/duc-cnzj/geekbang2md/cache"
	"github.com/duc-cnzj/geekbang2md/utils"
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

type ProductList []Product

func (p ProductList) Len() int {
	return len(p)
}

func (p ProductList) Less(i, j int) bool {
	return p[i].Type == ProductTypeZhuanlan
}

func (p ProductList) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

type ProjectResponseData struct {
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
	Products ProductList   `json:"products"`
	Page     struct {
		More   bool `json:"more"`
		Count  int  `json:"count"`
		Score  int  `json:"score"`
		Score0 int  `json:"score0"`
	} `json:"page"`
}
type ProjectResponse struct {
	Code  int                  `json:"code"`
	Data  *ProjectResponseData `json:"data"`
	Extra struct {
		Cost      float64 `json:"cost"`
		RequestID string  `json:"request-id"`
	} `json:"extra"`
}

type PType = string

const (
	ProductTypeZhuanlan PType = "c1"
	ProductTypeVideo    PType = "c3"
	ProductTypeAll      PType = ""
)

type InfosResponse struct {
	Code int `json:"code"`
	Data struct {
		Infos []struct {
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
				RecommendArticles        interface{} `json:"recommend_articles"`
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
				Cid       int   `json:"cid"`
				FirstAids []int `json:"first_aids"`
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
		} `json:"infos"`
		Articles []struct {
			ID           int    `json:"id"`
			Type         int    `json:"type"`
			Pid          int    `json:"pid"`
			ChapterID    int    `json:"chapter_id"`
			ChapterTitle string `json:"chapter_title"`
			Title        string `json:"title"`
			Subtitle     string `json:"subtitle"`
			ShareTitle   string `json:"share_title"`
			Summary      string `json:"summary"`
			Ctime        int    `json:"ctime"`
			Cover        struct {
				Default string `json:"default"`
			} `json:"cover"`
			Author struct {
				Name   string `json:"name"`
				Avatar string `json:"avatar"`
			} `json:"author"`
			Audio struct {
				Title       string   `json:"title"`
				Dubber      string   `json:"dubber"`
				DownloadURL string   `json:"download_url"`
				Md5         string   `json:"md5"`
				Size        int      `json:"size"`
				Time        string   `json:"time"`
				TimeArr     []string `json:"time_arr"`
				URL         string   `json:"url"`
			} `json:"audio"`
			Video struct {
				ID       string `json:"id"`
				Duration int    `json:"duration"`
				Cover    string `json:"cover"`
				Width    int    `json:"width"`
				Height   int    `json:"height"`
				Size     int    `json:"size"`
				Time     string `json:"time"`
				Medias   []struct {
					Size    int    `json:"size"`
					Quality string `json:"quality"`
					URL     string `json:"url"`
				} `json:"medias"`
				HlsVid    string `json:"hls_vid"`
				HlsMedias []struct {
					Size    int    `json:"size"`
					Quality string `json:"quality"`
					URL     string `json:"url"`
				} `json:"hls_medias"`
				Subtitles []interface{} `json:"subtitles"`
				Tips      []interface{} `json:"tips"`
			} `json:"video"`
			VideoPreview struct {
				Duration int         `json:"duration"`
				Medias   interface{} `json:"medias"`
			} `json:"video_preview"`
			VideoPreviews []struct {
				Duration int `json:"duration"`
				Medias   []struct {
					Size    int    `json:"size"`
					Quality string `json:"quality"`
					URL     string `json:"url"`
				} `json:"medias"`
			} `json:"video_previews"`
			CouldPreview      bool   `json:"could_preview"`
			VideoCouldPreview bool   `json:"video_could_preview"`
			CoverHidden       bool   `json:"cover_hidden"`
			Content           string `json:"content"`
			IsRequired        bool   `json:"is_required"`
			Extra             struct {
				Rate []struct {
					Type           int  `json:"type"`
					CurVersion     int  `json:"cur_version"`
					CurRate        int  `json:"cur_rate"`
					MaxRate        int  `json:"max_rate"`
					TotalRate      int  `json:"total_rate"`
					LearnedSeconds int  `json:"learned_seconds"`
					IsFinished     bool `json:"is_finished"`
				} `json:"rate"`
				RatePercent int  `json:"rate_percent"`
				IsFinished  bool `json:"is_finished"`
				Fav         struct {
					Count   int  `json:"count"`
					HadDone bool `json:"had_done"`
				} `json:"fav"`
				IsUnlocked bool `json:"is_unlocked"`
				Learn      struct {
					Ucount int `json:"ucount"`
				} `json:"learn"`
				FooterCoverData struct {
					ImgURL  string `json:"img_url"`
					MpURL   string `json:"mp_url"`
					LinkURL string `json:"link_url"`
				} `json:"footer_cover_data"`
			} `json:"extra"`
			Score           int    `json:"score"`
			IsVideo         bool   `json:"is_video"`
			PosterWxlite    string `json:"poster_wxlite"`
			HadFreelyread   bool   `json:"had_freelyread"`
			FloatQrcode     string `json:"float_qrcode"`
			FloatAppQrcode  string `json:"float_app_qrcode"`
			FloatQrcodeJump string `json:"float_qrcode_jump"`
			InPvip          int    `json:"in_pvip"`
			CommentCount    int    `json:"comment_count"`
			Cshort          string `json:"cshort"`
			Like            struct {
				Count   int  `json:"count"`
				HadDone bool `json:"had_done"`
			} `json:"like"`
		} `json:"articles"`
		Labels interface{} `json:"labels"`
	} `json:"data"`
	Error struct {
	} `json:"error"`
	Extra struct {
		Cost      float64 `json:"cost"`
		RequestID string  `json:"request-id"`
	} `json:"extra"`
}

type IntString []string

func (is IntString) Len() int {
	return len(is)
}

func (is IntString) Less(i, j int) bool {
	ii, _ := strconv.Atoi(is[i])
	ji, _ := strconv.Atoi(is[j])
	return ii < ji
}

func (is IntString) Swap(i, j int) {
	is[i], is[j] = is[j], is[i]
}

func Infos(chunks IntString) (*InfosResponse, error) {
	var result *InfosResponse
	sort.Sort(chunks)
	idStr := strings.Join(chunks, ",")
	cacheKey := "infos-" + utils.Md5(strings.Join(chunks, "-"))
	file, err := c.Get(cacheKey)
	if err == nil && len(file) > 0 {
		err = json.NewDecoder(bytes.NewReader(file)).Decode(&result)
		if err == nil {
			return result, err
		}
	}
	res, err := HttpClient.Post("https://time.geekbang.org/serv/v3/product/infos", fmt.Sprintf(`{"ids":[%s],"with_first_articles":true}`, idStr), false)
	if err != nil {
		return nil, err
	}
	defer func() {
		io.Copy(io.Discard, res.Body)
		res.Body.Close()
	}()
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	if res.StatusCode < 400 {
		c.Set(cacheKey, result)
	}

	return result, nil
}

type SkusResponse struct {
	Error []interface{} `json:"error"`
	Extra []interface{} `json:"extra"`
	Data  struct {
		List []struct {
			SubCount          int  `json:"sub_count"`
			IsVip             bool `json:"is_vip"`
			ColumnType        int  `json:"column_type"`
			ID                int  `json:"id"`
			ColumnPriceMarket int  `json:"column_price_market"`
			ColumnPriceFirst  int  `json:"column_price_first"`
			TopLevel          int  `json:"top_level"`
			LastAid           int  `json:"last_aid"`
			HadSub            bool `json:"had_sub"`
			PriceType         int  `json:"price_type"`
			IsExperience      bool `json:"is_experience"`
			ColumnCtime       int  `json:"column_ctime"`
			ColumnSku         int  `json:"column_sku"`
			ColumnGroupbuy    int  `json:"column_groupbuy"`
			LastChapterID     int  `json:"last_chapter_id"`
			InPvip            int  `json:"in_pvip"`
			IsChannel         int  `json:"is_channel"`
			ColumnPrice       int  `json:"column_price"`
			IsRealSub         bool `json:"is_real_sub"`
		} `json:"list"`
		Page struct {
			Count int `json:"count"`
		} `json:"page"`
	} `json:"data"`
	Code int `json:"code"`
}

func Skus(p PType) (*SkusResponse, error) {
	var result *SkusResponse
	var tp int
	switch p {
	case ProductTypeVideo:
		tp = 3
	case ProductTypeAll:
		tp = 0
	case ProductTypeZhuanlan:
		tp = 1
	}
	cacheKey := fmt.Sprintf("skus-%d", tp)
	file, err := c.Get(cacheKey)
	if err == nil && len(file) > 0 {
		err = json.NewDecoder(bytes.NewReader(file)).Decode(&result)
		if err == nil {
			return result, err
		}
	}
	//https://time.geekbang.org/serv/v1/column/label_skus
	res, err := HttpClient.Post("https://time.geekbang.org/serv/v1/column/label_skus", fmt.Sprintf(`{"label_id":0,"type":%d}`, tp), false)
	if err != nil {
		return nil, err
	}
	defer func() {
		io.Copy(io.Discard, res.Body)
		res.Body.Close()
	}()
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	if res.StatusCode < 400 {
		c.Set(cacheKey, result)
	}

	return result, nil
}

func AllProducts(t PType) ([]Product, error) {
	var results []Product
	page := 1
	for page > 0 {
		products, err := Products(page, 100, t)
		if err != nil {
			return nil, err
		}
		results = append(results, products.Data.Products...)
		if products.Data.Page.More {
			page++
		} else {
			page = -1
		}
	}
	return results, nil
}

func Products(prev, size int, t PType) (ProjectResponse, error) {
	var result ProjectResponse

	res, err := HttpClient.Post("https://time.geekbang.org/serv/v3/learn/product", fmt.Sprintf(`{"desc":true,"expire":1,"last_learn":0,"learn_status":0,"prev":%d,"size":%d,"sort":1,"type":"%s","with_learn_count":1}`, prev, size, t), false)
	if err != nil {
		return ProjectResponse{}, err
	}
	defer func() {
		io.Copy(io.Discard, res.Body)
		res.Body.Close()
	}()
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return ProjectResponse{}, err
	}

	if result.Code == -1 {
		return ProjectResponse{}, errors.New("再等等吧, 不让抓了")
	}

	return result, nil
}

type Video struct {
	Sd struct {
		URL  string `json:"url"`
		Size int    `json:"size"`
	} `json:"sd"`
	Hd struct {
		URL  string `json:"url"`
		Size int    `json:"size"`
	} `json:"hd"`
	Ld struct {
		URL  string `json:"url"`
		Size int    `json:"size"`
	} `json:"ld"`
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
		ArticleCouldPreview bool        `json:"article_could_preview"`
		HlsVideos           interface{} `json:"hls_videos"`
		InPvip              int         `json:"in_pvip"`
		AudioDownloadURL    string      `json:"audio_download_url"`
		ArticleCtime        int         `json:"article_ctime"`
		ArticleSharetitle   string      `json:"article_sharetitle"`
	} `json:"data"`
	Code int `json:"code"`
}

func DeleteCache(key string) {
	c.Delete(key)
}

func DeleteArticleCache(id string) {
	DeleteCache("article-" + id)
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

func DeleteArticlesCache(cid int) {
	key := fmt.Sprintf("articles-%d", cid)
	DeleteCache(key)
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
		fmt.Sprintf(`{"cid":%d,"size":500,"prev":0,"order":"earliest","sample":false}`, cid), false)
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

func VideoKey(u string, vid string) ([]byte, error) {
	cacheKey := "keyurl-" + vid
	file, err := c.Get(cacheKey)
	if err == nil {
		return file, nil
	}
	request, _ := http.NewRequest("GET", u, nil)
	request.Header.Set("origin", "https://time.geekbang.org")

	get, err := (&http.Client{}).Do(request)
	if err != nil {
		return nil, err
	}
	defer get.Body.Close()
	all, _ := io.ReadAll(get.Body)
	if get.ContentLength > 0 {
		c.SetOrigin(cacheKey, all)
	}
	if get.StatusCode != 200 {
		log.Printf("video key response code != 200, data: '%s', code: %d\n", string(all), get.StatusCode)
	}
	return all, nil
}
