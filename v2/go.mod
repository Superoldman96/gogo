module github.com/chainreactors/gogo/v2

go 1.10

require (
	github.com/M09ic/go-ntlmssp v0.0.0-20230312133735-dcccd454dfe0
	github.com/chainreactors/files v0.0.0-20240716182835-7884ee1e77f0
	github.com/chainreactors/fingers v1.0.1-0.20240730173434-48ba7446b94b
	github.com/chainreactors/logs v0.0.0-20240207121836-c946f072f81f
	github.com/chainreactors/neutron v0.0.0-20240924094734-5be76870bfa2
	github.com/chainreactors/parsers v0.0.0-20241107171801-3f8b998dedc1
	github.com/chainreactors/utils v0.0.0-20241203074749-f12134efcf26
	github.com/google/btree v1.0.0 // indirect
	github.com/jessevdk/go-flags v1.5.0
	github.com/panjf2000/ants/v2 v2.9.1
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519 // indirect
	golang.org/x/net v0.21.0
	gopkg.in/yaml.v3 v3.0.1
	sigs.k8s.io/yaml v1.4.0 // generate only
)

replace (
	golang.org/x/net => golang.org/x/net v0.0.0-20200202094626-16171245cfb2
	golang.org/x/text => golang.org/x/text v0.3.3
)
