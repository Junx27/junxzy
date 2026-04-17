package generator

import (
	"fmt"
	"os"
	"strings"
)

func InjectRoute(moduleName string) error {
	filePath := "router/router.go"

	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	content := string(data)
	lower := strings.ToLower(moduleName)

	modulePath := GetModulePath()

	importLine := fmt.Sprintf(`"%s/modules/%s/route"`, modulePath, lower)
	registerLine := fmt.Sprintf("\troute.Register%sRoutes(r)", capitalize(moduleName))

	// inject import
	if !strings.Contains(content, importLine) {
		content = strings.Replace(
			content,
			"import (",
			"import (\n\t"+importLine,
			1,
		)
	}

	// inject register
	if !strings.Contains(content, registerLine) {
		content = strings.Replace(
			content,
			"func RegisterRoutes(r *gin.Engine) {",
			"func RegisterRoutes(r *gin.Engine) {\n"+registerLine,
			1,
		)
	}

	return os.WriteFile(filePath, []byte(content), os.ModePerm)
}
