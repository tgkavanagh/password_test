package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

var passwordTestSet = map[string]bool {
	"a":true,
	"tv":false,
	"ptoui":false,
	"bontres":false,
	"zoggax":false,
	"wiinq":false,
	"eep":true,
	"houctuh":true,
	"ei":true,
	"cd":false,
	"bcdfghjklmnpqrstvwxy":false,
	"ee":true,
	"oo":true,
	"jj":false,
	"aei":false,
	"vwx":false,
	"to":true,
	"in":true,
	"try":false,
	"ask":true,
	"bot":true,
	"abcodefgohijkolmonop":true,
	"jkloxoxoxoxoxoxoxoxo":false,
	"mamamamajklamamamama":false,
	"egegegegegegegegejkl":false,
	"aeiququququququququq":false,
	"orororaeirororororor":false,
	"zezezezezezezezezoua":false,
	"eeprop":true,
	"peerop":true,
	"propee":true,
	"ooplant":true,
	"poolant":true,
	"plantoo":true,
	"ssuper":false,
	"supper":false,
	"superr":false,
	"aarmada":false,
	"armaada":false,
	"armadaa":false,
	"banana":true,
	"rhythm":false,
	"breakneck":true,
}

/*
func main() {
	pwAcceptable := false

	// Loop through all passwords in the test set
	for passwd, valid := range passwordTestSet {
		pwlt := passwordLengthTest(passwd)

		pwact := passwordApprovedCharactersTest(passwd)

		pwcvt := passwordContainsVowelTest(passwd)

		pwdlt := passwordDuplicateLettersTest(passwd)

		pwccct := passwordContainsConsecutiveCharacterTest(passwd, ConsecutiveCharacterLimit)

		pwAcceptable = pwlt && pwact && pwcvt && !pwdlt && !pwccct

		if pwAcceptable != valid {
			// The tests do not match the expected result.  Log error
			fmt.Printf("Incorrect: password (%s) result (%v) expected (%v)\n",
				passwd, pwAcceptable, valid)

			fmt.Printf("pwlt (%v) pwact (%v) pwcvt (%v) pwdlt (%v) pwcct (%v)\n\n\n",
				pwlt, pwact, pwcvt, pwdlt, pwccct)
		}
	}
}
*/

func main1() {
	pwAcceptable := false

	// Loop through all passwords in the test set
	for passwd, valid := range passwordTestSet {
		pwAcceptable = false

		// Check the password length
		if passwordLengthTest(passwd) == true {
			// Check if the password contains only approved characters
			if passwordApprovedCharactersTest(passwd) == true {
				// Check if the password contains at least one vowel
				if passwordContainsVowelTest(passwd) == true {
					// Check if the password contains consecutive occurrences of the same letter
					if passwordDuplicateLettersTest(passwd) == false {
						// Check if the password contains consecutive vowles of consonants
						if passwordContainsConsecutiveCharacterTest(passwd, ConsecutiveCharacterLimit) == false {
							pwAcceptable = true
						}
					}
				}
			}
		}

		printResult(passwd, pwAcceptable)

		if pwAcceptable != valid {
			// The tests do not match the expected result.  Log error
			fmt.Printf("Result is inccorect: password (%s) result (%v) expected (%v)\n",
				passwd, pwAcceptable, valid)
		}
	}
}
