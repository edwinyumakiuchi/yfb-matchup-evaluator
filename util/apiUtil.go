package util

import (
    "io"
    "fmt"
    "net/http"
)

func GetAPI(endpoint string, accessToken string) ([]byte, error) {
    apiReq, apiErr := http.NewRequest("GET", endpoint, nil)
    if apiErr != nil {
        return nil, fmt.Errorf("Error creating request: %v", apiErr)
    }
    apiReq.Header.Set("Authorization", "Bearer " + accessToken)

    client := &http.Client{}
    apiResp, apiErr := client.Do(apiReq)
    if apiErr != nil {
        return nil, fmt.Errorf("Error making request: %v", apiErr)
    }
    defer apiResp.Body.Close()

    if apiResp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("API responded with status code %d", apiResp.StatusCode)
    }

    apiBody, apiErr := io.ReadAll(apiResp.Body)
    if apiErr != nil {
        return nil, fmt.Errorf("Error reading response body: %v", apiErr)
    }

    return apiBody, nil
}