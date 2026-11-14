// Command coverage runs service-agnostic isolated Go unit coverage.
// Generated source is temporarily materialized because Go's cover tool does
// not read virtual-only files from the build overlay.
package main

import (
    "bufio"
    "context"
    "encoding/json"
    "errors"
    "flag"
    "fmt"
    "io"
    "os"
    "os/exec"
    "path/filepath"
    "sort"
    "strconv"
    "strings"
)

type target struct {
    name, owner, pattern string
}

type overlayConfig struct {
    Replace map[string]string
}

func main() {
    if err := run(); err != nil {
        fmt.Fprintln(os.Stderr, "[ERROR]", err)
        os.Exit(1)
    }
}

func run() error {
    backend := flag.String("backend", ".", "backend module directory")
    mode := flag.String("mode", "--strict", "strict, report-only, or inventory")
    single := flag.String("package", "", "one service/core/launcher name")
    multiple := flag.String("packages", "", "comma separated or bracketed names")
    flag.Parse()
    root, err := filepath.Abs(*backend)
    if err != nil { return err }
    targets, err := resolve(root, *single, *multiple)
    if err != nil { return err }
    if *mode != "--strict" && *mode != "--report-only" && *mode != "--inventory" {
        return fmt.Errorf("unknown unit testing mode %q", *mode)
    }
    fmt.Printf("[INFO] Unit-test packages: %s\n", names(targets))
    if *mode == "--inventory" {
        for _, item := range targets {
            fmt.Printf("[INFO] Method inventory: %s\n", item.name)
            if err := command(context.Background(), root, "go", "run", "./libs/core/annotation/generation/testing/inventory", item.owner); err != nil { return err }
        }
        return nil
    }
    generated := filepath.Join(root, ".generated", "unit-coverage")
    if err := os.MkdirAll(generated, 0755); err != nil { return err }
    staleManifest := filepath.Join(generated, ".materialized-files")
    if err := recoverGenerated(staleManifest, root); err != nil { return err }
    if err := command(context.Background(), root, "bash", filepath.Join(filepath.Dir(root), "scripts", "backend", "generate-annotations.sh")); err != nil { return err }
    overlay := filepath.Join(root, ".generated", "overlay.json")
    cleanup, err := materialize(overlay, staleManifest, targets)
    if err != nil { return err }
    defer cleanup()
    reports := filepath.Join(filepath.Dir(root), ".runtime", "reports", "unit-coverage")
    if err := os.MkdirAll(reports, 0755); err != nil { return err }
    failed := false
    for _, item := range targets {
        if err := testTarget(root, overlay, generated, reports, item, *mode); err != nil {
            fmt.Fprintf(os.Stderr, "[ERROR] %s: %v\n", item.name, err)
            failed = true
        }
    }
    if failed { return errors.New("one or more selected unit coverage runs failed") }
    return nil
}

func names(items []target) string {
    values := make([]string, 0, len(items))
    for _, item := range items { values = append(values, item.name) }
    return strings.Join(values, ", ")
}

func parseNames(value string) []string {
    value = strings.TrimSpace(value)
    if strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]") {
        value = strings.TrimSpace(value[1:len(value)-1])
    }
    var names []string
    for _, part := range strings.Split(value, ",") {
        name := strings.Trim(strings.TrimSpace(part), " \t\r\n\"'")
        if name != "" { names = append(names, name) }
    }
    return names
}

func resolve(root, single, multiple string) ([]target, error) {
    if strings.TrimSpace(single) != "" && strings.TrimSpace(multiple) != "" {
        return nil, errors.New("use package= or packages=, not both")
    }
    requested := parseNames(single)
    if len(requested) == 0 { requested = parseNames(multiple) }
    if len(requested) == 0 || (len(requested) == 1 && requested[0] == "all") {
        requested = []string{"core", "launcher"}
        services, err := os.ReadDir(filepath.Join(root, "services"))
        if err != nil { return nil, err }
        for _, entry := range services {
            if entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") {
                requested = append(requested, entry.Name())
            }
        }
    }
    unique := make(map[string]bool)
    var result []target
    for _, name := range requested {
        if name == "all" { return nil, errors.New("all cannot be combined with selected packages") }
        for _, r := range name {
            if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_') {
                return nil, fmt.Errorf("invalid unit-coverage package %q", name)
            }
        }
        if unique[name] { continue }
        unique[name] = true
        item := target{name: name, owner: filepath.Join(root, "services", name), pattern: "./services/" + name + "/..."}
        switch name {
        case "core":
            item.owner, item.pattern = filepath.Join(root, "libs", "core"), "./libs/core/..."
        case "launcher":
            item.owner, item.pattern = filepath.Join(root, "launcher"), "./launcher/..."
        }
        info, err := os.Stat(item.owner)
        if err != nil || !info.IsDir() { return nil, fmt.Errorf("unknown package %q (missing %s)", name, item.owner) }
        result = append(result, item)
    }
    sort.Slice(result, func(i, j int) bool { return result[i].name < result[j].name })
    return result, nil
}

