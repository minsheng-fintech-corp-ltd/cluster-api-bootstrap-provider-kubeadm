package ignition

import (
	"encoding/json"
	"errors"
	"net/url"
	"reflect"
	"strconv"

	"github.com/coreos/ignition/config/util"
	"sigs.k8s.io/cluster-api-bootstrap-provider-kubeadm/api/v1alpha2"

	ignTypes "github.com/coreos/ignition/config/v2_2/types"

	"github.com/coreos/ignition/config/validate"
	"github.com/vincent-petithory/dataurl"
)

const (
	DefaultFileMode = 0644
	DefaultDirMode  = 0755
)

type Node struct {
	Files    []v1alpha2.File
	Services []ServiceUnit
}

func (node Node) ToUserData() (data []byte, err error) {
	out := getGeneratedIgnitionConfig()
	out.Systemd = getSystemd(node.Services)
	if out.Storage, err = getStorage(node.Files); err != nil {
		return []byte{}, err
	}
	//validate outputc
	validationReport := validate.ValidateWithoutSource(reflect.ValueOf(*out))
	if validationReport.IsFatal() {
		return []byte{}, errors.New(validationReport.String())
	}
	return json.Marshal(out)
}

func getStorage(files []v1alpha2.File) (out ignTypes.Storage, err error) {
	for _, file := range files {
		newFile := ignTypes.File{
			Node: ignTypes.Node{
				Filesystem: "root",
				Path:       file.Path,
				Overwrite:  boolToPtr(true),
			},
			FileEmbedded1: ignTypes.FileEmbedded1{
				Append: false,
				Mode:   intToPtr(DefaultFileMode),
			},
		}
		if file.Permissions != "" {
			value, err := strconv.ParseInt(file.Permissions, 8, 32)
			if err != nil {
				return ignTypes.Storage{}, err
			}
			newFile.FileEmbedded1.Mode = util.IntToPtr(int(value))
		}
		if file.Content != "" {
			newFile.Contents = ignTypes.FileContents{
				Source: (&url.URL{
					Scheme: "data",
					Opaque: "," + dataurl.EscapeString(file.Content),
				}).String(),
			}
		}
		out.Files = append(out.Files, newFile)
	}
	return out, nil
}

func getSystemd(services []ServiceUnit) (out ignTypes.Systemd) {
	for _, service := range services {
		newUnit := ignTypes.Unit{
			Name:     service.Name,
			Enabled:  boolToPtr(service.Enabled),
			Contents: service.Content,
		}

		for _, dropIn := range service.Dropins {
			newUnit.Dropins = append(newUnit.Dropins, ignTypes.SystemdDropin{
				Name:     dropIn.Name,
				Contents: dropIn.Content,
			})
		}

		out.Units = append(out.Units, newUnit)
	}
	return
}
