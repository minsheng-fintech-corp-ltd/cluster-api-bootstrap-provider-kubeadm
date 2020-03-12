package ignition

import (
	"bytes"
	"encoding/json"
	"flag"
	"net/url"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	ignTypes "github.com/coreos/ignition/config/v2_2/types"
	"github.com/google/uuid"
)

var (
	templateConfigUri string
	userdataDir       string
	userDataBucket    string
	ignitionLogger    = ctrl.Log.WithName("ignition")
)

func init() {
	flag.StringVar(
		&templateConfigUri,
		"ignition-base-url",
		"ignition-config/k8s-1.17.3-CoreOS-stable-2191.5.0.ign",
		"The address the base image ignition file resides",
	)
	flag.StringVar(
		&userDataBucket,
		"ignition-userdata-bucket",
		"container-service-demo",
		"The bucket the userdata ignition file resides",
	)
	flag.StringVar(
		&userdataDir,
		"ignition-userdata-dir",
		"node-userdata",
		"The bucket the userdata ignition file resides",
	)
}

func getGeneratedIgnitionConfig() *ignTypes.Config {
	baseIgnitionUrl := &url.URL{
		Scheme: "s3",
		Host:   userDataBucket,
		Path:   templateConfigUri,
	}
	return &ignTypes.Config{
		Ignition: ignTypes.Ignition{
			Version: "2.2.0",
			Config: ignTypes.IgnitionConfig{
				Append: []ignTypes.ConfigReference{
					{
						Source: baseIgnitionUrl.String(),
					},
				},
			},
		},
	}
}

func getCompressedIgnitionConfig(url string) *ignTypes.Config {
	return &ignTypes.Config{
		Ignition: ignTypes.Ignition{
			Version: "2.2.0",
			Config: ignTypes.IgnitionConfig{
				Replace: &ignTypes.ConfigReference{
					Source: url,
				},
			},
		},
	}
}

func GenerateUserData(node *Node) ([]byte, error) {
	userdata, err := node.ToUserData()
	if err != nil {
		ignitionLogger.Error(err, "failed to generate ignition file")
		return nil, err
	}
	session, err := session.NewSession()
	if err != nil {
		ignitionLogger.Error(err, "failed to initialize s3 session")
		return nil, err
	}
	uploader := s3manager.NewUploader(session)
	filePath := strings.Join([]string{userdataDir, uuid.New().String()}, "/")
	_, err = uploader.Upload(&s3manager.UploadInput{
		Body:         bytes.NewReader(userdata),
		Bucket:       aws.String(userDataBucket),
		Expires:      aws.Time(time.Now().Add(time.Hour * 168)),
		Key:          aws.String(filePath),
		StorageClass: aws.String(s3.StorageClassIntelligentTiering),
	})
	if err != nil {
		ignitionLogger.Error(err, "failed to upload ignition file to bucket")
		return nil, err
	}

	userDataUrl := url.URL{
		Scheme: "s3",
		Host:   userDataBucket,
		Path:   filePath,
	}

	return json.Marshal(getCompressedIgnitionConfig(userDataUrl.String()))
}

func intToPtr(i int) *int {
	return &i
}

func boolToPtr(b bool) *bool {
	return &b
}

func StringToPtr(s string) *string {
	return &s
}
