# Pipedrive API v2 endpoints

Generated from `openapi/upstream/v2.yaml` by `cmd/endpoint-docs`. Do not edit manually.

Total operations: 123

| Method | Path | Summary | Operation ID |
| --- | --- | --- | --- |
| GET | `/activities` | Get all activities | `getActivities` |
| POST | `/activities` | Add a new activity | `addActivity` |
| GET | `/activities/{id}` | Get details of an activity | `getActivity` |
| PATCH | `/activities/{id}` | Update an activity | `updateActivity` |
| DELETE | `/activities/{id}` | Delete an activity | `deleteActivity` |
| GET | `/activityFields` | Get all activity fields | `getActivityFields` |
| GET | `/activityFields/{field_code}` | Get one activity field | `getActivityField` |
| GET | `/dealFields` | Get all deal fields | `getDealFields` |
| POST | `/dealFields` | Create one deal field | `addDealField` |
| GET | `/dealFields/{field_code}` | Get one deal field | `getDealField` |
| PATCH | `/dealFields/{field_code}` | Update one deal field | `updateDealField` |
| DELETE | `/dealFields/{field_code}` | Delete one deal field | `deleteDealField` |
| POST | `/dealFields/{field_code}/options` | Add deal field options in bulk | `addDealFieldOptions` |
| PATCH | `/dealFields/{field_code}/options` | Update deal field options in bulk | `updateDealFieldOptions` |
| DELETE | `/dealFields/{field_code}/options` | Delete deal field options in bulk | `deleteDealFieldOptions` |
| GET | `/deals` | Get all deals | `getDeals` |
| POST | `/deals` | Add a new deal | `addDeal` |
| GET | `/deals/archived` | Get all archived deals | `getArchivedDeals` |
| GET | `/deals/installments` | List installments added to a list of deals | `getInstallments` |
| GET | `/deals/products` | Get deal products of several deals | `getDealsProducts` |
| GET | `/deals/search` | Search deals | `searchDeals` |
| GET | `/deals/{id}` | Get details of a deal | `getDeal` |
| PATCH | `/deals/{id}` | Update a deal | `updateDeal` |
| DELETE | `/deals/{id}` | Delete a deal | `deleteDeal` |
| POST | `/deals/{id}/convert/lead` | Convert a deal to a lead (BETA) | `convertDealToLead` |
| GET | `/deals/{id}/convert/status/{conversion_id}` | Get Deal conversion status (BETA) | `getDealConversionStatus` |
| GET | `/deals/{id}/discounts` | List discounts added to a deal | `getAdditionalDiscounts` |
| POST | `/deals/{id}/discounts` | Add a discount to a deal | `postAdditionalDiscount` |
| PATCH | `/deals/{id}/discounts/{discount_id}` | Update a discount added to a deal | `updateAdditionalDiscount` |
| DELETE | `/deals/{id}/discounts/{discount_id}` | Delete a discount from a deal | `deleteAdditionalDiscount` |
| GET | `/deals/{id}/followers` | List followers of a deal | `getDealFollowers` |
| POST | `/deals/{id}/followers` | Add a follower to a deal | `addDealFollower` |
| GET | `/deals/{id}/followers/changelog` | List followers changelog of a deal | `getDealFollowersChangelog` |
| DELETE | `/deals/{id}/followers/{follower_id}` | Delete a follower from a deal | `deleteDealFollower` |
| POST | `/deals/{id}/installments` | Add an installment to a deal | `postInstallment` |
| PATCH | `/deals/{id}/installments/{installment_id}` | Update an installment added to a deal | `updateInstallment` |
| DELETE | `/deals/{id}/installments/{installment_id}` | Delete an installment from a deal | `deleteInstallment` |
| GET | `/deals/{id}/products` | List products attached to a deal | `getDealProducts` |
| POST | `/deals/{id}/products` | Add a product to a deal | `addDealProduct` |
| DELETE | `/deals/{id}/products` | Delete many products from a deal | `deleteManyDealProducts` |
| POST | `/deals/{id}/products/bulk` | Add multiple products to a deal | `addManyDealProducts` |
| PATCH | `/deals/{id}/products/{product_attachment_id}` | Update the product attached to a deal | `updateDealProduct` |
| DELETE | `/deals/{id}/products/{product_attachment_id}` | Delete an attached product from a deal | `deleteDealProduct` |
| GET | `/itemSearch` | Perform a search from multiple item types | `searchItem` |
| GET | `/itemSearch/field` | Perform a search using a specific field from an item type | `searchItemByField` |
| GET | `/leads/search` | Search leads | `searchLeads` |
| POST | `/leads/{id}/convert/deal` | Convert a lead to a deal (BETA) | `convertLeadToDeal` |
| GET | `/leads/{id}/convert/status/{conversion_id}` | Get Lead conversion status (BETA) | `getLeadConversionStatus` |
| GET | `/organizationFields` | Get all organization fields | `getOrganizationFields` |
| POST | `/organizationFields` | Create one organization field | `addOrganizationField` |
| GET | `/organizationFields/{field_code}` | Get one organization field | `getOrganizationField` |
| PATCH | `/organizationFields/{field_code}` | Update one organization field | `updateOrganizationField` |
| DELETE | `/organizationFields/{field_code}` | Delete one organization field | `deleteOrganizationField` |
| POST | `/organizationFields/{field_code}/options` | Add organization field options in bulk | `addOrganizationFieldOptions` |
| PATCH | `/organizationFields/{field_code}/options` | Update organization field options in bulk | `updateOrganizationFieldOptions` |
| DELETE | `/organizationFields/{field_code}/options` | Delete organization field options in bulk | `deleteOrganizationFieldOptions` |
| GET | `/organizations` | Get all organizations | `getOrganizations` |
| POST | `/organizations` | Add a new organization | `addOrganization` |
| GET | `/organizations/search` | Search organizations | `searchOrganization` |
| GET | `/organizations/{id}` | Get details of a organization | `getOrganization` |
| PATCH | `/organizations/{id}` | Update a organization | `updateOrganization` |
| DELETE | `/organizations/{id}` | Delete a organization | `deleteOrganization` |
| GET | `/organizations/{id}/followers` | List followers of an organization | `getOrganizationFollowers` |
| POST | `/organizations/{id}/followers` | Add a follower to an organization | `addOrganizationFollower` |
| GET | `/organizations/{id}/followers/changelog` | List followers changelog of an organization | `getOrganizationFollowersChangelog` |
| DELETE | `/organizations/{id}/followers/{follower_id}` | Delete a follower from an organization | `deleteOrganizationFollower` |
| GET | `/personFields` | Get all person fields | `getPersonFields` |
| POST | `/personFields` | Create one person field | `addPersonField` |
| GET | `/personFields/{field_code}` | Get one person field | `getPersonField` |
| PATCH | `/personFields/{field_code}` | Update one person field | `updatePersonField` |
| DELETE | `/personFields/{field_code}` | Delete one person field | `deletePersonField` |
| POST | `/personFields/{field_code}/options` | Add person field options in bulk | `addPersonFieldOptions` |
| PATCH | `/personFields/{field_code}/options` | Update person field options in bulk | `updatePersonFieldOptions` |
| DELETE | `/personFields/{field_code}/options` | Delete person field options in bulk | `deletePersonFieldOptions` |
| GET | `/persons` | Get all persons | `getPersons` |
| POST | `/persons` | Add a new person | `addPerson` |
| GET | `/persons/search` | Search persons | `searchPersons` |
| GET | `/persons/{id}` | Get details of a person | `getPerson` |
| PATCH | `/persons/{id}` | Update a person | `updatePerson` |
| DELETE | `/persons/{id}` | Delete a person | `deletePerson` |
| GET | `/persons/{id}/followers` | List followers of a person | `getPersonFollowers` |
| POST | `/persons/{id}/followers` | Add a follower to a person | `addPersonFollower` |
| GET | `/persons/{id}/followers/changelog` | List followers changelog of a person | `getPersonFollowersChangelog` |
| DELETE | `/persons/{id}/followers/{follower_id}` | Delete a follower from a person | `deletePersonFollower` |
| GET | `/persons/{id}/picture` | Get picture of a person | `getPersonPicture` |
| GET | `/pipelines` | Get all pipelines | `getPipelines` |
| POST | `/pipelines` | Add a new pipeline | `addPipeline` |
| GET | `/pipelines/{id}` | Get one pipeline | `getPipeline` |
| PATCH | `/pipelines/{id}` | Update a pipeline | `updatePipeline` |
| DELETE | `/pipelines/{id}` | Delete a pipeline | `deletePipeline` |
| GET | `/productFields` | Get all product fields | `getProductFields` |
| POST | `/productFields` | Create one product field | `addProductField` |
| GET | `/productFields/{field_code}` | Get one product field | `getProductField` |
| PATCH | `/productFields/{field_code}` | Update one product field | `updateProductField` |
| DELETE | `/productFields/{field_code}` | Delete one product field | `deleteProductField` |
| POST | `/productFields/{field_code}/options` | Add product field options in bulk | `addProductFieldOptions` |
| PATCH | `/productFields/{field_code}/options` | Update product field options in bulk | `updateProductFieldOptions` |
| DELETE | `/productFields/{field_code}/options` | Delete product field options in bulk | `deleteProductFieldOptions` |
| GET | `/products` | Get all products | `getProducts` |
| POST | `/products` | Add a product | `addProduct` |
| GET | `/products/search` | Search products | `searchProducts` |
| GET | `/products/{id}` | Get one product | `getProduct` |
| PATCH | `/products/{id}` | Update a product | `updateProduct` |
| DELETE | `/products/{id}` | Delete a product | `deleteProduct` |
| POST | `/products/{id}/duplicate` | Duplicate a product | `duplicateProduct` |
| GET | `/products/{id}/followers` | List followers of a product | `getProductFollowers` |
| POST | `/products/{id}/followers` | Add a follower to a product | `addProductFollower` |
| GET | `/products/{id}/followers/changelog` | List followers changelog of a product | `getProductFollowersChangelog` |
| DELETE | `/products/{id}/followers/{follower_id}` | Delete a follower from a product | `deleteProductFollower` |
| GET | `/products/{id}/images` | Get image of a product | `getProductImage` |
| POST | `/products/{id}/images` | Upload an image for a product | `uploadProductImage` |
| PUT | `/products/{id}/images` | Update an image for a product | `updateProductImage` |
| DELETE | `/products/{id}/images` | Delete an image of a product | `deleteProductImage` |
| GET | `/products/{id}/variations` | Get all product variations | `getProductVariations` |
| POST | `/products/{id}/variations` | Add a product variation | `addProductVariation` |
| PATCH | `/products/{id}/variations/{product_variation_id}` | Update a product variation | `updateProductVariation` |
| DELETE | `/products/{id}/variations/{product_variation_id}` | Delete a product variation | `deleteProductVariation` |
| GET | `/stages` | Get all stages | `getStages` |
| POST | `/stages` | Add a new stage | `addStage` |
| GET | `/stages/{id}` | Get one stage | `getStage` |
| PATCH | `/stages/{id}` | Update stage details | `updateStage` |
| DELETE | `/stages/{id}` | Delete a stage | `deleteStage` |
| GET | `/users/{id}/followers` | List followers of a user | `getUserFollowers` |

