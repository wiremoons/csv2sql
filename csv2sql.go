/*

   csv2sql - conversion program to convert a csv file to sql format
   		to allow easy checking / validation, and import into a SQLite3
   		database using SQLite  '.read' command

	author: simon rowe <simon@wiremoons.com>
	license: open-source released under "New BSD License"

   version: 0.5
   created: 16th April 2014 - initial outline code written
   updated: 17th April 2014 - add flags and output file handling
   updated: 27th April 2014 - wrap in double quotes instead of single
   updated: 28th April 2014 - add flush io file buffer to fix SQL missing EOF
   updated: 19th July 2014 - add more help text, tidy up comments and code
   updated: 6th August 2014 - enabled the -k flag to alter the table header characters

*/
package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// set global variables

// set the version of the app here
var appversion string = "0.5"

// below used by flag for command line args
var tableName string
var csvFileName string
var keepOrigCols bool
var debugSwitch bool

// init() function - always runs before main() - used here to set-up required flags variables
// from the command line parameters provided by the user when they run the app
func init() {
	// IntVar; StringVar; BoolVar all required: variable, cmd line flag, initial value, description used by flag.Usage() on error / help
	flag.StringVar(&tableName, "t", "", "\tUSE: '-t tablename' where tablename is the name of the SQLite table to hold your CSV file data [MANDATORY]")
	flag.StringVar(&csvFileName, "f", "", "\tUSE: '-f filename.csv' where filename.csv is the name and path to a CSV file that contains your data for conversion [MANDATORY]")
	flag.BoolVar(&keepOrigCols, "k", false, "\tUSE: '-k=true' to keep original csv header fields as the SQL table column names")
	flag.BoolVar(&debugSwitch, "d", false, "\tUSE: '-d=true' to include additional debug output when run")
}

//
//  FUNCTION: create a filename string for the SQL data to be written too - return it
//
func SQLFileName() (filename string) {
	// include the name of the csv file from command line (ie csvFileName)
	// remove any path etc
	var justFileName = filepath.Base(csvFileName)
	var extension = filepath.Ext(csvFileName)
	// remove the file extn
	justFileName = justFileName[0 : len(justFileName)-len(extension)]
	// get a date and time stamp - use GoLang reference date of: Mon Jan 2 15:04:05 MST 2006
	// TODO: figure out how to make this work - so filename has timestamp too ??
	//fileDate, err := time.Parse("2006-01-02", time.Now().String())
	//if err != nil {
	//	panic(err)
	//}
	//fileDate := fileDate.String()
	//fmt.Printf("\n%s\n", fileDate)
	sqlOutFile := "SQL-" + justFileName + ".sql"
	return sqlOutFile
}

//
//  FUNCTION: display a banner and help information on the screen
//  information is displayed when the program is run without including
//  any command line parameters - so assumes you want help to run it
//
func printBanner() {
	// add the help and about text to the variable 'about'
	about := `
    ABOUT CSV2SQL
    ¯¯¯¯¯¯¯¯¯¯¯¯¯¯
    CVS2SQL is a small simple program specifically designed to quickly convert a 
    coma separated value (CSV) file into structured query language (SQL) 
    statements, that can then be used as an import source for an SQLite database.

    The CSV file is also integrity checked while being converted to ensure it
    has a consistent number of column values throughout the file.

    The first line of your CSV file will be designated as the header line - and 
    therefore will become the column names in your subsequent SQLite database 
    table.

    Please note that any spaces or the following characters | - + @ # / \ : ( ) '
    will be replaced in the column names with the underscore character (ie '_').
    
    This is to avoid SQL syntax import issues, and make any future SQL statements
    referencing these column names easier to construct. You can of course rename
    these characters in your CSV file first. Or use the command line switch
    ' -k=true ' to force them to be left as is.

    The rest of the CSV file will be split up on the comma character, on a per 
    line basis. The eventual contents in your new database table will therefore
    be aligned to the column values - so each table row is a line from the CSV
    file.

    The output filename (ie <sql-filename.sql>) will be created 
    automatically for you when you run the program. Note that it will also 
    overwrite / replace any existing file with the same name! The filename it
    will create will be based on your input filename, prefixed with 'SQL' and 
    the file extension changed to '.sql'. So 'test-123.csv' -> 'SQL-test-123.sql'. 

    The newly generated output file will contain the SQL statements to allow
    the contents of your CSV file to be imported into a new SQLite database 
    table. The table name to be used must be provide on the command line also
    as ' -t tablename ' - where tablename is the name of the SQLite table to 
    hold your CSV file data. 

    To import the table and it contents, open your SQLite database with the 
    sqlite3 program, and use:  .read <sql-filename.sql>
 
    Latest version is kept on GitHub here: https://github.com/wiremoons
    The program is written in Go - more information here: http://www.golang.org/
    More information on SQLite can be found here: http://www.sqlite.org/
    The program was written by Simon Rowe, licensed under "New BSD License"
	`
	// now display the information on screen
	fmt.Println("\n\t\t\tcsv2sql conversion program\n\t\t\t\tVersion:", appversion, "\n", about)
}

