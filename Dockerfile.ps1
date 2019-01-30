$WIN_PATH=Convert-Path .

#Convert for docker mount to be OK on Windows 10 and Windows 7 Powershell
# Exact conversion is : remove the ":" symbol, replace all "\" by "/", remove last "/" and minor case only the disk letter
# Then for Windows 7, add a /host_mnt/" at the begin of string => this way : c:\Users is translated to /host_mnt/c/Users
# For Windows 10, add "//" => c:\Users is translated to //c/Users
$MOUNT_PATH=(($WIN_PATH -replace "\\","/") -replace ":","").Trim("/")

[regex]$regex='^[a-zA-Z]/'
$MOUNT_PATH=$regex.Replace($MOUNT_PATH, {$args[0].Value.ToLower()})

#Win 10
if ([Environment]::OSVersion.Version -ge (new-object 'Version' 10,0)) {
$MOUNT_PATH="//$MOUNT_PATH"
}
elseif ([Environment]::OSVersion.Version -ge (new-object 'Version' 6,1)) {
$MOUNT_PATH="/host_mnt/$MOUNT_PATH"
}

$PROJECT_PATH="/go/src/github.com/taliesins/terraform-provider-hyperv"

docker run -d -it --name terraform-hyperv --entrypoint bash -v "${MOUNT_PATH}/examples:${PROJECT_PATH}/examples" --workdir $PROJECT_PATH terraform-hyperv
