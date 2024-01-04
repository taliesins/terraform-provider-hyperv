cd C:\Data\terraform-provider-hyperv
del %appdata%\terraform.d\plugins\registry.terraform.io\taliesins\hyperv\1.0.3\windows_amd64\terraform-provider-hyperv_1.0.3.exe
go build -o %appdata%\terraform.d\plugins\registry.terraform.io\taliesins\hyperv\1.0.3\windows_amd64\terraform-provider-hyperv_1.0.3.exe
cd C:\Data\terraform-provider-hyperv\examples\resources\hyperv_iso_image
del .terraform.lock.hcl /Q
RMDIR .terraform /S /Q
#set TF_LOG=TRACE
#set TF_LOG_CORE=TRACE
set TF_LOG_PROVIDER=TRACE
#set WINRMCP_DEBUG=1
terraform init