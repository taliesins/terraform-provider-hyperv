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
	HypervGpuAdapterClient
	HypervVmStatusClient
	HypervVmSwitchClient
	HypervIsoImageClient
}

type Provider struct {
	Client Client
}
