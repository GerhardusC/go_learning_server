package datacollection

import (
	"testing"
)

/** 
Should accept default value as input with no error. */
func TestValidateTopicDefaultGiven(t *testing.T) {
	baseTopic := "/#"
	topic, err := validateTopic(baseTopic)

	if err != nil {
		t.Error("Error should not exist when using default topic.")
	}

	if topic != "/#" {
		t.Error("Topic should default to '/#' when dafaulted to.")
	}
}

/** 
Should accept input with leading and trailing slashes. */
func TestValidateTopicLeadingAndTrailingSlashes(t *testing.T) {
	baseTopic := "/as/da/sd/"
	topic, err := validateTopic(baseTopic)

	if err != nil {
		t.Error("Should not throw error while validating.", err)
	}

	if topic == "/#" {
		t.Error("Should not retrieve default value.")
	}

	if topic != "/as/da/sd/#" {
		t.Error("Topic should be unchanged. Topic is now: ", topic)
	}
}

/** 
Should accept input with leading slash. */
func TestValidateTopicLeadingSlash(t *testing.T) {
	baseTopic := "/as/da/sd"
	topic, err := validateTopic(baseTopic)

	if err != nil {
		t.Error("Should not throw error while validating.", err)
	}

	if topic == "/#" {
		t.Error("Should not retrieve default value.")
	}

	if topic != "/as/da/sd/#" {
		t.Error("Topic should be unchanged. Topic is now: ", topic)
	}
}

/** 
Should accept input with trailing slash. */
func TestValidateTopicTrailingSlash(t *testing.T) {
	baseTopic := "as/da/sd/"
	topic, err := validateTopic(baseTopic)

	if err != nil {
		t.Error("Should not throw error while validating.", err)
	}

	if topic == "/#" {
		t.Error("Should not retrieve default value.")
	}

	if topic != "/as/da/sd/#" {
		t.Error("Topic should be unchanged. Topic is now: ", topic)
	}
}

/** 
Should accept input with neither trailing, nor leading slash. */
func TestValidateTopicNeitherEndSlashes(t *testing.T) {
	baseTopic := "as/da/sd"
	topic, err := validateTopic(baseTopic)

	if err != nil {
		t.Error("Should not throw error while validating.", err)
	}

	if topic == "/#" {
		t.Error("Should not retrieve default value.")
	}

	if topic != "/as/da/sd/#" {
		t.Error("Topic should be unchanged. Topic is now: ", topic)
	}
}

/** 
Should throw error and default with successive slashes. */
func TestValidateTopicDoubleSlashesNotAllowed(t *testing.T) {
	baseTopic := "/as/a/sd//"
	topic, err := validateTopic(baseTopic)

	if err == nil {
		t.Error("Error should exist when invalid char exists.")
	}

	if topic != "/#" {
		t.Error("Topic should default to '/#' when invalid char exists.")
	}
}

/** 
Should throw error and default with invalid '#' character. */
func TestValidateTopicHashNotAllowed(t *testing.T) {
	baseTopic := "/as/a/sd#"
	topic, err := validateTopic(baseTopic)

	if err == nil {
		t.Error("Error should exist when invalid char exists.")
	}

	if topic != "/#" {
		t.Error("Topic should default to '/#' when invalid char exists.")
	}
}

/** 
Should throw error and default with invalid '*' character. */
func TestValidateTopicAsteriskNotAllowed(t *testing.T) {
	baseTopic := "/asd/as/*"
	topic, err := validateTopic(baseTopic)

	if err == nil {
		t.Error("Error should exist when asterisk included.")
	}

	if topic != "/#" {
		t.Error("Topic should default to '/#' when asterisk included.")
	}
}

