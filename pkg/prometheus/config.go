package prometheus

type Config struct {
	ListenAddress string `json:"listen_address"`
	ReadTimeout   int    `json:"read_timeout"`
}
