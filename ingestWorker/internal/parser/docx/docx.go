package docx

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"ingestWorker/internal/parser"
	"io"
	"mime"
	"net/http"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	mimeDocxOfficial = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	mimeDocxAlt      = "application/docx"
)

var (
	whitespaceRegexp         = regexp.MustCompile(`\s+`)
	negativeNumberRegexp     = regexp.MustCompile(`(^|[\s(])-\s+(\d)`)
	inlineImageMarkerRegexp  = regexp.MustCompile(`\[IMAGE \d+\]`)
	imageOnlyParagraphRegexp = regexp.MustCompile(`^\[IMAGE \d+\]$`)
	numberedHeadingRegexp    = regexp.MustCompile(`^\d+(\.\d+)*[.)]?\s`)
	figureCaptionRegexp      = regexp.MustCompile(`(?i)^(рис(?:унок)?|figure|fig)\.?\s*\d+[.)]?\s*`)
	tableCaptionRegexp       = regexp.MustCompile(`(?i)^(табл(?:ица)?|table)\.?\s*\d+[.)-]?\s*`)
	calculationLabelRegexp   = regexp.MustCompile(`(?i)^(расчет|вычисление|определение|формула|коэрцитивная сила|остаточн\pL+ индукц\pL+|магнитн\pL+ проницаем\pL+)\s*[:.]?$`)
	localLabelRegexp         = regexp.MustCompile(`(?i)^[A-Za-zА-Яа-я][A-Za-zА-Яа-я0-9\s/-]{0,40}:\s*$`)
	pageNumberRegexp         = regexp.MustCompile(`^(стр\.?\s*)?\d+(\s*/\s*\d+)?$`)
	nameValueRegexp          = regexp.MustCompile(`(?i)^(факультет|группа|студент|выполнил|выполнила|преподаватель|руководитель|проверил|проверила|кафедра)\s*[:.-]\s*(.+)$`)
	axisLabelRegexp          = regexp.MustCompile(`^[A-Za-zА-Яа-яμΔαβχΩ/.,°\s-]{2,}$`)
	axisLabelHintRegexp      = regexp.MustCompile(`(?i)([hнbμ]\s*[,;]?\s*[aа]/[mм]|μ\s*\(\s*[hн]\s*\)|b\s*\(\s*[hн]\s*\)|магнитн\pL+\s+проницаем\pL+)`)
	logoLikeOCRRegexp        = regexp.MustCompile(`^[A-ZА-Я0-9 .,&-]{2,24}$`)
)

type DocxParser struct{}

type structuredDocument struct {
	Format            string             `json:"format"`
	Metadata          documentMetadata   `json:"metadata"`
	SemanticMetadata  *semanticMetadata  `json:"semanticMetadata,omitempty"`
	Sections          []section          `json:"sections"`
	Tables            []table            `json:"tables"`
	Images            []documentImage    `json:"images,omitempty"`
	AuxiliaryParts    []auxiliaryPart    `json:"auxiliaryParts,omitempty"`
	Stats             documentStats      `json:"stats"`
	RawText           string             `json:"rawText,omitempty"`
	NormalizedText    string             `json:"normalizedText,omitempty"`
	FormulaCandidates []formulaCandidate `json:"formulaCandidates,omitempty"`
	RawParts          rawParts           `json:"rawParts"`
}

type documentMetadata struct {
	Title          string `json:"title,omitempty"`
	Subject        string `json:"subject,omitempty"`
	Description    string `json:"description,omitempty"`
	Creator        string `json:"creator,omitempty"`
	LastModifiedBy string `json:"lastModifiedBy,omitempty"`
}

type semanticMetadata struct {
	Institution string `json:"institution,omitempty"`
	Faculty     string `json:"faculty,omitempty"`
	Group       string `json:"group,omitempty"`
	Student     string `json:"student,omitempty"`
	Instructor  string `json:"instructor,omitempty"`
}

type section struct {
	Title      string         `json:"title"`
	Paragraphs []string       `json:"paragraphs"`
	Blocks     []sectionBlock `json:"blocks,omitempty"`
	Confidence float64        `json:"confidence,omitempty"`
	Warnings   []string       `json:"warnings,omitempty"`
}

type table struct {
	Index         int        `json:"index"`
	Headers       []string   `json:"headers,omitempty"`
	Rows          [][]string `json:"rows"`
	Caption       string     `json:"caption,omitempty"`
	CaptionSource string     `json:"captionSource,omitempty"`
	Kind          string     `json:"kind,omitempty"`
	SourcePart    string     `json:"sourcePart,omitempty"`
	SourcePath    string     `json:"sourcePath,omitempty"`
	SourceKind    string     `json:"sourceKind,omitempty"`
	Confidence    float64    `json:"confidence,omitempty"`
	Warnings      []string   `json:"warnings,omitempty"`
}

type documentImage struct {
	Index           int      `json:"index"`
	RelationshipID  string   `json:"relationshipId,omitempty"`
	RelationshipIDs []string `json:"relationshipIds,omitempty"`
	Name            string   `json:"name,omitempty"`
	Path            string   `json:"path,omitempty"`
	ContentType     string   `json:"contentType,omitempty"`
	SizeBytes       int      `json:"sizeBytes,omitempty"`
	Width           int      `json:"width,omitempty"`
	Height          int      `json:"height,omitempty"`
	OCRText         string   `json:"ocrText,omitempty"`
	OCREngine       string   `json:"ocrEngine,omitempty"`
	Caption         string   `json:"caption,omitempty"`
	CaptionSource   string   `json:"captionSource,omitempty"`
	CaptionQuality  string   `json:"captionQuality,omitempty"`
	ImageType       string   `json:"imageType,omitempty"`
	ChartTitle      string   `json:"chartTitle,omitempty"`
	XAxisLabel      string   `json:"xAxisLabel,omitempty"`
	YAxisLabel      string   `json:"yAxisLabel,omitempty"`
	AxisHints       []string `json:"axisHints,omitempty"`
	SourcePart      string   `json:"sourcePart,omitempty"`
	SourcePath      string   `json:"sourcePath,omitempty"`
	SourceKind      string   `json:"sourceKind,omitempty"`
	UsageParts      []string `json:"usageParts,omitempty"`
	UsagePaths      []string `json:"usagePaths,omitempty"`
	Placement       string   `json:"placement,omitempty"`
	PlacementHints  []string `json:"placementHints,omitempty"`
	Referenced      bool     `json:"referenced,omitempty"`
	UsedInContent   bool     `json:"usedInContent,omitempty"`
	UsageCount      int      `json:"usageCount,omitempty"`
	Confidence      float64  `json:"confidence,omitempty"`
	Warnings        []string `json:"warnings,omitempty"`

	data            []byte              `json:"-"`
	relationshipSet map[string]struct{} `json:"-"`
	usagePartSet    map[string]struct{} `json:"-"`
	usagePathSet    map[string]struct{} `json:"-"`
	placementSet    map[string]struct{} `json:"-"`
}

type documentStats struct {
	ParagraphCount       int `json:"paragraphCount"`
	TableCount           int `json:"tableCount"`
	NonEmptyTables       int `json:"nonEmptyTables"`
	ImageCount           int `json:"imageCount"`
	ArchivedImageCount   int `json:"archivedImageCount,omitempty"`
	ReferencedImageCount int `json:"referencedImageCount,omitempty"`
	UsedImageCount       int `json:"usedImageCount,omitempty"`
	AuxiliaryPartCount   int `json:"auxiliaryPartCount,omitempty"`
}

type rawParts struct {
	Paragraphs          []string `json:"paragraphs"`
	AuxiliaryParagraphs []string `json:"auxiliaryParagraphs,omitempty"`
	ImageTexts          []string `json:"imageTexts,omitempty"`
}

type coreProperties struct {
	Title          string `xml:"title"`
	Subject        string `xml:"subject"`
	Description    string `xml:"description"`
	Creator        string `xml:"creator"`
	LastModifiedBy string `xml:"lastModifiedBy"`
}

type documentXML struct {
	Body bodyXML `xml:"body"`
}

type bodyXML struct {
	InnerXML string `xml:",innerxml"`
}

type paragraphXML struct {
	InnerXML string `xml:",innerxml"`
}

type tableXML struct {
	Rows []tableRowXML `xml:"tr"`
}

type tableRowXML struct {
	Cells []tableCellXML `xml:"tc"`
}

type tableCellXML struct {
	Paragraphs []paragraphXML `xml:"p"`
}

type innerXMLNode struct {
	InnerXML string `xml:",innerxml"`
}

type relationshipsXML struct {
	Relationships []relationshipXML `xml:"Relationship"`
}

type relationshipXML struct {
	ID         string `xml:"Id,attr"`
	Type       string `xml:"Type,attr"`
	Target     string `xml:"Target,attr"`
	TargetMode string `xml:"TargetMode,attr"`
}

type bodyBlock struct {
	Kind         string
	Text         string
	TableIndex   int
	ImageIndexes []int
	Formulas     []formulaCandidate
	Confidence   float64
	Warnings     []string
	SourcePart   string
	SourcePath   string
	SourceKind   string
}

type parsedParagraph struct {
	Text         string
	ImageIndexes []int
	Formulas     []formulaCandidate
	Confidence   float64
	Warnings     []string
}

type formulaCandidate struct {
	Type       string   `json:"type"`
	Raw        string   `json:"raw"`
	Normalized string   `json:"normalized"`
	Confidence float64  `json:"confidence,omitempty"`
	Warnings   []string `json:"warnings,omitempty"`
}

type sectionBlock struct {
	Type       string   `json:"type"`
	Text       string   `json:"text,omitempty"`
	Raw        string   `json:"raw,omitempty"`
	Normalized string   `json:"normalized,omitempty"`
	TableIndex int      `json:"tableIndex,omitempty"`
	ImageIndex int      `json:"imageIndex,omitempty"`
	Confidence float64  `json:"confidence,omitempty"`
	Warnings   []string `json:"warnings,omitempty"`
	SourcePart string   `json:"sourcePart,omitempty"`
	SourcePath string   `json:"sourcePath,omitempty"`
	SourceKind string   `json:"sourceKind,omitempty"`
}

