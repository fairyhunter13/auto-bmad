---
description: 'Expand test automation coverage after implementation or analyze existing codebase to generate comprehensive test suite'
---

IT IS CRITICAL THAT YOU FOLLOW THESE STEPS - while staying in character as the current agent persona you may have loaded:

<scope-resolution CRITICAL="true" PRIORITY="HIGHEST">
## Step 0: Parse Scope BEFORE Anything Else (PARALLEL-SAFE)

You MUST parse scope from the command invocation FIRST. This enables parallel execution across multiple IDE windows.

### 0a. Parse Inline Scope from Command (HIGHEST PRIORITY)

Check if the user passed a scope with this command. Supported formats:

```
/testarch-automate --scope auth        # Flag format
/testarch-automate --scope=auth        # Flag with equals
/testarch-automate auth                # Positional (last arg if matches scope pattern)
```

**Scope pattern**: alphanumeric with hyphens, e.g., `auth`, `payment-v2`, `feature-123`

If inline scope found:

- Store as {scope} for this conversation
- Echo: "**[SCOPE: {scope}]** Workflow starting..."
- Skip to Step 0e

### 0b. Check Conversation Memory (PARALLEL-SAFE)

If no inline scope, check if scope was set earlier in THIS conversation:

- User may have run `/scope auth` or `/agent-pm auth` previously
- If you remember a scope from this conversation, use it
- Echo: "**[SCOPE: {scope}]** Using scope from conversation context"

### 0c. Fallback to .bmad-scope File (NOT PARALLEL-SAFE)

If no inline scope AND no conversation memory:

1. Check for `.bmad-scope` file in {project-root}
2. If exists, read `active_scope` value
3. **WARNING**: This file is shared across all sessions!
4. Echo: "**[SCOPE: {scope}]** Using scope from .bmad-scope (shared file - not parallel-safe)"

### 0d. No Scope Mode (Backward Compatible)

If no scope from any source:

- Continue without scope (artifacts go to root `_bmad-output/`)
- Echo: "**[NO SCOPE]** Running without scope isolation"

### 0e. Store Scope for Conversation

Once scope is resolved, REMEMBER IT for the rest of this conversation.

- Do NOT re-read `.bmad-scope` file later
- Inline scope from any command updates the conversation scope
  </scope-resolution>

<steps CRITICAL="TRUE">
1. Always LOAD the FULL @_bmad/core/tasks/workflow.xml
2. READ its entire contents - this is the CORE OS for EXECUTING the specific workflow-config @_bmad/bmm/workflows/testarch/automate/workflow.yaml
3. Pass the yaml path _bmad/bmm/workflows/testarch/automate/workflow.yaml as 'workflow-config' parameter to the workflow.xml instructions
4. **CRITICAL**: Pass {scope} (from Step 0) to workflow.xml - it needs this for path resolution
5. Follow workflow.xml instructions EXACTLY as written to process and follow the specific workflow config and its instructions
6. Save outputs after EACH section when generating any documents from templates
7. **Echo scope in every output**: Include "[SCOPE: {scope}]" prefix in status messages
</steps>

<scope-variables>
## Scope Path Overrides (when scope is active)

When {scope} is set, these variables OVERRIDE config.yaml values:

```
{scope_path} = {output_folder}/{scope}
{planning_artifacts} = {scope_path}/planning-artifacts
{implementation_artifacts} = {scope_path}/implementation-artifacts
{scope_tests} = {scope_path}/tests
```

**Example**: If scope="auth" and output_folder="\_bmad-output":

- {planning_artifacts} = `_bmad-output/auth/planning-artifacts`
- {implementation_artifacts} = `_bmad-output/auth/implementation-artifacts`
  </scope-variables>
