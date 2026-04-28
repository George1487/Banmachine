package docx

import (
	"encoding/xml"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"
)

const lostMathStructureWarning = "lost_math_structure"

var formulaOperatorSpacingRegexp = regexp.MustCompile(`\s*([=+*/])\s*`)
var formulaNumberUnitSpacingRegexp = regexp.MustCompile(`(\d)([А-Яа-я])`)
var formulaSplitDecimalRegexp1 = regexp.MustCompile(`\b(\d)\s*\*\s*(\d+,\d+)\b`)
var formulaSplitDecimalRegexp2 = regexp.MustCompile(`\b(\d+,\d)\s*\*\s*(\d{2})\b`)
var formulaSplitDecimalRegexp3 = regexp.MustCompile(`\b(0,\d{3})\s*\*\s*(\d)\b`)
var formulaImpliedMulOpenParenRegexp = regexp.MustCompile(`([A-Za-zμΔαβχ]{2,})\(`)
var formulaImpliedMulCloseParenRegexp = regexp.MustCompile(`\)([A-Za-zμΔαβχ])`)
var formulaUnitReciprocalRegexp = regexp.MustCompile(`(\d+(?:,\d+)?)\s*\*\s*1\s*/\s*\(([^)]+)\)`)

type xmlNode struct {
	Name     xml.Name
	Attr     []xml.Attr
	Text     string
	Children []*xmlNode
}

func readXMLNode(decoder *xml.Decoder, start xml.StartElement) (*xmlNode, error) {
	node := &xmlNode{
		Name: start.Name,
		Attr: append([]xml.Attr(nil), start.Attr...),
	}

	for {
		token, err := decoder.Token()
		if err != nil {
			return nil, err
		}

		switch current := token.(type) {
		case xml.StartElement:
			child, err := readXMLNode(decoder, current)
			if err != nil {
				return nil, err
			}
			node.Children = append(node.Children, child)
		case xml.CharData:
			node.Text += string(current)
		case xml.EndElement:
			if current.Name == start.Name {
				return node, nil
			}
		}
	}
}

func buildFormulaCandidate(node *xmlNode) formulaCandidate {
	warningSet := make(map[string]struct{})
	raw := cleanFormulaText(flattenXMLNodeText(node))
	normalized := cleanFormulaText(normalizeMathNode(node, warningSet))
	if normalized == "" {
		normalized = raw
		warningSet[lostMathStructureWarning] = struct{}{}
	}

	warnings := warningSetToList(warningSet)
	confidence := 0.9
	if len(warnings) > 0 {
		confidence -= 0.18 * float64(len(warnings))
	}
	if raw == normalized {
		confidence -= 0.05
	}
	if confidence < 0.45 {
		confidence = 0.45
	}

	candidate := formulaCandidate{
		Type:       "formula",
		Raw:        raw,
		Normalized: normalized,
		Confidence: roundConfidence(confidence),
		Warnings:   warnings,
	}
	if !isUsefulFormulaCandidate(candidate) {
		return formulaCandidate{}
	}

	return candidate
}

func isUsefulFormulaCandidate(candidate formulaCandidate) bool {
	expr := cleanFormulaText(candidate.Normalized)
	if expr == "" {
		expr = cleanFormulaText(candidate.Raw)
	}
	if expr == "" {
		return false
	}
	if len([]rune(expr)) < 3 {
		return false
	}
	if formulaLooksIncomplete(expr) {
		return false
	}
	if simpleMathToken(expr) {
		return false
	}
	if !hasFormulaStructure(expr) {
		return false
	}
	return true
}

func formulaLooksIncomplete(value string) bool {
	value = strings.TrimSpace(value)
	if value == "" {
		return true
	}
	last, ok := lastMeaningfulRune(value)
	if !ok {
		return true
	}
	switch last {
	case '=', '+', '-', '*', '/', '^', '_', '(', '[', '{':
		return true
	default:
		return false
	}
}

func hasFormulaStructure(value string) bool {
	value = strings.TrimSpace(value)
	if value == "" {
		return false
	}
	if strings.ContainsAny(value, "=*/^") {
		return true
	}
	if strings.Contains(value, "+") || strings.Contains(value, " - ") {
		return true
	}
	if strings.Contains(value, "(") && strings.Contains(value, ")") {
		return true
	}
	if strings.Contains(strings.ToLower(value), "sqrt(") || strings.Contains(strings.ToLower(value), "root(") {
		return true
	}
	return false
}

