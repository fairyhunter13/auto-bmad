<!-- Powered by BMAD-CORE™ -->

# CI/CD Pipeline Setup

**Workflow ID**: `_bmad/bmm/testarch/ci`
**Version**: 5.0 (Language Agnostic)

---

## Overview

Scaffolds a production-ready CI/CD quality pipeline with test execution, burn-in loops for flaky test detection, parallel sharding, artifact collection, and notification configuration. This workflow creates platform-specific CI configuration optimized for fast feedback and reliable test execution.

**Language Agnostic**: This workflow supports ANY programming language and test framework. Detection happens automatically in Step 0.

Note: This is typically a one-time setup per repo; run it any time after the test framework exists, ideally before feature work starts.

---

## Step 0: Detect Project Environment (MANDATORY)

**CRITICAL**: Before proceeding with CI setup, you MUST detect the project's language and test framework.

### Actions

1. **Detect Language from Codebase**

   Scan project root for language indicators:
   - `package.json` → JavaScript/TypeScript (Node.js)
   - `requirements.txt` / `pyproject.toml` / `setup.py` → Python
   - `go.mod` → Go
   - `pom.xml` / `build.gradle` → Java (Maven/Gradle)
   - `Cargo.toml` → Rust
   - `*.csproj` / `*.sln` → C# (.NET)
   - `Gemfile` → Ruby
   - `composer.json` → PHP

2. **Detect Test Framework**

   Based on language, identify test framework:
   - **JavaScript/TypeScript**: `playwright.config.*`, `vitest.config.*`, `jest.config.*`, `cypress.config.*`
   - **Python**: `pytest.ini`, `conftest.py`, `tox.ini`
   - **Go**: `*_test.go` files (native testing)
   - **Java**: JUnit config in `pom.xml`/`build.gradle`
   - **Rust**: `#[test]` in source files (native testing)
   - **C#**: `*.Tests.csproj`, xUnit/NUnit/MSTest references
   - **Ruby**: `spec/` directory, RSpec configuration
   - **PHP**: `phpunit.xml`

3. **Fetch Latest CI/CD Patterns (MANDATORY)**

   **CRITICAL**: Before generating ANY CI configuration, fetch current documentation:

   ```
   Search: "[DETECTED_LANGUAGE] [DETECTED_FRAMEWORK] CI/CD GitHub Actions best practices [CURRENT_YEAR]"
   Search: "[DETECTED_CI_PLATFORM] [DETECTED_FRAMEWORK] workflow configuration [CURRENT_YEAR]"
   ```

   Include `NOW()` timestamp in generated configuration comments:

   ```yaml
   # Generated: [TIMESTAMP]
   # Framework: [DETECTED_FRAMEWORK]
   # Docs fetched: [CURRENT_DATE]
   ```

4. **Store Detection Results**

   ```yaml
   detected:
     language: '[JavaScript|Python|Go|Java|Rust|C#|Ruby|PHP|Other]'
     runtime_version: '[e.g., Node 22, Python 3.12, Go 1.22]'
     test_framework: '[e.g., Playwright, pytest, go test, JUnit 6]'
     test_command: '[e.g., npm test, pytest, go test ./...]'
     package_manager: '[e.g., npm, pip, go mod, maven]'
     ci_platform: '[GitHub Actions|GitLab CI|CircleCI|Jenkins]'
   ```

**Halt Condition**: If language cannot be detected, ask user: "What programming language and test framework does this project use?"

---

## Preflight Requirements

**Critical:** Verify these requirements before proceeding. If any fail, HALT and notify the user.

- ✅ Git repository is initialized (`.git/` directory exists)
- ✅ Local test suite passes (run detected test command)
- ✅ Test framework is configured (from `framework` workflow)
- ✅ Team agrees on target CI platform (GitHub Actions, GitLab CI, Circle CI, etc.)
- ✅ Access to CI platform settings/secrets available (if updating existing pipeline)

---

## Step 1: Run Preflight Checks

### Actions

1. **Verify Git Repository**
   - Check for `.git/` directory
   - Confirm remote repository configured (`git remote -v`)
   - If not initialized, HALT with message: "Git repository required for CI/CD setup"

