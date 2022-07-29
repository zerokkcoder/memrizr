package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"memrizr/model/apperrors"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

func Timeout(timeout time.Duration, errTimeout *apperrors.Error) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置gin writer 为 自定义的writer
		tw := &timeoutWriter{ResponseWriter: c.Writer, h: make(http.Header)}
		c.Writer = tw

		// 包装带超时的上下文
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		// 更新 gin 请求上下文
		c.Request = c.Request.WithContext(ctx)

		finished := make(chan struct{})        // 表示处理程序已完成
		panicChan := make(chan interface{}, 1) // 用于在我们无法恢复的情况下处理恐慌

		go func() {
			defer func() {
				if p := recover(); p != nil {
					panicChan <- p
				}
			}()

			c.Next() // 调用后续的中间件和处理程序
			finished <- struct{}{}
		}()

		select {
		case <-panicChan:
			// 如果我们无法从恐慌中恢复，请发送内部服务器错误
			e := apperrors.NewInternal()
			tw.ResponseWriter.WriteHeader(e.Status())
			eResp, _ := json.Marshal(gin.H{
				"error": e,
			})
			tw.ResponseWriter.Write(eResp)
		case <-finished:
			// 如果完成，设置 header 和 写入响应
			tw.mu.Lock()
			defer tw.mu.Unlock()
			// 将 Headers 从tw.Header()(由gin写入)映射到tw.ResponseWriter以进行响应
			dst := tw.ResponseWriter.Header()
			for k, vv := range tw.Header() {
				dst[k] = vv
			}
			tw.ResponseWriter.WriteHeader(tw.code)
			// 当 gin 写入 tw.Write() 时，tw.wbuf 将已经被写入
			tw.ResponseWriter.Write(tw.wbuf.Bytes())
		case <-ctx.Done():
			// 超时已经发生，发送 超时错误 和 写 headers
			tw.mu.Lock()
			defer tw.mu.Unlock()
			// 从 gin 响应写入
			tw.ResponseWriter.Header().Set("Content-Type", "application/json")
			tw.ResponseWriter.WriteHeader(errTimeout.Status())
			eResp, _ := json.Marshal(gin.H{
				"error": errTimeout,
			})
			tw.ResponseWriter.Write(eResp)
			c.Abort()
			tw.SetTimedOut()
		}
	}
}

// implements http.Writer, but tracks if Writer has timed out
// or has already written its header to prevent
// header and body overwrites
// also locks access to this writer to prevent race conditions
// holds the gin.ResponseWriter which we'll manually call Write()
// on in the middleware function to send response
type timeoutWriter struct {
	gin.ResponseWriter
	h    http.Header
	wbuf bytes.Buffer // The zero value for Buffer is an empty buffer ready to use.

	mu          sync.Mutex
	timeOut     bool
	wroteHeader bool
	code        int
}

// Writes the response, but first makes sure there
// hasn't already been a timeout
// In http.ResponseWriter interface
func (tw *timeoutWriter) Write(b []byte) (int, error) {
	tw.mu.Lock()
	defer tw.mu.Unlock()
	if tw.timeOut {
		return 0, nil
	}

	return tw.wbuf.Write(b)
}

// In http.ResponseWriter interface
func (tw *timeoutWriter) WriteHeader(code int) {
	checkWriteHeaderCode(code)
	tw.mu.Lock()
	defer tw.mu.Unlock()
	// We do not write the header if we've timed out or written the header
	if tw.timeOut || tw.wroteHeader {
		return
	}
	tw.writeHeader(code)
}

// set that the header has been written
func (tw *timeoutWriter) writeHeader(code int) {
	tw.wroteHeader = true
	tw.code = code
}

// Header "relays" the header, h, set in struct
// In http.ResponseWriter interface
func (tw *timeoutWriter) Header() http.Header {
	return tw.h
}

// SetTimeOut sets timedOut field to true
func (tw *timeoutWriter) SetTimedOut() {
	tw.timeOut = true
}

func checkWriteHeaderCode(code int) {
	if code < 100 || code > 999 {
		panic(fmt.Sprintf("invalid WriteHeader code %v", code))
	}
}
