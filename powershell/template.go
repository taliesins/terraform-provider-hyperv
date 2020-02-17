package powershell

import (
	"strings"
	"text/template"
)

type executePowershellFromCommandLineTemplateOptions struct {
	Powershell	string
}

var executePowershellFromCommandLineTemplate = template.Must(template.New("ExecuteCommandFromCommandLine").Funcs(template.FuncMap{
	"escapeDoubleQuotes": func(textToEscape string) string {
		textToEscape = strings.Replace(textToEscape, "\n", "", -1)
		textToEscape = strings.Replace(textToEscape, "\r", "", -1)
		textToEscape = strings.Replace(textToEscape, "\t", "", -1)
		textToEscape = strings.Replace(textToEscape, `"`, `\"`, -1)
		return textToEscape
	},
}).Parse(`powershell -NoProfile -ExecutionPolicy Bypass "{{escapeDoubleQuotes .Powershell}}"`))

type executeCommandTemplateOptions struct {
	Vars		string
	Path		string
}

var executeCommandTemplate = template.Must(template.New("ExecuteCommand").Parse(`if (Test-Path variable:global:ProgressPreference){$ProgressPreference='SilentlyContinue';};{{.Vars}};&"{{.Path}}";exit $LastExitCode;`))

type elevatedCommandTemplateOptions struct {
	User            			string
	Password        			string
	TaskName        			string
	TaskDescription 			string
	TaskExecutionTimeLimit 		string
	Vars            			string
	ScriptPath  				string
}

var elevatedCommandTemplate = template.Must(template.New("ElevatedCommand").Funcs(template.FuncMap{
	"escapeSingleQuotes": func(textToEscape string) string {
		return strings.Replace(textToEscape, `'`, `''`, -1)
	},
}).Parse(`
function GetTempFile($fileName) {
  $path = $env:TEMP;
  if (!$path){
    $path = 'c:\windows\Temp\';
  }
  return Join-Path -Path $path -ChildPath $fileName;
}

function SlurpStdout($outFile, $currentLine) {
  if (Test-Path $outFile) {
    get-content $outFile | select -skip $currentLine | %{
      $currentLine += 1;
      Write-Host "$_";
    }
  }
  return $currentLine;
}

function SanitizeFileName($fileName) {
    return $fileName.Replace(' ', '_').Replace('&', 'and').Replace('{', '(').Replace('}', ')').Replace('~', '-').Replace('#', '').Replace('%', '');
}

function RunAsScheduledTask($username, $password, $taskName, $taskDescription, $taskExecutionTimeLimit, $vars, $scriptPath)
{
  $stdoutFile = GetTempFile("$(SanitizeFileName($taskName))_stdout.log");
  if (Test-Path $stdoutFile) {
    Remove-Item $stdoutFile | Out-Null;
  }
  $taskXml = @'
<?xml version="1.0" encoding="UTF-16"?>
<Task version="1.2" xmlns="http://schemas.microsoft.com/windows/2004/02/mit/task">
    <RegistrationInfo>
	    <Description>{taskDescription}</Description>
    </RegistrationInfo>
    <Principals>
        <Principal id="Author">
        <UserId>{username}</UserId>
        <LogonType>Password</LogonType>
        <RunLevel>HighestAvailable</RunLevel>
        </Principal>
    </Principals>
    <Settings>
        <MultipleInstancesPolicy>IgnoreNew</MultipleInstancesPolicy>
        <DisallowStartIfOnBatteries>false</DisallowStartIfOnBatteries>
        <StopIfGoingOnBatteries>false</StopIfGoingOnBatteries>
        <AllowHardTerminate>true</AllowHardTerminate>
        <StartWhenAvailable>false</StartWhenAvailable>
        <RunOnlyIfNetworkAvailable>false</RunOnlyIfNetworkAvailable>
        <IdleSettings>
        <StopOnIdleEnd>false</StopOnIdleEnd>
        <RestartOnIdle>false</RestartOnIdle>
        </IdleSettings>
        <AllowStartOnDemand>true</AllowStartOnDemand>
        <Enabled>true</Enabled>
        <Hidden>false</Hidden>
        <RunOnlyIfIdle>false</RunOnlyIfIdle>
        <WakeToRun>false</WakeToRun>
        <ExecutionTimeLimit>{taskExecutionTimeLimit}</ExecutionTimeLimit>
        <Priority>4</Priority>
    </Settings>
    <Actions Context="Author">
        <Exec>
        <Command>cmd</Command>
        <Arguments>{arguments}</Arguments>
        </Exec>
    </Actions>
</Task>
'@;
  $powershellToExecute = 'if (Test-Path variable:global:ProgressPreference){$ProgressPreference=''SilentlyContinue''};' + $vars + ';&"' + $scriptPath + '";exit $LastExitCode;';
  $powershellToExecute = $powershellToExecute.Replace('"', '\"');

  $arguments = '/C powershell -NoProfile -ExecutionPolicy Bypass "' + $powershellToExecute + '" *> "' + $stdoutFile + '"';
  $taskXml = $taskXml.Replace("{arguments}", $arguments.Replace('&', '&amp;').Replace('<', '&lt;').Replace('>', '&gt;').Replace('"', '&quot;').Replace('''', '&apos;'));
  $taskXml = $taskXml.Replace("{username}", $username.Replace('&', '&amp;').Replace('<', '&lt;').Replace('>', '&gt;').Replace('"', '&quot;').Replace('''', '&apos;'));
  $taskXml = $taskXml.Replace("{taskDescription}", $taskDescription.Replace('&', '&amp;').Replace('<', '&lt;').Replace('>', '&gt;').Replace('"', '&quot;').Replace('''', '&apos;'));
  $taskXml = $taskXml.Replace("{taskExecutionTimeLimit}", $taskExecutionTimeLimit.Replace('&', '&amp;').Replace('<', '&lt;').Replace('>', '&gt;').Replace('"', '&quot;').Replace('''', '&apos;'));

  $schedule = New-Object -ComObject "Schedule.Service";
  $schedule.Connect();
  $task = $schedule.NewTask($null);
  $task.XmlText = $taskXml;

  $folder = $schedule.GetFolder('\');
  $folder.RegisterTaskDefinition($taskName, $task, 6, $username, $password, 1, $null) | Out-Null;
  $registeredTask = $folder.GetTask("\$taskName");
  $registeredTask.Run($null) | Out-Null;
  $timeout = 10;
  $sec = 0;
  while ((!($registeredTask.state -eq 4)) -and ($sec -lt $timeout)) {
    Start-Sleep -s 1;
    $sec++;
  }
  $stdoutCurrentLine = 0;
  do {
    Start-Sleep -m 100;
    $stdoutCurrentLine = SlurpStdout $stdoutFile $stdoutCurrentLine;
  } while (!($registeredTask.state -eq 3))
  Start-Sleep -m 100;
  $exit_code = $registeredTask.LastTaskResult;
  $stdoutCurrentLine = SlurpStdout $stdoutFile $stdoutCurrentLine;

  if (Test-Path $stdoutFile) {
    #Remove-Item $stdoutFile -ErrorAction SilentlyContinue | Out-Null;
  }

  if (Test-Path $scriptPath) {
    #Remove-Item $scriptPath -ErrorAction SilentlyContinue | Out-Null;
  }

  $folder.DeleteTask($taskName, 0) | Out-Null;
  [System.Runtime.Interopservices.Marshal]::ReleaseComObject($schedule) | Out-Null;

  return $exit_code;
}

$username = '{{escapeSingleQuotes .User}}'.Replace('\.\\', $env:computername+'\');
$password = '{{escapeSingleQuotes .Password}}';
$taskName = '{{escapeSingleQuotes .TaskName}}';
$taskDescription = '{{escapeSingleQuotes .TaskDescription}}';
$taskExecutionTimeLimit = '{{escapeSingleQuotes .TaskExecutionTimeLimit}}';
$vars = '{{escapeSingleQuotes .Vars}}';
$scriptPath = '{{escapeSingleQuotes .ScriptPath}}';
$exitCode = RunAsScheduledTask -username $username -password $password -taskName $taskName -taskDescription $taskDescription -taskExecutionTimeLimit $taskExecutionTimeLimit -vars $vars -scriptPath $scriptPath;
exit $exitCode;
`))

