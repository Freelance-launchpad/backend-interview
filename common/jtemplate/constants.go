package jtemplate

import (
	_ "embed"
	"encoding/base64"
	"html/template"
)

var (
	//go:embed assets/logo_jump.svg
	LogoJump []byte

	//go:embed assets/icon-external-link.svg
	IconExternalLink []byte

	//go:embed assets/illustration-file.png
	illustrationFile []byte
	IllustrationFile template.URL

	//go:embed assets/signature_nicolas_fayon.svg
	SignatureNicolasFayon []byte

	//go:embed assets/stamp_jump_blue.svg
	StampJumpBlue []byte

	//go:embed assets/stamp_jump_green.svg
	StampJumpGreen []byte

	//go:embed assets/stamp_jump_realty.svg
	StampJumpRealty []byte

	//go:embed assets/stamp_jump_formation.svg
	StampJumpFormation []byte

	//go:embed assets/tailwind.dist.css
	Stylesheet []byte
)

func init() {
	IllustrationFile = template.URL(
		"data:image/png;base64," +
			base64.StdEncoding.EncodeToString(illustrationFile),
	)
}
