package csv

import (
	"log"
	"os"

	"github.com/gocarina/gocsv"
	"github.com/nathanhollows/FresherFocus/models"
)

// Parse static/database.csv into Students
func Parse() (models.Students, error) {
	// open csv file
	f, err := os.Open("static/database.csv")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	students := models.Students{}
	if err := gocsv.UnmarshalFile(f, &students); err != nil { // Load clients from file
		log.Fatal(err)
	}

	// Sort students by last name
	students.Sort()

	return students, nil
}

// // Read in the data
// func readFYdata() FirstYears {
// 	items := []list.Item{}
// 	// List the files in the current directory
// 	files, err := ioutil.ReadDir("Raw Data")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Ask the user which CSV file to read
// 	for _, file := range files {
// 		if strings.HasSuffix(file.Name(), ".csv") {
// 			items = append(items, item(file.Name()))
// 		}
// 	}

// 	filename := view("Which file would you like to process?", items)

// 	f, err := os.Open("Raw Data/" + filename)
// 	if err != nil {
// 		return nil
// 	}

// 	// Make sure we close it later
// 	defer f.Close()

// 	firstYears := FirstYears{}

// 	if err := gocsv.UnmarshalFile(f, &firstYears); err != nil { // Load clients from file
// 		log.Fatal(err)
// 	}

// 	// Set each first year's college and local status
// 	for i := range firstYears {
// 		firstYears[i].Local = true
// 		firstYears[i].Admitted = true
// 		firstYears[i].SecondSemNotDeclared = false
// 		firstYears[i].ProfessionalProgram = false
// 		firstYears[i].Distance = false
// 		firstYears[i].SplitPapers()
// 	}

// 	firstYears.processData()
// 	return firstYears
// }

// // Filter by Locals
// func (firstYears FirstYears) processData() {
// 	items := []list.Item{}
// 	items = append(items, item("No"))
// 	items = append(items, item("Yes"))
// 	removeSS := view("Do you want to remove Second Semester student who haven't declared?", items)
// 	CheckedSecondSemester := removeSS == "Yes"

// 	removeID := view("Do you want to remove any specific student IDs?", items)
// 	if removeID == "Yes" {
// 		var input string
// 		// Ask the user for the student_ids
// 		fmt.Println("Enter the student_ids you want to remove, separated by spaces")
// 		fmt.Print("Student IDs: ")
// 		_, err := fmt.Scan(&input)
// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		studentIDs := strings.Split(input, " ")
// 		for i := range studentIDs {
// 			studentIDs[i] = strings.TrimSpace(studentIDs[i])
// 		}
// 	}

// 	for i, firstYear := range firstYears {
// 		// Mark college students as not local
// 		if firstYear.isInCollege() {
// 			(firstYears)[i].College = true
// 			(firstYears)[i].Local = false
// 		}

// 		// Remove students who aren't yet admitted
// 		if !firstYear.isAdmitted() {
// 			(firstYears)[i].Admitted = false
// 			(firstYears)[i].Local = false
// 		}

// 		// Ask the user if they want to check Second Semester Students
// 		// Check if a student is a second semester student and hasn't been admitted
// 		if CheckedSecondSemester && firstYear.isSecondSemesterUndeclared() {
// 			(firstYears)[i].SecondSemNotDeclared = true
// 		}

// 		// Check if the student is in a professional program
// 		if firstYear.isInProfessionalProgram() {
// 			(firstYears)[i].ProfessionalProgram = true
// 			(firstYears)[i].Local = false
// 		}

// 		// Check if the student is a distance student
// 		if !firstYear.isDistanceStudent() {
// 			(firstYears)[i].Distance = true
// 			(firstYears)[i].Local = false
// 		}

// 		// Check if the student is only enrolled in summer school
// 		if firstYear.isSummerSchoolOnly() {
// 			(firstYears)[i].SummerSchool = true
// 			(firstYears)[i].Local = false
// 		}

// 		// Check if the student is likely in their second year
// 		if firstYear.isSecondYearPlus() {
// 			(firstYears)[i].Local = false
// 			(firstYears)[i].SecondPlusYear = false
// 			continue
// 		}
// 	}
// }

// // Check if the first year is in a college
// func (firstYear *FirstYear) isInCollege() bool {
// 	return firstYear.Accommodation == "Living in a residential college"
// }

