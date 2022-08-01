package hyperv_winrm

import (
	"context"
	"encoding/json"
	"text/template"

	"github.com/taliesins/terraform-provider-hyperv/api"
)

type existsVhdArgs struct {
	Path string
}

var existsVhdTemplate = template.Must(template.New("ExistsVhd").Parse(`
$ErrorActionPreference = 'Stop'
$path='{{.Path}}'

if (Test-Path $path) {
	$exists = ConvertTo-Json -InputObject @{Exists=$true}
	$exists
} else {
	$exists = ConvertTo-Json -InputObject @{Exists=$false}
	$exists
}
`))

func (c *ClientConfig) VhdExists(ctx context.Context, path string) (result api.VhdExists, err error) {
	err = c.WinRmClient.RunScriptWithResult(ctx, existsVhdTemplate, existsVhdArgs{
		Path: path,
	}, &result)

	return result, err
}

type createOrUpdateVhdArgs struct {
	Source     string
	SourceVm   string
	SourceDisk int
	VhdJson    string
}

var createOrUpdateVhdTemplate = template.Must(template.New("CreateOrUpdateVhd").Parse(`
$ErrorActionPreference = 'Stop'

Import-Module Hyper-V
$source='{{.Source}}'
$sourceVm='{{.SourceVm}}'
$sourceDisk={{.SourceDisk}}
$vhd = '{{.VhdJson}}' | ConvertFrom-Json
$vhdType = [Microsoft.Vhd.PowerShell.VhdType]$vhd.VhdType

function Get-TarPath {
	if (Get-Command "tar" -ErrorAction SilentlyContinue) {
		return "tar"
	} elseif (test-path "$env:SystemRoot\system32\tar.exe") {
		return "$env:SystemRoot\system32\tar.exe"
	} else {
		return ""
	}
}

function Get-7ZipPath {
	if (Get-Command "7z" -ErrorAction SilentlyContinue) {
		return "7z"
	} elseif (test-path "$env:ProgramFiles\7-Zip\7z.exe") {
		return "$env:ProgramFiles\7-Zip\7z.exe"
	} elseif (test-path "${env:ProgramFiles(x86)}\7-Zip\7z.exe") {
		return "${env:ProgramFiles(x86)}\7-Zip\7z.exe"
	} else {
		return ""
	}
}

function Expand-Downloads {
    param(
        [Parameter(Mandatory = $true, Position = 0)]
        [string]
        [Alias('Folder')]
        $FolderPath
    )
    process {
		Push-Location $FolderPath

        get-item *.zip | % {
			$tempPath = join-path $FolderPath "temp"

			$7zPath = Get-7ZipPath
			if ($7zPath) {
				$command = """$7zPath"" x ""$($_.FullName)"" -o""$tempPath""" 
				& cmd.exe /C $command
			} else {
				Add-Type -AssemblyName System.IO.Compression.FileSystem
    			if (!(Test-Path $tempPath)) {
        			New-Item -ItemType Directory -Force -Path $tempPath
    			}
            	[System.IO.Compression.ZipFile]::ExtractToDirectory($_.FullName, $tempPath)
			}

            if (Test-Path "$tempPath\Virtual Hard Disks") {
        		Move-Item "$tempPath\Virtual Hard Disks\*.*" $FolderPath
			} else {
				Move-Item "$tempPath\*.*" $FolderPath
			}

			Remove-Item $tempPath -Force -Recurse
			Remove-Item $_.FullName -Force
        }

        get-item *.7z | % {
			$7zPath = Get-7ZipPath
			if (-not $7zPath) {
 				throw "7z.exe needed"
			}
			$tempPath = join-path $FolderPath "temp"
			$command = """$7zPath"" x ""$($_.FullName)"" -o""$tempPath""" 
			& cmd.exe /C $command

			if (Test-Path "$tempPath\Virtual Hard Disks") {
        		Move-Item "$tempPath\Virtual Hard Disks\*.*" $FolderPath
			} else {
				Move-Item "$tempPath\*.*" $FolderPath
			}

			Remove-Item $tempPath -Force -Recurse
			Remove-Item $_.FullName -Force
        }

        get-item *.box | % {
			$tarPath = Get-TarPath
			if (-not $tarPath) {
				throw "tar.exe needed"
			}
			$tempPath = join-path $FolderPath "temp"

			if (!(Test-Path $tempPath)) {
				New-Item -ItemType Directory -Force -Path $tempPath
			}
			$command = """$tarPath"" -C ""$tempPath"" -x -f ""$($_.FullName)"""
			& cmd.exe /C $command

			if (Test-Path "$tempPath\Virtual Hard Disks") {
        		Move-Item "$tempPath\Virtual Hard Disks\*.*" $FolderPath
			} else {
				Move-Item "$tempPath\*.*" $FolderPath
			}

			Remove-Item $tempPath -Force -Recurse
			Remove-Item $_.FullName -Force
        }

		Pop-Location
    }
}

function Get-FileFromUri {
    param(
        [Parameter(Mandatory = $true, Position = 0, ValueFromPipeline = $true, ValueFromPipelineByPropertyName = $true)]
        [string]
        [Alias('Uri')]
        $Url,
        [Parameter(Mandatory = $false, Position = 1)]
        [string]
        [Alias('Folder')]
        $FolderPath
    )
    process {
        $req = [System.Net.HttpWebRequest]::Create($Url)
        $req.Method = "HEAD"
        $response = $req.GetResponse()
        $fUri = $response.ResponseUri
        $filename = [System.IO.Path]::GetFileName($fUri.LocalPath)
        $response.Close()

        $origExt = [System.IO.Path]::GetExtension($Url)
        $newExt = [System.IO.Path]::GetExtension($filename)
        if ($newExt -ne $origExt) {
            $filename += $origExt
        }

        $destination = (Get-Item -Path ".\" -Verbose).FullName
        if ($FolderPath) { $destination = $FolderPath }
        if ($destination.EndsWith('\')) {
            $destination += $filename
        }
        else {
            $destination += '\' + $filename
        }
        $webclient = New-Object System.Net.WebClient
        $webclient.DownloadFile($fUri.AbsoluteUri, $destination)
    }
}

function Test-Uri {
    param(
        [Parameter(Mandatory = $true, Position = 0, ValueFromPipeline = $true, ValueFromPipelineByPropertyName = $true)]
        [string]
        [Alias('Uri')]
        $Url
    )
    process {
        $testUri = $Url -as [System.URI]
        $null -ne $testUri.AbsoluteURI -and $testUri.Scheme -match '[http|https]' -and ($testUri.ToString().ToLower().StartsWith("http://") -or $testUri.ToString().ToLower().StartsWith("https://"))
    }
}

if (!(Test-Path -Path $vhd.Path)) {
    $pathDirectory = [System.IO.Path]::GetDirectoryName($vhd.Path)
    $pathFilename = [System.IO.Path]::GetFileName($vhd.Path)

    if (!(Test-Path $pathDirectory)) {
        New-Item -ItemType Directory -Force -Path $pathDirectory
    }

    if ($sourceVm) {
        Export-VM -Name $sourceVm -Path $pathDirectory
        $targetName = (split-path $vhd.Path -Leaf)
        $targetName = $targetName.Substring(0,$targetName.LastIndexOf('.')).split('\')[-1]
        Get-ChildItem -Path "$pathDirectory\$sourceVm\Virtual Hard Disks" |?{$_.BaseName.StartsWith($sourceVm)} | %{
            $targetNamePath = "$($pathDirectory)\$($_.Name.Replace($sourceVm, $targetName))"
            Move-Item $_.FullName $targetNamePath
        }

        Remove-Item "$pathDirectory\$sourceVm" -Force -Recurse
        Get-VHD -path $vhd.Path
    } elseif ($source) {
        Push-Location $pathDirectory
        
        if (Test-Uri -Url $source) {
            Get-FileFromUri -Url $source -FolderPath $pathDirectory
        }
        else {
            Copy-Item $source "$pathDirectory\$pathFilename" -Force
        }

        Expand-Downloads -FolderPath $pathDirectory

        Pop-Location
    } else {
        $NewVhdArgs = @{}
        $NewVhdArgs.Path = $vhd.Path

        if ($sourceDisk) {
            $NewVhdArgs.SourceDisk = $sourceDisk
        }
        elseif ($vhdType -eq [Microsoft.Vhd.PowerShell.VhdType]::Differencing) {
            $NewVhdArgs.Differencing = $true
            $NewVhdArgs.ParentPath = $vhd.ParentPath
        }
        else {
            if ($vhdType -eq [Microsoft.Vhd.PowerShell.VhdType]::Dynamic) {
                $NewVhdArgs.Dynamic = $true
            }
            elseif ($vhdType -eq [Microsoft.Vhd.PowerShell.VhdType]::Fixed) {
                $NewVhdArgs.Fixed = $true
            }

            if ($vhd.Size -gt 0) {
                $NewVhdArgs.SizeBytes = $vhd.Size
            }

            if ($vhd.BlockSize -gt 0) {
                $NewVhdArgs.BlockSizeBytes = $vhd.BlockSize
            }

            if ($vhd.LogicalSectorSize -gt 0) {
                $NewVhdArgs.LogicalSectorSizeBytes = $vhd.LogicalSectorSize
            }

            if ($vhd.PhysicalSectorSize -gt 0) {
                $NewVhdArgs.PhysicalSectorSizeBytes = $vhd.PhysicalSectorSize
            }
        }

        New-VHD @NewVhdArgs
    }
}
`))

