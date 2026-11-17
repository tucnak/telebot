package telebot

import "encoding/json"

// PassportData is Telegram Passport data a user shared with the bot. It is
// encrypted; see the Telegram Passport documentation for decrypting it.
//
// Bot API: PassportData.
type PassportData struct {
	Data        []EncryptedPassportElement `json:"data"`
	Credentials EncryptedCredentials       `json:"credentials"`
}

// EncryptedPassportElement is a document or other element of a user's
// Telegram Passport.
//
// Bot API: EncryptedPassportElement.
type EncryptedPassportElement struct {
	// "personal_details", "passport", "driver_license", "identity_card",
	// "internal_passport", "address", "utility_bill", "bank_statement",
	// "rental_agreement", "passport_registration",
	// "temporary_registration", "phone_number" or "email".
	Type string `json:"type"`

	// (Optional) Base64-encoded encrypted element data.
	Data string `json:"data,omitempty"`

	// (Optional) Verified phone number or email address.
	PhoneNumber string `json:"phone_number,omitempty"`
	Email       string `json:"email,omitempty"`

	// (Optional) Encrypted files of the element.
	Files       []PassportFile `json:"files,omitempty"`
	FrontSide   *PassportFile  `json:"front_side,omitempty"`
	ReverseSide *PassportFile  `json:"reverse_side,omitempty"`
	Selfie      *PassportFile  `json:"selfie,omitempty"`
	Translation []PassportFile `json:"translation,omitempty"`

	// Base64-encoded hash of the element.
	Hash string `json:"hash"`
}

// PassportFile is an encrypted JPEG file of a Telegram Passport element.
//
// Bot API: PassportFile.
type PassportFile struct {
	FileID       string `json:"file_id"`
	UniqueID     string `json:"file_unique_id"`
	FileSize     int64  `json:"file_size"`
	FileUnixtime int64  `json:"file_date"`
}

// EncryptedCredentials are the data needed to decrypt and authenticate the
// elements of PassportData.
//
// Bot API: EncryptedCredentials.
type EncryptedCredentials struct {
	Data   string `json:"data"`
	Hash   string `json:"hash"`
	Secret string `json:"secret"`
}

// Sources of PassportElementError.
const (
	PassportErrorDataField        = "data"
	PassportErrorFrontSide        = "front_side"
	PassportErrorReverseSide      = "reverse_side"
	PassportErrorSelfie           = "selfie"
	PassportErrorFile             = "file"
	PassportErrorFiles            = "files"
	PassportErrorTranslationFile  = "translation_file"
	PassportErrorTranslationFiles = "translation_files"
	PassportErrorUnspecified      = "unspecified"
)

// PassportElementError is an error in a Telegram Passport element the user
// must fix. Source says where the error is and which hash field applies.
//
// Bot API: PassportElementError, PassportElementErrorDataField, PassportElementErrorFrontSide, PassportElementErrorReverseSide, PassportElementErrorSelfie, PassportElementErrorFile, PassportElementErrorFiles, PassportElementErrorTranslationFile, PassportElementErrorTranslationFiles, PassportElementErrorUnspecified.
type PassportElementError struct {
	// One of the PassportError* sources.
	Source string `json:"source"`

	// The element type with the error, as in EncryptedPassportElement.
	Type string `json:"type"`

	// For "data", the field name and the hash of the data.
	FieldName string `json:"field_name,omitempty"`
	DataHash  string `json:"data_hash,omitempty"`

	// For a single file, its hash; for a list of files, their hashes.
	FileHash   string   `json:"file_hash,omitempty"`
	FileHashes []string `json:"file_hashes,omitempty"`

	// For "unspecified", the hash of the element.
	ElementHash string `json:"element_hash,omitempty"`

	Message string `json:"message"`
}

// SetPassportDataErrors tells user that some of the Telegram Passport
// elements they shared have errors; they can't share them again until the
// errors are fixed.
//
// Bot API: setPassportDataErrors.
func (b *Bot) SetPassportDataErrors(user Recipient, errs []PassportElementError) error {
	data, err := json.Marshal(errs)
	if err != nil {
		return err
	}
	_, err = b.Raw("setPassportDataErrors", map[string]string{
		"user_id": user.Recipient(),
		"errors":  string(data),
	})
	return err
}