type convertBase64FileToTextFileTemplateOptions struct {
	Base64FilePath			string
	FilePath				string
}

var convertBase64FileToTextFileTemplate = template.Must(template.New("ConvertBase64FileToTextFile").Parse(`
if (Test-Path variable:global:ProgressPreference) {
	$ProgressPreference='SilentlyContinue';
}
$base64FilePath = [System.IO.Path]::GetFullPath("{{.Base64FilePath}}");
$filePath = [System.IO.Path]::GetFullPath("{{.FilePath}}".Trim("'"));
if (Test-Path $filePath) {
	if (Test-Path -Path $filePath -PathType container) {
		Exit 1;
	} else {
		Remove-Item $filePath | Out-Null;
	}
} else {
	$destinationFolder = ([System.IO.Path]::GetDirectoryName($filePath));
	New-Item -ItemType directory -Force -ErrorAction SilentlyContinue -Path $destinationFolder | Out-Null;
}

if (Test-Path $base64FilePath) {
	$reader = [System.IO.File]::OpenText($base64FilePath);
	try {
		$writer = [System.IO.File]::OpenWrite($filePath);
		try {
			for(;;) {
				$base64_line = $reader.ReadLine();
				if ($base64_line -eq $null) { 
					break;
				}
				$bytes = [System.Convert]::FromBase64String($base64_line);
				$writer.write($bytes, 0, $bytes.Length);
			}
		} finally {
			$writer.Close();
		}
	} finally {
		$reader.Close();
	}
} else {
	Write-Output $null > $filePath;
}

$filePath;
exit $LastExitCode;
`))

type resolvePathTemplateOptions struct {
	FilePath	string
}

var resolvePathTemplate = template.Must(template.New("ResolvePath").Parse(`if (Test-Path variable:global:ProgressPreference){$ProgressPreference='SilentlyContinue';};[System.IO.Path]::GetFullPath("{{.FilePath}}");exit $LastExitCode;`))

type deleteFileTemplateOptions struct {
	FilePath	string
}

var deleteFileTemplate = template.Must(template.New("DeleteFile").Parse(`if (Test-Path variable:global:ProgressPreference){$ProgressPreference='SilentlyContinue';};if (Test-Path "{{.FilePath}}") {Remove-Item "{{.FilePath}}" -ErrorAction SilentlyContinue;};exit $LastExitCode;`))

type appendFileTemplateOptions struct {
	FilePath	string
	Content		string
}

//This is not a Powershell script
var appendFileTemplate = template.Must(template.New("AppendFile").Parse(`echo {{.Content}} >> "{{.FilePath}}"`))
