package substratum

import "testing"

func TestGetABINameForNumberRegister_Success(t *testing.T) {
	tests := map[string]string{
		"x0":  "zero",
		"x1":  "ra",
		"x8":  "fp", // special case as it has two names, s0 and fp. We assume "fp" is always returned,
		"x31": "t6",
	}

	for k, v := range tests {
		result, err := GetABINameForNumberRegister(k)
		if err != nil {
			t.Errorf("got an error trying to get an ABI name: %v", err)
			continue
		}

		if result != v {
			t.Errorf("expected GetABINameForNumberRegister(%s) == %s, but got %s", k, v, result)
			continue
		}
	}
}

func TestGetABINameForNumberRegister_Failure(t *testing.T) {
	tests := []string{"a", "x32", "x-1", "x0x0"}
	for _, test := range tests {
		_, err := GetABINameForNumberRegister(test)
		if err == nil {
			t.Errorf("expected an error for GetABINameforNumberRegister(%s) but didn't get one", test)
			continue
		}
	}
}

func TestGetNumberRegisterForABIName_Success(t *testing.T) {
	tests := map[string]string{
		"zero": "x0",
		"fp":   "x8",
		"s0":   "x8",
		"t6":   "x31",
		"ZERO": "x0",
	}

	for k, v := range tests {
		result, err := GetNumberRegisterForABIName(k)

		if err != nil {
			t.Errorf("got an error from GetNumberRegisterForABIName(%s): %v", k, err)
			continue
		}

		if result != v {
			t.Errorf("expected GetNumberRegisterForABIName(%s) == %s, but got %s", k, v, result)
			continue
		}
	}
}

func TestGetNumberRegisterForABIName_Failure(t *testing.T) {
	tests := []string{"a", "0", "pf", "zer0"}
	for _, test := range tests {
		_, err := GetABINameForNumberRegister(test)
		if err == nil {
			t.Errorf("expected an error for GetNumberRegisterForABIName(%s) but didn't get one", test)
			continue
		}
	}
}
