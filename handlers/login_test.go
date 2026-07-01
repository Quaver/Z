package handlers

import "testing"

func TestParseHardwareIds(t *testing.T) {
	tests := []struct {
		name          string
		hardware      string
		expectedCpu   string
		expectedDisk  string
		expectedCombo string
	}{
		{
			name:          "all values",
			hardware:      "cpu|disk|combo",
			expectedCpu:   "cpu",
			expectedDisk:  "disk",
			expectedCombo: "combo",
		},
		{
			name:          "missing cpu",
			hardware:      "|disk|combo",
			expectedCpu:   "none",
			expectedDisk:  "disk",
			expectedCombo: "combo",
		},
		{
			name:          "missing disk",
			hardware:      "cpu||combo",
			expectedCpu:   "cpu",
			expectedDisk:  "none",
			expectedCombo: "combo",
		},
		{
			name:          "missing combo",
			hardware:      "cpu|disk",
			expectedCpu:   "cpu",
			expectedDisk:  "disk",
			expectedCombo: "none",
		},
		{
			name:          "empty",
			hardware:      "",
			expectedCpu:   "none",
			expectedDisk:  "none",
			expectedCombo: "none",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			hardware := parseHardwareIds(test.hardware)

			if hardware.CpuId != test.expectedCpu {
				t.Fatalf("expected cpu id %q, got %q", test.expectedCpu, hardware.CpuId)
			}

			if hardware.DiskId != test.expectedDisk {
				t.Fatalf("expected disk id %q, got %q", test.expectedDisk, hardware.DiskId)
			}

			if hardware.CpuDiskId != test.expectedCombo {
				t.Fatalf("expected cpu disk id %q, got %q", test.expectedCombo, hardware.CpuDiskId)
			}
		})
	}
}