type auxiliaryPart struct {
	SourcePart        string             `json:"sourcePart"`
	SourcePath        string             `json:"sourcePath,omitempty"`
	SourceKind        string             `json:"sourceKind,omitempty"`
	Category          string             `json:"category,omitempty"`
	Utility           string             `json:"utility,omitempty"`
	Paragraphs        []string           `json:"paragraphs,omitempty"`
	Tables            []table            `json:"tables,omitempty"`
	Blocks            []sectionBlock     `json:"blocks,omitempty"`
	FormulaCandidates []formulaCandidate `json:"formulaCandidates,omitempty"`
	Warnings          []string           `json:"warnings,omitempty"`
}

type documentPart struct {
	Path          string
	Name          string
	Kind          string
	Relationships map[string]relationshipXML
	Paragraphs    []string
	Tables        []table
	Blocks        []bodyBlock
}

type partScope struct {
	Path string
	Name string
	Kind string
}

type imageCatalog struct {
	byPath               map[string]*documentImage
	byScopedRelationship map[string]*documentImage
	ordered              []*documentImage
	archivedCount        int
}

func (d *DocxParser) SupportMimeType(mimeType string) bool {
	mimeType = strings.TrimSpace(strings.ToLower(mimeType))
	return mimeType == mimeDocxOfficial || mimeType == mimeDocxAlt
}

func (d *DocxParser) ParseDocument(ctx context.Context, filePath string) (*parser.ResultParser, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	reader, err := zip.OpenReader(filePath)
	if err != nil {
		return nil, fmt.Errorf("open docx archive: %w", err)
	}
	defer reader.Close()

	coreXMLBytes, _ := readZipFileOptional(reader.File, "docProps/core.xml")

	var props coreProperties
	if len(coreXMLBytes) > 0 {
		_ = xml.Unmarshal(coreXMLBytes, &props)
	}

	parts, err := discoverDocumentParts(reader.File)
	if err != nil {
		return nil, err
	}

	catalog, err := buildImageCatalog(ctx, reader.File, parts)
	if err != nil {
		return nil, err
	}

	for idx := range parts {
		if err := parseDocumentPart(reader.File, &parts[idx], catalog); err != nil {
			return nil, err
		}
	}

	mainPart, auxiliaryDocParts := splitParts(parts)
	images := catalog.allImages()
	images = annotateImages(parts, images)
	formulaCandidates := collectFormulaCandidates(mainPart.Blocks)
	sections := splitIntoSections(mainPart.Blocks, mainPart.Tables, images)
	rawText := buildRawText(mainPart.Blocks, mainPart.Tables, images)
	normalizedText := buildNormalizedText(mainPart.Blocks, mainPart.Tables, images)
	auxiliaryPayload := buildAuxiliaryParts(auxiliaryDocParts, images)
	semantic := extractSemanticMetadata(parts)

	title := cleanText(props.Title)
	if title == "" {
		title = filepath.Base(filePath)
	}

	payload := structuredDocument{
		Format: "docx",
		Metadata: documentMetadata{
			Title:          title,
			Subject:        cleanText(props.Subject),
			Description:    cleanText(props.Description),
			Creator:        cleanText(props.Creator),
			LastModifiedBy: cleanText(props.LastModifiedBy),
		},
		SemanticMetadata: semantic,
		Sections:         sections,
		Tables:           mainPart.Tables,
		Images:           images,
		AuxiliaryParts:   auxiliaryPayload,
		Stats: documentStats{
			ParagraphCount:       len(mainPart.Paragraphs),
			TableCount:           len(mainPart.Tables),
			NonEmptyTables:       countNonEmptyTables(mainPart.Tables),
			ImageCount:           len(images),
			ArchivedImageCount:   catalog.archivedCount,
			ReferencedImageCount: catalog.referencedImageCount(),
			UsedImageCount:       catalog.usedImageCount(),
			AuxiliaryPartCount:   len(auxiliaryPayload),
		},
		RawText:           rawText,
		NormalizedText:    normalizedText,
		FormulaCandidates: formulaCandidates,
		RawParts: rawParts{
			Paragraphs:          mainPart.Paragraphs,
			AuxiliaryParagraphs: collectAuxiliaryParagraphs(auxiliaryDocParts),
			ImageTexts:          collectImageTexts(images),
		},
	}

	structuredData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal structured data: %w", err)
	}

	return &parser.ResultParser{
		RawText:        rawText,
		StructuredData: structuredData,
	}, nil
}

func discoverDocumentParts(files []*zip.File) ([]documentPart, error) {
	if _, err := readZipFile(files, "word/document.xml"); err != nil {
		return nil, err
	}

	queue := []documentPart{{
		Path: "word/document.xml",
		Name: "document",
		Kind: "document",
	}}
	seen := make(map[string]struct{})
	parts := make([]documentPart, 0)

	for len(queue) > 0 {
		part := queue[0]
		queue = queue[1:]
		part.Path = path.Clean(part.Path)
		if _, ok := seen[part.Path]; ok {
			continue
		}
		seen[part.Path] = struct{}{}

		relsXMLBytes, _ := readZipFileOptional(files, relationshipsPath(part.Path))
		relationships, err := parseRelationships(relsXMLBytes)
		if err != nil {
			return nil, fmt.Errorf("parse relationships for %s: %w", part.Path, err)
		}
		part.Relationships = relationships
		parts = append(parts, part)

		nextParts := make([]documentPart, 0)
		for _, rel := range relationships {
			targetPath, ok := resolveRelationshipTarget(part.Path, rel)
			if !ok || !shouldTraversePart(rel, targetPath) {
				continue
			}
			name, kind := classifyDocumentPart(targetPath, rel.Type)
			nextParts = append(nextParts, documentPart{
				Path: targetPath,
				Name: name,
				Kind: kind,
			})
		}
		sort.Slice(nextParts, func(i, j int) bool {
			return nextParts[i].Path < nextParts[j].Path
		})
		queue = append(queue, nextParts...)
	}

	return parts, nil
}

func splitParts(parts []documentPart) (documentPart, []documentPart) {
	mainPart := documentPart{}
	auxiliary := make([]documentPart, 0, len(parts))
	for _, part := range parts {
		if path.Clean(part.Path) == "word/document.xml" {
			mainPart = part
			continue
		}
		auxiliary = append(auxiliary, part)
	}
	return mainPart, auxiliary
}

func parseDocumentPart(files []*zip.File, part *documentPart, catalog *imageCatalog) error {
	if part == nil {
		return nil
	}

	data, err := readZipFile(files, part.Path)
	if err != nil {
		return err
	}

	scope := partScope{
		Path: part.Path,
		Name: part.Name,
		Kind: part.Kind,
	}

	paragraphs, tables, blocks, err := parsePartXML(data, scope, catalog)
	if err != nil {
		return err
	}

	part.Paragraphs = paragraphs
	part.Blocks = blocks
	part.Tables = annotateTables(blocks, tables)
	return nil
}

func parsePartXML(data []byte, scope partScope, catalog *imageCatalog) ([]string, []table, []bodyBlock, error) {
	var root innerXMLNode
	if err := xml.Unmarshal(data, &root); err != nil {
		return nil, nil, nil, fmt.Errorf("decode %s: %w", scope.Path, err)
	}

	tableIndex := 0
	if path.Clean(scope.Path) == "word/document.xml" {
		var doc documentXML
		if err := xml.Unmarshal(data, &doc); err != nil {
			return nil, nil, nil, fmt.Errorf("decode %s: %w", scope.Path, err)
		}
		return parseBlockContainer(doc.Body.InnerXML, scope, catalog, &tableIndex)
	}

	return parseBlockContainer(root.InnerXML, scope, catalog, &tableIndex)
}

func parseBlockContainer(innerXML string, scope partScope, catalog *imageCatalog, tableIndex *int) ([]string, []table, []bodyBlock, error) {
	decoder := xml.NewDecoder(strings.NewReader(wrapInnerXML(innerXML)))

	paragraphs := make([]string, 0)
	tables := make([]table, 0)
	blocks := make([]bodyBlock, 0)

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, nil, nil, fmt.Errorf("decode %s: %w", scope.Path, err)
		}

		start, ok := token.(xml.StartElement)
		if !ok {
			continue
		}

		switch start.Name.Local {
		case "p":
			var paragraphNode paragraphXML
			if err := decoder.DecodeElement(&paragraphNode, &start); err != nil {
				return nil, nil, nil, fmt.Errorf("decode paragraph: %w", err)
			}

			parsed, err := parseParagraph(paragraphNode.InnerXML, scope, catalog)
			if err != nil {
				return nil, nil, nil, err
			}
			if parsed.Text == "" {
				continue
			}

			paragraphs = append(paragraphs, parsed.Text)
			blocks = append(blocks, bodyBlock{
				Kind:         "paragraph",
				Text:         parsed.Text,
				ImageIndexes: parsed.ImageIndexes,
				Formulas:     parsed.Formulas,
				Confidence:   parsed.Confidence,
				Warnings:     parsed.Warnings,
				SourcePart:   scope.Name,
				SourcePath:   scope.Path,
				SourceKind:   scope.Kind,
			})
		case "tbl":
			var tableNode tableXML
			if err := decoder.DecodeElement(&tableNode, &start); err != nil {
				return nil, nil, nil, fmt.Errorf("decode table: %w", err)
			}

			nextTableIndex := len(tables) + 1
			if tableIndex != nil {
				*tableIndex = *tableIndex + 1
				nextTableIndex = *tableIndex
			}
			parsedTable, hasContent, err := parseTable(tableNode, nextTableIndex, scope, catalog)
			if err != nil {
				return nil, nil, nil, err
			}
			if !hasContent {
				continue
			}

			tables = append(tables, parsedTable)
			blocks = append(blocks, bodyBlock{
				Kind:       "table",
				TableIndex: parsedTable.Index,
				SourcePart: scope.Name,
				SourcePath: scope.Path,
				SourceKind: scope.Kind,
			})
		default:
			if !isRecursiveContentContainer(start.Name.Local) {
				continue
			}

			var node innerXMLNode
			if err := decoder.DecodeElement(&node, &start); err != nil {
				return nil, nil, nil, fmt.Errorf("decode %s container %s: %w", scope.Path, start.Name.Local, err)
			}

			nestedParagraphs, nestedTables, nestedBlocks, err := parseBlockContainer(node.InnerXML, scope, catalog, tableIndex)
			if err != nil {
				return nil, nil, nil, err
			}

			paragraphs = append(paragraphs, nestedParagraphs...)
			tables = append(tables, nestedTables...)
			blocks = append(blocks, nestedBlocks...)
		}
	}

	return paragraphs, tables, blocks, nil
}

