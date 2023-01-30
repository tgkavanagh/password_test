package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
)

/*******************************************************************************
 * Constants
 ******************************************************************************/

// Using "\r\n" as this is used in the included sample results file
const ValidResultString string = "<%s> is acceptable.\r\n"

const InvalidResultString string = "<%s> is not acceptable.\r\n"

const EndOfInputDelimiter string = "end"

const PasswordMinLength int = 1

const PasswordMaxLength int = 20

const ConsecutiveCharacterLimit uint8 = 3

const AllowedDuplicateList string = "eo"

/*******************************************************************************
 * Regex Expressions
 ******************************************************************************/
var ApprovedCharacterSetRegexPtr *regexp.Regexp = regexp.MustCompile(`(^[a-z]+$)`)

var VowelsRegexPtr *regexp.Regexp = regexp.MustCompile(`[aeiou]{1,}`)

var ConsecutiveVowelCharactersRegexPtr *regexp.Regexp = regexp.MustCompile(fmt.Sprintf("[aeiou]{%d}", ConsecutiveCharacterLimit))

var ConsecutiveConsonantCharactersRegexPtr *regexp.Regexp = regexp.MustCompile(fmt.Sprintf("[b-df-hj-np-tv-z]{%d}", ConsecutiveCharacterLimit))

/*******************************************************************************
 *
 * Name: printResult
 *
 * Description: Print whether or not the password is acceptable
 *
 * Input:
 *         password - string password
 *				 result - boolean indicating if the password is acceptable
 *
 * Return Value: N/A
 *
 ******************************************************************************/
func printResult(outfp *os.File, password string, result bool) {
	rStr := ""

	if result == true {
		rStr = fmt.Sprintf(ValidResultString, password)
	} else {
		rStr = fmt.Sprintf(InvalidResultString, password)
	}

	_, err := outfp.WriteString(rStr)
	if err != nil {
		fmt.Printf("Failed to write results for (%s) - %v\n", password, err)
	}
}

/*******************************************************************************
 *
 * Name: passwordLengthTest
 *
 * Description: Validate the pasword length is within spec
 *
 * Input:
 *         password - string password to check
 *
 * Return Value:
 *         true - password length is valid
 *         false  - password length is invalid
 *
 ******************************************************************************/
func passwordLengthTest(password string) bool {
	pwLen := len(password)
	return (pwLen >= PasswordMinLength) && (pwLen <= PasswordMaxLength)
}

/*******************************************************************************
 *
 * Name: passwordApprovedCharactersTest
 *
 * Description: Check whether the supplied password contains only approved
 *              characters - no numbers, no special characters, and no uppercase
 *              characters
 *
 * Input:
 *         password - string password to check
 *
 * Return Value:
 *         true - password contains Approved characters
 *         false  - password contains Unapproved characters
 *
 ******************************************************************************/
func passwordApprovedCharactersTest(password string) bool {
	return ApprovedCharacterSetRegexPtr.MatchString(password)
}

/*******************************************************************************
 *
 * Name: passwordDuplicateLettersTest
 *
 * Description: Check whether the supplied password contains two consecutive
 *              occurrences of the same letter - exceptions are ee and oo
 *
 * Input:
 *         password - string password to check
 *
 * Return Value:
 *         true - password does contain consecutive occurrences of the same letter
 *         false  - password does not contain consecutive occurrences of the same letter
 *
 ******************************************************************************/
func passwordDuplicateLettersTest(password string) bool {
	prevChar := ""
	curChar := ""

	// Cycle through each character to the password and see if we have consecutive
	// characters other than those is the allowed list
	for _, c := range password {
		curChar = string(c)
		if prevChar == curChar {
			if strings.Contains(AllowedDuplicateList, curChar) == false {
				return true
			}
		}

		prevChar = curChar
	}

	return false
}

/*******************************************************************************
 *
 * Name: passwordContainsVowelTest
 *
 * Description: Check whether the supplied password contains at least 1 vowel
 *
 * Input:
 *         password - string password to check
 *
 * Return Value:
 *         true - password contains at least one vowel
 *         false  - password does not contain a vowel
 *
 ******************************************************************************/
func passwordContainsVowelTest(password string) bool {
	return VowelsRegexPtr.MatchString(password)
}

/*******************************************************************************
 *
 * Name: passwordContainsConsecutiveCharacterTest
 *
 * Description: Check whether the supplied password contains consecutive vowels
 *              or consonants
 *
 * Input:
 *         password - string password to check
 *         count - number of consecutive characters
 *
 * Return Value:
 *         true - password contains at 3 consecutive vowels or consonants
 *         false  - password does not contains at 3 consecutive vowels or consonants
 *
 ******************************************************************************/
func passwordContainsConsecutiveCharacterTest(password string, count uint8) bool {
	return ConsecutiveVowelCharactersRegexPtr.MatchString(password) ||
				 ConsecutiveConsonantCharactersRegexPtr.MatchString(password)
}

