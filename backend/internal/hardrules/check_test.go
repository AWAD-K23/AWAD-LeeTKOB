package hardrules

import (
	"context"
	"os"
	"testing"

	"leetkob/backend/internal/sandbox"
)

func TestBanListFormat(t *testing.T) {
	for _, test := range []struct {
		text  string
		count int
	}{
		{"Do not call `std::sort`", 1}, {"Do not use `qsort` or `std::stable_sort`", 2},
		{"Use merge sort", 0}, {"Use `std::sort`", 0},
	} {
		if got := bannedNames(test.text); len(got) != test.count {
			t.Fatalf("%q: %v", test.text, got)
		}
	}
}

func TestSyntaxTraps(t *testing.T) {
	if os.Getenv("SANDBOX_TESTS") != "1" {
		t.Skip("requires built sandbox image")
	}

	for _, test := range []struct{ name, source, verdict string }{
		{"real call", "#include <algorithm>\nint main(){int a[2]; std::sort(a,a+2);}", "violation"},
		{"comment", "// std::sort(a,a+n);\nint main(){}", "pass"},
		{"string", `int main(){const char* s="std::sort(a,a+n)";}`, "pass"},
		{"own namespace", `namespace mine {void sort(int*,int*){}} int main(){int a[2];mine::sort(a,a+2);}`, "pass"},
		{"include for swap", "#include <algorithm>\nint main(){int a=1,b=2;std::swap(a,b);}", "pass"},
		{"using declaration", "#include <algorithm>\nusing std::sort; int main(){int a[2];sort(a,a+2);}", "violation"},
		{"macro", "#include <algorithm>\n#define S std::sort\nint main(){int a[2];S(a,a+2);}", "violation"},
	} {
		t.Run(test.name, func(t *testing.T) {
			got, err := Check(context.Background(), sandbox.Runner{}, []byte(test.source), []Constraint{{ID: "C1", Text: "Do not call `std::sort`"}, {ID: "C2", Text: "Use merge sort"}})

			if err != nil {
				t.Fatal(err)
			}

			if len(got) != 2 || got[0].Verdict != test.verdict || got[1].Source != "semantic_pending" {
				t.Fatalf("%+v", got)
			}

			if test.verdict == "violation" && (len(got[0].Evidence) != 1 || got[0].Evidence[0].Text == "") {
				t.Fatal("missing exact evidence")
			}
		})
	}
}

func TestStatementAndTypeBans(t *testing.T) {
	if os.Getenv("SANDBOX_TESTS") != "1" {
		t.Skip("requires built sandbox image")
	}

	for _, test := range []struct{ name, source, ban, verdict string }{
		{"for loop", "int main(){for(int i=0;i<2;i++) {}}", "for", "violation"},
		{"range loop", "int main(){int a[2]={};for(int x:a) {}}", "for", "violation"},
		{"while loop", "int main(){int x=2;while(x--) {}}", "while", "violation"},
		{"do loop", "int main(){int x=2;do {} while(x--);}", "do", "violation"},
		{"goto", "int main(){goto done;done:;}", "goto", "violation"},
		{"loop text trap", "// for while goto\nint main(){auto s=\"for while goto\";}", "for", "pass"},
		{"vector", "#include <vector>\nint main(){std::vector<int> a;}", "std::vector", "violation"},
		{"container alias", "#include <vector>\nusing Numbers=std::vector<int>;\nint main(){Numbers a;}", "std::vector", "violation"},
		{"own container", "namespace mine{template<class T>struct vector{};}\nint main(){mine::vector<int> a;}", "std::vector", "pass"},
		{"unused include", "#include <vector>\nint main(){int a[2]={};}", "std::vector", "pass"},
	} {
		t.Run(test.name, func(t *testing.T) {
			got, err := Check(context.Background(), sandbox.Runner{}, []byte(test.source),
				[]Constraint{{ID: "rule", Text: "Do not use `" + test.ban + "`"}},
			)
			if err != nil || len(got) != 1 || got[0].Verdict != test.verdict {
				t.Fatalf("checks=%+v err=%v", got, err)
			}
		})
	}
}