2. **Validate Test Framework** (Language-Agnostic)
   - Use detection results from Step 0
   - Read framework configuration to extract:
     - Test directory location
     - Test command (language-specific)
     - Reporter configuration
     - Timeout settings
   - If not found, HALT with message: "Run `framework` workflow first to set up test infrastructure"

3. **Run Local Tests** (Use Detected Command)
   - Execute the detected test command from Step 0:
     - **JavaScript/TypeScript**: `npm run test:e2e` or `yarn test`
     - **Python**: `pytest` or `python -m pytest`
     - **Go**: `go test ./...`
     - **Java**: `mvn test` or `gradle test`
     - **Rust**: `cargo test`
     - **C#**: `dotnet test`
     - **Ruby**: `bundle exec rspec`
     - **PHP**: `vendor/bin/phpunit`
   - Ensure tests pass before CI setup
   - If tests fail, HALT with message: "Fix failing tests before setting up CI/CD"

4. **Detect CI Platform**
   - Check for existing CI configuration:
     - `.github/workflows/*.yml` (GitHub Actions)
     - `.gitlab-ci.yml` (GitLab CI)
     - `.circleci/config.yml` (Circle CI)
     - `Jenkinsfile` (Jenkins)
     - `azure-pipelines.yml` (Azure DevOps)
   - If found, ask user: "Update existing CI configuration or create new?"
   - If not found, detect platform from git remote:
     - `github.com` → GitHub Actions (default)
     - `gitlab.com` → GitLab CI
     - `dev.azure.com` → Azure DevOps
     - Ask user if unable to auto-detect

5. **Read Environment Configuration** (Language-Specific)
   - **JavaScript/TypeScript**: Use `.nvmrc` for Node version, read `package.json`
   - **Python**: Use `.python-version` or `pyproject.toml` for Python version
   - **Go**: Read `go.mod` for Go version
   - **Java**: Read `pom.xml` or `build.gradle` for Java version
   - **Rust**: Read `rust-toolchain.toml` or use stable
   - **C#**: Read `.csproj` for target framework version
   - **Ruby**: Use `.ruby-version` for Ruby version
   - **PHP**: Read `composer.json` for PHP version requirements

**Halt Condition:** If preflight checks fail, stop immediately and report which requirement failed.

---

## Step 2: Scaffold CI Pipeline

### Actions

1. **Select CI Platform Template**

   Based on detection or user preference, use the appropriate template:

   **GitHub Actions** (`.github/workflows/test.yml`):
   - Most common platform
   - Excellent caching and matrix support
   - Free for public repos, generous free tier for private

   **GitLab CI** (`.gitlab-ci.yml`):
   - Integrated with GitLab
   - Built-in registry and runners
   - Powerful pipeline features

   **Circle CI** (`.circleci/config.yml`):
   - Fast execution with parallelism
   - Docker-first approach
   - Enterprise features

   **Jenkins** (`Jenkinsfile`):
   - Self-hosted option
   - Maximum customization
   - Requires infrastructure management

   **Azure DevOps** (`azure-pipelines.yml`):
   - Microsoft ecosystem integration
   - Built-in test reporting
   - Enterprise features

2. **Fetch Latest CI Patterns (MANDATORY)**

   **CRITICAL**: Before generating configuration, fetch current best practices:

   ```
   Search: "[DETECTED_LANGUAGE] [CI_PLATFORM] workflow best practices [CURRENT_YEAR]"
   Search: "[DETECTED_FRAMEWORK] CI configuration [CURRENT_YEAR]"
   ```

   Use official documentation for the detected language/framework combination.

3. **Generate Pipeline Configuration** (Language-Agnostic)

   **Key pipeline stages (universal pattern):**

   ```yaml
   stages:
     - lint # Code quality checks (language-specific linter)
     - test # Test execution (parallel shards)
     - burn-in # Flaky test detection (10 iterations)
     - report # Aggregate results and publish
   ```

