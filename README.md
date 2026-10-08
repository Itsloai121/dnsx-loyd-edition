# dnsx-loyd-edition


# DNSX Audit Summary and Collaboration Proposal

Hello DNSX development team,

I completed an AI-assisted source-code and local-test review of the supplied `dnsx (loyd edition)` archive, using the methodology in `ai_security_audit.md`. I am sharing the results to help improve the tool’s security, reliability, and the confidence users can place in its scan results.

This was a review of a specific supplied snapshot—not a complete penetration test, a security certification, or proof that the current upstream release has all the same issues.

## What I did

I reviewed the application’s CLI and library entrypoints, input handling, filesystem operations, DNS transport configuration, test helpers, dependencies, and CI/release workflows.

The work included:

- Safely extracting and inventorying the archive, without running the software during extraction.
- Reviewing 54 non-Git project files, including Go source, tests, scripts, dependency manifests, and nine workflows.
- Recording the archive’s SHA-256 fingerprint and bundled repository commit to identify the reviewed snapshot.
- Tracing input to relevant execution, network, and filesystem operations, and checking whether the visible controls were effective.
- Inspecting selected pinned dependency source where security behavior was delegated outside the application.
- Building the software with checksum-verified Go 1.26.9.
- Running selected local tests, a local runner suite with race detection, `go vet`, and module-integrity verification.
- Checking dependency advisories with OSV and `govulncheck`, then assessing reachability and actual application context rather than treating every advisory match as an exploitable vulnerability.
- Running harmless local reproductions using loopback servers, temporary files, and stub programs.
- Challenging the findings to remove unsupported claims and distinguish confirmed behavior from conditional or unverified security impact.

I produced a 16-page PDF report and an evidence package containing structured findings, source locations, test logs, reproduction tests, and proposed remediation.

## What changed—and what did not

**I did not patch, upgrade, or deploy the software.** The work produced an audit report, supporting evidence, and proposed fixes. It did not produce a repaired release.

The original extracted application files were preserved. An incidental Git index refresh was detected during integrity checking and restored from the archive; all 82 extracted regular files then matched their recorded hashes.

The audit also changed how I assess the existing validation results: a successful build or passing test suite is not enough to establish trustworthy scan behavior. Some local reproductions passed precisely because they demonstrated a defect, while parts of the existing functional harness could accept incorrect results.

Any fixes should therefore be reviewed, implemented, and tested separately before a finding is considered closed.

## Problems identified

The report records **nine findings: six Medium and three Low**, spanning security defects, conditional risks, release configuration, and reliability. These are not nine demonstrated remote vulnerabilities.

| Finding | Problem | Important qualification |
| --- | --- | --- |
| F01 — Medium | DNS-over-HTTPS accepts untrusted TLS certificates through the pinned resolver dependency. | Reproduced locally; applies when DoH is configured and an attacker can impersonate or intercept the endpoint. |
| F02 — Low | The automatic version-check client disables TLS certificate verification. | Self-signed acceptance was reproduced. This does not establish malicious replacement-binary installation. |
| F03 — Medium | Test helpers construct shell commands from targets and arguments. | Command execution was reproduced with a harmless marker. Security impact requires less-trusted input reaching a more-privileged harness; this is not a production CLI remote-code-execution finding. |
| F04 — Medium | Stdin temporary-file handling creates, deletes, and reopens a filename, losing atomicity and restrictive permissions. | A replacement-symlink interleaving was reproduced under controlled timing. Cross-user exploitation and natural race success were not established. |
| F05 — Medium | Privileged CI jobs reference third-party actions through mutable tags. | The configuration risk is visible, but no action compromise or CI takeover was demonstrated. |
| F06 — Medium | Docker publication is not explicitly tied to a successful release and its corresponding immutable commit. | The missing checks are visible in the workflow; an actually mislabelled published image was not verified. |
| F07 — Low | Functional comparisons can accept different outputs if they contain the same number of lines; the supposed release baseline is built from the same checkout. | Incorrect same-length output was accepted in a local reproduction. |
| F08 — Medium | Missing input, zero workers, and discarded output errors can make incomplete work appear successful. | Actual CLI checks reproduced silent success for missing input and zero workers; an isolated output-handler test reproduced ignored flush failure. |
| F09 — Low | Stream mode without stdin can panic instead of reporting a controlled input error. | Reproduced locally; treated as a local availability defect, not remote service denial of service. |

Dependency screening also matched five Go advisories across three modules. I documented their reachability and context separately; I did not count them automatically as additional exploitable findings.

## Audit limitations and process issues

The initial archive tools could not handle the supplied RAR format, and the environment initially lacked a Go compiler and vulnerability scanner. Those obstacles were resolved before the build, tests, and scanning were performed.

However, important limits remain:

- I did not validate live deployment behavior, repository permission settings, branch protections, third-party action internals, or published image provenance.
- One existing runner test skipped because the ASN API key was not configured; a separate ASN expansion test was intentionally excluded.
- An attempted CLI output-failure check was inconclusive, so I used a separate output-handler reproduction rather than claiming the CLI attempt proved data loss.
- I did not perform comprehensive historical secret scanning, cross-platform testing, full fuzzing, or a multi-user race campaign.
- I did not authenticate the bundled Git history against upstream or audit an original deployed executable.
- Proposed fixes have not been applied or validated as a complete repaired release.

I did not establish a Critical or High finding within this scope. That should not be interpreted as proof that none exist.

## My offer to collaborate

I would like to work with the development team to validate these observations and help turn the confirmed problems into practical improvements.

I can share the report and evidence package, walk through the reproduction steps, and help coordinate a focused remediation and follow-up testing effort. A useful starting point would be:

1. Confirm which findings still apply to the current upstream code and intended deployment.
2. Restore certificate verification in the affected transport paths.
3. Make invalid input, zero-worker runs, and output failures return clear failure signals.
4. Replace shell-string construction in test helpers and preserve safe temporary-file creation.
5. Improve functional comparisons and build an independent release baseline.
6. Strengthen action pinning and bind release commits, version tags, and Docker image metadata.
7. Re-run the agreed regression tests on the patched release and document closure evidence.

I am open to discussing the collaboration format, responsibilities, priorities, and scope with the team. I would prefer to coordinate security-sensitive details privately before any public disclosure, and I welcome corrections or additional context that changes a finding’s applicability or severity.

If you are interested, please let me know the appropriate maintainer or security contact and your preferred process for reviewing the report and agreeing on next steps.

Thank you for your work on DNSX. My aim is to support the project constructively and help make its results safer and more dependable.