func flattenXMLNodeText(node *xmlNode) string {
	if node == nil {
		return ""
	}

	var parts []string
	var walk func(*xmlNode)
	walk = func(current *xmlNode) {
		if current == nil {
			return
		}
		if text := strings.TrimSpace(current.Text); text != "" {
			parts = append(parts, text)
		}
		for _, child := range current.Children {
			walk(child)
		}
	}
	walk(node)
	return strings.Join(parts, " ")
}

func normalizeMathNode(node *xmlNode, warnings map[string]struct{}) string {
	if node == nil {
		return ""
	}

	switch node.Name.Local {
	case "oMath", "oMathPara", "num", "den", "e", "sub", "sup", "deg":
		return joinMathParts(normalizeMathChildren(node, warnings))
	case "r", "t":
		text := strings.TrimSpace(node.Text)
		if text != "" {
			return text
		}
		return joinMathParts(normalizeMathChildren(node, warnings))
	case "sSub":
		base := normalizeMathChild(node, "e", warnings)
		sub := normalizeMathChild(node, "sub", warnings)
		if base == "" || sub == "" {
			warnings[lostMathStructureWarning] = struct{}{}
			return joinMathParts(normalizeMathChildren(node, warnings))
		}
		return formatMathIndex(base, sub)
	case "sSup":
		base := normalizeMathChild(node, "e", warnings)
		sup := normalizeMathChild(node, "sup", warnings)
		if base == "" || sup == "" {
			warnings[lostMathStructureWarning] = struct{}{}
			return joinMathParts(normalizeMathChildren(node, warnings))
		}
		return formatMathSuperscript(base, sup)
	case "sSubSup":
		base := normalizeMathChild(node, "e", warnings)
		sub := normalizeMathChild(node, "sub", warnings)
		sup := normalizeMathChild(node, "sup", warnings)
		if base == "" {
			warnings[lostMathStructureWarning] = struct{}{}
			return joinMathParts(normalizeMathChildren(node, warnings))
		}
		result := base
		if sub != "" {
			result = formatMathIndex(result, sub)
		}
		if sup != "" {
			result = formatMathSuperscript(result, sup)
		}
		return result
	case "f":
		num := normalizeMathChild(node, "num", warnings)
		den := normalizeMathChild(node, "den", warnings)
		if num == "" || den == "" {
			warnings[lostMathStructureWarning] = struct{}{}
			return joinMathParts(normalizeMathChildren(node, warnings))
		}
		return formatMathFraction(num, den)
	case "d":
		content := normalizeMathChild(node, "e", warnings)
		if content == "" {
			warnings[lostMathStructureWarning] = struct{}{}
			return joinMathParts(normalizeMathChildren(node, warnings))
		}
		return formatMathDelimited(node, content)
	case "rad":
		content := normalizeMathChild(node, "e", warnings)
		if content == "" {
			warnings[lostMathStructureWarning] = struct{}{}
			return joinMathParts(normalizeMathChildren(node, warnings))
		}
		deg := normalizeMathChild(node, "deg", warnings)
		if deg == "" {
			return "sqrt(" + content + ")"
		}
		return "root(" + deg + ", " + content + ")"
	case "limLow":
		base := normalizeMathChild(node, "e", warnings)
		limit := normalizeMathChild(node, "lim", warnings)
		if base == "" {
			warnings[lostMathStructureWarning] = struct{}{}
			return joinMathParts(normalizeMathChildren(node, warnings))
		}
		if limit == "" {
			return base
		}
		return fmt.Sprintf("%s_{%s}", base, limit)
	case "limUpp":
		base := normalizeMathChild(node, "e", warnings)
		limit := normalizeMathChild(node, "lim", warnings)
		if base == "" {
			warnings[lostMathStructureWarning] = struct{}{}
			return joinMathParts(normalizeMathChildren(node, warnings))
		}
		if limit == "" {
			return base
		}
		return fmt.Sprintf("%s^{%s}", base, limit)
	case "m":
		return joinMathPartsWithSeparator(normalizeMathChildren(node, warnings), "; ")
	case "mr":
		return joinMathPartsWithSeparator(normalizeMathChildren(node, warnings), ", ")
	case "acc", "bar", "box", "borderBox", "func", "funcPr", "groupChr", "nary", "phant", "eqArr":
		return joinMathParts(normalizeMathChildren(node, warnings))
	default:
		if strings.HasSuffix(node.Name.Local, "Pr") || node.Name.Local == "ctrlPr" || node.Name.Local == "rPr" {
			return ""
		}
		if len(node.Children) > 0 {
			return joinMathParts(normalizeMathChildren(node, warnings))
		}
		text := strings.TrimSpace(node.Text)
		if text != "" {
			return text
		}
		warnings[lostMathStructureWarning] = struct{}{}
		return ""
	}
}

