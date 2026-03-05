package main

import (
	"os"
	"path/filepath"
	"testing"
)

// --- parseArgs ---

func TestParseArgs_BasicPrompt(t *testing.T) {
	opts, err := parseArgs([]string{"a", "cat"})
	if err != nil {
		t.Fatal(err)
	}
	if opts.Prompt != "a cat" {
		t.Errorf("prompt = %q, want %q", opts.Prompt, "a cat")
	}
	if opts.Size != "1K" {
		t.Errorf("size = %q, want %q", opts.Size, "1K")
	}
	if opts.OutputMode != modeHuman {
		t.Errorf("mode = %q, want %q", opts.OutputMode, modeHuman)
	}
}

func TestParseArgs_AllFlags(t *testing.T) {
	opts, err := parseArgs([]string{
		"-o", "out", "-s", "2K", "-a", "16:9",
		"-m", "pro", "-d", "/tmp", "-r", "ref.png",
		"-t", "--api-key", "key123", "--json",
		"hello", "world",
	})
	if err != nil {
		t.Fatal(err)
	}
	if opts.Output != "out" {
		t.Errorf("output = %q", opts.Output)
	}
	if opts.Size != "2K" {
		t.Errorf("size = %q", opts.Size)
	}
	if opts.AspectRatio != "16:9" {
		t.Errorf("aspect = %q", opts.AspectRatio)
	}
	if opts.Model != "gemini-3-pro-image-preview" {
		t.Errorf("model = %q", opts.Model)
	}
	if opts.OutputDir != "/tmp" {
		t.Errorf("dir = %q", opts.OutputDir)
	}
	if len(opts.References) != 1 || opts.References[0] != "ref.png" {
		t.Errorf("refs = %v", opts.References)
	}
	if !opts.Transparent {
		t.Error("transparent should be true")
	}
	if opts.APIKey != "key123" {
		t.Errorf("apiKey = %q", opts.APIKey)
	}
	if opts.OutputMode != modeJSON {
		t.Errorf("mode = %q", opts.OutputMode)
	}
	if opts.Prompt != "hello world" {
		t.Errorf("prompt = %q", opts.Prompt)
	}
}

func TestParseArgs_PlainMode(t *testing.T) {
	opts, err := parseArgs([]string{"--plain", "test"})
	if err != nil {
		t.Fatal(err)
	}
	if opts.OutputMode != modePlain {
		t.Errorf("mode = %q, want %q", opts.OutputMode, modePlain)
	}
}

func TestParseArgs_CostsMode(t *testing.T) {
	opts, err := parseArgs([]string{"--costs"})
	if err != nil {
		t.Fatal(err)
	}
	if !opts.ShowCosts {
		t.Error("ShowCosts should be true")
	}
}

func TestParseArgs_JQ(t *testing.T) {
	opts, err := parseArgs([]string{"--jq", ".files", "test"})
	if err != nil {
		t.Fatal(err)
	}
	if opts.JQ != ".files" {
		t.Errorf("jq = %q", opts.JQ)
	}
}

func TestParseArgs_NoPrompt(t *testing.T) {
	_, err := parseArgs([]string{"-s", "1K"})
	if err == nil {
		t.Error("expected error for no prompt")
	}
}

func TestParseArgs_UnknownOption(t *testing.T) {
	_, err := parseArgs([]string{"--bogus", "test"})
	if err == nil {
		t.Error("expected error for unknown option")
	}
}

func TestParseArgs_MissingValue(t *testing.T) {
	flags := []string{"--output", "--size", "--aspect", "--model", "--dir", "--ref", "--api-key", "--jq"}
	for _, f := range flags {
		_, err := parseArgs([]string{f})
		if err == nil {
			t.Errorf("expected error for %s without value", f)
		}
	}
}

func TestParseArgs_InvalidSize(t *testing.T) {
	_, err := parseArgs([]string{"-s", "3K", "test"})
	if err == nil {
		t.Error("expected error for invalid size")
	}
}

func TestParseArgs_InvalidAspect(t *testing.T) {
	_, err := parseArgs([]string{"-a", "7:3", "test"})
	if err == nil {
		t.Error("expected error for invalid aspect")
	}
}

func TestParseArgs_ProModel512Upgrade(t *testing.T) {
	opts, err := parseArgs([]string{"-s", "512", "-m", "pro", "test"})
	if err != nil {
		t.Fatal(err)
	}
	if opts.Size != "1K" {
		t.Errorf("size = %q, want 1K (pro model 512 upgrade)", opts.Size)
	}
}

// --- resolveModel ---

