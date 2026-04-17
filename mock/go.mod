module gitee.com/liujit/shop/mock

go 1.23.3

replace gitee.com/liujit/shop/server => ../server

require (
	gitee.com/liujit/shop/server v0.0.0-00010101000000-000000000000
	go.newcapec.cn/nctcommon/nmslib v0.0.7
	go.newcapec.cn/ncttools/nmskit-bootstrap/conf v0.2.7
	go.newcapec.cn/ncttools/nmskit-bootstrap/oss v0.0.1
	go.newcapec.cn/ncttools/nmskit-bootstrap/sqldb v0.0.3
	gorm.io/driver/mysql v1.5.7
	gorm.io/gen v0.3.26
	gorm.io/gorm v1.25.12
)

require (
	cloud.google.com/go v0.110.7 // indirect
	cloud.google.com/go/bigquery v1.53.0 // indirect
	cloud.google.com/go/compute/metadata v0.6.0 // indirect
	cloud.google.com/go/iam v1.1.1 // indirect
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/andybalholm/brotli v1.0.5 // indirect
	github.com/apache/arrow/go/v12 v12.0.0 // indirect
	github.com/apache/thrift v0.21.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/emirpasic/gods v1.18.1 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-playground/form/v4 v4.2.1 // indirect
	github.com/go-sql-driver/mysql v1.9.2 // indirect
	github.com/goccy/go-json v0.10.3 // indirect
	github.com/godoes/gorm-oracle v1.6.14 // indirect
	github.com/golang-jwt/jwt/v5 v5.2.1 // indirect
	github.com/golang-sql/civil v0.0.0-20220223132316-b832511892a9 // indirect
	github.com/golang-sql/sqlexp v0.1.0 // indirect
	github.com/golang/groupcache v0.0.0-20241129210726-2c02b8208cf8 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/flatbuffers v2.0.8+incompatible // indirect
	github.com/google/gnostic v0.7.0 // indirect
	github.com/google/gnostic-models v0.6.9-0.20230804172637-c7be7c783f49 // indirect
	github.com/google/s2a-go v0.1.4 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.2.3 // indirect
	github.com/googleapis/gax-go/v2 v2.11.0 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.4.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20231201235250-de7065d80cb9 // indirect
	github.com/jackc/pgx/v5 v5.5.5 // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/klauspost/asmfmt v1.3.2 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/klauspost/cpuid/v2 v2.2.8 // indirect
	github.com/mattn/go-sqlite3 v1.14.22 // indirect
	github.com/microsoft/go-mssqldb v1.7.2 // indirect
	github.com/minio/asm2plan9s v0.0.0-20200509001527-cdd76441f9d8 // indirect
	github.com/minio/c2goasm v0.0.0-20190812172519-36a3d3bbc4f3 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/pierrec/lz4/v4 v4.1.21 // indirect
	github.com/prometheus/client_golang v1.22.0 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.62.0 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	github.com/redis/go-redis/extra/rediscmd/v9 v9.7.0 // indirect
	github.com/redis/go-redis/extra/redisotel/v9 v9.7.0 // indirect
	github.com/redis/go-redis/v9 v9.7.0 // indirect
	github.com/shamsher31/goimgext v1.0.0 // indirect
	github.com/sijms/go-ora/v2 v2.8.21 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/thoas/go-funk v0.9.3 // indirect
	github.com/zeebo/xxh3 v1.0.2 // indirect
	go.newcapec.cn/mirror/dameng v0.0.1 // indirect
	go.newcapec.cn/mirror/dm v0.1.0 // indirect
	go.newcapec.cn/ncttools/nmskit v0.1.3 // indirect
	go.newcapec.cn/ncttools/nmskit-auth v0.1.1 // indirect
	go.newcapec.cn/ncttools/nmskit-auth/authn v0.0.1 // indirect
	go.newcapec.cn/ncttools/nmskit-auth/authn/middleware v0.0.1 // indirect
	go.newcapec.cn/ncttools/nmskit-auth/authz v0.0.1 // indirect
	go.newcapec.cn/ncttools/nmskit-auth/authz/middleware v0.0.1 // indirect
	go.newcapec.cn/ncttools/nmskit-bootstrap/cache v0.0.6 // indirect
	go.newcapec.cn/ncttools/nmskit-bootstrap/database/gorm v0.0.3 // indirect
	go.newcapec.cn/ncttools/nmskit-bootstrap/database/gorm/driver v0.0.1 // indirect
	go.newcapec.cn/ncttools/nmskit-bootstrap/database/gorm/driver/bigquery v0.0.3 // indirect
	go.newcapec.cn/ncttools/nmskit-bootstrap/database/gorm/driver/dameng v0.0.3 // indirect
	go.newcapec.cn/ncttools/nmskit-bootstrap/database/gorm/driver/mysql v0.0.1 // indirect
	go.newcapec.cn/ncttools/nmskit-bootstrap/database/gorm/driver/oracle v0.0.3 // indirect
	go.newcapec.cn/ncttools/nmskit-bootstrap/database/gorm/driver/postgres v0.0.3 // indirect
	go.newcapec.cn/ncttools/nmskit-bootstrap/database/gorm/driver/sqlite v0.0.3 // indirect
	go.newcapec.cn/ncttools/nmskit-bootstrap/database/gorm/driver/sqlserver v0.0.3 // indirect
	go.newcapec.cn/ncttools/nmskit-bootstrap/utils v0.0.2 // indirect
	go.newcapec.cn/ncttools/nmskit/log v0.1.2 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/otel v1.34.0 // indirect
	go.opentelemetry.io/otel/metric v1.34.0 // indirect
	go.opentelemetry.io/otel/trace v1.34.0 // indirect
	golang.org/x/crypto v0.37.0 // indirect
	golang.org/x/mod v0.24.0 // indirect
	golang.org/x/net v0.39.0 // indirect
	golang.org/x/oauth2 v0.26.0 // indirect
	golang.org/x/sync v0.13.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
	golang.org/x/text v0.24.0 // indirect
	golang.org/x/tools v0.32.0 // indirect
	golang.org/x/xerrors v0.0.0-20220907171357-04be3eba64a2 // indirect
	google.golang.org/api v0.126.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20230822172742-b8732ec3820d // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250428153025-10db94c68c34 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250425173222-7b384671a197 // indirect
	google.golang.org/grpc v1.72.0 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	gorm.io/datatypes v1.2.4 // indirect
	gorm.io/driver/bigquery v1.2.0 // indirect
	gorm.io/driver/postgres v1.5.0 // indirect
	gorm.io/driver/sqlite v1.5.7 // indirect
	gorm.io/driver/sqlserver v1.5.4 // indirect
	gorm.io/hints v1.1.0 // indirect
	gorm.io/plugin/dbresolver v1.5.3 // indirect
	gorm.io/plugin/opentelemetry v0.1.11 // indirect
	gorm.io/plugin/prometheus v0.1.0 // indirect
)