4. **Configure Test Execution** (Use Detected Commands)

   **Parallel Sharding (GitHub Actions example - adapt for detected language):**

   ```yaml
   strategy:
     fail-fast: false
     matrix:
       shard: [1, 2, 3, 4]

   steps:
     - name: Run tests
       run: ${{ env.TEST_COMMAND }} ${{ env.SHARD_FLAG }}
       env:
         # Set based on detected framework:
         # JavaScript/Playwright: npm run test:e2e -- --shard=${{ matrix.shard }}/4
         # Python/pytest: pytest --splits=4 --group=${{ matrix.shard }}
         # Go: go test ./... -run "TestShard${{ matrix.shard }}"
         # Java/JUnit: mvn test -Dsurefire.shardIndex=${{ matrix.shard }}
   ```

   **Purpose:** Splits tests into N parallel jobs for faster execution (target: <10 min per shard)

5. **Add Burn-In Loop** (Language-Agnostic Pattern)

   **Critical pattern from production systems:**

   ```yaml
   burn-in:
     name: Flaky Test Detection
     runs-on: ubuntu-latest
     steps:
       - uses: actions/checkout@v4

       # Language-specific setup (use detected values)
       - name: Setup runtime
         # JavaScript: actions/setup-node@v4 with node-version-file
         # Python: actions/setup-python@v5 with python-version-file
         # Go: actions/setup-go@v5 with go-version-file
         # Java: actions/setup-java@v4 with java-version
         # Rust: dtolnay/rust-toolchain@stable
         # C#: actions/setup-dotnet@v4

       - name: Install dependencies
         run: ${{ env.INSTALL_COMMAND }}
         # JavaScript: npm ci
         # Python: pip install -r requirements.txt
         # Go: go mod download
         # Java: mvn dependency:resolve
         # Rust: cargo fetch

       - name: Run burn-in loop (10 iterations)
         run: |
           for i in {1..10}; do
             echo "🔥 Burn-in iteration $i/10"
             ${{ env.TEST_COMMAND }} || exit 1
           done

       - name: Upload failure artifacts
         if: failure()
         uses: actions/upload-artifact@v4
         with:
           name: burn-in-failures
           path: ${{ env.TEST_RESULTS_PATH }}
           retention-days: 30
   ```

   **Purpose:** Runs tests multiple times to catch non-deterministic failures before they reach main branch.

   **When to run:**
   - On pull requests to main/develop
   - Weekly on cron schedule
   - After significant test infrastructure changes

6. **Configure Caching** (Language-Specific)

   **Fetch latest caching patterns for detected language:**

   ```
   Search: "[DETECTED_LANGUAGE] GitHub Actions cache [CURRENT_YEAR]"
   ```

   **Language-specific cache configurations:**

   | Language      | Cache Path                 | Hash File           |
   | ------------- | -------------------------- | ------------------- |
   | JavaScript    | `~/.npm` or `node_modules` | `package-lock.json` |
   | Python        | `~/.cache/pip`             | `requirements.txt`  |
   | Go            | `~/go/pkg/mod`             | `go.sum`            |
   | Java (Maven)  | `~/.m2/repository`         | `pom.xml`           |
   | Java (Gradle) | `~/.gradle/caches`         | `build.gradle`      |
   | Rust          | `~/.cargo`                 | `Cargo.lock`        |
   | C#            | `~/.nuget/packages`        | `*.csproj`          |
   | Ruby          | `vendor/bundle`            | `Gemfile.lock`      |
   | PHP           | `vendor`                   | `composer.lock`     |

   **Browser cache (for E2E frameworks):**
   - Playwright: `~/.cache/ms-playwright`
   - Cypress: `~/.cache/Cypress`
   - Selenium: Driver-specific paths

   **Purpose:** Reduces CI execution time by 2-5 minutes per run.

7. **Configure Artifact Collection** (Language-Agnostic)

   **Failure artifacts only:**

   ```yaml
   - name: Upload test results
     if: failure()
     uses: actions/upload-artifact@v4
     with:
       name: test-results-${{ matrix.shard }}
       path: ${{ env.TEST_RESULTS_PATH }}
       # Language-specific paths (set in env):
       # JavaScript: test-results/, playwright-report/, coverage/
       # Python: .pytest_cache/, htmlcov/, test-results/
       # Go: coverage.out, test-results/
       # Java: target/surefire-reports/, target/site/jacoco/
       # Rust: target/debug/, test-results/
       retention-days: 30
   ```

   **Artifacts to collect (by test type):**
   - **E2E tests**: Traces, screenshots, videos, HAR files
   - **Unit tests**: Coverage reports, failure logs
   - **API tests**: Request/response logs, HAR recordings
   - **All tests**: HTML reports, console logs, error messages

