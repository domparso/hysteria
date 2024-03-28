package cmd

import (
	"go.uber.org/zap"
	"net"
	"net/http"
	"strconv"
)

func (l *ServeLogger) Connect(addr net.Addr, id string, tx uint64) {
	Logger("info", "client connected", zap.String("addr", addr.String()), zap.String("id", id), zap.Uint64("tx", tx))
}

func (l *ServeLogger) Disconnect(addr net.Addr, id string, err error) {
	Logger("info", "client disconnected", zap.String("addr", addr.String()), zap.String("id", id), zap.Error(err))
}

func (l *ServeLogger) TCPRequest(addr net.Addr, id, reqAddr string) {
	Logger("debug", "TCP request", zap.String("addr", addr.String()), zap.String("id", id), zap.String("reqAddr", reqAddr))
}

func (l *ServeLogger) TCPError(addr net.Addr, id, reqAddr string, err error) {
	if err == nil {
		Logger("debug", "TCP closed", zap.String("addr", addr.String()), zap.String("id", id), zap.String("reqAddr", reqAddr))
	} else {
		Logger("error", "TCP error", zap.String("addr", addr.String()), zap.String("id", id), zap.String("reqAddr", reqAddr), zap.Error(err))
	}
}

func (l *ServeLogger) UDPRequest(addr net.Addr, id string, sessionID uint32, reqAddr string) {
	Logger("debug", "UDP request", zap.String("addr", addr.String()), zap.String("id", id), zap.Uint32("sessionID", sessionID), zap.String("reqAddr", reqAddr))
}

func (l *ServeLogger) UDPError(addr net.Addr, id string, sessionID uint32, err error) {
	if err == nil {
		Logger("debug", "UDP closed", zap.String("addr", addr.String()), zap.String("id", id), zap.Uint32("sessionID", sessionID))
	} else {
		Logger("error", "UDP error", zap.String("addr", addr.String()), zap.String("id", id), zap.Uint32("sessionID", sessionID), zap.Error(err))
	}
}

type masqHandlerLogWrapperS struct {
	H    http.Handler
	QUIC bool
}

func (m *masqHandlerLogWrapperS) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	Logger("Debug", "masquerade request",
		zap.String("addr", r.RemoteAddr),
		zap.String("method", r.Method),
		zap.String("host", r.Host),
		zap.String("url", r.URL.String()),
		zap.Bool("quic", m.QUIC))
	m.H.ServeHTTP(w, r)
}

func extractPortFromAdd(addr string) int {
	_, portStr, err := net.SplitHostPort(addr)
	if err != nil {
		return 0
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return 0
	}
	return port
}
