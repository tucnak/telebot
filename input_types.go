package telebot

// InputMessageContent objects represent the content of a message to be sent
// as a result of an inline query.
//
// Bot API: InputMessageContent.
type InputMessageContent interface {
	IsInputMessageContent() bool
}

// InputTextMessageContent represents the content of a text message to be
// sent as the result of an inline query.
//
// Bot API: InputTextMessageContent.
type InputTextMessageContent struct {
	// Text of the message to be sent, 1-4096 characters.
	Text string `json:"message_text"`

	// (Optional) Send Markdown or HTML, if you want Telegram apps to show
	// bold, italic, fixed-width text or inline URLs in your bot's message.
	ParseMode string `json:"parse_mode,omitempty"`

	// (Optional) Link preview generation options for the message.
	PreviewOptions *PreviewOptions `json:"link_preview_options,omitempty"`

	// Optional. Entities of the text, instead of ParseMode.
	Entities []MessageEntity `json:"entities,omitempty"`
}

func (input *InputTextMessageContent) IsInputMessageContent() bool {
	return true
}

// InputLocationMessageContent represents the content of a location message
// to be sent as the result of an inline query.
//
// Bot API: InputLocationMessageContent.
type InputLocationMessageContent struct {
	Lat float32 `json:"latitude"`
	Lng float32 `json:"longitude"`

	// Optional. Radius of uncertainty in meters (0-1500), and for live
	// locations the period in seconds, the direction in degrees and the
	// proximity alert radius in meters.
	HorizontalAccuracy float32 `json:"horizontal_accuracy,omitempty"`
	LivePeriod         int     `json:"live_period,omitempty"`
	Heading            int     `json:"heading,omitempty"`
	AlertRadius        int     `json:"proximity_alert_radius,omitempty"`
}

func (input *InputLocationMessageContent) IsInputMessageContent() bool {
	return true
}

// InputVenueMessageContent represents the content of a venue message to
// be sent as the result of an inline query.
//
// Bot API: InputVenueMessageContent.
type InputVenueMessageContent struct {
	Lat float32 `json:"latitude"`
	Lng float32 `json:"longitude"`

	// Name of the venue.
	Title string `json:"title"`

	// Address of the venue.
	Address string `json:"address"`

	// Optional. Foursquare identifier of the venue, if known.
	FoursquareID string `json:"foursquare_id,omitempty"`

	// Optional. Foursquare type, and Google Places identifier and type.
	FoursquareType  string `json:"foursquare_type,omitempty"`
	GooglePlaceID   string `json:"google_place_id,omitempty"`
	GooglePlaceType string `json:"google_place_type,omitempty"`
}

func (input *InputVenueMessageContent) IsInputMessageContent() bool {
	return true
}

// InputContactMessageContent represents the content of a contact
// message to be sent as the result of an inline query.
//
// Bot API: InputContactMessageContent.
type InputContactMessageContent struct {
	// Contact's phone number.
	PhoneNumber string `json:"phone_number"`

	// Contact's first name.
	FirstName string `json:"first_name"`

	// Optional. Contact's last name.
	LastName string `json:"last_name,omitempty"`

	// Optional. Additional data about the contact as a vCard, 0-2048 bytes.
	VCard string `json:"vcard,omitempty"`
}

func (input *InputContactMessageContent) IsInputMessageContent() bool {
	return true
}

// InputInvoiceMessageContent is an invoice to send as the result of an
// inline query.
//
// Bot API: InputInvoiceMessageContent.
type InputInvoiceMessageContent struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Payload     string  `json:"payload"`
	Currency    string  `json:"currency"`
	Prices      []Price `json:"prices"`

	// Optional. Payment provider token; empty for Telegram Stars.
	Token string `json:"provider_token,omitempty"`

	MaxTipAmount        int    `json:"max_tip_amount,omitempty"`
	SuggestedTipAmounts []int  `json:"suggested_tip_amounts,omitempty"`
	Data                string `json:"provider_data,omitempty"`

	PhotoURL    string `json:"photo_url,omitempty"`
	PhotoSize   int    `json:"photo_size,omitempty"`
	PhotoWidth  int    `json:"photo_width,omitempty"`
	PhotoHeight int    `json:"photo_height,omitempty"`

	NeedName            bool `json:"need_name,omitempty"`
	NeedPhoneNumber     bool `json:"need_phone_number,omitempty"`
	NeedEmail           bool `json:"need_email,omitempty"`
	NeedShippingAddress bool `json:"need_shipping_address,omitempty"`
	SendPhoneNumber     bool `json:"send_phone_number_to_provider,omitempty"`
	SendEmail           bool `json:"send_email_to_provider,omitempty"`
	Flexible            bool `json:"is_flexible,omitempty"`
}

// IsInputMessageContent marks InputInvoiceMessageContent as
// InputMessageContent.
func (input *InputInvoiceMessageContent) IsInputMessageContent() bool {
	return true
}
