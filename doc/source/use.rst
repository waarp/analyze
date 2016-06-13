Usage
#####

Execution
=========

.. program:: waarp-analyze

:program:`waarp-analyze` accepts the following options

.. option:: -v, --verbose

   Verbose output. All information about the collect of data and the
   generation of the report is written on the standard error.

.. option:: -V, --version

   Prints version and exits

.. option:: -o, --output=FILE

   Write report to this location. Use '-' for standard output (default
   is standard output)

.. option:: -H, --hostid=HOSTID

   Limit analyze to this Waarp instance

.. option:: -h, --help

   Show the help message

.. note::

   Depending on which system user runs :program:`waarp-analyze`, it
   might not have the rights to collect all information. It might not
   have the required UNIX permissions.

   To work around that limitation, make sure you run
   :program:`waarp-analyze` with the same user that runs Waarp
   instances.


One-Time Run
============

Waarp Analyze can be run when needed to get a full picture of the
of an instance configuration along with its system environment.

It is useful to debug the configuration or to have all the elements
necessary to fine-tune it.

Generate HTML reports
---------------------

Waarp Analyze generates its report in the `ReStructured Text`_ format.

This report can be used as is, or converted to HTML, or any format
supported by the various converter that exist (docutils_, pandoc_,
etc.).

To generate a HTML report, you can run the following command:

.. code-block:: bash

   ./waarp-analyze -o report.txt
   rst2html report.txt report.html

or use the one-liner :

.. code-block:: bash

   ./waarp-analyze | rst2html > report.html


.. _ReStructured Text: http://docutils.sourceforge.net/docs/ref/rst/restructuredtext.html
.. _docutils: http://docutils.sourceforge.net/
.. _pandoc: http://pandoc.org/


Run as error task in Waarp R66
==============================

It can be useful to run Waarp Analyze as a error task in Waarp R66 to
ease the debug of failed transfers by capturing the system environment
when the error occured.

To do so, add the following XML to your transfer rules in a
`<rerrortasks>` and/or `<serrortasks>` (whether you want to run it on
the receiver and/or the sender of the file):

.. code-block:: xml

   <tasks>
   [...]
     <task>
       <type>EXEC</type>
       <path>/path/to/waarp-analyze --output=/path/to/reports/#DATE##HOUR#-#TRANSFERID#.txt</path>
       <delay>30000</delay>
     </task>
   [...]
   </tasks>
