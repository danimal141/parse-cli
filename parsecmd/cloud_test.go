package parsecmd

import (
    "io/ioutil"
    "net/http"
    "testing"

    "github.com/back4app/parse-cli/parsecli"
    "github.com/facebookgo/ensure"
    "github.com/facebookgo/parse"
)

func newCloudHarness(t testing.TB) *parsecli.Harness {
    h := parsecli.NewHarness(t)
    ht := parsecli.TransportFunc(func(r *http.Request) (*http.Response, error) {
        ensure.DeepEqual(t, r.URL.Path, "/1/functions/echo")
        return &http.Response{
            StatusCode: http.StatusOK,
            Body:       ioutil.NopCloser(r.Body),
        }, nil
    })
    h.Env.ParseAPIClient = &parsecli.ParseAPIClient{APIClient: &parse.Client{Transport: ht}}
    return h
}

func TestHelloWorld(t *testing.T) {
    t.Parallel()
    h := newCloudHarness(t)
    defer h.Stop()
    j := cloudCmd{
        function: "echo",
        data: `{"hello": "world"}`,
    }
    response, err := j.executeCloudFunction(h.Env)
    ensure.Nil(t, err)

    expected := make(map[string]interface{})
    expected["hello"] = "world"

    ensure.DeepEqual(t, response, expected)
}
