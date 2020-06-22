Usage: dupe_detector \<dir\>

Simple go application to recursively and quickly find duplicate files.

A map is used to keep track of files are the same size. If there are multiple
files that are the same size, they are hashed to check for collisions.
