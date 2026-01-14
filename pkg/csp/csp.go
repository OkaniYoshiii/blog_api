package csp

import (
	"fmt"
	"strings"
)

type CSP struct {
	ChildSrc      string
	ConnectSrc    string
	DefaultSrc    string
	FontSrc       string
	FrameSrc      string
	ImgSrc        string
	ManifestSrc   string
	MediaSrc      string
	ObjectSrc     string
	ScriptSrc     string
	ScriptSrcElem string
	ScriptSrcAttr string
	StyleSrc      string
	StyleSrcElem  string
	StyleSrcAttr  string
	WorkerSrc     string
}

func (csp *CSP) String() string {
	buffer := strings.Builder{}
	policyFmt := "%s '%s'; "

	directives := map[string]string{
		"child-src":       csp.ChildSrc,
		"connect-src":     csp.ConnectSrc,
		"default-src":     csp.DefaultSrc,
		"font-src":        csp.FontSrc,
		"frame-src":       csp.FrameSrc,
		"img-src":         csp.ImgSrc,
		"manifest-src":    csp.ManifestSrc,
		"media-src":       csp.MediaSrc,
		"object-src":      csp.ObjectSrc,
		"script-src":      csp.ScriptSrc,
		"script-src-elem": csp.ScriptSrcElem,
		"script-src-attr": csp.ScriptSrcAttr,
		"style-src":       csp.StyleSrc,
		"style-src-elem":  csp.StyleSrcElem,
		"style-src-attr":  csp.StyleSrcAttr,
		"worker-src":      csp.WorkerSrc,
	}

	for key, val := range directives {
		if val != "" {
			buffer.WriteString(fmt.Sprintf(policyFmt, key, val))
		}
	}

	return buffer.String()
}

func Strict() CSP {
	return CSP{
		ChildSrc:      "none",
		ConnectSrc:    "none",
		DefaultSrc:    "none",
		FontSrc:       "none",
		FrameSrc:      "none",
		ImgSrc:        "none",
		ManifestSrc:   "none",
		MediaSrc:      "none",
		ObjectSrc:     "none",
		ScriptSrc:     "none",
		ScriptSrcElem: "none",
		ScriptSrcAttr: "none",
		StyleSrc:      "none",
		StyleSrcElem:  "none",
		StyleSrcAttr:  "none",
		WorkerSrc:     "none",
	}
}