func parseParagraph(innerXML string, scope partScope, catalog *imageCatalog) (parsedParagraph, error) {
	decoder := xml.NewDecoder(strings.NewReader(wrapInnerXML(innerXML)))

	segments := make([]paragraphSegment, 0)
	var textBuffer strings.Builder
	imageIndexes := make([]int, 0)
	seenImages := make(map[int]struct{})
	formulas := make([]formulaCandidate, 0)
	warningSet := make(map[string]struct{})
	textElementDepth := 0
	elementStack := make([]string, 0, 16)

	flushText := func() {
		text := normalizeParagraphText(textBuffer.String(), warningSet)
		textBuffer.Reset()
		if text == "" {
			return
		}
		segments = append(segments, paragraphSegment{Kind: "text", Text: text})
	}

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return parsedParagraph{}, fmt.Errorf("decode paragraph text: %w", err)
		}

		switch current := token.(type) {
		case xml.StartElement:
			elementStack = append(elementStack, current.Name.Local)
			switch current.Name.Local {
			case "t", "instrText", "delText":
				textElementDepth++
			case "oMath", "oMathPara":
				flushText()
				node, err := readXMLNode(decoder, current)
				if err != nil {
					return parsedParagraph{}, fmt.Errorf("decode paragraph math: %w", err)
				}
				elementStack = elementStack[:len(elementStack)-1]
				formula := buildFormulaCandidate(node)
				if formula.Normalized == "" && formula.Raw == "" {
					continue
				}
				formulas = append(formulas, formula)
				text := formula.Normalized
				if text == "" {
					text = formula.Raw
				}
				segments = append(segments, paragraphSegment{Kind: "formula", Text: text})
				for _, warning := range formula.Warnings {
					warningSet[warning] = struct{}{}
				}
			case "tab":
				textBuffer.WriteString(" | ")
			case "br", "cr":
				textBuffer.WriteString(" ")
			case "blip":
				flushText()
				relID := attrValue(current.Attr, "embed")
				if relID == "" {
					relID = attrValue(current.Attr, "link")
				}
				if relID == "" {
					continue
				}
				placement, hints := detectImagePlacement(elementStack)
				image := catalog.registerUsage(scope, relID, placement, hints)
				if image == nil {
					continue
				}
				segments = append(segments, paragraphSegment{Kind: "image", Text: formatInlineImageMarker(image.Index)})
				if _, ok := seenImages[image.Index]; !ok {
					seenImages[image.Index] = struct{}{}
					imageIndexes = append(imageIndexes, image.Index)
				}
			case "imagedata":
				flushText()
				relID := attrValue(current.Attr, "id")
				if relID == "" {
					continue
				}
				placement, hints := detectImagePlacement(elementStack)
				image := catalog.registerUsage(scope, relID, placement, hints)
				if image == nil {
					continue
				}
				segments = append(segments, paragraphSegment{Kind: "image", Text: formatInlineImageMarker(image.Index)})
				if _, ok := seenImages[image.Index]; !ok {
					seenImages[image.Index] = struct{}{}
					imageIndexes = append(imageIndexes, image.Index)
				}
			}
		case xml.EndElement:
			switch current.Name.Local {
			case "t", "instrText", "delText":
				if textElementDepth > 0 {
					textElementDepth--
				}
			}
			if len(elementStack) > 0 {
				elementStack = elementStack[:len(elementStack)-1]
			}
		case xml.CharData:
			if textElementDepth > 0 {
				textBuffer.WriteString(string(current))
			}
		}
	}

	flushText()
	text := joinParagraphSegments(segments)
	warnings := warningSetToList(warningSet)
	confidence := 0.96
	if len(warnings) > 0 {
		confidence -= 0.08 * float64(len(warnings))
	}
	if confidence < 0.55 {
		confidence = 0.55
	}

	return parsedParagraph{
		Text:         text,
		ImageIndexes: imageIndexes,
		Formulas:     formulas,
		Confidence:   roundConfidence(confidence),
		Warnings:     warnings,
	}, nil
}

func parseTable(node tableXML, index int, scope partScope, catalog *imageCatalog) (table, bool, error) {
	rows := make([][]string, 0, len(node.Rows))
	hasContent := false

	for _, row := range node.Rows {
		cells := make([]string, 0, len(row.Cells))
		for _, cell := range row.Cells {
			paragraphTexts := make([]string, 0, len(cell.Paragraphs))
			for _, paragraphNode := range cell.Paragraphs {
				parsed, err := parseParagraph(paragraphNode.InnerXML, scope, catalog)
				if err != nil {
					return table{}, false, err
				}
				if parsed.Text == "" {
					continue
				}
				paragraphTexts = append(paragraphTexts, parsed.Text)
			}

			cellText := mergeParagraphFragments(paragraphTexts)
			if cellText != "" {
				hasContent = true
			}
			cells = append(cells, cellText)
		}
		if len(cells) > 0 {
			rows = append(rows, cells)
		}
	}

	return table{
		Index:      index,
		Rows:       rows,
		SourcePart: scope.Name,
		SourcePath: scope.Path,
		SourceKind: scope.Kind,
	}, hasContent, nil
}

func parseRelationships(data []byte) (map[string]relationshipXML, error) {
	if len(data) == 0 {
		return map[string]relationshipXML{}, nil
	}

	var rels relationshipsXML
	if err := xml.Unmarshal(data, &rels); err != nil {
		return nil, fmt.Errorf("decode document relationships: %w", err)
	}

	result := make(map[string]relationshipXML, len(rels.Relationships))
	for _, rel := range rels.Relationships {
		result[rel.ID] = rel
	}
	return result, nil
}

func buildImageCatalog(ctx context.Context, files []*zip.File, parts []documentPart) (*imageCatalog, error) {
	catalog := &imageCatalog{
		byPath:               make(map[string]*documentImage),
		byScopedRelationship: make(map[string]*documentImage),
		ordered:              make([]*documentImage, 0),
		archivedCount:        countArchivedImages(files),
	}

	for _, part := range parts {
		scope := partScope{
			Path: part.Path,
			Name: part.Name,
			Kind: part.Kind,
		}
		for _, rel := range part.Relationships {
			if !strings.Contains(strings.ToLower(rel.Type), "/image") {
				continue
			}
			if err := catalog.addReferencedImage(files, scope, rel); err != nil {
				return nil, err
			}
		}
	}

	enrichImagesWithOCR(ctx, catalog.imagePointers())
	for _, imageInfo := range catalog.byPath {
		imageInfo.data = nil
	}

	return catalog, nil
}

func (c *imageCatalog) addReferencedImage(files []*zip.File, scope partScope, rel relationshipXML) error {
	targetPath, ok := resolveRelationshipTarget(scope.Path, rel)
	if !ok {
		return nil
	}

	data, err := readZipFile(files, targetPath)
	if err != nil {
		return fmt.Errorf("read image %s: %w", targetPath, err)
	}

	imageInfo := c.byPath[targetPath]
	if imageInfo == nil {
		imageInfo = &documentImage{
			Name:            path.Base(targetPath),
			Path:            targetPath,
			ContentType:     detectImageContentType(targetPath, data),
			SizeBytes:       len(data),
			SourcePart:      scope.Name,
			SourcePath:      scope.Path,
			SourceKind:      scope.Kind,
			Referenced:      true,
			data:            data,
			relationshipSet: make(map[string]struct{}),
			usagePartSet:    make(map[string]struct{}),
			usagePathSet:    make(map[string]struct{}),
			placementSet:    make(map[string]struct{}),
		}
		if cfg, _, err := image.DecodeConfig(bytes.NewReader(data)); err == nil {
			imageInfo.Width = cfg.Width
			imageInfo.Height = cfg.Height
		}
		c.byPath[targetPath] = imageInfo
	}
	if imageInfo.RelationshipID == "" {
		imageInfo.RelationshipID = rel.ID
	}
	if imageInfo.RelationshipIDs == nil {
		imageInfo.RelationshipIDs = make([]string, 0, 1)
	}
	if _, ok := imageInfo.relationshipSet[rel.ID]; !ok {
		imageInfo.relationshipSet[rel.ID] = struct{}{}
		imageInfo.RelationshipIDs = append(imageInfo.RelationshipIDs, rel.ID)
		sort.Strings(imageInfo.RelationshipIDs)
	}
	registerImageScope(imageInfo, scope)
	c.byScopedRelationship[scopedRelationshipKey(scope.Path, rel.ID)] = imageInfo
	return nil
}

func (c *imageCatalog) registerUsage(scope partScope, relID, placement string, placementHints []string) *documentImage {
	imageInfo := c.byScopedRelationship[scopedRelationshipKey(scope.Path, relID)]
	if imageInfo == nil {
		return nil
	}

	registerImageScope(imageInfo, scope)
	imageInfo.UsedInContent = true
	imageInfo.UsageCount++
	if imageInfo.Index == 0 {
		imageInfo.Index = len(c.ordered) + 1
		c.ordered = append(c.ordered, imageInfo)
	}
	if placement != "" && imageInfo.Placement == "" {
		imageInfo.Placement = placement
	}
	for _, hint := range placementHints {
		if hint == "" {
			continue
		}
		if _, ok := imageInfo.placementSet[hint]; ok {
			continue
		}
		imageInfo.placementSet[hint] = struct{}{}
		imageInfo.PlacementHints = append(imageInfo.PlacementHints, hint)
		sort.Strings(imageInfo.PlacementHints)
	}
	return imageInfo
}

func registerImageScope(imageInfo *documentImage, scope partScope) {
	if imageInfo == nil {
		return
	}
	if imageInfo.SourcePart == "" {
		imageInfo.SourcePart = scope.Name
	}
	if imageInfo.SourcePath == "" {
		imageInfo.SourcePath = scope.Path
	}
	if imageInfo.SourceKind == "" {
		imageInfo.SourceKind = scope.Kind
	}
	if imageInfo.usagePartSet == nil {
		imageInfo.usagePartSet = make(map[string]struct{})
	}
	if imageInfo.usagePathSet == nil {
		imageInfo.usagePathSet = make(map[string]struct{})
	}
	if scope.Name != "" {
		if _, ok := imageInfo.usagePartSet[scope.Name]; !ok {
			imageInfo.usagePartSet[scope.Name] = struct{}{}
			imageInfo.UsageParts = append(imageInfo.UsageParts, scope.Name)
			sort.Strings(imageInfo.UsageParts)
		}
	}
	if scope.Path != "" {
		if _, ok := imageInfo.usagePathSet[scope.Path]; !ok {
			imageInfo.usagePathSet[scope.Path] = struct{}{}
			imageInfo.UsagePaths = append(imageInfo.UsagePaths, scope.Path)
			sort.Strings(imageInfo.UsagePaths)
		}
	}
}

