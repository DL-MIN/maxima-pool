/*******************************************************************************
 * Model: maxima snapshots
 *
 * @author     Lars Thoms <lars@thoms.io>
 * @date       2023-05-11
 ******************************************************************************/

package models

import (
	"encoding/gob"
	"github.com/spf13/viper"
	"os"
	"path"
)

type MaximaSnapshot string
type MaximaSnapshotList []MaximaSnapshot

func (l *MaximaSnapshotList) Store() (err error) {
	file, err := os.Create(path.Join(viper.GetString("storage.data"), "maxima-versions.gob"))
	if err != nil {
		return
	}

	if err = gob.NewEncoder(file).Encode(l); err != nil {
		return
	}

	err = file.Close()
	return
}

func (l *MaximaSnapshotList) Load() (err error) {
	file, err := os.Open(path.Join(viper.GetString("storage.data"), "maxima-versions.gob"))
	if err != nil {
		return
	}

	if err = gob.NewDecoder(file).Decode(&l); err != nil {
		return
	}
	err = file.Close()
	return
}
