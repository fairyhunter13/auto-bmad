---
description: 'Analyzes and documents brownfield projects by scanning codebase, architecture, and patterns to create comprehensive reference documentation for AI-assisted development'
---

IT IS CRITICAL THAT YOU FOLLOW THESE STEPS - while staying in character as the current agent persona you may have loaded:

<steps CRITICAL="TRUE">
1. Always LOAD the FULL @_bmad/core/tasks/workflow.xml
2. READ its entire contents - this is the CORE OS for EXECUTING the specific workflow-config @_bmad/bmm/workflows/document-project/workflow.yaml
3. Pass the yaml path _bmad/bmm/workflows/document-project/workflow.yaml as 'workflow-config' parameter to the workflow.xml instructions
4. Follow workflow.xml instructions EXACTLY as written to process and follow the specific workflow config and its instructions
5. Save outputs after EACH section when generating any documents from templates
</steps>

<scope-resolution>
## Multi-Scope Support

This workflow supports multi-scope parallel artifacts. Scope resolution order:

1. **--scope flag**: If provided (e.g., `/document-project --scope auth`), use that scope
2. **Session context**: Check for `.bmad-scope` file in project root
3. **Environment variable**: Check `BMAD_SCOPE` env var
4. **Prompt user**: If workflow requires scope and none found, prompt to select/create

When a scope is active:
- Artifacts are isolated to `_bmad-output/{scope}/`
- Cross-scope reads are allowed, writes are blocked
- Use `bmad scope sync-up` to promote artifacts to shared layer
- Check for pending dependency updates at workflow start
</scope-resolution>
