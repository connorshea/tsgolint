# Performance Analysis

Analysis of the tsgolint codebase identifying performance improvement opportunities, ranked by estimated impact.

## High Impact

### 1. Repeated JSON Marshal/Unmarshal for Rule Options

**Location:** `internal/utils/utils.go:190-203`, called from ~30 rules' `Run` functions

**Issue:** `UnmarshalOptions[T]` performs `json.Marshal` → `json.Unmarshal` on every call. Each rule's `Run` function is called once per file per worker. With 30 rules using options × 1000 files = 30,000 redundant marshal/unmarshal cycles.

**Fix:** Cache the deserialized options. In headless mode, each file's rule config comes from `fileConfigs[fileName]` which maps to `[]headlessRule` — the options don't change per file. The `Run` function could be restructured to parse options once and reuse them, or a caching layer could be added to `UnmarshalOptions`.

### 2. Sequential Program Processing

**Location:** `internal/linter/linter.go:58-121`

**Issue:** When multiple tsconfig programs exist, they are processed strictly sequentially. Program creation (`utils.CreateProgram`) blocks before linting begins. With large monorepos containing many tsconfig files, this serialization is a bottleneck.

**Fix:** Overlap program creation with linting — start creating the next program while the current one is being linted. This is a larger architectural change but would significantly help monorepo scenarios.

### 3. Sequential Type Error Checking

**Location:** `internal/linter/linter.go:259-291`

**Issue:** When `typeErrors.ReportSyntactic` or `typeErrors.ReportSemantic` is enabled, diagnostics are collected sequentially for each file before the parallel worker pool starts. `GetSemanticDiagnostics` can be expensive.

**Fix:** Move type error collection into the parallel worker pool, or process it concurrently alongside rule execution.

## Medium Impact

### 4. Per-File Closure Re-creation in Linter

**Location:** `internal/linter/linter.go:360-433`

**Issue:** `runListeners`, `childVisitor`, and `patternVisitor` closures are re-created for every file processed by each worker. While Go closures are lightweight, this is in the hottest loop of the application and creates garbage for the GC to collect.

**Fix:** Move closure creation outside the per-file loop. `runListeners` only captures `registeredListeners` and `ctxBuilder` which are per-worker (not per-file). The `childVisitor`/`patternVisitor` can similarly be hoisted.

### 5. `getRulesForFile` Allocates New Slice Per File (Standalone Mode)

**Location:** `cmd/tsgolint/main.go:524-532`

**Issue:** In standalone mode, `getRulesForFile` calls `utils.Map(allRules, ...)` which allocates a new `[]ConfiguredRule` slice for every source file. Since standalone mode uses the same rules for every file, this slice could be computed once.

**Fix:** Pre-compute the `[]ConfiguredRule` slice once outside the callback and return a reference to it.

### 6. Heap Escape in `ruleToAny` / `internalToAny`

**Location:** `cmd/tsgolint/headless.go:151-157`

**Issue:** `ruleToAny` takes `RuleDiagnostic` by value then takes its address (`&d`), forcing the value to escape to the heap. Same for `internalToAny`. Every diagnostic emitted causes a heap allocation.

**Fix:** Accept pointer parameters or restructure `anyDiagnostic` to embed the diagnostic value directly rather than using pointers.

### 7. `NodeFactory` Allocation in `GetCommentsInRange`

**Location:** `internal/utils/utils.go:20`

**Issue:** `ast.NewNodeFactory(ast.NodeFactoryHooks{})` is called every time `GetCommentsInRange` is invoked. If this factory is stateless, it could be created once as a package-level variable.

**Fix:** Create a package-level `nodeFactory` variable initialized once.

## Low Impact

### 8. `Flatten` Doesn't Pre-allocate

**Location:** `internal/utils/utils.go:140-146`

**Issue:** `Flatten` appends to a nil slice without pre-computing total capacity, causing multiple re-allocations for large inputs.

**Fix:** Pre-compute total length and allocate once:
```go
func Flatten[T any](array [][]T) []T {
    n := 0
    for _, sub := range array {
        n += len(sub)
    }
    result := make([]T, 0, n)
    for _, sub := range array {
        result = append(result, sub...)
    }
    return result
}
```

### 9. `lineStarts`/`lineEnds` Allocated Per Diagnostic in `printDiagnostic`

**Location:** `cmd/tsgolint/main.go:285-286`

**Issue:** Two 13-element `[]int` slices are allocated on every `printDiagnostic` call. These are small and likely stack-allocated by the compiler, but could be moved to the caller for reuse if profiling shows they escape.

### 10. `allRulesByName` Rebuilt in Headless Mode

**Location:** `cmd/tsgolint/headless.go:277-280`

**Issue:** `allRulesByName` is rebuilt from `allRules` at the start of every headless run, despite `main.go:225-231` already maintaining a package-level `allRulesByName`. The headless function creates its own local copy.

**Fix:** Reuse the existing package-level `allRulesByName` map.

## Already Well-Optimized

The codebase already has several good performance patterns:

- **Worker pool with GOMAXPROCS workers** — good CPU utilization
- **Files sorted by length descending** — helps load balancing
- **Listener map slice reuse** — `registeredListeners[k][:0]` avoids reallocation
- **Per-worker `ruleContextBuilder`** — avoids closure allocation per rule
- **Buffered diagnostic channel (4096)** — prevents blocking workers
- **Large I/O buffer (409KB)** — efficient output writing
- **Streaming diagnostics** — no batching/sorting overhead

## Recommendations

1. **Profile first**: Run `tsgolint -cpuprof cpu.prof` on a real project and use `go tool pprof` to confirm which hotspots matter most.
2. **Start with #1** (options caching) — high impact, low risk, easy to benchmark.
3. **Consider #2** (parallel programs) for monorepo use cases — high impact but requires careful design.
4. **Batch the small wins** (#4, #5, #7, #8, #10) — individually minor but collectively meaningful.
