![carto](./_img/carto.png)

Carto can help you make typed maps quickly and easily.  

Out of the box, carto will generate a goroutine-safe struct with the specified key and value pairings

```
> carto --help
```

```
C A R T O

🌐 Maps made easy

usage:
carto -p <package> -s <structname> -k <keytype> -v <valuetype> [options]

-p    (string)      Package name (required)
-s    (string)      Struct name  (required)
-k    (string)      Key type     (required)
-v    (string)      Value type   (required)
-r    (string)      Receiver name (defaults to lowercase first char of
                    struct name)
-o    (string)      Output file path (if omitted, prints to STDOUT)
-i    (string)      Variable name for internal map (defaults to internal)
-rv   (bool)        Receivers are by value
-b    (bool)        "Get" return signature includes a bool value indicating
                    if the key exists in the internal map
-d    (bool)        "Get" signature has second parameter for default return
                    value when key does not exist in the internal map
-lz   (bool)        Will lazy-instantiate the internal map when a write
                    operation is used
-version            Print version and exit
```

Developing Carto

