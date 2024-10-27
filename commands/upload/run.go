package upload

import (
	"context"
	"errors"
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/simulot/immich-go/adapters"
	gp "github.com/simulot/immich-go/adapters/googlePhotos"
	"github.com/simulot/immich-go/commands/application"
	"github.com/simulot/immich-go/immich"
	"github.com/simulot/immich-go/internal/albums"
	"github.com/simulot/immich-go/internal/fileevent"
	"github.com/simulot/immich-go/internal/groups"
	"github.com/simulot/immich-go/internal/groups/groupby"
)

type UpCmd struct {
	Mode UpLoadMode
	*UploadOptions
	app *application.Application

	AssetIndex       *AssetIndex     // List of assets present on the server
	deleteServerList []*immich.Asset // List of server assets to remove

	// deleteLocalList  []*adapters.LocalAssetFile // List of local assets to remove
	// stacks        *stacking.StackBuilder
	adapter       adapters.Adapter
	DebugCounters bool // Enable CSV action counters per file

	// fsyss  []fs.FS                            // pseudo file system to browse
	Paths  []string                          // Path to explore
	albums map[string]immich.AlbumSimplified // Albums by title

	takeoutOptions *gp.ImportFlags
}

func newUpload(mode UpLoadMode, app *application.Application, options *UploadOptions) *UpCmd {
	upCmd := &UpCmd{
		UploadOptions: options,
		app:           app,
		Mode:          mode,
	}
	return upCmd
}

func (upCmd *UpCmd) setTakeoutOptions(options *gp.ImportFlags) *UpCmd {
	upCmd.takeoutOptions = options
	return upCmd
}

func (upCmd *UpCmd) run(ctx context.Context, adapter adapters.Adapter, app *application.Application) error {
	upCmd.adapter = adapter
	// if app.CommonFlags.StackBurstPhotos || app.CommonFlags.StackJpgWithRaw {
	// 	app.stacks = stacking.NewStackBuilder(app.ImmichServerFlags.Immich.SupportedMedia())
	// }

	// todo counters
	// defer func() {
	// 	if app.DebugCounters {
	// 		fn := strings.TrimSuffix(app.LogFile, filepath.Ext(app.LogFile)) + ".csv"
	// 		f, err := os.Create(fn)
	// 		if err == nil {
	// 			_ = app.Jnl.WriteFileCounts(f)
	// 			fmt.Println("\nCheck the counters file: ", f.Name())
	// 			f.Close()
	// 		}
	// 	}
	// }()

	if upCmd.NoUI {
		return upCmd.runNoUI(ctx, app)
	}
	_, err := tcell.NewScreen()
	if err != nil {
		upCmd.app.Log().Error("can't initialize the screen for the UI mode. Falling back to no-gui mode")
		fmt.Println("can't initialize the screen for the UI mode. Falling back to no-gui mode")
		return upCmd.runNoUI(ctx, app)
	}
	return upCmd.runUI(ctx, app)
}

func (upCmd *UpCmd) getImmichAlbums(ctx context.Context) error {
	serverAlbums, err := upCmd.app.Client().Immich.GetAllAlbums(ctx)
	upCmd.albums = map[string]immich.AlbumSimplified{}
	if err != nil {
		return fmt.Errorf("can't get the album list from the server: %w", err)
	}
	for _, a := range serverAlbums {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			upCmd.albums[a.AlbumName] = a
		}
	}
	return nil
}

func (upCmd *UpCmd) getImmichAssets(ctx context.Context, updateFn progressUpdate) error {
	statistics, err := upCmd.app.Client().Immich.GetAssetStatistics(ctx)
	if err != nil {
		return err
	}
	totalOnImmich := statistics.Total
	received := 0

	var list []*immich.Asset

	err = upCmd.app.Client().Immich.GetAllAssetsWithFilter(ctx, func(a *immich.Asset) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			received++
			list = append(list, a)
			if updateFn != nil {
				updateFn(received, totalOnImmich)
			}
			return nil
		}
	})
	if err != nil {
		return err
	}
	if updateFn != nil {
		updateFn(totalOnImmich, totalOnImmich)
	}
	upCmd.AssetIndex = &AssetIndex{
		assets: list,
	}
	upCmd.AssetIndex.ReIndex()
	return nil
}

