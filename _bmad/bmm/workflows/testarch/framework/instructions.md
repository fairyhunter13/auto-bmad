<!-- Powered by BMAD-CORE™ -->

# Test Framework Setup

**Workflow ID**: `_bmad/bmm/testarch/framework`
**Version**: 5.0 (BMad v6 - Language Agnostic)

---

## Overview

Initialize a production-ready test framework architecture with fixtures, helpers, configuration, and best practices. This workflow scaffolds the complete testing infrastructure for ANY programming language and test framework.

**Language Agnostic**: This workflow detects the project's language and framework, then generates appropriate scaffolding using real-time documentation fetching.

---

## Preflight Requirements

**Critical:** Verify these requirements before proceeding. If any fail, HALT and notify the user.

- ✅ Project has identifiable language (package manifest, source files, or config)
- ✅ No modern test harness is already configured for the detected framework
- ✅ Architectural/stack context available (project type, bundler, dependencies)

---

## Step 0: Detect Project Environment (MANDATORY)

### Purpose

Identify the project's programming language, existing test framework (if any), and conventions before scaffolding. This enables language-agnostic test generation.

### Actions

1. **Detect Programming Language**

   Scan project root for language indicators:
   - Check for package/manifest files (see `detection/language-hints.yaml`)
   - Identify primary language from file extensions
   - Note if multiple languages exist (e.g., TypeScript frontend + Python backend)

   **Common Detection Patterns:**
   | Files Found | Language |
   |-------------|----------|
   | package.json, tsconfig.json | TypeScript |
   | package.json (no tsconfig) | JavaScript |
   | requirements.txt, pyproject.toml | Python |
   | go.mod | Go |
   | Cargo.toml | Rust |
   | pom.xml, build.gradle | Java/Kotlin |
   | \*.csproj | C#/.NET |
   | Gemfile | Ruby |
   | composer.json | PHP |

2. **Check for Existing Test Framework**

   Search for test configuration files:
   - `playwright.config.*`, `cypress.config.*` → E2E frameworks
   - `jest.config.*`, `vitest.config.*` → Unit test frameworks
   - `pytest.ini`, `conftest.py` → Python testing
   - `*_test.go` files → Go testing
   - Other framework-specific configs (see `language-hints.yaml`)

   **If found:** HALT with message: "Existing test framework detected: [framework]. Use workflow `test-review` to audit or manually extend."

3. **Detect Project Conventions**

   If existing test files found:
   - Note directory structure (tests/, **tests**/, spec/, etc.)
   - Identify naming patterns (_.spec._, _.test._, _\_test._)
   - Check import/module style
   - Note assertion library preferences

4. **Web Fetch Latest Documentation (MANDATORY - Use NOW())**

   Before generating ANY scaffolding:

   ```
   ============================================================
   CRITICAL: Fetch with CURRENT TIMESTAMP - NOW()
   ============================================================

   1. Record fetch timestamp: fetch_time = NOW()
      - Use ISO 8601 format: YYYY-MM-DDTHH:MM:SSZ
      - Example: 2026-01-23T14:30:00Z

   2. Include timestamp in ALL web searches:
      - "[framework] official documentation January 2026"
      - "[framework] best practices 2026"
      - "[framework] latest version changelog"

   3. Fetch for detected/selected framework:
      - Installation guide (current version)
      - Configuration reference (current API)
      - Best practices / getting started
      - Breaking changes / migration guides

   4. Verify freshness:
      - Check documentation date/version
      - Look for "last updated" indicators
      - Prefer official sources over cached/archived pages

   5. If fetch fails:
      - Retry with alternative search terms
      - Try GitHub repository directly
      - Mark output as "[UNVERIFIED]" if all fetches fail

   WHY NOW() IS REQUIRED:
   - Framework APIs change frequently
   - Best practices evolve monthly
   - Yesterday's patterns may be deprecated today
   - Each workflow run must get fresh data
   ```

