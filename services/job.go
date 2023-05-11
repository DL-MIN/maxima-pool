/*******************************************************************************
 * Service: job
 *
 * @author     Lars Thoms <lars@thoms.io>
 * @date       2023-04-14
 ******************************************************************************/

package services

import (
	"Moodle_Maxima_Pool/models"
	"archive/zip"
	"bytes"
	"fmt"
	"github.com/spf13/viper"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"time"
)

func JobCreate(data *models.JobRequestQuery) (resp *models.JobResponse, err error) {
	resp = &models.JobResponse{
		Output: new(bytes.Buffer),
	}

	version, err := MaximaSnapshotGet(data.Version)
	if err != nil {
		return
	}

	stdOut, _, workspace, clean, err := CommandCreate(
		minDuration(viper.GetDuration("job.timeout"), time.Duration(data.Timeout)*time.Millisecond),
		fmt.Sprintf(`maxima_tempdir:getcurrentdirectory()$ IMAGE_DIR:getcurrentdirectory()$ URL_BASE:"%s"$\n%s`, data.PlotURLBase, data.Input),
		path.Join(viper.GetString("storage.data"), "maxima-"+version),
		"--quiet",
	)
	defer clean()
	err = jobResponse(resp, workspace, stdOut)
	return
}

func jobResponse(resp *models.JobResponse, workspace string, output []byte) (err error) {
	var file *zip.Writer

	err = filepath.WalkDir(workspace, func(itemPath string, item fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if item.IsDir() {
			return nil
		}

		if file == nil {
			file = zip.NewWriter(resp.Output)
			fileWriter, err := file.Create("OUTPUT")
			_, err = fileWriter.Write(output)
			if err != nil {
				return err
			}
		}

		fileReader, err := os.Open(path.Join(itemPath))
		if err != nil {
			return err
		}
		defer func(fileReader *os.File) {
			_ = fileReader.Close()
		}(fileReader)

		fileWriter, err := file.Create(item.Name())
		if err != nil {
			return err
		}

		_, err = io.Copy(fileWriter, fileReader)
		return err
	})

	if file == nil {
		_, err := resp.Output.Write(output)
		if err != nil {
			return err
		}
	} else {
		resp.IsZIP = true
		err = file.Close()
	}

	return
}

func minDuration(a time.Duration, b time.Duration) time.Duration {
	if a < b {
		return a
	}
	return b
}
