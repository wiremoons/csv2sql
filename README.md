## SUMMARYS

csv2sql - conversion program to convert a CSV file to SQL format to allow easy checking / validation, and subsequent import into a SQLite3 database using SQLite  '.read' command

## ABOUT CSV2SQL

CVS2SQL is a small simple program specifically designed to quickly convert a  coma separated value (CSV) file into structured query language (SQL) statements, that can then be used as an import source for an SQLite database.

The CSV file is also integrity checked while being converted to ensure it has a consistent number of column values throughout the file.

The first line of your CSV file will be designated as the header line - and therefore will become the column names in your subsequent SQLite database table.

Please note that any spaces or the following characters | - + @ # / \ : ( ) ' will be replaced in the column names with the underscore character (ie '_').

This is to avoid SQL syntax import issues, and make any future SQL statements referencing these column names easier to construct. You can of course rename these characters in your CSV file first. Or use the command line switch ' -k=true ' to force them to be left as is.

The rest of the CSV file will be split up on the comma character, on a per line basis. The eventual contents in your new database table will therefore be aligned to the column values - so each table row is a line from the CSV file.

The output filename (ie <sql-filename.sql>) will be created automatically for you when you run the program. Note that it will also overwrite / replace any existing file with the same name! The filename it will create will be based on your input filename, prefixed with 'SQL' and the file extension changed to '.sql'. So 'test-123.csv' -> 'SQL-test-123.sql'. 

The newly generated output file will contain the SQL statements to allow the contents of your CSV file to be imported into a new SQLite database table. The table name to be used must be provide on the command line also as ' -t tablename ' - where tablename is the name of the SQLite table to hold your CSV file data. 

To import the table and it contents, open your SQLite database with the sqlite3 program, and use:  .read <sql-filename.sql>

## OTHER INFORMATION

- Latest version is kept on GitHub here: https://github.com/wiremoons
- The program is written in Go - more information here: http://www.golang.org/
- More information on SQLite can be found here: http://www.sqlite.org/
- The program was written by Simon Rowe, licensed under "New BSD License
