package captcha

import (
	"image/color"
	"strconv"
	"sync"
	"time"

	logger "github.com/gbrayhan/microservices-go/src/infrastructure/lib/logger"
	"github.com/gbrayhan/microservices-go/src/shared/utils"
	"github.com/mojocn/base64Captcha"
)

// Captcha 验证码服务结构体
type Captcha struct {
	store  *MemoryStore
	driver base64Captcha.Driver // 修改为接口类型
	config Config
}

// Config 验证码配置
type Config struct {
	Width      int           `json:"width"`  // 图片宽度
	Height     int           `json:"height"` // 图片高度
	Length     int           `json:"length"` // 验证码长度
	Timeout    time.Duration `json:"-"`      // 过期时间
	Complexity int           `json:"-"`      // 复杂度 (0-2)
}

// CaptchaResponse 验证码响应
type CaptchaResponse struct {
	ID     string `json:"id"`   // 验证码ID
	B64s   string `json:"b64s"` // base64编码的图片
	Answer string `json:"-"`    // 答案（不对外暴露）
	Config Config `json:"config"`
}

// MemoryStore 内存存储实现
type MemoryStore struct {
	lock   sync.RWMutex
	data   map[string]string
	expire time.Duration
}

// NewMemoryStore 创建内存存储
func NewMemoryStore(expire time.Duration) *MemoryStore {
	ms := &MemoryStore{
		data:   make(map[string]string),
		expire: expire,
	}

	// 启动清理过期数据的goroutine
	go ms.gc()

	return ms
}

// Set 存储验证码
func (ms *MemoryStore) Set(id string, value string) {
	ms.lock.Lock()
	defer ms.lock.Unlock()
	ms.data[id] = value
}

// Get 获取验证码
func (ms *MemoryStore) Get(id string, clear bool) string {
	ms.lock.RLock()
	defer ms.lock.RUnlock()

	if value, ok := ms.data[id]; ok {
		if clear {
			delete(ms.data, id)
		}
		return value
	}
	return ""
}

// gc 清理过期数据
func (ms *MemoryStore) gc() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		// 内存存储不主动清理，由验证时自动清理
		// 这里可以扩展为带过期时间的存储
	}
}

// DefaultConfig 默认配置
func DefaultConfig(loggerInstance *logger.Logger) Config {

	width, err := strconv.Atoi(utils.GetEnv("CAPTCHA_WIDTH", "120"))
	if err != nil {
		loggerInstance.Error("error parse captcha width")
	}
	height, err := strconv.Atoi(utils.GetEnv("CAPTCHA_HEIGHT", "50"))
	if err != nil {
		loggerInstance.Error("error parse captcha height")
	}
	timeout, err := strconv.Atoi(utils.GetEnv("CAPTCHA_TIMEOUT", "5"))
	if err != nil {
		loggerInstance.Error("error parse captcha timeout")
	}
	length, err := strconv.Atoi(utils.GetEnv("CAPTCHA_LENGTH", "4"))
	if err != nil {
		loggerInstance.Error("error parse captcha length")
	}
	complexity, err := strconv.Atoi(utils.GetEnv("CAPTCHA_COMPLEXITY", "2"))
	if err != nil {
		loggerInstance.Error("error parse captcha complexity")

	}
	return Config{
		Width:      width,
		Height:     height,
		Length:     length,
		Timeout:    time.Duration(timeout) * time.Minute,
		Complexity: complexity,
	}
}

// New 创建验证码服务实例
func New(config Config) *Captcha {
	var driver base64Captcha.Driver

	// 根据复杂度调整
	switch config.Complexity {
	case 0:
		driver = base64Captcha.NewDriverDigit(
			config.Height,
			config.Width,
			config.Length,
			0.5,
			40,
		)
	case 1:
		driver = base64Captcha.NewDriverDigit(
			config.Height,
			config.Width,
			config.Length,
			0.7,
			80,
		)
	case 2:
		// 提供完整参数给 NewDriverMath
		driver = base64Captcha.NewDriverMath(
			config.Height,
			config.Width,
			0,                                     // 随机干扰点数量
			0,                                     // 随机干扰线数量
			&color.RGBA{R: 0, G: 0, B: 0, A: 255}, // 字体颜色
			base64Captcha.DefaultEmbeddedFonts,    // 字体存储
			[]string{"wqy-microhei.ttc"},          // 字体文件列表
		)
	}
	return &Captcha{
		store:  NewMemoryStore(config.Timeout),
		driver: driver,
		config: config,
	}
}

// Generate 生成验证码
func (c *Captcha) Generate() *CaptchaResponse {
	// 使用自定义存储生成验证码
	id, content, answer := c.driver.GenerateIdQuestionAnswer()

	// 存储答案
	c.store.Set(id, answer)

	// 生成图片
	item, err := c.driver.DrawCaptcha(content)
	if err != nil {
		return nil
	}

	// 转换为base64
	b64s := item.EncodeB64string()

	return &CaptchaResponse{
		ID:     id,
		B64s:   b64s,
		Answer: answer,
		Config: c.config,
	}
}

// Verify 验证验证码
func (c *Captcha) Verify(id, answer string) bool {
	if id == "" || answer == "" {
		return false
	}

	// 从存储中获取正确答案
	correctAnswer := c.store.Get(id, true) // 验证后清除
	if correctAnswer == "" {
		return false
	}

	// 比较答案（忽略大小写）
	return correctAnswer == answer
}

// Refresh 刷新验证码（重新生成）
func (c *Captcha) Refresh(id string) *CaptchaResponse {
	// 清除旧的验证码
	c.store.Get(id, true)

	// 生成新的验证码
	return c.Generate()
}
