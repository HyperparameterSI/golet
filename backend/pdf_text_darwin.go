// Author: HyperparameterSI

//go:build darwin

package golet

import "github.com/HyperparameterSI/golet/backend/internal/darwin"

// fixPrintedPDF gives the text of a PDF that WebKit printed the characters
// it names radicals for (pdf_text.go).
func fixPrintedPDF(pdf []byte) []byte { return fixPDFText(pdf, darwin.RadicalIdeographs) }