func (c *imageCatalog) ensureReferencedIndexes() {
	unassigned := make([]*documentImage, 0)
	for _, imageInfo := range c.byPath {
		if imageInfo == nil {
			continue
		}
		if imageInfo.Index == 0 {
			unassigned = append(unassigned, imageInfo)
		}
	}
	sort.Slice(unassigned, func(i, j int) bool {
		if unassigned[i].SourcePath == unassigned[j].SourcePath {
			return unassigned[i].Path < unassigned[j].Path
		}
		return unassigned[i].SourcePath < unassigned[j].SourcePath
	})
	for _, imageInfo := range unassigned {
		imageInfo.Index = len(c.ordered) + 1
		c.ordered = append(c.ordered, imageInfo)
	}
}

func (c *imageCatalog) allImages() []documentImage {
	c.ensureReferencedIndexes()
	result := make([]documentImage, 0, len(c.ordered))
	for _, imageInfo := range c.ordered {
		if imageInfo == nil {
			continue
		}
		result = append(result, *imageInfo)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Index < result[j].Index
	})
	return result
}

func (c *imageCatalog) imagePointers() []*documentImage {
	images := make([]*documentImage, 0, len(c.byPath))
	for _, imageInfo := range c.byPath {
		if imageInfo == nil {
			continue
		}
		images = append(images, imageInfo)
	}
	sort.Slice(images, func(i, j int) bool {
		return images[i].Path < images[j].Path
	})
	return images
}

func (c *imageCatalog) referencedImageCount() int {
	return len(c.byPath)
}

func (c *imageCatalog) usedImageCount() int {
	count := 0
	for _, imageInfo := range c.byPath {
		if imageInfo == nil || !imageInfo.UsedInContent {
			continue
		}
		count++
	}
	return count
}

func buildRawText(blocks []bodyBlock, tables []table, images []documentImage) string {
	return renderDocumentText(blocks, tables, images, false)
}

func buildNormalizedText(blocks []bodyBlock, tables []table, images []documentImage) string {
	return renderDocumentText(blocks, tables, images, true)
}

func renderDocumentText(blocks []bodyBlock, tables []table, images []documentImage, includeCaptions bool) string {
	imageByIndex := make(map[int]documentImage, len(images))
	for _, imageInfo := range images {
		imageByIndex[imageInfo.Index] = imageInfo
	}

	tableByIndex := make(map[int]table, len(tables))
	for _, tbl := range tables {
		tableByIndex[tbl.Index] = tbl
	}

	parts := make([]string, 0, len(blocks))
	emittedOCR := make(map[int]struct{})

	for _, block := range blocks {
		switch block.Kind {
		case "paragraph":
			if block.Text != "" {
				parts = append(parts, block.Text)
			}
			for _, imageIndex := range block.ImageIndexes {
				if _, ok := emittedOCR[imageIndex]; ok {
					continue
				}
				emittedOCR[imageIndex] = struct{}{}
				imageInfo, ok := imageByIndex[imageIndex]
				if !ok {
					continue
				}
				if includeCaptions && imageInfo.Caption != "" {
					parts = append(parts, fmt.Sprintf("[IMAGE %d CAPTION]\n%s", imageInfo.Index, imageInfo.Caption))
				}
				if imageInfo.OCRText == "" {
					continue
				}
				parts = append(parts, fmt.Sprintf("[IMAGE %d OCR]\n%s", imageInfo.Index, imageInfo.OCRText))
			}
		case "table":
			tbl, ok := tableByIndex[block.TableIndex]
			if !ok {
				continue
			}
			parts = append(parts, formatTable(tbl))
		}
	}

	return strings.TrimSpace(strings.Join(parts, "\n\n"))
}

func formatTable(tbl table) string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("[TABLE %d]\n", tbl.Index))
	for _, row := range tbl.Rows {
		buf.WriteString(strings.Join(row, " | "))
		buf.WriteString("\n")
	}
	return strings.TrimSpace(buf.String())
}

func collectImageTexts(images []documentImage) []string {
	texts := make([]string, 0, len(images))
	for _, imageInfo := range images {
		if imageInfo.OCRText == "" {
			continue
		}
		texts = append(texts, imageInfo.OCRText)
	}
	return texts
}

func collectFormulaCandidates(blocks []bodyBlock) []formulaCandidate {
	result := make([]formulaCandidate, 0)
	for _, block := range blocks {
		result = append(result, block.Formulas...)
	}
	return result
}

func buildAuxiliaryParts(parts []documentPart, images []documentImage) []auxiliaryPart {
	result := make([]auxiliaryPart, 0, len(parts))
	for _, part := range parts {
		if len(part.Paragraphs) == 0 && len(part.Tables) == 0 && len(part.Blocks) == 0 {
			continue
		}

		blocks := buildSectionBlocks(part.Blocks, part.Tables, images)
		warningSet := make(map[string]struct{})
		for _, block := range blocks {
			for _, warning := range block.Warnings {
				warningSet[warning] = struct{}{}
			}
		}
		category, utility := classifyAuxiliaryPart(part)
		if utility == "low" {
			warningSet["low_value_auxiliary"] = struct{}{}
		}

		result = append(result, auxiliaryPart{
			SourcePart:        part.Name,
			SourcePath:        part.Path,
			SourceKind:        part.Kind,
			Category:          category,
			Utility:           utility,
			Paragraphs:        append([]string(nil), part.Paragraphs...),
			Tables:            append([]table(nil), part.Tables...),
			Blocks:            blocks,
			FormulaCandidates: collectFormulaCandidates(part.Blocks),
			Warnings:          warningSetToList(warningSet),
		})
	}
	return result
}

func classifyAuxiliaryPart(part documentPart) (string, string) {
	if part.Kind == "footer" && allParagraphsMatch(part.Paragraphs, pageNumberRegexp) {
		return "pagination", "low"
	}
	if part.Kind == "header" {
		if hasBrandingLikeText(part.Paragraphs) || hasTableKind(part.Tables, "layout_table") {
			return "branding", "low"
		}
		return "related_content", "medium"
	}
	switch part.Kind {
	case "footnotes", "endnotes", "comments":
		return "notes", "medium"
	default:
		return "related_content", "medium"
	}
}

func hasBrandingLikeText(paragraphs []string) bool {
	for _, paragraph := range paragraphs {
		lower := strings.ToLower(cleanText(paragraph))
		if strings.Contains(lower, "университет") || strings.Contains(lower, "институт") || strings.Contains(lower, "кафедра") {
			return true
		}
	}
	return false
}

func allParagraphsMatch(paragraphs []string, pattern *regexp.Regexp) bool {
	if len(paragraphs) == 0 {
		return false
	}
	for _, paragraph := range paragraphs {
		paragraph = cleanText(paragraph)
		if paragraph == "" {
			continue
		}
		if !pattern.MatchString(paragraph) {
			return false
		}
	}
	return true
}

func hasTableKind(tables []table, kind string) bool {
	for _, tbl := range tables {
		if tbl.Kind == kind {
			return true
		}
	}
	return false
}

func collectAuxiliaryParagraphs(parts []documentPart) []string {
	paragraphs := make([]string, 0)
	for _, part := range parts {
		paragraphs = append(paragraphs, part.Paragraphs...)
	}
	return paragraphs
}

func extractSemanticMetadata(parts []documentPart) *semanticMetadata {
	meta := &semanticMetadata{}
	assign := func(field *string, value string) {
		value = cleanText(value)
		if value == "" || *field != "" {
			return
		}
		*field = value
	}

	for _, part := range parts {
		for _, paragraph := range part.Paragraphs {
			paragraph = cleanText(paragraph)
			if paragraph == "" {
				continue
			}
			if matches := nameValueRegexp.FindStringSubmatch(paragraph); len(matches) == 3 {
				key := strings.ToLower(matches[1])
				value := matches[2]
				switch key {
				case "факультет", "кафедра":
					assign(&meta.Faculty, value)
				case "группа":
					assign(&meta.Group, value)
				case "студент", "выполнил", "выполнила":
					assign(&meta.Student, value)
				case "преподаватель", "руководитель", "проверил", "проверила":
					assign(&meta.Instructor, value)
				}
				continue
			}
			if meta.Institution == "" && (part.Kind == "header" || part.Kind == "document") {
				lower := strings.ToLower(paragraph)
				if strings.Contains(lower, "университет") || strings.Contains(lower, "институт") {
					assign(&meta.Institution, paragraph)
				}
			}
		}
	}

	if *meta == (semanticMetadata{}) {
		return nil
	}
	return meta
}

func buildSectionBlocks(blocks []bodyBlock, tables []table, images []documentImage) []sectionBlock {
	tableByIndex := make(map[int]table, len(tables))
	for _, tbl := range tables {
		tableByIndex[tbl.Index] = tbl
	}

	imageByIndex := make(map[int]documentImage, len(images))
	for _, imageInfo := range images {
		imageByIndex[imageInfo.Index] = imageInfo
	}

	result := make([]sectionBlock, 0, len(blocks))
	for _, block := range blocks {
		switch block.Kind {
		case "paragraph":
			if block.Text != "" {
				blockType := "paragraph"
				switch classifyParagraphRole(block.Text) {
				case "figure_label":
					blockType = "figure_label"
				case "table_label":
					blockType = "table_label"
				case "local_label":
					blockType = "local_label"
				case "calculation_label":
					blockType = "calculation_label"
				}
				result = append(result, sectionBlock{
					Type:       blockType,
					Text:       block.Text,
					Confidence: block.Confidence,
					Warnings:   append([]string(nil), block.Warnings...),
					SourcePart: block.SourcePart,
					SourcePath: block.SourcePath,
					SourceKind: block.SourceKind,
				})
			}
			for _, formula := range block.Formulas {
				result = append(result, sectionBlock{
					Type:       "formula",
					Raw:        formula.Raw,
					Normalized: formula.Normalized,
					Confidence: formula.Confidence,
					Warnings:   append([]string(nil), formula.Warnings...),
					SourcePart: block.SourcePart,
					SourcePath: block.SourcePath,
					SourceKind: block.SourceKind,
				})
			}
			for _, imageIndex := range block.ImageIndexes {
				imageInfo := imageByIndex[imageIndex]
				result = append(result, sectionBlock{
					Type:       "image_ref",
					ImageIndex: imageIndex,
					Text:       imageInfo.Caption,
					Confidence: imageInfo.Confidence,
					Warnings:   append([]string(nil), imageInfo.Warnings...),
					SourcePart: block.SourcePart,
					SourcePath: block.SourcePath,
					SourceKind: block.SourceKind,
				})
			}
		case "table":
			tbl := tableByIndex[block.TableIndex]
			result = append(result, sectionBlock{
				Type:       "table_ref",
				TableIndex: block.TableIndex,
				Text:       strings.Join(tbl.Headers, " | "),
				Confidence: tbl.Confidence,
				Warnings:   append([]string(nil), tbl.Warnings...),
				SourcePart: block.SourcePart,
				SourcePath: block.SourcePath,
				SourceKind: block.SourceKind,
			})
		}
	}

	return result
}