8. **Add Retry Logic** (Language-Agnostic)

   ```yaml
   - name: Run tests with retries
     uses: nick-invision/retry@v2
     with:
       timeout_minutes: 30
       max_attempts: 3
       retry_on: error
       command: ${{ env.TEST_COMMAND }}
   ```

   **Purpose:** Handles transient failures (network issues, race conditions)

9. **Configure Notifications** (Optional)

   If `notify_on_failure` is enabled:

   ```yaml
   - name: Notify on failure
     if: failure()
     uses: 8398a7/action-slack@v3
     with:
       status: ${{ job.status }}
       text: 'Test failures detected in PR #${{ github.event.pull_request.number }}'
       webhook_url: ${{ secrets.SLACK_WEBHOOK }}
   ```

10. **Generate Helper Scripts** (Language-Agnostic)

    **CRITICAL**: Fetch latest script patterns for detected language before generating:

    ```
    Search: "[DETECTED_LANGUAGE] test script best practices [CURRENT_YEAR]"
    ```

    **Selective testing script** (`scripts/test-changed.sh`):

    ```bash
    #!/bin/bash
    # Run only tests for changed files
    # IMPORTANT: Adapt FILE_PATTERN and TEST_COMMAND for detected language

    CHANGED_FILES=$(git diff --name-only HEAD~1)

    # Language-specific patterns:
    # JavaScript/TypeScript: src/.*\.(ts|js)$
    # Python: .*\.py$
    # Go: .*\.go$
    # Java: .*\.java$
    # Rust: .*\.rs$
    FILE_PATTERN="${DETECTED_FILE_PATTERN}"

    if echo "$CHANGED_FILES" | grep -qE "$FILE_PATTERN"; then
      echo "Running affected tests..."
      ${TEST_COMMAND_WITH_FILTER}
    else
      echo "No test-affecting changes detected"
    fi
    ```

    **Local mirror script** (`scripts/ci-local.sh`):

    ```bash
    #!/bin/bash
    # Mirror CI execution locally for debugging
    # IMPORTANT: Replace commands with detected language equivalents

    echo "🔍 Running CI pipeline locally..."

    # Lint (language-specific: $LINT_COMMAND)
    ${LINT_COMMAND} || exit 1

    # Tests (language-specific: $TEST_COMMAND)
    ${TEST_COMMAND} || exit 1

    # Burn-in (reduced iterations)
    for i in {1..3}; do
      echo "🔥 Burn-in $i/3"
      ${TEST_COMMAND} || exit 1
    done

    echo "✅ Local CI pipeline passed"
    ```

    **Language-specific command reference:**

    | Language   | Lint Command                        | Test Command         |
    | ---------- | ----------------------------------- | -------------------- |
    | JavaScript | `npm run lint`                      | `npm run test:e2e`   |
    | Python     | `ruff check .` or `flake8`          | `pytest`             |
    | Go         | `golangci-lint run`                 | `go test ./...`      |
    | Java       | `mvn checkstyle:check`              | `mvn test`           |
    | Rust       | `cargo clippy`                      | `cargo test`         |
    | C#         | `dotnet format --verify-no-changes` | `dotnet test`        |
    | Ruby       | `rubocop`                           | `bundle exec rspec`  |
    | PHP        | `vendor/bin/phpcs`                  | `vendor/bin/phpunit` |

11. **Generate Documentation**

    **CI README** (`docs/ci.md`):
    - Pipeline stages and purpose
    - How to run locally
    - Debugging failed CI runs
    - Secrets and environment variables needed
    - Notification setup
    - Badge URLs for README

    **Secrets checklist** (`docs/ci-secrets-checklist.md`):
    - Required secrets list (SLACK_WEBHOOK, etc.)
    - Where to configure in CI platform
    - Security best practices

