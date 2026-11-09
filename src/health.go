package main

import (
	"crypto/subtle"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"runtime"
	"sort"
	"strings"
	"sync"
	"syscall"
	"time"
)

var startedAt = time.Now()

// pingTargets are dialled over TCP to measure outbound latency from the
// server (ICMP needs raw sockets, which the container does not have).
var pingTargets = []string{"1.1.1.1:443", "8.8.8.8:443"}

// reqWindow is how many recent requests the latency figures are computed over.
const reqWindow = 500

// reqStats holds a rolling window of handler durations plus lifetime counters.
var reqStats struct {
	sync.Mutex
	recent []time.Duration // ring buffer of the last reqWindow durations
	next   int
	total  int
	errors int // 5xx responses
}

// recordRequest is registered as a request observer. Health requests are
// skipped so polling the endpoint does not skew the numbers.
func recordRequest(path string, status int, d time.Duration) {
	if strings.HasPrefix(path, "/health") {
		return
	}
	reqStats.Lock()
	defer reqStats.Unlock()
	if len(reqStats.recent) < reqWindow {
		reqStats.recent = append(reqStats.recent, d)
	} else {
		reqStats.recent[reqStats.next] = d
		reqStats.next = (reqStats.next + 1) % reqWindow
	}
	reqStats.total++
	if status >= 500 {
		reqStats.errors++
	}
}

// requestSummary reports lifetime counts and avg/p95 handler time over the
// recent window. This is server time only, not what a visitor's network adds.
func requestSummary() map[string]any {
	reqStats.Lock()
	d := append([]time.Duration(nil), reqStats.recent...)
	total, errs := reqStats.total, reqStats.errors
	reqStats.Unlock()

	out := map[string]any{"total": total, "errors_5xx": errs, "window": len(d)}
	if len(d) == 0 {
		return out
	}
	sort.Slice(d, func(i, j int) bool { return d[i] < d[j] })
	var sum time.Duration
	for _, v := range d {
		sum += v
	}
	ms := func(v time.Duration) float64 { return float64(v.Microseconds()) / 1000 }
	out["avg_ms"] = ms(sum / time.Duration(len(d)))
	out["p95_ms"] = ms(d[len(d)*95/100])
	out["max_ms"] = ms(d[len(d)-1])
	return out
}

func init() {
	requestObservers = append(requestObservers, recordRequest)
	http.HandleFunc("/health", middlewareFunc(healthHandler))
	http.HandleFunc("/health/details", middlewareFunc(healthDetailsHandler))
}

// healthHandler is the public liveness probe. It reveals nothing but "up".
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Robots-Tag", "noindex")
	w.Write([]byte(`{"status":"ok"}` + "\n"))
}

// healthDetailsHandler serves server metrics behind basic auth. It fails
// closed: with HEALTH_USER or HEALTH_PASS unset the route does not exist.
func healthDetailsHandler(w http.ResponseWriter, r *http.Request) {
	user, pass := os.Getenv("HEALTH_USER"), os.Getenv("HEALTH_PASS")
	if user == "" || pass == "" {
		renderNotFound(w, r)
		return
	}
	u, p, ok := r.BasicAuth()
	good := subtle.ConstantTimeCompare([]byte(u), []byte(user)) &
		subtle.ConstantTimeCompare([]byte(p), []byte(pass))
	if !ok || good != 1 {
		w.Header().Set("WWW-Authenticate", `Basic realm="health"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	out := map[string]any{
		"status":         "ok",
		"time":           time.Now().UTC().Format(time.RFC3339),
		"uptime_seconds": int(time.Since(startedAt).Seconds()),
		"goroutines":     runtime.NumGoroutine(),
		"memory_mb":      mem.Alloc / 1024 / 1024,
		"load_avg":       loadAvg(),
		"disk_free_mb":   diskFreeMB("."),
		"requests":       requestSummary(),
		"render_ms":      renderMillis(),
		"ping_ms":        pingMillis(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Robots-Tag", "noindex")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.Encode(out)
}

// loadAvg returns the 1, 5 and 15 minute load averages (Linux only).
func loadAvg() []string {
	b, err := os.ReadFile("/proc/loadavg")
	if err != nil {
		return nil
	}
	f := strings.Fields(string(b))
	if len(f) < 3 {
		return nil
	}
	return f[:3]
}

func diskFreeMB(path string) uint64 {
	var s syscall.Statfs_t
	if err := syscall.Statfs(path, &s); err != nil {
		return 0
	}
	return uint64(s.Bavail) * uint64(s.Bsize) / 1024 / 1024
}

// renderMillis times a full render of the home page, excluding the network.
func renderMillis() float64 {
	start := time.Now()
	indexHandler(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))
	return float64(time.Since(start).Microseconds()) / 1000
}

// pingMillis maps each target to the mean of three TCP connect times, or -1
// if any attempt fails.
func pingMillis() map[string]float64 {
	const tries = 3
	res := map[string]float64{}
	for _, t := range pingTargets {
		var sum time.Duration
		for i := 0; i < tries; i++ {
			start := time.Now()
			c, err := net.DialTimeout("tcp", t, 2*time.Second)
			if err != nil {
				sum = -1
				break
			}
			c.Close()
			sum += time.Since(start)
		}
		if sum < 0 {
			res[t] = -1
			continue
		}
		res[t] = float64((sum / tries).Microseconds()) / 1000
	}
	return res
}
