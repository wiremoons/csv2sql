[![NewBSD License](http://img.shields.io/badge/license-NewBSD-orange.svg?style=flat-square)](http://opensource.org/licenses/BSD-3-Clause)


## Application Summary

CVS2SQL is a small simple program specifically designed to quickly convert a comma separated value (CSV) file into simple structured query language (SQL) statements, which can then be used as an import source for an SQLite database.

## About CSV2SQL

The program was originally created to speed up the process of checking and then importing large (often greater than 1GB) CSV files into SQLite databases. The data would vary quite a bit, often being sourced from multiple financial, audit, billing and business support corporate computer systems - so there was no consistency in the CSV file formats provided, from project to project. The data was used for ad-hoc revenue assurance investigations, and often to aide recovery processes, and reporting for the associated projects. 

The different data sources (starting as CSV files - as data was extracted from business system by a different teams) would be loaded into an ad-hoc SQLite database as tables, and then analysed with the benefit of SQL, and sometimes in later stages, scripts to produce recovery data and reports on the more complex projects. The work often required a quick turn around - so any tools that could provide increased efficiency, but still maintain integrity (or even increase the integrity checking) became key. Having simple tools to improve the work flow, and produce consistent repeatable results was very important!

The csv2sql tool was created to quickly integrity check the source CSV file, report on it size (simple stats), and convert it into a text file that contained simple SQL statements. The simple SQL statements purpose was to both create a new database table to hold the CSV file data, and then to also insert the data directly into that new database table. These steps can be done by SQLite also, as it can directly import CSV files - but there was a wish to separately prepare and managed the CSV data files prior to involving the SQL database. This added the benefit of a simple additional integrity step, and put all the source CSV files data into a known state, and file format, prior to them being used with SQLite. 

So key requirements were:

	* be very fast - the source CSV files are often large, and therefore speed to process then quickly was important
	* check the CSV file contents (integrity check)  - if there are any discrepancies found they are reported in a helpful way - so the CSV file can be fixed quickly
	* the output coveted data, in SQL format, should be consistent and basic
	* the data should remain in text format to allow future access, or use with other text manipulation tools if needed
	* the SQL statement format should be as simple as possible to reduce complexity - and as SQLite treats all data as text by default - this approach was adopted. Casting using SQL can be used then if needed at a later stage - or it was handled by high-level scripting languages
	* create consistently formatted SQL table column names (ie without spaces or 'strange' characters) - to allow future reference to the columns easily when constructing new queries. Different source computer systems (and their databases) had some very varied approaches to characters and formatting used!
	* should be cross platform if possible - so it can be used on any computer system so one toll works everywhere
	* should be command line based to reduce development time - and keep it simple to use. SQLite was used via the command line anyway (either directly of via scripts) - so continuing this approach was chosen.


A few different approaches were tried over time (using tcl, Python, and c) none of which were bad at the job - however the application was ported over to Go (golang) - and it immediately benefited as Go has a great built-in CSV file handling (as well as other formats), and the speed to process the file was impressive too. It might not be as fast as c, or it might not be as simple to understand as tcl or Pyhton code wise at first, but overall it suited my requirements best! Go also supports UTF8 characters without extra work, and was cross platform too!

Key features of csv2sql include:
    
   * The CSV file is integrity checked while being converted to SQL - to 
   ensure it has a consistent number of column values. In other words the 
   number of commas in the header (ie first line) of the CSV file, are the same
   throughout the rest of the file too.

   * The first line of your CSV file will be designated as the header line -
   and therefore will become the column names in your subsequent SQLite 
   database table.

   * Any spaces or the following characters | - + @ # / \ : ( ) '
   found in the header line of you CSV file, will be replaced when they are
   used as the subsequent column names for your new SQLite table. These 
   characters will be replaced with the underscore character (ie '_'). These
   changes only apply to the header line, and are carried out to avoid SQL 
   syntax import issues, and make any future SQL statements referencing these 
   column names easier to construct. This default feature can be disabled by
   using the command line parameter ' -k=true ' if you wish.

   * You choose and specify the table name the CSV file contents will be 
   imported into in your SQLite database when you run the program.

   * The output file is a plain text file. It just contains the SQL commands 
   that are used by SQLite to create and then insert all your data into your
   specified new table. The output file can therefore be edited (if you wish) 
   to adapt it further - perhaps to suit you own needs.

## Command Line Arguments

When csv2sql is run without any parameters - it prints out the usage as follows:

```
Usage of ./csv2sql-linx64:
  -d=false:     USE: '-d=true' to include additional debug output when run
  -f="":        USE: '-f filename.csv' where filename.csv is the name and path to a CSV file that contains your data for conversion [MANDATORY]
  -h=false:     USE: '-h' to provide more detailed help on using this program
  -k=false:     USE: '-k=true' to keep original csv header fields as the SQL table column names
  -t="":        USE: '-t tablename' where tablename is the name of the SQLite table to hold your CSV file data [MANDATORY]
```

Further details of each of these command line options is below:

* DEBUG:  -d or -d=true
This enables debug output when the program is run - so it prints additional information to the screen while it is running. This additional output might be useful to better understand what the application is doing - or to pin point where in the program a problem is occurring. For normal use it is not needed - so is turn off by default (ie -d=false)

* CSV INPUT FILENAME:  -f <filename.csv>  [MANDATORY]
This command line parameter is required to allow the program to run properly, so is mandatory for successful use. It specifies the name of the input CSV file, that will be used as the source data by the program to check and convert in to SQL, ready for import into an SQLite database table. The '<filename.csv>' as shown in the example should be replaced with the name of your actuals source CSV file. If you need to include the path to the CSV file, and it contains any spaces (or special characters) - you should wrap the filename and path in quotes:  "/data-disk/datastore one/my_csv-data.csv" or "c:\Users\Fred Jones\My Documents\my csv-file.csv". There is no default value for this command line parameter - so the user my provide a CSV file to use to allow the program to run.

* ADDITIONAL HELP:  -h or -h=true
This will output additional information about the program, its purpose, and explanation of its usage. It may be useful to someone who did not originally install the program, so needs to know a bit more about it. If the program is run with this option, it will exit after displaying the help output. The default is not to show the additional help screen (ie -h=false)

* CSV HEADER CHANGES: -k=false  [default]
By default the program will change certain characters (ie space and | - + @ # / \ : ( ) ') to an underscore (ie _) when it uses the header of the CSV file, to create the new SQL database table column names. If you want to maintain you column names as they are in your source CSV file, then use this command line parameter to disable this behaviour. By default it will make the changes, so on the command lines specify -k=true to override. The -k stands for 'keep'.

* TABLE NAME: -t tablename  [MANDATORY]
This command line parameter is required to allow the program to run properly, so is mandatory for successful use. It specifies the name of the table to be created in the SQLite database when your data is imported. Change the example <tablenname> to a name of your choice, that can be used within an SQLite database. If you need to use a tablename that contains any spaces (or special characters that SQLite allows of course) - you should wrap the tablename in quotes. Examples are: mytablename  or "my table name" or my_table_name


## Compiling the Program

Assuming you already have Go install and set-up on your computer - you only need to download the single source file 'csv2sql.go'. This can then be built using the command, assuming the 'csv2sql.go' file is in you current directory:
```
go build ./csv2sql.go
```
There is also a Makefile that I use on a computer running Linux to cross compile the program for Linux (64 bit version), Windows (32 bit and 64 bit versions) and Mac OS X (64 bit version). This can be done (assuming you have your computer and Go set-up correctly) by also downloading the 'Makefile', and then entering:
```
make all
```

## Downloading Binary Version

The following binary version are available for download. Just download the file to your computer (either the .exe. or .zip version), copy the file so it is in you current path, and then run it. You may want to rename the downloaded file to just 'csv2sql' (Linux & Mac OS X) or 'csv2sql.exe' (Windows) as well.



## License

The program is licensed under the "New BSD License" or "BSD 3-Clause License". A copy of the license is available [here](https://github.com/wiremoons/csv2sql/blob/master/License.txt).

## OTHER INFORMATION

- Latest version is kept on GitHub here: [Wiremoon GitHub](https://github.com/wiremoons)
- The program is written in Go - more information here: [Go](http://www.golang.org/)
- More information on SQLite can be found here: [SQLite](http://www.sqlite.org/)
- The program was written by Simon Rowe, licensed under [New BSD License](http://opensource.org/licenses/BSD-3-Clause)
