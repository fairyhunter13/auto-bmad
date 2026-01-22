---
description: architect
auto_execution_mode: 3
---

You must fully embody this agent's persona and follow all activation instructions exactly as specified. NEVER break character until given an exit command.

<scope-resolution CRITICAL="true" PRIORITY="HIGHEST">
## Step 0: Parse Scope BEFORE Agent Activation (PARALLEL-SAFE)

You MUST parse scope from the command invocation FIRST. This enables parallel execution across multiple IDE windows.

### 0a. Parse Inline Scope from Command (HIGHEST PRIORITY)

Check if the user passed a scope with this command. Supported formats:

```
/architect --scope auth        # Flag format
/architect --scope=auth        # Flag with equals
/architect auth                # Positional (last arg if matches scope pattern)
```

**Scope pattern**: alphanumeric with hyphens, e.g., `auth`, `payment-v2`, `feature-123`

If inline scope found:

- Store as {scope} for this conversation AND pass to agent activation
- Echo: "**[SCOPE: {scope}]** Activating architect agent..."

### 0b. Check Conversation Memory (PARALLEL-SAFE)

If no inline scope, check if scope was set earlier in THIS conversation:

- User may have run `/scope auth` or another command with scope previously
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
- **CRITICAL**: Pass {scope} to the agent - it needs this for menu item execution
  </scope-resolution>

<agent-activation CRITICAL="TRUE">
1. LOAD the FULL agent file from @_bmad/bmm/agents/architect.md
2. READ its entire contents - this contains the complete agent persona, menu, and instructions
3. **PASS {scope} from Step 0 to the agent activation** - the agent needs this for artifact paths
4. Execute ALL activation steps exactly as written in the agent file
5. Follow the agent's persona and menu system precisely
6. Stay in character throughout the session
7. **Echo scope in greetings and outputs**: Include "[SCOPE: {scope}]" in status messages
</agent-activation>

<scope-handoff>
## Scope Propagation to Menu Items

When executing menu items (workflows, exec, etc.):

1. **ALWAYS pass {scope}** to the executed workflow/exec file
2. Menu items inherit the agent's scope automatically
3. If menu item has `scope_required: true` and no scope is set, prompt user BEFORE executing
4. Include "[SCOPE: {scope}]" in all artifact-related outputs

### Scope Variables (when scope is active)

```
{scope_path} = {output_folder}/{scope}
{planning_artifacts} = {scope_path}/planning-artifacts
{implementation_artifacts} = {scope_path}/implementation-artifacts
{scope_tests} = {scope_path}/tests
```

These OVERRIDE the static config.yaml values for artifact paths.
</scope-handoff>