//////////////////////////////////////////////////////////////////////////////
//
//                     MAIN STARTS HERE
//
//////////////////////////////////////////////////////////////////////////////

func main() {
	//-------------------------------------------------------------------------
	// sort out the command line arguments
	//-------------------------------------------------------------------------
	// get the command line args passed to the program
	flag.Parse()
	// if debug is enabled - confirm the command line parameters received
	if debugSwitch {
		fmt.Println("Command Line Arguments provided are:")
		fmt.Println("\tCSV file to use:", csvFileName)
		fmt.Println("\tSQL table name to use:", tableName)
		fmt.Println("\tKeep original csv header fields:", strconv.FormatBool(keepOrigCols))
		fmt.Println("\tDisplay additional debug output when run:", strconv.FormatBool(debugSwitch))
	}

	// check we have a table name and csv file to work with - otherwise abort
	if csvFileName == "" || tableName == "" {
		printBanner()
		fmt.Println("ERROR: please provide both a 'table name' and the input 'CSV filename' to use\nrun 'csv2sql --help' for more information")
		//fmt.Println("Usage:",flag.Usage,"Command Line:",flag.CommandLine)
		os.Exit(-2)
	}

	//-------------------------------------------------------------------------
	// open and prepare the CSV input file
	//-------------------------------------------------------------------------
	// TODO:  manage multiple input files (ie csv2sql -f * -t testtable[1..x]) ??
	if debugSwitch {
		fmt.Println("Opening the CSV file:", csvFileName)
	}
	// open the CSV file - name provided via command line input - handle 'file'
	file, err := os.Open(csvFileName)
	// error - if we have one exit as CSV file not right
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		os.Exit(-3)
	}
	// now file is open - defer the close of CSV file handle until we return
	defer file.Close()
	// connect a CSV reader to the file handle - which is the actual opened
	// CSV file
	// TODO : is there an error from this to check?
	reader := csv.NewReader(file)

	//-------------------------------------------------------------------------
	// open and prepare the SQL output file
	//-------------------------------------------------------------------------
	// get a new filename to write the SQl converted data into - call our
	// function SQLFileName() to obtain a suitable string for the new filename
	// TODO : ad option to output to stdout instead of a file only
	sqlOutFile := SQLFileName()
	if debugSwitch {
		fmt.Println("Opening the SQL output file:", sqlOutFile)
	}

	// open the new file using the name we obtained above - handle 'filesql'
	filesql, err := os.Create(sqlOutFile)
	// error - if we have one when trying open & create the new file
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	// now new file is open - defer the close of the file handle until we return
	defer filesql.Close()
	// attach the opened new sql file handle to a buffered file writer
	// the buffered file writer has the handle 'sqlFileBuffer'
	sqlFileBuffer := bufio.NewWriter(filesql)

	//-------------------------------------------------------------------------
	// prepare to read the each line of the CSV file - and write out to the SQl
	//-------------------------------------------------------------------------
	// track the number of lines in the csv file
	lineCount := 0
	// track number of fields in csv file
	csvFields := 0

	// grab time now - so can calculate how long it takes to process the file
	start := time.Now()

	// create a buffer to hold each line of the SQL file as we build it
	// handle to this buffer is called 'strbuffer'
	var strbuffer bytes.Buffer

	// START - processing of each line in the CSV input file
	//-------------------------------------------------------------------------
	// loop through the csv file until EOF - or until we hit an error in parsing it.
	// Data is read in for each line of the csv file and held in the variable
	// 'record'.  Build a string for each line - wrapped with the SQL and
	// then output to the SQL file writer in its completed new form
	//-------------------------------------------------------------------------
	for {
		record, err := reader.Read()

		// if we hit end of file (EOF) or another unexpected error
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error:", err)
			return
		}

		// get the number of fields in the CSV file on this line
		csvFields = len(record)

		// if we are processing the first line - use the record field contents
		// as the SQL table column names - add to the temp string 'strbuffer'
		// use the tablename provided by the user
		// TODO - add option to skip this line if user is adding data to an existing table?
		if lineCount == 0 {
			strbuffer.WriteString("PRAGMA foreign_keys=OFF;\nBEGIN TRANSACTION;\nCREATE TABLE " + tableName + " (")
		}

		// if any line except the first one :
		// print the start of the SQL insert statement for the record
		// and  - add to the temp string 'strbuffer'
		// use the tablename provided by the user
		if lineCount > 0 {
			strbuffer.WriteString("INSERT INTO " + tableName + " VALUES (")
		}
		// loop through each of the csv lines individual fields held in 'record'
		// len(record) tells us how many fields are on this line - so we loop right number of times
		for i := 0; i < len(record); i++ {
			// if we are processing the first line used for the table column name - update the
			// record field contents to remove the characters: space | - + @ # / \ : ( ) '
			// from the SQL table column names. Can be overridden on command line with '-k true'
			if (lineCount == 0) && (keepOrigCols == false) {
				// for debug - output info so we can see current field being processed
				if debugSwitch {
					fmt.Printf("Running header clean up for '%s' ", record[i])
				}
				// call the function cleanHeader to do clean up on this field
				record[i] = cleanHeader(record[i])
				// for debug - output info so we can see any changes now made
				if debugSwitch {
					fmt.Printf("changed to '%s'\n", record[i])
				}
			}
			// if a csv record field is empty or has the text "NULL" - replace it with actual NULL field in SQLite
			// otherwise just wrap the existing content with ''
			// TODO : make sure we don't try to create a 'NULL' table column name?
			if len(record[i]) == 0 || record[i] == "NULL" {
				strbuffer.WriteString("NULL")
			} else {
				strbuffer.WriteString("\"" + record[i] + "\"")
			}
			// if we have not reached the last record yet - add a coma also to the output
			if i < len(record)-1 {
				strbuffer.WriteString(",")
			}
		}
		// end of the line - so output SQL format required ');' and newline
		strbuffer.WriteString(");\n")
		// line of SQL is complete - so push out to the new SQL file
		bWritten, err := sqlFileBuffer.WriteString(strbuffer.String())
		// check it wrote data ok - otherwise report the error giving the line number affected
		if (err != nil) || (bWritten != len(strbuffer.Bytes())) {
			fmt.Printf("WARNING: Error writing to SQL file line %d: %s", lineCount, err)
			return
		}
		// reset the string buffer - so it is empty ready for the next line to build
		strbuffer.Reset()
		// for debug - show the line number we are processing from the CSV file
		if debugSwitch {
			fmt.Print("..", lineCount)
		}
		// increment the line count - and loop back around for next line of the CSV file
		lineCount += 1
	}
	// END - reached the end of processing each line of the input CSV file
	//
	if debugSwitch {
		fmt.Println("\ncsv file processing complete - outputted to the new SQL file: ", sqlOutFile)
	}
	// finished processing the csv input file lines - so close off the SQL statements
	strbuffer.WriteString("COMMIT;\n")
	// write out final line to the SQL file
	bWritten, err := sqlFileBuffer.WriteString(strbuffer.String())
	// check it wrote data ok - otherwise report the error giving the line number affected
	if (err != nil) || (bWritten != len(strbuffer.Bytes())) {
		fmt.Printf("WARNING: Error outputting final line of the SQL file: line %d: %s", lineCount, err)
		return
	}
	if debugSwitch {
		fmt.Println("SQL file write complete")
	}
	fmt.Println("\nDONE\n\tCSV file processing complete, and the new SQL file format was written to: ", sqlOutFile)
	// finished the SQl file data writing - flush any IO buffers
	// NB below flush required as the data was being lost otherwise - maybe a bug in go version 1.2 only?
	sqlFileBuffer.Flush()
	// reset the string buffer - so it is empty as it is no longer needed
	strbuffer.Reset()
	// stop the timer for the SQL file creation process
	end := time.Now()

	// print out some stats about the csv file processed
	fmt.Println("\nSTATS\n\tCSV file", csvFileName, "has", lineCount, "lines with", csvFields, "CSV fields per record")
	fmt.Println("\tThe conversion took", end.Sub(start), "to run.\n\nAll is well.\n")
}