func (c *ClientConfig) CreateOrUpdateVhd(ctx context.Context, path string, source string, sourceVm string, sourceDisk int, vhdType api.VhdType, parentPath string, size uint64, blockSize uint32, logicalSectorSize uint32, physicalSectorSize uint32) (err error) {
	vhdJson, err := json.Marshal(api.Vhd{
		Path:               path,
		VhdType:            vhdType,
		ParentPath:         parentPath,
		Size:               size,
		BlockSize:          blockSize,
		LogicalSectorSize:  logicalSectorSize,
		PhysicalSectorSize: physicalSectorSize,
	})

	if err != nil {
		return err
	}

	err = c.WinRmClient.RunFireAndForgetScript(ctx, createOrUpdateVhdTemplate, createOrUpdateVhdArgs{
		Source:     source,
		SourceVm:   sourceVm,
		SourceDisk: sourceDisk,
		VhdJson:    string(vhdJson),
	})

	return err
}

type resizeVhdArgs struct {
	Path string
	Size uint64
}

var resizeVhdTemplate = template.Must(template.New("ResizeVhd").Parse(`
$ErrorActionPreference = 'Stop'
$vhd = Get-VHD -Path '{{.Path}}'
if ($vhd.Size -ne {{.Size}}){
	Resize-VHD -Path '{{.Path}}' -SizeBytes {{.Size}}
}
`))

