module github.com/flacatus/che-inspector

go 1.14

require (
	github.com/go-playground/universal-translator v0.17.0 // indirect
	github.com/go-playground/validator v9.31.0+incompatible
	github.com/google/go-cmp v0.5.2 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/mitchellh/go-homedir v1.1.0
	github.com/sirupsen/logrus v1.4.2
	github.com/slack-go/slack v0.8.2
	github.com/spf13/cobra v1.1.3
	github.com/spf13/viper v1.7.1
	go.uber.org/zap v1.10.0 // indirect
	gopkg.in/yaml.v2 v2.4.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	k8s.io/api v0.18.2
	k8s.io/apimachinery v0.18.2
	k8s.io/client-go v12.0.0+incompatible
	sigs.k8s.io/controller-runtime v0.6.0
)

// Pinned to kubernetes-1.16.2
replace (
	k8s.io/api => k8s.io/api v0.0.0-20191016110408-35e52d86657a
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20191004115801-a2eda9f80ab8
	k8s.io/client-go => k8s.io/client-go v0.0.0-20191016111102-bec269661e48
)