5. **Record Detection Results**

   Store for use in subsequent steps (include fetch timestamp):

   ```yaml
   detected:
     detection_timestamp: '{NOW() in ISO 8601}' # REQUIRED
     language: '{language}'
     language_version: '{version if detectable}'
     existing_framework: '{framework or none}'
     test_directory: '{existing test dir or default}'
     naming_convention: '{pattern}'
     docs_fetched:
       - url: '{official docs url}'
         fetched_at: '{NOW() timestamp}' # REQUIRED
         verified: true/false
   ```

   REQUIRED: Fetch current documentation for detected/selected framework
   1. Identify framework (from detection or user selection)
   2. Web fetch official documentation:
      - Installation guide
      - Configuration reference
      - Best practices / getting started
   3. Check for latest version and breaking changes
   4. Note any deprecation warnings

   Example searches:
   - "[framework] official documentation [year]"
   - "[framework] configuration reference"
   - "[framework] v[version] migration guide"

   ```

   ```

6. **Record Detection Results**

   Store for use in subsequent steps:

   ```yaml
   detected:
     language: '{language}'
     language_version: '{version if detectable}'
     existing_framework: '{framework or none}'
     test_directory: '{existing test dir or default}'
     naming_convention: '{pattern}'
     docs_fetched:
       - url: '{official docs url}'
         fetched_at: '{timestamp}'
   ```

**Halt Condition:** If language cannot be detected and user doesn't provide guidance, HALT with message: "Unable to detect project language. Please specify the primary language and preferred test framework."

---

## Step 1: Run Preflight Checks

### Actions

1. **Validate package.json**
   - Read `{project-root}/package.json`
   - Extract project type (React, Vue, Angular, Next.js, Node, etc.)
   - Identify bundler (Vite, Webpack, Rollup, esbuild)
   - Note existing test dependencies

2. **Check for Existing Framework**
   - Search for `playwright.config.*`, `cypress.config.*`, `cypress.json`
   - Check `package.json` for `@playwright/test` or `cypress` dependencies
   - If found, HALT with message: "Existing test framework detected. Use workflow `upgrade-framework` instead."

3. **Gather Context**
   - Look for architecture documents (`architecture.md`, `tech-spec*.md`)
   - Check for API documentation or endpoint lists
   - Identify authentication requirements

**Halt Condition:** If preflight checks fail, stop immediately and report which requirement failed.

---

## Step 2: Scaffold Framework

### Actions

1. **Framework Selection (Language-Aware)**

   Based on detected language from Step 0, recommend appropriate test frameworks:

   **Framework Selection by Language:**

   | Language      | E2E/Integration      | Unit/Component        | Recommended Default |
   | ------------- | -------------------- | --------------------- | ------------------- |
   | TypeScript/JS | Playwright, Cypress  | Jest, Vitest          | Playwright + Vitest |
   | Python        | Playwright, Selenium | pytest                | pytest + Playwright |
   | Java          | Selenium, Playwright | JUnit 5, TestNG       | JUnit 5 + Selenium  |
   | C#            | Playwright, Selenium | xUnit, NUnit          | xUnit + Playwright  |
   | Go            | Chromedp, Rod        | testing (built-in)    | testing + testify   |
   | Ruby          | Capybara, Selenium   | RSpec, Minitest       | RSpec + Capybara    |
   | Rust          | -                    | cargo test (built-in) | cargo test          |
   | PHP           | Codeception, Panther | PHPUnit, Pest         | PHPUnit/Pest        |

   **Selection Strategy:**
   1. Check detected language from Step 0
   2. Check for existing framework preferences in dependencies
   3. Consider project type (web app, API, CLI, library)
   4. Ask user if multiple valid options exist
   5. **CRITICAL**: Web fetch latest docs for selected framework before generating config

2. **Create Directory Structure**

   Generate structure appropriate for detected language:

   **Universal Pattern (adapt naming to language conventions):**

   ```
   {project-root}/
   ├── tests/                        # Root test directory (or language convention)
   │   ├── e2e/                      # End-to-end tests
   │   ├── integration/              # Integration tests
   │   ├── unit/                     # Unit tests
   │   ├── support/                  # Framework infrastructure
   │   │   ├── fixtures/             # Test fixtures (data, mocks)
   │   │   ├── helpers/              # Utility functions
   │   │   └── factories/            # Data factories
   │   └── README.md                 # Test suite documentation
   ```

   **Language-Specific Conventions:**
   - **Python**: `tests/`, `conftest.py` for fixtures
   - **Go**: `*_test.go` files alongside source, `testdata/` directory
   - **Java**: `src/test/java/`, `src/test/resources/`
   - **Ruby**: `spec/` for RSpec, `test/` for Minitest
   - **Rust**: `tests/` for integration, inline `#[cfg(test)]` for unit