func (upCmd *UpCmd) uploadLoop(ctx context.Context, groupChan chan *groups.AssetGroup) error {
	var err error
assetLoop:
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case g, ok := <-groupChan:
			if !ok {
				break assetLoop
			}
			err = upCmd.handleGroup(ctx, g)
			if err != nil {
				return err
			}
		}
	}

	// if app.StackBurstPhotos || app.StackJpgWithRaw {
	// 	stacks := app.stacks.Stacks()
	// 	if len(stacks) > 0 {
	// 		app.Root.Log.Info("Creating stacks")
	// 	nextStack:
	// 		for _, s := range stacks {
	// 			switch {
	// 			case !app.StackBurstPhotos && s.StackType == stacking.StackBurst:
	// 				continue nextStack
	// 			case !app.StackJpgWithRaw && s.StackType == stacking.StackRawJpg:
	// 				continue nextStack
	// 			}
	// 			app.Root.Message(fmt.Sprintf("Stacking %s...", strings.Join(s.Names, ", ")))
	// 			err = app.Server.Immich.StackAssets(ctx, s.CoverID, s.IDs)
	// 			if err != nil {
	// 				app.Root.Log.Error(fmt.Sprintf("Can't stack images: %s", err))
	// 			}
	// 		}
	// 	}
	// }

	// if app.CreateAlbums || app.CreateAlbumAfterFolder || (app.KeepPartner && app.PartnerAlbum != "") || app.ImportIntoAlbum != "" {
	// 	app.Log.Info("Managing albums")
	// 	err = app.ManageAlbums(ctx)
	// 	if err != nil {
	// 		app.Log.Error(err.Error())
	// 		err = nil
	// 	}
	// }

	if len(upCmd.deleteServerList) > 0 {
		ids := []string{}
		for _, da := range upCmd.deleteServerList {
			ids = append(ids, da.ID)
		}
		err := upCmd.DeleteServerAssets(ctx, ids)
		if err != nil {
			return fmt.Errorf("can't delete server's assets: %w", err)
		}
	}

	// if len(app.deleteLocalList) > 0 {
	// 	err = app.DeleteLocalAssets()
	// }

	return err
}

func (upCmd *UpCmd) handleGroup(ctx context.Context, g *groups.AssetGroup) error {
	var errGroup error

	// Upload assets from the group
	for _, a := range g.Assets {
		if a, ok := a.(*adapters.LocalAssetFile); ok {
			err := upCmd.handleAsset(ctx, g, a)
			errGroup = errors.Join(err)
		}
	}
	if errGroup != nil {
		return errGroup
	}

	switch g.Grouping {
	case groupby.GroupByNone:
	}

	// Manage albums
	if len(g.Albums) > 0 {
		upCmd.manageGroupAlbums(ctx, g)
	}
	return nil
}

func (upCmd *UpCmd) handleAsset(ctx context.Context, g *groups.AssetGroup, a *adapters.LocalAssetFile) error {
	defer func() {
		a.Close() // Close and clean resources linked to the local asset
	}()

	advice, err := upCmd.AssetIndex.ShouldUpload(a)
	if err != nil {
		return err
	}

	switch advice.Advice {
	case NotOnServer: // Upload and manage albums
		err = upCmd.uploadAsset(ctx, a)
		return err
	case SmallerOnServer: // Upload, manage albums and delete the server's asset
		upCmd.app.Jnl().Record(ctx, fileevent.UploadUpgraded, a, "reason", advice.Message)

		// Remember existing asset's albums, if any
		for _, al := range advice.ServerAsset.Albums {
			g.AddAlbum(albums.Album{
				Title:       al.AlbumName,
				Description: al.Description,
			})
		}

		// Upload the superior asset
		err = upCmd.uploadAsset(ctx, a)
		if err != nil {
			return err
		}

		// delete the existing lower quality asset
		err = upCmd.app.Client().Immich.DeleteAssets(ctx, []string{advice.ServerAsset.ID}, true)
		if err != nil {
			upCmd.app.Jnl().Record(ctx, fileevent.Error, fileevent.FileAndName{}, "error", err.Error())
		}
		return err

	case SameOnServer:
		a.ID = advice.ServerAsset.ID
		for _, al := range advice.ServerAsset.Albums {
			g.AddAlbum(albums.Album{
				Title:       al.AlbumName,
				Description: al.Description,
			})
		}
		upCmd.app.Jnl().Record(ctx, fileevent.UploadServerDuplicate, fileevent.AsFileAndName(a.FSys, a.FileName), "reason", advice.Message)

	case BetterOnServer: // and manage albums
		upCmd.app.Jnl().Record(ctx, fileevent.UploadServerBetter, fileevent.AsFileAndName(a.FSys, a.FileName), "reason", advice.Message)
	}

	return nil
}

