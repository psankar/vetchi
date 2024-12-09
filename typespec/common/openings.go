package common

type OpeningType string

const (
	// Any changes to here should be reflected also in the IsValid() method
	FullTimeOpening    OpeningType = "FULL_TIME_OPENING"
	PartTimeOpening    OpeningType = "PART_TIME_OPENING"
	ContractOpening    OpeningType = "CONTRACT_OPENING"
	InternshipOpening  OpeningType = "INTERNSHIP_OPENING"
	UnspecifiedOpening OpeningType = "UNSPECIFIED_OPENING"
)

func (o OpeningType) IsValid() bool {
	switch o {
	case FullTimeOpening,
		PartTimeOpening,
		ContractOpening,
		InternshipOpening,
		UnspecifiedOpening:
		return true
	}
	return false
}

type EducationLevel string

const (
	BachelorEducation    EducationLevel = "BACHELOR_EDUCATION"
	MasterEducation      EducationLevel = "MASTER_EDUCATION"
	DoctorateEducation   EducationLevel = "DOCTORATE_EDUCATION"
	NotMattersEducation  EducationLevel = "NOT_MATTERS_EDUCATION"
	UnspecifiedEducation EducationLevel = "UNSPECIFIED_EDUCATION"
)

func (e EducationLevel) IsValid() bool {
	return e == BachelorEducation || e == MasterEducation ||
		e == DoctorateEducation ||
		e == NotMattersEducation ||
		e == UnspecifiedEducation
}