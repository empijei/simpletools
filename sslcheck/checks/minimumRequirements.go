package checks

import "crypto/tls"

func CanConnect(hostPort string) (out Result) {
	conn, err := tls.Dial("tcp", hostPort, nil)
	if err == nil {
		_ = conn.Close()
	}
	return Result{err: err, hostPort: hostPort}
}