func normalizeMathChildren(node *xmlNode, warnings map[string]struct{}) []string {
	parts := make([]string, 0, len(node.Children))
	for _, child := range node.Children {
		part := normalizeMathNode(child, warnings)
		if part == "" {
			continue
		}
		parts = append(parts, part)
	}
	return parts
}

func normalizeMathChild(node *xmlNode, local string, warnings map[string]struct{}) string {
	for _, child := range node.Children {
		if child.Name.Local != local {
			continue
		}
		return normalizeMathNode(child, warnings)
	}
	return ""
}

func joinMathPartsWithSeparator(parts []string, separator string) string {
	filtered := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		filtered = append(filtered, part)
	}
	return strings.Join(filtered, separator)
}

func formatMathDelimited(node *xmlNode, content string) string {
	left := mathPropertyValue(node, "begChr")
	right := mathPropertyValue(node, "endChr")
	if left == "" && right == "" {
		return content
	}
	return left + content + right
}

func mathPropertyValue(node *xmlNode, local string) string {
	for _, child := range node.Children {
		if !strings.HasSuffix(child.Name.Local, "Pr") {
			continue
		}
		for _, grandChild := range child.Children {
			if grandChild.Name.Local != local {
				continue
			}
			if value := attrValue(grandChild.Attr, "val"); value != "" {
				return normalizeMathPropertyValue(value)
			}
			if value := strings.TrimSpace(grandChild.Text); value != "" {
				return normalizeMathPropertyValue(value)
			}
		}
	}
	return ""
}

func normalizeMathPropertyValue(value string) string {
	value = strings.TrimSpace(value)
	switch strings.ToLower(value) {
	case "_x0028_":
		return "("
	case "_x0029_":
		return ")"
	case "_x005b_":
		return "["
	case "_x005d_":
		return "]"
	case "_x007b_":
		return "{"
	case "_x007d_":
		return "}"
	default:
		return value
	}
}

func joinMathParts(parts []string) string {
	var buf strings.Builder
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		if buf.Len() == 0 {
			buf.WriteString(part)
			continue
		}

		current := buf.String()
		prevRune, prevOK := lastMeaningfulRune(current)
		nextRune, nextOK := firstMeaningfulRune(part)
		switch {
		case shouldConcatenateMathTokens(lastMathToken(current), part):
			// Keep broken numeric fragments together when Word split one number across runs.
		case prevOK && nextOK && shouldInsertMathMultiplication(prevRune, nextRune, part):
			buf.WriteString(" * ")
		case prevOK && nextOK && shouldInsertMathSpace(prevRune, nextRune):
			buf.WriteString(" ")
		}
		buf.WriteString(part)
	}
	return buf.String()
}

func shouldInsertMathMultiplication(prev, next rune, nextPart string) bool {
	if next == '(' || next == ')' || next == ',' || next == '.' {
		return false
	}
	if prev == '(' || prev == '^' || prev == '_' || prev == '/' || prev == '*' || prev == '+' || prev == '-' || prev == '=' {
		return false
	}
	if unicode.IsDigit(prev) && isLikelyUnitToken(nextPart) {
		return false
	}
	if strings.HasPrefix(nextPart, "(") || strings.HasPrefix(nextPart, "^") {
		return false
	}
	return isMathAtomRune(prev) && isMathAtomRune(next)
}

func shouldInsertMathSpace(prev, next rune) bool {
	if unicode.IsSpace(prev) || unicode.IsSpace(next) {
		return false
	}
	if prev == '(' || next == ')' {
		return false
	}
	return false
}

func isMathAtomRune(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r)
}

func shouldConcatenateMathTokens(prevToken, nextPart string) bool {
	prevToken = strings.TrimSpace(prevToken)
	nextPart = strings.TrimSpace(nextPart)
	if prevToken == "" || nextPart == "" {
		return false
	}
	if !isNumericFragment(prevToken) || !isNumericFragment(nextPart) {
		return false
	}
	if strings.Contains(nextPart, ",") || strings.Contains(nextPart, ".") {
		return true
	}
	return strings.HasPrefix(prevToken, "0,")
}

func isNumericFragment(value string) bool {
	if value == "" {
		return false
	}
	for _, r := range value {
		if unicode.IsDigit(r) || r == ',' || r == '.' || r == '-' {
			continue
		}
		return false
	}
	return true
}

func lastMathToken(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}

	start := len(value)
	for start > 0 {
		r, size := utf8.DecodeLastRuneInString(value[:start])
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == ',' || r == '.' || r == '-' {
			start -= size
			continue
		}
		break
	}
	return value[start:]
}

