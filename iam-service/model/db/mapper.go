package db

/*import (
	"errors"
	"fmt"
	"iam-service/model"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"log"
)

/*func getResource(cypherResult interface{}) *domain.Resource {
	resource, err := domain.NewResource("", "")
	if err != nil {
		return nil
	}
	resource.Attributes = make([]domain.Attribute, 0)
	attrs := cypherResult.([]*neo4j.Record)[0].Values[1].([]interface{})
	for _, attr := range attrs {
		a := attr.(map[string]interface{})
		name := a["name"].(string)
		kind := domain.AttributeKind(a["kind"].(int64))
		value := a["value"]
		if name == "id" {
			resource.SetId(value.(string))
		}
		if name == "kind" {
			resource.SetKind(value.(string))
		}
		attrId, err := domain.NewAttributeId(name)
		if err != nil {
			return nil
		}
		attribute, err := domain.NewAttribute(*attrId, kind, value)
		if err != nil {
			return nil
		}
		resource.Attributes = append(resource.Attributes, *attribute)
	}
	return resource
}*/


