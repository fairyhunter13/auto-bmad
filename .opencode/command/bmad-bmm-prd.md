---
description: 'PRD tri-modal workflow - Create, Validate, or Edit comprehensive PRDs'
---

IT IS CRITICAL THAT YOU FOLLOW THESE STEPS IN ORDER:

<scope-resolution CRITICAL="true" PRIORITY="HIGHEST">
## Step 0: Parse Scope BEFORE Anything Else (PARALLEL-SAFE)

You MUST parse scope from the command invocation FIRST. This enables parallel execution across multiple IDE windows.

### 0a. Parse Inline Scope from Command (HIGHEST PRIORITY)

Check if the user passed a scope with this command. Supported formats:

```
/prd --scope auth        # Flag format
/prd --scope=auth        # Flag with equals
/prd auth                # Positional (last arg if matches scope pattern)
```

**Scope pattern**: alphanumeric with hyphens, e.g., `auth`, `payment-v2`, `feature-123`

If inline scope found:

- Store as {scope} for this conversation
- Echo: "**[SCOPE: {scope}]** Workflow starting..."
- Skip to Step 0d

### 0b. Check Conversation Memory (PARALLEL-SAFE)

If no inline scope, check if scope was set earlier in THIS conversation:

- User may have run `/scope auth` or `/agent-pm auth` previously
- If you remember a scope from this conversation, use it
- Echo: "**[SCOPE: {scope}]** Using scope from conversation context"
- Skip to Step 0d

### 0c. Fallback to .bmad-scope File (NOT PARALLEL-SAFE)

If no inline scope AND no conversation memory:

1. Check for `.bmad-scope` file in {project-root}
2. If exists, read `active_scope` value
3. **WARNING**: This file is shared across all sessions!
4. Echo: "**[SCOPE: {scope}]** Using scope from .bmad-scope (shared file - not parallel-safe)"
5. If `.bmad-scope` does not exist, continue without scope (backward compatible)

### 0d. Override Config Paths (CRITICAL - if scope is set)

After loading config.yaml but BEFORE using any paths, you MUST override these variables:

```
{scope_path} = {output_folder}/{scope}
{planning_artifacts} = {scope_path}/planning-artifacts
{implementation_artifacts} = {scope_path}/implementation-artifacts
{scope_tests} = {scope_path}/tests
```

**Example:** If config.yaml has `output_folder: "_bmad-output"` and scope is "auth":

- {scope_path} = `_bmad-output/auth`
- {planning_artifacts} = `_bmad-output/auth/planning-artifacts`
- {implementation_artifacts} = `_bmad-output/auth/implementation-artifacts`

**WARNING:** Config.yaml contains pre-resolved static paths. You MUST override them with the scope-aware paths above. DO NOT use the config.yaml values directly for these variables when a scope is active.

### 0e. Load Scope Context

If scope is set:

- Load global context: `{output_folder}/_shared/project-context.md`
- Load scope context if exists: `{scope_path}/project-context.md`
- Merge: scope-specific content extends/overrides global

### 0f. Store Scope for Conversation

Once scope is resolved, REMEMBER IT for the rest of this conversation.

- Do NOT re-read `.bmad-scope` file later
- Inline scope from any command updates the conversation scope
  </scope-resolution>

## Step 1: Execute Workflow

NOW: LOAD the FULL @_bmad/bmm/workflows/2-plan-workflows/prd/workflow.md, READ its entire contents and follow its directions exactly!

When the workflow instructs you to use `{planning_artifacts}` or `{implementation_artifacts}`, use YOUR OVERRIDDEN VALUES from Step 0d, not the static config.yaml values.

**Echo scope in every output**: Include "[SCOPE: {scope}]" prefix in status messages.
