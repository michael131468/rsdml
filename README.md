# RSDML (Recursively Set Directory MTime to Latest)

A small tool to update the modification time ([mtime][1]) of directories in a
directory tree to match the newest mtime attributes of the files within. 

## Why?

The mtime of a directory is [only updated when a file directly within it is
added, removed, or renamed][2]. This means files in deeper subdirectories do not
affect the mtime of a top level parent directory. It also means that changes
to the files in the directories do not update the parent directories mtime
value.

The effect of this is you can have a directory with an older last modified
timestamp than its children directories or the files within.

This can be a problem when you want to use the last modified time of a
directory as an indicator for users about changes within them.

For example, if you allow users to browse a directory tree with
[Apache AutoIndex][3] and display the last modified times, a user starting from
the index may not be able to find recently modified files.

A user looking at the initial index page can be confused why the top level links
show older timestamps compared to the contents themselves. The tree below shows
how the timestamps might be set by default. If a user is browsing from the top
directory they may not be able to find the most recently changed file,
receipt2.txt.

```
$ tree -D
[Dec 11 08:50]  .
└── [Dec 11 08:50]  example
    ├── [Dec 11 08:58]  index.html
    ├── [Dec 11 08:59]  invoices
    │   ├── [Dec 11 09:00]  invoice1.txt
    │   └── [Dec 11 09:00]  invoice2.txt
    └── [Dec 11 08:59]  receipts
        ├── [Dec 11 09:00]  receipt1.txt
        └── [Dec 12 09:00]  receipt2.txt
```

This tool (rsdml) can be run over a set of directory trees to make the mtime
values of the directories within be synchronised to the newest child values. If
rsdml was run on the example above the tree would be updated to look like the
following. The result being the date/time "Dec 12 09:00" is propagated to the
parent directories: receipts, example and the root directory.

```
$ tree -D
[Dec 12 09:00]  .
└── [Dec 12 09:00]  example
    ├── [Dec 11 08:58]  index.html
    ├── [Dec 11 09:00]  invoices
    │   ├── [Dec 11 09:00]  newest.txt
    │   └── [Dec 11 09:00]  oldest.txt
    └── [Dec 12 09:00]  receipts
        ├── [Dec 11 09:00]  newest.txt
        └── [Dec 12 09:00]  oldest.txt
```
## Why Not a Shell Script?

There are many alternatives to do this like using [shell scripts][4]. In my opinion,
it can be fragile to do this via shell scripts as it is dependent on the host system
utils (busybox vs util-linux, bash vs dash, etc). Having a statically compiled tool
makes it highly portable and less fragile.

## Isn't This Race Condition Prone?

Yes. If files in the directory tree are being updated/changed while the program runs,
it doesn't handle that. The program works in two stages. First by determining the set
of directories in the directory tree from bottom to top. Second by iterating through
that determined set of directories (bottom to top) and updating the mtime values with
consideration to the immediate children of the directory at that point in time.

New directories created after the first stage are not accounted for. Files and
directories changed after step 2 occurs for a directory is run are not accounted for.

Deleted directories after step 1 are handled aas gracefully as possible.

In my opinion, the result might be not perfect but it does a best effort attempt that is
reasonable.

## Usage

Coming soon.

## License

See [LICENSE](LICENSE) file.

[1]: https://www.gnu.org/software/coreutils/manual/html_node/File-timestamps.html
[2]: https://stackoverflow.com/questions/3620684/directory-last-modified-date
[3]: https://httpd.apache.org/docs/2.4/mod/mod_autoindex.html
[4]: https://unix.stackexchange.com/questions/1524/how-do-i-change-folder-timestamps-recursively-to-the-newest-file