/*******************************************************************************
 *
 * Name: validatePasswordBrief
 *
 * Description: Check if the specified password meets all of our requirement
 *
 * Note: This version of the validate API will return on the first failure to
 *       minimize execution time
 *
 * Input:
 *         password - string password to check
 *
 * Return Value:
 *         true - password meets requirements
 *         false  - password does not meet requirements
 *
 ******************************************************************************/
func validatePasswordBrief(password string) bool {
	pwAcceptable := false

	// Check the password length
	if passwordLengthTest(password) == true {
		// Check if the password contains only approved characters
		if passwordApprovedCharactersTest(password) == true {
			// Check if the password contains at least one vowel
			if passwordContainsVowelTest(password) == true {
				// Check if the password contains consecutive occurrences of the same letter
				if passwordDuplicateLettersTest(password) == false {
					// Check if the password contains consecutive vowles of consonants
					if passwordContainsConsecutiveCharacterTest(password, ConsecutiveCharacterLimit) == false {
						pwAcceptable = true
					}
				}
			}
		}
	}

	return pwAcceptable
}

/*******************************************************************************
 *
 * Name: validatePasswordFull
 *
 * Description: Check if the specified password meets all of our requirement
 *
 * Note: This version of the validate API run and interpret all checks before
 *       returning.  To be used as a debugging call
 *
 * Input:
 *         password - string password to check
 *
 * Return Value:
 *         true - password meets requirements
 *         false  - password does not meet requirements
 *
 ******************************************************************************/
func validatePasswordFull(outfp *os.File, password string) bool {
	pwAcceptable := false

	pwlt := passwordLengthTest(password)

	pwact := passwordApprovedCharactersTest(password)

	pwcvt := passwordContainsVowelTest(password)

	pwdlt := passwordDuplicateLettersTest(password)

	pwccct := passwordContainsConsecutiveCharacterTest(password, ConsecutiveCharacterLimit)

	pwAcceptable = pwlt && pwact && pwcvt && !pwdlt && !pwccct

	rStr:= ""
	rStr = fmt.Sprintf("Password (%s) pwlt (%v) pwact (%v) pwcvt (%v) pwdlt (%v) pwcct (%v)\n\n\n",
		password, pwlt, pwact, pwcvt, pwdlt, pwccct)


	_, err := outfp.WriteString(rStr)
	if err != nil {
		fmt.Printf("Failed to write results for (%s) - %v\n", password, err)
	}

	return pwAcceptable
}

/*******************************************************************************
 *
 * Name: main
 *
 * Description: Main control function
 *
 * Input: Command-line arguements
 *         inFile - file containing the passwords to check
 *         outFile - file to store/log results
 *
 * Return Value:
 *         N/A
 *
 ******************************************************************************/
func main() {
	var ifp *os.File
	var ofp *os.File

	//  Let's track whetehr or not the end of the input keyword "end" was found
	// as is mentioned in the requirements
	endOfInputFound := false

	// The user needs to specify the inFile and outFile as command-line arguments.
	// Total number of command-line arguments must be 2
	inputFile := flag.String("inFile", "", "File containing passwords to test")
	outputFile := flag.String("outFile", "", "File to store log results")

	flag.Parse()

	// If the input or output file name is missing, error and exit
	if (*inputFile == "") || (*outputFile == "") {
		log.Fatal("Missing command-line arguments\n")
	}

	// Check that the inputFile exists.  We should sanitize the filenames but for
	// exercise, we will skip it.
	_, err := os.Stat(*inputFile)
	if err != nil {
		log.Fatal(err)
	}

	// Open the input file and create the output file
	ifp, err = os.Open(*inputFile)
	if err != nil {
		log.Fatal(err)
	}

	ofp, err = os.Create(*outputFile)
	if err != nil {
		ifp.Close()
		log.Fatal(err)
	}

  // Create a scanner to read in the passwords from the input file.
	fscanner := bufio.NewScanner(ifp)
	fscanner.Split(bufio.ScanLines)

	// Get the start time to measure execution time
	start := time.Now()

	// Loop through all lines in the file and validate the password
	passwd := ""

	for fscanner.Scan() {
		passwd = fscanner.Text()

		// If the password is "end" then this marks the end of the input and we need
		// to break out of this loop
		if passwd == EndOfInputDelimiter {
			endOfInputFound = true
			break
		}

		// Validate the current password and save the result
		printResult(ofp, passwd, validatePasswordBrief(passwd))
	}

	// Calculate and output the execution time
	elapsed := time.Since(start)
	fmt.Printf("Execution time: %s\n", elapsed)

	// Log if we did not find the end of input keyword
  if endOfInputFound == false {
		fmt.Printf("End of input keyword (%s) not found before EOF\n", EndOfInputDelimiter)
	}

	// Close all open files
  ifp.Close()
	ofp.Close()
}
