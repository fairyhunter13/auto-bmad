<!-- Powered by BMAD-CORE™ -->

# Test Automation Expansion

**Workflow ID**: `_bmad/bmm/testarch/automate`
**Version**: 5.0 (Language Agnostic)

---

## Overview

Expands test automation coverage by generating comprehensive test suites at appropriate levels (E2E, API, Component, Unit) with supporting infrastructure. This workflow operates in **dual mode**:

1. **BMad-Integrated Mode**: Works WITH BMad artifacts (story, tech-spec, PRD, test-design) to expand coverage after story implementation
2. **Standalone Mode**: Works WITHOUT BMad artifacts - analyzes existing codebase and generates tests independently

**Language Agnostic**: This workflow supports ANY programming language and test framework. Detection happens automatically in Step 0.

**Core Principle**: Generate prioritized, deterministic tests that avoid duplicate coverage and follow testing best practices.

---

## Step 0: Detect Project Environment (MANDATORY)

**CRITICAL**: Before generating ANY test code, you MUST detect the project's language and test framework.

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
   - **Python**: `pytest.ini`, `conftest.py`, `tox.ini`, `unittest` in code
   - **Go**: `*_test.go` files (native testing), `testify` imports
   - **Java**: JUnit 5/6 in `pom.xml`, TestNG, Mockito
   - **Rust**: `#[test]` attributes (native), `cargo test`
   - **C#**: xUnit, NUnit, MSTest in `*.csproj`
   - **Ruby**: RSpec (`spec/`), Minitest
   - **PHP**: PHPUnit (`phpunit.xml`)

3. **Fetch Latest Framework Documentation (MANDATORY)**

   **CRITICAL**: Before generating ANY test code, fetch current documentation:

   ```
   Search: "[DETECTED_FRAMEWORK] test patterns best practices [CURRENT_YEAR]"
   Search: "[DETECTED_FRAMEWORK] fixtures factories setup [CURRENT_YEAR]"
   Search: "[DETECTED_LANGUAGE] testing Given-When-Then [CURRENT_YEAR]"
   ```

   Include fetch timestamp in generated code comments:

   ```
   // Generated: [TIMESTAMP]
   // Framework: [DETECTED_FRAMEWORK] v[VERSION]
   // Docs fetched: [CURRENT_DATE]
   ```

4. **Store Detection Results**

   ```yaml
   detected:
     language: '[JavaScript|Python|Go|Java|Rust|C#|Ruby|PHP|Other]'
     test_framework: '[e.g., Playwright, pytest, go test, JUnit 6]'
     assertion_library: '[e.g., expect, pytest.assert, testify/assert]'
     fixture_pattern: '[e.g., test.extend, pytest fixtures, testify suites]'
     factory_library: '[e.g., faker-js, faker, gofakeit, java-faker]'
     test_directory: '[e.g., tests/, test/, spec/]'
   ```

**Halt Condition**: If language cannot be detected, ask user: "What programming language and test framework does this project use?"

---

## Preflight Requirements

**Flexible:** This workflow can run with minimal prerequisites. Only HALT if framework is completely missing.

### Required (Always)

- ✅ Framework scaffolding configured (run `framework` workflow if missing)
- ✅ Test framework configuration detected in Step 0

### Optional (BMad-Integrated Mode)

- Story markdown with acceptance criteria (enhances coverage targeting)
- Tech spec or PRD (provides architectural context)
- Test design document (provides risk/priority context)

### Optional (Standalone Mode)

- Source code to analyze (feature implementation)
- Existing tests (for gap analysis)

**If framework is missing:** HALT with message: "Framework scaffolding required. Run `bmad tea *framework` first."

---

## Step 1: Determine Execution Mode and Load Context

### Actions

1. **Verify Step 0 Detection Complete**

   Ensure language and framework detection from Step 0 is available:
   - Detected language (JavaScript, Python, Go, Java, etc.)
   - Detected test framework (Playwright, pytest, go test, JUnit, etc.)
   - Test directory structure
   - Fixture/factory patterns for this framework

2. **Detect Execution Mode**

   Check if BMad artifacts are available:
   - If `{story_file}` variable is set → BMad-Integrated Mode
   - If `{target_feature}` or `{target_files}` set → Standalone Mode
   - If neither set → Auto-discover mode (scan codebase for features needing tests)

3. **Load BMad Artifacts (If Available)**

   **BMad-Integrated Mode:**
   - Read story markdown from `{story_file}`
   - Extract acceptance criteria and technical requirements
   - Load tech-spec.md if `{use_tech_spec}` is true
   - Load test-design.md if `{use_test_design}` is true
   - Load PRD.md if `{use_prd}` is true
   - Note: These are **optional enhancements**, not hard requirements

   **Standalone Mode:**
   - Skip BMad artifact loading
   - Proceed directly to source code analysis

4. **Load Framework Configuration** (Language-Agnostic)

   Read test framework config based on detected language:
   - **JavaScript/TypeScript**: `playwright.config.*`, `jest.config.*`, `vitest.config.*`
   - **Python**: `pytest.ini`, `conftest.py`, `pyproject.toml`
   - **Go**: Test files in `*_test.go`, go.mod for dependencies
   - **Java**: `pom.xml` or `build.gradle` test configuration
   - **Rust**: `Cargo.toml` test configuration
   - **C#**: `*.csproj` test project configuration
   - **Ruby**: `.rspec`, `spec_helper.rb`
   - **PHP**: `phpunit.xml`

   Extract from config:
   - Test directory structure from `{test_dir}`
   - Existing test patterns
   - Test runner capabilities (parallel execution, fixtures, etc.)

5. **Analyze Existing Test Coverage**

   If `{analyze_coverage}` is true:
   - Search `{test_dir}` for existing test files
   - Identify tested features vs untested features
   - Map tests to source files (coverage gaps)
   - Check existing fixture and factory patterns (language-specific)