func (c *ClientConfig) ResizeVhd(ctx context.Context, path string, size uint64) (err error) {
	err = c.WinRmClient.RunFireAndForgetScript(ctx, resizeVhdTemplate, resizeVhdArgs{
		Path: path,
		Size: size,
	})

	return err
}

type getVhdArgs struct {
	Path string
}

var getVhdTemplate = template.Must(template.New("GetVhd").Parse(`
$ErrorActionPreference = 'Stop'
$path='{{.Path}}'

$vhdObject = $null
if (Test-Path $path) {
	$vhdObject = Get-VHD -path $path | %{ @{
		Path=$_.Path;
		BlockSize=$_.BlockSize;
		LogicalSectorSize=$_.LogicalSectorSize;
		PhysicalSectorSize=$_.PhysicalSectorSize;
		ParentPath=$_.ParentPath;
		FileSize=$_.FileSize;
		Size=$_.Size;
		MinimumSize=$_.MinimumSize;
		Attached=$_.Attached;
		DiskNumber=$_.DiskNumber;
		Number=$_.Number;
		FragmentationPercentage=$_.FragmentationPercentage;
		Alignment=$_.Alignment;
		DiskIdentifier=$_.DiskIdentifier;
		VhdType=$_.VhdType;
		VhdFormat=$_.VhdFormat;
	}}
}

if ($vhdObject){
	$vhd = ConvertTo-Json -InputObject $vhdObject
	$vhd
} else {
	"{}"
}
`))

func (c *ClientConfig) GetVhd(ctx context.Context, path string) (result api.Vhd, err error) {
	err = c.WinRmClient.RunScriptWithResult(ctx, getVhdTemplate, getVhdArgs{
		Path: path,
	}, &result)

	return result, err
}

type deleteVhdArgs struct {
	Path string
}

var deleteVhdTemplate = template.Must(template.New("DeleteVhd").Parse(`
$ErrorActionPreference = 'Stop'

$targetDirectory = (split-path '{{.Path}}' -Parent)
$targetName = (split-path '{{.Path}}' -Leaf)
$targetName = $targetName.Substring(0,$targetName.LastIndexOf('.')).split('\')[-1]

Get-ChildItem -Path $targetDirectory |?{$_.BaseName.StartsWith($targetName)} | %{
	Remove-Item $_.FullName -Force
}
`))

func (c *ClientConfig) DeleteVhd(ctx context.Context, path string) (err error) {
	err = c.WinRmClient.RunFireAndForgetScript(ctx, deleteVhdTemplate, deleteVhdArgs{
		Path: path,
	})

	return err
}
