package llm

import (
	"context"
	"io"

	openai "github.com/sashabaranov/go-openai"
)

type LLM struct {
	client *openai.Client
}

func NewLLM() *LLM {
	return &LLM{client: openai.NewClient("")}
}

func (l *LLM) StreamComplete(ctx context.Context, prompt string, onToken func(string)) (string, error) {
	req := openai.ChatCompletionRequest{
		Model: openai.GPT4oMini,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: SystemPrompt},
			{Role: openai.ChatMessageRoleUser, Content: prompt},
		},
		Stream: true,
	}

	stream, err := l.client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		return "", err
	}
	defer stream.Close()

	var full string
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
		token := resp.Choices[0].Delta.Content
		full += token
		onToken(token)
	}
	return full, nil
}