func splitIntoSections(blocks []bodyBlock, tables []table, images []documentImage) []section {
	if len(blocks) == 0 {
		return nil
	}

	tableByIndex := make(map[int]table, len(tables))
	for _, tbl := range tables {
		tableByIndex[tbl.Index] = tbl
	}

	imageByIndex := make(map[int]documentImage, len(images))
	for _, imageInfo := range images {
		imageByIndex[imageInfo.Index] = imageInfo
	}

	sections := []section{}
	current := section{
		Title:      "document",
		Confidence: 0.98,
	}

	for _, block := range blocks {
		switch block.Kind {
		case "paragraph":
			role := classifyParagraphRole(block.Text)
			if role == "section_heading" {
				if len(current.Paragraphs) > 0 || len(current.Blocks) > 0 || current.Title != "document" {
					current = finalizeSection(current)
					sections = append(sections, current)
				}
				current = section{
					Title:      block.Text,
					Confidence: roundConfidence(maxFloat(block.Confidence, 0.84)),
				}
				current.Blocks = append(current.Blocks, sectionBlock{
					Type:       "section_title",
					Text:       block.Text,
					Confidence: block.Confidence,
					Warnings:   append([]string(nil), block.Warnings...),
					SourcePart: block.SourcePart,
					SourcePath: block.SourcePath,
					SourceKind: block.SourceKind,
				})
				continue
			}

			if block.Text != "" {
				current.Paragraphs = append(current.Paragraphs, block.Text)
				blockType := "paragraph"
				switch role {
				case "figure_label":
					blockType = "figure_label"
				case "table_label":
					blockType = "table_label"
				case "local_label":
					blockType = "local_label"
				case "calculation_label":
					blockType = "calculation_label"
				}
				current.Blocks = append(current.Blocks, sectionBlock{
					Type:       blockType,
					Text:       block.Text,
					Confidence: block.Confidence,
					Warnings:   append([]string(nil), block.Warnings...),
					SourcePart: block.SourcePart,
					SourcePath: block.SourcePath,
					SourceKind: block.SourceKind,
				})
			}
			for _, formula := range block.Formulas {
				current.Blocks = append(current.Blocks, sectionBlock{
					Type:       "formula",
					Raw:        formula.Raw,
					Normalized: formula.Normalized,
					Confidence: formula.Confidence,
					Warnings:   append([]string(nil), formula.Warnings...),
					SourcePart: block.SourcePart,
					SourcePath: block.SourcePath,
					SourceKind: block.SourceKind,
				})
			}
			for _, imageIndex := range block.ImageIndexes {
				imageInfo := imageByIndex[imageIndex]
				current.Blocks = append(current.Blocks, sectionBlock{
					Type:       "image_ref",
					ImageIndex: imageIndex,
					Text:       imageInfo.Caption,
					Confidence: imageInfo.Confidence,
					Warnings:   append([]string(nil), imageInfo.Warnings...),
					SourcePart: block.SourcePart,
					SourcePath: block.SourcePath,
					SourceKind: block.SourceKind,
				})
			}
		case "table":
			tbl := tableByIndex[block.TableIndex]
			current.Blocks = append(current.Blocks, sectionBlock{
				Type:       "table_ref",
				TableIndex: block.TableIndex,
				Text:       strings.Join(tbl.Headers, " | "),
				Confidence: tbl.Confidence,
				Warnings:   append([]string(nil), tbl.Warnings...),
				SourcePart: block.SourcePart,
				SourcePath: block.SourcePath,
				SourceKind: block.SourceKind,
			})
		}
	}

	if len(current.Paragraphs) > 0 || len(current.Blocks) > 0 || current.Title != "document" {
		sections = append(sections, finalizeSection(current))
	}

	return sections
}

func looksLikeHeading(value string) bool {
	return classifyParagraphRole(value) == "section_heading"
}

func classifyParagraphRole(value string) string {
	value = cleanText(value)
	if value == "" || imageOnlyParagraphRegexp.MatchString(value) {
		return "paragraph"
	}
	if len([]rune(value)) > 90 {
		return "paragraph"
	}
	if figureCaptionRegexp.MatchString(value) {
		return "figure_label"
	}
	if tableCaptionRegexp.MatchString(value) {
		return "table_label"
	}
	if calculationLabelRegexp.MatchString(value) {
		return "calculation_label"
	}
	if localLabelRegexp.MatchString(value) {
		return "local_label"
	}

	keywords := []string{
		"цель работы",
		"задачи, решаемые при выполнении работы",
		"объект исследования",
		"метод экспериментального исследования",
		"рабочие формулы",
		"исходные данные",
		"измерительные приборы",
		"схема установки",
		"результаты прямых измерений",
		"расчет результатов косвенных измерений",
		"расчет погрешностей измерений",
		"графики",
		"окончательные результаты",
		"выводы и анализ результатов работы",
		"дополнительные задания",
		"выполнение дополнительных заданий",
		"замечания преподавателя",
		"примечание",
		"приложение",
	}

	lower := strings.ToLower(value)
	for _, keyword := range keywords {
		if strings.Contains(lower, keyword) {
			return "section_heading"
		}
	}
	if numberedHeadingRegexp.MatchString(value) && len([]rune(value)) <= 48 && strings.HasSuffix(value, ".") {
		return "section_heading"
	}
	return "paragraph"
}

func countNonEmptyTables(tables []table) int {
	count := 0
	for _, tbl := range tables {
		found := false
		for _, row := range tbl.Rows {
			for _, cell := range row {
				if cleanText(cell) == "" {
					continue
				}
				count++
				found = true
				break
			}
			if found {
				break
			}
		}
	}
	return count
}

func mergeParagraphFragments(paragraphs []string) string {
	if len(paragraphs) == 0 {
		return ""
	}

	merged := paragraphs[0]
	for _, next := range paragraphs[1:] {
		if next == "" {
			continue
		}
		if shouldConcatenateWithoutSpace(merged, next) {
			merged += next
			continue
		}
		merged += " " + next
	}
	return cleanText(merged)
}

type paragraphSegment struct {
	Kind string
	Text string
}

func joinParagraphSegments(segments []paragraphSegment) string {
	var buf strings.Builder
	prevKind := ""

	for _, segment := range segments {
		text := strings.TrimSpace(segment.Text)
		if text == "" {
			continue
		}
		if buf.Len() == 0 {
			buf.WriteString(text)
			prevKind = segment.Kind
			continue
		}

		if segment.Kind == "formula" && prevKind == "formula" {
			buf.WriteString(" ; ")
		} else {
			buf.WriteString(" ")
		}
		buf.WriteString(text)
		prevKind = segment.Kind
	}

	return normalizeParagraphText(buf.String(), nil)
}

func shouldConcatenateWithoutSpace(left, right string) bool {
	leftRune, ok := lastMeaningfulRune(left)
	if !ok {
		return true
	}
	rightRune, ok := firstMeaningfulRune(right)
	if !ok {
		return true
	}
	if unicode.IsDigit(leftRune) && unicode.IsDigit(rightRune) {
		return true
	}
	if leftRune == ',' && unicode.IsDigit(rightRune) {
		return true
	}
	return false
}

func firstMeaningfulRune(value string) (rune, bool) {
	for _, r := range value {
		if unicode.IsSpace(r) {
			continue
		}
		return r, true
	}
	return 0, false
}

func lastMeaningfulRune(value string) (rune, bool) {
	for len(value) > 0 {
		r, size := utf8.DecodeLastRuneInString(value)
		if r == utf8.RuneError && size == 0 {
			break
		}
		value = value[:len(value)-size]
		if unicode.IsSpace(r) {
			continue
		}
		return r, true
	}
	return 0, false
}

func cleanText(value string) string {
	if value == "" {
		return ""
	}

	value = strings.ReplaceAll(value, "\u00a0", " ")
	value = strings.ReplaceAll(value, "\u00ad", "")
	value = strings.ReplaceAll(value, "\u200b", "")
	value = strings.ReplaceAll(value, "\u200c", "")
	value = strings.ReplaceAll(value, "\u200d", "")
	value = strings.ReplaceAll(value, "\r", " ")
	value = strings.ReplaceAll(value, "\n", " ")
	value = strings.ReplaceAll(value, "\t", " ")
	value = whitespaceRegexp.ReplaceAllString(value, " ")
	value = strings.ReplaceAll(value, "( ", "(")
	value = strings.ReplaceAll(value, "[ ", "[")
	value = strings.ReplaceAll(value, "{ ", "{")
	value = strings.ReplaceAll(value, " )", ")")
	value = strings.ReplaceAll(value, " ]", "]")
	value = strings.ReplaceAll(value, " }", "}")
	value = strings.ReplaceAll(value, " ,", ",")
	value = strings.ReplaceAll(value, " .", ".")
	value = strings.ReplaceAll(value, " :", ":")
	value = strings.ReplaceAll(value, " ;", ";")
	value = strings.ReplaceAll(value, " !", "!")
	value = strings.ReplaceAll(value, " ?", "?")
	value = strings.ReplaceAll(value, " %", "%")
	value = negativeNumberRegexp.ReplaceAllString(value, "$1-$2")
	value = normalizePipeSeparators(value)
	return strings.TrimSpace(value)
}

