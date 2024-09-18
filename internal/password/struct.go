package password

// Config struct to configure password hashing.
type Config struct {
	Time      uint32
	Memory    uint32
	Threads   uint8
	KeyLength uint32
}