func (upCmd *UpCmd) uploadAsset(ctx context.Context, a *adapters.LocalAssetFile) error {
	defer upCmd.app.Log().Debug("", "file", a)
	ar, err := upCmd.app.Client().Immich.AssetUpload(ctx, a)
	if err != nil {
		upCmd.app.Jnl().Record(ctx, fileevent.UploadServerError, fileevent.AsFileAndName(a.FSys, a.FileName), "error", err.Error())
		return err // Must signal the error to the caller
	}
	if ar.Status == immich.UploadDuplicate {
		upCmd.app.Jnl().Record(ctx, fileevent.UploadServerDuplicate, fileevent.AsFileAndName(a.FSys, a.FileName), "reason", "the server has this file")
	} else {
		upCmd.app.Jnl().Record(ctx, fileevent.Uploaded, fileevent.AsFileAndName(a.FSys, a.FileName))
	}
	a.ID = ar.ID
	return nil
}

// manageGroupAlbums add the assets to the albums listed in the group.
// If an album does not exist, it is created.
// Errors are logged.
func (upCmd *UpCmd) manageGroupAlbums(ctx context.Context, g *groups.AssetGroup) {
	assetIDs := []string{}
	for _, a := range g.Assets {
		if a, ok := a.(*adapters.LocalAssetFile); ok {
			assetIDs = append(assetIDs, a.ID)
		}
	}

	for _, album := range g.Albums {
		title := album.Title
		l, exist := upCmd.albums[title]
		if !exist {
			newAl, err := upCmd.app.Client().Immich.CreateAlbum(ctx, title, album.Description, assetIDs)
			if err != nil {
				upCmd.app.Jnl().Record(ctx, fileevent.Error, fileevent.FileAndName{}, err)
			}
			upCmd.albums[title] = newAl
			l = newAl
		} else {
			_, err := upCmd.app.Client().Immich.AddAssetToAlbum(ctx, l.ID, assetIDs)
			if err != nil {
				upCmd.app.Jnl().Record(ctx, fileevent.Error, fileevent.FileAndName{}, err)
				return
			}
		}

		// Log the action
		for _, a := range g.Assets {
			if a, ok := a.(*adapters.LocalAssetFile); ok {
				upCmd.app.Jnl().Record(ctx, fileevent.UploadAddToAlbum, fileevent.AsFileAndName(a.FSys, a.FileName), "Album", title)
			}
		}
	}
}

func (upCmd *UpCmd) DeleteServerAssets(ctx context.Context, ids []string) error {
	upCmd.app.Log().Message("%d server assets to delete.", len(ids))
	return upCmd.app.Client().Immich.DeleteAssets(ctx, ids, false)
}

/*
func (app *UpCmd) DeleteLocalAssets() error {
	app.RootImmichFlags.Message(fmt.Sprintf("%d local assets to delete.", len(app.deleteLocalList)))

	for _, a := range app.deleteLocalList {
		if !app.DryRun {
			app.Log.Info(fmt.Sprintf("delete file %q", a.Title))
			err := a.Remove()
			if err != nil {
				return err
			}
		} else {
			app.Log.Info(fmt.Sprintf("file %q not deleted, dry run mode.", a.Title))
		}
	}
	return nil
}
*/
