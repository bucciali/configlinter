package rules

import (
	"configlinter/internal/domain"
	"strings"
)

var weakAlgorithms = map[string]string{
	"md5":  "MD5 — коллизии находятся за секунды",
	"sha1": "SHA1 — коллизии доказаны (SHAttered, 2017)",
	"des":  "DES — ключ 56 бит, брутфорсится",
	"3des": "3DES — уязвим к Sweet32",
	"rc4":  "RC4 — статистические уязвимости, запрещён в TLS 1.3",
	"rc2":  "RC2 — устаревший, небезопасен",
}

type WeakCryptoRule struct{}

func NewWeakCryptoRule() *WeakCryptoRule {
	return &WeakCryptoRule{}
}

func (r *WeakCryptoRule) ID() string {
	return "weak_crypto"
}

func (r *WeakCryptoRule) Description() string {
	return "Обнаружение слабых криптографических алгоритмов"
}

func (r *WeakCryptoRule) Analyze(root *domain.ConfigNode) []domain.Finding {
	var findings []domain.Finding

	root.Walk(func(node *domain.ConfigNode) {
		val, ok := node.StringValue()
		if !ok {
			return
		}

		lower := strings.ToLower(val)

		reason, found := weakAlgorithms[lower]
		if !found {
			return
		}

		findings = append(findings, domain.Finding{
			RuleID:        r.ID(),
			Severity:      domain.HIGH,
			Path:          node.Path,
			Message:       "Слабый алгоритм: " + reason,
			Recomendation: "Используйте SHA-256+, AES-256, ChaCha20",
		})
	})

	return findings
}
