package parsecmd

import (
    "fmt"
    "net/url"
    "encoding/json"

    "github.com/back4app/parse-cli/parsecli"
    "github.com/facebookgo/stackerr"
    "github.com/spf13/cobra"
)

type cloudCmd struct {
    function string
    data     string
}

func (s *cloudCmd) executeCloudFunction(e *parsecli.Env) (interface{}, error) {
    u := url.URL{
        Path: "functions/" + s.function,
    }
    var data interface{}
    json.Unmarshal([]byte(s.data), &data)

    var response interface{}
    if _, err := e.ParseAPIClient.Post(&u, data, &response); err != nil {
        return nil, stackerr.Wrap(err)
    }
    return response, nil
}

func (s *cloudCmd) run(e *parsecli.Env, c *parsecli.Context, args []string) error {
    if len(args) == 0 {
        return stackerr.New("Function argument is required.")
    }
    s.function = args[0]
    response, _ := s.executeCloudFunction(e)
    cloudResponse, _ := json.Marshal(&response)
    fmt.Println(string(cloudResponse))
    return nil
}

func NewCloudCmd(e *parsecli.Env) *cobra.Command {
    var s cloudCmd
    cmd := &cobra.Command{
        Use:   "cloud [app] <function>",
        Short: "Call cloud code function",
        Long: `...`,
        Run: parsecli.RunWithArgsClient(e, s.run),
    }
    cmd.Flags().StringVarP(&s.data, "data", "d", s.data,
        "Data")
    return cmd
}
