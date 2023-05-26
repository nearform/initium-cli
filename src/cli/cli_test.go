// package cli
//
// import (
// 	"os"
// 	"testing"
// 	"embed"
// )
//
// func TestEnvConfig(t *testing.T) {
//     cwd, err := os.Getwd()
// 	if err != nil {
// 		t.Fatalf("Error %s", err)
// 	}
//
// 	os.Args = []string{"./bin/kka-cli", "template"}
//     //go:embed ../../assets
//     var resources embed.FS
//
//     cli := CLI{
//         CWD: cwd,
//         Resources: resources,
//     }
//
//     cli.Run()
// }