func TestResolveModel_Alias(t *testing.T) {
	tests := map[string]string{
		"flash":  "gemini-3.1-flash-image-preview",
		"Flash":  "gemini-3.1-flash-image-preview",
		"pro":    "gemini-3-pro-image-preview",
		"nb2":    "gemini-3.1-flash-image-preview",
		"nb-pro": "gemini-3-pro-image-preview",
	}
	for input, want := range tests {
		got := resolveModel(input)
		if got != want {
			t.Errorf("resolveModel(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestResolveModel_Passthrough(t *testing.T) {
	input := "gemini-custom-model"
	if got := resolveModel(input); got != input {
		t.Errorf("resolveModel(%q) = %q, want passthrough", input, got)
	}
}

// --- imageSize ---

func TestImageSize_512(t *testing.T) {
	if got := imageSize("512"); got != "1K" {
		t.Errorf("imageSize(512) = %q, want 1K", got)
	}
}

func TestImageSize_Passthrough(t *testing.T) {
	for _, s := range []string{"1K", "2K", "4K"} {
		if got := imageSize(s); got != s {
			t.Errorf("imageSize(%q) = %q, want %q", s, got, s)
		}
	}
}

// --- calculateCost ---

func TestCalculateCost_Flash(t *testing.T) {
	cost := calculateCost("gemini-3.1-flash-image-preview", 1000, 1000)
	if cost <= 0 {
		t.Errorf("cost = %f, want > 0", cost)
	}
	// 1000 tokens: (1000/1e6)*0.25 + (1000/1e6)*60 = 0.00025 + 0.06 = 0.06025
	expected := 0.06025
	if diff := cost - expected; diff > 1e-9 || diff < -1e-9 {
		t.Errorf("cost = %.10f, want %.10f", cost, expected)
	}
}

func TestCalculateCost_Pro(t *testing.T) {
	cost := calculateCost("gemini-3-pro-image-preview", 1000, 1000)
	// (1000/1e6)*2.0 + (1000/1e6)*120 = 0.002 + 0.12 = 0.122
	expected := 0.122
	if diff := cost - expected; diff > 1e-9 || diff < -1e-9 {
		t.Errorf("cost = %.10f, want %.10f", cost, expected)
	}
}

func TestCalculateCost_UnknownModel(t *testing.T) {
	cost := calculateCost("unknown-model", 1000, 1000)
	// Falls back to defaultModel rates (flash)
	expected := 0.06025
	if diff := cost - expected; diff > 1e-9 || diff < -1e-9 {
		t.Errorf("cost = %.10f, want %.10f (fallback to default)", cost, expected)
	}
}

// --- extFromMime ---

func TestExtFromMime_Known(t *testing.T) {
	got := extFromMime("image/png")
	if got != ".png" {
		t.Errorf("extFromMime(image/png) = %q, want .png", got)
	}
}

func TestExtFromMime_Fallback(t *testing.T) {
	got := extFromMime("image/webp")
	// Should return something reasonable
	if got == "" {
		t.Error("extFromMime(image/webp) returned empty")
	}
}

func TestExtFromMime_NoSlash(t *testing.T) {
	got := extFromMime("bogus")
	if got != ".png" {
		t.Errorf("extFromMime(bogus) = %q, want .png", got)
	}
}

func TestExtFromMime_UnknownSubtype(t *testing.T) {
	got := extFromMime("image/xyzzy")
	if got != ".xyzzy" {
		t.Errorf("extFromMime(image/xyzzy) = %q, want .xyzzy", got)
	}
}

// --- mustWd ---

func TestMustWd(t *testing.T) {
	wd := mustWd()
	if wd == "" {
		t.Error("mustWd() returned empty")
	}
}

// --- nullable ---

func TestNullable_Empty(t *testing.T) {
	if nullable("") != nil {
		t.Error("nullable(\"\") should return nil")
	}
}

func TestNullable_NonEmpty(t *testing.T) {
	p := nullable("hello")
	if p == nil {
		t.Fatal("nullable(\"hello\") should not return nil")
	}
	if *p != "hello" {
		t.Errorf("nullable(\"hello\") = %q", *p)
	}
}

// --- readDotEnvValue ---

func TestReadDotEnvValue_ReadsKey(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	content := "# comment\n\nGEMINI_API_KEY=test-key-123\nOTHER=val\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	got := readDotEnvValue(path, "GEMINI_API_KEY")
	if got != "test-key-123" {
		t.Errorf("readDotEnvValue = %q, want %q", got, "test-key-123")
	}
}

func TestReadDotEnvValue_QuotedValue(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	content := `GEMINI_API_KEY="quoted-key"` + "\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	got := readDotEnvValue(path, "GEMINI_API_KEY")
	if got != "quoted-key" {
		t.Errorf("readDotEnvValue = %q, want %q", got, "quoted-key")
	}
}

func TestReadDotEnvValue_MissingKey(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	content := "OTHER=val\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	got := readDotEnvValue(path, "GEMINI_API_KEY")
	if got != "" {
		t.Errorf("readDotEnvValue = %q, want empty", got)
	}
}

func TestReadDotEnvValue_MissingFile(t *testing.T) {
	got := readDotEnvValue("/nonexistent/path/.env", "GEMINI_API_KEY")
	if got != "" {
		t.Errorf("readDotEnvValue = %q, want empty", got)
	}
}

// --- readStdinPipe ---

func TestReadStdinPipe_FromPipe(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	orig := os.Stdin
	os.Stdin = r
	defer func() { os.Stdin = orig }()

	go func() {
		w.Write([]byte("hello from pipe\n"))
		w.Close()
	}()

	got, err := readStdinPipe()
	if err != nil {
		t.Fatal(err)
	}
	if got != "hello from pipe" {
		t.Errorf("readStdinPipe() = %q, want %q", got, "hello from pipe")
	}
}