func isLikelyUnitToken(value string) bool {
	value = strings.TrimSpace(value)
	if value == "" {
		return false
	}

	hasCyrillic := false
	for _, r := range value {
		switch {
		case unicode.IsSpace(r), r == '/', r == '·', r == '^', r == ',', r == '.', r == '-', r == '°':
			continue
		case unicode.IsDigit(r):
			continue
		case unicode.In(r, unicode.Cyrillic):
			hasCyrillic = true
		case unicode.IsLetter(r):
			return false
		default:
			return false
		}
	}
	return hasCyrillic
}

func formatMathIndex(base, sub string) string {
	sub = strings.TrimSpace(sub)
	if sub == "" {
		return base
	}
	if simpleMathToken(sub) {
		return base + sub
	}
	return fmt.Sprintf("%s_{%s}", base, sub)
}

func formatMathSuperscript(base, sup string) string {
	sup = strings.TrimSpace(sup)
	if sup == "" {
		return base
	}
	if simpleMathToken(sup) {
		return base + "^" + sup
	}
	return fmt.Sprintf("%s^(%s)", base, sup)
}

func formatMathFraction(num, den string) string {
	num = strings.TrimSpace(num)
	den = strings.TrimSpace(den)
	if num == "" || den == "" {
		return strings.TrimSpace(num + " / " + den)
	}
	return fmt.Sprintf("%s / %s", wrapFormulaOperand(num), wrapFormulaOperand(den))
}

func wrapFormulaOperand(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return value
	}
	if strings.HasPrefix(value, "(") && strings.HasSuffix(value, ")") {
		return value
	}
	if simpleMathToken(value) {
		return value
	}
	return "(" + value + ")"
}

func simpleMathToken(value string) bool {
	if value == "" {
		return false
	}
	for _, r := range value {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			continue
		}
		if r == ',' || r == '.' || r == '-' {
			continue
		}
		return false
	}
	return true
}

func cleanFormulaText(value string) string {
	value = cleanText(value)
	if value == "" {
		return ""
	}
	value = normalizeSplitDecimalFragments(value)
	value = formulaOperatorSpacingRegexp.ReplaceAllString(value, " $1 ")
	value = normalizeScientificNotation(value)
	value = normalizeImpliedMultiplication(value)
	value = strings.ReplaceAll(value, "sqrt * (", "sqrt(")
	value = strings.ReplaceAll(value, "root * (", "root(")
	value = strings.ReplaceAll(value, " 1 / (", " 1/(")
	value = normalizeUnitAwareFormulaText(value)
	value = strings.ReplaceAll(value, "( ", "(")
	value = strings.ReplaceAll(value, " )", ")")
	value = strings.ReplaceAll(value, "[ ", "[")
	value = strings.ReplaceAll(value, " ]", "]")
	value = strings.ReplaceAll(value, "{ ", "{")
	value = strings.ReplaceAll(value, " }", "}")
	value = whitespaceRegexp.ReplaceAllString(value, " ")
	return strings.TrimSpace(value)
}

func normalizeSplitDecimalFragments(value string) string {
	value = formulaSplitDecimalRegexp1.ReplaceAllString(value, "$1$2")
	value = formulaSplitDecimalRegexp2.ReplaceAllString(value, "$1$2")
	value = formulaSplitDecimalRegexp3.ReplaceAllString(value, "${1}$2")
	return value
}

func normalizeScientificNotation(value string) string {
	value = strings.ReplaceAll(value, "^ -", "^-")
	value = strings.ReplaceAll(value, "^ +", "^+")
	value = strings.ReplaceAll(value, " × 10 ^ ", " * 10^")
	value = strings.ReplaceAll(value, " x 10 ^ ", " * 10^")
	value = strings.ReplaceAll(value, "·10 ^", "*10^")
	value = strings.ReplaceAll(value, " 10 ^ ", " 10^")
	return value
}

func normalizeImpliedMultiplication(value string) string {
	value = formulaImpliedMulOpenParenRegexp.ReplaceAllString(value, "$1 * (")
	value = formulaImpliedMulCloseParenRegexp.ReplaceAllString(value, ") * $1")
	value = formulaUnitReciprocalRegexp.ReplaceAllString(value, "$1 1 / ($2)")
	return value
}

func normalizeUnitAwareFormulaText(value string) string {
	if value == "" {
		return ""
	}
	value = formulaNumberUnitSpacingRegexp.ReplaceAllString(value, "$1 $2")
	return value
}

func warningSetToList(values map[string]struct{}) []string {
	if len(values) == 0 {
		return nil
	}
	result := make([]string, 0, len(values))
	for value := range values {
		result = append(result, value)
	}
	sort.Strings(result)
	return result
}
