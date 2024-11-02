package immich

import (
	"context"
	"fmt"
)

type TagSimplified struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

type TagAssetsResponse struct {
	Error   string `json:"error,omitempty"` // [duplicate, no_permission, not_found, unknown]
	ID      string `json:"id"`
	Success bool   `json:"success"`
}

func (ic *ImmichClient) UpsertTags(ctx context.Context, tags []string) ([]TagSimplified, error) {
	var resp []TagSimplified
	body := struct {
		Tags []string `json:"tags"`
	}{Tags: tags}
	err := ic.newServerCall(ctx, EndPointUpsertTags).
		do(putRequest("/tags", setJSONBody(body), setAcceptJSON()), responseJSON(&resp))
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (ic *ImmichClient) TagAssets(
	ctx context.Context,
	tagID string,
	assetIDs []string,
) ([]TagAssetsResponse, error) {
	var resp []TagAssetsResponse

	body := struct {
		IDs []string `json:"ids"`
	}{IDs: assetIDs}
	err := ic.newServerCall(ctx, EndPointTagAssets).
		do(putRequest(fmt.Sprintf("/tags/%s/assets", tagID), setJSONBody(body), setAcceptJSON()), responseJSON(&resp))

	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (ic *ImmichClient) BulkTagAssets(
	ctx context.Context,
	tagIDs []string,
	assetIDs []string,
) (struct {
	Count int `json:"count"`
}, error) {
	var resp struct {
		Count int `json:"count"`
	}

	body := struct {
		TagIDs   []string `json:"tagIds"`
		AssetIDs []string `json:"assetIds"`
	}{
		TagIDs:   tagIDs,
		AssetIDs: assetIDs,
	}
	err := ic.newServerCall(ctx, EndPointBulkTagAssets).
		do(putRequest("/tags/assets", setJSONBody(body)), responseJSON(&resp))

	return resp, err
}

func (ic *ImmichClient) UntagAssets(ctx context.Context, tagID string, assetIDs []string) error {
	body := struct {
		IDs []string `json:"ids"`
	}{IDs: assetIDs}

	// TODO response body https://immich.app/docs/api/untag-assets
	return ic.newServerCall(ctx, EndPointUntagAssets).
		do(deleteRequest("/tags/"+tagID+"/assets", setJSONBody(body)))
}