// // Check the student is admitted
// func (firstYear *FirstYear) isAdmitted() bool {
// 	// Check if Admission contains "pending"
// 	if strings.Contains(strings.ToLower(firstYear.Admission), "pending") {
// 		return false
// 	}
// 	// Check if Admission contains "not verified"
// 	if strings.Contains(strings.ToLower(firstYear.Admission), "not verified") {
// 		return false
// 	}
// 	// Check if Admission contains "non-matriculated"
// 	if strings.Contains(strings.ToLower(firstYear.Admission), "non-matriculated") {
// 		return false
// 	}

// 	return true
// }

// // Check if the student has declared
// func (firstYear *FirstYear) hasDeclared() bool {
// 	// If we're not past the 28th of January of the current year, we can't check if they've declared
// 	if time.Now().Month() == time.January && time.Now().Day() < 28 {
// 		return true
// 	}

// 	return firstYear.Declared == "Y"
// }

// // Check if the student is a second semester student and hasn't declared
// func (firstYear *FirstYear) isSecondSemesterUndeclared() bool {
// 	return firstYear.S2_start == "Y" && !firstYear.hasDeclared()
// }

// // Check if the student is in a professional program
// func (firstYear *FirstYear) isInProfessionalProgram() bool {
// 	professionalPrograms := []string{"BDS", "BPharm", "BPhty", "BMLSc", "MB ChB"}

// 	for _, program := range professionalPrograms {
// 		if strings.Contains(firstYear.Prog1, program) {
// 			return true
// 		}
// 		if strings.Contains(firstYear.Prog2, program) {
// 			return true
// 		}
// 		if strings.Contains(firstYear.Prog3, program) {
// 			return true
// 		}
// 	}
// 	return false

// }

// // Check if the student is a distance student
// func (firstYear *FirstYear) isDistanceStudent() bool {
// 	return firstYear.Study_country != "New Zealand"
// }

// // Split the papers into a slice
// func (firstYear *FirstYear) SplitPapers() {
// 	firstYear.CurrentPapers = strings.Split(firstYear.Papers, "; ")
// }

// // Check if the student is likely in their second year
// func (firstYear *FirstYear) isSecondYearPlus() bool {
// 	if strings.Contains(strings.ToLower(firstYear.Prior_activity), "university") {
// 		return false
// 	}
// 	// Check how many papers are higher than 100
// 	var papersOver100 int
// 	for _, paper := range firstYear.CurrentPapers {
// 		if paper == "" {
// 			continue
// 		}
// 		paperNumber := strings.TrimSpace(paper)[4:7]
// 		level, err := strconv.Atoi(paperNumber)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		if level > 199 {
// 			papersOver100++
// 		}
// 	}

// 	// If there are more than 5 papers over 100, the student is likely in their second year or beyond
// 	return papersOver100 > 4
// }

// // Check if the student is only enrolled in summer school
// func (firstYear *FirstYear) isSummerSchoolOnly() bool {
// 	for _, paper := range firstYear.CurrentPapers {
// 		if paper == "" {
// 			continue
// 		}
// 		period := strings.TrimSpace(paper)[len(paper)-4:]
// 		if period != "(SS)" {
// 			return false
// 		}
// 	}
// 	return true
// }

// // FilterLocals returns a slice of local students
// func (firstYears FirstYears) filterLocals() FirstYears {
// 	var locals FirstYears
// 	for _, firstYear := range firstYears {
// 		if firstYear.Local {
// 			locals = append(locals, firstYear)
// 		}
// 	}
// 	return locals
// }

// // outputCSV outputs the data to a CSV file
// func (firstYears FirstYears) outputCSV(filename string) {
// 	date := time.Now().Format("2006-01-02")
// 	filename = fmt.Sprintf("Reports and Databases/%s %s.csv", filename, date)
// 	file, err := os.Create(filename)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer file.Close()
// 	err = gocsv.MarshalFile(&firstYears, file)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	var style = quitTextStyle.Padding(0, 0, 0, 0)
// 	message := fmt.Sprintf("Wrote %d records to %s", len(firstYears), filename)
// 	fmt.Println(style.Render(message))

// }

// // FindStudentByID finds a student by their ID
// func (firstYears FirstYears) findStudentByID(id string) FirstYear {
// 	for _, firstYear := range firstYears {
// 		if firstYear.Student_id == id {
// 			return *firstYear
// 		}
// 	}
// 	return FirstYear{}
// }

// // FindStudentsByPaper finds students by a paper code
// func (firstYears FirstYears) findStudentsByPaper(paper string) FirstYears {
// 	var students FirstYears
// 	for _, firstYear := range firstYears {
// 		for _, currentPaper := range firstYear.CurrentPapers {
// 			if currentPaper[:7] == paper {
// 				students = append(students, firstYear)
// 			}
// 		}
// 	}
// 	return students
// }
