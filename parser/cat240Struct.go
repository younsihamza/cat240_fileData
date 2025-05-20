package cat240

type VideosHeaders struct {
    StartAzimuth  float64 		`json:"start_azimuth"`
    EndAzimuth    float64 		`json:"end_azimuth"`
    StartRange    int     		`json:"start_range"`
    CellDuration  float64 		`json:"cell_duration"`
}
type VideoCellsResolution struct {
	CompressionIndicator byte 	`json:"compression_indicator"`
	BitResolution        int 	`json:"bit_resolution"`
}
type VideoOctetsVideoCellCounters struct {
	ValidOctetsInVideoBlock int `json:"valid_octets_in_video_block"`
	ValidCellsInVideoBlock  int `json:"valid_cells_in_video_block"`
}
type DataSourceIndentifier struct {
	SAC int `json:"sac"`
	SIC int `json:"sic"`
}
// VideoDataItem is a struct that represents the data structure of the message
type VideoDataItem struct {
    DataSourceIndentifier        *DataSourceIndentifier  			`json:"data_source_identifier"`
    MessageType                  int             					`json:"message_type"`
    RecodeHeader                 int             					`json:"recode_header"`
    VideoSummary                 []byte          					`json:"video_summary"`
    VideoHeaderNano              *VideosHeaders  					`json:"video_header_nano"`
    VideoHeaderFemto             *VideosHeaders  					`json:"video_header_femto"`
    VideoCellsResolution         *VideoCellsResolution  			`json:"video_cells_resolution"`
    VideoOctetsVideoCellCounters *VideoOctetsVideoCellCounters  	`json:"video_octets_video_cell_counters"`
    VideoBlockLowDataVolume      []byte          					`json:"video_block_low_data_volume"`
    VideoBlockMediumDataVolume   []byte          					`json:"video_block_medium_data_volume"`
    VideoBlockHighDataVolume     []byte          					`json:"video_block_high_data_volume"`
    TimeOfDay                    float64         					`json:"time_of_day"`
    ReservedExpansionField       string          					`json:"reserved_expansion_field"`
    SpecialPurposeField          []byte          					`json:"special_purpose_field"`
}

type ValidData struct {
    DataSourceIndentifier        *DataSourceIndentifier           `json:"data_source_identifier"`
    MessageType                  int                              `json:"message_type"`
    RecodeHeader                 int                              `json:"recode_header"`
    VideoHeader                  *VideosHeaders                   `json:"video_header"`
    VideoCellsResolution         *VideoCellsResolution            `json:"video_cells_resolution"`
    VideoOctetsVideoCellCounters *VideoOctetsVideoCellCounters    `json:"video_octets_video_cell_counters"`
    VideoBlock                   []byte                           `json:"video_block"`
    TimeOfDay                    float64                          `json:"time_of_day"`
    ReservedExpansionField       string                           `json:"reserved_expansion_field"`
    SpecialPurposeField          []byte                           `json:"special_purpose_field"`
}


