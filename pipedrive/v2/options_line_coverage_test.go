package v2

import (
	"testing"
	"time"

	"github.com/juhokoskela/pipedrive-go/pipedrive"
)

func requireOneRequestOption(t *testing.T, name string, opts []pipedrive.RequestOption) {
	t.Helper()

	if len(opts) != 1 {
		t.Fatalf("%s request options count = %d, want 1", name, len(opts))
	}
}

func TestProductRequestOptionsApplyAllTargets(t *testing.T) {
	t.Parallel()

	opt := WithProductRequestOptions(pipedrive.WithHeader("X-Test", "1"))

	var get getProductOptions
	opt.applyGetProduct(&get)
	requireOneRequestOption(t, "get", get.requestOptions)

	var list listProductsOptions
	opt.applyListProducts(&list)
	requireOneRequestOption(t, "list", list.requestOptions)

	var create createProductOptions
	opt.applyCreateProduct(&create)
	requireOneRequestOption(t, "create", create.requestOptions)

	var update updateProductOptions
	opt.applyUpdateProduct(&update)
	requireOneRequestOption(t, "update", update.requestOptions)

	var deleteProduct deleteProductOptions
	opt.applyDeleteProduct(&deleteProduct)
	requireOneRequestOption(t, "delete", deleteProduct.requestOptions)

	var search searchProductsOptions
	opt.applySearchProducts(&search)
	requireOneRequestOption(t, "search", search.requestOptions)

	var duplicate duplicateProductOptions
	opt.applyDuplicateProduct(&duplicate)
	requireOneRequestOption(t, "duplicate", duplicate.requestOptions)

	var listVariations listProductVariationsOptions
	opt.applyListProductVariations(&listVariations)
	requireOneRequestOption(t, "list variations", listVariations.requestOptions)

	var createVariation createProductVariationOptions
	opt.applyCreateProductVariation(&createVariation)
	requireOneRequestOption(t, "create variation", createVariation.requestOptions)

	var updateVariation updateProductVariationOptions
	opt.applyUpdateProductVariation(&updateVariation)
	requireOneRequestOption(t, "update variation", updateVariation.requestOptions)

	var deleteVariation deleteProductVariationOptions
	opt.applyDeleteProductVariation(&deleteVariation)
	requireOneRequestOption(t, "delete variation", deleteVariation.requestOptions)

	var getImage getProductImageOptions
	opt.applyGetProductImage(&getImage)
	requireOneRequestOption(t, "get image", getImage.requestOptions)

	var uploadImage uploadProductImageOptions
	opt.applyUploadProductImage(&uploadImage)
	requireOneRequestOption(t, "upload image", uploadImage.requestOptions)

	var updateImage updateProductImageOptions
	opt.applyUpdateProductImage(&updateImage)
	requireOneRequestOption(t, "update image", updateImage.requestOptions)

	var deleteImage deleteProductImageOptions
	opt.applyDeleteProductImage(&deleteImage)
	requireOneRequestOption(t, "delete image", deleteImage.requestOptions)

	var followers getProductFollowersOptions
	opt.applyGetProductFollowers(&followers)
	requireOneRequestOption(t, "followers", followers.requestOptions)

	var addFollower addProductFollowerOptions
	opt.applyAddProductFollower(&addFollower)
	requireOneRequestOption(t, "add follower", addFollower.requestOptions)

	var deleteFollower deleteProductFollowerOptions
	opt.applyDeleteProductFollower(&deleteFollower)
	requireOneRequestOption(t, "delete follower", deleteFollower.requestOptions)

	var changelog getProductFollowersChangelogOptions
	opt.applyGetProductFollowersChangelog(&changelog)
	requireOneRequestOption(t, "followers changelog", changelog.requestOptions)
}

