# go_engage_eoy

Go version of the Engage EOY project.  Provides a spreadsheet that details end-of-year data for an Engage client.

# Background

Salsa's Engage product does not yet have extensive reporting.  Because of that, Engage clients are missing out on
the typical end-of-year reporting that occurs after every Giving Tuesday.

# Summary

This application accepts an Engage API Token, then provides the usual end-of-year reports that Salsa's clients want.  The 
report is in the form of an Excel spreadsheet.  Each sheet in the spreadsheet represents a single set of metrics.
- Current Year
- Months in the current year
- All Years
- Each month and all years
- All donors
- Top donors (nominally the top 20)
- Donation form name

Each row of metrics contains identifiers about the metric (year, month, donor, etc.) and a common set of statistics.

These are the common statistics collected for each metric.

-   Number of donations
-   Total donations
-   Number of one-time donations
-   Total one-time donations
-   Number of recurring donation payments
-   Total recurring donations
-   Number of offline donations
-   Total offline donations
-   Number of refunds
-   Total refunds
-   Largest donation
-   Smallest donation
-   Average donation

# Installation

## Summary

1. Install the Go language if it's not installed.
1. Create the requred directory structure.
1. Add ~/go/bin to the PATH variable.
1. Install the app.
1. Resolve dependencies.
1. Build the executable.

## Prerequisites
The only prerequisite is the most recent version of the Go language.  If you already have Go installed, then skip to "Environment variables" (below).

You can install Go by a variety of methods.  Please [click here](https://golang.org/dl/)
to see the official download page.

## Directory
The next step is to create the correct directory hierarchy.  This *must* appear
in your home directory on your computer.
```
HOMEDIR
    |
    + go
       |
       + bin
       |
       + pkg
       |
       + src
```

## Environment variables
Add `go/bin` in your home dir to the PATH environment variable.  If you already have `go/bin` in the PATH
environment variable, then skip this section.

In Linux and MacOSX, you can use these steps to add
`go/bin` to your environment variables.

1.  Open a console.
1.  Edit `.bashrc` in your home dir.
1.  Paste this text to the end of `.bashrc`.
```
export PATH=~/go/bin:$PATH
```
1.  Save the file.
1.  Log out.
1.  Login to apply the path changes.

In Windows, you'll need to change the PATH environment variable.  Please use
Cortana or the Googles to search for "environment variables".

### Install the app
The application is stored in a Github repository. Open
a console window and type

```go get github.com/salsalabs/goengage-eoy```

When you're done, you should see a directory structure like this

```
HOMEDIR
|
+ go
   |
   + bin
   |
   + pkg
   |
   + src
      |
      + github.com
        |
        + (other directories)
        |
        + salsalabs
            |
            + goengage-eoy
```

### Resolve dependencies
Next, install the dependencies for the `goengage-eoy` Go package.

Still using the console, change the directory to
`goengage-eoy`, then type

```go get ./...```

Go will find all of the dependencies and install them.  This may take a while.
Be patient.

### Build the executable
The last step is to build the executable. Stay in the `goengage-eoy` directory.
Type this

```go build -o $HOME/go/bin/goengage-eoy cmd/main.go```

That will create a new file named `goengage-eoy` (or `goengage-eoy.exe`) in the `go/bin`
directory in your home dir.

# Operations

The  `goengage-eoy` application runs by accepting command line arguments.  You'll need to 
open a terminal window (or remain the one opened during installation).  Here's a summary
of the command-line arguments.

```
goengage-eoy --help
usage: goengage-eoy --login=LOGIN --org=ORG [<flags>]

A command-line app to create an Engage EOY

Flags:
  --help                Show context-sensitive help (also try --help-long and --help-man).
  --login=LOGIN         YAML file with API token
  --org=ORG             Organization name (for output file)
  --year=2019           Year to use for reporting, default is 2019
  --top=20              Number in top donors sheet
  --timezone="Eastern"  Choose 'Eastern' (the default), 'Central', 'Mountain', 'Pacific' or 'Alaska'
  ```

## Command-line options

|Flag|Description|
|----------|----------|
|--login|Filename of a YAML file containing the API token to use.  More info below.|
|--org|The organization's name.  Used to name the spreadsheet|
|--year|The year to report on.  Defaults to the current year.|
|--top|Number of reords to appear in the top donor list.  Defaults to 20.|
|--timezone|Timezone abbreviation for the organization.  Useful for accurately reporting donations made in the last minute of a year.|

## YAML file contents
The program requires a YAML file containing the Engage API token to use to retrieve data.  The file has this general format.
```
token: your-incredibly-long-engage-api-token
```

## Output
The output is a spreadsheet that contains a page for each of the metrics described at the top of this page.  The spreadsheet is stored in 
a file named like this.

```
[[organization]] [[year]] EOY.xlsx
```

For example, a command line like this
```bash
goengage-eoy --login=mylogin.yaml --org="Save the Ding Danged Mules" --year 2018 --timezone Central
```
will write the spreadsheet to 
```
Save the Ding Danged Mules 2018 EOY.xlsx
```

Note that using the same name and year more than once will overwrite the older version of the spreadsheet.

# Questions?  Comments?

If you have questions, then plesae use the [Issues](https://github.com/salsalabs/goengage-eoy/issues) tab.  Don't
waste your time bothering the folks at Salsalabs support.   Not worth your trouble...
