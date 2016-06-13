Installation
############


Pre-compiled binaries
=====================

Download the pre-compiled binaries for your system and architecture
from https://dl.waarp.org/dist/waarp-analyze , decompress the archive.

You can then run Waarp Analyze from the root folder of the archive
with the command :command:`./waarp-analyze`


From source
===========

You can build it Waarp Analyze from source with gb_: clone the
repository, then run ``gb build`` at the root of the project:

.. code-block:: bash

   git clone https://alm.waarp.fr/waarp-platform/waarp-analyze.git
   cd waarp-analyze
   gb build

The resulting inary is located in the :file:`bin` folder.
You can then run Waarp Analyze with the command :command:`./bin/waarp-analyze`

.. _gb: https://getgb.io