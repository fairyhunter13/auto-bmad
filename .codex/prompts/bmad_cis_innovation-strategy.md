---
description: 'Identify disruption opportunities and architect business model innovation. This workflow guides strategic analysis of markets, competitive dynamics, and business model innovation to uncover sustainable competitive advantages and breakthrough opportunities.'
---

IT IS CRITICAL THAT YOU FOLLOW THESE STEPS - while staying in character as the current agent persona you may have loaded:

<scope-resolution>
## Step 0: Resolve Scope (PARALLEL-SAFE)

Check for scope in priority order:

1. Inline flag: `/innovation-strategy --scope auth` or `/innovation-strategy auth`
2. Conversation memory (scope set earlier in this chat)
3. BMAD_SCOPE environment variable
4. .bmad-scope file (WARNING: shared across sessions)

If scope found, store as {scope} and echo: "**[SCOPE: {scope}]**"
If no scope, echo: "**[NO SCOPE]**" and continue (backward compatible)

When {scope} is set, override paths:

- {scope_path} = {output_folder}/{scope}
- {planning_artifacts} = {scope_path}/planning-artifacts
- {implementation_artifacts} = {scope_path}/implementation-artifacts
  </scope-resolution>

<steps CRITICAL="TRUE">
1. Always LOAD the FULL @_bmad/core/tasks/workflow.xml
2. READ its entire contents - this is the CORE OS for EXECUTING the specific workflow-config @_bmad/cis/workflows/innovation-strategy/workflow.yaml
3. Pass the yaml path _bmad/cis/workflows/innovation-strategy/workflow.yaml as 'workflow-config' parameter to the workflow.xml instructions
4. Pass {scope} (from Step 0) to workflow.xml for path resolution
5. Follow workflow.xml instructions EXACTLY as written to process and follow the specific workflow config and its instructions
6. Save outputs after EACH section when generating any documents from templates
</steps>
