package utils

import "testing"

func TestPasswordValid(t *testing.T) {
	mockPwd := "LPLL!!@SSDDDasdasd111"
	
	err := ValidatePwd(mockPwd)

	if err != nil {
		t.Error("Should be no error when password meets all criteria", err)
	}
}

func TestPasswordWithoutCapitalLettersRejected (t *testing.T) {
	mockPwd := "ooooppddd222!!"
	
	err := ValidatePwd(mockPwd)

	if err == nil {
		t.Error("Error should exist when no capital letter")
	}
}

func TestPasswordWithoutSmallLettersRejected (t *testing.T) {
	mockPwd := "OOOIIIDDD523!!"
	
	err := ValidatePwd(mockPwd)

	if err == nil {
		t.Error("Error should exist when no small letter")
	}
}

func TestPasswordWithoutSpecialChars (t *testing.T) {
	mockPwd := "KUOSSSDDLmnmnmn093"
	
	err := ValidatePwd(mockPwd)

	if err == nil {
		t.Error("Error should exist when no special char")
	}
}

