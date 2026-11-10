# Migrating to SDK v2.0.0

SDK v2.0.0 introduces the changes below. Existing v1.17.x consumers keep their
current dependency until they change module and import paths. Under the existing
[security policy](../SECURITY.md), only the latest release receives security
updates, so v1 releases are no longer supported once v2.0.0 ships.

## Update the module and imports

Install the new major version:

```sh
go get github.com/juhokoskela/pipedrive-go/v2@v2.0.0
```

Insert `/v2` after `pipedrive-go` in imports:

```go
import (
	"github.com/juhokoskela/pipedrive-go/v2/pipedrive"
	v1 "github.com/juhokoskela/pipedrive-go/v2/pipedrive/v1"
	v2 "github.com/juhokoskela/pipedrive-go/v2/pipedrive/v2"
)
```

The first `/v2` selects the SDK major version. The final `/v1` or `/v2` selects
the Pipedrive API. After updating imports and the calls listed below, run
`go mod tidy` to remove the old module if nothing still uses it.

## Use v2 for overlapping operations

The following 17 v1 methods have been removed. A deprecation comment alone would
leave both versions callable, so these methods and their unused supporting types
are removed together. Create a `v2.Client` with the same authentication and use
these replacements:

| Removed v1 method | v2 replacement |
| --- | --- |
| `Projects.List` | `Projects.List` |
| `Projects.Create` | `Projects.Create` |
| `Projects.Get` | `Projects.Get` |
| `Projects.Update` | `Projects.Update` |
| `Projects.Delete` | `Projects.Delete` |
| `Projects.Archive` | `Projects.Archive` |
| `Projects.ListBoards` | `ProjectBoards.List` |
| `Projects.GetBoard` | `ProjectBoards.Get` |
| `Projects.ListPhases` | `ProjectPhases.List` |
| `Projects.GetPhase` | `ProjectPhases.Get` |
| `Tasks.List` | `Tasks.List` |
| `Tasks.Get` | `Tasks.Get` |
| `Tasks.Create` | `Tasks.Create` |
| `Tasks.Update` | `Tasks.Update` |
| `Tasks.Delete` | `Tasks.Delete` |
| `ProjectTemplates.List` | `ProjectTemplates.List` |
| `ProjectTemplates.Get` | `ProjectTemplates.Get` |

Use the ID types, models, and functional options from `pipedrive/v2`. For example,
replace a task update payload map with `v2.WithTaskTitle("New title")` passed to
`client.Tasks.Update(ctx, v2.TaskID(id), ...)`. Project and task updates use PATCH
in v2. List methods for projects and tasks return a cursor instead of offset
pagination. `ProjectPhases.List` requires a `v2.ProjectBoardID` argument.

The v1 `Tasks` and `ProjectTemplates` client fields, services, and options are
removed. So are the v1 `Project`, `ProjectBoard`, `ProjectPhase`, `ProjectTemplate`,
and `Task` models and `ProjectBoardID`, `ProjectPhaseID`, and `ProjectTemplateID`
types. `v1.ProjectID` and `v1.TaskID` remain for legacy associations.

The v1 `Projects` service still provides `ListActivities`, `ListGroups`, `GetPlan`,
`UpdatePlanActivity`, `UpdatePlanTask`, and `ListTasks`. These operations have no
v2 equivalents in the checked-in specification. The raw client remains available
for direct API calls.

## Corrected response types

| Method | Previous return data | New return data |
| --- | --- | --- |
| `v1.Roles.ListPipelines` | `[]map[string]any` | `map[string]any` |
| `v1.Roles.ListSettings` | `[]map[string]any` | `map[string]any` |
| `v1.Users.ListRoleSettings` | `[]map[string]any` | `map[string]any` |
| `v1.Projects.GetPlan` | `map[string]any` | `[]map[string]any` |

Read pipeline visibility from `result["pipeline_ids"]` and `result["visible"]`.
Read settings by their keys, such as `result["deal_default_visibility"]`. Iterate
project plan items and read their `item_id`, `item_type`, `phase_id`, and `group_id`.
Values inside these maps retain the standard `encoding/json` representations,
including `float64` for numbers.
