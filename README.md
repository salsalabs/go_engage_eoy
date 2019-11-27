# go_engage_eoy

Go version of the Engage EOY project.  Provides a spreadsheet that details end-of-year data for an Engage client.

# Background

Salsa's Engage product does not yet have extensive reporting.  Because of that, Engage clients are missing out on
the typical end-of-year reporting that occurs after every Giving Tuesday.

# Summary

This application accepts an Engage API Token, then provides the usual end-of-year reports that Salsa's client want.

## How did we do this year?

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

## How did we do this year by month?

-   Month
-   Donations
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

## Fundraising, Year-over-year

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

## Fundraising, month-over-month

-   Month
-   Year
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

## Show me who donated, how much, how given.

-   First Name
-   Last Name
-   Email
-   Number of donations
-   Total donations
-   Number of one-time donations.
-   Total one-time donations
-   Number of recurring donations.
-   Total recurring donations
-   Number of offline donations.
-   Total offline donations
-   Largest donation and date
-   Smallest donation and date
-   Average donation

## Who were our biggest donors?

-   Supporter
    -   First Name
    -   Last Name
    -   Email
-   Number of donations
-   Total donations
-   Number of one-time donations.
-   Total one-time donations
-   Number of recurring donations.
-   Total recurring donations
-   Number of offline donations.
-   Total offline donations
-   Largest donation and date
-   Smallest donation and date
-   Average donation

## Show me the revenue by activity page.

-   Activity page name
-   Activity page start date
-   Activity page end date
-   Number of tickets
-   Total ticket revenue (outside of
-   Number of purchases
-   Total purchase revenue
-   Number of donations
-   Total of donations

## Show me the projected revenue from recurring donations

-   First Name
-   Last Name
-   Email
-   Recurring donation
-   Start date
-   End date (if recorded)
-   Credit card expiration date.
-   Recurring amount.
-   Number of months until the lesser of end date or credit card expiration.
-   Total projected revenue (recurring amount x number of months)

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
Cortana or the Googles to search for "Environment variables".

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

# Questions?  Comments?