func TestProductOptionsIgnoreEmptyInputs(t *testing.T) {
	t.Parallel()

	list := newListProductsOptions([]ListProductsOption{
		nil,
		WithProductsIDs(),
		WithProductsPageSize(0),
		WithProductsCursor(""),
		WithProductsCustomFields(),
	})
	if list.params.Ids != nil || list.params.Limit != nil || list.params.Cursor != nil || list.params.CustomFields != nil {
		t.Fatalf("expected empty list params, got %#v", list.params)
	}

	search := newSearchProductsOptions([]SearchProductsOption{
		nil,
		WithProductSearchFields(),
		WithProductSearchPageSize(0),
		WithProductSearchCursor(""),
	})
	if search.params.Fields != nil || search.params.Limit != nil || search.params.Cursor != nil {
		t.Fatalf("expected empty search params, got %#v", search.params)
	}

	variations := newListProductVariationsOptions([]ListProductVariationsOption{
		nil,
		WithProductVariationsPageSize(0),
		WithProductVariationsCursor(""),
	})
	if variations.params.Limit != nil || variations.params.Cursor != nil {
		t.Fatalf("expected empty variation params, got %#v", variations.params)
	}

	followers := newGetProductFollowersOptions([]GetProductFollowersOption{
		nil,
		WithProductFollowersPageSize(0),
		WithProductFollowersCursor(""),
	})
	if followers.params.Limit != nil || followers.params.Cursor != nil {
		t.Fatalf("expected empty follower params, got %#v", followers.params)
	}

	changelog := newGetProductFollowersChangelogOptions([]GetProductFollowersChangelogOption{
		nil,
		WithProductFollowersChangelogPageSize(0),
		WithProductFollowersChangelogCursor(""),
	})
	if changelog.params.Limit != nil || changelog.params.Cursor != nil {
		t.Fatalf("expected empty changelog params, got %#v", changelog.params)
	}

	create := newCreateProductOptions([]CreateProductOption{nil, WithProductPrices()})
	if len(create.payload.prices) != 0 {
		t.Fatalf("expected no product prices, got %#v", create.payload.prices)
	}

	variationCreate := newCreateProductVariationOptions([]CreateProductVariationOption{nil, WithProductVariationPrices()})
	if len(variationCreate.payload.prices) != 0 {
		t.Fatalf("expected no variation prices, got %#v", variationCreate.payload.prices)
	}

	if _, _, err := (productImagePayload{}).toMultipart(); err == nil {
		t.Fatal("expected missing product image file error")
	}

	_ = newGetProductOptions([]GetProductOption{nil})
	_ = newUpdateProductOptions([]UpdateProductOption{nil})
	_ = newDeleteProductOptions([]DeleteProductOption{nil})
	_ = newDuplicateProductOptions([]DuplicateProductOption{nil})
	_ = newGetProductImageOptions([]GetProductImageOption{nil})
	_ = newUploadProductImageOptions([]UploadProductImageOption{nil})
	_ = newUpdateProductImageOptions([]UpdateProductImageOption{nil})
	_ = newDeleteProductImageOptions([]DeleteProductImageOption{nil})
	_ = newUpdateProductVariationOptions([]UpdateProductVariationOption{nil})
	_ = newDeleteProductVariationOptions([]DeleteProductVariationOption{nil})
	_ = newAddProductFollowerOptions([]AddProductFollowerOption{nil})
	_ = newDeleteProductFollowerOptions([]DeleteProductFollowerOption{nil})
}

func TestPersonRequestOptionsApplyAllTargets(t *testing.T) {
	t.Parallel()

	opt := WithPersonRequestOptions(pipedrive.WithHeader("X-Test", "1"))

	var get getPersonOptions
	opt.applyGetPerson(&get)
	requireOneRequestOption(t, "get", get.requestOptions)

	var list listPersonsOptions
	opt.applyListPersons(&list)
	requireOneRequestOption(t, "list", list.requestOptions)

	var create createPersonOptions
	opt.applyCreatePerson(&create)
	requireOneRequestOption(t, "create", create.requestOptions)

	var update updatePersonOptions
	opt.applyUpdatePerson(&update)
	requireOneRequestOption(t, "update", update.requestOptions)

	var deletePerson deletePersonOptions
	opt.applyDeletePerson(&deletePerson)
	requireOneRequestOption(t, "delete", deletePerson.requestOptions)

	var search searchPersonsOptions
	opt.applySearchPersons(&search)
	requireOneRequestOption(t, "search", search.requestOptions)

	var followers getPersonFollowersOptions
	opt.applyGetPersonFollowers(&followers)
	requireOneRequestOption(t, "followers", followers.requestOptions)

	var addFollower addPersonFollowerOptions
	opt.applyAddPersonFollower(&addFollower)
	requireOneRequestOption(t, "add follower", addFollower.requestOptions)

	var deleteFollower deletePersonFollowerOptions
	opt.applyDeletePersonFollower(&deleteFollower)
	requireOneRequestOption(t, "delete follower", deleteFollower.requestOptions)

	var changelog getPersonFollowersChangelogOptions
	opt.applyGetPersonFollowersChangelog(&changelog)
	requireOneRequestOption(t, "followers changelog", changelog.requestOptions)

	var picture getPersonPictureOptions
	opt.applyGetPersonPicture(&picture)
	requireOneRequestOption(t, "picture", picture.requestOptions)
}

