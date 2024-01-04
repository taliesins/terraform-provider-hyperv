package hyperv_winrm

import (
	"context"
	"encoding/json"
	"fmt"
	"text/template"

	"github.com/taliesins/terraform-provider-hyperv/api"
)

func (c *ClientConfig) RemoteFileUpload(ctx context.Context, filePath string, remoteFilePath string) (err error) {
	_, err = c.WinRmClient.UploadFile(ctx, filePath, remoteFilePath)
	return err
}

func (c *ClientConfig) RemoteFileDelete(ctx context.Context, remoteFilePath string) (err error) {
	err = c.WinRmClient.DeleteFileOrDirectory(ctx, remoteFilePath)
	return err
}

func (c *ClientConfig) RemoteFileExists(ctx context.Context, remoteFilePath string) (exists bool, err error) {
	exists, err = c.WinRmClient.FileExists(ctx, remoteFilePath)
	return exists, err
}

type createOrUpdateIsoImageArgs struct {
	IsoImageJson string
}

var createOrUpdateIsoImageTemplate = template.Must(template.New("CreateOrUpdateIsoImage").Parse(`
$ErrorActionPreference = 'Stop'
$isoImageJson = @'
{{.IsoImageJson}}
'@ 
$isoImage = $isoImageJson | ConvertFrom-Json

$mediaType = @{}

$fileSystemType = @{}

function New-TemporaryDirectory {
  $parent = [System.IO.Path]::GetTempPath()
  do {
    $name = [System.IO.Path]::GetRandomFileName()
    $item = New-Item -Path $parent -Name $name -ItemType "directory" -ErrorAction SilentlyContinue
  } while (-not $item)
  return $item.FullName
}

function Save-IsoImage {
    [CmdletBinding(SupportsShouldProcess = $true, ConfirmImpact = "Low")]
    Param
    (
        [parameter(Mandatory = $false, ValueFromPipeline = $false)]
        [string]$SourceIsoFilePath = "",
        [parameter(Mandatory = $false, ValueFromPipeline = $false)]
        [string]$SourceIsoFilePathHash = "",
        [parameter(Mandatory = $false, ValueFromPipeline = $false)]
        [string]$SourceZipFilePath = "",
        [parameter(Mandatory = $false, ValueFromPipeline = $false)]
        [string]$SourceZipFilePathHash = "",
        [parameter(Mandatory = $false, ValueFromPipeline = $false)]
        [string]$SourceBootFilePath = "",
        [parameter(Mandatory = $false, ValueFromPipeline = $false)]
        [string]$SourceBootFilePathHash = "",
        [parameter(Mandatory = $true, ValueFromPipeline = $false)]
        [string]$DestinationIsoFilePath,
        [parameter(Mandatory = $false, ValueFromPipeline = $false)]
        [string]$DestinationZipFilePath = "",
        [parameter(Mandatory = $false, ValueFromPipeline = $false)]
        [string]$DestinationBootFilePath = "",
        [Parameter(Mandatory = $false, ValueFromPipeline = $false)]
        [ValidateSet(0,0x1,0x2,0x3,0x4,0x5,0x6,0x7,0x8,0x9,0xa,0xb,0xc,0xd,0xe,0xf,0x10,0x11,0x12,0x13)]
        [int]$Media = 0xd,
        [Parameter(Mandatory = $false, ValueFromPipeline = $false)]
        [ValidateSet(0,0x1,0x2,0x3,0x4,0x6,0x7,0x40000000)]
        [int]$FileSystem = 0x40000000,
        [Parameter(Mandatory = $false, ValueFromPipeline = $false)]
        [string]$VolumeName = "UNTITLED",
        [parameter(Mandatory = $true, ValueFromPipeline = $false)]
        [string]$ResolveDestinationIsoFilePath = "",
        [parameter(Mandatory = $false, ValueFromPipeline = $false)]
        [string]$ResolveDestinationZipFilePath = "",
        [parameter(Mandatory = $false, ValueFromPipeline = $false)]
        [string]$ResolveDestinationBootFilePath = "",
        [Parameter(Mandatory = $false, ValueFromPipeline = $false)]
        [switch]$Force
    )
    $typeDefinition = @'
        public class ISOFile  {
            public unsafe static void Create(string Path, object Stream, int BlockSize, int TotalBlocks) {
                int bytes = 0;
                byte[] buf = new byte[BlockSize];
                var ptr = (System.IntPtr)(&bytes);
                var o = System.IO.File.OpenWrite(Path);
                var i = Stream as System.Runtime.InteropServices.ComTypes.IStream;

                if (o != null) {
                    while (TotalBlocks-- > 0) {
                        i.Read(buf, BlockSize, ptr); o.Write(buf, 0, bytes);
                    }

                    o.Flush(); o.Close();
                }
            }
        }
'@

    if (!('ISOFile' -as [type])) {

        ## Add-Type works a little differently depending on PowerShell version.
        ## https://docs.microsoft.com/en-us/powershell/module/microsoft.powershell.utility/add-type
        switch ($PSVersionTable.PSVersion.Major) {

            ## 7 and (hopefully) later versions
            { $_ -ge 7 } {
                Add-Type -CompilerOptions "/unsafe" -TypeDefinition $typeDefinition
            }

            ## 5, and only 5. We aren't interested in previous versions.
            5 {
                $compOpts = New-Object System.CodeDom.Compiler.CompilerParameters
                $compOpts.CompilerOptions = "/unsafe"

                Add-Type -CompilerParameters $compOpts -TypeDefinition $typeDefinition
            }

            default {
                ## If it's not 7 or later, and it's not 5, then we aren't doing it.
                throw ("Unsupported PowerShell version.")
            }
        }
    }

	$expandedResolveDestinationIsoFilePath = $ExecutionContext.InvokeCommand.ExpandString($ResolveDestinationIsoFilePath)
    if (!$expandedResolveDestinationIsoFilePath) {
        throw ("must specify a value for ResolveDestinationIsoFilePath")
    }
	$expandedResolveDestinationZipFilePath = $ExecutionContext.InvokeCommand.ExpandString($ResolveDestinationZipFilePath)
	$expandedResolveDestinationBootFilePath = $ExecutionContext.InvokeCommand.ExpandString($ResolveDestinationBootFilePath)

	if (!(Test-Path -Path $expandedResolveDestinationIsoFilePath) -and !$SourceIsoFilePath) {
        if (!$expandedResolveDestinationZipFilePath) {
            throw ("must specify a value for ResolveDestinationZipFilePath if no SourceIsoFilePath is provided")
        }

		if (!(Test-Path -Path $expandedResolveDestinationZipFilePath)) {
			throw ("Could not find $($expandedResolveDestinationZipFilePath) for specified SourceZipFilePath=$($SourceZipFilePath)")
		} 

		if ($SourceBootFilePath) {
			if ($Media -eq 0x11 -or $Media -eq 0x12 -or $Media -eq 0x13) {
				throw ("Selected boot image may not work with BDR/BDRE media types.")
			}

			if (!(Test-Path -Path $expandedResolveDestinationBootFilePath)) {
				throw ("Could not find $($expandedResolveDestinationBootFilePath) for specified SourceBootFilePath=$($SourceBootFilePath)")
			} 
		}

		$expandedResolveDestinationUnzipDirectoryPath = New-TemporaryDirectory
		try
		{
			Expand-Archive -Path $expandedResolveDestinationZipFilePath -DestinationPath $expandedResolveDestinationUnzipDirectoryPath
	
			if ($SourceBootFilePath) {
				try {
					$stream = New-Object -ComObject ADODB.Stream -Property @{Type = 1} -ErrorAction Stop
					$stream.Open()
					$stream.LoadFromFile((Get-Item -LiteralPath $expandedResolveDestinationBootFilePath).Fullname)
				}
				catch {
					throw ("Failed to open boot file. " + $_.exception.message)
				}
	
				try {
					$boot = New-Object -ComObject IMAPI2FS.BootOptions -ErrorAction Stop
					$boot.AssignBootImage($stream)
				}
				catch {
					throw ("Failed to apply boot file. " + $_.exception.message)
				}
			}
	
			try {
				$image = New-Object -ComObject IMAPI2FS.MsftFileSystemImage -Property @{VolumeName = $VolumeName} -ErrorAction Stop
				$image.ChooseImageDefaultsForMediaType($Media)
				if ($FileSystem -ne 0x40000000) {
					$image.FileSystemsToCreate = $FileSystem
				}
			}
			catch {
				throw ("Failed to initialise image. Media=$($Media), FileSystem=$($FileSystem), isoImageJson=$($isoImageJson).  " + $_.exception.Message)
			}
	
			if (!($targetFile = New-Item -Path $expandedResolveDestinationIsoFilePath -ItemType File -Force:$Force -ErrorAction SilentlyContinue)) {
				throw ("Cannot create file " + $expandedResolveDestinationIsoFilePath + ". Use -Force parameter to overwrite if the target file already exists.")
			}
	
			try {
				$sourceItems = Get-ChildItem -LiteralPath $expandedResolveDestinationUnzipDirectoryPath -ErrorAction Stop
			}
			catch {
				throw ("Failed to get source items. ExpandedResolveDestinationUnzipDirectoryPath=$($expandedResolveDestinationUnzipDirectoryPath), isoImageJson=$($isoImageJson). " + $_.exception.message)
			}
	
			foreach ($sourceItem in $sourceItems) {
				try {
					$image.Root.AddTree($sourceItem.FullName, $true)
				}
				catch {
					throw ("Failed to add " + $sourceItem.fullname + ". " + $_.exception.message)
				}
			} 
		
			if ($boot) {
				$Image.BootImageOptions = $boot
			}
		
			try {
				$result = $image.CreateResultImage()
				[ISOFile]::Create($targetFile.FullName, $result.ImageStream, $result.BlockSize, $result.TotalBlocks)
			}
			catch {
				throw ("Failed to write ISO file. " + $_.exception.Message)
			}
		} finally {
			Remove-Item $expandedResolveDestinationUnzipDirectoryPath -Force -Recurse -ErrorAction SilentlyContinue
		}
	}

	Save-IsoImageMetaData -SourceIsoFilePath $SourceIsoFilePath -SourceIsoFilePathHash $SourceIsoFilePathHash -SourceZipFilePath $SourceZipFilePath -SourceZipFilePathHash $SourceZipFilePathHash -SourceBootFilePath $SourceBootFilePath -SourceBootFilePathHash $SourceBootFilePathHash -DestinationIsoFilePath $DestinationIsoFilePath -DestinationZipFilePath $DestinationZipFilePath -DestinationBootFilePath $DestinationBootFilePath -Media $Media -FileSystem $FileSystem -VolumeName $VolumeName -ResolveDestinationIsoFilePath $ResolveDestinationIsoFilePath -ResolveDestinationZipFilePath $ResolveDestinationZipFilePath -ResolveDestinationBootFilePath $ResolveDestinationBootFilePath -Force:$true
}

function Save-IsoImageMetaData {
    [CmdletBinding(SupportsShouldProcess = $true, ConfirmImpact = "Low")]
    Param
    (
        [parameter(Mandatory = $false, ValueFromPipeline = $false)]
        [string]$SourceIsoFilePath = "",
        [parameter(Mandatory = $false, ValueFromPipeline = $false)]
        [string]$SourceIsoFilePathHash = "",
        [parameter(Mandatory = $false, ValueFromPipeline = $false)]
        [string]$SourceZipFilePath = "",
        [parameter(Mandatory = $false, ValueFromPipeline = $false)]
        [string]$SourceZipFilePathHash = "",
        [parameter(Mandatory = $false, ValueFromPipeline = $false)]
        [string]$SourceBootFilePath = "",
        [parameter(Mandatory = $false, ValueFromPipeline = $false)]
        [string]$SourceBootFilePathHash = "",
        [parameter(Mandatory = $true, ValueFromPipeline = $false)]
        [string]$DestinationIsoFilePath,
        [parameter(Mandatory = $false, ValueFromPipeline = $false)]
        [string]$DestinationZipFilePath = "",
        [parameter(Mandatory = $false, ValueFromPipeline = $false)]
        [string]$DestinationBootFilePath = "",
        [Parameter(Mandatory = $false, ValueFromPipeline = $false)]
        [ValidateSet(0,0x1,0x2,0x3,0x4,0x5,0x6,0x7,0x8,0x9,0xa,0xb,0xc,0xd,0xe,0xf,0x10,0x11,0x12,0x13)]
        [int]$Media = 0xd,
        [Parameter(Mandatory = $false, ValueFromPipeline = $false)]
        [ValidateSet(0,0x1,0x2,0x3,0x4,0x6,0x7,0x40000000)]
        [int]$FileSystem = 0x40000000,
        [Parameter(Mandatory = $false, ValueFromPipeline = $false)]
        [string]$VolumeName = "UNTITLED",
        [parameter(Mandatory = $false, ValueFromPipeline = $false)]
        [string]$ResolveDestinationIsoFilePath = "",
        [parameter(Mandatory = $true, ValueFromPipeline = $false)]
        [string]$ResolveDestinationZipFilePath = "",
        [parameter(Mandatory = $false, ValueFromPipeline = $false)]
        [string]$ResolveDestinationBootFilePath = "",
        [Parameter(Mandatory = $false, ValueFromPipeline = $false)]
        [switch]$Force
    )

	$expandedResolveDestinationIsoFilePath = $ExecutionContext.InvokeCommand.ExpandString($ResolveDestinationIsoFilePath)

	$IsoImageMetadata = @{}
	$IsoImageMetadata.SourceIsoFilePath=$SourceIsoFilePath
	$IsoImageMetadata.SourceIsoFilePathHash=$SourceIsoFilePathHash
	$IsoImageMetadata.SourceZipFilePath=$SourceZipFilePath
	$IsoImageMetadata.SourceZipFilePathHash=$SourceZipFilePathHash
	$IsoImageMetadata.SourceBootFilePath=$SourceBootFilePath
	$IsoImageMetadata.SourceBootFilePathHash=$SourceBootFilePathHash
	$IsoImageMetadata.DestinationIsoFilePath=$DestinationIsoFilePath
	$IsoImageMetadata.DestinationZipFilePath=$DestinationZipFilePath
	$IsoImageMetadata.DestinationBootFilePath=$DestinationBootFilePath
	$IsoImageMetadata.Media=$Media
	$IsoImageMetadata.FileSystem=$FileSystem
	$IsoImageMetadata.VolumeName=$VolumeName
	$IsoImageMetadata.ResolveDestinationIsoFilePath=$ResolveDestinationIsoFilePath
	$IsoImageMetadata.ResolveDestinationZipFilePath=$ResolveDestinationZipFilePath
	$IsoImageMetadata.ResolveDestinationBootFilePath=$ResolveDestinationBootFilePath

	$IsoImageMetadata | ConvertTo-Json -depth 100 | Out-File "$($expandedResolveDestinationIsoFilePath).json" -Force:$Force
}

$SaveIsoImageArgs = @{}
$SaveIsoImageArgs.SourceIsoFilePath=$isoImage.SourceIsoFilePath
$SaveIsoImageArgs.SourceIsoFilePathHash=$isoImage.SourceIsoFilePathHash
$SaveIsoImageArgs.SourceZipFilePath=$isoImage.SourceZipFilePath
$SaveIsoImageArgs.SourceZipFilePathHash=$isoImage.SourceZipFilePathHash
$SaveIsoImageArgs.SourceBootFilePath=$isoImage.SourceBootFilePath
$SaveIsoImageArgs.SourceBootFilePathHash=$isoImage.SourceBootFilePathHash
$SaveIsoImageArgs.DestinationIsoFilePath=$isoImage.DestinationIsoFilePath
$SaveIsoImageArgs.DestinationZipFilePath=$isoImage.DestinationZipFilePath
$SaveIsoImageArgs.DestinationBootFilePath=$isoImage.DestinationBootFilePath
$SaveIsoImageArgs.Media=$isoImage.Media
$SaveIsoImageArgs.FileSystem=$isoImage.FileSystem
$SaveIsoImageArgs.VolumeName=$isoImage.VolumeName
$SaveIsoImageArgs.ResolveDestinationIsoFilePath=$isoImage.ResolveDestinationIsoFilePath
$SaveIsoImageArgs.ResolveDestinationZipFilePath=$isoImage.ResolveDestinationZipFilePath
$SaveIsoImageArgs.ResolveDestinationBootFilePath=$isoImage.ResolveDestinationBootFilePath
$SaveIsoImageArgs.Force=$true

Save-IsoImage @SaveIsoImageArgs
`))

