package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/narvdeshwar/IOT-rag/internal/cache"
	"github.com/narvdeshwar/IOT-rag/internal/embedder"
	"github.com/narvdeshwar/IOT-rag/internal/llm"
	"github.com/narvdeshwar/IOT-rag/internal/retriever"
)

type QueryRequest struct {
	Query string `json:"query"`
}

type QueryResponse struct {
	Answer string `json:"answer"`
	Cached bool   `json:"cached"`
}

type Service struct {
	Embedder  *embedder.OpenAIEmbedder
	Retriever *retriever.Retriever
	LLM       *llm.LLM
}

func QueryHandler(svc *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req QueryRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Cache check
		if cached, err := cache.Get(c.Request.Context(), req.Query); err == nil {
			c.JSON(http.StatusOK, QueryResponse{Answer: cached, Cached: true})
			return
		}

		// Embed query
		vec, err := svc.Embedder.Embed(c.Request.Context(), req.Query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Retrieve
		chunks, err := svc.Retriever.Search(c.Request.Context(), vec, 10)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Build prompt
		contextStr := ""
		for _, ch := range chunks {
			contextStr += ch.Content + "\n"
		}
		prompt := fmt.Sprintf("Context:\n%s\n\nQuestion: %s", contextStr, req.Query)

		// Stream
		c.Header("Content-Type", "text/event-stream")
		c.Header("Cache-Control", "no-cache")
		c.Header("Connection", "keep-alive")

		full, err := svc.LLM.StreamComplete(c.Request.Context(), prompt, func(token string) {
			fmt.Fprintf(c.Writer, "data: %s\n\n", token)
			c.Writer.Flush()
		})
		if err != nil {
			return
		}

		// Cache
		cache.Set(c.Request.Context(), req.Query, full, 5*time.Minute)
	}
}