func normalizeParagraphText(value string, warnings map[string]struct{}) string {
	before := cleanText(value)
	if before == "" {
		return ""
	}

	after := before
	replacements := map[string]string{
		"проницаем ость":           "проницаемость",
		"магнитная проницаем ость": "магнитная проницаемость",
	}
	for oldValue, newValue := range replacements {
		after = strings.ReplaceAll(after, oldValue, newValue)
	}

	if warnings != nil && after != before {
		warnings["split_word_detected"] = struct{}{}
	}

	return after
}

func normalizePipeSeparators(value string) string {
	if !strings.Contains(value, "|") {
		return value
	}

	parts := strings.Split(value, "|")
	cleaned := make([]string, 0, len(parts))
	for _, part := range parts {
		part = whitespaceRegexp.ReplaceAllString(strings.TrimSpace(part), " ")
		if part == "" {
			continue
		}
		cleaned = append(cleaned, part)
	}
	if len(cleaned) == 0 {
		return ""
	}
	return strings.Join(cleaned, " | ")
}

func annotateTables(blocks []bodyBlock, tables []table) []table {
	result := make([]table, 0, len(tables))
	for _, tbl := range tables {
		result = append(result, annotateTable(tbl))
	}
	for blockIdx, block := range blocks {
		if block.Kind != "table" {
			continue
		}
		for idx := range result {
			if result[idx].Index != block.TableIndex {
				continue
			}
			prevText := nearestParagraphText(blocks, blockIdx, -1)
			nextText := nearestParagraphText(blocks, blockIdx, 1)
			caption, source := chooseTableCaption(prevText, nextText)
			if caption != "" {
				result[idx].Caption = caption
				result[idx].CaptionSource = source
			}
			assignTableKind(&result[idx], prevText, nextText)
			break
		}
	}
	return result
}

func annotateTable(tbl table) table {
	maxCols := 0
	for _, row := range tbl.Rows {
		if len(row) > maxCols {
			maxCols = len(row)
		}
	}
	if maxCols == 0 {
		return tbl
	}

	warningSet := make(map[string]struct{})
	normalizedRows := make([][]string, 0, len(tbl.Rows))
	for _, row := range tbl.Rows {
		normalizedRow := append([]string(nil), row...)
		if len(normalizedRow) < maxCols {
			warningSet["merged_cells_ambiguous"] = struct{}{}
			for len(normalizedRow) < maxCols {
				normalizedRow = append(normalizedRow, "")
			}
		}
		for idx, cell := range normalizedRow {
			normalizedRow[idx] = normalizeParagraphText(cell, warningSet)
		}
		normalizedRows = append(normalizedRows, normalizedRow)
	}

	headers := []string(nil)
	if len(normalizedRows) > 0 {
		headers = append([]string(nil), normalizedRows[0]...)
		for idx, header := range headers {
			repaired := repairTableHeader(header)
			if repaired != header {
				warningSet["possible_broken_header"] = struct{}{}
				headers[idx] = repaired
				normalizedRows[0][idx] = repaired
			}
		}
	}

	confidence := 0.94
	if len(warningSet) > 0 {
		confidence -= 0.12 * float64(len(warningSet))
	}
	if confidence < 0.5 {
		confidence = 0.5
	}

	tbl.Headers = headers
	tbl.Rows = normalizedRows
	tbl.Confidence = roundConfidence(confidence)
	tbl.Warnings = warningSetToList(warningSet)
	return tbl
}

func chooseTableCaption(prevText, nextText string) (string, string) {
	prevText = cleanText(prevText)
	nextText = cleanText(nextText)
	switch {
	case tableCaptionRegexp.MatchString(prevText):
		return prevText, "direct_preceding"
	case tableCaptionRegexp.MatchString(nextText):
		return nextText, "direct_following"
	default:
		return "", ""
	}
}

func assignTableKind(tbl *table, prevText, nextText string) {
	if tbl == nil {
		return
	}

	warningSet := make(map[string]struct{}, len(tbl.Warnings))
	for _, warning := range tbl.Warnings {
		warningSet[warning] = struct{}{}
	}

	kind := classifyTableKind(*tbl, prevText, nextText)
	switch kind {
	case "layout_table":
		warningSet["layout_table_detected"] = struct{}{}
	case "metadata_table":
		warningSet["metadata_table_detected"] = struct{}{}
	}
	tbl.Kind = kind
	tbl.Warnings = warningSetToList(warningSet)
}

func classifyTableKind(tbl table, prevText, nextText string) string {
	rows := len(tbl.Rows)
	cols := maxTableColumns(tbl.Rows)
	numericDensity := tableNumericDensity(tbl.Rows)
	imageMarkers := tableContainsImageMarker(tbl.Rows)
	labelDensity := tableLabelDensity(tbl.Rows)

	switch {
	case tbl.SourceKind == "header" || tbl.SourceKind == "footer":
		if rows <= 4 && cols <= 2 {
			return "layout_table"
		}
		if labelDensity >= 0.5 && numericDensity < 0.15 {
			return "metadata_table"
		}
		return "layout_table"
	case imageMarkers:
		return "layout_table"
	case labelDensity >= 0.5 && numericDensity < 0.2 && cols <= 2:
		return "metadata_table"
	case tableCaptionRegexp.MatchString(cleanText(prevText)) || tableCaptionRegexp.MatchString(cleanText(nextText)):
		return "data_table"
	case rows >= 2 && cols >= 2 && numericDensity >= 0.12:
		return "data_table"
	case rows >= 3 && cols >= 3:
		return "data_table"
	default:
		return "layout_table"
	}
}

func maxTableColumns(rows [][]string) int {
	maxCols := 0
	for _, row := range rows {
		if len(row) > maxCols {
			maxCols = len(row)
		}
	}
	return maxCols
}

func tableNumericDensity(rows [][]string) float64 {
	total := 0
	numeric := 0
	for _, row := range rows {
		for _, cell := range row {
			cell = cleanText(cell)
			if cell == "" {
				continue
			}
			total++
			if looksNumericCell(cell) {
				numeric++
			}
		}
	}
	if total == 0 {
		return 0
	}
	return float64(numeric) / float64(total)
}

func tableLabelDensity(rows [][]string) float64 {
	total := 0
	labels := 0
	for _, row := range rows {
		if len(row) == 0 {
			continue
		}
		cell := cleanText(row[0])
		if cell == "" {
			continue
		}
		total++
		if looksLikeLabelCell(cell) {
			labels++
		}
	}
	if total == 0 {
		return 0
	}
	return float64(labels) / float64(total)
}

func tableContainsImageMarker(rows [][]string) bool {
	for _, row := range rows {
		for _, cell := range row {
			if inlineImageMarkerRegexp.MatchString(cell) {
				return true
			}
		}
	}
	return false
}

func looksNumericCell(value string) bool {
	value = cleanText(value)
	if value == "" {
		return false
	}
	digits := 0
	letters := 0
	for _, r := range value {
		switch {
		case unicode.IsDigit(r):
			digits++
		case unicode.IsLetter(r):
			letters++
		}
	}
	return digits > 0 && digits >= letters
}

