package init

import (
	"path"

	"github.com/go-macaron/cache"
	_ "github.com/go-macaron/cache/ledis"
	"github.com/go-macaron/captcha"
	"github.com/go-macaron/session"
	_ "github.com/jicg/easypos/modules/log"
	mlog "github.com/jicg/easypos/modules/log"
	temfunc "github.com/jicg/easypos/modules/temfuc"
	"gopkg.in/macaron.v1"
)

const (
	CACHE_DIR          = "./data/cache.db,db=0"
	LOG_PATH           = "./log/easypos.log"
	XORM_LOG_PATH      = "./log/xorm.log"
	DEFAULT_LOG_PREFIX = "[easypos]"
	HTML_STATIC        = "static"
)

var (
	Logger *mlog.Logger
	Evn    = 1
)

func NewSession() macaron.Handler {
	return session.Sessioner(session.Options{
		Provider:       "file",
		ProviderConfig: "data/sessions",
		CookieName:     "easypos",
	})
}

func NewCache() macaron.Handler {
	var op cache.Options
	op = cache.Options{
		Adapter:       "ledis",
		AdapterConfig: "data_dir=" + CACHE_DIR,
	}
	return cache.Cacher(op)
}

func NewStatic(dir string, pre string) macaron.Handler {
	return macaron.Static(dir, macaron.StaticOptions{
		// 请求静态资源时的 URL 前缀，默认没有前缀
		Prefix: dir,
		// 禁止记录静态资源路由日志，默认为不禁止记录
		SkipLogging: true,
		// 当请求目录时的默认索引文件，默认为 "index.html"
		//IndexFile: "index.html",
		// 用于返回自定义过期响应头，默认为不设置
		// https://developers.google.com/speed/docs/insights/LeverageBrowserCaching
		//Expires: func() string {
		//	return time.Now().Add(24 * 60 * time.Minute).UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT")
		//},
	})
}

func NewCap() macaron.Handler {
	return captcha.Captchaer(captcha.Options{
		// 获取验证码图片的 URL 前缀，默认为 "/captcha/"
		URLPrefix: "/captcha/",
		// 表单隐藏元素的 ID 名称，默认为 "captcha_id"
		FieldIdName: "captcha_id",
		// 用户输入验证码值的元素 ID，默认为 "captcha"
		FieldCaptchaName: "captcha",
		// 验证字符的个数，默认为 6
		ChallengeNums: 6,
		// 验证码图片的宽度，默认为 240 像素
		Width: 240,
		// 验证码图片的高度，默认为 80 像素
		Height: 80,
		// 验证码过期时间，默认为 600 秒
		Expiration: 600,
		// 用于存储验证码正确值的 Cache 键名，默认为 "captcha_"
		CachePrefix: "hw_",
	})
}

func NewRender() macaron.Handler {
	return macaron.Renderer(macaron.RenderOptions{
		Directory:  path.Join("views"),
		Extensions: []string{".tmpl", ".html"},
		Funcs:      temfunc.NewFuncMap(),
		Delims:     macaron.Delims{"{{", "}}"},
		// 追加的 Content-Type 头信息，默认为 "UTF-8"
		Charset:    "UTF-8",
		IndentJSON: macaron.Env != macaron.PROD,
	})
}
