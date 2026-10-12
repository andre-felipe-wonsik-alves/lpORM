package form

import _ "embed"

//go:embed style/global.css
var defaultCSS string

// CSS devolve a style sheet padrão dos formulários gerados
func CSS() string {
	return defaultCSS
}