3. **Generate Configuration File (Framework-Specific)**

   **CRITICAL**: Use freshly fetched documentation to generate config.
   Do NOT rely on memorized examples - APIs and best practices change.

   **Generation Process:**

   ```
   1. Web fetch: "[framework] configuration reference [year]"
   2. Web fetch: "[framework] recommended settings"
   3. Check for version-specific configuration options
   4. Generate config using CURRENT API from fetched docs
   5. Include comments linking to official documentation
   ```

   **Universal Configuration Principles (all frameworks):**
   - Timeout settings: action ~15s, test ~60s
   - Parallel execution enabled (where supported)
   - Retry on failure in CI (typically 2 retries)
   - Failure artifacts: screenshots, traces, logs
   - Reporter configuration: console + CI-compatible (JUnit XML)
   - Environment-based configuration (local vs CI)

   **Example Structure (generate from fetched docs):**

   ```
   // Generated: {timestamp}
   // Framework: {framework} v{version}
   // Docs: {official_docs_url}

   {config_content_from_fetched_docs}
   ```

   **For TypeScript/JavaScript (Playwright example - fetch latest):**
   Web fetch <https://playwright.dev/docs/test-configuration> then generate.

   **For Python (pytest example - fetch latest):**
   Web fetch <https://docs.pytest.org/en/stable/reference/customize.html> then generate.

   **For Java (JUnit 5 example - fetch latest):**
   Web fetch <https://junit.org/junit5/docs/current/user-guide/> then generate.

   **For other languages:**
   Web search "[framework] configuration guide" and generate from official docs.

4. **Generate Environment Configuration**

   Create environment config appropriate for the language:

   **For all languages - `.env.example`:**

   ```bash
   # Test Environment Configuration
   TEST_ENV=local
   BASE_URL=http://localhost:3000
   API_URL=http://localhost:3001/api

   # Authentication (if applicable)
   TEST_USER_EMAIL=test@example.com
   TEST_USER_PASSWORD=

   # Feature Flags (if applicable)
   FEATURE_FLAG_NEW_UI=true

   # API Keys (if applicable)
   TEST_API_KEY=
   ```

   **Language-specific version files:**
   - **Node.js**: `.nvmrc` with LTS version
   - **Python**: `.python-version` or `pyproject.toml` python version
   - **Ruby**: `.ruby-version`
   - **Go**: `go.mod` already specifies version
   - **Java**: `.java-version` or build tool config

5. **Implement Fixture Architecture (Language-Agnostic)**

   **Knowledge Base Reference**: `testarch/knowledge/fixture-architecture.md`

   **Universal Fixture Principles:**
   - Fixtures provide reusable setup/teardown
   - Auto-cleanup: resources created during setup are deleted in teardown
   - Composable: combine multiple fixtures without inheritance
   - Type-safe: leverage language's type system where available

   **CRITICAL**: Web fetch fixture patterns for detected framework:

   ```
   1. Web fetch: "[framework] fixtures guide"
   2. Web fetch: "[framework] test setup teardown"
   3. Generate using CURRENT patterns from fetched docs
   ```

   **Generate fixtures following language conventions:**

   | Language     | Fixture Pattern              | Example Location          |
   | ------------ | ---------------------------- | ------------------------- |
   | TypeScript   | `test.extend<T>()`           | `tests/support/fixtures/` |
   | Python       | `@pytest.fixture`            | `conftest.py`             |
   | Java         | `@BeforeEach/@AfterEach`     | Test classes              |
   | Go           | `TestMain`, helper functions | `*_test.go`               |
   | Ruby (RSpec) | `let`, `before/after`        | `spec/support/`           |
   | C#           | `[SetUp]/[TearDown]`         | Test classes              |