func TestPersonOptionsIgnoreEmptyInputs(t *testing.T) {
	t.Parallel()

	get := newGetPersonOptions([]GetPersonOption{
		nil,
		WithPersonIncludeFields(),
		WithPersonCustomFields(),
	})
	if get.params.IncludeFields != nil || get.params.CustomFields != nil {
		t.Fatalf("expected empty get params, got %#v", get.params)
	}

	list := newListPersonsOptions([]ListPersonsOption{
		nil,
		WithPersonsIncludeFields(),
		WithPersonsCustomFields(),
		WithPersonsIDs(),
		WithPersonsPageSize(0),
		WithPersonsCursor(""),
	})
	if list.params.IncludeFields != nil || list.params.CustomFields != nil || list.params.Ids != nil || list.params.Limit != nil || list.params.Cursor != nil {
		t.Fatalf("expected empty list params, got %#v", list.params)
	}

	search := newSearchPersonsOptions([]SearchPersonsOption{
		nil,
		WithPersonSearchFields(),
		WithPersonSearchIncludeFields(),
		WithPersonSearchPageSize(0),
		WithPersonSearchCursor(""),
	})
	if search.params.Fields != nil || search.params.IncludeFields != nil || search.params.Limit != nil || search.params.Cursor != nil {
		t.Fatalf("expected empty search params, got %#v", search.params)
	}

	create := newCreatePersonOptions([]CreatePersonOption{
		nil,
		WithPersonEmails(),
		WithPersonPhones(),
		WithPersonLabelIDs(),
	})
	if len(create.payload.emails) != 0 || len(create.payload.phones) != 0 || len(create.payload.labelIDs) != 0 {
		t.Fatalf("expected empty person payload, got %#v", create.payload)
	}

	followers := newGetPersonFollowersOptions([]GetPersonFollowersOption{
		nil,
		WithPersonFollowersPageSize(0),
		WithPersonFollowersCursor(""),
	})
	if followers.params.Limit != nil || followers.params.Cursor != nil {
		t.Fatalf("expected empty follower params, got %#v", followers.params)
	}

	changelog := newGetPersonFollowersChangelogOptions([]GetPersonFollowersChangelogOption{
		nil,
		WithPersonFollowersChangelogPageSize(0),
		WithPersonFollowersChangelogCursor(""),
	})
	if changelog.params.Limit != nil || changelog.params.Cursor != nil {
		t.Fatalf("expected empty changelog params, got %#v", changelog.params)
	}

	_ = newUpdatePersonOptions([]UpdatePersonOption{nil})
	_ = newDeletePersonOptions([]DeletePersonOption{nil})
	_ = newAddPersonFollowerOptions([]AddPersonFollowerOption{nil})
	_ = newDeletePersonFollowerOptions([]DeletePersonFollowerOption{nil})
	_ = newGetPersonPictureOptions([]GetPersonPictureOption{nil})
}

