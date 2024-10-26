package groups_test

import (
	"context"
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/simulot/immich-go/internal/filenames"
	"github.com/simulot/immich-go/internal/groups"
	"github.com/simulot/immich-go/internal/groups/burst"
	"github.com/simulot/immich-go/internal/groups/groupby"
	"github.com/simulot/immich-go/internal/groups/series"
	"github.com/simulot/immich-go/internal/metadata"
)

type mockedAsset struct {
	nameInfo  filenames.NameInfo
	dateTaken time.Time
}

func (m mockedAsset) NameInfo() filenames.NameInfo {
	return m.nameInfo
}

func (m mockedAsset) DateTaken() time.Time {
	return m.dateTaken
}

func mockAsset(ic *filenames.InfoCollector, name string, dateTaken time.Time) groups.Asset {
	return mockedAsset{
		nameInfo:  ic.GetInfo(name),
		dateTaken: dateTaken,
	}
}

func TestGroup(t *testing.T) {
	ic := filenames.NewInfoCollector(time.Local, metadata.DefaultSupportedMedia)
	baseTime := time.Date(2021, 1, 1, 0, 0, 0, 0, time.Local)

	assets := []groups.Asset{
		mockAsset(ic, "photo1.jpg", baseTime.Add(50*time.Hour)),
		mockAsset(ic, "photo2.jpg", baseTime.Add(55*time.Hour)),
		mockAsset(ic, "IMG_001.jpg", baseTime),                            // Group 1
		mockAsset(ic, "IMG_002.jpg", baseTime.Add(200*time.Millisecond)),  // Group 1
		mockAsset(ic, "IMG_003.jpg", baseTime.Add(400*time.Millisecond)),  // Group 1
		mockAsset(ic, "IMG_004.jpg", baseTime.Add(600*time.Millisecond)),  // Group 1
		mockAsset(ic, "IMG_005.jpg", baseTime.Add(800*time.Millisecond)),  // Group 1
		mockAsset(ic, "IMG_006.jpg", baseTime.Add(1000*time.Millisecond)), // Group 1
		mockAsset(ic, "IMG_007.jpg", baseTime.Add(1200*time.Millisecond)), // Group 1
		mockAsset(ic, "IMG_008.jpg", baseTime.Add(1400*time.Millisecond)), // Group 1
		mockAsset(ic, "IMG_009.jpg", baseTime.Add(1600*time.Millisecond)), // Group 1
		mockAsset(ic, "photo3.jpg", baseTime.Add(5*time.Hour)),
		mockAsset(ic, "photo4.jpg", baseTime.Add(6*time.Hour)),
		mockAsset(ic, "IMG_001.jpg", baseTime.Add(7*time.Hour)),
		mockAsset(ic, "IMG_20231014_183246_BURST001_COVER.jpg", time.Date(2023, 10, 14, 18, 32, 46, 0, time.Local)), // Group 2
		mockAsset(ic, "IMG_20231014_183246_BURST002.jpg", time.Date(2023, 10, 14, 18, 32, 46, 0, time.Local)),       // Group 2
		mockAsset(ic, "IMG_003.jpg", baseTime.Add(9*time.Hour)),                                                     // Group 3
		mockAsset(ic, "IMG_003.raw", baseTime.Add(9*time.Hour)),                                                     // Group 3
		mockAsset(ic, "IMG_004.heic", baseTime.Add(10*time.Hour)),                                                   // Group 4
		mockAsset(ic, "IMG_004.jpg", baseTime.Add(10*time.Hour+100*time.Millisecond)),                               // Group 4
		mockAsset(ic, "IMG_005.raw", baseTime.Add(100*time.Hour)),
		mockAsset(ic, "IMG_005.jpg", baseTime.Add(101*time.Hour)),
		mockAsset(ic, "00001IMG_00001_BURST20210101153000.jpg", time.Date(2021, 1, 1, 15, 30, 0, 0, time.Local)),       // Group 5
		mockAsset(ic, "00002IMG_00002_BURST20210101153000_COVER.jpg", time.Date(2021, 1, 1, 15, 30, 0, 0, time.Local)), // Group 5
		mockAsset(ic, "00003IMG_00003_BURST20210101153000.jpg", time.Date(2021, 1, 1, 15, 30, 0, 0, time.Local)),       // Group 5
		mockAsset(ic, "IMG_006.heic", baseTime.Add(110*time.Hour)),
		mockAsset(ic, "photo5.jpg", baseTime.Add(120*time.Hour)),
		mockAsset(ic, "photo6.jpg", baseTime.Add(130*time.Hour)),
	}

	expectedGroup := []*groups.AssetGroup{
		groups.NewAssetGroup(groupby.GroupByBurst,
			mockAsset(ic, "00001IMG_00001_BURST20210101153000.jpg", time.Date(2021, 1, 1, 15, 30, 0, 0, time.Local)),
			mockAsset(ic, "00002IMG_00002_BURST20210101153000_COVER.jpg", time.Date(2021, 1, 1, 15, 30, 0, 0, time.Local)),
			mockAsset(ic, "00003IMG_00003_BURST20210101153000.jpg", time.Date(2021, 1, 1, 15, 30, 0, 0, time.Local)),
		).SetCover(1),
		groups.NewAssetGroup(groupby.GroupByBurst,
			mockAsset(ic, "IMG_001.jpg", baseTime),
			mockAsset(ic, "IMG_002.jpg", baseTime.Add(200*time.Millisecond)),
			mockAsset(ic, "IMG_003.jpg", baseTime.Add(400*time.Millisecond)),
			mockAsset(ic, "IMG_004.jpg", baseTime.Add(600*time.Millisecond)),
			mockAsset(ic, "IMG_005.jpg", baseTime.Add(800*time.Millisecond)),
			mockAsset(ic, "IMG_006.jpg", baseTime.Add(1000*time.Millisecond)),
			mockAsset(ic, "IMG_007.jpg", baseTime.Add(1200*time.Millisecond)),
			mockAsset(ic, "IMG_008.jpg", baseTime.Add(1400*time.Millisecond)),
			mockAsset(ic, "IMG_009.jpg", baseTime.Add(1600*time.Millisecond)),
		).SetCover(0),
		groups.NewAssetGroup(groupby.GroupByBurst,
			mockAsset(ic, "IMG_20231014_183246_BURST001_COVER.jpg", time.Date(2023, 10, 14, 18, 32, 46, 0, time.Local)),
			mockAsset(ic, "IMG_20231014_183246_BURST002.jpg", time.Date(2023, 10, 14, 18, 32, 46, 0, time.Local)),
		).SetCover(0),
		groups.NewAssetGroup(groupby.GroupByHeicJpg,
			mockAsset(ic, "IMG_004.heic", baseTime.Add(10*time.Hour)),
			mockAsset(ic, "IMG_004.jpg", baseTime.Add(10*time.Hour+100*time.Millisecond)),
		),
		groups.NewAssetGroup(groupby.GroupByRawJpg,
			mockAsset(ic, "IMG_003.jpg", baseTime.Add(9*time.Hour)),
			mockAsset(ic, "IMG_003.raw", baseTime.Add(9*time.Hour)),
		),
	}

	expectedAssets := []groups.Asset{
		mockAsset(ic, "photo1.jpg", baseTime.Add(50*time.Hour)),
		mockAsset(ic, "photo2.jpg", baseTime.Add(55*time.Hour)),
		mockAsset(ic, "photo3.jpg", baseTime.Add(5*time.Hour)),
		mockAsset(ic, "photo4.jpg", baseTime.Add(6*time.Hour)),
		mockAsset(ic, "IMG_001.jpg", baseTime.Add(7*time.Hour)),
		mockAsset(ic, "IMG_005.raw", baseTime.Add(100*time.Hour)),
		mockAsset(ic, "IMG_005.jpg", baseTime.Add(101*time.Hour)),
		mockAsset(ic, "IMG_006.heic", baseTime.Add(110*time.Hour)),
		mockAsset(ic, "photo5.jpg", baseTime.Add(120*time.Hour)),
		mockAsset(ic, "photo6.jpg", baseTime.Add(130*time.Hour)),
	}

	in := make(chan groups.Asset)
	out := make(chan groups.Asset)
	gOut := make(chan *groups.AssetGroup)

	go func() {
		for _, a := range assets {
			in <- a
		}
		close(in)
	}()

	ctx := context.Background()
	go func() {
		groups.NewGrouperPipeline(ctx, series.Group, burst.Group).Group(ctx, in, out, gOut)
		close(out)
		close(gOut)
	}()

	var gotGroups []*groups.AssetGroup
	var gotAssets []groups.Asset

	doneGroups := false
	doneAssets := false

	for !doneGroups || !doneAssets {
		select {
		case group, ok := <-gOut:
			if !ok {
				doneGroups = true
				continue
			}
			gotGroups = append(gotGroups, group)
		case asset, ok := <-out:
			if !ok {
				doneAssets = true
				continue
			}
			gotAssets = append(gotAssets, asset)
		}
	}

	sortGroupFn := func(s []*groups.AssetGroup) func(i, j int) bool {
		return func(i, j int) bool {
			if s[i].Assets[0].NameInfo().Radical == s[j].Assets[0].NameInfo().Radical {
				return s[i].Assets[0].DateTaken().Before(s[j].Assets[0].DateTaken())
			}
			return s[i].Assets[0].NameInfo().Radical < s[j].Assets[0].NameInfo().Radical
		}
	}

	sort.Slice(expectedGroup, sortGroupFn(expectedGroup))
	sort.Slice(gotGroups, sortGroupFn(gotGroups))
	if len(gotGroups) != len(expectedGroup) {
		t.Errorf("Expected %d group, got %d", len(expectedGroup), len(gotGroups))
	} else {
		for i := range gotGroups {
			for j := range gotGroups[i].Assets {
				got := gotGroups[i].Assets[j].(mockedAsset)
				expected := expectedGroup[i].Assets[j].(mockedAsset)
				if !reflect.DeepEqual(got, expected) {
					t.Errorf("Expected group %d asset %d \n%#v got\n%#v", i, j, expected, got)
				}
			}
		}
	}

	sortAssetFn := func(s []groups.Asset) func(i, j int) bool {
		return func(i, j int) bool {
			if s[i].NameInfo().Radical == s[j].NameInfo().Radical {
				if s[i].NameInfo().Index == s[j].NameInfo().Index {
					return s[i].DateTaken().Before(s[j].DateTaken())
				}
				return s[i].NameInfo().Index < s[j].NameInfo().Index
			}
			return s[i].NameInfo().Radical < s[j].NameInfo().Radical
		}
	}

	sort.Slice(expectedAssets, sortAssetFn(expectedAssets))
	sort.Slice(gotAssets, sortAssetFn(gotAssets))
	if len(gotAssets) != len(expectedAssets) {
		t.Errorf("Expected %d assets, got %d", len(expectedAssets), len(gotAssets))
	} else {
		for i := range gotAssets {
			if !reflect.DeepEqual(gotAssets[i], expectedAssets[i]) {
				t.Errorf("Expected asset \n%#v got asset \n%#v", expectedAssets[i], gotAssets[i])
			}
		}
	}
}
