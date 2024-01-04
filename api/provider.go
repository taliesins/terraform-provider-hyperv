package api

type Client interface {
	HypervVhdClient
	HypervVmClient
	HypervVmDvdDriveClient
	HypervVmFirmwareClient
	HypervVmHardDiskDriveClient
	HypervVmIntegrationServiceClient
	HypervVmNetworkAdapterClient
	HypervVmProcessorClient
	HypervVmStatusClient
	HypervVmSwitchClient
	HypervIsoImageClient
}

type Provider struct {
	Client Client
}
