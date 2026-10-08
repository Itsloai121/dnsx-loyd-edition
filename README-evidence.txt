DNSX AUDIT EVIDENCE — 08 OCTOBER 2026

Read the PDF for scope, severity, limitations, and remediation order.
findings.json contains nine structured findings; source quotes are exact.
proposed-fixes.txt contains targeted proposed diffs/snippets, NOT an
apply-ready or closure-tested patch set. External action SHAs are deliberately
not fabricated. No repairs were made to original software.

Commands/checks performed
-------------------------
Archive inspection: file; 7z l (missing); 7za l (RAR v5 unsupported).
Safe extraction: local libarchive through extract.py; regular files only,
path/size limits, per-file SHA-256 manifest. 116 entries, 82 regular files.
Source inventory and sink searches: find, rg, nl/sed and read of supplied code.
Git: log -1, diff --ignore-space-at-eol, status --short. Bundled HEAD was not
externally authenticated. Git index refresh was detected and restored.
Secret screening: working-tree heuristic patterns, not an exhaustive/history
secret scanner; five candidates were GitHub secret references, not secrets.

Compiler: official Go 1.26.9 archive checksum verified against go.dev JSON.
Environment: GOPATH set to audit-local cache; GOTOOLCHAIN=local.
go mod download
go build -o <audit>/dnsx ./cmd/dnsx
go test ./internal/runner -run <selected local test names> -count=1 -timeout=120s -json
go test -race ./internal/runner -run <all runner Test names except
  TestRunner_asnInput_prepareInput> -count=1 -timeout=120s -json
go vet ./...
go mod verify

Vulnerability screening:
OSV querybatch on 105 declared module requirements; five advisories matched.
go install golang.org/x/vuln/cmd/govulncheck@latest (resolved v1.8.0)
govulncheck -json ./...
JSON mode returned exit 0 despite findings; consult records, not exit alone.
Scanner build version and DB update timestamp are in scanner-version.txt.

Benign reproductions in a separate audit copy:
go test ./internal/testutils ./cmd/functional-test ./libs/dnsx
  ./internal/runner -run TestAudit -count=1 -timeout=120s -json
go test ./internal/runner -run TestAuditOutputErrorDiscarded -count=1 -json

Place the included audit_security_test.go files at their corresponding paths
in a disposable COPY of the supplied source to repeat those reproductions.
They are proof-of-defect tests: passing means the defect/behavior happened.
After remediation, invert the expectations and use negative regression tests.
The zero-thread unit test is a control-flow check; actual CLI behavior is also
recorded in cli-checks.json. Temp-race interleaving was controlled, not timed
against a different user or privileged process. Tests touch only local stubs,
loopback TLS/DNS listeners, disposable files and /dev/full.

CLI checks with -duc -silent -nc:
  no input, stdin /dev/null -> exit 0, no output (undesired).
  -l example.test -r 127.0.0.1:9 -t 0 -> exit 0, no output (undesired).
  -stream, stdin /dev/null -> panic, exit 2.
  -l localhost -hostsfile -o /dev/full -> no response obtained; INCONCLUSIVE
  as an output-loss proof. A separate HandleOutput test supplied a record
  and proved its flush error is discarded.

Package exclusions and interpretation
-------------------------------------
Original executable, full Git history, toolchain and caches are excluded.
No production credentials, live exploitation or deployment access was used.
Raw scanner/test records can include workspace-local paths and ANSI codes.
No malware-free verdict, complete penetration test, or production safety
certification is implied. Full historical secrets/cross-platform/fuzzing and
live CI settings/action internals remain unverified.