func (c *ClientConfig) CreateOrUpdateIsoImage(ctx context.Context, sourceIsoFilePath string, sourceIsoFilePathHash string, sourceZipFilePath string, sourceZipFilePathHash string, sourceBootFilePath string, sourceBootFilePathHash string, destinationIsoFilePath string, destinationZipFilePath string, destinationBootFilePath string, media api.IsoMediaType, fileSystem api.IsoFileSystemType, volumeName string, resolveDestinationIsoFilePath string, resolveDestinationZipFilePath string, resolveDestinationBootFilePath string) (err error) {

	isoImageJson, err := json.Marshal(api.IsoImage{
		SourceIsoFilePath:              sourceIsoFilePath,
		SourceIsoFilePathHash:          sourceIsoFilePathHash,
		SourceZipFilePath:              sourceZipFilePath,
		SourceZipFilePathHash:          sourceZipFilePathHash,
		SourceBootFilePath:             sourceBootFilePath,
		SourceBootFilePathHash:         sourceBootFilePathHash,
		DestinationIsoFilePath:         destinationIsoFilePath,
		DestinationZipFilePath:         destinationZipFilePath,
		DestinationBootFilePath:        destinationBootFilePath,
		Media:                          media,
		FileSystem:                     fileSystem,
		VolumeName:                     volumeName,
		ResolveDestinationIsoFilePath:  resolveDestinationIsoFilePath,
		ResolveDestinationZipFilePath:  resolveDestinationZipFilePath,
		ResolveDestinationBootFilePath: resolveDestinationBootFilePath,
	})

	if err != nil {
		return fmt.Errorf("error converting object to json: %s", err)
	}

	err = c.WinRmClient.RunFireAndForgetScript(ctx, createOrUpdateIsoImageTemplate, createOrUpdateIsoImageArgs{
		IsoImageJson: string(isoImageJson),
	})

	if err != nil {
		return fmt.Errorf("error creating or updating iso image: %v", err)
	}

	return err //This will return the error from deferred functions
}

