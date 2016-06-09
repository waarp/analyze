Waarp Analyze
=============

Waarp Analyze is a tool that gathers information about Waarp instances
and about the server it runs on.

It is useful to generate reports about the state of the system when a
transfer is in error, or to verify the global configuration.

Reports are generated in Restructured Text and can be converted to any
format supported by docutils_ (or pandoc_, or any RST converter).

.. _docutils: http://docutils.sourceforge.net/
.. _pandoc: http://pandoc.org/


Warning about compatibility
---------------------------

For now, Waarp Analyze *only* runs on GNU/Linux 32 and 64 bits.
It *only* detects **running instances of Waarp R66 Server**.


Installation
------------

Download the pre-compiled binaries for your system and architecture
from https://dl.waarp.org/dist/waarp-analyze , decompress the archive.

You can then run Waarp Analyze with the command::

  /path/to/waarp-analyze-{version}/waarp-analyze


Execution
---------

You can run waarp-analyze directly or with following options::

  $ ./bin/waarp-analyze --help
  Usage:
    waarp-analyze [OPTIONS]

  Application Options:
    -v, --verbose  Verbose output
    -V, --version  Prints version and exits
    -o, --output=  Write report to this location. Use '- for stdout
    -H, --hostid=  Limit analyze to this Waarp instance

  Help Options:
    -h, --help     Show this help message


Roadmap
-------

Following developments are planned:

- Support for stopped instances
- Support for Waarp R66 clients
- Support for Waarp Gateway FTP
- Support for Windows


Data Gathered
-------------

Waarp Analyze includes following data in reports:

- System information:

  - Linux distribution and version
  - Kernel version
  - CPU information
  - RAM information
  - Mount points and available space on each partitions
  - System load
  - System processes

- For each found instance of supported application:

  - Its HostId
  - Its configuration
  - Its start arguments
  - The PID, user, group and environment of the process
  - the kernel limits of the process
  - The version of the JVM used
  - Its RAM, CPU threads and open FD consumption
  - The jars used, and their version
  - Its open sockets
  - Its logs (the last open file if log rotation is enabled)