package exporter

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/jomei/notionapi"
)

type Exporter struct {
	client notionapi.DatabaseService
}

func NewExporter(token string) Exporter {
	return Exporter{
		client: notionapi.NewClient(notionapi.Token(token)).Database,
	}
}

func NewExporterWithClient(client *notionapi.Client) Exporter {
	return Exporter{
		client: client.Database,
	}
}

type Options struct {
	SortKey string
	Order   string
}

func (o Options) buildRequestParameter(cursor notionapi.Cursor) *notionapi.DatabaseQueryRequest {
	var sorts []notionapi.SortObject
	if o.SortKey != "" {
		sorts = append(sorts, notionapi.SortObject{
			Property:  o.SortKey,
			Direction: notionapi.SortOrder(o.Order),
		})
	} else {
		sorts = append(sorts, notionapi.SortObject{
			Timestamp: "created_time",
			Direction: notionapi.SortOrder(o.Order),
		})
	}
	return &notionapi.DatabaseQueryRequest{
		StartCursor: cursor,
		Sorts:       sorts,
	}
}

func (e *Exporter) ExportDatabase(ctx context.Context, databaseID string, options Options, writter io.Writer) error {

	w := csv.NewWriter(writter)
	database, err := e.client.Get(ctx, notionapi.DatabaseID(databaseID))
	if err != nil {
		fmt.Println("Failed to get database:", err)
		return err
	}

	header := make([]string, 0, len(database.Properties))
	for name, propertyConfig := range database.Properties {
		if EnableDownloadPropertyConfig(propertyConfig) {
			header = append(header, name)
		}
	}
	w.Write(header)

	var cursor notionapi.Cursor

	for {
		res, err := e.client.Query(ctx, notionapi.DatabaseID(databaseID), options.buildRequestParameter(cursor))
		if err != nil {
			return err
		}
		for _, row := range res.Results {
			values := make([]string, len(header))
			for i, key := range header {
				v, ok := row.Properties[key]
				if !ok {
					return err
				}
				value, err := GetStringValueByProperty(v)
				if err != nil {
					return err
				}
				values[i] = value
			}
			w.Write(values)
		}
		w.Flush()
		if !res.HasMore {
			return nil
		}
		cursor = res.NextCursor
	}
}

