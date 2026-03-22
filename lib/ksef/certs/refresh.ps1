#!/usr/bin/env pwsh

$PROD_BASE_URL="https://api.ksef.mf.gov.pl/api/v2"
$DEMO_BASE_URL="https://api-demo.ksef.mf.gov.pl/api/v2"
$TEST_BASE_URL="https://api-test.ksef.mf.gov.pl/api/v2"

$mapping = @{
  "prod.crt" = $PROD_BASE_URL
  "test.crt" = $TEST_BASE_URL
  "demo.crt" = $DEMO_BASE_URL
}

$ENDPOINT="/security/public-key-certificates"

$EXPECTED_TYPE="SymmetricKeyEncryption"

function Get-Certificate(
  [Parameter(Mandatory = $true)]
  [string]$baseEndpoint
)
{
  $certificates = curl "$PROD_BASE_URL$ENDPOINT" | ConvertFrom-Json

  foreach ($cert in $certificates)
  {
    if ($cert.usage -eq $EXPECTED_TYPE)
    {
      return $cert
    }
  }

  return $null
}

function ConvertTo-Pem(
  [Parameter(Mandatory = $true)]
  [string]$derCertificate
)
{
  return [Convert]::FromBase64String($derCertificate) | openssl x509 
}

foreach ($key in $mapping.Keys)
{
  $baseUrl = $mapping[$key]
  $cert = Get-Certificate($baseUrl)

  if ($null -eq $cert)
  {
    write-error "$EXPECTED_TYPE certificate not found in response!"
    exit 1
  }

  ConvertTo-Pem($cert.certificate) | Out-File "$PSScriptRoot/$key"
}
