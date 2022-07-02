package jwt

type Config struct {
	SigningKey  string `yaml:"signing-key"`  // jwt签名
	ExpiresTime int64  `yaml:"expires-time"` // 过期时间
	BufferTime  int64  `yaml:"buffer-time"`  // 缓冲时间
	Issuer      string `yaml:"issuer"`       // 签发者
}
