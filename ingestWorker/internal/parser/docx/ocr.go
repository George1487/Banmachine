package docx

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
)

type ocrResult struct {
	Path string `json:"path"`
	Text string `json:"text"`
}

func enrichImagesWithOCR(ctx context.Context, images []*documentImage) {
	if strings.TrimSpace(os.Getenv("INGESTWORKER_DISABLE_OCR")) != "" {
		return
	}

	ordered := orderedImages(images)
	if len(ordered) == 0 {
		return
	}

	tmpDir, fileMap, err := writeTempImages(ordered)
	if err != nil {
		return
	}
	defer os.RemoveAll(tmpDir)

	_ = runTesseractOCR(ctx, fileMap)
	_ = runVisionOCR(ctx, tmpDir, fileMap)
}

func orderedImages(images []*documentImage) []*documentImage {
	result := make([]*documentImage, 0, len(images))
	for _, imageInfo := range images {
		if imageInfo == nil || len(imageInfo.data) == 0 {
			continue
		}
		result = append(result, imageInfo)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Path < result[j].Path
	})
	return result
}

func writeTempImages(images []*documentImage) (string, map[string]*documentImage, error) {
	tmpDir, err := os.MkdirTemp("", "ingestworker-docx-images-")
	if err != nil {
		return "", nil, err
	}

	fileMap := make(map[string]*documentImage, len(images))
	for idx, imageInfo := range images {
		fileName := imageInfo.Name
		if fileName == "" {
			fileName = fmt.Sprintf("image-%d.bin", idx+1)
		}
		fileName = fmt.Sprintf("%03d_%s", idx+1, filepath.Base(fileName))
		filePath := filepath.Join(tmpDir, fileName)
		if err := os.WriteFile(filePath, imageInfo.data, 0o600); err != nil {
			return "", nil, err
		}
		fileMap[filePath] = imageInfo
	}
	return tmpDir, fileMap, nil
}

func runTesseractOCR(ctx context.Context, fileMap map[string]*documentImage) bool {
	if _, err := exec.LookPath("tesseract"); err != nil {
		return false
	}

	hadSuccess := false
	for filePath, imageInfo := range fileMap {
		cmd := exec.CommandContext(ctx, "tesseract", filePath, "stdout", "-l", "rus+eng", "--psm", "6")
		output, err := cmd.Output()
		if err != nil {
			continue
		}
		text := normalizeOCRText(string(output))
		if text == "" {
			continue
		}
		applyOCRResult(imageInfo, text, "tesseract")
		hadSuccess = true
	}
	return hadSuccess
}

func runVisionOCR(ctx context.Context, tmpDir string, fileMap map[string]*documentImage) bool {
	if runtime.GOOS != "darwin" {
		return false
	}
	if _, err := exec.LookPath("swift"); err != nil {
		return false
	}

	scriptPath := filepath.Join(tmpDir, "ocr_vision.swift")
	if err := os.WriteFile(scriptPath, []byte(visionOCRScript), 0o600); err != nil {
		return false
	}

	paths := make([]string, 0, len(fileMap))
	for filePath := range fileMap {
		paths = append(paths, filePath)
	}
	sort.Strings(paths)

	moduleCacheDir := filepath.Join(tmpDir, "swift-module-cache")
	_ = os.MkdirAll(moduleCacheDir, 0o700)

	args := append([]string{scriptPath}, paths...)
	cmd := exec.CommandContext(ctx, "swift", args...)
	cmd.Env = append(os.Environ(),
		"CLANG_MODULE_CACHE_PATH="+moduleCacheDir,
		"SWIFT_MODULECACHE_PATH="+moduleCacheDir,
	)
	output, err := cmd.Output()
	if err != nil {
		return false
	}

	var results []ocrResult
	if err := json.Unmarshal(output, &results); err != nil {
		return false
	}

	hadSuccess := false
	for _, result := range results {
		imageInfo := fileMap[result.Path]
		if imageInfo == nil {
			continue
		}
		text := normalizeOCRText(result.Text)
		if text == "" {
			continue
		}
		applyOCRResult(imageInfo, text, "vision")
		hadSuccess = true
	}

	return hadSuccess
}

func applyOCRResult(imageInfo *documentImage, text, engine string) {
	if imageInfo == nil || text == "" {
		return
	}
	if ocrUtilityScore(text) < ocrUtilityScore(imageInfo.OCRText) {
		return
	}
	imageInfo.OCRText = text
	imageInfo.OCREngine = engine
}

func ocrUtilityScore(value string) int {
	lines := strings.Split(value, "\n")
	score := 0
	for _, line := range lines {
		line = cleanText(line)
		if line == "" {
			continue
		}
		score += len([]rune(line))
		if axisLabelHintRegexp.MatchString(line) {
			score += 8
		}
		if figureCaptionRegexp.MatchString(line) {
			score += 4
		}
		if looksLikeAxisTick(line) {
			score += 2
		}
	}
	return score
}

func normalizeOCRText(value string) string {
	value = strings.ReplaceAll(value, "\r\n", "\n")
	value = strings.ReplaceAll(value, "\r", "\n")

	lines := strings.Split(value, "\n")
	cleaned := make([]string, 0, len(lines))
	for _, line := range lines {
		line = cleanText(line)
		if line == "" {
			continue
		}
		cleaned = append(cleaned, line)
	}
	return strings.Join(cleaned, "\n")
}

const visionOCRScript = `
import Foundation
import Vision
import CoreGraphics
import ImageIO

struct OCRResult: Encodable {
    let path: String
    let text: String
}

func recognizeText(at path: String) -> String {
    let url = URL(fileURLWithPath: path)
    guard let source = CGImageSourceCreateWithURL(url as CFURL, nil),
          let image = CGImageSourceCreateImageAtIndex(source, 0, nil) else {
        return ""
    }

    let request = VNRecognizeTextRequest()
    request.recognitionLevel = .accurate
    request.usesLanguageCorrection = true
    request.recognitionLanguages = ["ru-RU", "en-US"]

    let handler = VNImageRequestHandler(cgImage: image, options: [:])
    do {
        try handler.perform([request])
    } catch {
        return ""
    }

    let lines = (request.results ?? []).compactMap { observation in
        observation.topCandidates(1).first?.string
    }
    return lines.joined(separator: "\n")
}

let paths = Array(CommandLine.arguments.dropFirst())
let results = paths.map { path in
    OCRResult(path: path, text: recognizeText(at: path))
}
let encoder = JSONEncoder()
if let data = try? encoder.encode(results) {
    FileHandle.standardOutput.write(data)
}
`
