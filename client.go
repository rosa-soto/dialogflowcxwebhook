package dialogflowcxwebhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Client struct {
	HTTPClient *http.Client
}

func NewClient() *Client {
	return &Client{
		HTTPClient: &http.Client{},
	}
}

const dialogFlowAI = "https://dialogflow.cloud.google.com/v1/cx/integrations/messenger/webhook/b3984f5d-5e30-42f2-9ec9-8fc5feb28606/sessions/1234"

type QueryInput struct {
	Text         Text   `json:"text"`
	LanguageCode string `json:"languageCode"`
}

type Text struct {
	Text string `json:"text"`
}

type MessageRequest struct {
	QueryInput QueryInput `json:"queryInput"`
}

type Responses struct {
	Text     string `json:"text"`
	Document string `json:"document"`
	URL      string `json:"url"`
	DebugID  string `json:"debugId"`
}
type QueryResponse struct {
	ResponseID  string `json:"responseId"`
	QueryResult struct {
		Text             string `json:"text"`
		LanguageCode     string `json:"languageCode"`
		ResponseMessages []struct {
			Text struct {
				Text                                  []string `json:"text"`
				RedactedText                          []string `json:"redactedText"`
				ContainAiGeneratedContent             bool     `json:"containAiGeneratedContent"`
				ContainDataStoreContent               bool     `json:"containDataStoreContent"`
				ContainAttemptedDataStoreInteractions bool     `json:"containAttemptedDataStoreInteractions"`
			} `json:"text,omitempty"`
			ResponseType string `json:"responseType"`
			Source       string `json:"source"`
			Payload      struct {
				RichContent [][]struct {
					Title      string `json:"title"`
					ActionLink string `json:"actionLink"`
					Type       string `json:"type"`
					Subtitle   string `json:"subtitle"`
				} `json:"richContent"`
			} `json:"payload,omitempty"`
		} `json:"responseMessages"`
		CurrentPage struct {
			Name        string `json:"name"`
			DisplayName string `json:"displayName"`
		} `json:"currentPage"`
		IntentDetectionConfidence float64 `json:"intentDetectionConfidence"`
		DiagnosticInfo            struct {
			ExecutionSequence []struct {
				Step1 struct {
					InitialState struct {
						FlowState struct {
							PageState struct {
								Name   string `json:"Name"`
								PageID string `json:"PageId"`
								Status string `json:"Status"`
							} `json:"PageState"`
							Version float64 `json:"Version"`
							FlowID  string  `json:"FlowId"`
							Name    string  `json:"Name"`
						} `json:"FlowState"`
						Event string `json:"Event"`
					} `json:"InitialState"`
					Type string `json:"Type"`
				} `json:"Step 1,omitempty"`
				Step2 struct {
					Type              string `json:"Type"`
					FunctionExecution struct {
						Responses []struct {
							Source string `json:"source"`
							Text   struct {
								ContainDataStoreContent               bool     `json:"containDataStoreContent"`
								Text                                  []string `json:"text"`
								RedactedText                          []string `json:"redactedText"`
								ContainAttemptedDataStoreInteractions bool     `json:"containAttemptedDataStoreInteractions"`
								ContainAiGeneratedContent             bool     `json:"containAiGeneratedContent"`
							} `json:"text,omitempty"`
							ResponseType string `json:"responseType"`
							Payload      struct {
								RichContent [][]struct {
									ActionLink string `json:"actionLink"`
									Subtitle   string `json:"subtitle"`
									Title      string `json:"title"`
									Type       string `json:"type"`
								} `json:"richContent"`
							} `json:"payload,omitempty"`
						} `json:"Responses"`
					} `json:"FunctionExecution"`
					StateMachine struct {
						FlowState struct {
							Name      string `json:"Name"`
							FlowID    string `json:"FlowId"`
							PageState struct {
								Name   string `json:"Name"`
								Status string `json:"Status"`
								PageID string `json:"PageId"`
							} `json:"PageState"`
							Version float64 `json:"Version"`
						} `json:"FlowState"`
						TriggeredTransitionRouteID string `json:"TriggeredTransitionRouteId"`
					} `json:"StateMachine"`
				} `json:"Step 2,omitempty"`
				Step3 struct {
					StateMachine struct {
						FlowState struct {
							Name      string  `json:"Name"`
							Version   float64 `json:"Version"`
							PageState struct {
								Name   string `json:"Name"`
								Status string `json:"Status"`
								PageID string `json:"PageId"`
							} `json:"PageState"`
							FlowID string `json:"FlowId"`
						} `json:"FlowState"`
					} `json:"StateMachine"`
					Type string `json:"Type"`
				} `json:"Step 3,omitempty"`
			} `json:"Execution Sequence"`
			DataStoreExecutionSequence struct {
				AdditionalInfo struct {
					UserQuery          string `json:"user_query"`
					UcsProjectNumber   string `json:"ucs_project_number"`
					RewrittenQuery     string `json:"rewritten_query"`
					AgentProjectNumber string `json:"agent_project_number"`
				} `json:"additionalInfo"`
				Steps []struct {
					Info      string      `json:"info,omitempty"`
					Responses []Responses `json:"responses"`
					Name      string      `json:"name"`
					Status    struct {
						Code string `json:"code"`
					} `json:"status"`
				} `json:"steps"`
			} `json:"DataStore Execution Sequence"`
			TriggeredTransitionNames []string `json:"Triggered Transition Names"`
			SessionID                string   `json:"Session Id"`
			ResponseID               string   `json:"Response Id"`
		} `json:"diagnosticInfo"`
		Match struct {
			MatchType                string  `json:"matchType"`
			Confidence               float64 `json:"confidence"`
			KnowledgeConnectorAnswer []struct {
				Answer               string  `json:"answer"`
				MatchConfidence      float64 `json:"matchConfidence"`
				GenerativePrediction struct {
					Snippets []struct {
						URI   string `json:"uri"`
						Text  string `json:"text"`
						Title string `json:"title"`
					} `json:"snippets"`
					GroundingInfo struct {
						Score    string `json:"score"`
						Accepted bool   `json:"accepted"`
					} `json:"groundingInfo"`
				} `json:"generativePrediction"`
			} `json:"knowledgeConnectorAnswer"`
		} `json:"match"`
		AdvancedSettings struct {
			LoggingSettings struct {
			} `json:"loggingSettings"`
		} `json:"advancedSettings"`
		CurrentFlow struct {
			Name        string `json:"name"`
			DisplayName string `json:"displayName"`
		} `json:"currentFlow"`
	} `json:"queryResult"`
	ResponseType string `json:"responseType"`
	UlmCalls     int    `json:"ulmCalls"`
}

func (s *Client) Post(questionText string) ([]Responses, error) {
	var resp *http.Response
	requestBody := map[string]interface{}{
		"queryInput": map[string]interface{}{
			"text": map[string]string{
				"text": "'" + questionText + "'",
			},
			"languageCode": "en",
		},
	}
	requestBodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("Error encoding request body: %s", err)
	}
	req, err := http.NewRequest(
		http.MethodPost,
		dialogFlowAI,
		bytes.NewBuffer(requestBodyJSON))

	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %s", err)
	}
	req.Header.Set("content-type", "application/json")
	resp, err = s.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do HTTP request: %s", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Received non-200 response %s", fmt.Sprint(resp.StatusCode))
	}
	defer resp.Body.Close()

	c, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	a := string(c[4:])
	x := []byte(a)
	var rs *QueryResponse
	err = json.Unmarshal(x, &rs)
	if err != nil {
		fmt.Printf("Error decoding %s", err)
	}

	return rs.QueryResult.DiagnosticInfo.DataStoreExecutionSequence.Steps[6].Responses, nil
}
