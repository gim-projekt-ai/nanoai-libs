/*MIT licence
(c) Bartosz Deptuła
*/

//package main
package convert_new
import (
	//podstawowa biblioteka
	"bufio"
	"fmt"
	"os"
	//biblioteka stemmująca
	"github.com/kljensen/snowball"
	//biblioteka stringów
	"strings"
)

//funkcje pomonicnicze do biblioteki "strings"
func GetLines(s string) []string {
	return strings.Split(s, " ")
}
func GetLine(s string, LineIndex int) string {
	return GetLines(s)[LineIndex]
}

//inne funkcje pomocnicze

//funkcja obsługująca komendy od użytkownika
func komendy(kom string) {
	if kom == "/help" {
		fmt.Println("NANO-AI HELP:\n/info - information about project\n/exit - exit the program\nend")
	}
	if kom == "/info" {
		fmt.Println("MIT licence\n(c)\nNLP: Bartosz Deptuła\nNANO-AI: Jan Piskorski\nend")
	}
	if kom == "/exit" {
		fmt.Println("Exiting...")
		os.Exit(0)
	}
}

//funkcja odpowiadająca na powitania
func witanie(zdanie string) bool {
	strings.Replace(zdanie, ".", "", -1)
	zdanie = strings.ToLower(zdanie)
	powitania := []string{"hi", "hello", "good morning", "good afternoon",
		"good evening", "welcome", "guten morgen", "guten tag", "guten abend",
		"hi nano", "hi nano-ai"}
	for i := range powitania {
		if powitania[i] == zdanie {
			return true
		}
	}
	return false
}

//funkcja odpowiadająca na pożegnania
func zegnanie(zdanie string) {
	strings.Replace(zdanie, ".", "", -1)
	zdanie = strings.ToLower(zdanie)
	pozegnania := []string{"good bye", "goodbye", "bye", "see you", "f*ck you",
		"ja piórkuję", "auf wiedersehen", "auf wiederschauen", "ruhe", "hau ab"}
	for i := range pozegnania {
		if pozegnania[i] == zdanie {
			fmt.Println("Good bye. Thank you for conversation. ", zdanie, " too.")
			fmt.Println("Exiting...")
			os.Exit(0)
		}
	}
}

//formatowanie znacznikowe w zmiennych y
func yEditing(y, wynikEdycji string) string {
	if len(y) == 0 {
		return wynikEdycji
	} else {
		return wynikEdycji + "," + y
	}
}

//dodawanie nawiasu do słowa
func naw(slowo, nawias string) string {
	slowoWNawiasie := slowo + "(" + nawias + ")"
	return slowoWNawiasie
}

func snowballer(tablica []string) ([]string, error) {
	var stemmed string
	var err error
	for elementTablicy := range tablica {
		stemmed, err = snowball.Stem(tablica[elementTablicy], "english", true)
		tablica[elementTablicy] = stemmed
	}
	return tablica, err

}

//usówanie przedrostków
func usunPrzedrostek(tablica, obiektyWykluczone []string) (int, []string) {
	miejscePrzedrostka := -1
	iloscPrzedrostkow := 0
	for xTablicy := range tablica {
		for xObiektowWykluczonych := range obiektyWykluczone {
			if tablica[xTablicy] == obiektyWykluczone[xObiektowWykluczonych] {
				copy(tablica[xTablicy:], tablica[xTablicy+1:])
				tablica[len(tablica)-1] = ""
				iloscPrzedrostkow++
				if xTablicy != 0 {
					miejscePrzedrostka = xTablicy
				}
			}
		}
	}
	tablicaOut := make([]string, len(tablica)-iloscPrzedrostkow)
	copy(tablicaOut, tablica[:len(tablicaOut)])
	return miejscePrzedrostka, tablicaOut
}

//FUNKCJA GŁÓWNA NLP