6. **Implement Data Factories (Language-Agnostic)**

   **Knowledge Base Reference**: `testarch/knowledge/data-factories.md`

   **Universal Factory Principles:**
   - Generate realistic fake data (no hardcoded values)
   - Support overrides for specific test scenarios
   - Track created resources for cleanup
   - Parallel-safe (unique IDs, no collisions)

   **Fake Data Libraries by Language:**
   | Language | Library | Docs |
   |----------|---------|------|
   | TypeScript/JS | @faker-js/faker | <https://fakerjs.dev> |
   | Python | Faker | <https://faker.readthedocs.io> |
   | Java | JavaFaker, DataFaker | <https://github.com/datafaker-net/datafaker> |
   | Go | gofakeit | <https://github.com/brianvoe/gofakeit> |
   | Ruby | Faker | <https://github.com/faker-ruby/faker> |
   | C# | Bogus | <https://github.com/bchavez/Bogus> |
   | PHP | FakerPHP | <https://fakerphp.github.io> |
   | Rust | fake-rs | <https://github.com/cksac/fake-rs> |

   **CRITICAL**: Web fetch factory patterns for detected language:

   ```
   1. Web fetch: "[language] test data factory pattern"
   2. Web fetch: "[faker_library] documentation"
   3. Generate factories using CURRENT API from fetched docs
   ```

7. **Generate Sample Tests (Language-Specific)**

   Generate sample tests using detected framework and conventions:

   **Universal Test Structure (Given-When-Then):**

   ```
   // GIVEN: Initial state/setup
   // WHEN: Action being tested
   // THEN: Expected outcome (assertion)
   ```

   **CRITICAL**: Generate samples using freshly fetched framework docs:

   ```
   1. Web fetch: "[framework] writing tests guide"
   2. Web fetch: "[framework] assertions reference"
   3. Generate tests matching project conventions
   4. Include comments linking to relevant docs
   ```

   **Sample test should demonstrate:**
   - Basic test structure for the framework
   - Fixture usage
   - Factory usage
   - Assertion patterns
   - Async handling (if applicable)

8. **Update Build/Package Scripts**

   Add test scripts appropriate for the language:

   **JavaScript/TypeScript (package.json):**

   ```json
   { "scripts": { "test": "[framework] test", "test:e2e": "[e2e-framework] test" } }
   ```

   **Python (pyproject.toml or setup.cfg):**

   ```toml
   [tool.pytest.ini_options]
   testpaths = ["tests"]
   ```

   **Go (Makefile or go task):**

   ```makefile
   test: go test ./...
   ```

   **Java (Maven pom.xml or Gradle):**
   Add surefire/failsafe plugin configuration

   **Fetch latest recommended scripts:**

   ```
   Web fetch: "[framework] npm scripts" or "[framework] recommended commands"
   ```

9. **Generate Documentation**

   Create `tests/README.md` with:
   - Setup instructions for detected language/framework
   - Commands to run tests (fetched from official docs)
   - Links to official documentation
   - Architecture overview (fixtures, factories, conventions)

---

## Step 3: Deliverables

### Primary Artifacts Created

1. **Configuration File**
   - `playwright.config.ts` or `cypress.config.ts`
   - Timeouts: action 15s, navigation 30s, test 60s
   - Reporters: HTML + JUnit XML

2. **Directory Structure**
   - `tests/` with `e2e/`, `api/`, `support/` subdirectories
   - `support/fixtures/` for test fixtures
   - `support/helpers/` for utility functions

3. **Environment Configuration**
   - `.env.example` with `TEST_ENV`, `BASE_URL`, `API_URL`
   - `.nvmrc` with Node version

4. **Test Infrastructure**
   - Fixture architecture (`mergeTests` pattern)
   - Data factories (faker-based, with auto-cleanup)
   - Sample tests demonstrating patterns

5. **Documentation**
   - `tests/README.md` with setup instructions
   - Comments in config files explaining options

### README Contents

The generated `tests/README.md` should include:

- **Setup Instructions**: How to install dependencies, configure environment
- **Running Tests**: Commands for local execution, headed mode, debug mode
- **Architecture Overview**: Fixture pattern, data factories, page objects
- **Best Practices**: Selector strategy (data-testid), test isolation, cleanup
- **CI Integration**: How tests run in CI/CD pipeline
- **Knowledge Base References**: Links to relevant TEA knowledge fragments

---

## Important Notes

### Language-Agnostic Approach

**This workflow adapts to ANY programming language.** The key principle is:

1. **Detect** the project's language and existing conventions
2. **Fetch** latest documentation for the selected framework
3. **Generate** code using current patterns from official docs
4. **Apply** universal testing principles regardless of language

