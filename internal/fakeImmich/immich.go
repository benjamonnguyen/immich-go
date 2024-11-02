package fakeimmich

import (
	"context"
	"io"

	"github.com/simulot/immich-go/browser"
	"github.com/simulot/immich-go/immich"
)

type MockedCLient struct{}

func (c *MockedCLient) GetAllAssetsWithFilter(context.Context, func(*immich.Asset) error) error {
	return nil
}

func (c *MockedCLient) GetAsset(ctx context.Context, ID string) (*immich.Asset, error) {
	return nil, nil
}

func (c *MockedCLient) AssetUpload(
	context.Context,
	*browser.LocalAssetFile,
) (immich.AssetResponse, error) {
	return immich.AssetResponse{}, nil
}

func (c *MockedCLient) DeleteAssets(context.Context, []string, bool) error {
	return nil
}

func (c *MockedCLient) GetAllAlbums(context.Context) ([]immich.AlbumSimplified, error) {
	return nil, nil
}

func (c *MockedCLient) AddAssetToAlbum(
	context.Context,
	string,
	[]string,
) ([]immich.UpdateAlbumResult, error) {
	return nil, nil
}

func (c *MockedCLient) CreateAlbum(
	context.Context,
	string,
	string,
	[]string,
) (immich.AlbumSimplified, error) {
	return immich.AlbumSimplified{}, nil
}

func (c *MockedCLient) UpdateAssets(
	ctx context.Context,
	ids []string,
	isArchived bool,
	isFavorite bool,
	latitude float64,
	longitude float64,
	removeParent bool,
	stackParentID string,
) error {
	return nil
}

func (c *MockedCLient) StackAssets(ctx context.Context, cover string, ids []string) error {
	return nil
}

func (c *MockedCLient) UpdateAsset(
	ctx context.Context,
	id string,
	a *browser.LocalAssetFile,
) (*immich.Asset, error) {
	return nil, nil
}

func (c *MockedCLient) EnableAppTrace(w io.Writer) {}

func (c *MockedCLient) GetServerStatistics(ctx context.Context) (immich.ServerStatistics, error) {
	return immich.ServerStatistics{}, nil
}

func (c *MockedCLient) PingServer(ctx context.Context) error {
	return nil
}

func (c *MockedCLient) SetDeviceUUID(string) {}

func (c *MockedCLient) SetEndPoint(string) {}

func (c *MockedCLient) ValidateConnection(ctx context.Context) (immich.User, error) {
	return immich.User{}, nil
}

func (c *MockedCLient) GetAssetAlbums(
	ctx context.Context,
	id string,
) ([]immich.AlbumSimplified, error) {
	return nil, nil
}

func (c *MockedCLient) GetAllAssets(ctx context.Context) ([]*immich.Asset, error) {
	return nil, nil
}

func (c *MockedCLient) DeleteAlbum(ctx context.Context, id string) error {
	return nil
}

func (c *MockedCLient) SupportedMedia() immich.SupportedMedia {
	return immich.DefaultSupportedMedia
}

func (c *MockedCLient) GetAssetStatistics(ctx context.Context) (immich.UserStatistics, error) {
	return immich.UserStatistics{
		Images: 1,
		Videos: 1,
		Total:  1,
	}, nil
}

func (c *MockedCLient) GetJobs(ctx context.Context) (map[string]immich.Job, error) {
	return nil, nil
}

func (c *MockedCLient) SendJobCommand(
	ctx context.Context,
	jobID immich.JobID,
	command immich.JobCommand,
	force bool,
) (immich.SendJobCommandResponse, error) {
	return immich.SendJobCommandResponse{}, nil
}

func (c *MockedCLient) GetAlbumInfo(context.Context, string, bool) (immich.AlbumContent, error) {
	return immich.AlbumContent{}, nil
}

func (c *MockedCLient) UpsertTags(
	ctx context.Context,
	tags []string,
) ([]immich.TagSimplified, error) {
	return nil, nil
}

func (c *MockedCLient) TagAssets(
	ctx context.Context,
	tagID string,
	assetIDs []string,
) ([]immich.TagAssetsResponse, error) {
	return nil, nil
}

func (c *MockedCLient) BulkTagAssets(
	ctx context.Context,
	tagIDs []string,
	assetIDs []string,
) (struct {
	Count int `json:"count"`
}, error) {
	return struct {
		Count int `json:"count"`
	}{}, nil
}

func (ic *MockedCLient) UntagAssets(ctx context.Context, tagID string, assetIDs []string) error {
	return nil
}

func (c *MockedCLient) CreateJob(ctx context.Context, name immich.JobName) error {
	return nil
}
