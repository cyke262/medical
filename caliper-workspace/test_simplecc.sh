npm init -y

npm install --only=prod @hyperledger/caliper-cli@0.4.2

npx caliper bind --caliper-bind-sut fabric:2.2

npx caliper launch manager --caliper-workspace ./ --caliper-networkconfig networks/networkConfig.json --caliper-benchconfig benchmarks/simpleccBenchmark.yaml --caliper-flow-only-test --caliper-fabric-gateway-enabled --caliper-fabric-gateway-discovery