6. **Check Framework-Specific Utils Flag**

   Read `{config_source}` and check `config.tea_use_playwright_utils`.

   **Note**: This flag is specific to JavaScript/TypeScript Playwright projects. For other languages, use equivalent framework-specific patterns.

7. **Load Knowledge Base Fragments** (Language-Agnostic Principles)

   **Critical:** Consult `{project-root}/_bmad/bmm/testarch/tea-index.csv` to load:

   **Core Testing Patterns (Always load - universal principles):**
   - `test-levels-framework.md` - Test level selection (E2E vs API vs Component vs Unit)
   - `test-priorities-matrix.md` - Priority classification (P0-P3)
   - `data-factories.md` - Factory patterns (adapt to detected language's faker library)
   - `selective-testing.md` - Targeted test execution strategies
   - `ci-burn-in.md` - Flaky test detection patterns
   - `test-quality.md` - Test design principles (Given-When-Then, isolation, cleanup)

   **If JavaScript/TypeScript with `config.tea_use_playwright_utils: true`:**
   - Load Playwright-specific utility fragments (overview.md, api-request.md, etc.)

   **For other languages, fetch framework-specific patterns:**

   ```
   Search: "[DETECTED_FRAMEWORK] fixture patterns [CURRENT_YEAR]"
   Search: "[DETECTED_FRAMEWORK] test factories [CURRENT_YEAR]"
   Search: "[DETECTED_FRAMEWORK] async testing patterns [CURRENT_YEAR]"
   ```

   **Healing Knowledge (If `{auto_heal_failures}` is true):**
   - `test-healing-patterns.md` - Common failure patterns (adapt to detected framework)
   - `selector-resilience.md` - Selector debugging (for E2E frameworks)
   - `timing-debugging.md` - Race condition identification

---

## Step 2: Identify Automation Targets

### Actions

1. **Determine What Needs Testing**

   **BMad-Integrated Mode (story available):**
   - Map acceptance criteria from story to test scenarios
   - Identify features implemented in this story
   - Check if story has existing ATDD tests (from `*atdd` workflow)
   - Expand beyond ATDD with edge cases and negative paths

   **Standalone Mode (no story):**
   - If `{target_feature}` specified: Analyze that specific feature
   - If `{target_files}` specified: Analyze those specific files
   - If `{auto_discover_features}` is true: Scan `{source_dir}` for features
   - Prioritize features with:
     - No test coverage (highest priority)
     - Complex business logic
     - External integrations (API calls, database, auth)
     - Critical user paths (login, checkout, etc.)

2. **Apply Test Level Selection Framework**

   **Knowledge Base Reference**: `test-levels-framework.md`

   For each feature or acceptance criterion, determine appropriate test level:

   **E2E (End-to-End)**:
   - Critical user journeys (login, checkout, core workflows)
   - Multi-system integration
   - Full user-facing scenarios
   - Characteristics: High confidence, slow, brittle

   **API (Integration)**:
   - Business logic validation
   - Service contracts and data transformations
   - Backend integration without UI
   - Characteristics: Fast feedback, stable, good balance

   **Component**:
   - UI component behavior (buttons, forms, modals)
   - Interaction testing (click, hover, keyboard)
   - State management within component
   - Characteristics: Fast, isolated, granular

   **Unit**:
   - Pure business logic and algorithms
   - Edge cases and error handling
   - Minimal dependencies
   - Characteristics: Fastest, most granular

3. **Avoid Duplicate Coverage**

   **Critical principle:** Don't test same behavior at multiple levels unless necessary
   - Use E2E for critical happy path only
   - Use API tests for business logic variations
   - Use component tests for UI interaction edge cases
   - Use unit tests for pure logic edge cases

   **Example:**
   - E2E: User can log in with valid credentials → Dashboard loads
   - API: POST /auth/login returns 401 for invalid credentials
   - API: POST /auth/login returns 200 and JWT token for valid credentials
   - Component: LoginForm disables submit button when fields are empty
   - Unit: validateEmail() returns false for malformed email addresses

4. **Assign Test Priorities**

   **Knowledge Base Reference**: `test-priorities-matrix.md`

   **P0 (Critical - Every commit)**:
   - Critical user paths that must always work
   - Security-critical functionality (auth, permissions)
   - Data integrity scenarios
   - Run in pre-commit hooks or PR checks

   **P1 (High - PR to main)**:
   - Important features with high user impact
   - Integration points between systems
   - Error handling for common failures
   - Run before merging to main branch

   **P2 (Medium - Nightly)**:
   - Edge cases with moderate impact
   - Less-critical feature variations
   - Performance/load testing
   - Run in nightly CI builds

   **P3 (Low - On-demand)**:
   - Nice-to-have validations
   - Rarely-used features
   - Exploratory testing scenarios
   - Run manually or weekly

   **Priority Variables:**
   - `{include_p0}` - Always include (default: true)
   - `{include_p1}` - High priority (default: true)
   - `{include_p2}` - Medium priority (default: true)
   - `{include_p3}` - Low priority (default: false)

5. **Create Test Coverage Plan**

   Document what will be tested at each level with priorities:

   ```markdown
   ## Test Coverage Plan

   ### E2E Tests (P0)

   - User login with valid credentials → Dashboard loads
   - User logout → Redirects to login page

   ### API Tests (P1)

   - POST /auth/login - valid credentials → 200 + JWT token
   - POST /auth/login - invalid credentials → 401 + error message
   - POST /auth/login - missing fields → 400 + validation errors

   ### Component Tests (P1)

   - LoginForm - empty fields → submit button disabled
   - LoginForm - valid input → submit button enabled

   ### Unit Tests (P2)

   - validateEmail() - valid email → returns true
   - validateEmail() - malformed email → returns false
   ```

---

## Step 3: Generate Test Infrastructure

**CRITICAL**: Before generating ANY infrastructure code, fetch latest patterns for detected framework:

```
Search: "[DETECTED_FRAMEWORK] fixture patterns [CURRENT_YEAR]"
Search: "[DETECTED_LANGUAGE] test factory library [CURRENT_YEAR]"
```

### Actions

1. **Enhance Fixture Architecture** (Language-Agnostic Pattern)

   **Knowledge Base Reference**: `fixture-architecture.md`

   **Universal fixture principles (apply to ALL languages):**
   - Setup runs BEFORE each test
   - Teardown runs AFTER each test (cleanup)
   - Fixtures provide test data/state to tests
   - Auto-cleanup prevents test pollution

   **Language-specific fixture patterns:**

   | Language              | Fixture Pattern                | Example                        |
   | --------------------- | ------------------------------ | ------------------------------ |
   | JavaScript/Playwright | `test.extend()`                | Extend base test with fixtures |
   | JavaScript/Jest       | `beforeEach`/`afterEach`       | Setup/teardown hooks           |
   | Python/pytest         | `@pytest.fixture`              | Decorator-based fixtures       |
   | Go                    | `TestMain`, `t.Cleanup()`      | Test setup and cleanup         |
   | Java/JUnit            | `@BeforeEach`/`@AfterEach`     | Lifecycle annotations          |
   | Rust                  | `#[test]` with setup functions | Function-based setup           |
   | C#/xUnit              | `IClassFixture<T>`             | Shared fixtures                |
   | Ruby/RSpec            | `let`, `before`, `after`       | DSL-based fixtures             |

   **Common fixtures to create (adapt for detected language):**
   - **authenticatedUser**: User with valid session (auto-cleanup)
   - **apiClient**: Authenticated HTTP client
   - **mockNetwork**: Network mocking/stubbing
   - **testDatabase**: Database with test data (auto-cleanup)

   **Fetch framework-specific fixture example:**

   ```
   Search: "[DETECTED_FRAMEWORK] fixture example with cleanup [CURRENT_YEAR]"
   ```

   Include generation timestamp in fixture comments.

2. **Enhance Data Factories** (Language-Agnostic Pattern)

   **Knowledge Base Reference**: `data-factories.md`

   **Universal factory principles (apply to ALL languages):**
   - Generate realistic fake data (no hardcoded values)
   - Support overrides for specific scenarios
   - Create nested/related objects
   - Provide cleanup helpers

   **Language-specific faker libraries:**

   | Language   | Faker Library     | Install Command                       |
   | ---------- | ----------------- | ------------------------------------- |
   | JavaScript | `@faker-js/faker` | `npm install @faker-js/faker`         |
   | Python     | `faker`           | `pip install faker`                   |
   | Go         | `gofakeit`        | `go get github.com/brianvoe/gofakeit` |
   | Java       | `javafaker`       | Maven/Gradle dependency               |
   | Rust       | `fake`            | `cargo add fake`                      |
   | C#         | `Bogus`           | `dotnet add package Bogus`            |
   | Ruby       | `faker`           | `gem install faker`                   |
   | PHP        | `fakerphp/faker`  | `composer require fakerphp/faker`     |

   **Fetch factory example for detected language:**

   ```
   Search: "[DETECTED_LANGUAGE] [FAKER_LIBRARY] factory pattern [CURRENT_YEAR]"
   ```

   **Common factories to create (adapt for detected language):**
   - User factory (email, password, name, role)
   - Product factory (name, price, description, SKU)
   - Order factory (items, total, status, customer)

3. **Create/Enhance Helper Utilities** (Language-Agnostic)

   If `{update_helpers}` is true:

   **Common utilities to create (adapt for detected language):**
   - **waitFor/poll**: Polling helper for async conditions
   - **retry**: Retry helper for flaky operations
   - **testData**: Test data generation helpers
   - **assertions**: Custom assertion helpers

   **Fetch helper patterns for detected language:**

   ```
   Search: "[DETECTED_FRAMEWORK] async wait helper [CURRENT_YEAR]"
   Search: "[DETECTED_LANGUAGE] retry pattern testing [CURRENT_YEAR]"
   ```

---

## Step 4: Generate Test Files

**CRITICAL**: Before generating test code, fetch latest syntax for detected framework:

```
Search: "[DETECTED_FRAMEWORK] test syntax [CURRENT_YEAR]"
Search: "[DETECTED_FRAMEWORK] assertion examples [CURRENT_YEAR]"
```

### Actions

1. **Create Test File Structure** (Language-Specific Naming)

   **Adapt file extensions and naming conventions for detected language:**

   | Language              | E2E Tests             | API Tests       | Unit Tests   |
   | --------------------- | --------------------- | --------------- | ------------ |
   | JavaScript/TypeScript | `*.spec.ts`           | `*.api.spec.ts` | `*.test.ts`  |
   | Python                | `test_*.py`           | `test_*_api.py` | `test_*.py`  |
   | Go                    | `*_test.go`           | `*_api_test.go` | `*_test.go`  |
   | Java                  | `*Test.java`          | `*ApiTest.java` | `*Test.java` |
   | Rust                  | `mod tests` in source | `*_api_test.rs` | `mod tests`  |
   | C#                    | `*Tests.cs`           | `*ApiTests.cs`  | `*Tests.cs`  |
   | Ruby                  | `*_spec.rb`           | `*_api_spec.rb` | `*_spec.rb`  |
   | PHP                   | `*Test.php`           | `*ApiTest.php`  | `*Test.php`  |

2. **Write E2E Tests (If Applicable)**

   **CRITICAL**: Fetch framework-specific E2E syntax before generating:

   ```
   Search: "[DETECTED_E2E_FRAMEWORK] test example [CURRENT_YEAR]"
   ```

   **Universal Given-When-Then pattern (adapt syntax for detected framework):**

   ```
   // Generated: [TIMESTAMP]
   // Framework: [DETECTED_FRAMEWORK]

   [FRAMEWORK_TEST_STRUCTURE] 'User Authentication' {

     [TEST_FUNCTION] '[P0] should login with valid credentials and load dashboard' {
       // GIVEN: User is on login page
       [NAVIGATE_TO]('/login')

       // WHEN: User submits valid credentials
       [FILL_INPUT]('[data-testid="email-input"]', 'user@example.com')
       [FILL_INPUT]('[data-testid="password-input"]', 'Password123!')
       [CLICK]('[data-testid="login-button"]')

       // THEN: User is redirected to dashboard
       [ASSERT_URL]('/dashboard')
       [ASSERT_VISIBLE]('[data-testid="user-name"]')
     }
   }
   ```

   **Critical patterns (universal for ALL frameworks):**
   - Tag tests with priority: `[P0]`, `[P1]`, `[P2]`, `[P3]` in test name
   - One assertion per test (atomic tests)
   - Explicit waits (no hard waits/sleeps)
   - Network-first approach (setup mocks before navigation)
   - Stable selectors (data-testid, ARIA roles, semantic HTML)
   - Clear Given-When-Then structure

3. **Write API Tests (If Applicable)**

   **Fetch framework-specific API testing syntax:**

   ```
   Search: "[DETECTED_FRAMEWORK] API testing HTTP requests [CURRENT_YEAR]"
   ```

   **Universal API test pattern (adapt for detected framework):**

   ```
   [TEST_FUNCTION] '[P1] POST /api/auth/login - should return token for valid credentials' {
     // GIVEN: Valid user credentials
     credentials = { email: 'user@example.com', password: 'Password123!' }

     // WHEN: Logging in via API
     response = [HTTP_POST]('/api/auth/login', credentials)

     // THEN: Returns 200 and JWT token
     [ASSERT_STATUS](response, 200)
     [ASSERT_HAS_PROPERTY](response.body, 'token')
     [ASSERT_MATCHES](response.body.token, JWT_PATTERN)
   }
   ```

4. **Write Component Tests (If Applicable)**

   **Knowledge Base Reference**: `component-tdd.md`

   **Fetch framework-specific component testing syntax:**

   ```
   Search: "[DETECTED_LANGUAGE] component testing [CURRENT_YEAR]"
   ```

   **Component testing frameworks by language:**
   - **JavaScript/React**: Playwright CT, Testing Library, Vitest
   - **JavaScript/Vue**: Vue Test Utils, Vitest
   - **Python/Django**: pytest-django
   - **Go**: Native testing with template rendering
   - **Java/Spring**: MockMvc, WebTestClient
   - **C#/Blazor**: bUnit

5. **Write Unit Tests (If Applicable)**

   **Fetch framework-specific unit test syntax:**

   ```
   Search: "[DETECTED_FRAMEWORK] unit test example [CURRENT_YEAR]"
   ```

   **Universal unit test pattern:**

   ```
   [TEST_FUNCTION] '[P2] should return true for valid email' {
     // GIVEN: Valid email address
     email = 'user@example.com'

     // WHEN: Validating email
     result = validateEmail(email)

     // THEN: Returns true
     [ASSERT_TRUE](result)
   }
   ```

6. **Apply Network-First Pattern (E2E tests)**

   **Knowledge Base Reference**: `network-first.md`

   **Critical pattern (adapt for detected E2E framework):**
   - Playwright: `page.route()` before `page.goto()`
   - Cypress: `cy.intercept()` before `cy.visit()`
   - Selenium: Proxy setup or mock server
   - Others: Framework-specific network mocking

   **Fetch framework-specific mocking syntax:**

   ```
   Search: "[DETECTED_E2E_FRAMEWORK] network mocking [CURRENT_YEAR]"
   ```

7. **Enforce Quality Standards** (Universal)

   **For every test (ALL languages):**
   - ✅ Uses Given-When-Then format (comments or structure)
   - ✅ Has clear, descriptive name with priority tag
   - ✅ One assertion per test (atomic)
   - ✅ No hard waits or sleeps (use explicit waits)
   - ✅ Self-cleaning (uses fixtures with auto-cleanup)
   - ✅ Deterministic (no flaky patterns)
   - ✅ Fast (under {max_test_duration} seconds)
   - ✅ Lean (test file under {max_file_lines} lines)

   **Forbidden patterns:**
   - ❌ Hard waits: `await page.waitForTimeout(2000)`
   - ❌ Conditional flow: `if (await element.isVisible()) { ... }`
   - ❌ Try-catch for test logic (use for cleanup only)
   - ❌ Hardcoded test data (use factories)
   - ❌ Page objects (keep tests simple and direct)
   - ❌ Shared state between tests

---

## Step 5: Execute, Validate & Heal Generated Tests (NEW - Phase 2.5)

**Purpose**: Automatically validate generated tests and heal common failures before delivery

### Actions

1. **Validate Generated Tests**

   Always validate (auto_validate is always true):
   - Run generated tests to verify they work
   - Continue with healing if config.tea_use_mcp_enhancements is true

2. **Run Generated Tests**

   Execute the full test suite that was just generated:

   ```bash
   npx playwright test {generated_test_files}
   ```

   Capture results:
   - Total tests run
   - Passing tests count
   - Failing tests count
   - Error messages and stack traces for failures

3. **Evaluate Results**

   **If ALL tests pass:**
   - ✅ Generate report with success summary
   - Proceed to Step 6 (Documentation and Scripts)

   **If tests FAIL:**
   - Check config.tea_use_mcp_enhancements setting
   - If true: Enter healing loop (Step 5.4)
   - If false: Document failures for manual review, proceed to Step 6

4. **Healing Loop (If config.tea_use_mcp_enhancements is true)**

   **Iteration limit**: 3 attempts per test (constant)

   **For each failing test:**

   **A. Load Healing Knowledge Fragments**

   Consult `tea-index.csv` to load healing patterns:
   - `test-healing-patterns.md` - Common failure patterns and fixes
   - `selector-resilience.md` - Selector debugging and refactoring
   - `timing-debugging.md` - Race condition identification and fixes

   **B. Identify Failure Pattern**

   Analyze error message and stack trace to classify failure type:

   **Stale Selector Failure:**
   - Error contains: "locator resolved to 0 elements", "element not found", "unable to find element"
   - Extract selector from error message
   - Apply selector healing (knowledge from `selector-resilience.md`):
     - If CSS class → Replace with `page.getByTestId()`
     - If nth() → Replace with `filter({ hasText })`
     - If ID → Replace with data-testid
     - If complex XPath → Replace with ARIA role

   **Race Condition Failure:**
   - Error contains: "timeout waiting for", "element not visible", "timed out retrying"
   - Detect missing network waits or hard waits in test code
   - Apply timing healing (knowledge from `timing-debugging.md`):
     - Add network-first interception before navigate
     - Replace `waitForTimeout()` with `waitForResponse()`
     - Add explicit element state waits (`waitFor({ state: 'visible' })`)

   **Dynamic Data Failure:**
   - Error contains: "Expected 'User 123' but received 'User 456'", timestamp mismatches
   - Identify hardcoded assertions
   - Apply data healing (knowledge from `test-healing-patterns.md`):
     - Replace hardcoded IDs with regex (`/User \d+/`)
     - Replace hardcoded dates with dynamic generation
     - Capture dynamic values and use in assertions

   **Network Error Failure:**
   - Error contains: "API call failed", "500 error", "network error"
   - Detect missing route interception
   - Apply network healing (knowledge from `test-healing-patterns.md`):
     - Add `page.route()` or `cy.intercept()` for API mocking
     - Mock error scenarios (500, 429, timeout)

   **Hard Wait Detection:**
   - Scan test code for `page.waitForTimeout()`, `cy.wait(number)`, `sleep()`
   - Apply hard wait healing (knowledge from `timing-debugging.md`):
     - Replace with event-based waits
     - Add network response waits
     - Use element state changes

   **C. MCP Healing Mode (If MCP Tools Available)**

   If Playwright MCP tools are available in your IDE:

   Use MCP tools for interactive healing:
   - `playwright_test_debug_test`: Pause on failure for visual inspection
   - `browser_snapshot`: Capture visual context at failure point
   - `browser_console_messages`: Retrieve console logs for JS errors
   - `browser_network_requests`: Analyze network activity
   - `browser_generate_locator`: Generate better selectors interactively

   Apply MCP-generated fixes to test code.

   **D. Pattern-Based Healing Mode (Fallback)**

   If MCP unavailable, use pattern-based analysis:
   - Parse error message and stack trace
   - Match against failure patterns from knowledge base
   - Apply fixes programmatically:
     - Selector fixes: Use suggestions from `selector-resilience.md`
     - Timing fixes: Apply patterns from `timing-debugging.md`
     - Data fixes: Use patterns from `test-healing-patterns.md`

   **E. Apply Healing Fix**
   - Modify test file with healed code
   - Re-run test to validate fix
   - If test passes: Mark as healed, move to next failure
   - If test fails: Increment iteration count, try different pattern

   **F. Iteration Limit Handling**

   After 3 failed healing attempts:

   Always mark unfixable tests:
   - Mark test with `test.fixme()` instead of `test()`
   - Add detailed comment explaining:
     - What failure occurred
     - What healing was attempted (3 iterations)
     - Why healing failed
     - Manual investigation needed

   ```typescript
   test.fixme('[P1] should handle complex interaction', async ({ page }) => {
     // FIXME: Test healing failed after 3 attempts
     // Failure: "Locator 'button[data-action="submit"]' resolved to 0 elements"
     // Attempted fixes:
     //   1. Replaced with page.getByTestId('submit-button') - still failing
     //   2. Replaced with page.getByRole('button', { name: 'Submit' }) - still failing
     //   3. Added waitForLoadState('networkidle') - still failing
     // Manual investigation needed: Selector may require application code changes
     // TODO: Review with team, may need data-testid added to button component
     // Original test code...
   });
   ```

   **Note**: Workflow continues even with unfixable tests (marked as test.fixme() for manual review)

5. **Generate Healing Report**

   Document healing outcomes:

   ```markdown
   ## Test Healing Report

   **Auto-Heal Enabled**: {auto_heal_failures}
   **Healing Mode**: {use_mcp_healing ? "MCP-assisted" : "Pattern-based"}
   **Iterations Allowed**: {max_healing_iterations}

   ### Validation Results

   - **Total tests**: {total_tests}
   - **Passing**: {passing_tests}
   - **Failing**: {failing_tests}

   ### Healing Outcomes

   **Successfully Healed ({healed_count} tests):**

   - `tests/e2e/login.spec.ts:15` - Stale selector (CSS class → data-testid)
   - `tests/e2e/checkout.spec.ts:42` - Race condition (added network-first interception)
   - `tests/api/users.spec.ts:28` - Dynamic data (hardcoded ID → regex pattern)

   **Unable to Heal ({unfixable_count} tests):**

   - `tests/e2e/complex-flow.spec.ts:67` - Marked as test.fixme() with manual investigation needed
     - Failure: Locator not found after 3 healing attempts
     - Requires application code changes (add data-testid to component)

   ### Healing Patterns Applied

   - **Selector fixes**: 2 (CSS class → data-testid, nth() → filter())
   - **Timing fixes**: 1 (added network-first interception)
   - **Data fixes**: 1 (hardcoded ID → regex)

   ### Knowledge Base References

   - `test-healing-patterns.md` - Common failure patterns
   - `selector-resilience.md` - Selector refactoring guide
   - `timing-debugging.md` - Race condition prevention
   ```

6. **Update Test Files with Healing Results**
   - Save healed test code to files
   - Mark unfixable tests with `test.fixme()` and detailed comments
   - Preserve original test logic in comments (for debugging)

---

## Step 6: Update Documentation and Scripts

### Actions

1. **Update Test README**

   If `{update_readme}` is true:

   Create or update `tests/README.md` with:
   - Overview of test suite structure
   - How to run tests (all, specific files, by priority)
   - Fixture and factory usage examples
   - Priority tagging convention ([P0], [P1], [P2], [P3])
   - How to write new tests
   - Common patterns and anti-patterns

   **Example section:**

   ````markdown
   ## Running Tests

   ```bash
   # Run all tests
   npm run test:e2e

   # Run by priority
   npm run test:e2e -- --grep "@P0"
   npm run test:e2e -- --grep "@P1"

   # Run specific file
   npm run test:e2e -- user-authentication.spec.ts

   # Run in headed mode
   npm run test:e2e -- --headed

   # Debug specific test
   npm run test:e2e -- user-authentication.spec.ts --debug
   ```
   ````

   ## Priority Tags
   - **[P0]**: Critical paths, run every commit
   - **[P1]**: High priority, run on PR to main
   - **[P2]**: Medium priority, run nightly
   - **[P3]**: Low priority, run on-demand

   ```

   ```

2. **Update package.json Scripts**

   If `{update_package_scripts}` is true:

   Add or update test execution scripts:

   ```json
   {
     "scripts": {
       "test:e2e": "playwright test",
       "test:e2e:p0": "playwright test --grep '@P0'",
       "test:e2e:p1": "playwright test --grep '@P1|@P0'",
       "test:api": "playwright test tests/api",
       "test:component": "playwright test tests/component",
       "test:unit": "vitest"
     }
   }
   ```

3. **Run Test Suite**

   If `{run_tests_after_generation}` is true:
   - Run full test suite locally
   - Capture results (passing/failing counts)
   - Verify no flaky patterns (tests should be deterministic)
   - Document any setup requirements or known issues

---

## Step 6: Generate Automation Summary

### Actions

1. **Create Automation Summary Document**

   Save to `{output_summary}` with:

   **BMad-Integrated Mode:**

   ````markdown
   # Automation Summary - {feature_name}

   **Date:** {date}
   **Story:** {story_id}
   **Coverage Target:** {coverage_target}

   ## Tests Created

   ### E2E Tests (P0-P1)

   - `tests/e2e/user-authentication.spec.ts` (2 tests, 87 lines)
     - [P0] Login with valid credentials → Dashboard loads
     - [P1] Display error for invalid credentials

   ### API Tests (P1-P2)

   - `tests/api/auth.api.spec.ts` (3 tests, 102 lines)
     - [P1] POST /auth/login - valid credentials → 200 + token
     - [P1] POST /auth/login - invalid credentials → 401 + error
     - [P2] POST /auth/login - missing fields → 400 + validation

   ### Component Tests (P1)

   - `tests/component/LoginForm.test.tsx` (2 tests, 45 lines)
     - [P1] Empty fields → submit button disabled
     - [P1] Valid input → submit button enabled

   ## Infrastructure Created

   ### Fixtures

   - `tests/support/fixtures/auth.fixture.ts` - authenticatedUser with auto-cleanup

   ### Factories

   - `tests/support/factories/user.factory.ts` - createUser(), deleteUser()

   ### Helpers

   - `tests/support/helpers/wait-for.ts` - Polling helper for complex conditions

   ## Test Execution

   ```bash
   # Run all new tests
   npm run test:e2e

   # Run by priority
   npm run test:e2e:p0  # Critical paths only
   npm run test:e2e:p1  # P0 + P1 tests
   ```
   ````

   ## Coverage Analysis

   **Total Tests:** 7
   - P0: 1 test (critical path)
   - P1: 5 tests (high priority)
   - P2: 1 test (medium priority)

   **Test Levels:**
   - E2E: 2 tests (user journeys)
   - API: 3 tests (business logic)
   - Component: 2 tests (UI behavior)

   **Coverage Status:**
   - ✅ All acceptance criteria covered
   - ✅ Happy path covered (E2E + API)
   - ✅ Error cases covered (API)
   - ✅ UI validation covered (Component)
   - ⚠️ Edge case: Password reset flow not yet covered (future story)

   ## Definition of Done
   - [x] All tests follow Given-When-Then format
   - [x] All tests use data-testid selectors
   - [x] All tests have priority tags
   - [x] All tests are self-cleaning (fixtures with auto-cleanup)
   - [x] No hard waits or flaky patterns
   - [x] Test files under 300 lines
   - [x] All tests run under 1.5 minutes each
   - [x] README updated with test execution instructions
   - [x] package.json scripts updated

   ## Next Steps
   1. Review generated tests with team
   2. Run tests in CI pipeline: `npm run test:e2e`
   3. Integrate with quality gate: `bmad tea *gate`
   4. Monitor for flaky tests in burn-in loop

   ````

   **Standalone Mode:**
   ```markdown
   # Automation Summary - {target_feature}

   **Date:** {date}
   **Target:** {target_feature} (standalone analysis)
   **Coverage Target:** {coverage_target}

   ## Feature Analysis

   **Source Files Analyzed:**
   - `src/auth/login.ts` - Login logic and validation
   - `src/auth/session.ts` - Session management
   - `src/auth/validation.ts` - Email/password validation

   **Existing Coverage:**
   - E2E tests: 0 found
   - API tests: 0 found
   - Component tests: 0 found
   - Unit tests: 0 found

   **Coverage Gaps Identified:**
   - ❌ No E2E tests for login flow
   - ❌ No API tests for /auth/login endpoint
   - ❌ No component tests for LoginForm
   - ❌ No unit tests for validateEmail()

   ## Tests Created

   {Same structure as BMad-Integrated Mode}

   ## Recommendations

   1. **High Priority (P0-P1):**
      - Add E2E test for password reset flow
      - Add API tests for token refresh endpoint
      - Add component tests for logout button

   2. **Medium Priority (P2):**
      - Add unit tests for session timeout logic
      - Add E2E test for "remember me" functionality

   3. **Future Enhancements:**
      - Consider contract testing for auth API
      - Add visual regression tests for login page
      - Set up burn-in loop for flaky test detection

   ## Definition of Done

   {Same checklist as BMad-Integrated Mode}
   ````

2. **Provide Summary to User**

   Output concise summary:

   ```markdown
   ## Automation Complete

   **Coverage:** {total_tests} tests created across {test_levels} levels
   **Priority Breakdown:** P0: {p0_count}, P1: {p1_count}, P2: {p2_count}, P3: {p3_count}
   **Infrastructure:** {fixture_count} fixtures, {factory_count} factories
   **Output:** {output_summary}

   **Run tests:** `npm run test:e2e`
   **Next steps:** Review tests, run in CI, integrate with quality gate
   ```

---

## Important Notes

### Dual-Mode Operation

**BMad-Integrated Mode** (story available):

- Uses story acceptance criteria for coverage targeting
- Aligns with test-design risk/priority assessment
- Expands ATDD tests with edge cases and negative paths
- Updates BMad status tracking

**Standalone Mode** (no story):

- Analyzes source code independently
- Identifies coverage gaps automatically
- Generates tests based on code analysis
- Works with any project (BMad or non-BMad)

**Auto-discover Mode** (no targets specified):

- Scans codebase for features needing tests
- Prioritizes features with no coverage
- Generates comprehensive test plan

### Avoid Duplicate Coverage

**Critical principle:** Don't test same behavior at multiple levels

**Good coverage:**

- E2E: User can login → Dashboard loads (critical happy path)
- API: POST /auth/login returns correct status codes (variations)
- Component: LoginForm validates input (UI edge cases)

**Bad coverage (duplicate):**

- E2E: User can login → Dashboard loads
- E2E: User can login with different emails → Dashboard loads (unnecessary duplication)
- API: POST /auth/login returns 200 (already covered in E2E)

Use E2E sparingly for critical paths. Use API/Component for variations and edge cases.

### Priority Tagging (Universal Pattern)

**Tag every test with priority in test name (ALL frameworks):**

```
// Pattern: '[P{N}] should {behavior}'
[TEST] '[P0] should login with valid credentials' { ... }
[TEST] '[P1] should display error for invalid credentials' { ... }
[TEST] '[P2] should remember login preference' { ... }
```

**Selective execution by framework:**

| Framework       | Run P0 Only             | Run P0+P1           |
| --------------- | ----------------------- | ------------------- |
| Playwright/Jest | `--grep "@P0"`          | `--grep "@P0\|@P1"` |
| pytest          | `-k "P0"`               | `-k "P0 or P1"`     |
| go test         | `-run "P0"`             | `-run "P0\|P1"`     |
| JUnit           | `@Tag("P0")` annotation | Tags annotation     |
| RSpec           | `--tag p0`              | `--tag p0 --tag p1` |

### No Page Objects (Recommendation)

**Prefer direct test code over page object abstractions.** Keep tests simple:

```
// ✅ CORRECT: Direct test (adapt syntax for detected framework)
[TEST] 'should login' {
  [NAVIGATE]('/login')
  [FILL]('[data-testid="email"]', 'user@example.com')
  [CLICK]('[data-testid="login-button"]')
  [ASSERT_URL]('/dashboard')
}

// ❌ AVOID: Page object abstraction adds indirection
class LoginPage {
  login(email, password) { ... }
}
```

Use fixtures for setup/teardown, not page objects for actions.

### Deterministic Tests Only (Universal)

**No flaky patterns allowed (applies to ALL frameworks):**

```
// ❌ WRONG: Hard wait (sleep/delay)
[WAIT_FIXED_TIME](2000)  // Never use this!

// ✅ CORRECT: Explicit wait for condition
[WAIT_FOR_ELEMENT]('[data-testid="user-name"]')
[ASSERT_VISIBLE]('[data-testid="user-name"]')

// ❌ WRONG: Conditional flow
if ([IS_VISIBLE](element)) {
  [CLICK](element)
}

// ✅ CORRECT: Deterministic assertion
[ASSERT_VISIBLE](element)
[CLICK](element)

// ❌ WRONG: Swallowing errors
try {
  [CLICK](element)
} catch (e) {
  // Test shouldn't catch errors
}

// ✅ CORRECT: Let test fail if element not found
[CLICK](element)  // Will fail clearly if not found
```

### Self-Cleaning Tests (Universal)

**Every test must clean up its data (adapt for detected framework):**

```
// ✅ CORRECT: Fixture/setup with auto-cleanup
[FIXTURE] testUser {
  [SETUP] {
    user = createUser()
    return user
  }
  [TEARDOWN] {
    deleteUser(user.id)  // Auto-cleanup
  }
}

// ❌ WRONG: Manual cleanup (can be forgotten)
[TEST] 'should login' {
  user = createUser()
  // ... test logic ...
  // Forgot to delete user! Test pollution!
}
```

### File Size Limits

**Keep test files lean (under {max_file_lines} lines):**

- If file exceeds limit, split into multiple files by feature area
- Group related tests in describe blocks
- Extract common setup to fixtures

### Knowledge Base Integration (Language-Agnostic Principles)

**Core Fragments (Auto-loaded in Step 1 - universal principles):**

- `test-levels-framework.md` - E2E vs API vs Component vs Unit decision framework
- `test-priorities-matrix.md` - P0-P3 classification with risk mapping
- `data-factories.md` - Factory patterns (adapt for detected faker library)
- `selective-testing.md` - Tag-based, spec filters, diff-based selection
- `ci-burn-in.md` - 10-iteration burn-in loop, sharding
- `test-quality.md` - Deterministic tests, isolation, cleanup, assertions

**Framework-Specific Patterns (fetch via web search):**

```
Search: "[DETECTED_FRAMEWORK] fixture patterns [CURRENT_YEAR]"
Search: "[DETECTED_FRAMEWORK] test factories [CURRENT_YEAR]"
Search: "[DETECTED_FRAMEWORK] network mocking [CURRENT_YEAR]"
```

**Healing Fragments (If `{auto_heal_failures}` enabled):**

- `test-healing-patterns.md` - Common failure patterns (adapt to detected framework)
- `selector-resilience.md` - Selector best practices (for E2E frameworks)
- `timing-debugging.md` - Race condition prevention

**Manual Reference (Optional):**

- Use `tea-index.csv` for additional specialized fragments

---

## Output Summary

After completing this workflow, provide a summary:

````markdown
## Automation Complete

**Mode:** {standalone_mode ? "Standalone" : "BMad-Integrated"}
**Target:** {story_id || target_feature || "Auto-discovered features"}
**Language:** {DETECTED_LANGUAGE}
**Framework:** {DETECTED_TEST_FRAMEWORK}
**Generated:** {TIMESTAMP}

**Tests Created:**

- E2E: {e2e_count} tests ({p0_count} P0, {p1_count} P1, {p2_count} P2)
- API: {api_count} tests ({p0_count} P0, {p1_count} P1, {p2_count} P2)
- Component: {component_count} tests ({p1_count} P1, {p2_count} P2)
- Unit: {unit_count} tests ({p2_count} P2, {p3_count} P3)

**Infrastructure:**

- Fixtures: {fixture_count} created/enhanced
- Factories: {factory_count} created/enhanced (using {FAKER_LIBRARY})
- Helpers: {helper_count} created/enhanced

**Documentation Updated:**

- ✅ Test README with execution instructions
- ✅ Test scripts configured for detected framework

**Test Execution (use detected commands):**

```bash
# Run all tests
{TEST_COMMAND}

# Run by priority (framework-specific filtering)
{TEST_COMMAND_P0}  # Critical paths only
{TEST_COMMAND_P1}  # P0 + P1 tests

# Run specific file
{TEST_COMMAND} {first_test_file}
```

**Documentation Fetched:**

- {DETECTED_FRAMEWORK} test patterns ({CURRENT_YEAR})
- {DETECTED_LANGUAGE} {FAKER_LIBRARY} factory examples
````

**Coverage Status:**

- ✅ {coverage_percentage}% of features covered
- ✅ All P0 scenarios covered
- ✅ All P1 scenarios covered
- ⚠️ {gap_count} coverage gaps identified (documented in summary)

**Quality Checks:**

- ✅ All tests follow Given-When-Then format
- ✅ All tests have priority tags
- ✅ All tests use stable selectors (data-testid, ARIA, semantic)
- ✅ All tests are self-cleaning (fixtures with auto-cleanup)
- ✅ No hard waits or flaky patterns
- ✅ All test files under {max_file_lines} lines

**Output File:** {output_summary}

**Next Steps:**

1. Review generated tests with team
2. Run tests in CI pipeline
3. Monitor for flaky tests in burn-in loop
4. Integrate with quality gate: `bmad tea *gate`

**Knowledge Base References Applied:**

- Test level selection framework (E2E vs API vs Component vs Unit)
- Priority classification (P0-P3)
- {DETECTED_FRAMEWORK} fixture patterns (fetched {CURRENT_DATE})
- {FAKER_LIBRARY} factory patterns (fetched {CURRENT_DATE})
- Selective testing strategies
- Test quality principles

```

---

## Validation

After completing all steps, verify:

- [ ] Execution mode determined (BMad-Integrated, Standalone, or Auto-discover)
- [ ] BMad artifacts loaded if available (story, tech-spec, test-design, PRD)
- [ ] Framework configuration loaded
- [ ] Existing test coverage analyzed (gaps identified)
- [ ] Knowledge base fragments loaded (test-levels, test-priorities, fixture-architecture, data-factories, selective-testing)
- [ ] Automation targets identified (what needs testing)
- [ ] Test levels selected appropriately (E2E, API, Component, Unit)
- [ ] Duplicate coverage avoided (same behavior not tested at multiple levels)
- [ ] Test priorities assigned (P0, P1, P2, P3)
- [ ] Fixture architecture created/enhanced (with auto-cleanup)
- [ ] Data factories created/enhanced (using faker)
- [ ] Helper utilities created/enhanced (if needed)
- [ ] E2E tests written (Given-When-Then, priority tags, data-testid selectors)
- [ ] API tests written (Given-When-Then, priority tags, comprehensive coverage)
- [ ] Component tests written (Given-When-Then, priority tags, UI behavior)
- [ ] Unit tests written (Given-When-Then, priority tags, pure logic)
- [ ] Network-first pattern applied (route interception before navigation)
- [ ] Quality standards enforced (no hard waits, no flaky patterns, self-cleaning, deterministic)
- [ ] Test README updated (execution instructions, priority tagging, patterns)
- [ ] package.json scripts updated (test execution commands)
- [ ] Test suite run locally (results captured)
- [ ] Tests validated (if auto_validate enabled)
- [ ] Failures healed (if auto_heal_failures enabled)
- [ ] Healing report generated (if healing attempted)
- [ ] Unfixable tests marked with test.fixme() (if any)
- [ ] Automation summary created (tests, infrastructure, coverage, healing, DoD)
- [ ] Output file formatted correctly

Refer to `checklist.md` for comprehensive validation criteria.
```
