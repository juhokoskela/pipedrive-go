package v1

import (
	"net/url"
	"testing"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

func requireOneV1RequestOption(t *testing.T, name string, opts []pipedrive.RequestOption) {
	t.Helper()

	if len(opts) != 1 {
		t.Fatalf("%s request options count = %d, want 1", name, len(opts))
	}
}

func TestUsersRequestOptionsApplyAllTargets(t *testing.T) {
	t.Parallel()

	opt := WithUsersRequestOptions(pipedrive.WithHeader("X-Test", "1"))

	var list listUsersOptions
	opt.applyListUsers(&list)
	requireOneV1RequestOption(t, "list users", list.requestOptions)

	var get getUserOptions
	opt.applyGetUser(&get)
	requireOneV1RequestOption(t, "get user", get.requestOptions)

	var current getCurrentUserOptions
	opt.applyGetCurrentUser(&current)
	requireOneV1RequestOption(t, "current user", current.requestOptions)

	var permissions getUserPermissionsOptions
	opt.applyGetUserPermissions(&permissions)
	requireOneV1RequestOption(t, "user permissions", permissions.requestOptions)

	_ = newListUsersOptions([]ListUsersOption{nil})
	_ = newGetUserOptions([]GetUserOption{nil})
	_ = newGetCurrentUserOptions([]GetCurrentUserOption{nil})
	_ = newGetUserPermissionsOptions([]GetUserPermissionsOption{nil})
}

func TestFiltersRequestOptionsApplyAllTargets(t *testing.T) {
	t.Parallel()

	opt := WithFiltersRequestOptions(pipedrive.WithHeader("X-Test", "1"))

	var list listFiltersOptions
	opt.applyListFilters(&list)
	requireOneV1RequestOption(t, "list filters", list.requestOptions)

	var get getFilterOptions
	opt.applyGetFilter(&get)
	requireOneV1RequestOption(t, "get filter", get.requestOptions)

	var create createFilterOptions
	opt.applyCreateFilter(&create)
	requireOneV1RequestOption(t, "create filter", create.requestOptions)

	var update updateFilterOptions
	opt.applyUpdateFilter(&update)
	requireOneV1RequestOption(t, "update filter", update.requestOptions)

	var deleteFilter deleteFilterOptions
	opt.applyDeleteFilter(&deleteFilter)
	requireOneV1RequestOption(t, "delete filter", deleteFilter.requestOptions)

	var deleteFilters deleteFiltersOptions
	opt.applyDeleteFilters(&deleteFilters)
	requireOneV1RequestOption(t, "delete filters", deleteFilters.requestOptions)

	var helpers listFilterHelpersOptions
	opt.applyListFilterHelpers(&helpers)
	requireOneV1RequestOption(t, "filter helpers", helpers.requestOptions)

	create = newCreateFilterOptions([]CreateFilterOption{nil, WithFilterConditions(nil)})
	if create.payload.conditions != nil {
		t.Fatalf("expected nil filter conditions, got %#v", create.payload.conditions)
	}

	_ = newListFiltersOptions([]ListFiltersOption{nil})
	_ = newGetFilterOptions([]GetFilterOption{nil})
	_ = newUpdateFilterOptions([]UpdateFilterOption{nil})
	_ = newDeleteFilterOptions([]DeleteFilterOption{nil})
	_ = newDeleteFiltersOptions([]DeleteFiltersOption{nil})
	_ = newListFilterHelpersOptions([]ListFilterHelpersOption{nil})
}

func TestNotesRequestOptionsApplyAllTargets(t *testing.T) {
	t.Parallel()

	opt := WithNotesRequestOptions(pipedrive.WithHeader("X-Test", "1"))

	var list listNotesOptions
	opt.applyListNotes(&list)
	requireOneV1RequestOption(t, "list notes", list.requestOptions)

	var get getNoteOptions
	opt.applyGetNote(&get)
	requireOneV1RequestOption(t, "get note", get.requestOptions)

	var create createNoteOptions
	opt.applyCreateNote(&create)
	requireOneV1RequestOption(t, "create note", create.requestOptions)

	var update updateNoteOptions
	opt.applyUpdateNote(&update)
	requireOneV1RequestOption(t, "update note", update.requestOptions)

	var deleteNote deleteNoteOptions
	opt.applyDeleteNote(&deleteNote)
	requireOneV1RequestOption(t, "delete note", deleteNote.requestOptions)

	var comments listNoteCommentsOptions
	opt.applyListNoteComments(&comments)
	requireOneV1RequestOption(t, "list note comments", comments.requestOptions)

	var createComment createNoteCommentOptions
	opt.applyCreateNoteComment(&createComment)
	requireOneV1RequestOption(t, "create note comment", createComment.requestOptions)

	var getComment getNoteCommentOptions
	opt.applyGetNoteComment(&getComment)
	requireOneV1RequestOption(t, "get note comment", getComment.requestOptions)

	var updateComment updateNoteCommentOptions
	opt.applyUpdateNoteComment(&updateComment)
	requireOneV1RequestOption(t, "update note comment", updateComment.requestOptions)

	var deleteComment deleteNoteCommentOptions
	opt.applyDeleteNoteComment(&deleteComment)
	requireOneV1RequestOption(t, "delete note comment", deleteComment.requestOptions)

	_ = newListNotesOptions([]ListNotesOption{nil})
	_ = newGetNoteOptions([]GetNoteOption{nil})
	_ = newCreateNoteOptions([]CreateNoteOption{nil})
	_ = newUpdateNoteOptions([]UpdateNoteOption{nil})
	_ = newDeleteNoteOptions([]DeleteNoteOption{nil})
	_ = newListNoteCommentsOptions([]ListNoteCommentsOption{nil})
	_ = newCreateNoteCommentOptions([]CreateNoteCommentOption{nil})
	_ = newGetNoteCommentOptions([]GetNoteCommentOption{nil})
	_ = newUpdateNoteCommentOptions([]UpdateNoteCommentOption{nil})
	_ = newDeleteNoteCommentOptions([]DeleteNoteCommentOption{nil})
}

func TestActivityTypesRequestOptionsApplyAllTargets(t *testing.T) {
	t.Parallel()

	opt := WithActivityTypesRequestOptions(pipedrive.WithHeader("X-Test", "1"))

	var list listActivityTypesOptions
	opt.applyListActivityTypes(&list)
	requireOneV1RequestOption(t, "list activity types", list.requestOptions)

	var create createActivityTypeOptions
	opt.applyCreateActivityType(&create)
	requireOneV1RequestOption(t, "create activity type", create.requestOptions)

	var update updateActivityTypeOptions
	opt.applyUpdateActivityType(&update)
	requireOneV1RequestOption(t, "update activity type", update.requestOptions)

	var deleteActivityType deleteActivityTypeOptions
	opt.applyDeleteActivityType(&deleteActivityType)
	requireOneV1RequestOption(t, "delete activity type", deleteActivityType.requestOptions)

	_ = newListActivityTypesOptions([]ListActivityTypesOption{nil})
	_ = newCreateActivityTypeOptions([]CreateActivityTypeOption{nil})
	_ = newUpdateActivityTypeOptions([]UpdateActivityTypeOption{nil})
	_ = newDeleteActivityTypeOptions([]DeleteActivityTypeOption{nil})
}

func TestLeadLabelsRequestOptionsApplyAllTargets(t *testing.T) {
	t.Parallel()

	opt := WithLeadLabelsRequestOptions(pipedrive.WithHeader("X-Test", "1"))

	var list listLeadLabelsOptions
	opt.applyListLeadLabels(&list)
	requireOneV1RequestOption(t, "list lead labels", list.requestOptions)

	var create createLeadLabelOptions
	opt.applyCreateLeadLabel(&create)
	requireOneV1RequestOption(t, "create lead label", create.requestOptions)

	var update updateLeadLabelOptions
	opt.applyUpdateLeadLabel(&update)
	requireOneV1RequestOption(t, "update lead label", update.requestOptions)

	var deleteLeadLabel deleteLeadLabelOptions
	opt.applyDeleteLeadLabel(&deleteLeadLabel)
	requireOneV1RequestOption(t, "delete lead label", deleteLeadLabel.requestOptions)

	_ = newListLeadLabelsOptions([]ListLeadLabelsOption{nil})
	_ = newCreateLeadLabelOptions([]CreateLeadLabelOption{nil})
	_ = newUpdateLeadLabelOptions([]UpdateLeadLabelOption{nil})
	_ = newDeleteLeadLabelOptions([]DeleteLeadLabelOption{nil})
}

func TestWebhooksRequestOptionsApplyAllTargets(t *testing.T) {
	t.Parallel()

	opt := WithWebhooksRequestOptions(pipedrive.WithHeader("X-Test", "1"))

	var list listWebhooksOptions
	opt.applyListWebhooks(&list)
	requireOneV1RequestOption(t, "list webhooks", list.requestOptions)

	var create createWebhookOptions
	opt.applyCreateWebhook(&create)
	requireOneV1RequestOption(t, "create webhook", create.requestOptions)

	var deleteWebhook deleteWebhookOptions
	opt.applyDeleteWebhook(&deleteWebhook)
	requireOneV1RequestOption(t, "delete webhook", deleteWebhook.requestOptions)

	_ = newListWebhooksOptions([]ListWebhooksOption{nil})
	_ = newCreateWebhookOptions([]CreateWebhookOption{nil})
	_ = newDeleteWebhookOptions([]DeleteWebhookOption{nil})
}

func TestOrganizationRelationshipsRequestOptionsApplyAllTargets(t *testing.T) {
	t.Parallel()

	opt := WithOrganizationRelationshipsRequestOptions(pipedrive.WithHeader("X-Test", "1"))

	var list listOrganizationRelationshipsOptions
	opt.applyListOrganizationRelationships(&list)
	requireOneV1RequestOption(t, "list organization relationships", list.requestOptions)

	var get getOrganizationRelationshipOptions
	opt.applyGetOrganizationRelationship(&get)
	requireOneV1RequestOption(t, "get organization relationship", get.requestOptions)

	var create createOrganizationRelationshipOptions
	opt.applyCreateOrganizationRelationship(&create)
	requireOneV1RequestOption(t, "create organization relationship", create.requestOptions)

	var update updateOrganizationRelationshipOptions
	opt.applyUpdateOrganizationRelationship(&update)
	requireOneV1RequestOption(t, "update organization relationship", update.requestOptions)

	var deleteRelationship deleteOrganizationRelationshipOptions
	opt.applyDeleteOrganizationRelationship(&deleteRelationship)
	requireOneV1RequestOption(t, "delete organization relationship", deleteRelationship.requestOptions)

	_ = newListOrganizationRelationshipsOptions([]ListOrganizationRelationshipsOption{nil})
	_ = newGetOrganizationRelationshipOptions([]GetOrganizationRelationshipOption{nil})
	_ = newCreateOrganizationRelationshipOptions([]CreateOrganizationRelationshipOption{nil})
	_ = newUpdateOrganizationRelationshipOptions([]UpdateOrganizationRelationshipOption{nil})
	_ = newDeleteOrganizationRelationshipOptions([]DeleteOrganizationRelationshipOption{nil})
}

func TestFilesOptionsApplyQueryAndRequestOptions(t *testing.T) {
	t.Parallel()

	query := url.Values{"deal_id": []string{"7"}}
	cfg := newFilesOptions([]FilesOption{
		nil,
		WithFilesQuery(query),
		WithFilesRequestOptions(pipedrive.WithHeader("X-Test", "1")),
	})

	if got := cfg.query.Get("deal_id"); got != "7" {
		t.Fatalf("unexpected deal_id query: %q", got)
	}
	requireOneV1RequestOption(t, "files", cfg.requestOptions)
}

func TestV1OptionFuncsInvokeCallbacks(t *testing.T) {
	t.Parallel()

	called := false
	listUsersOptionFunc(func(*listUsersOptions) { called = true }).applyListUsers(&listUsersOptions{})
	if !called {
		t.Fatal("expected list users callback")
	}

	called = false
	getUserOptionFunc(func(*getUserOptions) { called = true }).applyGetUser(&getUserOptions{})
	if !called {
		t.Fatal("expected get user callback")
	}

	called = false
	getCurrentUserOptionFunc(func(*getCurrentUserOptions) { called = true }).applyGetCurrentUser(&getCurrentUserOptions{})
	if !called {
		t.Fatal("expected get current user callback")
	}

	called = false
	getUserPermissionsOptionFunc(func(*getUserPermissionsOptions) { called = true }).applyGetUserPermissions(&getUserPermissionsOptions{})
	if !called {
		t.Fatal("expected get user permissions callback")
	}

	called = false
	getFilterOptionFunc(func(*getFilterOptions) { called = true }).applyGetFilter(&getFilterOptions{})
	if !called {
		t.Fatal("expected get filter callback")
	}

	called = false
	deleteFilterOptionFunc(func(*deleteFilterOptions) { called = true }).applyDeleteFilter(&deleteFilterOptions{})
	if !called {
		t.Fatal("expected delete filter callback")
	}

	called = false
	deleteFiltersOptionFunc(func(*deleteFiltersOptions) { called = true }).applyDeleteFilters(&deleteFiltersOptions{})
	if !called {
		t.Fatal("expected delete filters callback")
	}

	called = false
	listFilterHelpersOptionFunc(func(*listFilterHelpersOptions) { called = true }).applyListFilterHelpers(&listFilterHelpersOptions{})
	if !called {
		t.Fatal("expected list filter helpers callback")
	}

	called = false
	listActivityTypesOptionFunc(func(*listActivityTypesOptions) { called = true }).applyListActivityTypes(&listActivityTypesOptions{})
	if !called {
		t.Fatal("expected list activity types callback")
	}

	called = false
	listLeadLabelsOptionFunc(func(*listLeadLabelsOptions) { called = true }).applyListLeadLabels(&listLeadLabelsOptions{})
	if !called {
		t.Fatal("expected list lead labels callback")
	}

	called = false
	getOrganizationRelationshipOptionFunc(func(*getOrganizationRelationshipOptions) { called = true }).applyGetOrganizationRelationship(&getOrganizationRelationshipOptions{})
	if !called {
		t.Fatal("expected get organization relationship callback")
	}

	called = false
	deleteOrganizationRelationshipOptionFunc(func(*deleteOrganizationRelationshipOptions) { called = true }).applyDeleteOrganizationRelationship(&deleteOrganizationRelationshipOptions{})
	if !called {
		t.Fatal("expected delete organization relationship callback")
	}
}