---

## Step 3: Deliverables

### Primary Artifacts Created

1. **CI Configuration File** (Based on detected CI platform)
   - `.github/workflows/test.yml` (GitHub Actions)
   - `.gitlab-ci.yml` (GitLab CI)
   - `.circleci/config.yml` (Circle CI)
   - `Jenkinsfile` (Jenkins)
   - `azure-pipelines.yml` (Azure DevOps)

2. **Pipeline Stages** (Language-agnostic pattern)
   - **Lint**: Code quality checks (language-specific linter)
   - **Test**: Parallel test execution (4 shards)
   - **Burn-in**: Flaky test detection (10 iterations)
   - **Report**: Result aggregation and publishing

3. **Helper Scripts** (Adapted for detected language)
   - `scripts/test-changed.sh` - Selective testing
   - `scripts/ci-local.sh` - Local CI mirror
   - `scripts/burn-in.sh` - Standalone burn-in execution

4. **Documentation**
   - `docs/ci.md` - CI pipeline guide
   - `docs/ci-secrets-checklist.md` - Required secrets
   - Inline comments in CI configuration

5. **Optimization Features** (Language-Agnostic)
   - Dependency caching (language-specific package manager)
   - Browser binary caching (for E2E frameworks)
   - Parallel sharding (4 jobs default)
   - Retry logic (2 retries on failure)
   - Failure-only artifact upload

### Performance Targets

- **Lint stage**: <2 minutes
- **Test stage** (per shard): <10 minutes
- **Burn-in stage**: <30 minutes (10 iterations)
- **Total pipeline**: <45 minutes

**Speedup:** 20× faster than sequential execution through parallelism and caching.

---

## Important Notes

### Language-Agnostic CI/CD Approach

**CRITICAL**: This workflow generates CI configuration for ANY language. Always:

