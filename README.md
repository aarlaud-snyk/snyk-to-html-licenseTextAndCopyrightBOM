# Small client to extract bill of material from Snyk API into an html

Pull the package name and version, usage count, and license + legal text into an html file.

## Prerequisites
1. Your Snyk org ID => Check settings to find the Org ID
2. Your Snyk token => account => settings
3. A paid Snyk account (API are a paid thing. Check snyk.io/plans)

## Usage
licenseTextAndCopyrightBOM -orgID=<your_org_ID> -token=<your_snyk_token>

Generates an output.html file locally

Optional for Snyk Private instances  
`-api https://your-instance-hostname/api`

Optional to ignore TLS errors (private CA certs, expired certificates)
`-insecure`

### To do
- Implement some tests
- Set build pipeline
- Add copyrights extraction