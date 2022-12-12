# RSDML (Recursively Set Directory MTime to Latest)

A small statically compiled tool to recursively update a tree of directories to
have an mtime that matches the most recently modified mtime of the directory
contents.

The mtime of a directory is only updated when a file directly within it is
added, removed, or renamed. This means files in deeper subdirectories do not
affect the mtime of a top level parent directory. It also means that changes to
the files themselves do not update the parent directories mtime value.

This poses a problem sometimes when you want to reflect updates to files in a
deep directory structure through the mtime of the top level directories. An
example is where you have a webserver hosting a file serve from a directory
structure with just the names and last modified timestamp visible. A user
looking at the initial index page can be confused why the top level links show
older timestamps compared to the contents themselves.

```
$ tree -D
[Dec 12 08:58]  .
└── [Dec 12 08:58]  example
    ├── [Dec 12 08:58]  index.html
    ├── [Dec 12 08:59]  invoices
    │   ├── [Dec 12 09:00]  newest.txt
    │   └── [Dec 12 09:00]  oldest.txt
    └── [Dec 12 08:59]  receipts
        ├── [Dec 12 09:00]  newest.txt
        └── [Dec 12 09:00]  oldest.txt
```

Example: See how the parent directories have older timestamps than the contents
of said directories.

This tool (rsdml) can be run over a set of directory trees to make their mtime
values synchronised to the newest child values. If rsdml was run on the example
above the tree would be updated to look like the following.

```
$ tree -D
[Dec 12 09:00]  .
└── [Dec 12 09:00]  example
    ├── [Dec 12 08:58]  index.html
    ├── [Dec 12 09:00]  invoices
    │   ├── [Dec 12 09:00]  newest.txt
    │   └── [Dec 12 09:00]  oldest.txt
    └── [Dec 12 09:00]  receipts
        ├── [Dec 12 09:00]  newest.txt
        └── [Dec 12 09:00]  oldest.txt
```