//
//  cleanHeader receives a string and removes the characters: space | - + @ # / \ : ( ) '
//  Function is used to clean up the CSV file header fields as they will be used for column table names
//  in our SQLIte database. Therefore we don't want any odd characters for our table column names
//
//  TODO:  consider using: strings.NewReplacer function instead?
//
func cleanHeader(headField string) string {
	// ok - remove any spaces and replace with _
	headField = strings.Replace(headField, " ", "_", -1)
	// ok - remove any | and replace with _
	headField = strings.Replace(headField, "|", "_", -1)
	// ok - remove any - and replace with _
	headField = strings.Replace(headField, "-", "_", -1)
	// ok - remove any + and replace with _
	headField = strings.Replace(headField, "+", "_", -1)
	// ok - remove any @ and replace with _
	headField = strings.Replace(headField, "@", "_", -1)
	// ok - remove any # and replace with _
	headField = strings.Replace(headField, "#", "_", -1)
	// ok - remove any / and replace with _
	headField = strings.Replace(headField, "/", "_", -1)
	// ok - remove any \ and replace with _
	headField = strings.Replace(headField, "\\", "_", -1)
	// ok - remove any : and replace with _
	headField = strings.Replace(headField, ":", "_", -1)
	// ok - remove any ( and replace with _
	headField = strings.Replace(headField, "(", "_", -1)
	// ok - remove any ) and replace with _
	headField = strings.Replace(headField, ")", "_", -1)
	// ok - remove any ' and replace with _
	headField = strings.Replace(headField, "'", "_", -1)
	return headField
}
