package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/narvdeshwar/IOT-rag/internal/cache"
	"github.com/narvdeshwar/IOT-rag/internal/telemetry"
	"github.com/narvdeshwar/IOT-rag/internal/types"
)

type QueryRequest struct {
	Query string `json:"query"`
}

type QueryResponse struct {
	Answer string `json:"answer"`
	Cached bool   `json:"cached"`
}

type Service struct {
	Embedder  types.Embedder
	Retriever Retriever
	LLM       types.LLM
}

type Retriever interface {
	Search(ctx context.Context, vec []float32, k int) ([]telemetry.EmbeddedChunk, error)
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

		fmt.Printf("\nDEBUG: Found %d relevant chunks in DB\n", len(chunks))

		// Build prompt
		contextStr := ""
		for i, ch := range chunks {
			contextStr += ch.Content + "\n"
			fmt.Printf("Chunk %d: %s\n", i, ch.Content)
		}
		prompt := fmt.Sprintf("Context:\n%s\n\nQuestion: %s", contextStr, req.Query)
		
		fmt.Println("\nDEBUG: Final Prompt sent to LLM:")
		fmt.Println(prompt)

		// Stream
		c.Header("Content-Type", "text/event-stream")
		c.Header("Cache-Control", "no-cache")
		c.Header("Connection", "keep-alive")

		full, err := svc.LLM.StreamComplete(c.Request.Context(), prompt, func(token string) {
			b, _ := json.Marshal(map[string]string{"token": token})
			fmt.Fprintf(c.Writer, "data: %s\n\n", string(b))
			c.Writer.Flush()
		})
		if err != nil {
			return
		}

		// Cache
		cache.Set(c.Request.Context(), req.Query, full, 5*time.Minute)
	}
}