func Format(zdanie string) string {
	komendy(zdanie)
	if witanie(zdanie) {
		return "*hi"
	}
	zegnanie(zdanie)
	if len(GetLines(zdanie)) == 1 {
		if zdanie == "" {
			return "*empty"
		}
		return strings.ToLower(zdanie)
	}
	var x1, x2, x3, y1, y2, y3 string
	znaki := strings.NewReplacer(".", "", ",", "", ":", "", ";", "", "(", "", ")", "",
		"!", "", "?", "")
	zdanie = znaki.Replace(zdanie)
	//Sprawdź czy pytanie
	//formatowanie znacznikowe --> OKOLICZNIKOWE
	znacznikOkolicznikowy := []string{"Where", "When", "How"}
	x1 = GetLine(zdanie, 0)
	dodajZnacznikOkolicznika := false
	for i := range znacznikOkolicznikowy {
		if znacznikOkolicznikowy[i] == x1 {
			usunZaimek := strings.NewReplacer("Where ", "", "When ", "", "How ", "")
			zdanie = usunZaimek.Replace(zdanie)
			dodajZnacznikOkolicznika = true
		}
	}
	//rozwinięcie skrótów
	zdanie = strings.ToLower(zdanie)
	zmienSkrot := strings.NewReplacer("i'm", "i am", "they're", "they are",
		"we're", "we are", "you're", "you are", "she's", "she is", "he's", "he is",
		"it's", "it is", "i've", "i have", "they've", "they have",
		"we've", "we have", "you've", "you have", "n't", " *not", "not", "*not")
	zdanie = zmienSkrot.Replace(zdanie)
	//Oznaczanie zaprzeczenia
	y1 = ""
	if GetLine(zdanie, 1) == "*not" {
		y1 = "*not"
		usunZaprzeczenie := strings.NewReplacer(" *not", "")
		zdanie = usunZaprzeczenie.Replace(zdanie)
	}
	copyx2 := ""
	if len(GetLines(zdanie)) > 2 {
		copyx2 = GetLine(zdanie, 2)
	}
	//Wyrzucenie zdania do tablicy
	zdanieTablicain := strings.Fields(zdanie)
	//Stemowanie zdania
	zdanieTablicain, _ = snowballer(zdanieTablicain)
	//Usuwanie przedrostków
	przedrostki := []string{"A", "a", "An", "an", "The", "the"}
	miejscePrzedrostka, zdanieTablica := usunPrzedrostek(zdanieTablicain, przedrostki)
	dlugoscZdanie := len(zdanieTablica) - 1
	x1 = zdanieTablica[0]
	x2, y2 = zdanieTablica[1], ""

	//formatowanie znacznikowe
	//przypisywanie tablic znacznikowych
	znacznikBe := []string{"am", "is", "are", "do", "doe"}
	znacznikOsobowy := []string{"what", "who"}
	znacznikOdnosnikowy := []string{"my", "your", "his", "her", "its", "our",
		"your", "them"}
	dodajZnacznikBe := false
	//osobaWhy := "#"
	//	if dlugoscZdanie > 2 {
	//		osobaWhy = zdanieTablicain[2]
	//	}
	//odnonik do-osobowy
	for i := range znacznikOdnosnikowy {
		if znacznikOdnosnikowy[i] == x1 {
			y1 = yEditing(y1, x1)
			x1 = x2
			if dlugoscZdanie > 2 {
				x2 = zdanieTablica[2]
				//osobaWhy = zdanieTablicain[3]
			} else {
				x2 = "#"
			}
		}
	}

	//pytanie typu tak/nie
	for i := range znacznikBe {
		if x1 == znacznikBe[i] {
			x1 = x2
			x2 = znacznikBe[i]
			dodajZnacznikBe = true
		}
		if x2 == znacznikBe[i] {
			x2 = "*be"
		}
	}
	whyEdited := false
	if x1 == "whi" {
		/*		x1 = osobaWhy
				y2 = yEditing(y2, "*why")
				if y2 == ",*why" {
					y2 = "*why"
				}
		*/ //dodajZnacznikBe = true
		y1 = "*why"
		x1 = zdanieTablicain[2]
		whyEdited = true
		//zdanieTablicain[1] = ""
		/*
		copy(zdanieTablicain[2:], zdanieTablicain[3:])
		zdanieTablicain[2] = ""
		tablicaTemp := make([]string, len(zdanieTablicain)-1)
		copy(tablicaTemp, zdanieTablicain[:len(tablicaTemp)])
		zdanieTablicain = make([]string, len(tablicaTemp))
		copy(zdanieTablicain, tablicaTemp)
		dlugoscZdanie--
		*/
	}
	//inne
	if x1 == "i" { //"I" w x1
		x1 = "*i"
	} else if x1 == "pleas" { //"Please" w x1
		x1 = "*r"
	}
	//dla wszystkich parametrów tablicy znacznikOsobowy zamień x1 na "*q"
	for i := range znacznikOsobowy {
		if x1 == znacznikOsobowy[i] {
			x1 = "*q"
		}

	}

	if dlugoscZdanie > 2 {
		fmt.Println(zdanieTablicain[2])
	}

	/*	copy(zdanieTablicain[2:], zdanieTablicain[3:])
		zdanieTablicain[2] = ""
		tablicaTemp := make([]string, len(zdanieTablicain)-1)
		copy(tablicaTemp, zdanieTablicain[:len(tablicaTemp)])
		zdanieTablicain = make([]string, len(tablicaTemp))
		copy(zdanieTablicain, tablicaTemp)
		//dlugoscZdanie
	*/
	//Wyrzucanie do stringa
	if dlugoscZdanie > 2 {
		if miejscePrzedrostka != -1 {
			y2 = yEditing(y2, strings.Join(zdanieTablicain[2:miejscePrzedrostka], ","))
			x3 = zdanieTablicain[dlugoscZdanie]
			y3 = strings.Join(zdanieTablicain[miejscePrzedrostka:dlugoscZdanie], ",")
		} else if miejscePrzedrostka == -1 {
			y2 = yEditing(y2, strings.Join(zdanieTablicain[2:dlugoscZdanie-1], ","))
			x3 = zdanieTablicain[dlugoscZdanie]
			y3 = zdanieTablicain[dlugoscZdanie-1]
		}
	} else if dlugoscZdanie == 2 {
		x3, y3 = zdanieTablicain[dlugoscZdanie], ""
	} else if dlugoscZdanie == 1 {
		x3, y3 = "#", ""
	} else {
		x2, x3, y3 = "#", "#", ""
	}

	strings.Replace(y2, ",", "", 0)
	
	if whyEdited {
		whyWy2 := strings.NewReplacer(copyx2+",", "", copyx2, "")
		y2 = whyWy2.Replace(y2)
		whyWy3 := strings.NewReplacer(copyx2+",", "", copyx2, "")
		y3 = whyWy3.Replace(y3)
	} 
	//oddanie wyników
	var zdanieout string //zmienna zdanieout jest wynikiem wyjściowym
	if dodajZnacznikOkolicznika {
		if x3 == "#" {
			x3 = ""
		}
		zdanieout = naw(x1, y1) + " " + naw(x2, y2) + " " + naw("*q", yEditing(y3, x3))
	} else {
		if dodajZnacznikBe {
			zdanieout = naw(x1, y1) + " " + naw(x2, yEditing(y2, "*rly")) + " " + naw(x3, y3)
		} else {
			zdanieout = naw(x1, y1) + " " + naw(x2, y2) + " " + naw(x3, y3)
		}
	}
	return zdanieout
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

/* BŁĘDY


*/
func main() {
	for true {
		testowe_zdanie := GetQuery()
		fmt.Println(Format(testowe_zdanie))
	}
	//fmt.Println(len(GetLines("")))
}
