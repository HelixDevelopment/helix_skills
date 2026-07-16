package api

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
	"go.uber.org/zap"
)

// setupHTTP3 creates and starts an HTTP/3 server alongside the HTTP/2 server.
// It requires a valid TLS certificate; HTTP/3 cannot run without TLS.
func (s *Server) setupHTTP3(handler http.Handler) error {
	if s.cfg.TLSCert == "" || s.cfg.TLSKey == "" {
		return fmt.Errorf("HTTP/3 requires TLS certificate and key")
	}

	// Load TLS certificate
	cert, err := tls.LoadX509KeyPair(s.cfg.TLSCert, s.cfg.TLSKey)
	if err != nil {
		return fmt.Errorf("failed to load TLS certificates for HTTP/3: %w", err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		// Prefer modern TLS versions
		MinVersion: tls.VersionTLS13,
	}

	// Configure QUIC for optimal performance.
	// 0-RTT is a QUIC-level setting in quic-go (quic.Config.Allow0RTT),
	// not a crypto/tls.Config field.
	quicConfig := &quic.Config{
		MaxIdleTimeout:        30 * time.Second,
		HandshakeIdleTimeout:  10 * time.Second,
		MaxIncomingStreams:    100,
		MaxIncomingUniStreams: 100,
		// Enable 0-RTT for faster session resumption handshakes.
		Allow0RTT: true,
	}

	addr := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.HTTP3Port)
	s.http3Server = &http3.Server{
		Addr:       addr,
		Handler:    handler,
		TLSConfig:  tlsConfig,
		QUICConfig: quicConfig,
	}

	// Start HTTP/3 server in a goroutine
	go func() {
		zap.L().Info("starting HTTP/3 server",
			zap.String("addr", addr),
		)
		if err := s.http3Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Error("HTTP/3 server error", zap.Error(err))
		}
	}()

	return nil
}

// shutdownHTTP3 gracefully shuts down the HTTP/3 server.
func (s *Server) shutdownHTTP3(ctx context.Context) error {
	if s.http3Server == nil {
		return nil
	}

	zap.L().Info("shutting down HTTP/3 server")

	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		done <- s.http3Server.Close()
	}()

	select {
	case err := <-done:
		if err != nil {
			return fmt.Errorf("HTTP/3 shutdown error: %w", err)
		}
		zap.L().Info("HTTP/3 server stopped")
		return nil
	case <-shutdownCtx.Done():
		return fmt.Errorf("HTTP/3 shutdown timed out")
	}
}