1. Detect language/framework in Step 0
2. Fetch latest CI patterns via web search before generating config
3. Use detected commands throughout (don't hardcode npm/pytest/etc.)
4. Include generation timestamp in config comments

### Knowledge Base Integration

**Critical:** Check configuration and load appropriate fragments.

Read `{config_source}` and check `config.tea_use_playwright_utils`.

**Core CI Patterns (Always load - language-agnostic principles):**

- `ci-burn-in.md` - Burn-in loop patterns: 10-iteration detection, CI workflow, shard orchestration, selective execution (678 lines, 4 examples)
- `selective-testing.md` - Changed test detection strategies: tag-based, spec filters, diff-based selection, promotion rules (727 lines, 4 examples)
- `visual-debugging.md` - Artifact collection best practices: trace viewer, HAR recording, custom artifacts (522 lines, 5 examples)
- `test-quality.md` - CI-specific test quality criteria: deterministic tests, isolated with cleanup, explicit assertions (658 lines, 5 examples)

**If `config.tea_use_playwright_utils: true` (JavaScript/TypeScript projects):**

Load Playwright-specific CI-relevant fragments:

- `playwright-config.md` - CI-optimized configuration: parallelization, artifact output, sharding
- `burn-in.md` - Smart test selection with git diff analysis
- `network-error-monitor.md` - Automatic HTTP 4xx/5xx detection

**For other languages, fetch framework-specific CI patterns:**

```
Search: "[DETECTED_FRAMEWORK] CI/CD configuration best practices [CURRENT_YEAR]"
```

### CI Platform-Specific Guidance

**GitHub Actions:**

- Use `actions/cache` for caching
- Matrix strategy for parallelism
- Secrets in repository settings
- Free 2000 minutes/month for private repos

**GitLab CI:**

- Use `.gitlab-ci.yml` in root
- `cache:` directive for caching
- Parallel execution with `parallel: 4`
- Variables in project CI/CD settings

**Circle CI:**

- Use `.circleci/config.yml`
- Docker executors recommended
- Parallelism with `parallelism: 4`
- Context for shared secrets

### Burn-In Loop Strategy

**When to run:**

- ✅ On PRs to main/develop branches
- ✅ Weekly on schedule (cron)
- ✅ After test infrastructure changes
- ❌ Not on every commit (too slow)

**Iterations:**

- **10 iterations** for thorough detection
- **3 iterations** for quick feedback
- **100 iterations** for high-confidence stability

**Failure threshold:**

- Even ONE failure in burn-in → tests are flaky
- Must fix before merging

### Artifact Retention

**Failure artifacts only:**

- Saves storage costs
- Maintains debugging capability
- 30-day retention default

**Artifact types by framework:**

| Framework  | Artifacts                   | Typical Size     |
| ---------- | --------------------------- | ---------------- |
| Playwright | Traces, screenshots, videos | 5-10 MB per test |
| Cypress    | Screenshots, videos, HAR    | 2-5 MB per test  |
| pytest     | HTML reports, coverage      | 1-2 MB per run   |
| JUnit      | XML reports, coverage       | 1-5 MB per run   |
| Go test    | Coverage profiles, JSON     | <1 MB per run    |
| RSpec      | HTML reports, screenshots   | 1-3 MB per run   |

### Selective Testing

**Detect changed files:**

```bash
git diff --name-only HEAD~1
```

**Run affected tests only:**

- Faster feedback for small changes
- Full suite still runs on main branch
- Reduces CI time by 50-80% for focused PRs

**Trade-off:**

- May miss integration issues
- Run full suite at least on merge

### Local CI Mirror

**Purpose:** Debug CI failures locally

**Usage:**

```bash
./scripts/ci-local.sh
```

**Mirrors CI environment:**

- Same runtime version (Node/Python/Go/Java/etc.)
- Same test command (detected in Step 0)
- Same stages (lint → test → burn-in)
- Reduced burn-in iterations (3 vs 10)

---

## Output Summary

After completing this workflow, provide a summary:

```markdown
## CI/CD Pipeline Complete

**Platform**: {DETECTED_CI_PLATFORM} (GitHub Actions, GitLab CI, etc.)
**Language**: {DETECTED_LANGUAGE}
**Framework**: {DETECTED_TEST_FRAMEWORK}
**Generated**: {TIMESTAMP}

**Artifacts Created**:

- ✅ Pipeline configuration: {CI_CONFIG_PATH}
- ✅ Burn-in loop: 10 iterations for flaky detection
- ✅ Parallel sharding: 4 jobs for fast execution
- ✅ Caching: Dependencies ({PACKAGE_MANAGER}) + test artifacts
- ✅ Artifact collection: Failure-only reports/screenshots/logs
- ✅ Helper scripts: test-changed.sh, ci-local.sh, burn-in.sh
- ✅ Documentation: docs/ci.md, docs/ci-secrets-checklist.md

**Performance:**

- Lint: <2 min
- Test (per shard): <10 min
- Burn-in: <30 min
- Total: <45 min (20× speedup vs sequential)

**Next Steps**:

1. Commit CI configuration: `git add {CI_CONFIG_PATH} && git commit -m "ci: add test pipeline"`
2. Push to remote: `git push`
3. Configure required secrets in CI platform settings (see docs/ci-secrets-checklist.md)
4. Open a PR to trigger first CI run
5. Monitor pipeline execution and adjust parallelism if needed

**Documentation Fetched**:

- {FRAMEWORK} CI best practices ({CURRENT_YEAR})
- {CI_PLATFORM} workflow configuration

**Knowledge Base References Applied**:

- Burn-in loop pattern (ci-burn-in.md)
- Selective testing strategy (selective-testing.md)
- Artifact collection (visual-debugging.md)
- Test quality criteria (test-quality.md)
```

---

## Validation

After completing all steps, verify:

- [ ] CI configuration file created and syntactically valid
- [ ] Burn-in loop configured (10 iterations)
- [ ] Parallel sharding enabled (4 jobs)
- [ ] Caching configured (dependencies + browsers)
- [ ] Artifact collection on failure only
- [ ] Helper scripts created and executable (`chmod +x`)
- [ ] Documentation complete (ci.md, secrets checklist)
- [ ] No errors or warnings during scaffold

Refer to `checklist.md` for comprehensive validation criteria.
