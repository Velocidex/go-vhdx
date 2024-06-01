# VHDX Parser

Based on documentation from https://learn.microsoft.com/en-us/openspecs/windows_protocols/ms-vhdx/83e061f8-f6e2-4de1-91bd-5d518a43d477

This parser is used by Velociraptor's vhdx accessor. Example of use in
VQL:

```sql
SELECT OSPath.Path AS OSPath, Size, Mode.String
FROM glob(
  globs="*", accessor="raw_ntfs", root=pathspec(
     Path="/",
     DelegateAccessor="offset",
     DelegatePath=pathspec(
         Path="/65536",
         DelegateAccessor="vhdx",
         DelegatePath="/tmp/test.vhdx")))
```

# Testing locally

There is a small tool that allows inspection of the VHDX volume:

The following will print some information about internal data structures and metadata.
```
govhdx parse test.vhdx
```

The following will dump the image into stdout and redirect to a flat dd file.
```
govhdx cat test.vhdx > test.dd
```