func TestOrganizationRequestOptionsApplyAllTargets(t *testing.T) {
	t.Parallel()

	opt := WithOrganizationRequestOptions(pipedrive.WithHeader("X-Test", "1"))

	var get getOrganizationOptions
	opt.applyGetOrganization(&get)
	requireOneRequestOption(t, "get", get.requestOptions)

	var list listOrganizationsOptions
	opt.applyListOrganizations(&list)
	requireOneRequestOption(t, "list", list.requestOptions)

	var create createOrganizationOptions
	opt.applyCreateOrganization(&create)
	requireOneRequestOption(t, "create", create.requestOptions)

	var update updateOrganizationOptions
	opt.applyUpdateOrganization(&update)
	requireOneRequestOption(t, "update", update.requestOptions)

	var deleteOrganization deleteOrganizationOptions
	opt.applyDeleteOrganization(&deleteOrganization)
	requireOneRequestOption(t, "delete", deleteOrganization.requestOptions)

	var search searchOrganizationsOptions
	opt.applySearchOrganizations(&search)
	requireOneRequestOption(t, "search", search.requestOptions)

	var followers getOrganizationFollowersOptions
	opt.applyGetOrganizationFollowers(&followers)
	requireOneRequestOption(t, "followers", followers.requestOptions)

	var addFollower addOrganizationFollowerOptions
	opt.applyAddOrganizationFollower(&addFollower)
	requireOneRequestOption(t, "add follower", addFollower.requestOptions)

	var deleteFollower deleteOrganizationFollowerOptions
	opt.applyDeleteOrganizationFollower(&deleteFollower)
	requireOneRequestOption(t, "delete follower", deleteFollower.requestOptions)

	var changelog getOrganizationFollowersChangelogOptions
	opt.applyGetOrganizationFollowersChangelog(&changelog)
	requireOneRequestOption(t, "followers changelog", changelog.requestOptions)
}

func TestOrganizationOptionsIgnoreEmptyInputs(t *testing.T) {
	t.Parallel()

	get := newGetOrganizationOptions([]GetOrganizationOption{
		nil,
		WithOrganizationIncludeFields(),
		WithOrganizationCustomFields(),
	})
	if get.params.IncludeFields != nil || get.params.CustomFields != nil {
		t.Fatalf("expected empty get params, got %#v", get.params)
	}

	until := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)
	listWithUntil := newListOrganizationsOptions([]ListOrganizationsOption{WithOrganizationsUpdatedUntil(until)})
	if listWithUntil.params.UpdatedUntil == nil {
		t.Fatal("expected updated_until option to be set")
	}

	list := newListOrganizationsOptions([]ListOrganizationsOption{
		nil,
		WithOrganizationsIncludeFields(),
		WithOrganizationsCustomFields(),
		WithOrganizationsIDs(),
		WithOrganizationsPageSize(0),
		WithOrganizationsCursor(""),
	})
	if list.params.IncludeFields != nil || list.params.CustomFields != nil || list.params.Ids != nil || list.params.Limit != nil || list.params.Cursor != nil {
		t.Fatalf("expected empty list params, got %#v", list.params)
	}

	search := newSearchOrganizationsOptions([]SearchOrganizationsOption{
		nil,
		WithOrganizationSearchFields(),
		WithOrganizationSearchPageSize(0),
		WithOrganizationSearchCursor(""),
	})
	if search.params.Fields != nil || search.params.Limit != nil || search.params.Cursor != nil {
		t.Fatalf("expected empty search params, got %#v", search.params)
	}

	create := newCreateOrganizationOptions([]CreateOrganizationOption{
		nil,
		WithOrganizationLabelIDs(),
	})
	if len(create.payload.labelIDs) != 0 {
		t.Fatalf("expected empty organization label IDs, got %#v", create.payload.labelIDs)
	}

	followers := newGetOrganizationFollowersOptions([]GetOrganizationFollowersOption{
		nil,
		WithOrganizationFollowersPageSize(0),
		WithOrganizationFollowersCursor(""),
	})
	if followers.params.Limit != nil || followers.params.Cursor != nil {
		t.Fatalf("expected empty follower params, got %#v", followers.params)
	}

	changelog := newGetOrganizationFollowersChangelogOptions([]GetOrganizationFollowersChangelogOption{
		nil,
		WithOrganizationFollowersChangelogPageSize(0),
		WithOrganizationFollowersChangelogCursor(""),
	})
	if changelog.params.Limit != nil || changelog.params.Cursor != nil {
		t.Fatalf("expected empty changelog params, got %#v", changelog.params)
	}

	_ = newUpdateOrganizationOptions([]UpdateOrganizationOption{nil})
	_ = newDeleteOrganizationOptions([]DeleteOrganizationOption{nil})
	_ = newAddOrganizationFollowerOptions([]AddOrganizationFollowerOption{nil})
	_ = newDeleteOrganizationFollowerOptions([]DeleteOrganizationFollowerOption{nil})
}