type getIsoImageArgs struct {
	ResolveDestinationIsoFilePath string
}

var getIsoImageTemplate = template.Must(template.New("GetIsoImage").Parse(`
$ErrorActionPreference = 'Stop'
$ResolveDestinationIsoFilePath='{{.ResolveDestinationIsoFilePath}}'
$isoImageObject = $null

$expandedResolveDestinationIsoFilePath = $ExecutionContext.InvokeCommand.ExpandString($ResolveDestinationIsoFilePath)

$metadataPath="$($expandedResolveDestinationIsoFilePath).json"

if (Test-Path $metadataPath) {
	$metadata = Get-Content -Raw -Path $metadataPath | ConvertFrom-Json

	$isoImageObject=@{}
	$isoImageObject.SourceIsoFilePath=$metadata.SourceIsoFilePath
	$isoImageObject.SourceIsoFilePathHash=$metadata.SourceIsoFilePathHash
	$isoImageObject.SourceZipFilePath=$metadata.SourceZipFilePath
	$isoImageObject.SourceZipFilePathHash=$metadata.SourceZipFilePathHash
	$isoImageObject.SourceBootFilePath=$metadata.SourceBootFilePath
	$isoImageObject.SourceBootFilePathHash=$metadata.SourceBootFilePathHash
	$isoImageObject.DestinationIsoFilePath=$metadata.DestinationIsoFilePath
	$isoImageObject.DestinationZipFilePath=$metadata.DestinationZipFilePath
	$isoImageObject.DestinationBootFilePath=$metadata.DestinationBootFilePath
	$isoImageObject.Media=$metadata.Media
	$isoImageObject.FileSystem=$metadata.FileSystem
	$isoImageObject.VolumeName=$metadata.VolumeName
	$isoImageObject.ResolveDestinationIsoFilePath=$metadata.ResolveDestinationIsoFilePath
	$isoImageObject.ResolveDestinationZipFilePath=$metadata.ResolveDestinationZipFilePath
	$isoImageObject.ResolveDestinationBootFilePath=$metadata.ResolveDestinationBootFilePath
} else {}

if ($isoImageObject){
	$isoImage = ConvertTo-Json -InputObject $isoImageObject
	$isoImage
} else {
	"{}"
}
`))

func (c *ClientConfig) GetIsoImage(ctx context.Context, resolveDestinationIsoFilePath string) (result api.IsoImage, err error) {
	err = c.WinRmClient.RunScriptWithResult(ctx, getIsoImageTemplate, getIsoImageArgs{
		ResolveDestinationIsoFilePath: resolveDestinationIsoFilePath,
	}, &result)

	return result, err
}