### Knowledge Base Integration

**Critical:** Load knowledge fragments based on detected language and framework.

Consult `{project-root}/_bmad/bmm/testarch/tea-index.csv` and load relevant fragments.

**Universal Fragments (load for all languages):**

- `fixture-architecture.md` - Fixture patterns and composition principles
- `data-factories.md` - Factory patterns with fake data generation
- `test-quality.md` - Test design principles (deterministic, isolated, explicit assertions)

**Language-Specific Guidance:**
After loading universal fragments, adapt patterns to detected language using:

1. Freshly fetched documentation for the framework
2. Existing project conventions (from Step 0 detection)
3. Language-specific idioms and best practices

### Real-Time Documentation Fetching

**MANDATORY for every code generation:**

```
Before generating config, fixtures, factories, or tests:
1. Web fetch official documentation for detected framework
2. Check for latest version and API changes
3. Generate code using CURRENT patterns (not memorized)
4. Include documentation links in generated code comments
```

**Why this matters:**

- Framework APIs change frequently
- Best practices evolve over time
- Static templates become outdated
- Real-time fetch ensures accuracy

### Universal Testing Principles (All Languages)

**Selector Strategy (for UI tests):**

- Use stable identifiers: `data-testid`, `data-test`, `test-id`
- Avoid brittle selectors: CSS classes, XPath, DOM structure
- Prefer accessibility attributes when meaningful

**Test Isolation:**

- Each test runs independently
- No shared mutable state between tests
- Auto-cleanup of created resources

**Deterministic Tests:**

- No flaky tests (race conditions, timing issues)
- Explicit waits instead of sleep/delays
- Mocked external dependencies where appropriate

**Failure Artifacts:**

- Capture on failure: screenshots, logs, traces
- Delete on success to save storage
- Include enough context for debugging

### Contract Testing

For microservices architectures, recommend contract testing:

- **Pact**: Consumer-driven contracts (multiple languages)
- **Spring Cloud Contract**: JVM-based services
- **Dredd**: API Blueprint / OpenAPI validation

### Framework Selection Guidance

**E2E Framework Selection (by ecosystem):**

| Scenario       | TypeScript/JS        | Python            | Java                | C#/.NET    | Go       |
| -------------- | -------------------- | ----------------- | ------------------- | ---------- | -------- |
| Web UI testing | Playwright           | Playwright        | Selenium/Playwright | Playwright | Chromedp |
| API testing    | Playwright/SuperTest | pytest + requests | REST Assured        | RestSharp  | net/http |
| Mobile         | Appium               | Appium            | Appium              | Appium     | gomobile |

**Unit Test Framework Selection:**

- Choose framework with best ecosystem integration
- Prefer frameworks with good IDE support
- Consider parallel execution capabilities
- Check community size and maintenance status

This reduces storage overhead while maintaining debugging capability.

---

## Output Summary

After completing this workflow, provide a summary:

```markdown
## Framework Scaffold Complete

**Framework Selected**: Playwright (or Cypress)

**Artifacts Created**:

- ✅ Configuration file: `playwright.config.ts`
- ✅ Directory structure: `tests/e2e/`, `tests/support/`
- ✅ Environment config: `.env.example`
- ✅ Node version: `.nvmrc`
- ✅ Fixture architecture: `tests/support/fixtures/`
- ✅ Data factories: `tests/support/fixtures/factories/`
- ✅ Sample tests: `tests/e2e/example.spec.ts`
- ✅ Documentation: `tests/README.md`

**Next Steps**:

1. Copy `.env.example` to `.env` and fill in environment variables
2. Run `npm install` to install test dependencies
3. Run `npm run test:e2e` to execute sample tests
4. Review `tests/README.md` for detailed setup instructions

**Knowledge Base References Applied**:

- Fixture architecture pattern (pure functions + mergeTests)
- Data factories with auto-cleanup (faker-based)
- Network-first testing safeguards
- Failure-only artifact capture
```

---

## Validation

After completing all steps, verify:

- [ ] Configuration file created and valid
- [ ] Directory structure exists
- [ ] Environment configuration generated
- [ ] Sample tests run successfully
- [ ] Documentation complete and accurate
- [ ] No errors or warnings during scaffold

Refer to `checklist.md` for comprehensive validation criteria.
