/*MIT licence
(c) Jan Piskorski
*/

package aiRequests
//package main

import (
	//sterowanie motorami, sensorami
	//"github.com/ldmberman/GoEV3/Motor"
	"github.com/janekjan/ugolibev3"
	//uruchamianie podprogramów
	"io/ioutil"
	"os/exec"
	"path/filepath"
	//operacje na słowach
	"strings"

	//do debugu
	"fmt"
	//yes/no
	"bufio"
	"os"
	//convert string->int
	"strconv"
)

//var waitForOutput bool = true
var waitForOutput bool = true

/*Runs the request.
 * r is the content of that request, prepared without brackets.
 */
func Run(r string, out bool) {
	waitForOutput = out
	_ = runRequest(r)
}

func runRequest(r string) int {
	r = strings.Trim(r, " \t\n().,!?")
	split := strings.Split(r, " ")
	if r != "" {
		if (in(split, "finish")) || (in(split, "stop")) {
			//finish the execution
			fmt.Println("I'm done!")
			return 1
		} else if (in(split, "repeat")) || (in(split, "again")) {
			//let us proceed to repeating
			return 3
			//use with care!!!
		} else if split[0] == "*base" {
			//base functions
			runBase(split[1])
			return 0
		} else if (in(split, "when")) || (in(split, "if")) {
			//conditional
			if (in(split, "see")) || (in(split, "near")) {
				if !(isNear()) {
					return 2
				} else {
					return 0
				}
			}
			if in(split, "touch") {
				if !(touches()) {
					return 2
				} else {
					return 0
				}
			}
			//end conditionals
		} else { //begin normal commands
			//the advanced way, disabled by now.
			/*
				commands := getCommands()
				contains := false
				var theCommand string
				//does command exist?
				for _, v := range split {
					if in(commands, v) {
						contains = true
						theCommand = v
					}
				}
				var theOption string
				if contains {
					opt := getOptions()
					contains = false
					for _, v := range split {
						if in(commands, v) {
							contains = true
							theOption = v
						}
					}
				}
				if !(contains) {
					learn(r)
				} else {
					//exisits, let's proceed

				} */
			commands := getCommands()
			if !(in(commands, split[0])) {
				learn(split[0], split[1])
				return 0
			}
			opts := getOptions(split[0])
			if !(in(opts, split[1])) {
				learn(split[0], split[1])
				return 0
			}
			theCommand := getInstructions(split[0], split[1])
			lastRet := 0
			for _, v := range theCommand {
				if lastRet == 1 {
					break
				}
				if lastRet == 2 {
					continue
				}
				if lastRet == 3 {
					_ = runRequest(r)
					return 0
				}
				/*
				if v[len(v)-1:] == "#" {
					v =v[:len(v)-1] + "\\#"
				}
				*/
				lastRet = runRequest(v)
			}
		}
	}

	return 0
}

func getCommands() []string {
	classiflist, err := filepath.Glob("comm/*")
	errorcheck(err)
	for i, v := range classiflist {
		if strings.Contains(v, "/") {
			classiflist[i] = v[strings.Index(v, "/")+1:]
		}
	}
	return classiflist

}
func getOptions(command string) []string {
	classiflist, err := filepath.Glob("comm/" + command + "/*")
	errorcheck(err)
	//informatycznie odpowiednia długość wycinka
	for i, v := range classiflist {
		if strings.Contains(v, "/") {
			classiflist[i] = v[strings.LastIndex(v, "/")+1:]
		}
	}
	return classiflist
}
func getInstructions(c, o string) []string {
	dat, err := ioutil.ReadFile("comm/" + c + "/" + o)
	errorcheck(err)
	data := string(dat)
	//informatycznie odpowiednia długość wycinka
	lines := make([]string, 16383)
	lines = strings.Split(data, "\n")
	//log("Pobrałem dane z bazy danych...")
	return lines
}

func isNear() bool {
	sensorval := ugolibev3.ReadSensorValue(0)
	val, err := strconv.ParseInt(sensorval[:len(sensorval)-1], 10, 16)
	errorcheck(err)
	if val < 50 {
		return true
	}
	return false
}
func touches() bool {
	sensorval := ugolibev3.ReadSensorValue(0)
	val, err := strconv.ParseInt(sensorval[:len(sensorval)-1], 10, 16)
	errorcheck(err)
	if val == 1 {
		return true
	}
	return false
}

func runBase(c string) {
	cmd := exec.Command("./base/" + c)
	if waitForOutput {
		out, err := cmd.Output()
		errorcheck(err)
		fmt.Printf("%s\n", out)
	} else {
		err := cmd.Start()
		errorcheck(err)
	}
}

func in(sl []string, s string) bool {
	var rlyin bool = false
	for _, v := range sl {
		if v == s {
			rlyin = true
		}
	}
	return rlyin
}
func errorcheck(e error) {
	if e != nil {
		panic(e)
	}
}
func GetQuery() string {
	var inp string
	fmt.Printf("$> ")
	//źródło to konsola
	scnr := bufio.NewScanner(os.Stdin)
	//skanujemy i wynik do zmiennej
	scnr.Scan()
	inp = scnr.Text()
	//fmt.Printf("%s\n", scnr.Text())
	return inp
}
func YesNoQuestion(q string) bool {
	var inp string
	fmt.Printf("%v ", q)
	scnr := bufio.NewScanner(os.Stdin)
	scnr.Scan()
	inp = scnr.Text()
	var o bool
	if ((inp[:1] == "t") || (inp[:1] == "y")) || ((inp[:1] == "Y") || (inp[:1] == "T")) {
		//log("Pytanie", q, " tak/nie. Udzieliłe(a)ś odp. twierdzącej!")
		o = true
	} else {
		//log("Pytanie", q, " tak/nie. Udzieliłe(a)ś odp. przeczącej!")
		o = false
	}
	return o
}
func learn(c, o string) {
	fmt.Println("I don't know how to "+c, o)
	if YesNoQuestion("Will you teach me now? ") {
		contents := make([]string, 160)
		count := 0
		for {
			line := GetQuery()
			if strings.Contains(line, "all") || strings.Contains(line, "everything") {
				break
			}
			contents[count] = line
			count += 1
		}
		newInstruction(contents, c, o)
		if YesNoQuestion("Would you like me to do it now? ") {
			runRequest(c+" "+o)
		}				
	}
}
func newInstruction(contents []string, c, o string) {
	preparedir(c)
	f, err := os.OpenFile("comm/"+c+"/"+o, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	defer f.Close()
	errorcheck(err)
	n2 := 0
	//zakodowanie linii
	for _, v := range contents {
		d2 := []byte(v + "\n")
		//piszemy
		n3, err := f.Write(d2)
		errorcheck(err)
		n2 += n3
	}
	defer fmt.Printf("wrote %d bytes\n", n2)
}
func preparedir(dir string) {
	if b, _ := dirExists("comm/"+dir); !(b) {
		cmd := exec.Command("mkdir", "comm/"+dir)
		err := cmd.Run()
		errorcheck(err)
	}
}
// exists returns whether the given file or directory exists or not
func dirExists(path string) (bool, error) {
    _, err := os.Stat(path)
    if err == nil { return true, nil }
    if os.IsNotExist(err) { return false, nil }
    return false, err
}

func main() {
	fmt.Println(getCommands())
	fmt.Println(getOptions("go"))
	fmt.Println(getInstructions("go", "forward"))
	Run("go forward first, please!", true)
	//newInstruction([]string{"*base stop-motors", ""}, "stop", "#")
	//fmt.Println(getInstructions("go","forward"))
	Run("go backwards", true)
}