func command(ctx context.Context, dir, name string, args ...string) error {
    cmd := exec.CommandContext(ctx, name, args...)
    cmd.Dir = dir
    cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
    if err := cmd.Run(); err != nil { return fmt.Errorf("%s %s: %w", name, strings.Join(args, " "), err) }
    return nil
}

func within(root, path string) bool {
    rel, err := filepath.Rel(root, path)
    return err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(os.PathSeparator))
}

func recoverGenerated(manifest, root string) error {
    data, err := os.ReadFile(manifest)
    if os.IsNotExist(err) { return nil }
    if err != nil { return err }
    for _, line := range strings.Split(string(data), "\n") {
        if line == "" { continue }
        if !within(root, line) { return fmt.Errorf("unsafe stale generated path %s", line) }
        content, err := os.ReadFile(line)
        if os.IsNotExist(err) { continue }
        if err != nil { return err }
        if !strings.HasPrefix(filepath.Base(line), "zz_") || !strings.Contains(string(content[:min(160, len(content))]), "Code generated") {
            return fmt.Errorf("refusing to remove a changed source file after interrupted coverage: %s", line)
        }
        if err := os.Remove(line); err != nil { return err }
    }
    return os.Remove(manifest)
}

func materialize(overlay, manifest string, targets []target) (func(), error) {
    data, err := os.ReadFile(overlay)
    if err != nil { return nil, err }
    var configuration overlayConfig
    if err := json.Unmarshal(data, &configuration); err != nil { return nil, err }
    logicals := make([]string, 0, len(configuration.Replace))
    for logical := range configuration.Replace { logicals = append(logicals, logical) }
    sort.Strings(logicals)
    created := []string{}
    clean := func() {
        for i := len(created)-1; i >= 0; i-- { _ = os.Remove(created[i]) }
        _ = os.Remove(manifest)
    }
    if err := os.WriteFile(manifest, []byte{}, 0600); err != nil { return nil, err }
    for _, logical := range logicals {
        selected := false
        for _, item := range targets {
            if within(item.owner, logical) { selected = true; break }
        }
        if !selected { continue }
        source := configuration.Replace[logical]
        if _, err := os.Stat(logical); err == nil {
            clean()
            return nil, fmt.Errorf("refusing to overwrite maintained source %s", logical)
        } else if !os.IsNotExist(err) {
            clean()
            return nil, err
        }
        if err := os.MkdirAll(filepath.Dir(logical), 0755); err != nil { clean(); return nil, err }
        input, err := os.Open(source)
        if err != nil { clean(); return nil, err }
        output, err := os.OpenFile(logical, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
        if err != nil { _ = input.Close(); clean(); return nil, err }
        created = append(created, logical)
        f, err := os.OpenFile(manifest, os.O_WRONLY|os.O_APPEND, 0600)
        if err == nil {
            _, err = fmt.Fprintln(f, logical)
            _ = f.Close()
        }
        if err != nil { _ = input.Close(); _ = output.Close(); clean(); return nil, err }
        _, err = io.Copy(output, input)
        _ = input.Close()
        if closeErr := output.Close(); err == nil { err = closeErr }
        if err != nil { clean(); return nil, err }
    }
    return clean, nil
}

func listProduction(root, overlay string, item target) ([]string, error) {
    cmd := exec.Command("go", "list", "-overlay="+overlay, "-f", "{{.ImportPath}}", item.pattern)
    cmd.Dir, cmd.Stderr = root, os.Stderr
    data, err := cmd.Output()
    if err != nil { return nil, err }
    var result []string
    for _, pkg := range strings.Fields(string(data)) {
        if strings.HasSuffix(pkg, "/tests") || strings.Contains(pkg, "/tests/") || strings.HasSuffix(pkg, "/testdata") || strings.Contains(pkg, "/testdata/") { continue }
        result = append(result, pkg)
    }
    if len(result) == 0 { return nil, errors.New("no production packages found") }
    return result, nil
}

func testTarget(root, overlay, generated, reports string, item target, mode string) error {
    production, err := listProduction(root, overlay, item)
    if err != nil { return err }
    rel, err := filepath.Rel(root, item.owner)
    if err != nil { return err }
    candidates := []string{
        filepath.Join(item.owner, "internal", "tests", "unit-test"),
        filepath.Join(item.owner, "tests", "unit-test"),
    }
    arguments := []string{"test", "-count=1", "-vet=off", "-overlay="+overlay, "-covermode=atomic", "-coverpkg="+strings.Join(production, ",")}
    profile := filepath.Join(generated, item.name+".out")
    arguments = append(arguments, "-coverprofile="+profile)
    arguments = append(arguments, production...)
    tests := 0
    for _, dir := range candidates {
        if info, err := os.Stat(dir); err == nil && info.IsDir() {
            sub, err := filepath.Rel(root, dir)
            if err != nil { return err }
            arguments = append(arguments, "./"+filepath.ToSlash(sub)+"/...")
            tests++
        }
    }
    if tests == 0 { fmt.Printf("[WARN] %s has no unit-test directory; its production code is still counted\n", item.name) }
    fmt.Printf("[INFO] Testing %s: %d production packages, %d unit-test trees (%s)\n", item.name, len(production), tests, rel)
    if err := command(context.Background(), root, "go", arguments...); err != nil { return err }
    total, covered, err := score(profile)
    if err != nil { return err }
    if total == 0 { return errors.New("coverage profile has zero instrumented statements") }
    percent := 100 * float64(covered) / float64(total)
    fmt.Printf("[INFO] %s: %.2f%% statement coverage (%d/%d), profile %s\n", item.name, percent, covered, total, profile)
    html := filepath.Join(reports, item.name+".html")
    if err := command(context.Background(), root, "go", "tool", "cover", "-html="+profile, "-o="+html); err == nil {
        fmt.Printf("[INFO] %s coverage HTML: %s\n", item.name, html)
    } else {
        fmt.Fprintf(os.Stderr, "[WARN] Could not render %s HTML coverage: %v\n", item.name, err)
    }
    if mode == "--strict" && covered != total {
        return fmt.Errorf("100%% coverage required; got %.2f%% (%d/%d)", percent, covered, total)
    }
    return nil
}

// score merges repeated coverage blocks from multiple go test binaries.
// A statement counts as covered if any selected unit-test suite executed it.
// Retaining duplicate zero-hit copies would make 100% unattainable even after
// every production statement was exercised in at least one suite.
func score(path string) (int64, int64, error) {
    file, err := os.Open(path)
    if err != nil { return 0, 0, err }
    scanner := bufio.NewScanner(file)
    if !scanner.Scan() {
        _ = file.Close()
        return 0, 0, fmt.Errorf("empty coverage profile: %s", path)
    }
    mode := scanner.Text()
    if !strings.HasPrefix(mode, "mode: ") {
        _ = file.Close()
        return 0, 0, fmt.Errorf("invalid coverage mode %q", mode)
    }
    type coverageBlock struct {
        statements int64
        hits int64
    }
    records := make(map[string]coverageBlock)
    order := make([]string, 0)
    for scanner.Scan() {
        columns := strings.Fields(scanner.Text())
        if len(columns) != 3 {
            _ = file.Close()
            return 0, 0, fmt.Errorf("invalid coverage record %q", scanner.Text())
        }
        statements, err := strconv.ParseInt(columns[1], 10, 64)
        if err != nil { _ = file.Close(); return 0, 0, err }
        hits, err := strconv.ParseInt(columns[2], 10, 64)
        if err != nil { _ = file.Close(); return 0, 0, err }
        if statements < 0 || hits < 0 {
            _ = file.Close()
            return 0, 0, fmt.Errorf("negative coverage values for %s", columns[0])
        }
        key := columns[0]
        block, found := records[key]
        if found && block.statements != statements {
            _ = file.Close()
            return 0, 0, fmt.Errorf("conflicting statement counts for %s", key)
        }
        if !found {
            order = append(order, key)
            block.statements = statements
        }
        block.hits += hits
        records[key] = block
    }
    if err := scanner.Err(); err != nil { _ = file.Close(); return 0, 0, err }
    if err := file.Close(); err != nil { return 0, 0, err }

    // Keep the first-seen order to retain Go cover's source-block ordering.
    // Rewrite the profile so HTML rendering and the strict gate use identical,
    // deduplicated instrumentation rather than different coverage denominators.
    var normalized strings.Builder
    normalized.WriteString(mode)
    normalized.WriteByte('\n')
    var total, covered int64
    for _, key := range order {
        block := records[key]
        total += block.statements
        if block.hits > 0 { covered += block.statements }
        fmt.Fprintf(&normalized, "%s %d %d\n", key, block.statements, block.hits)
    }
    temp, err := os.CreateTemp(filepath.Dir(path), "coverage-*.out")
    if err != nil { return 0, 0, err }
    defer os.Remove(temp.Name())
    if _, err := temp.WriteString(normalized.String()); err != nil { _ = temp.Close(); return 0, 0, err }
    if err := temp.Close(); err != nil { return 0, 0, err }
    if err := os.Rename(temp.Name(), path); err != nil { return 0, 0, err }
    return total, covered, nil
}
