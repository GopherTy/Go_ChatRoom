package configure

import (
	"path/filepath"
	"strings"
)

// GRPC configure grpc
type GRPC struct {
	Addr string

	CertFile string
	KeyFile  string
}

// H2 if tls return true
func (c *GRPC) H2() bool {
	return c.CertFile != "" && c.KeyFile != ""
}

// H2C if not use tls return true
func (c *GRPC) H2C() bool {
	return c.CertFile == "" || c.KeyFile == ""
}

// Format .
func (c *GRPC) Format(basePath string) (e error) {
	c.Addr = strings.TrimSpace(c.Addr)
	c.CertFile = strings.TrimSpace(c.CertFile)
	c.KeyFile = strings.TrimSpace(c.KeyFile)

	if c.H2() {
		if filepath.IsAbs(c.CertFile) {
			c.CertFile = filepath.Clean(c.CertFile)
		} else {
			c.CertFile = filepath.Clean(basePath + "/" + c.CertFile)
		}

		if filepath.IsAbs(c.KeyFile) {
			c.KeyFile = filepath.Clean(c.KeyFile)
		} else {
			c.KeyFile = filepath.Clean(basePath + "/" + c.KeyFile)
		}
	}
	return
}
