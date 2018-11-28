Some old go utility code that I created when I started to learn go.

BY NOW THIS REPOSITORY IS NO LONGER ACTIVELY MAINTAINED AND SHOULD NOT BE USED ANYMORE.

Please consider copying any code from this repository which is relevant to you
into your own projects (i.e. logic for string sets).

Here some advice from [the go blog][1]:

> Avoid meaningless package names. Packages named util, common, or misc provide
  clients with no sense of what the package contains. This makes it harder for
  clients to use the package and makes it harder for maintainers to keep the
  package focused. Over time, they accumulate dependencies that can make
  compilation significantly and unnecessarily slower, especially in large
  programs. And since such package names are generic, they are more likely to
  collide with other packages imported by client code, forcing clients to invent
  names to distinguish them.
  
> Break up generic packages. To fix such packages, look for types and functions
  with common name elements and pull them into their own package.
  
[1]: https://blog.golang.org/package-names
