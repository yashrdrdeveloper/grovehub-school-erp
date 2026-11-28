package impl_test

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net"
	"net/http"
	"sync/atomic"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
	mock "grovehub/launcher/internal/resources/tests/mocks/server"
	"grovehub/launcher/internal/src/server/impl"
)

// @Testing suite=LauncherServer
func ServerSpecs() {
	Describe("HTTP server lifecycle", func() {
		It("reports a bind error and gracefully drains a real listener on cancellation", func() {
			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(204)
			})
			draining := &atomic.Bool{}
			Expect(impl.New("not-a-listen-address", handler, logger, draining).Run(context.Background())).To(MatchError(ContainSubstring("bind HTTP listener")))
			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			Expect(impl.New("127.0.0.1:0", handler, logger, draining).Run(ctx)).To(Succeed())
			Expect(draining.Load()).To(BeTrue())
		})
		It("propagates serving errors and treats a closed server as normal termination", func() {
			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			for _, failure := range []error{nil, http.ErrServerClosed, errors.New("serve failure")} {
				listener, err := net.Listen("tcp", "127.0.0.1:0")
				Expect(err).NotTo(HaveOccurred())
				defer listener.Close()
				transport := mock.NewMockLifecycle(gomock.NewController(GinkgoT()))
				transport.EXPECT().Serve(listener).Return(failure)
				err = impl.RunLifecycle(context.Background(), listener, transport, logger, &atomic.Bool{})
				if failure == nil || errors.Is(failure, http.ErrServerClosed) {
					Expect(err).NotTo(HaveOccurred())
				} else {
					Expect(err).To(MatchError(ContainSubstring("serve HTTP: serve failure")))
				}
			}
		})
		It("forces close after a graceful shutdown failure and returns the original error", func() {
			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			listener, err := net.Listen("tcp", "127.0.0.1:0")
			Expect(err).NotTo(HaveOccurred())
			defer listener.Close()
			transport := mock.NewMockLifecycle(gomock.NewController(GinkgoT()))
			release, served := make(chan struct{}), make(chan struct{})
			transport.EXPECT().Serve(listener).DoAndReturn(func(net.Listener) error {
				<-release
				close(served)
				return http.ErrServerClosed
			})
			failure := errors.New("shutdown timeout")
			transport.EXPECT().Shutdown(gomock.Any()).Return(failure)
			transport.EXPECT().Close().DoAndReturn(func() error {
				close(release)
				return nil
			})
			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			draining := &atomic.Bool{}
			Expect(impl.RunLifecycle(ctx, listener, transport, logger, draining)).To(MatchError(failure))
			Eventually(served).Should(BeClosed())
			Expect(draining.Load()).To(BeTrue())
		})
	})
}
