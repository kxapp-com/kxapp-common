package filez

import "testing"

func TestRestoreSymlinks(t *testing.T) {
	srcDir := "/Users/admin/Downloads/bin/iPhoneOS.sdk/.LinkPatch.json"
	destDir := "/Users/admin/Downloads/bin/iPhoneOS.sdk"
	e := RestoreSymlinks(srcDir, destDir, true)
	if e != nil {
		t.Error(e)
	}
}