func looksLikeLabelCell(value string) bool {
	value = cleanText(value)
	if value == "" || len([]rune(value)) > 40 {
		return false
	}
	if pageNumberRegexp.MatchString(value) {
		return false
	}
	colonFree := strings.TrimSuffix(value, ":")
	if localLabelRegexp.MatchString(value) {
		return true
	}
	for _, r := range colonFree {
		if unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

func repairTableHeader(value string) string {
	value = cleanText(value)
	if value == "" {
		return value
	}
	if strings.HasPrefix(value, ", ") {
		return value
	}
	return value
}

type imageAnnotationContext struct {
	CurrentText         string
	PrevText            string
	NextText            string
	CurrentRole         string
	CurrentHeading      string
	SourcePart          string
	SourceKind          string
	ImageCountInContext int
}

func annotateImages(parts []documentPart, images []documentImage) []documentImage {
	if len(images) == 0 {
		return images
	}

	result := make([]documentImage, len(images))
	copy(result, images)

	contexts := make(map[int][]imageAnnotationContext)
	for _, part := range parts {
		currentHeading := ""
		for blockIdx, block := range part.Blocks {
			if block.Kind != "paragraph" {
				continue
			}

			role := classifyParagraphRole(block.Text)
			if role == "section_heading" {
				currentHeading = block.Text
			}
			if len(block.ImageIndexes) == 0 {
				continue
			}

			context := imageAnnotationContext{
				CurrentText:         block.Text,
				PrevText:            nearestParagraphText(part.Blocks, blockIdx, -1),
				NextText:            nearestParagraphText(part.Blocks, blockIdx, 1),
				CurrentRole:         role,
				CurrentHeading:      cleanText(currentHeading),
				SourcePart:          block.SourcePart,
				SourceKind:          block.SourceKind,
				ImageCountInContext: len(block.ImageIndexes),
			}
			for _, imageIndex := range block.ImageIndexes {
				contexts[imageIndex] = append(contexts[imageIndex], context)
			}
		}
	}

	for idx := range result {
		annotateSingleImage(&result[idx], contexts[result[idx].Index])
	}

	return result
}

func annotateSingleImage(imageInfo *documentImage, contexts []imageAnnotationContext) {
	if imageInfo == nil {
		return
	}

	warningSet := make(map[string]struct{}, len(imageInfo.Warnings))
	for _, warning := range imageInfo.Warnings {
		warningSet[warning] = struct{}{}
	}

	caption := cleanText(imageInfo.Caption)
	captionSource := cleanText(imageInfo.CaptionSource)
	captionQuality := cleanText(imageInfo.CaptionQuality)
	ocrText := imageInfo.OCRText

	if caption == "" {
		candidate := chooseCaptionCandidate(imageInfo.Index, imageInfo.SourceKind, contexts)
		if candidate.Text != "" {
			caption = candidate.Text
			captionSource = candidate.Source
			captionQuality = candidate.Quality
			if candidate.Source != "direct" {
				warningSet["caption_inferred"] = struct{}{}
			}
		}
	}

	ocrLines, ocrNormalized := normalizeOCRLinesForImage(imageInfo, ocrText)
	if ocrNormalized {
		warningSet["ocr_scientific_normalized"] = struct{}{}
	}

	if caption == "" {
		for _, line := range ocrLines {
			if !figureCaptionRegexp.MatchString(line) {
				continue
			}
			caption = line
			captionSource = "ocr"
			captionQuality = "medium"
			warningSet["caption_inferred"] = struct{}{}
			break
		}
	}

	filteredOCRLines := make([]string, 0, len(ocrLines))
	for _, line := range ocrLines {
		if figureCaptionRegexp.MatchString(line) {
			if caption == "" {
				caption = line
				captionSource = "ocr"
				captionQuality = "medium"
				warningSet["caption_inferred"] = struct{}{}
			}
			continue
		}
		if caption != "" && line == caption {
			continue
		}
		filteredOCRLines = append(filteredOCRLines, line)
	}
	ocrLines = filteredOCRLines

	imageType, chartTitle, xAxisLabel, yAxisLabel, axisHints := classifyImage(imageInfo, caption, ocrLines, contexts)
	if imageType == "" {
		imageType = "other"
	}
	if caption != "" && captionSource == "" {
		captionSource = "direct"
	}
	if caption != "" && captionQuality == "" {
		captionQuality = "high"
	}

	if shouldSuppressImageOCR(imageInfo, imageType, ocrLines) {
		ocrLines = nil
		warningSet["ocr_suppressed_low_value"] = struct{}{}
	}

	confidence := 0.88
	if caption != "" {
		confidence += 0.04
	}
	if chartTitle != "" || imageType == "chart" {
		confidence += 0.03
	}
	if seemsLowConfidenceOCR(ocrLines) {
		warningSet["ocr_low_confidence"] = struct{}{}
		confidence -= 0.16
	}
	if !imageInfo.UsedInContent {
		warningSet["not_anchored_in_parsed_parts"] = struct{}{}
		warningSet["unresolved_usage"] = struct{}{}
		confidence -= 0.1
	}
	if confidence > 0.98 {
		confidence = 0.98
	}
	if confidence < 0.45 {
		confidence = 0.45
	}

	imageInfo.Caption = caption
	imageInfo.OCRText = strings.Join(ocrLines, "\n")
	imageInfo.ImageType = imageType
	imageInfo.CaptionSource = captionSource
	imageInfo.CaptionQuality = captionQuality
	imageInfo.ChartTitle = chartTitle
	imageInfo.XAxisLabel = xAxisLabel
	imageInfo.YAxisLabel = yAxisLabel
	imageInfo.AxisHints = axisHints
	imageInfo.Confidence = roundConfidence(confidence)
	imageInfo.Warnings = warningSetToList(warningSet)
}

type captionCandidate struct {
	Text    string
	Source  string
	Quality string
}

func chooseCaptionCandidate(imageIndex int, sourceKind string, contexts []imageAnnotationContext) captionCandidate {
	for _, context := range contexts {
		if candidate := findCaptionNearImage(context.CurrentText, imageIndex, context.ImageCountInContext); candidate != "" {
			return captionCandidate{Text: candidate, Source: "direct", Quality: "high"}
		}
	}
	for _, context := range contexts {
		if context.CurrentRole != "figure_label" || context.ImageCountInContext != 1 {
			continue
		}
		if candidate := findCaptionInText(context.CurrentText); candidate != "" {
			return captionCandidate{Text: candidate, Source: "same_paragraph_label", Quality: "high"}
		}
	}
	if sourceKind == "header" || sourceKind == "footer" {
		return captionCandidate{}
	}
	for _, context := range contexts {
		if candidate := findCaptionInText(context.NextText); candidate != "" {
			return captionCandidate{Text: candidate, Source: "nearby_following", Quality: "medium"}
		}
	}
	for _, context := range contexts {
		if candidate := findCaptionInText(context.PrevText); candidate != "" {
			return captionCandidate{Text: candidate, Source: "nearby_preceding", Quality: "medium"}
		}
	}
	return captionCandidate{}
}

func nearestParagraphText(blocks []bodyBlock, start, direction int) string {
	for idx := start + direction; idx >= 0 && idx < len(blocks); idx += direction {
		block := blocks[idx]
		if block.Kind != "paragraph" {
			continue
		}
		text := cleanText(block.Text)
		if text == "" || imageOnlyParagraphRegexp.MatchString(text) {
			continue
		}
		return text
	}
	return ""
}

func findCaptionInText(value string) string {
	value = cleanText(inlineImageMarkerRegexp.ReplaceAllString(value, " "))
	if figureCaptionRegexp.MatchString(value) {
		return value
	}
	return ""
}

func findCaptionNearImage(value string, imageIndex, imageCount int) string {
	marker := formatInlineImageMarker(imageIndex)
	if !strings.Contains(value, marker) {
		return ""
	}
	if imageCount > 1 {
		return ""
	}

	parts := strings.Split(value, marker)
	for _, part := range parts {
		if candidate := findCaptionInText(part); candidate != "" {
			return candidate
		}
	}

	return findCaptionInText(strings.ReplaceAll(value, marker, " "))
}

func normalizeOCRLinesForImage(imageInfo *documentImage, value string) ([]string, bool) {
	if value == "" {
		return nil, false
	}

	changed := false
	seen := make(map[string]struct{})
	lines := make([]string, 0)
	for _, line := range strings.Split(value, "\n") {
		line = cleanText(line)
		if line == "" {
			continue
		}
		normalizedLine := normalizeScientificOCRLine(line, imageInfo)
		if normalizedLine != line {
			changed = true
		}
		key := strings.ToLower(normalizedLine)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		lines = append(lines, normalizedLine)
	}
	return lines, changed
}

func normalizeScientificOCRLine(line string, imageInfo *documentImage) string {
	normalized := cleanText(line)
	if normalized == "" {
		return ""
	}

	replacer := strings.NewReplacer("А/М", "A/m", "A/M", "A/m", "А/M", "A/m", "A/М", "A/m")
	upper := cleanText(strings.ToUpper(normalized))
	if upper == "И(H)" || upper == "И(Н)" || upper == "U(H)" || upper == "U(Н)" {
		return "μ(H)"
	}
	if axisLabelHintRegexp.MatchString(normalized) || strings.Contains(strings.ToLower(normalized), "зависим") {
		normalized = replacer.Replace(normalized)
		switch cleanText(strings.ToUpper(normalized)) {
		case "И(H)", "И(Н)", "U(H)", "U(Н)":
			normalized = "μ(H)"
		case "Н, A/M", "H, A/M", "Н, А/М":
			normalized = "H, A/m"
		}
	}
	if imageInfo != nil && (imageInfo.SourceKind == "header" || imageInfo.SourceKind == "footer") && logoLikeOCRRegexp.MatchString(normalized) {
		normalized = strings.TrimSpace(normalized)
	}
	return normalized
}

func classifyImage(imageInfo *documentImage, caption string, ocrLines []string, contexts []imageAnnotationContext) (string, string, string, string, []string) {
	content := strings.ToLower(strings.Join(append([]string{caption}, ocrLines...), "\n"))
	contextContent := strings.ToLower(joinImageContextTexts(contexts))
	chartTitle, xAxisLabel, yAxisLabel, axisHints, chartScore := extractChartSemantics(caption, ocrLines, contexts)
	brandingScore := brandingSignalScore(imageInfo, ocrLines, contexts)

	if brandingScore >= 3 {
		if hasDecorativeSignals(imageInfo, ocrLines) {
			return "decorative", "", "", "", nil
		}
		return "branding", "", "", "", nil
	}

	switch {
	case chartScore >= 3 && imageInfo.SourceKind != "header" && imageInfo.SourceKind != "footer":
		return "chart", chartTitle, xAxisLabel, yAxisLabel, axisHints
	case strings.Contains(content, "схема") || strings.Contains(content, "r1") || strings.Contains(content, "c1") || strings.Contains(content, "tr"):
		return "circuit", "", "", "", nil
	case strings.Contains(content, "схема") || strings.Contains(contextContent, "схема"):
		return "diagram", "", "", "", nil
	case strings.Contains(content, "общий вид") || strings.Contains(content, "установка"):
		return "photo", "", "", "", nil
	case strings.Contains(content, "магнитопровод") || strings.Contains(content, "сердечник"):
		return "diagram", "", "", "", nil
	case imageInfo.SourceKind == "header" || imageInfo.SourceKind == "footer":
		return "decorative", "", "", "", nil
	default:
		return "other", "", "", "", nil
	}
}

func joinImageContextTexts(contexts []imageAnnotationContext) string {
	parts := make([]string, 0, len(contexts)*2)
	for _, context := range contexts {
		if context.CurrentHeading != "" {
			parts = append(parts, context.CurrentHeading)
		}
		if context.CurrentText != "" {
			parts = append(parts, context.CurrentText)
		}
	}
	return strings.Join(parts, "\n")
}

func brandingSignalScore(imageInfo *documentImage, ocrLines []string, contexts []imageAnnotationContext) int {
	score := 0
	nameContent := strings.ToLower(imageInfo.Name + " " + imageInfo.Path)
	if strings.Contains(nameContent, "logo") || strings.Contains(nameContent, "brand") || strings.Contains(nameContent, "герб") || strings.Contains(nameContent, "emblem") {
		score += 3
	}
	if imageInfo.SourceKind == "header" || imageInfo.SourceKind == "footer" {
		score += 1
	}
	if isLogoLikeOCR(ocrLines) {
		score += 2
	}
	for _, context := range contexts {
		heading := strings.ToLower(context.CurrentHeading)
		if strings.Contains(heading, "график") {
			score--
		}
		if context.SourceKind == "header" && isLogoLikeText(context.CurrentText) {
			score += 1
		}
	}
	return score
}

func hasDecorativeSignals(imageInfo *documentImage, ocrLines []string) bool {
	nameContent := strings.ToLower(imageInfo.Name + " " + imageInfo.Path)
	if strings.Contains(nameContent, "logo") || strings.Contains(nameContent, "brand") || strings.Contains(nameContent, "герб") || strings.Contains(nameContent, "emblem") {
		return false
	}
	return (imageInfo.SourceKind == "header" || imageInfo.SourceKind == "footer") && len(ocrLines) == 0
}

func isLogoLikeOCR(lines []string) bool {
	if len(lines) == 0 || len(lines) > 2 {
		return false
	}
	for _, line := range lines {
		if !logoLikeOCRRegexp.MatchString(line) {
			return false
		}
	}
	return true
}

func isLogoLikeText(value string) bool {
	value = cleanText(inlineImageMarkerRegexp.ReplaceAllString(value, " "))
	return strings.Contains(strings.ToLower(value), "университет") || strings.Contains(strings.ToLower(value), "институт")
}

func containsAxisLikeText(lines []string) bool {
	if countAxisTicks(lines) >= 3 {
		return true
	}
	for _, line := range lines {
		lower := strings.ToLower(line)
		if axisLabelHintRegexp.MatchString(line) || strings.Contains(lower, "кривая") || strings.Contains(lower, "зависимость") {
			return true
		}
	}
	return false
}

func countAxisTicks(lines []string) int {
	count := 0
	for _, line := range lines {
		if looksLikeAxisTick(line) {
			count++
		}
	}
	return count
}

func extractChartSemantics(caption string, lines []string, contexts []imageAnnotationContext) (string, string, string, []string, int) {
	title := ""
	xAxisLabel := ""
	yAxisLabel := ""
	axisHints := make([]string, 0)
	score := 0

	if cleanText(caption) != "" {
		lowerCaption := strings.ToLower(caption)
		if strings.Contains(lowerCaption, "зависимость") || strings.Contains(lowerCaption, "кривая") || strings.Contains(lowerCaption, "график") {
			title = caption
			score += 2
		}
	}

	for _, line := range lines {
		lower := strings.ToLower(line)
		switch {
		case title == "" && (strings.Contains(lower, "зависимость") || strings.Contains(lower, "кривая") || strings.Contains(lower, "график")):
			title = line
			score += 2
		case xAxisLabel == "" && (strings.Contains(lower, "а/м") || strings.Contains(lower, "h,") || strings.Contains(lower, "н,")):
			xAxisLabel = line
			score += 2
		case yAxisLabel == "" && axisLabelRegexp.MatchString(line) && !strings.Contains(lower, "рис."):
			if strings.ContainsAny(line, "Bμβ") || strings.Contains(strings.ToLower(line), "проницаем") {
				yAxisLabel = line
				score += 2
			}
		case looksLikeAxisTick(line):
			axisHints = append(axisHints, line)
		}
	}

	if len(axisHints) > 4 {
		axisHints = axisHints[:4]
	}
	if len(axisHints) >= 3 {
		score++
	}
	for _, context := range contexts {
		heading := strings.ToLower(context.CurrentHeading)
		text := strings.ToLower(context.CurrentText)
		if strings.Contains(heading, "график") || strings.Contains(text, "график") {
			score += 2
			break
		}
		if strings.Contains(text, "зависимость") || strings.Contains(text, "кривая") {
			score++
		}
	}

	return title, xAxisLabel, yAxisLabel, axisHints, score
}

func looksLikeAxisTick(value string) bool {
	value = cleanText(value)
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

func seemsLowConfidenceOCR(lines []string) bool {
	if len(lines) == 0 {
		return false
	}
	shortLines := 0
	for _, line := range lines {
		if len([]rune(line)) <= 3 {
			shortLines++
		}
	}
	return shortLines >= len(lines)/2 && len(lines) >= 4
}

func shouldSuppressImageOCR(imageInfo *documentImage, imageType string, lines []string) bool {
	if len(lines) == 0 {
		return false
	}
	switch imageType {
	case "branding", "decorative":
		return isLowValueOCR(lines, imageInfo.SourceKind)
	default:
		return false
	}
}

func isLowValueOCR(lines []string, sourceKind string) bool {
	if len(lines) == 0 {
		return true
	}
	if isLogoLikeOCR(lines) && (sourceKind == "header" || sourceKind == "footer") {
		return true
	}
	shortLines := 0
	for _, line := range lines {
		if len([]rune(line)) <= 3 {
			shortLines++
		}
	}
	return shortLines == len(lines)
}

func finalizeSection(value section) section {
	warningSet := make(map[string]struct{})
	if len(value.Warnings) > 0 {
		for _, warning := range value.Warnings {
			warningSet[warning] = struct{}{}
		}
	}

	confidenceSum := 0.0
	confidenceCount := 0.0
	for _, block := range value.Blocks {
		if block.Confidence > 0 {
			confidenceSum += block.Confidence
			confidenceCount++
		}
		for _, warning := range block.Warnings {
			warningSet[warning] = struct{}{}
		}
	}

	if confidenceCount > 0 {
		value.Confidence = roundConfidence(confidenceSum / confidenceCount)
	}
	value.Warnings = warningSetToList(warningSet)
	return value
}

func roundConfidence(value float64) float64 {
	if value <= 0 {
		return 0
	}
	return float64(int(value*100+0.5)) / 100
}

func maxFloat(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func readZipFile(files []*zip.File, name string) ([]byte, error) {
	data, err := readZipFileOptional(files, name)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("required file %s not found in docx", name)
	}
	return data, nil
}

func readZipFileOptional(files []*zip.File, name string) ([]byte, error) {
	cleanName := path.Clean(name)
	for _, file := range files {
		if path.Clean(file.Name) != cleanName {
			continue
		}

		rc, err := file.Open()
		if err != nil {
			return nil, fmt.Errorf("open %s: %w", cleanName, err)
		}

		data, err := io.ReadAll(rc)
		_ = rc.Close()
		if err != nil {
			return nil, fmt.Errorf("read %s: %w", cleanName, err)
		}
		return data, nil
	}
	return nil, nil
}

func relationshipsPath(partPath string) string {
	partPath = path.Clean(partPath)
	return path.Join(path.Dir(partPath), "_rels", path.Base(partPath)+".rels")
}

func resolveRelationshipTarget(partPath string, rel relationshipXML) (string, bool) {
	target := strings.TrimSpace(rel.Target)
	if target == "" || strings.EqualFold(strings.TrimSpace(rel.TargetMode), "External") {
		return "", false
	}
	if strings.Contains(target, "://") {
		return "", false
	}
	if strings.HasPrefix(target, "/") {
		return path.Clean(strings.TrimPrefix(target, "/")), true
	}
	return path.Clean(path.Join(path.Dir(path.Clean(partPath)), target)), true
}

func shouldTraversePart(rel relationshipXML, targetPath string) bool {
	targetPath = strings.ToLower(path.Clean(targetPath))
	if !strings.HasPrefix(targetPath, "word/") || path.Ext(targetPath) != ".xml" {
		return false
	}

	lowerType := strings.ToLower(rel.Type)
	switch {
	case strings.Contains(lowerType, "/header"),
		strings.Contains(lowerType, "/footer"),
		strings.Contains(lowerType, "/footnotes"),
		strings.Contains(lowerType, "/endnotes"),
		strings.Contains(lowerType, "/comments"),
		strings.Contains(lowerType, "/glossarydocument"),
		strings.Contains(lowerType, "/subdocument"):
		return true
	}

	base := strings.ToLower(path.Base(targetPath))
	switch {
	case base == "document.xml",
		strings.HasPrefix(base, "header"),
		strings.HasPrefix(base, "footer"),
		strings.HasPrefix(base, "footnote"),
		strings.HasPrefix(base, "endnote"),
		strings.HasPrefix(base, "comment"),
		strings.HasPrefix(base, "glossarydocument"),
		strings.HasPrefix(base, "textbox"),
		strings.HasPrefix(base, "subdoc"):
		return true
	}

	return false
}

func classifyDocumentPart(partPath, relType string) (string, string) {
	base := strings.TrimSuffix(path.Base(path.Clean(partPath)), path.Ext(partPath))
	lowerBase := strings.ToLower(base)
	lowerType := strings.ToLower(relType)

	switch {
	case lowerBase == "document":
		return "document", "document"
	case strings.Contains(lowerType, "/header") || strings.HasPrefix(lowerBase, "header"):
		return base, "header"
	case strings.Contains(lowerType, "/footer") || strings.HasPrefix(lowerBase, "footer"):
		return base, "footer"
	case strings.Contains(lowerType, "/footnotes") || strings.HasPrefix(lowerBase, "footnote"):
		return base, "footnotes"
	case strings.Contains(lowerType, "/endnotes") || strings.HasPrefix(lowerBase, "endnote"):
		return base, "endnotes"
	case strings.Contains(lowerType, "/comments") || strings.HasPrefix(lowerBase, "comment"):
		return base, "comments"
	case strings.Contains(lowerType, "/glossarydocument") || strings.HasPrefix(lowerBase, "glossarydocument"):
		return base, "glossary"
	default:
		return base, "related"
	}
}

func isRecursiveContentContainer(local string) bool {
	switch local {
	case "sdt", "sdtContent", "customXml", "smartTag", "ins", "del", "moveFrom", "moveTo", "comment", "footnote", "endnote", "AlternateContent", "Choice", "Fallback", "hdr", "ftr", "txbxContent":
		return true
	default:
		return false
	}
}

func detectImagePlacement(stack []string) (string, []string) {
	if len(stack) == 0 {
		return "", nil
	}

	hints := make(map[string]struct{})
	has := func(name string) bool {
		for _, current := range stack {
			if current == name {
				return true
			}
		}
		return false
	}

	placement := ""
	switch {
	case has("anchor"):
		placement = "anchored"
	case has("inline"):
		placement = "inline"
	case has("pict") || has("imagedata"):
		placement = "vml"
	}

	if has("anchor") {
		hints["anchor"] = struct{}{}
	}
	if has("inline") {
		hints["inline"] = struct{}{}
	}
	if has("pict") || has("imagedata") {
		hints["vml"] = struct{}{}
	}
	if has("AlternateContent") {
		hints["alternate_content"] = struct{}{}
	}
	if has("txbxContent") || has("textbox") {
		hints["text_box"] = struct{}{}
	}

	return placement, warningSetToList(hints)
}

func scopedRelationshipKey(partPath, relID string) string {
	return path.Clean(partPath) + "#" + strings.TrimSpace(relID)
}

func countArchivedImages(files []*zip.File) int {
	count := 0
	for _, file := range files {
		name := strings.ToLower(path.Clean(file.Name))
		if !strings.HasPrefix(name, "word/media/") {
			continue
		}
		if strings.HasSuffix(name, "/") {
			continue
		}
		count++
	}
	return count
}

func detectImageContentType(filePath string, data []byte) string {
	if contentType := mime.TypeByExtension(strings.ToLower(path.Ext(filePath))); contentType != "" {
		return contentType
	}
	if len(data) == 0 {
		return ""
	}
	return strings.TrimSpace(strings.SplitN(http.DetectContentType(data), ";", 2)[0])
}

func attrValue(attrs []xml.Attr, local string) string {
	for _, attr := range attrs {
		if attr.Name.Local == local {
			return strings.TrimSpace(attr.Value)
		}
	}
	return ""
}

func wrapInnerXML(innerXML string) string {
	return "<root>" + innerXML + "</root>"
}

func formatInlineImageMarker(index int) string {
	return fmt.Sprintf("[IMAGE %d]", index)
}
