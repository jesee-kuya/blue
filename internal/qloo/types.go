package qloo

// Segment represents an audience segment with affinity score from Qloo's Taste AI™
type Segment struct {
    Name         string  `json:"name"`
    AffinityScore float64 `json:"affinity_score"`
}

// TasteProfileRequest represents the request structure for Qloo's Taste AI™ API
type TasteProfileRequest struct {
    Description string `json:"description"`
    Options     struct {
        MaxSegments int `json:"max_segments,omitempty"`
    } `json:"options,omitempty"`
}

// TasteProfileResponse represents the response structure from Qloo's Taste AI™ API
type TasteProfileResponse struct {
    Segments []Segment `json:"segments"`
    Status   string    `json:"status"`
    Message  string    `json:"message,omitempty"`
}