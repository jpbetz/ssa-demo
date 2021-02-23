module github.com/jpbetz/ssademo

go 1.15

require (
	github.com/go-logr/logr v0.4.0
	k8s.io/apimachinery v0.20.4
	sigs.k8s.io/controller-runtime v0.8.2
	sigs.k8s.io/controller-tools v0.0.0-00010101000000-000000000000 // indirect
)

replace (
	sigs.k8s.io/controller-runtime => ../controller-runtime
	sigs.k8s.io/controller-tools => ../controller-tools
)
