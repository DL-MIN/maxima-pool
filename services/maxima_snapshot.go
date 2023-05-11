/*******************************************************************************
 * Service: maxima snapshots
 *
 * @author     Lars Thoms <lars@thoms.io>
 * @date       2023-05-11
 ******************************************************************************/

package services

import (
	"Moodle_Maxima_Pool/models"
	_ "embed"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/hashicorp/go-version"
	"github.com/spf13/viper"
	"os"
	"path"
	"regexp"
)

type ErrNoSnapshotsFound struct{}

func (e ErrNoSnapshotsFound) Error() string {
	return "no snapshots are found in storage path"
}

type ErrVersionNotFound string

func (e ErrVersionNotFound) Error() string {
	return "could not find a version string in tag " + string(e)
}

var (
	//go:embed maximalocal.mac
	maximaLocal []byte

	maximaVersionRegex = regexp.MustCompile("stackmaximaversion:([0-9]{10})\\$")
	maximaSnapshotList models.MaximaSnapshotList
)

func MaximaSnapshotCreate() (err error) {
	// Remove old versions
	dir, err := os.ReadDir(viper.GetString("storage.data"))
	if err != nil {
		return
	}
	for _, item := range dir {
		if err = os.RemoveAll(path.Join(viper.GetString("storage.data"), item.Name())); err != nil {
			return
		}
	}

	// Workspace for repository
	workspace, err := os.MkdirTemp(viper.GetString("storage.workspace"), "maxima-")
	if err != nil {
		return
	}
	defer func() {
		_ = os.RemoveAll(workspace)
	}()

	// Minimum version constraint
	versionConstraint, err := version.NewConstraint(">=" + viper.GetString("maxima.min_version"))
	if err != nil {
		return
	}

	// Clone git repository
	repository, err := git.PlainClone(workspace, false, &git.CloneOptions{URL: viper.GetString("maxima.repository")})
	if err != nil {
		return
	}

	// Get repository's worktree
	worktree, err := repository.Worktree()
	if err != nil {
		return
	}

	// Get all tags from repository and process them
	tags, err := repository.TagObjects()
	if err != nil {
		return
	}
	err = tags.ForEach(func(item *object.Tag) (err error) {
		tag, err := version.NewVersion(item.Name)
		if !versionConstraint.Check(tag) || err != nil {
			return
		}

		// Reset repository to tag
		if err = worktree.Reset(&git.ResetOptions{Commit: item.Target, Mode: git.HardReset}); err != nil {
			return
		}

		return maximaSnapshotCreate(workspace, item)
	})

	err = maximaSnapshotList.Store()
	return
}

func maximaSnapshotCreate(workspace string, tag *object.Tag) (err error) {
	stackmaxima, err := os.ReadFile(path.Join(workspace, "stack", "maxima", "stackmaxima.mac"))
	if err != nil {
		return
	}

	match := maximaVersionRegex.FindSubmatch(stackmaxima)
	if match == nil {
		return ErrVersionNotFound(tag.Name)
	}

	batchString := fmt.Sprintf(
		`file_search_maxima:append([sconcat("%s")],file_search_maxima)$`+
			`file_search_lisp:append([sconcat("%s")],file_search_lisp)$`+
			`%s`+
			`:lisp (sb-ext:save-lisp-and-die "%s" `+
			`:toplevel #'run :executable t)`,
		path.Join(workspace, "stack", "maxima", "###.{mac,mc}"),
		path.Join(workspace, "stack", "maxima", "###.{lisp}"),
		maximaLocal,
		path.Join(viper.GetString("storage.data"), "maxima-"+string(match[1])))

	_, _, _, clean, err := CommandCreate(viper.GetDuration("job.timeout"), "", viper.GetString("maxima.command"), "--quiet", "--batch-string", batchString)
	defer clean()
	if err != nil {
		return
	}
	maximaSnapshotList = append(maximaSnapshotList, models.MaximaSnapshot(match[1]))
	return
}

func MaximaSnapshotGet(v string) (version string, err error) {
	if maximaSnapshotList == nil {
		_ = maximaSnapshotList.Load()
	}

	if len(maximaSnapshotList) == 0 {
		return "", &ErrNoSnapshotsFound{}
	}

	version = string((maximaSnapshotList)[len(maximaSnapshotList)-1])

	for _, item := range maximaSnapshotList {
		if string(item) == v {
			return v, nil
		}
	}

	return
}