func GetStringValueByProperty(property notionapi.Property) (string, error) {
	switch property.GetType() {
	case notionapi.PropertyTypeTitle:
		titleProperty := property.(*notionapi.TitleProperty)
		return strings.Join(ExtractValues(titleProperty.Title, func(v notionapi.RichText) string {
			return v.PlainText
		}), ""), nil
	case notionapi.PropertyTypeRichText:
		richTextProperty := property.(*notionapi.RichTextProperty)
		return strings.Join(ExtractValues(richTextProperty.RichText, func(v notionapi.RichText) string {
			return v.PlainText
		}), ""), nil
	case notionapi.PropertyTypeText:
		textProperty := property.(*notionapi.TextProperty)
		return strings.Join(ExtractValues(textProperty.Text, func(v notionapi.RichText) string {
			return v.PlainText
		}), ""), nil
	case notionapi.PropertyTypeNumber:
		numberProperty := property.(*notionapi.NumberProperty)
		return fmt.Sprintf("%f", numberProperty.Number), nil
	case notionapi.PropertyTypeSelect:
		selectProperty := property.(*notionapi.SelectProperty)
		return selectProperty.Select.Name, nil
	case notionapi.PropertyTypeMultiSelect:
		multiSelectProperty := property.(*notionapi.MultiSelectProperty)
		return strings.Join(ExtractValues(multiSelectProperty.MultiSelect, func(v notionapi.Option) string {
			return v.Name
		}), ", "), nil
	case notionapi.PropertyTypeDate:
		dateProperty := property.(*notionapi.DateProperty)
		if dateProperty.Date.Start == nil {
			return "", nil
		}
		return dateProperty.Date.Start.String(), nil
	case notionapi.PropertyTypeFormula:
		formulaProperty := property.(*notionapi.FormulaProperty)
		switch formulaProperty.Formula.Type {
		case notionapi.FormulaTypeString:
			return formulaProperty.Formula.String, nil
		case notionapi.FormulaTypeNumber:
			return fmt.Sprintf("%f", formulaProperty.Formula.Number), nil
		case notionapi.FormulaTypeBoolean:
			if formulaProperty.Formula.Boolean {
				return "true", nil
			}
			return "false", nil
		case notionapi.FormulaTypeDate:
			if formulaProperty.Formula.Date == nil {
				return "", nil
			}
			return formulaProperty.Formula.Date.Start.String(), nil
		default:
			return "", fmt.Errorf("unsupported formula type: %s", formulaProperty.Formula.Type)
		}
	case notionapi.PropertyTypeRelation:
		relationProperty := property.(*notionapi.RelationProperty)
		return strings.Join(ExtractValues(relationProperty.Relation, func(v notionapi.Relation) string {
			return v.ID.String()
		}), ", "), nil
	case notionapi.PropertyTypeRollup:
		rollupProperty := property.(*notionapi.RollupProperty)
		switch rollupProperty.Rollup.Type {
		case notionapi.RollupTypeNumber:
			return fmt.Sprintf("%f", rollupProperty.Rollup.Number), nil
		case notionapi.RollupTypeDate:
			if rollupProperty.Rollup.Date == nil {
				return "", nil
			}
			return rollupProperty.Rollup.Date.Start.String(), nil
		default:
			return "", fmt.Errorf("unsupported rollup type: %s", rollupProperty.Rollup.Type)
		}
	case notionapi.PropertyTypePeople:
		peopleProperty := property.(*notionapi.PeopleProperty)
		return strings.Join(ExtractValues(peopleProperty.People, func(v notionapi.User) string {
			return v.Name
		}), ", "), nil
	case notionapi.PropertyTypeFiles:
		filesProperty := property.(*notionapi.FilesProperty)
		return strings.Join(ExtractValues(filesProperty.Files, func(v notionapi.File) string {
			return v.Name
		}), ", "), nil
	case notionapi.PropertyTypeCheckbox:
		checkboxProperty := property.(*notionapi.CheckboxProperty)
		if checkboxProperty.Checkbox {
			return "true", nil
		}
		return "false", nil
	case notionapi.PropertyTypeURL:
		urlProperty := property.(*notionapi.URLProperty)
		return urlProperty.URL, nil
	case notionapi.PropertyTypeEmail:
		emailProperty := property.(*notionapi.EmailProperty)
		return emailProperty.Email, nil
	case notionapi.PropertyTypePhoneNumber:
		phoneNumberProperty := property.(*notionapi.PhoneNumberProperty)
		return phoneNumberProperty.PhoneNumber, nil
	case notionapi.PropertyTypeCreatedTime:
		createdTimeProperty := property.(*notionapi.CreatedTimeProperty)
		return createdTimeProperty.CreatedTime.Format(time.RFC3339), nil
	case notionapi.PropertyTypeCreatedBy:
		createdByProperty := property.(*notionapi.CreatedByProperty)
		return createdByProperty.CreatedBy.Name, nil
	case notionapi.PropertyTypeLastEditedTime:
		lastEditedTimeProperty := property.(*notionapi.LastEditedTimeProperty)
		return lastEditedTimeProperty.LastEditedTime.Format(time.RFC3339), nil
	case notionapi.PropertyTypeLastEditedBy:
		lastEditedByProperty := property.(*notionapi.LastEditedByProperty)
		return lastEditedByProperty.LastEditedBy.Name, nil
	case notionapi.PropertyTypeStatus:
		statusProperty := property.(*notionapi.StatusProperty)
		return statusProperty.Status.Name, nil
	case notionapi.PropertyTypeUniqueID:
		uniqueIDProperty := property.(*notionapi.UniqueIDProperty)
		return uniqueIDProperty.UniqueID.String(), nil
	case notionapi.PropertyTypeVerification:
		verificationProperty := property.(*notionapi.VerificationProperty)
		return verificationProperty.Verification.State.String(), nil
	case notionapi.PropertyTypeButton:
		return "", fmt.Errorf("button property is not supported")
	default:
		return "", fmt.Errorf("unsupported property type: %s", property.GetType())
	}
}

func EnableDownloadPropertyConfig(propertyConfig notionapi.PropertyConfig) bool {
	switch propertyConfig.GetType() {
	case notionapi.PropertyConfigButton:
		return false
	default:
		return true
	}
}

func ExtractValues[T any](elms []T, extractFunc func(val T) string) []string {
	results := make([]string, len(elms))
	for i, v := range elms {
		results[i] = extractFunc(v)
	}
	return results
}
