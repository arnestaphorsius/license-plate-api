# license-plate-api

## Cross compiling
To build the binary for raspberry pi on your local machine use the following command
```
  env GOOS=linux GOARCH=arm go build
```

## Runnning
To run the binary on the raspberry pi issue the following:
```
  sudo ./license-plate-api <pinnumber>
```
where the pinnumber corresponds to the gpio pin numbering, so for physical pin number 7 you would use pinnumber 4.

See https://pinout.xyz/ for the complete pinout
