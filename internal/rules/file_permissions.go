package rules

import (
	"configlinter/internal/domain"
	"fmt"
)

type FilePermissionsRule struct{}

func (r *FilePermissionsRule) ID() string { return "world_readable_config" }
func (r *FilePermissionsRule) Description() string {
	return "Проверка прав доступа к файлу"
}

func (r *FilePermissionsRule) Analyze(root *domain.ConfigNode) []domain.Finding {
	if root.FileMode == 0 {
		return nil
	}

	var findings []domain.Finding

	if root.FileMode&0o002 != 0 {
		findings = append(findings, domain.Finding{
			RuleID:        r.ID(),
			Severity:      domain.HIGH,
			Path:          root.FilePath,
			Message:       fmt.Sprintf("Конфиг доступен на запись всем (permissions: %s)", root.FileMode),
			Recomendation: "Установите chmod 600 или 640 для конфигурационных файлов",
		})
	} else if root.FileMode&0o004 != 0 {
		findings = append(findings, domain.Finding{
			RuleID:        r.ID(),
			Severity:      domain.MEDIUM,
			Path:          root.FilePath,
			Message:       fmt.Sprintf("Конфиг доступен на чтение всем (permissions: %s)", root.FileMode),
			Recomendation: "Установите chmod 600 или 640 для конфигурационных файлов",
		})
	}

	return findings
}
