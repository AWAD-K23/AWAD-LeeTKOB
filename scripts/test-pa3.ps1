# Run from the repository root. Read local settings without printing secrets.
$ErrorActionPreference = 'Stop'

if (Test-Path -LiteralPath '.env') {
    foreach ($line in Get-Content -LiteralPath '.env') {
        if ($line -match '^([A-Z_]+)=(.*)$') {
            [Environment]::SetEnvironmentVariable($matches[1], $matches[2].Trim('"', "'"), 'Process')
        }
    }
}

if (-not $env:TEST_DATABASE_URL) { $env:TEST_DATABASE_URL = $env:DATABASE_URL }
if (-not $env:TEST_DATABASE_URL) { throw 'Set TEST_DATABASE_URL in your ignored .env file.' }
$env:SANDBOX_TESTS = '1'

docker build -t leetkob-sandbox:pa3 sandbox
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }

node scripts/check-go-format.mjs
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }

go -C backend vet ./...
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }

go -C backend test ./... -count=1 -v
exit $LASTEXITCODE
