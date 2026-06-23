// Package verify checks whether a candidate proxy is actually reachable and, where
// possible, that it speaks the expected protocol — measuring latency along the way.
//
// We deliberately do NOT claim a proxy "works end to end" (that would require relaying
// a full MTProto session to a Telegram DC). We report exactly what we observed:
//
//	reachable    - DNS resolved and a TCP connection was established (RTT measured)
//	handshake_ok - additionally completed a TLS handshake for an ee/FakeTLS secret
package verify

import (
	"context"
	"crypto/tls"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/zakky8/mtproto-proxy-pro/internal/model"
)

// Result is the outcome of verifying one proxy.
type Result struct {
	OK        bool
	Status    string
	LatencyMS int
	IP        net.IP // resolved address, for geo lookup
}

// Options tunes a verification run.
type Options struct {
	Concurrency int
	Timeout     time.Duration
}

func (o Options) withDefaults() Options {
	if o.Concurrency <= 0 {
		o.Concurrency = 150
	}
	if o.Timeout <= 0 {
		o.Timeout = 6 * time.Second
	}
	return o
}

// Many verifies proxies on a bounded worker pool. The returned slice is index-aligned
// with the input. progress, if non-nil, is called after each check with (done, total).
func Many(ctx context.Context, proxies []model.Proxy, opts Options, progress func(done, total int)) []Result {
	opts = opts.withDefaults()
	results := make([]Result, len(proxies))
	sem := make(chan struct{}, opts.Concurrency)
	var wg sync.WaitGroup
	var done int64

	for i := range proxies {
		wg.Add(1)
		sem <- struct{}{}
		go func(i int) {
			defer wg.Done()
			defer func() { <-sem }()
			results[i] = One(ctx, proxies[i], opts.Timeout)
			if progress != nil {
				progress(int(atomic.AddInt64(&done, 1)), len(proxies))
			}
		}(i)
	}
	wg.Wait()
	return results
}

// One verifies a single proxy.
func One(ctx context.Context, p model.Proxy, timeout time.Duration) Result {
	ip := resolve(ctx, p.Server, timeout)
	if ip == nil {
		return Result{}
	}
	addr := net.JoinHostPort(ip.String(), strconv.Itoa(p.Port))

	dialer := net.Dialer{Timeout: timeout}
	start := time.Now()
	conn, err := dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		return Result{}
	}
	rtt := time.Since(start)
	defer conn.Close()

	res := Result{
		OK:        true,
		Status:    model.StatusReachable,
		LatencyMS: int(rtt.Milliseconds()),
		IP:        ip,
	}

	// For FakeTLS (ee) proxies, a completed TLS handshake to the fronted SNI is a strong
	// signal the proxy is live and speaking the protocol.
	if domain := p.FakeTLSDomain(); domain != "" {
		if tlsHandshake(ctx, conn, domain, timeout) {
			res.Status = model.StatusHandshakeOK
		}
	}
	return res
}

func resolve(ctx context.Context, host string, timeout time.Duration) net.IP {
	host = trimDot(host)
	if ip := net.ParseIP(host); ip != nil {
		return ip
	}
	rctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	ips, err := net.DefaultResolver.LookupIP(rctx, "ip", host)
	if err != nil || len(ips) == 0 {
		return nil
	}
	for _, ip := range ips { // prefer IPv4 for geo + connectivity
		if ip.To4() != nil {
			return ip
		}
	}
	return ips[0]
}

func tlsHandshake(ctx context.Context, conn net.Conn, sni string, timeout time.Duration) bool {
	_ = conn.SetDeadline(time.Now().Add(timeout))
	defer conn.SetDeadline(time.Time{})
	tc := tls.Client(conn, &tls.Config{
		ServerName:         sni,
		InsecureSkipVerify: true, // we only care that it completes a TLS handshake
		MinVersion:         tls.VersionTLS12,
	})
	return tc.HandshakeContext(ctx) == nil
}

func trimDot(s string) string {
	for len(s) > 0 && s[len(s)-1] == '.' {
		s = s[:len(s)-1]
	}
	return s
}
