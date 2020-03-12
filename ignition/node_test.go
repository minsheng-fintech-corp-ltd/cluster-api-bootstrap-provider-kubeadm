package ignition

import (
	"testing"

	"sigs.k8s.io/cluster-api-bootstrap-provider-kubeadm/api/v1alpha2"
)

func TestToUserData(t *testing.T) {
	for _, testcase := range testCases {
		if data, err := testcase.input.ToUserData(); err != nil {
			t.Error(err)
		} else if string(data) != testcase.expected {
			t.Errorf("got \n %s \n expected\n %s\n", string(data), testcase.expected)
		}
	}
}

var testCases = []struct {
	expected string
	input    Node
}{
	{
		expected: `{"ignition":{"config":{"append":[{"source":"s3://container-service-demo/ignition-config/k8s-1.17.3-CoreOS-stable-2191.5.0.ign","verification":{}}]},"security":{"tls":{}},"timeouts":{},"version":"2.2.0"},"networkd":{},"passwd":{},"storage":{},"systemd":{}}`,
		input:    Node{},
	},
	{
		expected: `{"ignition":{"config":{"append":[{"source":"s3://container-service-demo/ignition-config/k8s-1.17.3-CoreOS-stable-2191.5.0.ign","verification":{}}]},"security":{"tls":{}},"timeouts":{},"version":"2.2.0"},"networkd":{},"passwd":{},"storage":{"files":[{"filesystem":"root","overwrite":true,"path":"/etc/docker/daemon.json","contents":{"source":"data:,%7B%22bridge%22%3A%22none%22%2C%22log-driver%22%3A%20%22json-file%22%2C%22log-opts%22%3A%20%7B%22max-size%22%3A%20%2210m%22%2C%22max-file%22%3A%20%2210%22%7D%2C%22live-restore%22%3A%20true%2C%22max-concurrent-downloads%22%3A10%7D%0A","verification":{}},"mode":420}]},"systemd":{}}`,
		input: Node{
			Files: []v1alpha2.File{
				{
					Path:        "/etc/docker/daemon.json",
					Permissions: "0644",
					Content: `{"bridge":"none","log-driver": "json-file","log-opts": {"max-size": "10m","max-file": "10"},"live-restore": true,"max-concurrent-downloads":10}
`,
				},
			},
		},
	},
	{
		expected: `{"ignition":{"config":{"append":[{"source":"s3://container-service-demo/ignition-config/k8s-1.17.3-CoreOS-stable-2191.5.0.ign","verification":{}}]},"security":{"tls":{}},"timeouts":{},"version":"2.2.0"},"networkd":{},"passwd":{},"storage":{},"systemd":{"units":[{"contents":"[Unit]\nDescription=extract k8s files\n\n[Service]\nType=oneshot\nExecStart=/usr/bin/tar xzvf /opt/kubernetes.tar.gz -C /\n\n[Install]\nWantedBy=multi-user.target\n","enabled":true,"name":"extractk8s.service"}]}}`,
		input: Node{
			Services: []ServiceUnit{
				{
					Content: "[Unit]\nDescription=extract k8s files\n\n[Service]\nType=oneshot\nExecStart=/usr/bin/tar xzvf /opt/kubernetes.tar.gz -C /\n\n[Install]\nWantedBy=multi-user.target\n",
					Enabled: true,
					Name:    "extractk8s.service",
				},
			},
		},
	},
	{
		expected: `{"ignition":{"config":{"append":[{"source":"s3://container-service-demo/ignition-config/k8s-1.17.3-CoreOS-stable-2191.5.0.ign","verification":{}}]},"security":{"tls":{}},"timeouts":{},"version":"2.2.0"},"networkd":{},"passwd":{},"storage":{},"systemd":{"units":[{"contents":"[Unit]\nDescription=extract k8s files\n\n[Service]\nType=oneshot\nExecStart=/usr/bin/tar xzvf /opt/kubernetes.tar.gz -C /\n\n[Install]\nWantedBy=multi-user.target\n","dropins":[{"contents":"[Service]\nLimitMEMLOCK=infinity\n","name":"30-increase-ulimit.conf"}],"enabled":true,"name":"extractk8s.service"}]}}`,
		input: Node{
			Services: []ServiceUnit{
				{
					Content: "[Unit]\nDescription=extract k8s files\n\n[Service]\nType=oneshot\nExecStart=/usr/bin/tar xzvf /opt/kubernetes.tar.gz -C /\n\n[Install]\nWantedBy=multi-user.target\n",
					Enabled: true,
					Name:    "extractk8s.service",
					Dropins: []Dropin{
						{
							Name:    "30-increase-ulimit.conf",
							Content: "[Service]\nLimitMEMLOCK=infinity\n",
						},
					},
				},
			},
		},
	},
}
