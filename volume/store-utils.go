package volume

import (
	"encoding/json"

	"github.com/gluster/glusterd2/context"
	"github.com/gluster/glusterd2/store"

	log "github.com/Sirupsen/logrus"
)

const (
	volumePrefix string = store.GlusterPrefix + "volume/"
)

func init() {
	context.Store.InitPrefix(volumePrefix)
}

// AddOrUpdateVolume marshals to volume object and passes to store to add/update
func AddOrUpdateVolume(v *Volinfo) error {
	json, e := json.Marshal(v)
	if e != nil {
		log.WithField("error", e).Error("Failed to marshal the volinfo object")
		return e
	}

	e = context.Store.Put(v.Name, json, nil)
	if e != nil {
		log.WithField("error", e).Error("Couldn't add volume to store")
		return e
	}
	return nil
}

// GetVolume fetches the json object from the store and unmarshalls it into
// volinfo object
func GetVolume(name string) (*Volinfo, error) {
	var v Volinfo
	b, e := context.Store.Get(volumePrefix + name)
	if e != nil {
		log.WithField("error", e).Error("Couldn't retrive volume from store")
		return nil, e
	}
	if e = json.Unmarshal(b.Value, &v); e != nil {
		log.WithField("error", e).Error("Failed to unmarshal the data into volinfo object")
		return nil, e
	}
	return &v, nil
}

//DeleteVolume passes the volname to store to delete the volume object
func DeleteVolume(name string) error {
	return context.Store.Delete(volumePrefix + name)
}

//GetVolumes retrives the json objects from the store and converts them into
//respective volinfo objects
func GetVolumes() ([]Volinfo, error) {
	pairs, e := context.Store.List(volumePrefix)
	if e != nil {
		return nil, e
	}

	volumes := make([]Volinfo, len(pairs))

	for index, pair := range pairs {
		var vol Volinfo

		if err := json.Unmarshal(pair.Value, &vol); err != nil {
			log.WithFields(log.Fields{
				"volume": pair.Key,
				"error":  err,
			}).Error("Failed to unmarshal volume")
			continue
		}
		volumes[index] = vol
	}

	return volumes, nil

}

//Exists check whether a given volume exist or not
func Exists(name string) bool {
	b, e := context.Store.Exists(volumePrefix + name)
	if e != nil {
		return false
	}

	return b
}