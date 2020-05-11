cd cmd
go build
cd ..
mv cmd/cmd httplogger

export JAEGER_DISABLED='true'

./httplogger