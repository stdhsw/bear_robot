package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Server struct {
	port    string
	handler Handler

	router *gin.Engine
	srv    *http.Server
}

func NewServer(handler Handler, port string) (*Server, error) {
	if port == "" {
		port = "8080"
	}

	s := &Server{
		port:    port,
		handler: handler,
	}

	return s, s.init()
}

func (s *Server) init() error {
	s.router = gin.Default()

	s.router.POST("/create", func(c *gin.Context) {
		var accountReq AccountRequest
		if err := c.ShouldBindJSON(&accountReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := s.check(accountReq.Account, accountReq.Password); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := s.handler.Create(&accountReq); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{})
	})

	s.router.POST("/history", func(c *gin.Context) {
		var accountReq AccountRequest
		if err := c.ShouldBindJSON(&accountReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := s.check(accountReq.Account, accountReq.Password); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		accountInfo, err := s.handler.History(&accountReq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, accountInfo)
	})

	s.router.POST("/deposit", func(c *gin.Context) {
		var accountReq AccountRequest
		if err := c.ShouldBindJSON(&accountReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := s.check(accountReq.Account, accountReq.Password); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := s.handler.Deposit(&accountReq); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{})
	})

	s.router.POST("/withdraw", func(c *gin.Context) {
		var accountReq AccountRequest
		if err := c.ShouldBindJSON(&accountReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := s.check(accountReq.Account, accountReq.Password); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := s.handler.Withdraw(&accountReq); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{})
	})

	return nil
}

func (s *Server) check(account, pass string) error {
	if account == "" {
		return fmt.Errorf("account is empty")
	} else if pass == "" {
		return fmt.Errorf("password is empty")
	}

	return nil
}

func (s *Server) Start() error {
	s.srv = &http.Server{
		Addr:    ":" + s.port,
		Handler: s.router,
	}

	if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (s *Server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	s.srv.Shutdown(ctx)
}
