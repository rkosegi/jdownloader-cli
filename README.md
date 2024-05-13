# JDownloader CLI

Command line interface to [JDownloader](https://jdownloader.org/) based
on [jdownloader-go](https://github.com/rkosegi/jdownloader-go) library.

### Implemented commands

- Device
    - `jdcli device list` - list all devices associated with configured account


- Downloads
    - `jdcli download clean` - Clean completed downloads

    - Links - Manages download links
        - `jdcli download link list` - list links
        - `jdcli download link list` - remove link(s) from downloader

    - Packages - Manages download packages
        - `jdcli download link list` - list links


- Link collector
    - `jdcli links list` - list links in link collector
    - `jdcli links add` - add links into link collector


- Login
    - `jdcli login` - configure account
    - `jdcli logout` - discard any configured credentials

- Miscellaneous
    - `jdcli version` - display current program version
