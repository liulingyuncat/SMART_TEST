package prompts

import (
	"encoding/json"
	"fmt"

	"webtest/internal/mcp/transport"
)

// Notifier 提示词变更通知器
type Notifier struct {
	transport transport.Transport
}

// NewNotifier 创建新的Notifier
func NewNotifier(trans transport.Transport) *Notifier {
	return &Notifier{
		transport: trans,
	}
}

// NotifyPromptsChanged 向客户端发送prompts/list_changed通知
// 仅在stdio传输模式下有效
func (n *Notifier) NotifyPromptsChanged() error {
	notification := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "notifications/prompts/list_changed",
		"params":  map[string]interface{}{},
	}

	data, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}

	if err := n.transport.Send(data); err != nil {
		return fmt.Errorf("failed to send notification: %w", err)
	}

	return nil
}
