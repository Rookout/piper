package git_provider

import (
	"fmt"
	bitbucket "github.com/ktrysmt/go-bitbucket"
	"github.com/rookout/piper/pkg/conf"
	"github.com/rookout/piper/pkg/utils"
	"log"
	"net/http"
	"strings"
)

func ValidateBitbucketPermissions(client *bitbucket.Client, cfg *conf.GlobalConfig) error {

	repoAdminScopes := []string{"webhook", "repository:admin", "pullrequest:write"}
	repoGranularScopes := []string{"webhook", "repository", "pullrequest"}

	scopes, err := GetBitbucketTokenScopes(client, cfg)

	if err != nil {
		return fmt.Errorf("failed to get scopes: %v", err)
	}
	if len(scopes) == 0 {
		return fmt.Errorf("permissions error: no scopes found for the github client")
	}

	if utils.ListContains(repoAdminScopes, scopes) {
		return nil
	}
	if utils.ListContains(repoGranularScopes, scopes) {
		return nil
	}

	return fmt.Errorf("permissions error: %v is not a valid scopes", scopes)
}

func GetBitbucketTokenScopes(client *bitbucket.Client, cfg *conf.GlobalConfig) ([]string, error) {

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/repositories/%s", client.GetApiBaseURL(), cfg.GitProviderConfig.OrgName), nil)
	if err != nil {
		log.Println("Error creating request:", err)
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+cfg.GitProviderConfig.Token)

	resp, err := client.HttpClient.Do(req)
	if err != nil {
		log.Println("Error making request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("token validation failed: %v", resp.Status)
	}

	// Check the "X-OAuth-Scopes" header to get the token scopes
	acceptedScopes := resp.Header.Get("X-Accepted-OAuth-Scopes")
	scopes := resp.Header.Get("X-OAuth-Scopes")
	log.Println("Bitbucket Token Scopes are:", scopes, acceptedScopes)

	scopes = strings.ReplaceAll(scopes, " ", "")
	return append(strings.Split(scopes, ","), acceptedScopes), nil

}

func addHookToHashTable(hookUuid string, hookHashTable map[string]int64) {
	hookHashTable[hookUuid] = utils.StringToInt64(hookUuid)
}

func getHookByUUID(hookUuid string, hookHashTable map[string]int64) (int64, error) {
	res, ok := hookHashTable[hookUuid]
	if !ok {
		return 0, fmt.Errorf("hookUuid %s not found", hookUuid)
	}
	return res, nil
}
