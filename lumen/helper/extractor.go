package helper

import "github.com/hashicorp/terraform-plugin-sdk/v2/diag"

func ExtractDiagnosticWarnings(diagnostics diag.Diagnostics) diag.Diagnostics {
	warnings := diag.Diagnostics{}
	for _, d := range diagnostics {
		if d.Severity == diag.Warning {
			warnings = append(warnings, d)
		}
	}
	return warnings
}

func ExtractDiagnosticErrorIfPresent(diagnostics diag.Diagnostics) *diag.Diagnostic {
	if diagnostics.HasError() {
		for _, d := range diagnostics {
			if d.Severity == diag.Error {
				return &d
			}
		}
	}

	return nil
}
