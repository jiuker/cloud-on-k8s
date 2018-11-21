package version

import (
	"fmt"
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	type args struct {
		version string
	}
	tests := []struct {
		name    string
		args    args
		want    *Version
		wantErr bool
	}{
		{
			name: "simple version",
			args: args{version: "1.2.3"},
			want: &Version{Major: 1, Minor: 2, Patch: 3},
		},
		{
			name: "version with label",
			args: args{version: "1.2.3-foo"},
			want: &Version{Major: 1, Minor: 2, Patch: 3, Label: "foo"},
		},
		{
			name: "version with dotted label",
			args: args{version: "1.2.3-f.oo"},
			want: &Version{Major: 1, Minor: 2, Patch: 3, Label: "f.oo"},
		},
		{
			name: "version with dashed label",
			args: args{version: "1.2.3-f.o-o"},
			want: &Version{Major: 1, Minor: 2, Patch: 3, Label: "f.o-o"},
		},
		{
			name: "zero version",
			args: args{version: "0.0.0"},
			want: &Version{},
		},
		{
			name:    "invalid major version",
			args:    args{version: "a.0.0"},
			wantErr: true,
		},
		{
			name:    "invalid minor version",
			args:    args{version: "0.a.0"},
			wantErr: true,
		},
		{
			name:    "invalid patch version",
			args:    args{version: "0.0.a"},
			wantErr: true,
		},
		{
			name:    "invalid label",
			args:    args{version: "0.0.0.label"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMustParse(t *testing.T) {
	type args struct {
		version string
	}
	tests := []struct {
		name string
		args args
		want Version
	}{
		{name: "simple", args: args{"7.0.0"}, want: Version{Major: 7}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MustParse(tt.args.version); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MustParse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVersion_IsSameOrAfter(t *testing.T) {
	dimensions := 2
	sortedVersionsLength := dimensions * dimensions * dimensions
	sortedVersions := make([]Version, sortedVersionsLength)
	for major := 0; major < dimensions; major++ {
		for minor := 0; minor < dimensions; minor++ {
			for patch := 0; patch < dimensions; patch++ {
				index := major*dimensions*dimensions + minor*dimensions + patch
				sortedVersions[index] = Version{Major: major, Minor: minor, Patch: patch}
			}
		}
	}

	type test struct {
		name  string
		v     Version
		other Version
		want  bool
	}
	tests := make([]test, 0)

	for i, version := range sortedVersions {
		for j := 0; j < i; j++ {
			other := sortedVersions[j]
			tests = append(tests, test{
				name:  fmt.Sprintf("%v > %v", version, other),
				v:     version,
				other: other,
				want:  true,
			})
		}
		tests = append(tests, test{
			name:  fmt.Sprintf("%v=%v", version, version),
			v:     version,
			other: version,
			want:  true,
		})
		for j := i + 1; j < sortedVersionsLength; j++ {
			other := sortedVersions[j]
			tests = append(tests, test{
				name:  fmt.Sprintf("%v < %v", version, other),
				v:     version,
				other: other,
				want:  false,
			})
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.IsSameOrAfter(tt.other); got != tt.want {
				t.Errorf("Version.IsSameOrAfter() = %v, want %v", got, tt.want)
			}
		})
	}
}
