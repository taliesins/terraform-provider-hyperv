package hyperv_winrm

import (
	"context"
	"fmt"
	"text/template"

	"github.com/taliesins/terraform-provider-hyperv/api"
)

type createOrUpdateIsoImageArgs struct {
	SourceDirectoryPath string
	SourceBootFilePath  string
	DestinationIsoPath  string
	Media               api.IsoMediaType
	Title               string
}

var createOrUpdateIsoImageTemplate = template.Must(template.New("NewIso").Parse(`
function New-ISOFile {
    [CmdletBinding(SupportsShouldProcess = $true, ConfirmImpact = "Low")]
    Param
    (
        [parameter(Mandatory = $true, ValueFromPipeline = $false)]
        [string]$source,
        [parameter(Mandatory = $true, ValueFromPipeline = $false)]
        [string]$destinationIso,
        [parameter(Mandatory = $false, ValueFromPipeline = $false)]
        [string]$bootFile = $null,
        [Parameter(Mandatory = $false, ValueFromPipeline = $false)]
        [ValidateSet("CDR", "CDRW", "DVDRAM", "DVDPLUSR", "DVDPLUSRW", "DVDPLUSR_DUALLAYER", "DVDDASHR", "DVDDASHRW", "DVDDASHR_DUALLAYER", "DISK", "DVDPLUSRW_DUALLAYER", "BDR", "BDRE")]
        [string]$media = "DVDPLUSRW_DUALLAYER",
        [Parameter(Mandatory = $false, ValueFromPipeline = $false)]
        [string]$title = "untitled",
        [Parameter(Mandatory = $false, ValueFromPipeline = $false)]
        [switch]$force
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

    if ($bootFile) {
        if (@('BDR', 'BDRE') -contains $media) {
            throw ("Selected boot image may not work with BDR/BDRE media types.")
        }

        if (!(Test-Path -Path $bootFile)) {
            throw ($bootFile + " is not valid.")
        } 

        try {
            $stream = New-Object -ComObject ADODB.Stream -Property @{Type = 1 } -ErrorAction Stop
            $stream.Open()
            $stream.LoadFromFile((Get-Item -LiteralPath $bootFile).Fullname)
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

    ## Build array of media types
    $mediaType = @(
        "UNKNOWN",
        "CDROM",
        "CDR",
        "CDRW",
        "DVDROM",
        "DVDRAM",
        "DVDPLUSR",
        "DVDPLUSRW",
        "DVDPLUSR_DUALLAYER",
        "DVDDASHR",
        "DVDDASHRW",
        "DVDDASHR_DUALLAYER",
        "DISK",
        "DVDPLUSRW_DUALLAYER",
        "HDDVDROM",
        "HDDVDR",
        "HDDVDRAM",
        "BDROM",
        "BDR",
        "BDRE"
    )

    try {
        $image = New-Object -ComObject IMAPI2FS.MsftFileSystemImage -Property @{VolumeName = $title } -ErrorAction Stop
        $image.ChooseImageDefaultsForMediaType($mediaType.IndexOf($media))
    }
    catch {
        throw ("Failed to initialise image. " + $_.exception.Message)
    }

    ## Create target ISO, throw if file exists and -force parameter is not used.
    if ($PSCmdlet.ShouldProcess($destinationIso)) {
        if (!($targetFile = New-Item -Path $destinationIso -ItemType File -Force:$Force -ErrorAction SilentlyContinue)) {
            throw ("Cannot create file " + $destinationIso + ". Use -Force parameter to overwrite if the target file already exists.")
        }
    }

    try {
        $sourceItems = Get-ChildItem -LiteralPath $source -ErrorAction Stop
    }
    catch {
        throw ("Failed to get source items. " + $_.exception.message)
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

    return $targetFile
}

New-ISOFile -source "{{.Source}}" -destinationIso "{{.DestinationIso}}" -bootFile "{{.BootFile}}" -media "{{.Media}}" -title "{{.Title}}" -force $true
`))

func (c *ClientConfig) CreateOrUpdateIsoImage(ctx context.Context, sourceDirectoryPath string, sourceBootFilePath string, destinationIsoPath string, excludeList []string, media api.IsoMediaType, title string) (err error) {
	var err1 error
	var err2 error
	var err3 error

	remoteSourceBootFilePath := ""
	if sourceBootFilePath != "" {
		remoteSourceBootFilePath, err1 = c.WinRmClient.UploadFile(ctx, sourceBootFilePath)
	} else {
		remoteSourceBootFilePath = ""
	}

	remoteSourceDirectoryPath := ""
	if err1 == nil {
		remoteSourceDirectoryPath, _, err2 = c.WinRmClient.UploadDirectory(ctx, sourceDirectoryPath, excludeList)
	}

	if err1 == nil && err2 == nil {
		err3 = c.WinRmClient.RunFireAndForgetScript(ctx, createOrUpdateIsoImageTemplate, createOrUpdateIsoImageArgs{
			SourceDirectoryPath: remoteSourceDirectoryPath,
			SourceBootFilePath:  remoteSourceBootFilePath,
			DestinationIsoPath:  destinationIsoPath,
			Media:               media,
			Title:               title,
		})
	}

	if err2 == nil {
		err2 = c.WinRmClient.DeleteFileOrDirectory(ctx, remoteSourceDirectoryPath)
	}

	if err1 == nil {
		if remoteSourceBootFilePath != "" {
			err1 = c.WinRmClient.DeleteFileOrDirectory(ctx, remoteSourceBootFilePath)
		}
	}

	if err3 != nil {
		return fmt.Errorf("error for iso file %s: %v", destinationIsoPath, err3)
	}

	if err2 != nil {
		return fmt.Errorf("error for files %s: %v", remoteSourceDirectoryPath, err2)
	}

	if err3 != nil {
		return fmt.Errorf("error for boot file %s: %v", remoteSourceBootFilePath, err3)
	}

	return nil
}

type deleteIsoImageArgs struct {
	Path string
}

var deleteIsoImageTemplate = template.Must(template.New("DeleteIsoImage").Parse(`
$ErrorActionPreference = 'Stop'

$targetDirectory = (split-path '{{.Path}}' -Parent)
$targetName = (split-path '{{.Path}}' -Leaf)
$targetName = $targetName.Substring(0,$targetName.LastIndexOf('.')).split('\')[-1]

Get-ChildItem -Path $targetDirectory |?{$_.BaseName.StartsWith($targetName)} | %{
	Remove-Item $_.FullName -Force
}
`))

func (c *ClientConfig) DeleteIsoImage(ctx context.Context, path string) (err error) {
	err = c.WinRmClient.RunFireAndForgetScript(ctx, deleteIsoImageTemplate, deleteIsoImageArgs{
		Path: path,
	})

	return err
}
