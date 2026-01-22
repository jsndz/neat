# Neat is a clone of git

## Stage 1 : neat init

Git is Version control. It tracks changes overtime.
Everything about git is stored in .git which is a hidden file.
Create a .git folder with refs and objects as sub-folders and a HEAD file.
If you know linked list we name the main point as head.
This HEAD is similar it should point to something for now lets say HEAD points to "ref: refs/heads/main".

## Stage 2: neat add

You should know how git add works.
add puts files into staging area consider it as draft for the project.
Here you need to:

- read the file content from the file
- make a blob from it with the structure 
 blob [size]\0[raw content] -> which is stored as binary
- and from that blob you need to get hash will use sha1 algorithm
- with hash as filename and blob as file content we will store in .neat/objects
- while storing the file, first 2 letter as folder and others as file
    0b8e848f7890715dba5f4703a9160e1ef4debb07 -> .neat/objects/0b/8e848f7890715dba5f4703a9160e1ef4debb07
- then create an index for the file

## Stage 3: neat cat-file 

- cat-file is used for reading the content based on hash
- it can have -p as flag followed by hash
- so find the file based on hash and decompress it to read the file
- if -p then remove the things before null("\0") and return it