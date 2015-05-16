package temple

import (
	"testing"
)

func TestStructGen(t *testing.T) {
	testCases := []struct {
		paths         []string
		expectedDeclr string
		expectedInit  string
		expectError   bool
	}{
		{
			// This should cause an error becuase show is the name of a directory
			// and a file
			paths: []string{
				"people/show",
				"people/show/foo",
			},
			expectError: true,
		},
		{
			// This should cause an error because 5 is not a valid variable or field name
			paths: []string{
				"people/show/5",
			},
			expectError: true,
		},
		{
			// This test case has two paths of depth two each
			paths: []string{
				"people/show",
				"people/index",
			},
			expectedDeclr: "var Templates struct {\nPeople struct {\nShow string\nIndex string\n}\n}\n",
			expectedInit:  "Templates.People.Show = TemplatesMap[\"people/show\"]\nTemplates.People.Index = TemplatesMap[\"people/index\"]\n",
			expectError:   false,
		},
		{
			// This test case more paths with varying depths
			paths: []string{
				"people/show/small",
				"people/show/large",
				"people/index",
				"books/index",
				"books/create/multiple",
				"books/create/single",
				"home",
			},
			expectedDeclr: `var Templates struct {
People struct {
Show struct {
Small string
Large string
}
Index string
}
Books struct {
Index string
Create struct {
Multiple string
Single string
}
}
Home string
}
`,
			expectedInit: `Templates.People.Show.Small = TemplatesMap["people/show/small"]
Templates.People.Show.Large = TemplatesMap["people/show/large"]
Templates.People.Index = TemplatesMap["people/index"]
Templates.Books.Index = TemplatesMap["books/index"]
Templates.Books.Create.Multiple = TemplatesMap["books/create/multiple"]
Templates.Books.Create.Single = TemplatesMap["books/create/single"]
Templates.Home = TemplatesMap["home"]
`,
			expectError: false,
		},
	}
	for i, tc := range testCases {
		tree, err := createTreeFromPaths(tc.paths)
		if tc.expectError && err == nil {
			t.Errorf("Failure at test case %d: Expected error but got none.", i)
		} else if !tc.expectError && err != nil {
			t.Errorf("Failure at test case %d: Unexpected error: %s", i, err.Error())
		}
		gotDeclr := generateStructDeclr("Templates", tree)
		if gotDeclr != tc.expectedDeclr {
			t.Errorf("Failure at test case %d: Declaration was incorrect.\n\tExpected: `%s`\n\tBut got:  `%s`", i, tc.expectedDeclr, gotDeclr)
		}
		if tree != nil {
			gotInit := generateStructInit("Templates", "TemplatesMap", tree.OrderedChildren)
			if gotInit != tc.expectedInit {
				t.Errorf("Failure at test case %d: Declaration was incorrect.\n\tExpected: `%s`\n\tBut got:  `%s`", i, tc.expectedInit, gotInit)
			}
		}
	}
}
