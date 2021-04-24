module github.com/flacatus/che-inspector

go 1.14

require (
	github.com/Microsoft/go-winio v0.4.17 // indirect
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/docker/docker v0.7.3-0.20190327010347-be7ac8be2ae0
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/go-playground/universal-translator v0.17.0 // indirect
	github.com/go-playground/validator v9.31.0+incompatible
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.0.1 // indirect
	github.com/sirupsen/logrus v1.7.0
	github.com/slack-go/slack v0.8.2 // indirect
	github.com/spf13/cobra v1.1.3
	github.com/spf13/viper v1.7.1
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	gopkg.in/yaml.v2 v2.4.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	k8s.io/api v0.21.0
	k8s.io/apimachinery v0.21.0
	k8s.io/client-go v12.0.0+incompatible
	k8s.io/component-base v0.18.2
	sigs.k8s.io/controller-runtime v0.6.0
)

// Pinned to kubernetes-1.16.2
replace (
	k8s.io/api => k8s.io/api v0.0.0-20210416194706-86cef11b7287
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20210415154527-1ba67c107540
	k8s.io/client-go => k8s.io/client-go v0.0.0-20210416194932-d974964d1226
	k8s.io/component-base => k8s.io/component-base v0.0.0-20210412032905-a57cc3fac704
)
