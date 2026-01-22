# Pipedrive API v1 legacy endpoints

Generated from `openapi/derived/v1-legacy.yaml` by `cmd/endpoint-docs`. Do not edit manually.

Total operations: 194

| Method | Path | Summary | Operation ID |
| --- | --- | --- | --- |
| DELETE | `/activities` | Delete multiple activities in bulk | `deleteActivities` |
| GET | `/activities/collection` | Get all activities collection | `getActivitiesCollection` |
| GET | `/activityTypes` | Get all activity types | `getActivityTypes` |
| POST | `/activityTypes` | Add new activity type | `addActivityType` |
| DELETE | `/activityTypes` | Delete multiple activity types in bulk | `deleteActivityTypes` |
| PUT | `/activityTypes/{id}` | Update an activity type | `updateActivityType` |
| DELETE | `/activityTypes/{id}` | Delete an activity type | `deleteActivityType` |
| GET | `/billing/subscriptions/addons` | Get all add-ons for a single company | `getCompanyAddons` |
| GET | `/callLogs` | Get all call logs assigned to a particular user | `getUserCallLogs` |
| POST | `/callLogs` | Add a call log | `addCallLog` |
| GET | `/callLogs/{id}` | Get details of a call log | `getCallLog` |
| DELETE | `/callLogs/{id}` | Delete a call log | `deleteCallLog` |
| POST | `/callLogs/{id}/recordings` | Attach an audio file to the call log | `addCallLogAudioFile` |
| POST | `/channels` | Add a channel | `addChannel` |
| POST | `/channels/messages/receive` | Receives an incoming message | `receiveMessage` |
| DELETE | `/channels/{channel-id}/conversations/{conversation-id}` | Delete a conversation | `deleteConversation` |
| DELETE | `/channels/{id}` | Delete a channel | `deleteChannel` |
| GET | `/currencies` | Get all supported currencies | `getCurrencies` |
| DELETE | `/dealFields` | Delete multiple deal fields in bulk | `deleteDealFields` |
| DELETE | `/deals` | Delete multiple deals in bulk | `deleteDeals` |
| GET | `/deals/collection` | Get all deals collection | `getDealsCollection` |
| GET | `/deals/summary` | Get deals summary | `getDealsSummary` |
| GET | `/deals/summary/archived` | Get archived deals summary | `getArchivedDealsSummary` |
| GET | `/deals/timeline` | Get deals timeline | `getDealsTimeline` |
| GET | `/deals/timeline/archived` | Get archived deals timeline | `getArchivedDealsTimeline` |
| GET | `/deals/{id}/activities` | List activities associated with a deal | `getDealActivities` |
| GET | `/deals/{id}/changelog` | List updates about deal field values | `getDealChangelog` |
| POST | `/deals/{id}/duplicate` | Duplicate deal | `duplicateDeal` |
| GET | `/deals/{id}/files` | List files attached to a deal | `getDealFiles` |
| GET | `/deals/{id}/flow` | List updates about a deal | `getDealUpdates` |
| GET | `/deals/{id}/mailMessages` | List mail messages associated with a deal | `getDealMailMessages` |
| PUT | `/deals/{id}/merge` | Merge two deals | `mergeDeals` |
| GET | `/deals/{id}/participants` | List participants of a deal | `getDealParticipants` |
| POST | `/deals/{id}/participants` | Add a participant to a deal | `addDealParticipant` |
| DELETE | `/deals/{id}/participants/{deal_participant_id}` | Delete a participant from a deal | `deleteDealParticipant` |
| GET | `/deals/{id}/participantsChangelog` | List updates about participants of a deal | `getDealParticipantsChangelog` |
| GET | `/deals/{id}/permittedUsers` | List permitted users | `getDealUsers` |
| GET | `/deals/{id}/persons` | List all persons associated with a deal | `getDealPersons` |
| GET | `/files` | Get all files | `getFiles` |
| POST | `/files` | Add file | `addFile` |
| POST | `/files/remote` | Create a remote file and link it to an item | `addFileAndLinkIt` |
| POST | `/files/remoteLink` | Link a remote file to an item | `linkFileToItem` |
| GET | `/files/{id}` | Get one file | `getFile` |
| PUT | `/files/{id}` | Update file details | `updateFile` |
| DELETE | `/files/{id}` | Delete a file | `deleteFile` |
| GET | `/files/{id}/download` | Download one file | `downloadFile` |
| GET | `/filters` | Get all filters | `getFilters` |
| POST | `/filters` | Add a new filter | `addFilter` |
| DELETE | `/filters` | Delete multiple filters in bulk | `deleteFilters` |
| GET | `/filters/helpers` | Get all filter helpers | `getFilterHelpers` |
| GET | `/filters/{id}` | Get one filter | `getFilter` |
| PUT | `/filters/{id}` | Update filter | `updateFilter` |
| DELETE | `/filters/{id}` | Delete a filter | `deleteFilter` |
| POST | `/goals` | Add a new goal | `addGoal` |
| GET | `/goals/find` | Find goals | `getGoals` |
| PUT | `/goals/{id}` | Update existing goal | `updateGoal` |
| DELETE | `/goals/{id}` | Delete existing goal | `deleteGoal` |
| GET | `/goals/{id}/results` | Get result of a goal | `getGoalResult` |
| GET | `/leadFields` | Get all lead fields | `getLeadFields` |
| GET | `/leadLabels` | Get all lead labels | `getLeadLabels` |
| POST | `/leadLabels` | Add a lead label | `addLeadLabel` |
| PATCH | `/leadLabels/{id}` | Update a lead label | `updateLeadLabel` |
| DELETE | `/leadLabels/{id}` | Delete a lead label | `deleteLeadLabel` |
| GET | `/leadSources` | Get all lead sources | `getLeadSources` |
| GET | `/leads` | Get all leads | `getLeads` |
| POST | `/leads` | Add a lead | `addLead` |
| GET | `/leads/archived` | Get all archived leads | `getArchivedLeads` |
| GET | `/leads/{id}` | Get one lead | `getLead` |
| PATCH | `/leads/{id}` | Update a lead | `updateLead` |
| DELETE | `/leads/{id}` | Delete a lead | `deleteLead` |
| GET | `/leads/{id}/permittedUsers` | List permitted users | `getLeadUsers` |
| GET | `/legacyTeams` | Get all teams | `getTeams` |
| POST | `/legacyTeams` | Add a new team | `addTeam` |
| GET | `/legacyTeams/user/{id}` | Get all teams of a user | `getUserTeams` |
| GET | `/legacyTeams/{id}` | Get a single team | `getTeam` |
| PUT | `/legacyTeams/{id}` | Update a team | `updateTeam` |
| GET | `/legacyTeams/{id}/users` | Get all users in a team | `getTeamUsers` |
| POST | `/legacyTeams/{id}/users` | Add users to a team | `addTeamUser` |
| DELETE | `/legacyTeams/{id}/users` | Delete users from a team | `deleteTeamUser` |
| GET | `/mailbox/mailMessages/{id}` | Get one mail message | `getMailMessage` |
| GET | `/mailbox/mailThreads` | Get mail threads | `getMailThreads` |
| GET | `/mailbox/mailThreads/{id}` | Get one mail thread | `getMailThread` |
| PUT | `/mailbox/mailThreads/{id}` | Update mail thread details | `updateMailThreadDetails` |
| DELETE | `/mailbox/mailThreads/{id}` | Delete mail thread | `deleteMailThread` |
| GET | `/mailbox/mailThreads/{id}/mailMessages` | Get all mail messages of mail thread | `getMailThreadMessages` |
| POST | `/meetings/userProviderLinks` | Link a user with the installed video call integration | `saveUserProviderLink` |
| DELETE | `/meetings/userProviderLinks/{id}` | Delete the link between a user and the installed video call integration | `deleteUserProviderLink` |
| GET | `/noteFields` | Get all note fields | `getNoteFields` |
| GET | `/notes` | Get all notes | `getNotes` |
| POST | `/notes` | Add a note | `addNote` |
| GET | `/notes/{id}` | Get one note | `getNote` |
| PUT | `/notes/{id}` | Update a note | `updateNote` |
| DELETE | `/notes/{id}` | Delete a note | `deleteNote` |
| GET | `/notes/{id}/comments` | Get all comments for a note | `getNoteComments` |
| POST | `/notes/{id}/comments` | Add a comment to a note | `addNoteComment` |
| GET | `/notes/{id}/comments/{commentId}` | Get one comment | `getComment` |
| PUT | `/notes/{id}/comments/{commentId}` | Update a comment related to a note | `updateCommentForNote` |
| DELETE | `/notes/{id}/comments/{commentId}` | Delete a comment related to a note | `deleteComment` |
| GET | `/oauth/authorize` | Requesting authorization | `authorize` |
| POST | `/oauth/token` | Getting the tokens | `get-tokens` |
| POST | `/oauth/token/` | Refreshing the tokens | `refresh-tokens` |
| DELETE | `/organizationFields` | Delete multiple organization fields in bulk | `deleteOrganizationFields` |
| GET | `/organizationRelationships` | Get all relationships for organization | `getOrganizationRelationships` |
| POST | `/organizationRelationships` | Create an organization relationship | `addOrganizationRelationship` |
| GET | `/organizationRelationships/{id}` | Get one organization relationship | `getOrganizationRelationship` |
| PUT | `/organizationRelationships/{id}` | Update an organization relationship | `updateOrganizationRelationship` |
| DELETE | `/organizationRelationships/{id}` | Delete an organization relationship | `deleteOrganizationRelationship` |
| DELETE | `/organizations` | Delete multiple organizations in bulk | `deleteOrganizations` |
| GET | `/organizations/collection` | Get all organizations collection | `getOrganizationsCollection` |
| GET | `/organizations/{id}/activities` | List activities associated with an organization | `getOrganizationActivities` |
| GET | `/organizations/{id}/changelog` | List updates about organization field values | `getOrganizationChangelog` |
| GET | `/organizations/{id}/deals` | List deals associated with an organization | `getOrganizationDeals` |
| GET | `/organizations/{id}/files` | List files attached to an organization | `getOrganizationFiles` |
| GET | `/organizations/{id}/flow` | List updates about an organization | `getOrganizationUpdates` |
| GET | `/organizations/{id}/mailMessages` | List mail messages associated with an organization | `getOrganizationMailMessages` |
| PUT | `/organizations/{id}/merge` | Merge two organizations | `mergeOrganizations` |
| GET | `/organizations/{id}/permittedUsers` | List permitted users | `getOrganizationUsers` |
| GET | `/organizations/{id}/persons` | List persons of an organization | `getOrganizationPersons` |
| GET | `/permissionSets` | Get all permission sets | `getPermissionSets` |
| GET | `/permissionSets/{id}` | Get one permission set | `getPermissionSet` |
| GET | `/permissionSets/{id}/assignments` | List permission set assignments | `getPermissionSetAssignments` |
| DELETE | `/personFields` | Delete multiple person fields in bulk | `deletePersonFields` |
| DELETE | `/persons` | Delete multiple persons in bulk | `deletePersons` |
| GET | `/persons/collection` | Get all persons collection | `getPersonsCollection` |
| GET | `/persons/{id}/activities` | List activities associated with a person | `getPersonActivities` |
| GET | `/persons/{id}/changelog` | List updates about person field values | `getPersonChangelog` |
| GET | `/persons/{id}/deals` | List deals associated with a person | `getPersonDeals` |
| GET | `/persons/{id}/files` | List files attached to a person | `getPersonFiles` |
| GET | `/persons/{id}/flow` | List updates about a person | `getPersonUpdates` |
| GET | `/persons/{id}/mailMessages` | List mail messages associated with a person | `getPersonMailMessages` |
| PUT | `/persons/{id}/merge` | Merge two persons | `mergePersons` |
| GET | `/persons/{id}/permittedUsers` | List permitted users | `getPersonUsers` |
| POST | `/persons/{id}/picture` | Add person picture | `addPersonPicture` |
| DELETE | `/persons/{id}/picture` | Delete person picture | `deletePersonPicture` |
| GET | `/persons/{id}/products` | List products associated with a person | `getPersonProducts` |
| GET | `/pipelines/{id}/conversion_statistics` | Get deals conversion rates in pipeline | `getPipelineConversionStatistics` |
| GET | `/pipelines/{id}/deals` | Get deals in a pipeline | `getPipelineDeals` |
| GET | `/pipelines/{id}/movement_statistics` | Get deals movements in pipeline | `getPipelineMovementStatistics` |
| DELETE | `/productFields` | Delete multiple product fields in bulk | `deleteProductFields` |
| GET | `/products/{id}/deals` | Get deals where a product is attached to | `getProductDeals` |
| GET | `/products/{id}/files` | List files attached to a product | `getProductFiles` |
| GET | `/products/{id}/permittedUsers` | List permitted users | `getProductUsers` |
| GET | `/projectTemplates` | Get all project templates | `getProjectTemplates` |
| GET | `/projectTemplates/{id}` | Get details of a template | `getProjectTemplate` |
| GET | `/projects` | Get all projects | `getProjects` |
| POST | `/projects` | Add a project | `addProject` |
| GET | `/projects/boards` | Get all project boards | `getProjectsBoards` |
| GET | `/projects/boards/{id}` | Get details of a board | `getProjectsBoard` |
| GET | `/projects/phases` | Get project phases | `getProjectsPhases` |
| GET | `/projects/phases/{id}` | Get details of a phase | `getProjectsPhase` |
| GET | `/projects/{id}` | Get details of a project | `getProject` |
| PUT | `/projects/{id}` | Update a project | `updateProject` |
| DELETE | `/projects/{id}` | Delete a project | `deleteProject` |
| GET | `/projects/{id}/activities` | Returns project activities | `getProjectActivities` |
| POST | `/projects/{id}/archive` | Archive a project | `archiveProject` |
| GET | `/projects/{id}/groups` | Returns project groups | `getProjectGroups` |
| GET | `/projects/{id}/plan` | Returns project plan | `getProjectPlan` |
| PUT | `/projects/{id}/plan/activities/{activityId}` | Update activity in project plan | `putProjectPlanActivity` |
| PUT | `/projects/{id}/plan/tasks/{taskId}` | Update task in project plan | `putProjectPlanTask` |
| GET | `/projects/{id}/tasks` | Returns project tasks | `getProjectTasks` |
| GET | `/recents` | Get recents | `getRecents` |
| GET | `/roles` | Get all roles | `getRoles` |
| POST | `/roles` | Add a role | `addRole` |
| GET | `/roles/{id}` | Get one role | `getRole` |
| PUT | `/roles/{id}` | Update role details | `updateRole` |
| DELETE | `/roles/{id}` | Delete a role | `deleteRole` |
| GET | `/roles/{id}/assignments` | List role assignments | `getRoleAssignments` |
| POST | `/roles/{id}/assignments` | Add role assignment | `addRoleAssignment` |
| DELETE | `/roles/{id}/assignments` | Delete a role assignment | `deleteRoleAssignment` |
| GET | `/roles/{id}/pipelines` | List pipeline visibility for a role | `getRolePipelines` |
| PUT | `/roles/{id}/pipelines` | Update pipeline visibility for a role | `updateRolePipelines` |
| GET | `/roles/{id}/settings` | List role settings | `getRoleSettings` |
| POST | `/roles/{id}/settings` | Add or update role setting | `addOrUpdateRoleSetting` |
| DELETE | `/stages` | Delete multiple stages in bulk | `deleteStages` |
| GET | `/stages/{id}/deals` | Get deals in a stage | `getStageDeals` |
| GET | `/tasks` | Get all tasks | `getTasks` |
| POST | `/tasks` | Add a task | `addTask` |
| GET | `/tasks/{id}` | Get details of a task | `getTask` |
| PUT | `/tasks/{id}` | Update a task | `updateTask` |
| DELETE | `/tasks/{id}` | Delete a task | `deleteTask` |
| GET | `/userConnections` | Get all user connections | `getUserConnections` |
| GET | `/userSettings` | List settings of an authorized user | `getUserSettings` |
| GET | `/users` | Get all users | `getUsers` |
| POST | `/users` | Add a new user | `addUser` |
| GET | `/users/find` | Find users by name | `findUsersByName` |
| GET | `/users/me` | Get current user data | `getCurrentUser` |
| GET | `/users/{id}` | Get one user | `getUser` |
| PUT | `/users/{id}` | Update user details | `updateUser` |
| GET | `/users/{id}/permissions` | List user permissions | `getUserPermissions` |
| GET | `/users/{id}/roleAssignments` | List role assignments | `getUserRoleAssignments` |
| GET | `/users/{id}/roleSettings` | List user role settings | `getUserRoleSettings` |
| GET | `/webhooks` | Get all Webhooks | `getWebhooks` |
| POST | `/webhooks` | Create a new Webhook | `addWebhook` |
| DELETE | `/webhooks/{id}` | Delete existing Webhook | `deleteWebhook` |

